package emails_test

import (
	"os"
	"testing"

	"github.com/go-bolo/bolo"
	emails "github.com/go-bolo/emails"
	"github.com/pkg/errors"
)

var appInstance bolo.App

func TestMain(m *testing.M) {
	GetAppInstance()

	if code := m.Run(); code != 0 {
		os.Exit(code)
	}
}

func GetAppInstance() bolo.App {
	if appInstance != nil {
		return appInstance
	}

	app := bolo.Init(&bolo.AppOptions{})
	// start this plugin:
	app.RegisterPlugin(emails.NewPlugin(&emails.PluginCfg{}))

	err := app.Bootstrap()
	if err != nil {
		panic(err)
	}

	err = app.GetDB().AutoMigrate(
		&emails.EmailTemplateModel{},
		&emails.EmailModel{},
	)

	if err != nil {
		panic(errors.Wrap(err, "emails.GetAppInstance Error on run auto migration"))
	}

	return app
}
