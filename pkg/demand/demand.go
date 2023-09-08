// demand package.
package demand

import (
	"net/url"

	"github.com/adgear/go-commons/pkg/metric"

	"github.com/adgear/go-commons/pkg/httpclient"
	"github.com/adgear/go-commons/pkg/log"
)

const SAM_FLAGS_PATH string = "imp.0.pmp.deals.0.ext.sam_flags"
const SAM_HB_TAG_PATH string = "imp.0.pmp.deals.0.ext.sam_hb_tag"
const DEAL_ID_PATH string = "imp.0.pmp.deals.0.id"
const DEALS_PATH string = "imp.0.pmp.deals.#"
const PUBLICA_DEAL_ID_PATH string = "seatbid.0.bid.dealid"

// Demand parameters.
type Demand struct {
	//logger service
	logger log.Service
	//httpclient service
	httpClient httpclient.Service
	//Publica Demand Url
	publicaUrl string
}

// Demand Samsung Ext params
type DemandExtParams struct {
	// Sam pbid flag
	SamFlag string
	// Sam Tag id or site_id
	SamHbTag string
	// request body
	Body []byte
}

// This statement forcing the module to
// implements the Demand `Service` interface.
var _ Service = (*Demand)(nil)

// New is a constructor function which get logger implementation and
// httpclient params arguments and return implementation of Demand interface.
func New(l log.Service, publicaUrl string, httpClient httpclient.Service) *Demand {
	return &Demand{
		logger:     l,
		httpClient: httpClient,
		publicaUrl: publicaUrl,
	}
}

// EmptyDemandParams returns Empty/Init Demand Params
func EmptyDemandParams() DemandExtParams {
	return DemandExtParams{
		SamFlag:  "",
		SamHbTag: "",
		Body:     []byte{},
	}
}

// MakePublicaReq function takes body, extSamFlag, extSamHbTag and makes a publica request if the pbid fag is true
// returns the openrtb response body with 200 and if pbid flag is false returns 204 no content
func (d *Demand) BidOrtbReq(demandParams DemandExtParams) (respCode int, respBody []byte) {
	if demandParams.SamFlag != "pbid" {
		statusNoContent := 204
		return statusNoContent, []byte{}
	}

	lvs := map[string]string{
		"key":   "pbid",
		"value": "true",
	}

	if demandParams.SamHbTag == "" {
		statusMissingParam := 200
		return statusMissingParam, []byte(`{"error": "sam_hb_tag(site_id) is missing"}`)
	}

	metric.Incr("http_requests_sps_total", lvs)
	publicaUrl, _ := url.Parse(d.publicaUrl)
	values := publicaUrl.Query()
	values.Set("site_id", demandParams.SamHbTag)
	publicaUrl.RawQuery = values.Encode()
	req := &httpclient.Request{
		Url: publicaUrl.String(),
		Headers: httpclient.Headers{
			ContentType: []byte("application/json"),
		},
		Body: demandParams.Body,
	}
	resp, err := d.httpClient.Post(req)

	if err != nil {
		d.logger.Error("Error sending request to publica", log.Metadata{"url": d.publicaUrl, "error": err})
	}

	if resp == nil {
		return 500, []byte("Internal error")
	}

	d.logger.Debug("Successfully get resp from Publica", log.Metadata{"url": d.publicaUrl, "resp_status": resp.StatusCode, "resp_body": string(resp.Body)})

	return resp.StatusCode, resp.Body
}
