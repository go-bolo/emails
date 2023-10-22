package emails

import (
	"github.com/go-bolo/bolo"
	"github.com/gookit/event"
	"github.com/sirupsen/logrus"
)

type EmailPlugin struct {
	App                     bolo.App
	Name                    string
	EmailTypes              EmailTypes
	EmailController         *EmailController
	EmailTemplateController *EmailTemplateController
}

func (p *EmailPlugin) GetName() string {
	return p.Name
}

func (p *EmailPlugin) Init(app bolo.App) error {
	logrus.Debug(p.GetName() + " Init")

	p.App = app
	p.EmailController = NewEmailController(app)
	p.EmailTemplateController = NewEmailTemplateController(app)

	app.GetEvents().On("bindRoutes", event.ListenerFunc(func(e event.Event) error {
		return p.BindRoutes(app)
	}), event.Normal)

	return nil
}

func (p *EmailPlugin) BindRoutes(app bolo.App) error {
	logrus.Debug(p.GetName() + " BindRoutes")

	routerEmailApi := app.SetRouterGroup("email-api", "/api/email")
	app.SetResource("email", p.EmailController, routerEmailApi)

	routerEmailTemplateApi := app.SetRouterGroup("email-template-api", "/api/email-template")
	app.SetResource("email-template", p.EmailTemplateController, routerEmailTemplateApi)

	router := app.GetRouter()
	router.GET("/api/email-template-types", p.EmailController.GetEmailTemplateTypes)

	return nil
}

func (p *EmailPlugin) AddEmailTemplate(name string, t *EmailType) error {
	p.EmailTypes[name] = t
	return nil
}

func (p *EmailPlugin) GetMigrations() []*bolo.Migration {
	return []*bolo.Migration{}
}

type PluginCfg struct{}

func NewPlugin(cfg *PluginCfg) *EmailPlugin {
	p := EmailPlugin{Name: "emails", EmailTypes: EmailTypes{}}
	return &p
}
