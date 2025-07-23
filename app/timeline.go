package app

import (
	"fmt"
	"sync"
	"termfedi/component"
	"termfedi/config"
	"termfedi/layer"
	utils "termfedi/utils"

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

const (
	CURRUNTTL_GLOBAL int = iota + 1
	CURRUNTTL_HOME
	CURRUNTTL_LOCAL
)

type TimelineScreen struct {
	//somethibgs brrrrr
	Timelines    *utils.ItemAutoDemandPagination[layer.Note]
	timelinelock sync.RWMutex

	FetchLayer layer.DataFetch
}

func NewTimelineScreen(config config.Configuration) *TimelineScreen {
	ts := new(TimelineScreen)
	ts.Timelines = utils.NewItemAutoDemandPagination[layer.Note](20, 5)
	ts.FetchLayer = layer.NewDataFetchAction(layer.NewMastodonFetch(config.Session.Token, config.Session.Url))

	items := ts.FetchLayer.GetHomeTimeline()
	for _, item := range items {
		ts.Timelines.PutItem(item)
	}

	return ts
}

func (ts *TimelineScreen) InitScene(screen tcell.Screen, ctx ApplicationContext) {
	_, h := screen.Size()
	ts.Timelines.SetMaxItemPerPage(int(h / 6))
}

func (ts *TimelineScreen) WindowChangedScene(screen tcell.Screen, ctx ApplicationContext) {
	_, h := screen.Size()
	ts.Timelines.SetMaxItemPerPage(int(h / 6))
}

func (ts *TimelineScreen) DoScene(screen tcell.Screen, event tcell.Event, ctx ApplicationContext) {
	textStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)

	items := ts.Timelines.GetCurruntPage()
	for i, notes := range items {
		component.DrawNoteComponent(0, i*6, notes, screen, textStyle, 6)
	}

	footer := fmt.Sprintf(" HomeTL Page %d/%d | [Quit] C-p | [Rfrh] C-r | [Noti] C-n [Comp] C-q | [Prev] <- [Next] -> ", ts.Timelines.GetCurruntPagePointer()+1, ts.Timelines.GetTotalPage()+1)
	ctx.DrawFooterbar(footer)

	switch ev := event.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyCtrlP:
			ctx.Exit(0)

		case tcell.KeyCtrlR:
			screen.Clear()
			utils.WriteTo(screen, 0, 0, "Refreshing data... Please Wait.", textStyle)
			ts.Timelines.Clear()

			items := ts.FetchLayer.GetHomeTimeline()
			for _, item := range items {
				ts.Timelines.PutItem(item)
			}
			ctx.RequestFullRedraw()

		case tcell.KeyLeft:
			ts.Timelines.GoPrev()
			ctx.RequestFullRedraw()

		case tcell.KeyRight:
			ts.Timelines.GoNext()
			ctx.RequestFullRedraw()

		}
	}
}
