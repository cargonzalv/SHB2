package headerbidder

import (
	"fmt"

	"github.com/adgear/go-commons/pkg/log"
	"github.com/adgear/sps-header-bidder/pkg/demand"
	"github.com/adgear/sps-header-bidder/pkg/kafkaclient"
	"github.com/adgear/sps-header-bidder/pkg/privacy"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/valyala/fasthttp"
)

const SAM_FLAGS_FW_PATH string = "imp.0.pmp.deals.0.ext.sam_flags"
const SAM_HB_TAG_FW_PATH string = "imp.0.pmp.deals.0.ext.sam_hb_tag"
const DEAL_ID_PATH string = "imp.0.pmp.deals.0.id"
const DEALS_PATH string = "imp.0.pmp.deals.#"
const BIDS_PATH = "seatbid.0.bid.#"
const PUBLICA_DEAL_ID_PATH string = "seatbid.0.bid.0.dealid"

// FreeWheelHandler parameters.
type FreeWheelHandler struct {
	//logger service
	logger log.Service
	//kafkaClient service
	kafkaClient kafkaclient.Service
	//demandClient service
	demandClient demand.Service
	//privacyClient service
	privacyClient privacy.Service
}

// FwNewHandler is a constructor function which get logger, demandClient implementation and
// kafkaClient service params arguments and return implementation of FreeWheelHandler interface.
func FwNewHandler(l log.Service, k kafkaclient.Service, d demand.Service, p privacy.Service) *FreeWheelHandler {
	return &FreeWheelHandler{
		logger:        l,
		kafkaClient:   k,
		demandClient:  d,
		privacyClient: p,
	}
}

// Handler for processing incoming post request for freewheel api
func (f *FreeWheelHandler) Handler(ctx *fasthttp.RequestCtx) {
	defer func() {
		if recover := recover(); recover != nil {
			f.logger.Debug("Recover, Panic Detected", log.Metadata{"recover": recover})
		}
	}()
	ctx.Response.Header.Set("Connection", "keep-alive")
	body, err := ctx.Request.BodyGunzip()
	if err != nil {
		body = ctx.PostBody()
		f.logger.Debug("FW Request Body Not GZipped")
	}
	jsonBody := string(body)

	if !isValidJSON(body) {
		code := fasthttp.StatusBadRequest
		ctx.Response.SetStatusCode(code)
		f.logger.Error("Invalid FW Request body", log.Metadata{"body": jsonBody})
		fmt.Fprintf(ctx, `{"error": "invalid FW Request body"}`)
		return
	}

	if f.privacyClient.Optout(body) {
		body = f.privacyClient.SetOptout(body)
	}

	headers := ctx.Request.Header.Header()
	path := ctx.Path()
	URI := ctx.URI()
	qs := URI.QueryString()
	metadata := log.Metadata{
		"path":   string(path),
		"qs":     string(qs),
		"header": string(headers),
		"body":   jsonBody}
	f.logger.Info("FW Request processing", metadata)

	go f.ProcessAdEventToKafka(body, path, qs, headers)
	go f.ProcessInvEventToKafka(body)

	demandParams, dealId := fwDemandParams(body)
	code, respBody := f.demandClient.BidOrtbReq(demandParams)
	updatedBody := updateResponseDemandId(respBody, dealId)
	ctx.SetBody(updatedBody)
	ctx.SetStatusCode(code)
	ctx.Response.Header.Set("Content-Type", "application/json")
}

func updateResponseDemandId(json []byte, dealId string) []byte {
	if gjson.GetBytes(json, BIDS_PATH).Num < 1 {
		// if we have got unexpected or invalid response from Publica,
		// return response as is
		return json
	}

	modifiedJson, err := sjson.SetBytes(json, PUBLICA_DEAL_ID_PATH, dealId)

	if err != nil {
		return json
	}

	return modifiedJson
}

func fwDemandParams(json []byte) (demand.DemandExtParams, string) {
	if gjson.GetBytes(json, DEALS_PATH).Num < 1 {
		return demand.EmptyDemandParams(), ""
	}

	dealId := gjson.GetBytes(json, DEAL_ID_PATH).String()

	if len(dealId) == 0 {
		return demand.EmptyDemandParams(), ""
	}

	demandParams := demand.DemandExtParams{
		Body:     json,
		SamFlag:  gjson.GetBytes(json, SAM_FLAGS_FW_PATH).String(),
		SamHbTag: gjson.GetBytes(json, SAM_HB_TAG_FW_PATH).String(),
	}

	return demandParams, dealId
}

// ProcessAdEventToKafka publishes asynchronous Ad events.
func (f *FreeWheelHandler) ProcessAdEventToKafka(jsonBody []byte, path []byte, qs []byte, headers []byte) {
	jsonData, err := frameAdEvent(jsonBody, path, qs, headers)
	if err != nil {
		f.logger.Error("FrameAdEvent Error", log.Metadata{"error": err})
		return
	}

	f.kafkaClient.PublishAdEventAsync(jsonData)
}

// ProcessInvEventToKafka publishes asynchronous Inventtory events.
func (f *FreeWheelHandler) ProcessInvEventToKafka(jsonBody []byte) {
	jsonData, err := frameFWInvEvent(jsonBody)
	if err != nil {
		f.logger.Error("FrameAdEvent Error", log.Metadata{"error": err})
		return
	}
	f.kafkaClient.PublishInventoryEventAsync(jsonData)
}
