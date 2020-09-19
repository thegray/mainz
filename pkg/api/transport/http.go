package transport

import (
	"pbk-main/pkg/server/mw"
	"pbk-main/pkg/service"

	"github.com/labstack/echo"
)

type ServiceMap struct {
	credential *service.CredService
	safetybox  *service.SafetyBox
}

func NewServices(cred *service.CredService, sb *service.SafetyBox) *ServiceMap {
	srvs := ServiceMap{
		credential: cred,
		safetybox:  sb,
	}
	return &srvs
}

func (srvs *ServiceMap) Init(r *echo.Group) {

	// assign route to the services
	r.POST("/login", srvs.credential.Login)
	r.POST("/login/refresh", srvs.credential.Refresh)
	r.POST("/logout", mw.JWTAuthorization(srvs.credential.Logout))
	r.POST("/cred/change", srvs.credential.ChangePassword) //TODO

	// f*ck restful, all POST
	r.POST("/safetybox/list", srvs.safetybox.List) //TODO: Implement pagination
	// endpoint to view single item
	r.POST("/safetybox", mw.JWTAuthorization(srvs.safetybox.View))
	// endpoint to create
	r.POST("/safetybox/create", mw.JWTAuthorization(srvs.safetybox.Create))
	// endpoint to update
	r.POST("/safetybox/update", mw.JWTAuthorization(srvs.safetybox.Update))
	// endpoint to update secret
	r.POST("/safetybox/secret/update", mw.JWTAuthorization(srvs.safetybox.UpdateSecret))
	// endpoint to reveal/decrypt
	r.POST("/safetybox/secret/reveal", mw.JWTAuthorization(srvs.safetybox.RevealSecret))
	// endpoint to delete
	r.POST("/safetybox/delete", mw.JWTAuthorization(srvs.safetybox.Delete))

}
