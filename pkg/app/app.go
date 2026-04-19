package app

import (
	"os"
	"github.com/sirupsen/logrus"
	"devtool/pkg/config"
	"devtool/pkg/gui"
)

type App struct {
	Log     *logrus.Entry
	Config  *config.AppConfig
	Gui     *gui.Gui
}

func Setup(debug bool) (*App, error) {
	log := logrus.New()
	if debug {
		log.SetLevel(logrus.DebugLevel)
		f, _ := os.OpenFile("devtool.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		log.SetOutput(f)
	} else {
		log.SetLevel(logrus.WarnLevel)
	}

	logger := log.WithField("app", "devtool")

	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	g, err := gui.NewGui(logger, cfg)
	if err != nil {
		return nil, err
	}

	return &App{
		Log:     logger,
		Config:  cfg,
		Gui:     g,
	}, nil
}

func (app *App) Run() error {
	defer app.Gui.Close()
	return app.Gui.Run()
}
