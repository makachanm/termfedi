package app

import (
	"fmt"
	"sync"
	"termfedi/component"
	"termfedi/config"
	"termfedi/layer"
	"termfedi/utils"

	"github.com/gdamore/tcell/v2"
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
}

func NewNotificationScreen(config config.Configuration) *NotificationScreen {
	ts := new(NotificationScreen)
	ts.Notifications = utils.NewItemAutoDemandPagination[layer.Notification](20, 5)
	ts.FetchLayer = layer.NewDataFetchAction(layer.NewMastodonFetch(config.Session.Token, config.Session.Url))

	items := ts.FetchLayer.GetNotifications()
	for _, item := range items {
		ts.Notifications.PutItem(item)
	}

	return ts
}

func (ns *NotificationScreen) InitScene(screen tcell.Screen, ctx ApplicationContext) {
	_, h := screen.Size()
	ns.Notifications.SetMaxItemPerPage(int(h / 5))
	screen.Clear()
}

func (ns *NotificationScreen) WindowChangedScene(screen tcell.Screen, ctx ApplicationContext) {
	_, h := screen.Size()
	ns.Notifications.SetMaxItemPerPage(int(h / 5))
}

func (ns *NotificationScreen) DoScene(screen tcell.Screen, event tcell.Event, ctx ApplicationContext) {
	ns.drawNotis(screen, ctx)

	switch ev := event.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyCtrlP:
			ctx.Exit(0)

		case tcell.KeyCtrlR:
			ns.refreshData(screen, ctx)

		case tcell.KeyCtrlN:
			ctx.TranslateTo("main")

		case tcell.KeyLeft:
			ns.Notifications.GoPrev()
			ctx.RequestFullRedraw()

		case tcell.KeyRight:
			ns.Notifications.GoNext()
			ctx.RequestFullRedraw()

		}
	}
}

func (ns *NotificationScreen) refreshData(screen tcell.Screen, ctx ApplicationContext) {
	textStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)

	screen.Clear()
	utils.WriteTo(screen, 0, 0, "Refreshing data... Please Wait.", textStyle)
	ns.Notifications.Clear()

	items := ns.FetchLayer.GetNotifications()
	for _, item := range items {
		ns.Notifications.PutItem(item)
	}
	screen.Clear()
	ns.drawNotis(screen, ctx)

}

func (ns *NotificationScreen) drawNotis(screen tcell.Screen, ctx ApplicationContext) {
	textStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)

	items := ns.Notifications.GetCurruntPage()
	for i, notes := range items {
		component.DrawNotificationComponent(0, i*5, notes, screen, textStyle, 5)
	}

	footer := fmt.Sprintf(" Noti Page %d/%d | [Quit] C-p | [Rfrh] C-r | [Timeline] C-n | [Prev] <- [Next] -> ", ns.Notifications.GetCurruntPagePointer()+1, ns.Notifications.GetTotalPage()+1)

	ctx.DrawFooterbar(footer)
}
