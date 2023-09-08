// privacy package.
package privacy

import (
	"encoding/base64"
	"strconv"

	"github.com/adgear/go-commons/pkg/log"
	"github.com/adgear/go-commons/pkg/metric"
	"github.com/adgear/sps-header-bidder/pkg/tifascache"
	"github.com/prebid/go-gdpr/api"
	"github.com/prebid/go-gdpr/vendorconsent"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

const REDEACTED_GEO_STR string = `{"lat":0.1,"lon":-0.1,"country":"*REDACTED*","region":"*REDACTED*","metro":"*REDACTED*","city":"*REDACTED*","zip":"*REDACTED*"}`
const REDEACTED_STR string = "*REDACTED*"
const USER_CONSENT_ENABLED string = "1-Y-"

// Privacy parameters.
type Privacy struct {
	//logger service
	logger log.Service
	//tifascache service
	cache tifascache.Service
}

// This statement forcing the module to
// implements the Privacy `Service` interface.
var _ Service = (*Privacy)(nil)

// NewService is a constructor function which get logger implementation and
// tifascache params arguments and return implementation of Privacy interface.
func NewService(l log.Service, c tifascache.Service) Service {
	return &Privacy{
		logger: l,
		cache:  c,
	}
}

// Optout function takes requestBody as input and checks the various optout params like
// us_privacy, coppa, gdpr and checks the device id against the tifas stored in the tifascache
// and returns true/false based on the conditions of various optouts
func (p *Privacy) Optout(body []byte) bool {
	usPrivacy := gjson.GetBytes(body, "regs.ext.us_privacy").String()
	coppa := gjson.GetBytes(body, "regs.coppa").String()
	gdpr := gjson.GetBytes(body, "regs.ext.gdpr").String()
	gdprConsent := gjson.GetBytes(body, "user.ext.consent").String()
	tifa := gjson.GetBytes(body, "device.ifa").String()
	return p.ccpa_tifa(tifa) || p.ccpa(usPrivacy) || p.coppa(coppa) || p.gdpr(gdpr, gdprConsent)
}

// SetOptout function takes requestBody as input and sets the requestBody with redacted string to deviceid,
// ip and the geo params and sets the us_privacy field with user_consent yes with 1-Y-
func (p *Privacy) SetOptout(body []byte) []byte {
	jsonBody := string(body)
	jsonBody, _ = sjson.Set(jsonBody, "device.ifa", REDEACTED_STR)
	jsonBody, _ = sjson.Set(jsonBody, "device.ip", REDEACTED_STR)
	jsonBody, _ = sjson.Set(jsonBody, "device.geo", REDEACTED_GEO_STR)
	jsonBody, _ = sjson.Set(jsonBody, "regs.ext.us_privacy", USER_CONSENT_ENABLED)
	return []byte(jsonBody)
}

// ccpa function checks the user consent and returns bool
func (p *Privacy) ccpa(usPrivacy string) bool {
	var optoutResult bool
	if usPrivacy == "" {
		optoutResult = false
	} else {
		userOptout := usPrivacy[2:3]
		switch userOptout {
		case "y":
			optoutResult = true
		case "Y":
			optoutResult = true
		default:
			optoutResult = false
		}
	}
	addMetrics("usprivacy", strconv.FormatBool(optoutResult))
	return optoutResult
}

// ccpa_tifa function will check the device ids against the tifas in the
// tifascache and returns true if exists
func (p *Privacy) ccpa_tifa(tifa string) bool {
	exists := p.cache.GetTifa(tifa)
	p.logger.Info("checking tifa in Cache", log.Metadata{"exists": exists})
	addMetrics("ccpatifa", strconv.FormatBool(exists))
	return exists
}

// coppa function checks the coppa flag and returns bool
func (p *Privacy) coppa(coppa string) bool {
	optOutResult := coppa == "1"
	addMetrics("coppa", strconv.FormatBool(optOutResult))
	return optOutResult
}

// gdpr function checks the gdpr flag and returns bool
// calls gdpr_consent if the gdpr is 1
func (p *Privacy) gdpr(gdpr string, consent string) bool {
	var optoutResult = false
	if gdpr == "1" {
		optoutResult = p.gdpr_consent(consent)
	}
	addMetrics("gdpr", strconv.FormatBool(optoutResult))
	return optoutResult
}

// gdpr_consent function checks if the gdpr_consent str is other than empty,0,1
// gdpr_consent func decodes the string and checks against vendor consent
func (p *Privacy) gdpr_consent(consentStr string) bool {
	if consentStr == "" {
		return false
	}
	if consentStr == "0" {
		return false
	}
	if consentStr == "1" {
		return true
	}
	data, _ := base64.RawURLEncoding.DecodeString(consentStr)
	consent, err := vendorconsent.Parse(data)
	if err != nil {
		p.logger.Info("Data was not a valid consent string",
			log.Metadata{"consent": consent, "error": err})
		//  to be safe, remove pii if can't parse consent string,
		//  same as rtb_delivery_privacy
		return true
	}
	if consent != nil {
		return p.gdpr_optout(consent)
	}
	return false

}

// temp impl before clarifying how to identify vis optout.
func (p *Privacy) gdpr_optout(consent api.VendorConsents) bool {
	return true
}

// addMetrics add the count to display on the grafana panel
func addMetrics(optoutType string, optoutResult string) {
	lvs := map[string]string{
		"type":   optoutType,
		"optout": optoutResult,
	}
	metric.Incr("sps_privacy_optouts_total", lvs)
}
