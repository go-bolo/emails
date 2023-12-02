package emails

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aymerick/raymond"
	"github.com/go-bolo/bolo"
	"github.com/go-bolo/bolo/helpers"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/vanng822/go-premailer/premailer"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EmailTemplateModel struct {
	ID        uint64    `gorm:"primary_key;column:id" json:"id" filter:"param:id;type:number"`
	Subject   string    `gorm:"column:subject" json:"subject" filter:"param:subject"`
	Text      string    `gorm:"type:MEDIUMTEXT;column:text" json:"text" filter:"param:text"`
	Css       string    `gorm:"column:css" json:"css" filter:"param:css"`
	Html      string    `gorm:"type:MEDIUMTEXT;column:html" json:"html" filter:"param:html"`
	Type      string    `gorm:"column:type;unique" json:"type" filter:"param:type"`
	CreatedAt time.Time `gorm:"column:createdAt;" json:"createdAt" filter:"param:createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt;" json:"updatedAt" filter:"param:updatedAt"`
}

func (r *EmailTemplateModel) ToJSON() string {
	jsonString, _ := json.Marshal(r)
	return string(jsonString)
}

// TableName - Table name
func (EmailTemplateModel) TableName() string {
	return "email_templates"
}

func (r *EmailTemplateModel) Render(ctx TemplateVariables, email *Email) error {
	if ctx == nil {
		ctx = TemplateVariables{}
	}

	if r.Subject != "" {
		html, err := raymond.Render(r.Subject, ctx)
		if err != nil {
			return err
		}

		email.Subject = html
	} else if true {
		// TODO! add a default subject
	}

	if r.Text != "" {
		html, err := raymond.Render(r.Text, ctx)
		if err != nil {
			return err
		}

		email.Text = html
	}

	if r.Html != "" {
		html, err := raymond.Render(r.Html, ctx)
		if err != nil {
			return err
		}

		email.HTML = html
	}

	if email.HTML != "" && r.Css != "" {
		email.HTML = "<style>" + r.Css + "</style>" + email.HTML

		prem, err := premailer.NewPremailerFromString(email.HTML, premailer.NewOptions())
		if err != nil {
			log.Fatal(err)
		}

		html, err := prem.Transform()
		if err != nil {
			log.Fatal(err)
		}

		email.HTML = html
	}

	return nil
}

func (r *EmailTemplateModel) LoadData() error {
	return nil
}

// Save - Create if is new or update
func (m *EmailTemplateModel) Save() error {
	var err error
	db := bolo.GetDefaultDatabaseConnection()

	if m.ID == 0 {
		// create ....
		err = db.Create(&m).Error
		if err != nil {
			return err
		}
	} else {
		// update ...
		err = db.Save(&m).Error
		if err != nil {
			return err
		}
	}

	return nil
}

// Delete - Delete
func (r *EmailTemplateModel) Delete() error {
	db := bolo.GetDefaultDatabaseConnection()
	return db.Unscoped().Delete(&r).Error
}

func EmailTemplateFindOne(id string, record *EmailTemplateModel) error {
	db := bolo.GetDefaultDatabaseConnection()

	err := db.First(&record, id).Error
	if err != nil {
		return err
	}

	return nil
}

func TemplateFindOneByType(t string, record *EmailTemplateModel) error {
	db := bolo.GetDefaultDatabaseConnection()

	err := db.
		Where("type = ?", t).
		First(&record).Error
	if err != nil {
		return err
	}

	return nil
}

type TemplateVariables map[string]string

type EmailTemplateQueryOpts struct {
	Records *[]*EmailTemplateModel
	Count   *int64
	Limit   int
	Offset  int
	C       echo.Context
	IsHTML  bool
}

func EmailTemplateQueryAndCountReq(opts *EmailTemplateQueryOpts) error {
	db := bolo.GetDefaultDatabaseConnection()

	c := opts.C
	q := c.QueryParam("q")
	query := db
	ctx := c.(*bolo.RequestContext)

	queryI, err := ctx.Query.SetDatabaseQueryForModel(query, &EmailTemplateModel{})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": fmt.Sprintf("%+v\n", err),
		}).Error("EmailTemplateQueryAndCountReq error")
	}
	query = queryI.(*gorm.DB)

	if q != "" {
		query = query.Where(
			db.Where("subject LIKE ?", "%"+q+"%").Or(db.Where("text LIKE ?", "%"+q+"%")),
		)
	}

	orderColumn, orderIsDesc, orderValid := helpers.ParseUrlQueryOrder(c.QueryParam("order"), c.QueryParam("sort"), c.QueryParam("sortDirection"))

	if orderValid {
		query = query.Order(clause.OrderByColumn{
			Column: clause.Column{Table: clause.CurrentTable, Name: orderColumn},
			Desc:   orderIsDesc,
		})
	} else {
		query = query.Order("id DESC")
	}

	query = query.Limit(opts.Limit).
		Offset(opts.Offset)

	r := query.Find(opts.Records)
	if r.Error != nil {
		return r.Error
	}

	return EmailTemplateCountReq(opts)
}
func EmailTemplateCountReq(opts *EmailTemplateQueryOpts) error {
	db := bolo.GetDefaultDatabaseConnection()

	c := opts.C
	q := c.QueryParam("q")

	ctx := c.(*bolo.RequestContext)

	// Count ...
	queryCount := db

	if q != "" {
		queryCount = queryCount.Or(
			db.Where("subject LIKE ?", "%"+q+"%"),
			db.Where("text LIKE ?", "%"+q+"%"),
		)
	}

	queryICount, err := ctx.Query.SetDatabaseQueryForModel(queryCount, &EmailTemplateModel{})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": fmt.Sprintf("%+v\n", err),
		}).Error("EmailTemplateCountReq error")
	}
	queryCount = queryICount.(*gorm.DB)

	return queryCount.
		Table(EmailTemplateModel{}.TableName()).
		Count(opts.Count).Error
}
