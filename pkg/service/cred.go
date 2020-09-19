package service

import (
	"errors"
	"log"
	"net/http"
	"pbk-main/pkg/application"
	"pbk-main/pkg/common/errcode"
	"pbk-main/pkg/model/common"
	"pbk-main/pkg/model/request"
	"pbk-main/pkg/model/response"
	"pbk-main/pkg/server"
	"pbk-main/pkg/util/jwtutil"

	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v9"
)

type CredService struct {
	app *application.CredApp
}

func NewCredService(app *application.CredApp) *CredService {
	credsvc := CredService{app: app}
	return &credsvc
}

func (cred *CredService) Login(c echo.Context) error {
	req := &request.LoginReq{}
	if err := c.Bind(req); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			log.Printf("[Service][Cred][Login] Request Validation fail: %+v \n", e)
		}
		return c.JSON(http.StatusBadRequest, response.Error{Code: errcode.GBadRequest, Message: "invalid body request"})
	}

	userAgent := c.Request().Header.Get("User-Agent")
	ip := c.RealIP()
	res, err := cred.app.CheckLogin(req.User, req.Pass, userAgent, ip)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, response.Error{Code: errcode.CREDLoginError, Message: "invalid login"})
	}

	return c.JSON(http.StatusOK, res)
}

func (cred *CredService) Logout(c echo.Context) error {
	userInt := c.(*server.CustomContext).GetCustomInfo("userinfo")
	if userInt == nil {
		log.Println("[Service][Cred][Logout] Cannot find 'userinfo' in request")
		return c.JSON(http.StatusUnprocessableEntity, response.Error{Code: errcode.GAuthError, Message: "invalid data"})
	}

	user, ok := userInt.(*common.UserInfo)
	if !ok {
		log.Println("[Service][Cred][Logout] invalid 'userinfo'")
		return c.JSON(http.StatusUnprocessableEntity, response.Error{Code: errcode.GAuthError, Message: "invalid data user"})
	}

	err := cred.app.Logout(user.UserName, user.UUIDAccess)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, response.Error{Code: errcode.GAuthError, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, response.General{Message: "success logout"})
}

func (cred *CredService) Refresh(c echo.Context) error {
	req := &request.RefreshReq{}
	if err := c.Bind(req); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			log.Printf("[Service][Cred][Login] Request Validation fail: %+v \n", e)
		}
		return c.JSON(http.StatusBadRequest, response.Error{Code: errcode.GBadRequest, Message: "invalid body request"})
	}

	// check if refresh token is valid
	token, err := jwtutil.VerifyToken(req.RefreshToken, "refresh")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, response.Error{Code: errcode.GAuthError, Message: "invalid token"})
	}

	// then initiated create new refresh token
	rclaim, err := jwtutil.ExtractTokenToRefreshClaim(token)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.Error{Code: errcode.GAuthError, Message: err.Error()})
	}

	// get required request identity
	userAgent := c.Request().Header.Get("User-Agent")
	ip := c.RealIP()

	res, err := cred.app.CreateRefreshToken(rclaim, userAgent, ip)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.Error{Code: errcode.CREDRefreshError_FailCreateToken, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

func (cred *CredService) ChangePassword(c echo.Context) error {
	return errors.New("Change password")
}
