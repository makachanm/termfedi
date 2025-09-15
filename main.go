package main

import (
	"termfedi/app"
	"termfedi/config"

	"github.com/makachanm/flogger-lib"
)

func main() {
	flogger.Println("Starting termfedi")

	cnf, err := config.LoadConfig()
	if err != nil {
		flogger.Errorf("Failed to load config: %v", err)
		panic(err)
	}

	m := app.NewTerminalScreen()
	err = m.InitTerminalScreen()
	if err != nil {
		flogger.Errorf("Failed to init terminal screen: %v", err)
		panic(err)
	}

	m.RegisterScene("main", app.NewTimelineScreen(cnf))
	m.RegisterScene("noti", app.NewNotificationScreen(cnf))
	m.RegisterScene("action", app.NewActionScreen(cnf))
	m.Mainloop()
}