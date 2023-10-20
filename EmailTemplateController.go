package emails

import (
	"net/http"

	"github.com/go-bolo/bolo"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func NewEmailTemplateController(app bolo.App) *EmailTemplateController {
	return &EmailTemplateController{
		App: app,
	}
}

type EmailTemplateBodyRequest struct {
	Record *EmailTemplateModel `json:"email-template"`
}

type EmailTemplateJSONResponse struct {
	bolo.BaseListReponse
	Record []*EmailTemplateModel `json:"email-template"`
}

type EmailTemplateFindOneJSONResponse struct {
	Record *EmailTemplateModel `json:"email-template"`
}

type emailTemplateCountJSONResponse struct {
	Count int64 `json:"count"`
}

type EmailTemplateController struct {
	App bolo.App
}

// Query - Query email templates
func (ctl *EmailTemplateController) Query(c echo.Context) error {
	var err error
	ctx := c.(*bolo.RequestContext)

	var count int64

	records := []*EmailTemplateModel{}
	err = EmailTemplateQueryAndCountReq(&EmailTemplateQueryOpts{
		Records: &records,
		Count:   &count,
		Limit:   ctx.GetLimit(),
		Offset:  ctx.GetOffset(),
		C:       c,
	})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Debug("EmailTemplateController query Error on find records")
		return err
	}

	ctx.Pager.Count = count

	logrus.WithFields(logrus.Fields{
		"count":             count,
		"len_records_found": len(records),
	}).Debug("ContentFindAll count result")

	for i := range records {
		records[i].LoadData()
	}

	resp := EmailTemplateJSONResponse{
		Record: records,
	}

	resp.Meta.Count = count

	return c.JSON(200, &resp)
}

func (ctl *EmailTemplateController) Create(c echo.Context) error {
	var err error
	ctx := c.(*bolo.RequestContext)

	can := ctx.Can("create_email-template")
	if !can {
		return echo.NewHTTPError(http.StatusForbidden, "Forbidden")
	}

	var body EmailTemplateBodyRequest

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
		}).Debug("EmailTemplateController query Error on find records")
		return err
	}

	err = record.LoadData()
	if err != nil {
		return err
	}

	resp := EmailTemplateFindOneJSONResponse{
		Record: record,
	}

	return c.JSON(http.StatusCreated, &resp)
}

func (ctl *EmailTemplateController) Count(c echo.Context) error {
	var err error
	RequestContext := c.(*bolo.RequestContext)

	var count int64
	err = EmailTemplateCountReq(&EmailTemplateQueryOpts{
		Count:  &count,
		Limit:  RequestContext.GetLimit(),
		Offset: RequestContext.GetOffset(),
		C:      c,
	})

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Debug("EmailTemplateController query Error on find records")
	}

	RequestContext.Pager.Count = count

	resp := emailCountJSONResponse{}
	resp.Count = count

	return c.JSON(200, &resp)
}

func (ctl *EmailTemplateController) FindOne(c echo.Context) error {
	id := c.Param("id")
	// ctx := c.(*bolo.RequestContext)

	logrus.WithFields(logrus.Fields{
		"id": id,
	}).Debug("EmailTemplateController FindOne id from params")

	record := EmailTemplateModel{}
	err := EmailTemplateFindOne(id, &record)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"id":    id,
			"error": err,
		}).Debug("EmailTemplateController FindOne Error on find record")
		return err
	}

	if record.ID == 0 {
		logrus.WithFields(logrus.Fields{
			"id": id,
		}).Debug("EmailTemplateController FindOne record not found")

		return echo.NotFoundHandler(c)
	}

	record.LoadData()

	resp := EmailTemplateFindOneJSONResponse{
		Record: &record,
	}

	return c.JSON(200, &resp)
}

func (ctl *EmailTemplateController) Update(c echo.Context) error {
	var err error

	id := c.Param("id")

	RequestContext := c.(*bolo.RequestContext)

	logrus.WithFields(logrus.Fields{
		"id":    id,
		"roles": RequestContext.GetAuthenticatedRoles(),
	}).Debug("EmailTemplateController Update")

	can := RequestContext.Can("update_email-template")
	if !can {
		return echo.NewHTTPError(http.StatusForbidden, "Forbidden")
	}

	record := EmailTemplateModel{}
	err = EmailTemplateFindOne(id, &record)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"id":    id,
			"error": err,
		}).Debug("EmailTemplateController Update Error on find record")
		return errors.Wrap(err, "EmailTemplateController Update Error on find record")
	}

	record.LoadData()

	body := EmailTemplateFindOneJSONResponse{Record: &record}

	if err := c.Bind(&body); err != nil {
		logrus.WithFields(logrus.Fields{
			"id":    id,
			"error": err,
		}).Debug("EmailTemplateController Update Error on bind record")

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
		return err
	}

	resp := EmailTemplateFindOneJSONResponse{
		Record: &record,
	}

	return c.JSON(http.StatusOK, &resp)
}

func (ctl *EmailTemplateController) Delete(c echo.Context) error {
	var err error
	RequestContext := c.(*bolo.RequestContext)
	id := c.Param("id")

	logrus.WithFields(logrus.Fields{
		"id": id,
	}).Debug("EmailTemplateController id from params")

	record := EmailTemplateModel{}
	err = EmailTemplateFindOne(id, &record)
	if err != nil {
		return err
	}

	can := RequestContext.Can("delete_email-template")
	if !can {
		return echo.NewHTTPError(http.StatusForbidden, "Forbidden")
	}

	if record.ID == 0 {
		return c.JSON(http.StatusNotFound, make(map[string]string))
	}

	err = record.Delete()
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
