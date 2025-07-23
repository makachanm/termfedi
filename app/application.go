package app

import (
	"os"
	utils "termfedi/utils"

	"github.com/gdamore/tcell/v2"
)

type TerminalMainApp struct {
	screen        tcell.Screen
	scenes        map[string]ApplicationScene
	currunt_scene string
	appctx        *MainAppContexts

	termination_signal   chan int
	transision_signal    chan string
	redraw_signal        chan bool
	footer_redraw_signal chan bool
}

func NewTerminalScreen() *TerminalMainApp {
	app := new(TerminalMainApp)

	app.currunt_scene = "main"
	app.scenes = make(map[string]ApplicationScene)
	app.termination_signal = make(chan int, 1)
	app.transision_signal = make(chan string, 1)
	app.redraw_signal = make(chan bool, 1)
	app.footer_redraw_signal = make(chan bool, 1)

	app.appctx = NewMainAppCtx(app.termination_signal, app.transision_signal, app.redraw_signal, app.footer_redraw_signal)
	return app
}

func (t *TerminalMainApp) InitTerminalScreen() error {
	var e error
	t.screen, e = tcell.NewScreen()

	if e != nil {
		return e
	}

	e = t.screen.Init()
	if e != nil {
		return e
	}

	// TODO: make color to customizeable
	color := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)
	t.screen.SetStyle(color)

	return nil
}

func (t *TerminalMainApp) RegisterScene(name string, scene ApplicationScene) {
	t.scenes[name] = scene
}

// TODO: add global event handling
func (t *TerminalMainApp) DrawStatusBar() {
	w, h := t.screen.Size()
	textStyle := tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite)

	text := t.appctx.GetFooterbar()
	nullspace := w - len(text)

	for i := 0; i < nullspace; i++ {
		text += " "
	}

	utils.WriteTo(t.screen, 0, h-1, text, textStyle)
}

func (t *TerminalMainApp) Mainloop() {
	t.scenes[t.currunt_scene].InitScene(t.screen, t.appctx)

	for {
		need_full_redraw := false

		select {
		case <-t.termination_signal:
			close(t.termination_signal)
			close(t.transision_signal)
			close(t.redraw_signal)
			t.screen.Fini()
			os.Exit(0)

		case target := <-t.transision_signal:
			if _, ok := t.scenes[target]; ok {
				t.currunt_scene = target
				need_full_redraw = true
			}
			t.scenes[t.currunt_scene].InitScene(t.screen, t.appctx)

		case <-t.redraw_signal:
			need_full_redraw = true

		case <-t.footer_redraw_signal:
			t.DrawStatusBar()

		default:
		}

		event := t.screen.PollEvent()

		switch event.(type) {
		case *tcell.EventResize:
			need_full_redraw = true
			t.scenes[t.currunt_scene].WindowChangedScene(t.screen, t.appctx)
		}

		if need_full_redraw {
			t.screen.Clear()
			t.screen.Sync()
		}

		t.scenes[t.currunt_scene].DoScene(t.screen, event, t.appctx)
		t.DrawStatusBar()

		t.screen.Show()
	}
}

type MainAppContexts struct {
	termination_signal chan int
	transision_signal  chan string
	redraw_signal      chan bool
	footer_redraw      chan bool

	Footer string
}

func NewMainAppCtx(termsig chan int, transsig chan string, redrawsig chan bool, footer chan bool) *MainAppContexts {
	Ctx := new(MainAppContexts)

	Ctx.termination_signal = termsig
	Ctx.transision_signal = transsig
	Ctx.redraw_signal = redrawsig
	Ctx.footer_redraw = footer

	return Ctx
}

func (m *MainAppContexts) Exit(exitcode int) {
	m.termination_signal <- exitcode
}

func (m *MainAppContexts) TranslateTo(name string) {
	m.transision_signal <- name
}

func (m *MainAppContexts) RequestFullRedraw() {
	m.redraw_signal <- true
}

func (m *MainAppContexts) RequestFooterbarRedraw() {
	m.footer_redraw <- true
}

func (m *MainAppContexts) DrawFooterbar(content string) {
	m.Footer = content
}

func (m *MainAppContexts) GetFooterbar() string {
	return m.Footer
}
