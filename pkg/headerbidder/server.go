package headerbidder

import (
	"context"
	"time"

	gocommonshttp "github.com/adgear/go-commons/pkg/fasthttp"
	"github.com/adgear/go-commons/pkg/log"
	"github.com/adgear/sps-header-bidder/pkg/demand"
	"github.com/adgear/sps-header-bidder/pkg/health"
	"github.com/adgear/sps-header-bidder/pkg/kafkaclient"
	"github.com/adgear/sps-header-bidder/pkg/privacy"
	"github.com/valyala/fasthttp"
)

// Handlers Params
type Handlers struct {
	// springServeHandler to handle SpringServe requests
	springServeHandler *SpringServeHandler
	// freeWheelHandler to handle FreeWheel requests
	freeWheelHandler *FreeWheelHandler
}

// HeaderBidder Params
type HeaderBidder struct {
	// handler struct
	handlers *Handlers
	//http web server
	webserver gocommonshttp.WebServer
	// server port
	port string
	// logger service
	logger log.Service
}

// CreateServer is a constructor function which get demandclient, port, logger implementation and
// kafkaClient params arguments and return implementation of Header Bidder interface.
func CreateServer(logger log.Service, kafka kafkaclient.Service, demandClient demand.Service,
	privacyClient privacy.Service, port string) *HeaderBidder {

	fwNewHandler := FwNewHandler(logger, kafka, demandClient, privacyClient)
	ssNewHandler := SsNewHandler(logger, kafka, demandClient, privacyClient)
	handlers := &Handlers{
		springServeHandler: ssNewHandler,
		freeWheelHandler:   fwNewHandler,
	}
	routes := routes(handlers)
	params := &gocommonshttp.Params{
		Compress:         true,
		EnablePrometheus: true,
		MaxWaitTime:      5 * time.Second,
		Routes:           routes,
	}
	webserver := gocommonshttp.New(logger, params)
	return &HeaderBidder{
		handlers:  handlers,
		port:      port,
		logger:    logger,
		webserver: webserver,
	}
}

// Run function will take  context as argument and start the server to listen on given port
func (s *HeaderBidder) Run(ctx context.Context) error {
	return s.webserver.ListenAndServe(ctx, "0.0.0.0:"+s.port)
}

// routes function will takes handlers as arguments and return all the routes with respective handlers
func routes(handlers *Handlers) *gocommonshttp.HttpRoutes {
	return &gocommonshttp.HttpRoutes{
		GetBindings: map[string]fasthttp.RequestHandler{
			"/health/liveness":  health.LivenessHandler,
			"/health/readiness": health.ReadinessHandler},
		PostBindings: map[string]fasthttp.RequestHandler{
			"/hb": handlers.springServeHandler.Handler,
			"/fw": handlers.freeWheelHandler.Handler},
		StaticFilesBindings: map[string]string{
			"/docs/{filepath:*}": "./pkg/swagger/"},
	}
}
