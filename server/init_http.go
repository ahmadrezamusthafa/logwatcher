package server

import (
	"github.com/ahmadrezamusthafa/logwatcher/pkg/commonhandlers"
	"github.com/ahmadrezamusthafa/logwatcher/pkg/logger"
	"github.com/ahmadrezamusthafa/logwatcher/config"
	"github.com/ahmadrezamusthafa/logwatcher/server/http/health"
	"github.com/ahmadrezamusthafa/logwatcher/server/http/thirdparty"
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
	Health     *health.Handler     `inject:"healthHandler"`
	ThirdParty *thirdparty.Handler `inject:"thirdPartyHandler"`
}

func (svr *HttpServer) Serve() {
	route := createRouter(svr.RootHandler)
	commonHandlers := commonhandlers.New()
	middleware := commonhandlers.Chain(
		commonHandlers.RecoverHandler,
		commonHandlers.LoggingHandler,
	)

	logger.Info("Web is serving on " + "http://localhost:" + svr.Config.HttpPort)
	err := http.ListenAndServe(":"+svr.Config.HttpPort, middleware.Then(route))
	if err != nil {
		logger.Err(err.Error())
	}
}

func createRouter(rh *RootHandler) *gin.Engine {
	router := gin.Default()

	router.Use(static.Serve("/", static.LocalFile("./web", true)))
	router.GET("/health", rh.Health.HealthCheck)

	router.POST("/generate_query", rh.ThirdParty.GenerateQuery)
	router.POST("/query", rh.ThirdParty.Query)
	router.GET("/attributes", rh.ThirdParty.GetLogAttributes)

	return router
}
