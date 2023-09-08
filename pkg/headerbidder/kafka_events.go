package headerbidder

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"strconv"
	"time"
)

// HbRequest Params
type AdEvent struct {
	// HbRequest Body
	Body string `json:"body"`
	// Type
	Type string `json:"type"`
	// Version
	Ver string `json:"version"`
	// Http request Headers
	Headers string `json:"headers"`
	// Query string
	Qs string `json:"qs"`
	// Request Api Path
	Path string `json:"path"`
}

type InventoryEvent struct {
	// Bid Request Id
	BidRequestId string `json:"bidRequestId"`
	// Timestamp of the Bid request
	Timestamp string `json:"timestamp"`
	// Samsung sessio id in uuid format
	SamsungSessionId string `json:"samsungSessionId"`
	// Samsung App Id
	SamsungAppId string `json:"samsungAppId"`
	// Samsung Hb Tag or Site id
	SamsungHbTag int `json:"samsungHbTag"`
	// Publica Tag Id
	PublicaTagId int `json:"publicaTagId"`
	// Content Genre Typee
	ContentGenre string `json:"contentGenre"`
	// Content Cat Iab
	ContentCatIab string `json:"contentCatIab"`
	// Content Live Stream
	IsLive int `json:"isLive"`
	// Device Ip from the heeaders
	Ip string `json:"ip"`
	// Device UA from the headers
	Ua string `json:"ua"`
	// Device Id
	DeviceId string `json:"deviceId"`
	// Device Lattitide
	Lat int `json:"lat"`
	// Device Coppa
	Coppa int `json:"coppa"`
	// GDPR
	Gdpr int `json:"gdpr"`
	// Us Privacy
	UsPrivacy string `json:"usPrivacy"`
	// Conseent for GDPR
	GdprConsent string `json:"gdprConsent"`
	// App Bundle
	AppBundle string `json:"appBundle"`
	// App Name
	AppName string `json:"appName"`
	// Request Slots
	RequestSlots int `json:"requestSlots"`
	// Request Durattion
	RequestDuration int `json:"requestDuration"`
}

// frameAdEvent frames the Ad event by taking body, path, qs and headers as arguments
func frameAdEvent(jsonBody []byte, path []byte, qs []byte, headers []byte) ([]byte, error) {
	var event AdEvent
	event.Body = string(jsonBody)

	event.Headers = string(headers)
	event.Path = string(path)
	event.Qs = string(qs)
	json_data, err := json.Marshal(&event)
	return json_data, err
}

// frameSSInvEvent frames the Inventory event by taking body as arguments
func frameSSInvEvent(jsonBody []byte) ([]byte, error) {
	var event InventoryEvent
	event.SamsungSessionId = gjson.GetBytes(jsonBody, "ext.sam_session_id").String()
	event.SamsungAppId = gjson.GetBytes(jsonBody, "ext.sam_app_id").String()
	event.SamsungHbTag = int(gjson.GetBytes(jsonBody, "ext.sam_hb_tag").Int())
	event.PublicaTagId = int(gjson.GetBytes(jsonBody, "ext.tag_id").Int())
	event.ContentGenre = gjson.GetBytes(jsonBody, "content.genre").String()
	event.ContentCatIab = gjson.GetBytes(jsonBody, "content.cat").String()
	event.IsLive = int(gjson.GetBytes(jsonBody, "content.livestream").Int())
	event.UsPrivacy = gjson.GetBytes(jsonBody, "regs.ext.us_privacy").String()
	frameInvEvent(jsonBody, &event)

	json_data, err := json.Marshal(&event)
	return json_data, err
}

// frameFWInvEvent frames the Inventory event by taking body as arguments
func frameFWInvEvent(jsonBody []byte) ([]byte, error) {
	var event InventoryEvent
	event.SamsungSessionId = gjson.GetBytes(jsonBody, "imp.0.pmp.deals.0.ext.sam_session_id").String()
	event.SamsungAppId = gjson.GetBytes(jsonBody, "imp.0.pmp.deals.0.ext.sam_app_id").String()
	event.SamsungHbTag = int(gjson.GetBytes(jsonBody, "imp.0.pmp.deals.0.ext.sam_hb_tag").Int())
	event.PublicaTagId = int(gjson.GetBytes(jsonBody, "imp.0.pmp.deals.0.ext.sam_hb_tag").Int())
	event.ContentGenre = gjson.GetBytes(jsonBody, "site.content.genre").String()
	event.ContentCatIab = gjson.GetBytes(jsonBody, "site.content.cat").String()
	event.IsLive = int(gjson.GetBytes(jsonBody, "site.content.livestream").Int())
	event.UsPrivacy = gjson.GetBytes(jsonBody, "regs.us_privacy").String()
	frameInvEvent(jsonBody, &event)

	json_data, err := json.Marshal(&event)
	return json_data, err
}

func frameInvEvent(jsonBody []byte, event *InventoryEvent) {
	event.BidRequestId = gjson.GetBytes(jsonBody, "id").String()
	event.Timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	event.Ip = gjson.GetBytes(jsonBody, "device.ip").String()
	event.Ua = gjson.GetBytes(jsonBody, "device.ua").String()
	event.DeviceId = gjson.GetBytes(jsonBody, "device.ifa").String()
	event.Lat = int(gjson.GetBytes(jsonBody, "device.lmt").Int())
	event.Coppa = int(gjson.GetBytes(jsonBody, "regs.coppa").Int())
	event.Gdpr = int(gjson.GetBytes(jsonBody, "regs.ext.gdpr").Int())
	event.GdprConsent = gjson.GetBytes(jsonBody, "user.ext.consent").String()
	event.AppBundle = gjson.GetBytes(jsonBody, "app.bundle").String()
	event.AppName = gjson.GetBytes(jsonBody, "app.name").String()
	event.RequestSlots = int(gjson.GetBytes(jsonBody, "imp.#").Int())
	var sum = 0
	count := int(gjson.Get(gjson.GetBytes(jsonBody, "imp.#.video").String(), "#").Int())
	for i := 0; i <= count; i++ {
		sum = sum + int(gjson.Get(gjson.GetBytes(jsonBody, "imp.#.video").String(), strconv.Itoa(i)+".maxduration").Int())
	}
	event.RequestDuration = sum
}
