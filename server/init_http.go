package server

import (
	"github.com/ahmadrezamusthafa/logwatcher/common/commonhandlers"
	"github.com/ahmadrezamusthafa/logwatcher/common/logger"
	"github.com/ahmadrezamusthafa/logwatcher/config"
	"github.com/ahmadrezamusthafa/logwatcher/server/http/health"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"net/http"
)

type HttpServer struct {
	Config      config.Config
	RootHandler *RootHandler
	router      *mux.Router
}

type RootHandler struct {
	Health *health.Handler `inject:"healthHandler"`
}

func (svr *HttpServer) Serve() {
	route := createRouter(svr.RootHandler)
	commonHandlers := commonhandlers.New()
	middleware := commonhandlers.Chain(
		commonHandlers.RecoverHandler,
		commonHandlers.LoggingHandler,
	)

	logger.Info("Http server is serving on :" + svr.Config.HttpPort)
	err := http.ListenAndServe(":"+svr.Config.HttpPort, middleware.Then(route))
	if err != nil {
		logger.Err(err.Error())
	}
}

func createRouter(rh *RootHandler) *gin.Engine {
	router := gin.Default()

	router.Use(static.Serve("/", static.LocalFile("./web/dist", true)))
	router.GET("/health", rh.Health.HealthCheck)

	return router
}
