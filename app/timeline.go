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
		return "GlobTL"
	case CURRUNTTL_HOME:
		return "HomeTL"
	case CURRUNTTL_LOCAL:
		return "LoclTL"
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
}

func NewTimelineScreen(config config.Configuration) *TimelineScreen {
	ts := new(TimelineScreen)
	ts.Timelines = utils.NewItemAutoDemandPagination[layer.Note](20, 5)
	ts.FetchLayer = layer.NewDataFetchAction(layer.NewMastodonFetch(config.Session.Token, config.Session.Url))

	ts.currunt_tl = CURRUNTTL_HOME
	items := getTimeline(&ts.FetchLayer, ts.currunt_tl)
	for _, item := range items {
		ts.Timelines.PutItem(item)
	}

	return ts
}

func (ts *TimelineScreen) InitScene(screen tcell.Screen, ctx ApplicationContext) {
	_, h := screen.Size()
	ts.Timelines.SetMaxItemPerPage(int(h / 6))

	autoRef := func() { ts.autoRefresh(ctx) }
	time.AfterFunc(time.Second*30, autoRef)
}

func (ts *TimelineScreen) WindowChangedScene(screen tcell.Screen, ctx ApplicationContext) {
	_, h := screen.Size()
	ts.Timelines.SetMaxItemPerPage(int(h / 6))
}

func (ts *TimelineScreen) DoScene(screen tcell.Screen, event tcell.Event, ctx ApplicationContext) {
	ts.drawNotes(screen, ctx)

	switch ev := event.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyCtrlP:
			ctx.Exit(0)

		case tcell.KeyCtrlR:
			ts.refreshData(screen, ctx)

		case tcell.KeyCtrlE:
			switch ts.currunt_tl {
			case CURRUNTTL_GLOBAL:
				ts.currunt_tl = CURRUNTTL_LOCAL
			case CURRUNTTL_HOME:
				ts.currunt_tl = CURRUNTTL_GLOBAL
			case CURRUNTTL_LOCAL:
				ts.currunt_tl = CURRUNTTL_HOME
			}

			ts.refreshData(screen, ctx)

		case tcell.KeyCtrlN:
			ctx.TranslateTo("noti")

		case tcell.KeyLeft:
			ts.Timelines.GoPrev()
			ctx.RequestFullRedraw()

		case tcell.KeyRight:
			ts.Timelines.GoNext()
			ctx.RequestFullRedraw()

		}
	}
}

func (ts *TimelineScreen) refreshData(screen tcell.Screen, ctx ApplicationContext) {
	textStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)

	screen.Clear()
	utils.WriteTo(screen, 0, 0, "Refreshing data... Please Wait.", textStyle)

	ts.Timelines.Clear()

	items := getTimeline(&ts.FetchLayer, ts.currunt_tl)
	for _, item := range items {
		ts.Timelines.PutItem(item)
	}

	screen.Clear()
	ts.drawNotes(screen, ctx)

}

func (ts *TimelineScreen) autoRefresh(ctx ApplicationContext) {
	ts.Timelines.Clear()
	items := getTimeline(&ts.FetchLayer, ts.currunt_tl)
	for _, item := range items {
		ts.Timelines.PutItem(item)
	}

	autoRef := func() { ts.autoRefresh(ctx) }
	time.AfterFunc(time.Second*30, autoRef)
}

func (ts *TimelineScreen) drawNotes(screen tcell.Screen, ctx ApplicationContext) {
	textStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)

	items := ts.Timelines.GetCurruntPage()
	for i, notes := range items {
		component.DrawNoteComponent(0, i*6, notes, screen, textStyle, 6)
	}

	footer := fmt.Sprintf(" %s Page %d/%d | [Quit] C-p | [Rfrh] C-r | [TL] C-e [Noti] C-n [Comp] C-q | [Prev] <- [Next] -> ", curruntTimeline(ts.currunt_tl), ts.Timelines.GetCurruntPagePointer()+1, ts.Timelines.GetTotalPage()+1)
	ctx.DrawFooterbar(footer)
}
