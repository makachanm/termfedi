package app

import (
	"fmt"
	"termfedi/component"
	"termfedi/config"
	"termfedi/layer"
	utils "termfedi/utils"

	"github.com/gdamore/tcell/v2"
)

var ActionTimelineNoteData []layer.Note

func init() {
	ActionTimelineNoteData = make([]layer.Note, 0)
}

func insertActionTargets(notes []layer.Note) {
	ActionTimelineNoteData = append(ActionTimelineNoteData, notes...)
}

func clearActionTargets() {
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

	var vlayer layer.FetchActionBase
	switch cfg.Session.Type {
	case config.Mastodon:
		vlayer = layer.NewMastodonFetch(cfg.Session.Token, cfg.Session.Url)

	case config.Misskey:
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
	as.drawNotes(screen, ctx)
}

func (as *ActionScreen) WindowChangedScene(screen tcell.Screen, ctx ApplicationContext) {
}

func (as *ActionScreen) DoScene(screen tcell.Screen, event tcell.Event, ctx ApplicationContext) {
	as.drawNotes(screen, ctx)

	switch ev := event.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyCtrlX:
			ctx.TranslateTo("main")

		case tcell.KeyLeft:
			if as.selection > 0 {
				as.selection--
			}

		case tcell.KeyRight:
			if as.selection < len(ActionTimelineNoteData)-1 {
				as.selection++
			}

		case tcell.KeyF1:
			as.action_mode = REACT

		case tcell.KeyF2:
			as.action_mode = RENOTE

		case tcell.KeyEnter:
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

	footer := fmt.Sprintf("Action Mode: %s | [Back] C-x | [Enter] Do | [Prev] <- [Next] -> | [F1] React | [F2] Renote", selectedmode)

	ctx.DrawFooterbar(footer)
	screen.Show()
}
