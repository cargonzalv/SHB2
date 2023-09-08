package privacy

import (
	"os"
	"testing"

	"github.com/adgear/go-commons/pkg/log"
	"github.com/adgear/sps-header-bidder/config"
	"github.com/adgear/sps-header-bidder/pkg/importer"
	"github.com/adgear/sps-header-bidder/pkg/mockservices"
	"github.com/adgear/sps-header-bidder/pkg/tifascache"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

const hbReqBody = `{
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
        "ifa": "4e5e288c-284d-5187-0b46-76fcdd8d148e"
    },
    "user": {
        "id": "575883e2-8fb1-4147-b6d3-0a5896b94e54"
    },
    "at": 1,
    "tmax": 2800,
	"regs": {
        "ext": {
            "us_privacy": "1NNN"
        }
    },
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

func TestUsPrivacyOptoutValues(t *testing.T) {
	cfg, _ := config.NewConfig("../../config/config.yml")

	logger := log.GetLogger()
	cacheClient := tifascache.NewService(logger, cfg)
	privacyClient := NewService(logger, cacheClient)

	hbReqBodyWithUsPrivacyOptout, _ := sjson.Set(hbReqBody, "regs.ext.us_privacy", "1YYY")
	hbReqBodyWithoutUsPrivacyOptout, _ := sjson.Set(hbReqBody, "regs.ext.us_privacy", "1YNY")
	optout1 := privacyClient.Optout([]byte(hbReqBodyWithUsPrivacyOptout))
	optout2 := privacyClient.Optout([]byte(hbReqBodyWithoutUsPrivacyOptout))

	assert.Equal(t, optout1, true)
	assert.Equal(t, optout2, false)
}

func TestCoppaOptoutValues(t *testing.T) {
	cfg, _ := config.NewConfig("../../config/config.yml")

	logger := log.GetLogger()
	cacheClient := tifascache.NewService(logger, cfg)
	privacyClient := NewService(logger, cacheClient)

	hbReqBodyWithCoppaEnabled, _ := sjson.Set(hbReqBody, "regs.coppa", "1")
	hbReqBodyWithCoppadisabled, _ := sjson.Set(hbReqBody, "regs.coppa", "0")
	optout1 := privacyClient.Optout([]byte(hbReqBodyWithCoppaEnabled))
	optout2 := privacyClient.Optout([]byte(hbReqBodyWithCoppadisabled))

	assert.Equal(t, optout1, true)
	assert.Equal(t, optout2, false)
}

func TestGdprOptoutValues(t *testing.T) {
	cfg, _ := config.NewConfig("../../config/config.yml")

	logger := log.GetLogger()
	cacheClient := tifascache.NewService(logger, cfg)
	privacyClient := NewService(logger, cacheClient)

	hbReqBodyWithGdprEnabled, _ := sjson.Set(hbReqBody, "regs.ext.gdpr", "1")
	hbReqBodyWithGdprDisabled, _ := sjson.Set(hbReqBody, "regs.ext.gdpr", "0")
	hbReqBodyWithGdprEnabledWith0Consent, _ := sjson.Set(hbReqBodyWithGdprEnabled, "user.ext.consent", "0")
	hbReqBodyWithGdprEnabledWith1Consent, _ := sjson.Set(hbReqBodyWithGdprEnabled, "user.ext.consent", "1")

	optout0 := privacyClient.Optout([]byte(hbReqBodyWithGdprEnabledWith0Consent))
	optout1 := privacyClient.Optout([]byte(hbReqBodyWithGdprEnabledWith1Consent))
	optout2 := privacyClient.Optout([]byte(hbReqBodyWithGdprDisabled))

	assert.Equal(t, optout0, false)
	assert.Equal(t, optout1, true)
	assert.Equal(t, optout2, false)
}

func TestOptoutCcpaTifaFromDbExists(t *testing.T) {
	cfg, _ := config.NewConfig("../../config/config.yml")
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	logger := log.GetLogger()

	sourceSuccessFilePath := "../../testing/importer/_SUCCESS"
	sourceFilePath := "../../testing/importer/compliance.json.gz"
	successFile := "compliance_opt_out_singles_v1__SUCCESS"
	complianceGzfile := "opt_out_singles_compliance.json.gz"

	destFilePath := "/tmp/" + complianceGzfile
	destSucessFilePath := "/tmp/" + successFile
	createMockDownloadFile(sourceSuccessFilePath, destSucessFilePath)
	createMockDownloadFile(sourceFilePath, destFilePath)

	var content types.Object
	content.Key = &complianceGzfile
	contents := make([]types.Object, 0, 1)
	contents = append(contents, content)

	mockUdwS3Client := mockservices.NewMockudws3Client(mockCtrl)
	cacheClient := tifascache.NewService(logger, cfg)
	privacyClient := NewService(logger, cacheClient)
	mockUdwS3Client.EXPECT().FetchGzFiles().Return(contents, nil).AnyTimes()
	mockUdwS3Client.EXPECT().DownloadGzFile(gomock.Any(), gomock.Any()).Return(true, nil).AnyTimes()

	importerClient := importer.NewService(logger, cfg, mockUdwS3Client, cacheClient)
	importerClient.LoadTifas()

	optout := privacyClient.Optout([]byte(hbReqBody))
	redactedHbReqBody := privacyClient.SetOptout([]byte(hbReqBody))
	redactedTifa := gjson.GetBytes(redactedHbReqBody, "device.ifa").String()
	redactedIp := gjson.GetBytes(redactedHbReqBody, "device.ip").String()
	redactedLocation := gjson.GetBytes(redactedHbReqBody, "device.geo").String()
	usPrivacy := gjson.GetBytes(redactedHbReqBody, "regs.ext.us_privacy").String()

	hbReqBodyWithTestTifa, _ := sjson.Set(hbReqBody, "device.ifa", "4e5e288c-284d-5187-0b46-76fcdd8dtest")

	optout2 := privacyClient.Optout([]byte(hbReqBodyWithTestTifa))

	assert.Equal(t, optout, true)
	assert.Equal(t, optout2, false)
	assert.Equal(t, redactedTifa == REDEACTED_STR, true)
	assert.Equal(t, redactedIp == REDEACTED_STR, true)
	assert.Equal(t, REDEACTED_GEO_STR == redactedLocation, true)
	assert.Equal(t, usPrivacy == USER_CONSENT_ENABLED, true)
}
func TestOptoutCcpaTifaFromDbDoesNotExists(t *testing.T) {
	cfg, _ := config.NewConfig("../../config/config.yml")
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logger := log.GetLogger()
	cacheClient := tifascache.NewService(logger, cfg)
	privacyClient := NewService(logger, cacheClient)

	optout := privacyClient.Optout([]byte(hbReqBody))
	assert.Equal(t, optout, false)
}

func createMockDownloadFile(sourceFilePath string, destFilePath string) {
	logger := log.GetLogger()
	input, err := os.ReadFile(sourceFilePath)
	if err != nil {
		logger.Error("Read Error", log.Metadata{"err": err})
		return
	}

	err = os.WriteFile(destFilePath, input, 0644)
	if err != nil {
		logger.Error("Write Error", log.Metadata{"err": err})
		return
	}
}
