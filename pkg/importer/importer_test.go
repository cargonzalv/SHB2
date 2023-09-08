package importer

import (
	"os"
	"testing"

	"github.com/adgear/go-commons/pkg/log"
	"github.com/adgear/sps-header-bidder/config"
	"github.com/adgear/sps-header-bidder/pkg/mockservices"
	"github.com/adgear/sps-header-bidder/pkg/tifascache"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestLoadTifas(t *testing.T) {
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

	mockUdwS3Client.EXPECT().FetchGzFiles().Return(contents, nil).AnyTimes()
	mockUdwS3Client.EXPECT().DownloadGzFile(gomock.Any(), gomock.Any()).Return(true, nil).AnyTimes()

	importerClient := NewService(logger, cfg, mockUdwS3Client, cacheClient)
	importerClient.LoadTifas()

	output1 := cacheClient.GetTifa("4e5e288c-284d-5187-0b46-76fcdd8d148e")
	output2 := cacheClient.GetTifa("970baefa-e68c-c43f-036a-fadf2b9a910a")
	output3 := cacheClient.GetTifa("e88821e4-1bbe-f150-b05d-cb95b4cb5c2e")

	assert.Equal(t, output1, true)
	assert.Equal(t, output2, true)
	assert.Equal(t, output3, true)

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
