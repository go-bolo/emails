package emails

import (
	"fmt"

	"github.com/go-bolo/bolo"
	"github.com/go-bolo/bolo/helpers"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EmailModel struct {
	ID        uint64      `gorm:"primary_key;column:id" json:"id"`
	EmailId   string      `gorm:"column:emailId" json:"emailId"`
	From      string      `gorm:"type:TEXT;column:from" json:"from"`
	To        string      `gorm:"type:TEXT;column:to" json:"to"`
	Cc        string      `gorm:"type:TEXT;column:cc" json:"cc"`
	Bcc       string      `gorm:"type:TEXT;column:bcc" json:"bcc"`
	ReplyTo   string      `gorm:"type:TEXT;column:replyTo" json:"replyTo"`
	InReplyTo string      `gorm:"type:TEXT;column:inReplyTo" json:"inReplyTo"`
	Variables interface{} `gorm:"-" json:"variables"`
	Subject   string      `gorm:"type:TEXT;column:subject" json:"subject"`
	Text      string      `gorm:"type:MEDIUMTEXT;column:text" json:"text"`
	Html      string      `gorm:"type:MEDIUMTEXT;column:html" json:"html"`
	Type      string      `gorm:"column:type" json:"type"`
	Status    string      `gorm:"column:status;not null;default:added" json:"status"`
}

func (EmailModel) TableName() string {
	return "emails"
}

func (EmailModel) LoadData() error {
	return nil
}

func (r *EmailModel) Delete() error {
	db := bolo.GetDefaultDatabaseConnection()
	return db.Unscoped().Delete(&r).Error
}

// Save - Create if is new or update
func (m *EmailModel) Save() error {
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

func EmailFindOne(id string, record *EmailModel) error {
	db := bolo.GetDefaultDatabaseConnection()

	err := db.First(&record, id).Error
	if err != nil {
		return err
	}

	return nil
}

type EmailQueryOpts struct {
	Records *[]*EmailModel
	Count   *int64
	Limit   int
	Offset  int
	C       echo.Context
	IsHTML  bool
}

func EmailQueryAndCountReq(opts *EmailQueryOpts) error {
	db := bolo.GetDefaultDatabaseConnection()

	c := opts.C

	q := c.QueryParam("q")

	query := db
	// query.DryRun = true

	ctx := c.(*bolo.RequestContext)

	queryI, err := ctx.Query.SetDatabaseQueryForModel(query, &EmailModel{})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": fmt.Sprintf("%+v\n", err),
		}).Error("EmailQueryAndCountReq error")
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

	return EmailCountReq(opts)
}

func EmailCountReq(opts *EmailQueryOpts) error {
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
		}).Error("EmailCountReq error")
	}
	queryCount = queryICount.(*gorm.DB)

	return queryCount.
		Table(EmailModel{}.TableName()).
		Count(opts.Count).Error
}
