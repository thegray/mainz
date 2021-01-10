package service

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"pbk-main/pkg/application"
	"pbk-main/pkg/common/errcode"
	"pbk-main/pkg/model/request"
	"pbk-main/pkg/model/response"

	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v9"
)

type SafetyBox struct {
	app *application.SafetyBoxApp
}

func NewSafetyBoxService(a *application.SafetyBoxApp) *SafetyBox {
	sfb := SafetyBox{app: a}
	return &sfb
}

func (sb *SafetyBox) List(c echo.Context) error {
	return errors.New("list fail")
}

func (sb *SafetyBox) View(c echo.Context) error {
	req := &request.ByIdReq{}
	if err := c.Bind(req); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			log.Printf("[Service][SafetyBox][View] Request Validation fail: %+v \n", e)
		}
		return c.JSON(http.StatusBadRequest, response.Error{Code: errcode.GBadRequest, Message: "invalid body request"})
	}

	data, err := sb.app.GetSafetyBoxById(*req.ID)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.Error{Code: errcode.SBNotExist, Message: "item not exist"})
	}
	return c.JSON(http.StatusOK, data)
}

func (sb *SafetyBox) Create(c echo.Context) error {
	req := &request.CreateAccSBReq{}
	if err := c.Bind(req); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			log.Printf("[Service][SafetyBox][Create] Request Validation fail: %+v \n", e)
		}
		return c.JSON(http.StatusBadRequest, response.Error{Code: errcode.GBadRequest, Message: "invalid body request"})
	}
	// **** special validation ****
	var isSecretInputed bool
	if isSecretInputed = (req.SecretInfo != "") || (req.SecretPass != ""); !isSecretInputed {
		return c.JSON(http.StatusBadRequest, response.Error{Code: errcode.GBadRequest, Message: "please fill the secret"})
	}
	isKeyInputed := req.Key != ""
	if !((isKeyInputed || isSecretInputed) && !(isKeyInputed && isSecretInputed)) { // XOR expression
		return c.JSON(http.StatusBadRequest, response.Error{Code: errcode.GBadRequest, Message: "please fill both the secrets AND the key, or none at all"})
	}
	// **** end ****
	createdId, err := sb.app.CreateSafetyBox(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error{Code: errcode.SBFailInsert, Message: "failed insert new item"})
	}

	return c.JSON(http.StatusOK, &response.General{Message: fmt.Sprintf("created with id: %d", createdId)})
}

func (sb *SafetyBox) Update(c echo.Context) error {
	req := &request.UpdateSBReq{}
	if err := c.Bind(req); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			log.Printf("[Service][SafetyBox][Update] Request Validation fail: %+v \n", e)
		}
		return c.JSON(http.StatusBadRequest, response.Error{Code: errcode.GBadRequest, Message: "invalid body request"})
	}
	err := sb.app.UpdateSafetyBox(*req.ID, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error{Code: errcode.SBFailUpdate, Message: "failed update item"})
	}

	return c.JSON(http.StatusOK, &response.General{Message: fmt.Sprintf("item updated, id: %d", *req.ID)})
}

func (sb *SafetyBox) UpdateSecret(c echo.Context) error {
	req := &request.UpdateSecretReq{}
	if err := c.Bind(req); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			log.Printf("[Service][SafetyBox][UpdSecret] Request Validation fail: %+v \n", e)
		}
		return c.JSON(http.StatusBadRequest, response.Error{Code: errcode.GBadRequest, Message: "invalid body request"})
	}
	err := sb.app.UpdateSecret(*req.ID, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error{Code: errcode.SBFailUpdateSecret, Message: "failed update secret"})
	}

	return c.JSON(http.StatusOK, &response.General{Message: fmt.Sprintf("secret updated, id: %d", *req.ID)})
}

func (sb *SafetyBox) RevealSecret(c echo.Context) error {
	req := &request.RevealSecretReq{}
	if err := c.Bind(req); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			log.Printf("[Service][SafetyBox][Reveal] Request Validation fail: %+v \n", e)
		}
		return c.JSON(http.StatusBadRequest, response.Error{Code: errcode.GBadRequest, Message: "invalid body request"})
	}
	sp, si, err := sb.app.RevealSecret(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error{Code: errcode.SBFailReveal, Message: "failed reveal"})
	}

	return c.JSON(http.StatusOK, &response.Revealed{SecretPass: sp, SecretInfo: si})
}

func (sb *SafetyBox) Delete(c echo.Context) error {
	req := &request.ByIdReq{}
	if err := c.Bind(req); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			log.Printf("[Service][SafetyBox][Delete] Request Validation fail: %+v \n", e)
		}
		return c.JSON(http.StatusBadRequest, response.Error{Code: errcode.GBadRequest, Message: "invalid body request"})
	}

	result, err := sb.app.DeleteSafetyBox(*req.ID)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.Error{Code: errcode.SBNotExist, Message: "item not exist"})
	}
	if result < 1 {
		return c.JSON(http.StatusOK, &response.General{Message: fmt.Sprintf("no item deleted, id: %d", *req.ID)})
	}
	return c.JSON(http.StatusOK, &response.General{Message: fmt.Sprintf("item deleted, id: %d", *req.ID)})
}
