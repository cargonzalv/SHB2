package headerbidder

import (
	"encoding/json"
	"fmt"

	"github.com/adgear/go-commons/pkg/log"
	"github.com/adgear/sps-header-bidder/pkg/demand"
	"github.com/adgear/sps-header-bidder/pkg/kafkaclient"
	"github.com/adgear/sps-header-bidder/pkg/privacy"
	"github.com/tidwall/gjson"

	"github.com/valyala/fasthttp"
)

const SAM_FLAGS_SS_PATH string = "ext.sam_flags"
const SAM_HB_TAG_SS_PATH string = "ext.sam_hb_tag"

// SpringServeHandler parameters.
type SpringServeHandler struct {
	//logger service
	logger log.Service
	//kafkaClient service
	kafkaClient kafkaclient.Service
	//demandClient service
	demandClient demand.Service
	//privacyClient service
	privacyClient privacy.Service
}

// NewHandler is a constructor function which get logger, demandClient implementation and
// kafkaClient service params arguments and return implementation of SpringServeHandler interface.
func SsNewHandler(l log.Service, k kafkaclient.Service, d demand.Service, p privacy.Service) *SpringServeHandler {
	return &SpringServeHandler{
		logger:        l,
		kafkaClient:   k,
		demandClient:  d,
		privacyClient: p,
	}
}

// Handler for processing incoming post request for springserve api
func (s *SpringServeHandler) Handler(ctx *fasthttp.RequestCtx) {
	defer func() {
		if rV := recover(); rV != nil {
			s.logger.Debug("Recover, Panic Detected", log.Metadata{"recover": rV})
		}
	}()
	ctx.Response.Header.Set("Connection", "keep-alive")
	body, err := ctx.Request.BodyGunzip()
	if err != nil {
		body = ctx.PostBody()
		s.logger.Debug("HB Request Body Not GZipped")
	}

	if !isValidJSON(body) {
		code := fasthttp.StatusBadRequest
		ctx.Response.SetStatusCode(code)
		s.logger.Error("Invalid Request body", log.Metadata{"body": string(body)})
		fmt.Fprintf(ctx, `{"error": "invalid Request body"}`)
		return
	}

	if s.privacyClient.Optout(body) {
		body = s.privacyClient.SetOptout(body)
	}

	headers := ctx.Request.Header.Header()
	path := ctx.Path()
	URI := ctx.URI()
	qs := URI.QueryString()
	metadata := log.Metadata{
		"path":   string(path),
		"qs":     string(qs),
		"header": string(headers),
		"body":   string(body)}
	s.logger.Info("HB Request processing", metadata)

	go s.ProcessAdEventToKafka(body, path, qs, headers)
	go s.ProcessInvEventToKafka(body)

	demandParams := ssDemandParams(body)
	code, respBody := s.demandClient.BidOrtbReq(demandParams)
	ctx.SetBody(respBody)
	ctx.SetStatusCode(code)
	ctx.Response.Header.Set("Content-Type", "application/json")
}

func ssDemandParams(json []byte) demand.DemandExtParams {
	return demand.DemandExtParams{
		Body:     json,
		SamFlag:  gjson.GetBytes(json, SAM_FLAGS_SS_PATH).String(),
		SamHbTag: gjson.GetBytes(json, SAM_HB_TAG_SS_PATH).String(),
	}
}

// ProcessAdEventToKafka publishes asynchronous Ad events.
func (s *SpringServeHandler) ProcessAdEventToKafka(jsonBody []byte, path []byte, qs []byte, headers []byte) {
	jsonData, err := frameAdEvent(jsonBody, path, qs, headers)
	if err != nil {
		s.logger.Error("FrameAdEvent Error", log.Metadata{"error": err})
		return
	}

	s.kafkaClient.PublishAdEventAsync(jsonData)
}

// ProcessAdEventToKafka publishes asynchronous Inventtory events.
func (s *SpringServeHandler) ProcessInvEventToKafka(jsonBody []byte) {
	jsonData, err := frameSSInvEvent(jsonBody)
	if err != nil {
		s.logger.Error("FrameAdEvent Error", log.Metadata{"error": err})
		return
	}
	s.kafkaClient.PublishInventoryEventAsync(jsonData)
}

// checks for invalid json and returns bool
func isValidJSON(body []byte) bool {
	var jsn interface{}

	return json.Unmarshal(body, &jsn) == nil
}
