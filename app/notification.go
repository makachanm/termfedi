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
	"github.com/makachanm/flogger-lib"
)

/*
UNDER CONSTRUCTION
ONLY TESTING PERPOSE ONLY
*/

/*
 MUST REMOVE BEFORE RELEASE
*/

type NotificationScreen struct {
	//somethibgs brrrrr
	Notifications *utils.ItemAutoDemandPagination[layer.Notification]
	timelinelock  sync.RWMutex

	FetchLayer layer.DataFetch

	config config.Configuration
}

func NewNotificationScreen(cfg config.Configuration) *NotificationScreen {
	ts := new(NotificationScreen)
	flogger.Println("Initializing NotificationScreen")
	ts.Notifications = utils.NewItemAutoDemandPagination[layer.Notification](40, 5)

	var vlayer layer.FetchActionBase
	switch cfg.Session.Type {
	case config.Mastodon:
		flogger.Println("Using Mastodon layer")
		vlayer = layer.NewMastodonFetch(cfg.Session.Token, cfg.Session.Url)

	case config.Misskey:
		flogger.Println("Using Misskey layer")
		vlayer = layer.NewMisskeyFetch(cfg.Session.Token, cfg.Session.Url)

	default:
		panic("session type invalid")

	}
	ts.FetchLayer = layer.NewDataFetchAction(vlayer)

	items := ts.FetchLayer.GetNotifications()
	flogger.Printf("Fetched %d notifications", len(items))
	for _, item := range items {
		ts.Notifications.PutItem(item)
	}

	ts.config = cfg

	return ts
}

func (ns *NotificationScreen) InitScene(screen tcell.Screen, ctx ApplicationContext) {
	flogger.Println("NotificationScreen: InitScene")
	_, h := screen.Size()
	ns.Notifications.SetMaxItemPerPage(int(h / ns.config.UI.MaxItemHeight))

	autoRef := func() { ns.autoRefresh(screen, ctx) }
	time.AfterFunc(time.Second*30, autoRef)

	ns.drawNotis(screen, ctx)
}

func (ns *NotificationScreen) WindowChangedScene(screen tcell.Screen, ctx ApplicationContext) {
	flogger.Println("NotificationScreen: WindowChangedScene")
	_, h := screen.Size()
	ns.Notifications.SetMaxItemPerPage(int(h / ns.config.UI.MaxItemHeight))
}

func (ns *NotificationScreen) DoScene(screen tcell.Screen, event tcell.Event, ctx ApplicationContext) {
	ns.drawNotis(screen, ctx)

	switch ev := event.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyCtrlP:
			flogger.Println("NotificationScreen: Ctrl+P pressed, exiting")
			ctx.Exit(0)

		case tcell.KeyCtrlR:
			flogger.Println("NotificationScreen: Ctrl+R pressed, refreshing data")
			ns.refreshData()

		case tcell.KeyCtrlN:
			flogger.Println("NotificationScreen: Ctrl+N pressed, switching to main scene")
			ctx.TranslateTo("main")

		case tcell.KeyLeft:
			flogger.Println("NotificationScreen: Left arrow pressed, going to previous page")
			ns.Notifications.GoPrev()
			screen.Clear()
			ns.drawNotis(screen, ctx)
			//ctx.RequestFullRedraw()

		case tcell.KeyRight:
			flogger.Println("NotificationScreen: Right arrow pressed, going to next page")
			ns.Notifications.GoNext()
			screen.Clear()
			ns.drawNotis(screen, ctx)
			//ctx.RequestFullRedraw()

		}
	}
}

func (ns *NotificationScreen) refreshData() {
	flogger.Println("NotificationScreen: refreshing data")
	ns.Notifications.Clear()

	currunt_pos := ns.Notifications.GetCurruntPagePointer()

	items := ns.FetchLayer.GetNotifications()
	flogger.Printf("NotificationScreen: fetched %d new notifications", len(items))
	for _, item := range items {
		ns.Notifications.PutItem(item)
	}

	for i := 0; i < currunt_pos; i++ {
		ns.Notifications.GoNext()
	}

}

func (ns *NotificationScreen) autoRefresh(screen tcell.Screen, ctx ApplicationContext) {
	flogger.Println("NotificationScreen: auto-refreshing data")
	ns.refreshData()

	if ctx.GetCurruntScene() == "noti" {
		screen.Clear()
		ns.drawNotis(screen, ctx)
		screen.Show()
	}

	autoRef := func() { ns.autoRefresh(screen, ctx) }
	time.AfterFunc(time.Second*30, autoRef)
}

func (ns *NotificationScreen) drawNotis(screen tcell.Screen, ctx ApplicationContext) {
	textStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)

	items := ns.Notifications.GetCurruntPage()
	for i, notes := range items {
		component.DrawNotificationComponent(0, i*ns.config.UI.MaxItemHeight, notes, screen, textStyle, ns.config.UI.MaxItemHeight)
	}

	footer := fmt.Sprintf(" Noti Page %d/%d | [Quit] C-p | [Refresh] C-r | [Back] C-n | [Prev] <- [Next] -> ", ns.Notifications.GetCurruntPagePointer()+1, ns.Notifications.GetTotalPage()+1)

	ctx.DrawFooterbar(footer)
}