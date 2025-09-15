package app

import (
	"fmt"
	"termfedi/component"
	"termfedi/config"
	"termfedi/layer"
	utils "termfedi/utils"

	"github.com/gdamore/tcell/v2"
	"github.com/makachanm/flogger-lib"
)

var ActionTimelineNoteData []layer.Note

func init() {
	ActionTimelineNoteData = make([]layer.Note, 0)
}

func insertActionTargets(notes []layer.Note) {
	flogger.Printf("Inserting %d notes as action targets", len(notes))
	ActionTimelineNoteData = append(ActionTimelineNoteData, notes...)
}

func clearActionTargets() {
	flogger.Println("Clearing action targets")
	ActionTimelineNoteData = make([]layer.Note, 0)
}

type ActionType int

const (
	RENOTE ActionType = iota
	REACT
)

type ActionScreen struct {
	FetchLayer layer.DataFetch

	config      config.Configuration
	action_mode ActionType
	selection   int
}

func NewActionScreen(cfg config.Configuration) *ActionScreen {
	as := new(ActionScreen)
	flogger.Println("Initializing ActionScreen")

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
	as.FetchLayer = layer.NewDataFetchAction(vlayer)

	as.action_mode = REACT
	as.config = cfg
	as.selection = 0

	return as
}

func (as *ActionScreen) InitScene(screen tcell.Screen, ctx ApplicationContext) {
	flogger.Println("ActionScreen: InitScene")
	as.drawNotes(screen, ctx)
}

func (as *ActionScreen) WindowChangedScene(screen tcell.Screen, ctx ApplicationContext) {
	flogger.Println("ActionScreen: WindowChangedScene")
}

func (as *ActionScreen) DoScene(screen tcell.Screen, event tcell.Event, ctx ApplicationContext) {
	as.drawNotes(screen, ctx)

	switch ev := event.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyCtrlX:
			flogger.Println("ActionScreen: Ctrl+X pressed, returning to main scene")
			clearActionTargets()
			ctx.TranslateTo("main")

		case tcell.KeyLeft:
			if as.selection > 0 {
				as.selection--
				flogger.Printf("ActionScreen: Left arrow pressed, selection is now %d", as.selection)
			}

		case tcell.KeyRight:
			if as.selection < len(ActionTimelineNoteData)-1 {
				as.selection++
				flogger.Printf("ActionScreen: Right arrow pressed, selection is now %d", as.selection)
			}

		case tcell.KeyF1:
			flogger.Println("ActionScreen: F1 pressed, action mode set to REACT")
			as.action_mode = REACT

		case tcell.KeyF2:
			flogger.Println("ActionScreen: F2 pressed, action mode set to RENOTE")
			as.action_mode = RENOTE

		case tcell.KeyEnter:
			flogger.Printf("ActionScreen: Enter pressed, performing action %v on note %s", as.action_mode, ActionTimelineNoteData[as.selection].Id)
			switch as.action_mode {
			case REACT:
				as.FetchLayer.PostReaction(ActionTimelineNoteData[as.selection].Id)
				clearActionTargets()
				ctx.TranslateTo("main")
			case RENOTE:
				as.FetchLayer.PostRenote(ActionTimelineNoteData[as.selection].Id)
				clearActionTargets()
				ctx.TranslateTo("main")
			}
		}
	}

}

func (as *ActionScreen) drawNotes(screen tcell.Screen, ctx ApplicationContext) {
	screen.Clear()

	textStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)

	for i, notes := range ActionTimelineNoteData {
		if i == as.selection {
			utils.WriteTo(screen, 0, i*as.config.UI.MaxItemHeight, " > ", textStyle)
		}
		component.DrawNoteComponent(3, i*as.config.UI.MaxItemHeight, notes, screen, textStyle, as.config.UI.MaxItemHeight, true)
	}

	var selectedmode string
	switch as.action_mode {
	case REACT:
		selectedmode = "React"
	case RENOTE:
		selectedmode = "Renote"

	}

	footer := fmt.Sprintf(" Action Mode: %s | [Back] C-x | [Enter] Do | [Prev] <- [Next] -> | [F1] React | [F2] Renote", selectedmode)

	ctx.DrawFooterbar(footer)
	screen.Show()
}
