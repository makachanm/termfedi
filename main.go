package main

import (
	"termfedi/app"
	"termfedi/config"
)

func main() {
	cnf, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	m := app.NewTerminalScreen()
	err = m.InitTerminalScreen()
	if err != nil {
		panic(err)
	}

	m.RegisterScene("main", app.NewTimelineScreen(cnf))
	m.Mainloop()
}
