package emails

import (
	"net/http"

	"github.com/go-bolo/bolo"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func NewEmailController(app bolo.App) *EmailController {
	return &EmailController{
		App: app,
	}
}

type EmailBodyRequest struct {
	Record *EmailModel `json:"email"`
}

type EmailJSONResponse struct {
	bolo.BaseListReponse
	Record []*EmailModel `json:"email"`
}

type EmailFindOneJSONResponse struct {
	Record *EmailModel `json:"email"`
}

type emailCountJSONResponse struct {
	Count int64 `json:"count"`
}

type EmailTypesResponse struct {
	EmailTypes EmailTypes `json:"emailTypes"`
}

func NewEmailTypesResponse() *EmailTypesResponse {
	return &EmailTypesResponse{
		EmailTypes: EmailTypes{},
	}
}

type EmailController struct {
	App bolo.App
}

// GetEmailTemplateTypes - Get all email template types
func (ctl *EmailController) GetEmailTemplateTypes(c echo.Context) error {
	resp := NewEmailTypesResponse()

	emailPlugin := ctl.App.GetPlugin("emails").(*EmailPlugin)
	resp.EmailTypes = emailPlugin.EmailTypes

	return c.JSON(200, resp)
}

func (ctl *EmailController) Query(c echo.Context) error {
	var err error
	RequestContext := c.(*bolo.RequestContext)

	var count int64
	records := make([]*EmailModel, 0)
	err = EmailQueryAndCountReq(&EmailQueryOpts{
		Records: &records,
		Count:   &count,
		Limit:   RequestContext.GetLimit(),
		Offset:  RequestContext.GetOffset(),
		C:       c,
	})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Debug("EmailController query Error on find records")
		return err
	}

	RequestContext.Pager.Count = count

	logrus.WithFields(logrus.Fields{
		"count":             count,
		"len_records_found": len(records),
	}).Debug("ContentFindAll count result")

	for i := range records {
		records[i].LoadData()
	}

	resp := EmailJSONResponse{
		Record: records,
	}

	resp.Meta.Count = count

	return c.JSON(200, &resp)
}

func (ctl *EmailController) Create(c echo.Context) error {
	var err error
	ctx := c.(*bolo.RequestContext)

	can := ctx.Can("create_email")
	if !can {
		return echo.NewHTTPError(http.StatusForbidden, "Forbidden")
	}

	var body EmailBodyRequest

	if err := c.Bind(&body); err != nil {
		if _, ok := err.(*echo.HTTPError); ok {
			return err
		}
		return c.NoContent(http.StatusNotFound)
	}

	record := body.Record
	record.ID = 0

	if err := c.Validate(record); err != nil {
		if _, ok := err.(*echo.HTTPError); ok {
			return err
		}
		return err
	}

	err = record.Save()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Debug("EmailController Create Error on save record")
		return err
	}

	err = record.LoadData()
	if err != nil {
		return err
	}

	resp := EmailFindOneJSONResponse{
		Record: record,
	}

	return c.JSON(http.StatusCreated, &resp)
}

func (ctl *EmailController) Count(c echo.Context) error {
	var err error
	ctx := c.(*bolo.RequestContext)

	can := ctx.Can("find_email")
	if !can {
		return echo.NewHTTPError(http.StatusForbidden, "Forbidden")
	}

	var count int64
	err = EmailCountReq(&EmailQueryOpts{
		Count:  &count,
		Limit:  ctx.GetLimit(),
		Offset: ctx.GetOffset(),
		C:      c,
	})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Debug("EmailController.Count Error on find contents")
	}

	resp := emailCountJSONResponse{
		Count: count,
	}

	return c.JSON(200, &resp)
}

func (ctl *EmailController) FindOne(c echo.Context) error {
	id := c.Param("id")
	ctx := c.(*bolo.RequestContext)

	logrus.WithFields(logrus.Fields{
		"id": id,
	}).Debug("ContentFindOne id from params")

	can := ctx.Can("find_email")
	if !can {
		return echo.NewHTTPError(http.StatusForbidden, "Forbidden")
	}

	record := EmailModel{}
	err := EmailFindOne(id, &record)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Debug("ContentFindOne Error on find content")
		return err
	}

	err = record.LoadData()

	resp := EmailFindOneJSONResponse{
		Record: &record,
	}

	return c.JSON(200, &resp)
}

func (ctl *EmailController) Update(c echo.Context) error {
	var err error
	ctx := c.(*bolo.RequestContext)

	can := ctx.Can("update_email")
	if !can {
		return echo.NewHTTPError(http.StatusForbidden, "Forbidden")
	}

	id := c.Param("id")

	logrus.WithFields(logrus.Fields{
		"id": id,
	}).Debug("EmailController.Update id from params")

	record := EmailModel{}
	err = EmailFindOne(id, &record)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"id":    id,
			"error": err,
		}).Debug("EmailController.Update error on find one")
		return errors.Wrap(err, "EmailController.Update error on find one")
	}

	err = record.LoadData()
	if err != nil {
		return errors.Wrap(err, "EmailController.Update error on LoadData")
	}

	body := EmailBodyRequest{Record: &record}

	if err := c.Bind(&body); err != nil {
		if _, ok := err.(*echo.HTTPError); ok {
			return err
		}
		return c.NoContent(http.StatusNotFound)
	}

	if err := c.Validate(record); err != nil {
		if _, ok := err.(*echo.HTTPError); ok {
			return err
		}
		return err
	}

	err = record.Save()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Debug("EmailController Update Error on save record")
		return err
	}

	resp := EmailFindOneJSONResponse{
		Record: &record,
	}

	return c.JSON(200, &resp)
}

func (ctl *EmailController) Delete(c echo.Context) error {
	var err error
	ctx := c.(*bolo.RequestContext)
	id := c.Param("id")

	logrus.WithFields(logrus.Fields{
		"id": id,
	}).Debug("EmailController.Delete id from params")

	record := EmailModel{}
	err = EmailFindOne(id, &record)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"id":    id,
			"error": err,
		}).Debug("EmailController.Delete error on find one")
		return errors.Wrap(err, "EmailController.Delete error on find one")
	}

	if record.ID == 0 {
		return c.JSON(http.StatusNotFound, make(map[string]string))
	}

	can := ctx.Can("delete_email")
	if !can {
		return echo.NewHTTPError(http.StatusForbidden, "Forbidden")
	}

	err = record.Delete()
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
