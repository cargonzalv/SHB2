package demand_test

import (
	"errors"
	"testing"

	"github.com/adgear/go-commons/pkg/log"

	"github.com/adgear/go-commons/pkg/httpclient"
	"github.com/adgear/sps-header-bidder/pkg/demand"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const hbReqBodyWithoutSamFlag = `{
    "id": "24234dfa-ef05-45fc-9dd8-7f37898294b6",
    "imp": [
        {
            "id": "1",
            "video": {
                "mimes": [
                    "video/mp4",
                    "video/ogg",
                    "video/webm"
                ],
                "minduration": 1,
                "maxduration": 300,
                "protocols": [
                    1,
                    2,
                    3,
                    4,
                    5,
                    6
                ],
                "w": 1920,
                "h": 1080,
                "startdelay": -1,
                "linearity": 1,
                "sequence": 1,
                "minbitrate": 1,
                "maxbitrate": 280000
            },
            "tagid": "666310",
            "bidfloor": 0.01,
            "secure": 1
        }
    ],
    "app": {
        "id": "666310",
        "name": "test_app",
        "bundle": "test",
        "domain": "test.com",
        "publisher": {
            "id": "1681"
        }
    },
    "device": {
        "ua": "Mozilla/5.0 (SMART-TV; LINUX; Tizen 3.0) AppleWebKit/538.1 (KHTML, like Gecko) Version/3.0 TV Safari/538.1",
        "geo": {
            "lat": 37.751,
            "lon": -97.822,
            "type": 2,
            "country": "USA",
            "ipservice": 3
        },
        "dnt": 0,
        "lmt": 0,
        "ip": "11.22.33.44",
        "devicetype": 3,
        "make": "Samsung",
        "model": "Tizen TV 2017",
        "os": "Tizen",
        "language": "en",
        "ifa": "11112222-3333-4444-5555-666677778888"
    },
    "user": {
        "id": "575883e2-8fb1-4147-b6d3-0a5896b94e54"
    },
    "at": 1,
    "tmax": 2800,
    "source": {
        "fd": 0,
        "pchain": ""
    },
    "ext": {
        "sam_hb_tag": "1234",
        "sam_app_id": "TEST-App-Id01",
        "sam_session_id": "888888e2-621f-4f6a-844c-d21aed450c2a"
    }
}`

const hbReqBodyWithSamFlag = `{
	"id": "24234dfa-ef05-45fc-9dd8-7f37898294b6",
	"imp": [
	  {
		"id": "1",
		"video": {
		  "mimes": [
			"video/mp4",
			"video/ogg",
			"video/webm"
		  ],
		  "minduration": 1,
		  "maxduration": 300,
		  "protocols": [
			1,
			2,
			3,
			4,
			5,
			6
		  ],
		  "w": 1920,
		  "h": 1080,
		  "startdelay": -1,
		  "linearity": 1,
		  "sequence": 1,
		  "minbitrate": 1,
		  "maxbitrate": 280000
		},
		"tagid": "666310",
		"bidfloor": 0.01,
		"secure": 1
	  }
	],
	"app": {
	  "id": "666310",
	  "name": "test_app",
	  "bundle": "test",
	  "domain": "test.com",
	  "publisher": {
		"id": "1681"
	  }
	},
	"device": {
	  "ua": "Mozilla/5.0 (SMART-TV; LINUX; Tizen 3.0) AppleWebKit/538.1 (KHTML, like Gecko) Version/3.0 TV Safari/538.1",
	  "geo": {
		"lat": 37.751,
		"lon": -97.822,
		"type": 2,
		"country": "USA",
		"ipservice": 3
	  },
	  "dnt": 0,
	  "lmt": 0,
	  "ip": "11.22.33.44",
	  "devicetype": 3,
	  "make": "Samsung",
	  "model": "Tizen TV 2017",
	  "os": "Tizen",
	  "language": "en",
	  "ifa": "11112222-3333-4444-5555-666677778888"
	},
	"user": {
	  "id": "575883e2-8fb1-4147-b6d3-0a5896b94e54"
	},
	"at": 1,
	"tmax": 2800,
	"source": {
	  "fd": 0,
	  "pchain": ""
	},
	"ext": {
	  "sam_flags": "pbid",
	  "sam_hb_tag": "1234",
	  "sam_app_id": "TEST-App-Id01",
	  "sam_session_id": "888888e2-621f-4f6a-844c-d21aed450c2a"
	}
} `

const hbReqBodyWithoutHbTag = `{
	"id": "24234dfa-ef05-45fc-9dd8-7f37898294b6",
	"imp": [
	  {
		"id": "1",
		"video": {
		  "mimes": [
			"video/mp4",
			"video/ogg",
			"video/webm"
		  ],
		  "minduration": 1,
		  "maxduration": 300,
		  "protocols": [
			1,
			2,
			3,
			4,
			5,
			6
		  ],
		  "w": 1920,
		  "h": 1080,
		  "startdelay": -1,
		  "linearity": 1,
		  "sequence": 1,
		  "minbitrate": 1,
		  "maxbitrate": 280000
		},
		"tagid": "666310",
		"bidfloor": 0.01,
		"secure": 1
	  }
	],
	"app": {
	  "id": "666310",
	  "name": "test_app",
	  "bundle": "test",
	  "domain": "test.com",
	  "publisher": {
		"id": "1681"
	  }
	},
	"device": {
	  "ua": "Mozilla/5.0 (SMART-TV; LINUX; Tizen 3.0) AppleWebKit/538.1 (KHTML, like Gecko) Version/3.0 TV Safari/538.1",
	  "geo": {
		"lat": 37.751,
		"lon": -97.822,
		"type": 2,
		"country": "USA",
		"ipservice": 3
	  },
	  "dnt": 0,
	  "lmt": 0,
	  "ip": "11.22.33.44",
	  "devicetype": 3,
	  "make": "Samsung",
	  "model": "Tizen TV 2017",
	  "os": "Tizen",
	  "language": "en",
	  "ifa": "11112222-3333-4444-5555-666677778888"
	},
	"user": {
	  "id": "575883e2-8fb1-4147-b6d3-0a5896b94e54"
	},
	"at": 1,
	"tmax": 2800,
	"source": {
	  "fd": 0,
	  "pchain": ""
	},
	"ext": {
	  "sam_flags": "pbid",
	  "sam_app_id": "TEST-App-Id01",
	  "sam_session_id": "888888e2-621f-4f6a-844c-d21aed450c2a"
	}
} `

const publicaUrl string = "http://test_publica_url"
const validBody = "valid"
const invalidBody = "invalid"

// Happy path. Request with sam_flags
func TestRequestSamFlagHappyPath(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logger := log.GetLogger()
	mockhttpClient := httpclient.NewMockService(mockCtrl)
	demandClient := demand.New(logger, publicaUrl, mockhttpClient)
	response := &httpclient.Response{
		Headers:    httpclient.Headers{},
		Body:       []byte(validBody),
		StatusCode: 200,
	}

	mockhttpClient.EXPECT().Post(gomock.Any()).Return(response, nil)

	var demandParams demand.DemandExtParams
	demandParams.Body = []byte(hbReqBodyWithSamFlag)
	demandParams.SamFlag = "pbid"
	demandParams.SamHbTag = "1234"
	respCode, respBody := demandClient.BidOrtbReq(demandParams)
	strBody := string(respBody)
	assert.Equal(t, respCode, 200)
	assert.Equal(t, strBody, validBody)
}

// Happy path. Request without sam_flags
func TestRequestNoSamFlagHappyPath(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logger := log.GetLogger()
	mockhttpClient := httpclient.NewMockService(mockCtrl)
	demandClient := demand.New(logger, publicaUrl, mockhttpClient)

	var demandParams demand.DemandExtParams
	demandParams.Body = []byte(hbReqBodyWithoutSamFlag)
	demandParams.SamFlag = ""
	demandParams.SamHbTag = "1234"
	respCode, respBody := demandClient.BidOrtbReq(demandParams)
	strBody := string(respBody)
	assert.Equal(t, respCode, 204)
	assert.Equal(t, strBody, "")
}

// Request with sam_hb_tag_missing
func TestRequestNoSamHbTag(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logger := log.GetLogger()
	mockhttpClient := httpclient.NewMockService(mockCtrl)
	demandClient := demand.New(logger, publicaUrl, mockhttpClient)

	var demandParams demand.DemandExtParams
	demandParams.Body = []byte(hbReqBodyWithoutHbTag)
	demandParams.SamFlag = "pbid"
	demandParams.SamHbTag = ""
	respCode, respBody := demandClient.BidOrtbReq(demandParams)
	strBody := string(respBody)
	assert.Equal(t, respCode, 200)
	assert.Equal(t, strBody, "{\"error\": \"sam_hb_tag(site_id) is missing\"}")
}

// Http client error. Http error response
func TestRequestWhenHttpClientErrorResponse(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logger := log.GetLogger()
	mockhttpClient := httpclient.NewMockService(mockCtrl)
	demandClient := demand.New(logger, publicaUrl, mockhttpClient)

	errorResponse := &httpclient.Response{
		Headers:    httpclient.Headers{},
		Body:       []byte(invalidBody),
		StatusCode: 501,
	}
	err := errors.New("http error")
	mockhttpClient.EXPECT().Post(gomock.Any()).Return(errorResponse, err)

	var demandParams demand.DemandExtParams
	demandParams.Body = []byte(hbReqBodyWithSamFlag)
	demandParams.SamFlag = "pbid"
	demandParams.SamHbTag = "1234"
	respCode, respBody := demandClient.BidOrtbReq(demandParams)
	assert.Equal(t, respCode, 501)
	strBody := string(respBody)
	assert.Equal(t, strBody, invalidBody)
}

// Http client error. Http client error
func TestRequestWhenHttpClientError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logger := log.GetLogger()
	mockhttpClient := httpclient.NewMockService(mockCtrl)
	demandClient := demand.New(logger, publicaUrl, mockhttpClient)

	err := errors.New("http client error")
	mockhttpClient.EXPECT().Post(gomock.Any()).Return(nil, err)

	var demandParams demand.DemandExtParams
	demandParams.Body = []byte(hbReqBodyWithSamFlag)
	demandParams.SamFlag = "pbid"
	demandParams.SamHbTag = "1234"
	respCode, respBody := demandClient.BidOrtbReq(demandParams)

	assert.Equal(t, respCode, 500)
	strBody := string(respBody)
	assert.Equal(t, strBody, "Internal error")
}
