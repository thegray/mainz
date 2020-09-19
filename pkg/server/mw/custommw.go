package mw

import (
	"net/http"
	"pbk-main/pkg/common/errcode"
	"pbk-main/pkg/model/common"
	"pbk-main/pkg/model/response"
	"pbk-main/pkg/server"
	"pbk-main/pkg/store/memstore"
	"pbk-main/pkg/util/jwtutil"
	"strings"

	"github.com/labstack/echo"
)

func JWTAuthorization(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authorizationHeader := c.Request().Header.Get("Authorization")

		if !strings.Contains(authorizationHeader, "Bearer") {
			return c.JSON(http.StatusBadRequest, response.Error{Code: errcode.GAuthError, Message: "no token"})
		}

		tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)

		// // check validity of the access token
		// err := jwtutil.TokenValid(tokenString, "access")
		// if err != nil {
		// 	return c.JSON(http.StatusUnauthorized, response.Error{Code: errcode.GAuthError, Message: "invalid auth token"})
		// }

		// BIG TODO: Implement check to memstore
		token, err := jwtutil.VerifyToken(tokenString, "access")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, response.Error{Code: errcode.GAuthError, Message: "invalid token"})
		}

		claim, err := jwtutil.ExtractTokenToAccessClaim(token)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, response.Error{Code: errcode.GAuthError, Message: err.Error()})
		}

		authinfo, err := memstore.CheckAccs(claim.Name, claim.UUID)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, response.Error{Code: errcode.GAuthError, Message: err.Error()})
		}
		c.(*server.CustomContext).SetCustomInfo("authinfo", authinfo)
		userinfo := &common.UserInfo{
			// UserID: claim.UserID,
			UserName:   claim.Name,
			UUIDAccess: claim.UUID}
		c.(*server.CustomContext).SetCustomInfo("userinfo", userinfo)

		return next(c)
	}
}
