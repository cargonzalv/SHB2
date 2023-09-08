// udws3 package.
package udws3

import (
	"context"
	"os"

	"github.com/adgear/go-commons/pkg/log"
	"github.com/adgear/sps-header-bidder/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	s3config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// Udws3 parameters.
type Udws3 struct {
	logger   log.Service
	s3Client *s3.Client
	config   *config.Config
}

// This statement forcing the module to
// implements the Udws3 `Service` interface.
var _ Service = (*Udws3)(nil)

// NewService is a constructor function which get logger implementation and
// config params arguments and return implementation of Udws3 interface.
func NewService(l log.Service, cfg *config.Config) Service {
	// Adding the default config with region
	assumecnf, _ := s3config.LoadDefaultConfig(context.TODO(), s3config.WithRegion(cfg.S3.Region))
	// Create the credentials from AssumeRoleProvider to assume the role
	// referenced by the "myRoleARN" ARN.
	stsSvc := sts.NewFromConfig(assumecnf)
	creds := stscreds.NewAssumeRoleProvider(stsSvc, cfg.S3.AssumeRole)
	assumecnf.Credentials = aws.NewCredentialsCache(creds)
	// Create service client value configured for credentials
	// from assumed role.
	s3Client := s3.NewFromConfig(assumecnf)
	return &Udws3{
		logger:   l,
		s3Client: s3Client,
		config:   cfg,
	}
}

// FetchGzFiles fetches the list of compressed GZ files from s3 bucket
func (u *Udws3) FetchGzFiles() ([]types.Object, error) {
	results, err := u.s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(u.config.S3.Bucket),
		Prefix: aws.String(u.config.S3.OptoutPrefix),
	})
	var contents []types.Object
	if err != nil {
		u.logger.Error("Cannot fetch the Compliance files from bucket", log.Metadata{"bucketName": u.config.S3.Bucket, "error": err})
	} else {
		u.logger.Info("List objects in bucket", log.Metadata{"bucketName": u.config.S3.Bucket, "Contents": results.Contents})
		contents = results.Contents
	}
	return contents, err
}

// DownloadGzFile downloads the GZfile with bucket key and download path
func (u *Udws3) DownloadGzFile(key string, filePath string) (bool, error) {
	downloadFile, err := os.Create(filePath)
	if err != nil {
		u.logger.Error("Error unable to open file ", log.Metadata{"error": err, "filePath": filePath})
		return false, err
	}

	defer downloadFile.Close()

	downloader := manager.NewDownloader(u.s3Client)
	numBytes, err := downloader.Download(context.TODO(), downloadFile, &s3.GetObjectInput{
		Bucket: aws.String(u.config.S3.Bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		u.logger.Error("Error unable to download file", log.Metadata{"error": err, "filePath": filePath})
		return false, err
	}
	log.Info("Download file ", log.Metadata{"numBytes": numBytes, "filePath": filePath})

	return true, nil
}
