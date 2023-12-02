package emails

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"strings"

	"github.com/go-bolo/bolo"
	"github.com/go-bolo/msgbroker"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	gomail "gopkg.in/mail.v2"
)

type Email struct {
	Identifier string `json:"identifier"`

	Template *string `json:"template"`

	To      string `json:"to"`
	From    string `json:"from"`
	ReplyTo string `json:"replyTo"`
	CC      string `json:"cc"`
	CCO     string `json:"cco"`

	Subject string `json:"subject"`
	Text    string `json:"Text"`
	HTML    string `json:"HTML"`

	DeliveryAttempts int `json:"deliveryAttempts"`
}

func (r *Email) ToJSON() []byte {
	jsonString, _ := json.MarshalIndent(r, "", "  ")
	return jsonString
}

func (r *Email) Send() error {
	app := bolo.GetApp()
	cfgs := app.GetConfiguration()

	m := gomail.NewMessage()

	userName := cfgs.GetF("SMTP_USER", "")
	password := cfgs.GetF("SMTP_PASSWORD", "")
	from := cfgs.GetF("SMTP_FROM", "")

	host := cfgs.GetF("SMTP_HOST", "")
	port := cfgs.GetIntF("SMTP_PORT", 587)

	enableEmail := cfgs.GetF("ENABLE_EMAIL_DELIVERY", "")

	if userName == "" || password == "" || from == "" || host == "" || port == 0 {
		log.Println("Email.Send Email delivery configuration not found", r.To, r.Subject)
		return nil
	}

	if enableEmail == "" {
		logrus.WithFields(logrus.Fields{
			"template": r.Template,
			"to":       r.To,
			"subject":  r.Subject,
			"text":     r.Text,
		}).Warn("Email.Send Email delivery disabled, skiping")
		return nil
	}

	// Set E-Mail sender
	m.SetHeader("From", from)

	recipients := strings.Split(r.To, ",")
	addresses := make([]string, len(recipients))
	for i, recipient := range recipients {
		addresses[i] = m.FormatAddress(recipient, "")
	}

	// Set E-Mail receivers
	m.SetHeader("To", addresses...)

	if r.ReplyTo != "" {
		m.SetHeader("ReplyTo", r.ReplyTo)
	}

	// Set E-Mail subject
	m.SetHeader("Subject", r.Subject)

	if r.HTML != "" {
		m.SetBody("text/html", r.HTML)
	} else {
		m.SetBody("text/html", r.Text)
	}

	d := gomail.NewDialer(host, port, userName, password)

	env := app.GetConfiguration().GetF("GO_ENV", "development")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: env != "production"}

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// TODO! add support to only log with a configuration
	if err := d.DialAndSend(m); err != nil {
		log.Println("Email.Send Error on send email", err)
		return err
	}
	return nil
}

func (r *Email) QueueToSend() error {
	app := bolo.GetApp()
	cfgs := app.GetConfiguration()

	enableEmail := cfgs.GetF("ENABLE_EMAIL_DELIVERY", "")

	if enableEmail == "" {
		logrus.WithFields(logrus.Fields{
			"template": r.Template,
			"to":       r.To,
			"subject":  r.Subject,
			"text":     r.Text,
			"html":     r.HTML,
		}).Info("QueueToSend: Skipping email delivery")

		return nil
	} else {
		logrus.WithFields(logrus.Fields{
			"template": r.Template,
			"to":       r.To,
			"subject":  r.Subject,
		}).Debug("QueueToSend: Will send email to")
	}

	msgBI := app.GetPlugin("msg-broker")
	switch msgBrokerPlugin := msgBI.(type) {
	case *msgbroker.MSGBrokerPlugin:
		c := msgBrokerPlugin.Client

		if r.ReplyTo == "" {
			r.ReplyTo = cfgs.GetF("SMTP_REPLY_TO", "Monitor do Mercado <monitordomercado@linkysystems.com>")
		}

		// publish in rabbit mq ...
		return c.Publish("notification-email-delivery", r.ToJSON())
	default:
		return r.Send()
	}
}

func (r *Email) Requeue() error {
	r.DeliveryAttempts++

	if r.DeliveryAttempts < 3 {
		log.Println("notification.Email requeuing email", r.To, r.DeliveryAttempts, r.Subject)
		return r.QueueToSend()
	} else {
		log.Println("notification.Email max requeue limit, skiping requeue", r.To, r.DeliveryAttempts, r.Subject)
	}

	return nil
}

type EmailOpts struct {
	To      string `json:"to"`
	From    string `json:"from"`
	ReplyTo string `json:"replyTo"`
	CC      string `json:"cc"`
	CCO     string `json:"cco"`

	TemplateName string
	Variables    TemplateVariables
}

func NewEmailWithTemplate(opts *EmailOpts) (*Email, error) {
	e := Email{
		To:      opts.To,
		From:    opts.From,
		ReplyTo: opts.ReplyTo,
		CC:      opts.CC,
		CCO:     opts.CCO,
	}

	if opts.TemplateName == "" {
		return &e, nil
	}

	template := EmailTemplateModel{}
	err := TemplateFindOneByType(opts.TemplateName, &template)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrap(err, "Email.NewEmail error on find email template by type")
	}

	if template.ID == 0 {
		return &e, nil
	}

	logrus.WithFields(logrus.Fields{
		"type": template.Type,
		"text": template.Text,
		"html": template.Html,
	}).Debug("NewEmailWithTemplate tempplate found")

	err = template.Render(opts.Variables, &e)
	if err != nil {
		return nil, errors.Wrap(err, "Email.NewEmail error on render email template")
	}

	return &e, nil
}
