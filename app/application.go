package app

import (
	"os"
	utils "termfedi/utils"

	"github.com/gdamore/tcell/v2"
	"github.com/makachanm/flogger-lib"
)

type TerminalMainApp struct {
	screen        tcell.Screen
	scenes        map[string]ApplicationScene
	currunt_scene string
	appctx        *MainAppContexts

	termination_signal chan int
	transision_signal  chan string
}

func NewTerminalScreen() *TerminalMainApp {
	app := new(TerminalMainApp)

	flogger.Println("Initializing TerminalMainApp")
	app.currunt_scene = "main"
	app.scenes = make(map[string]ApplicationScene)
	app.termination_signal = make(chan int, 1)
	app.transision_signal = make(chan string, 1)

	app.appctx = NewMainAppCtx(app.termination_signal, app.transision_signal, app.DrawStatusBar, &app.currunt_scene)
	return app
}

func (t *TerminalMainApp) InitTerminalScreen() error {
	var e error
	flogger.Println("Initializing tcell screen")
	t.screen, e = tcell.NewScreen()

	if e != nil {
		return flogger.Errorf("Failed to create new screen: %v", e)
	}

	e = t.screen.Init()
	if e != nil {
		return flogger.Errorf("Failed to init screen: %v", e)
	}

	// TODO: make color to customizeable
	color := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)
	t.screen.SetStyle(color)

	return nil
}

func (t *TerminalMainApp) RegisterScene(name string, scene ApplicationScene) {
	flogger.Printf("Registering scene: %s", name)
	t.scenes[name] = scene
}

// TODO: add global event handling
func (t *TerminalMainApp) DrawStatusBar(input string) {
	w, h := t.screen.Size()
	textStyle := tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite)

	text := input
	nullspace := w - len(text)

	for i := 0; i < nullspace; i++ {
		text += " "
	}

	utils.WriteTo(t.screen, 0, h-1, text, textStyle)
}

func (t *TerminalMainApp) Mainloop() {
	flogger.Printf("Starting main loop, initial scene: %s", t.currunt_scene)
	t.scenes[t.currunt_scene].InitScene(t.screen, t.appctx)

	for {
		select {
		case exitCode := <-t.termination_signal:
			flogger.Printf("Termination signal received with exit code: %d", exitCode)
			close(t.termination_signal)
			close(t.transision_signal)
			t.screen.Fini()
			os.Exit(exitCode)

		case target := <-t.transision_signal:
			flogger.Printf("Transitioning from scene '%s' to '%s'", t.currunt_scene, target)
			if _, ok := t.scenes[target]; ok {
				t.currunt_scene = target
			} else {
				flogger.Printf("Transition target scene '%s' not found, staying in '%s'", target, t.currunt_scene)
			}

			t.screen.Clear()
			t.screen.Show()

			t.scenes[t.currunt_scene].InitScene(t.screen, t.appctx)
			t.screen.Show()
			continue

		default:
		}

		event := t.screen.PollEvent()

		switch ev := event.(type) {
		case *tcell.EventResize:
			flogger.Printf("Window resized")
			t.scenes[t.currunt_scene].WindowChangedScene(t.screen, t.appctx)
			t.screen.Clear()
			t.screen.Sync()
		default:
			flogger.Printf("Received event: %T", ev)
		}

		t.scenes[t.currunt_scene].DoScene(t.screen, event, t.appctx)
		t.screen.Show()
	}
}

type MainAppContexts struct {
	termination_signal chan int
	transision_signal  chan string
	scene              *string

	FooterFunc func(string)
}

func NewMainAppCtx(termsig chan int, transsig chan string, footerfunc func(string), scene *string) *MainAppContexts {
	Ctx := new(MainAppContexts)

	Ctx.termination_signal = termsig
	Ctx.transision_signal = transsig
	Ctx.FooterFunc = footerfunc
	Ctx.scene = scene

	return Ctx
}

func (m *MainAppContexts) Exit(exitcode int) {
	m.termination_signal <- exitcode
}

func (m *MainAppContexts) TranslateTo(name string) {
	m.transision_signal <- name
}

func (m *MainAppContexts) DrawFooterbar(content string) {
	m.FooterFunc(content)
}

func (m *MainAppContexts) GetCurruntScene() string {
	return *m.scene
}