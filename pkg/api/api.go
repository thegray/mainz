package api

import (
	"log"
	"pbk-main/pkg/api/transport"
	"pbk-main/pkg/application"
	"pbk-main/pkg/db"
	"pbk-main/pkg/server"
	"pbk-main/pkg/service"
	"pbk-main/pkg/store"
)

func Start() {
	s := server.New()
	v1Route := s.Group("/main/api/v1")

	// read config here, and init other things
	//

	dbImpl := db.DBImpl{}
	// test connection to db
	err := dbImpl.TestConnect()
	if err != nil {
		log.Fatalf("[Server][API] Cannot start server, cannot init db: %+v", err)
		return
	}

	credStore := store.NewCredentialStore(&dbImpl)
	accStore := store.NewAccountStore(&dbImpl)
	sftbStore := store.NewSafetyBoxStore(&dbImpl)

	credApp := application.NewCredApp(credStore)
	accApp := application.NewAccountApp(accStore)
	sbApp := application.NewSafetyBoxApp(sftbStore)

	credSvc := service.NewCredService(credApp)
	accSvc := service.NewAccountService(accApp)
	sftBoxSvc := service.NewSafetyBoxService(sbApp)
	services := transport.NewServices(credSvc, sftBoxSvc, accSvc)

	services.Init(v1Route)

	server.Run(s)
}
