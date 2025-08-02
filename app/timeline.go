package app

import (
	"fmt"
	"sync"
	"termfedi/component"
	"termfedi/config"
	"termfedi/layer"
	"termfedi/utils"
	"time"

	"github.com/gdamore/tcell/v2"
)

/*
UNDER CONSTRUCTION
ONLY TESTING PERPOSE ONLY
*/
const INSTANCE = ""
const TOKEN = ""

/*
 MUST REMOVE BEFORE RELEASE
*/

type CurruntTL int

const (
	CURRUNTTL_GLOBAL CurruntTL = iota + 1
	CURRUNTTL_HOME
	CURRUNTTL_LOCAL
)

func curruntTimeline(tl CurruntTL) string {
	switch tl {
	case CURRUNTTL_GLOBAL:
		return "Glob"
	case CURRUNTTL_HOME:
		return "Home"
	case CURRUNTTL_LOCAL:
		return "Locl"
	default:
		return "Unknown"
	}
}

func getTimeline(layer *layer.DataFetch, tl CurruntTL) []layer.Note {
	switch tl {
	case CURRUNTTL_GLOBAL:
		return layer.GetGlobalTimeline()
	case CURRUNTTL_HOME:
		return layer.GetHomeTimeline()
	case CURRUNTTL_LOCAL:
		return layer.GetLocalTimeline()
	default:
		return nil
	}
}

type TimelineScreen struct {
	//somethibgs brrrrr
	Timelines    *utils.ItemAutoDemandPagination[layer.Note]
	timelinelock sync.RWMutex

	FetchLayer layer.DataFetch
	currunt_tl CurruntTL

	config config.Configuration
	showcw bool
}

func NewTimelineScreen(cfg config.Configuration) *TimelineScreen {
	ts := new(TimelineScreen)

	var vlayer layer.FetchActionBase
	switch cfg.Session.Type {
	case config.Mastodon:
		vlayer = layer.NewMastodonFetch(cfg.Session.Token, cfg.Session.Url)
		ts.Timelines = utils.NewItemAutoDemandPagination[layer.Note](20, 5)

	case config.Misskey:
		vlayer = layer.NewMisskeyFetch(cfg.Session.Token, cfg.Session.Url)
		ts.Timelines = utils.NewItemAutoDemandPagination[layer.Note](100, 5)

	default:
		panic("session type invalid")

	}
	ts.FetchLayer = layer.NewDataFetchAction(vlayer)

	ts.currunt_tl = CURRUNTTL_HOME
	items := getTimeline(&ts.FetchLayer, ts.currunt_tl)
	for _, item := range items {
		ts.Timelines.PutItem(item)
	}

	ts.config = cfg
	ts.showcw = false

	return ts
}

func (ts *TimelineScreen) InitScene(screen tcell.Screen, ctx ApplicationContext) {
	_, h := screen.Size()
	ts.Timelines.SetMaxItemPerPage(int(h / ts.config.UI.MaxItemHeight))

	autoRef := func() {
		ts.autoRefresh(screen, ctx)
	}
	time.AfterFunc(time.Second*20, autoRef)

	ts.drawNotes(screen, ctx)
}

func (ts *TimelineScreen) WindowChangedScene(screen tcell.Screen, ctx ApplicationContext) {
	_, h := screen.Size()
	ts.Timelines.SetMaxItemPerPage(int(h / ts.config.UI.MaxItemHeight))
}

func (ts *TimelineScreen) DoScene(screen tcell.Screen, event tcell.Event, ctx ApplicationContext) {

	ts.drawNotes(screen, ctx)

	switch ev := event.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyCtrlP:
			ctx.Exit(0)

		case tcell.KeyCtrlR:
			ts.refreshData()

		case tcell.KeyCtrlE:
			switch ts.currunt_tl {
			case CURRUNTTL_GLOBAL:
				ts.currunt_tl = CURRUNTTL_LOCAL
			case CURRUNTTL_HOME:
				ts.currunt_tl = CURRUNTTL_GLOBAL
			case CURRUNTTL_LOCAL:
				ts.currunt_tl = CURRUNTTL_HOME
			}

			ts.refreshData()

		case tcell.KeyCtrlQ:
			if ts.showcw {
				ts.showcw = false
			} else {
				ts.showcw = true
			}
			screen.Clear()
			ts.drawNotes(screen, ctx)

		case tcell.KeyCtrlN:
			ctx.TranslateTo("noti")

		case tcell.KeyLeft:
			ts.Timelines.GoPrev()
			screen.Clear()
			ts.drawNotes(screen, ctx)

		case tcell.KeyRight:
			ts.Timelines.GoNext()
			screen.Clear()
			ts.drawNotes(screen, ctx)

		}
	}
}

func (ts *TimelineScreen) refreshData() {
	ts.timelinelock.Lock()
	currunt_pos := ts.Timelines.GetCurruntPagePointer()

	ts.Timelines.Clear()

	items := getTimeline(&ts.FetchLayer, ts.currunt_tl)
	for _, item := range items {
		ts.Timelines.PutItem(item)
	}

	for i := 0; i < currunt_pos; i++ {
		ts.Timelines.GoNext()
	}

	ts.timelinelock.Unlock()
}

func (ts *TimelineScreen) autoRefresh(screen tcell.Screen, ctx ApplicationContext) {
	ts.refreshData()

	if ctx.GetCurruntScene() == "main" {
		screen.Clear()
		ts.drawNotes(screen, ctx)
		screen.Show()
	}

	autoRef := func() { ts.autoRefresh(screen, ctx) }
	time.AfterFunc(time.Second*20, autoRef)
}

func (ts *TimelineScreen) drawNotes(screen tcell.Screen, ctx ApplicationContext) {
	ts.timelinelock.Lock()
	textStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)

	items := ts.Timelines.GetCurruntPage()
	for i, notes := range items {
		component.DrawNoteComponent(0, i*ts.config.UI.MaxItemHeight, notes, screen, textStyle, ts.config.UI.MaxItemHeight, ts.showcw)
	}

	footer := fmt.Sprintf(" %s Page %d/%d | [Quit] C-p | [Rfrh] C-r | [TL] C-e [Noti] C-n [CW] C-q | [Prev] <- [Next] -> ", curruntTimeline(ts.currunt_tl), ts.Timelines.GetCurruntPagePointer()+1, ts.Timelines.GetTotalPage()+1)

	ctx.DrawFooterbar(footer)
	ts.timelinelock.Unlock()
}
