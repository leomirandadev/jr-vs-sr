package capsulas

import (
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/leomirandadev/capsulas/internal/models"
	"github.com/leomirandadev/capsulas/internal/services"
	"github.com/leomirandadev/capsulas/pkg/customerr"
	"github.com/leomirandadev/capsulas/pkg/httprouter"
	"github.com/leomirandadev/capsulas/pkg/validator"
)

func Init(
	router httprouter.Router,
	srv *services.Container,
) {
	ctr := controllers{
		srv,
	}

	router.POST("/v1/capsulas", ctr.Create)
	router.GET("/v1/capsulas/{id}", ctr.GetByID)
	router.POST("/v1/capsulas/{id}/message", ctr.CreateMessage)
}

type controllers struct {
	srv *services.Container
}

// capsula swagger document
// @Description Create Capsula
// @Tags capsulas
// @Param capsula body models.CreateCapsulaReq true "create new capsula"
// @Accept json
// @Produce json
// @Success 201 {object} models.Capsula
// @Failure 400 {object} customerr.Error
// @Failure 401 {object} customerr.Error
// @Failure 500 {object} customerr.Error
// @Router /v1/capsulas [post]
func (ctr controllers) Create(c httprouter.Context) error {
	ctx := c.Context()

	var newCapsula models.CreateCapsulaReq
	if err := c.Decode(&newCapsula); err != nil {
		return customerr.WithStatus(http.StatusBadRequest, "decode error", err.Error())
	}

	if err := validator.Validate(newCapsula); err != nil {
		slog.WarnContext(ctx, "validator.validate", "err", err)
		return customerr.WithStatus(http.StatusBadRequest, "invalid payload", err.Error())
	}

	capsula, err := ctr.srv.Capsula.Create(ctx, newCapsula)
	if err != nil {
		slog.WarnContext(ctx, "create", "err", err)
		return err
	}

	return c.JSON(http.StatusCreated, capsula)
}

// capsula swagger document
// @Description Get one capsula
// @Tags capsulas
// @Param id path string true "CapsulaID"
// @Accept json
// @Produce json
// @Success 200 {object} models.Capsula
// @Failure 400 {object} customerr.Error
// @Failure 401 {object} customerr.Error
// @Failure 500 {object} customerr.Error
// @Router /v1/capsulas/{id} [get]
func (ctr controllers) GetByID(c httprouter.Context) error {
	ctx := c.Context()

	capsulaID := c.GetPathParam("id")
	if capsulaID == "" {
		slog.InfoContext(ctx, "you need to provide an id")
		return customerr.WithStatus(http.StatusBadRequest, "you need to provide an id", map[string]string{"capsula_id": capsulaID})
	}

	capsula, err := ctr.srv.Capsula.GetByID(ctx, capsulaID)
	if err != nil {
		slog.WarnContext(ctx, "get_by_id", "err", err)
		return err
	}

	return c.JSON(http.StatusOK, capsula)
}

// capsula swagger document
// @Description Create Capsula message
// @Tags capsulas
// @Param id path string true "CapsulaID"
// @Param message body models.CreateMessageReq true "create new capsula"
// @Accept json
// @Produce json
// @Success 201 {object} models.Capsula
// @Failure 400 {object} customerr.Error
// @Failure 401 {object} customerr.Error
// @Failure 500 {object} customerr.Error
// @Router /v1/capsulas/{id}/message [post]
func (ctr controllers) CreateMessage(c httprouter.Context) error {
	ctx := c.Context()

	var new models.CreateMessageReq
	if err := c.Decode(&new); err != nil {
		return customerr.WithStatus(http.StatusBadRequest, "decode error", err.Error())
	}

	if err := validator.Validate(new); err != nil {
		slog.WarnContext(ctx, "validator.validate", "err", err)
		return customerr.WithStatus(http.StatusBadRequest, "invalid payload", err.Error())
	}

	new.CapsulaID = c.GetPathParam("id")
	if new.CapsulaID == "" {
		slog.InfoContext(ctx, "you need to provide an id")
		return customerr.WithStatus(http.StatusBadRequest, "you need to provide an id", map[string]string{"capsula_id": new.CapsulaID})
	}

	capsula, err := ctr.srv.Messages.Create(ctx, new)
	if err != nil {
		slog.WarnContext(ctx, "create message", "err", err)
		return err
	}

	return c.JSON(http.StatusCreated, capsula)
}

// capsulas swagger document
// @Description Capsula upload image
// @Tags capsulas
// @Param image formData file true "image"
// @Accept json
// @Produce json
// @Success 204
// @Failure 400 {object} customerr.Error
// @Failure 401 {object} customerr.Error
// @Failure 500 {object} customerr.Error
// @Security BearerAuth
// @Router /v1/capsulas/photo/upload [post]
func (ctr controllers) Upload(c httprouter.Context) error {

	var (
		err error
	)
	data, ext, err := getFile(c.GetRequestReader(), "image")
	if err != nil {
		return err
	}

	slog.Debug("here", "data", data, "ext", ext)

	return c.JSON(http.StatusNoContent, nil)
}

func getFile(r *http.Request, param string) ([]byte, string, error) {
	file, handler, err := r.FormFile(param)
	if err != nil {
		slog.Warn("get form file fails", "err", err)
		return nil, "", customerr.WithStatus(http.StatusBadRequest, "get image fails", nil)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		slog.Warn("io.ReadAll fails", "err", err)
		return nil, "", customerr.WithStatus(http.StatusBadRequest, "get image fails", nil)
	}

	filenameSplitted := strings.Split(handler.Filename, ".")

	return data, filenameSplitted[len(filenameSplitted)-1], nil
}
