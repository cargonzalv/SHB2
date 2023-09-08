package udws3

import "github.com/aws/aws-sdk-go-v2/service/s3/types"

// `interface.go` module contains declaration of interfaces,
// provided by `udws3` package
type (
	Service interface {
		FetchGzFiles() ([]types.Object, error)
		DownloadGzFile(string, string) (bool, error)
	}
)
