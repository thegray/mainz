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

type Account struct {
	app *application.AccountApp
}

func NewAccountService(a *application.AccountApp) *Account {
	acc := Account{app: a}
	return &acc
}

func (ac *Account) List(c echo.Context) error {
	return errors.New("list fail")
}

func (ac *Account) View(c echo.Context) error {
	return c.JSON(http.StatusOK, "implement this")
}

func (ac *Account) Create(c echo.Context) error {
	req := &request.CreateAccSBReq{}
	if err := c.Bind(req); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			log.Printf("[Service][Account][Create] Request Validation fail: %+v \n", e)
		}
		return c.JSON(http.StatusBadRequest, response.Error{Code: errcode.GBadRequest, Message: "invalid body request"})
	}
	// **** special validation ****
	isSecretInputed := (req.SecretInfo != "") || (req.SecretPass != "")
	isKeyInputed := req.Key != ""
	if !((isKeyInputed || isSecretInputed) && !(isKeyInputed && isSecretInputed)) { // XOR expression
		return c.JSON(http.StatusBadRequest, response.Error{Code: errcode.GBadRequest, Message: "please fill both the secrets AND the key, or none at all"})
	}
	// **** end ****

	createdId, err := ac.app.CreateAccount(req, isKeyInputed)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error{Code: errcode.SBFailInsert, Message: "failed insert new item"})
	}

	return c.JSON(http.StatusOK, &response.General{Message: fmt.Sprintf("created with id: %d", createdId)})
}

func (ac *Account) Update(c echo.Context) error {
	return c.JSON(http.StatusOK, "implement this")
}

func (ac *Account) Delete(c echo.Context) error {
	return c.JSON(http.StatusOK, "implement this")
}
