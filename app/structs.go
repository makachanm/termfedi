package app

import "github.com/gdamore/tcell/v2"

type ApplicationScene interface {
	InitScene(screen tcell.Screen, ctx ApplicationContext)
	WindowChangedScene(screen tcell.Screen, ctx ApplicationContext)
	DoScene(screen tcell.Screen, event tcell.Event, ctx ApplicationContext)
}

type ApplicationContext interface {
	Exit(exitcode int)
	TranslateTo(name string)
	RequestFullRedraw()
	RequestFooterbarRedraw()
	DrawFooterbar(content string)
}

type Configuration struct {
}
