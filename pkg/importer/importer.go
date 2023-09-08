// nolint
package importer

import (
	"bufio"
	"compress/gzip"
	"context"
	"errors"
	"io"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/adgear/go-commons/pkg/log"
	"github.com/adgear/sps-header-bidder/config"
	"github.com/adgear/sps-header-bidder/pkg/tifascache"
	"github.com/adgear/sps-header-bidder/pkg/udws3"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/datadog/mmh3"
	"github.com/google/uuid"
	"github.com/tidwall/gjson"
)

type processed struct {
	numRows         int
	consentTypes    []string
	ids             []string
	idTypes         []string
	regulationTypes []string
}

// Importer parameters.
type Importer struct {
	//config service
	config *config.Config
	//logger service
	logger log.Service
	//udws3 service
	udws3 udws3.Service
	//tifascache service
	tifasCache tifascache.Service
}

// This statement forcing the module to
// implements the Importer `Service` interface.
var _ Service = (*Importer)(nil)

// NewService is a constructor function which get logger implementation,
// config params, udws3 params and tifascache params arguments
// and return implementation of Importer interface.
func NewService(l log.Service, cfg *config.Config, u udws3.Service, tc tifascache.Service) Service {

	return &Importer{
		config:     cfg,
		logger:     l,
		udws3:      u,
		tifasCache: tc,
	}
}

// LoadTifas func checks the success dump file if tifas dump expired/not
// if expires, then gets the dump from the s3 bucket and downloads to tmp folder
// decompresses it and reads the dump and loads the tifas to tifascache
// and deletes the downloaded files from tmp after completion
func (i *Importer) LoadTifas() {
	i.logger.Info("Load Tifas Initiated at ", log.Metadata{"time": time.Now().Format("01-02-2006 15:04:05 Monday")})
	updatedTimestamp, err := i.getSucessOptoutFileTimestamp(i.config.S3.OptputSuccessFileKey)
	if err != nil {
		i.logger.Error("Unable to process the SUCCESS file",
			log.Metadata{"error": err,
				"successfile": i.config.S3.OptputSuccessFileKey})
		return
	}
	if i.tifasCache.IsLastLoadTsExpired(updatedTimestamp) {
		i.tifasCache.SetTifa(i.config.Cache.LoadTimestampKey, updatedTimestamp, -1)
		gzFiles, err := i.udws3.FetchGzFiles()
		if err != nil {
			i.logger.Error("Error fetching compliance gz files from S3 bucket",
				log.Metadata{"error": err, "bucket": i.config.S3.Bucket})
			return
		}
		var wg sync.WaitGroup
		wg.Add(len(gzFiles))
		for _, gzFile := range gzFiles {
			go func(gzFile types.Object) {
				processed, _ := i.processTifasFromS3(aws.ToString(gzFile.Key))
				if err != nil {
					i.logger.Debug("Unable to process the gz file",
						log.Metadata{"processed": processed, "error": err,
							"gzfile": aws.ToString(gzFile.Key)})
				}
				wg.Done()
			}(gzFile)
		}
		wg.Wait()
	}
}

func (i *Importer) getSucessOptoutFileTimestamp(key string) (string, error) {
	fileName := strings.Replace(key, "/", "_", -1)
	downloadPath := "/tmp/" + fileName

	completed, err := i.udws3.DownloadGzFile(key, downloadPath)
	if completed {
		return i.processSuccessFile(downloadPath)
	}
	return "", err
}

func (i *Importer) processTifasFromS3(key string) (bool, error) {
	if !isOptoutFile(key) {
		return false, errors.New("non opt_out compliance file")
	}
	fileName := strings.Replace(key, "/", "_", -1)
	downloadPath := "/tmp/" + fileName

	completed, err := i.udws3.DownloadGzFile(key, downloadPath)
	if completed {
		return i.processGzFile(downloadPath)
	}
	return false, err
}

func (i *Importer) processGzFile(gzfile string) (bool, error) {
	_, err := i.extractGzFile(gzfile)
	if err != nil {
		i.logger.Error("Error unable to Decompress file ", log.Metadata{"error": err, "gzfile": gzfile})
		return false, err
	}
	jsonfilename := strings.TrimSuffix(gzfile, ".gz")
	_, err = i.processJsonFile(jsonfilename)
	if err != nil {
		i.logger.Error("Error unable to process Json file ", log.Metadata{"error": err, "filename": jsonfilename})
		return false, err
	}
	if err == nil {
		_ = os.Remove(gzfile)
		_ = os.Remove(jsonfilename)
	}
	return true, nil
}

func (i *Importer) extractGzFile(gzfile string) (bool, error) {
	gzipfile, err := os.Open(gzfile)
	if err != nil {
		i.logger.Error("Error unable to open gzfile file ", log.Metadata{"error": err, "filename": gzfile})
		return false, err
	}
	reader, err := gzip.NewReader(gzipfile)
	if err != nil {
		i.logger.Error("Error unable to read gzfile file ", log.Metadata{"error": err, "filename": gzfile})
		return false, err
	}
	defer reader.Close()

	newfilename := strings.TrimSuffix(gzfile, ".gz")
	writer, err := os.Create(newfilename)
	if err != nil {
		i.logger.Error("Error unable to write to json file ", log.Metadata{"error": err, "newfilename": newfilename})
		return false, err
	}
	defer writer.Close()

	if _, err = io.Copy(writer, reader); err != nil {
		i.logger.Error("Error unable to copy bytes to json file ", log.Metadata{"error": err, "newfilename": newfilename})
		return false, err
	}
	return true, nil

}

func (i *Importer) processSuccessFile(successFile string) (string, error) {
	var successTimestamp string
	f, err := os.Open(successFile)
	if err != nil {
		i.logger.Error("Error unable to open success file ", log.Metadata{"error": err, "successFile": successFile})
		return successTimestamp, err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		successTimestamp = gjson.GetBytes(scanner.Bytes(), "updated").String()
	}
	if err := scanner.Err(); err != nil {
		i.logger.Error("Unable to read success file", log.Metadata{"error": err, "successFile": successFile})
		return successTimestamp, err
	}
	i.logger.Info("Read success file", log.Metadata{"updated_timestamp": successTimestamp})
	return successTimestamp, nil
}

func (i *Importer) processJsonFile(jsonFile string) (bool, error) {
	f, err := os.Open(jsonFile)
	if err != nil {
		i.logger.Error("Error unable to open json file ", log.Metadata{"error": err, "jsonFile": jsonFile})
		return false, err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		i.processAndStore(scanner.Bytes())
	}
	if err := scanner.Err(); err != nil {
		i.logger.Error("Unable to read json file", log.Metadata{"error": err, "jsonFile": jsonFile})
		return false, err
	}
	return true, nil
}

func (i *Importer) processAndStore(data []byte) {
	consent_type := gjson.GetBytes(data, "consent_type").String()
	id := gjson.GetBytes(data, "id").String()
	id_type := gjson.GetBytes(data, "id_type").String()
	// r := gjson.GetBytes(data, "regulation_type").String()
	if (consent_type == "IBA") && ((id_type == "PSID") || (id_type == "TIFA")) {
		if id_type == "TIFA" {
			i.tifasCache.SetTifa(id, id, 0)
			return
		}
		tifa, err := psidToTifa(id)
		if err != nil {
			i.logger.Error("Unable to do psidToTifa conversion ", log.Metadata{"error": err, "psid": id})
		}
		i.tifasCache.SetTifa(tifa.String(), id, 0)
	}
}

func psidToTifa(psid string) (uuid.UUID, error) {
	key := strings.ToLower(psid)
	hash := mmh3.Hash128([]byte(key))
	b := make([]byte, 16)

	rh := *(*reflect.SliceHeader)(unsafe.Pointer(&b))
	rh.Len = 2

	hashBytes := *(*[]uint64)(unsafe.Pointer(&rh))
	high, low := hash.Values()
	hashBytes[0] = high
	hashBytes[1] = low
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return uuid.FromBytes(b)
}

// func isOptoutSuccessFile(str string) bool {
// 	optoutStrings := []string{"opt_out", "singles", "_SUCCESS"}
// 	match := true
// 	for _, sub := range optoutStrings {
// 		if !strings.Contains(str, sub) {
// 			return false
// 		}
// 	}
// 	return match
// }

func isOptoutFile(str string) bool {
	optoutStrings := []string{"opt_out", "singles", ".json.gz"}
	match := true
	for _, sub := range optoutStrings {
		if !strings.Contains(str, sub) {
			return false
		}
	}
	return match
}

// fetchComplianceGzFilesFromS3 lists the tifa objects in a bucket.
func (i *Importer) fetchComplianceGzFilesFromS3(client *s3.Client) ([]types.Object, error) {
	results, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(i.config.S3.Bucket),
	})
	var contents []types.Object
	if err != nil {
		i.logger.Error("Cannot fetch the Compliance files from bucket", log.Metadata{"bucketName": i.config.S3.Bucket, "error": err})
	} else {
		i.logger.Info("List objects in bucket", log.Metadata{"bucketName": i.config.S3.Bucket, "Contents": results.Contents})
		contents = results.Contents
	}
	return contents, err
}

func extractGzipFile(file string) (io.Reader, error) {
	// open archived file
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	return gzip.NewReader(f)
}

// concurrent processes a file by splitting the file
// processing the files concurrently and returning the result.
func (i *Importer) ProcessFile(file string, numWorkers int, batchSize int) bool {

	// open file
	f, err := extractGzipFile(file)
	if err != nil {
		i.logger.Error("Unable to do extract gzip file ", log.Metadata{"error": err, "file": file})
	}

	reader := reader(f, batchSize)
	worker := worker()
	combiner := combiner()

	// create a main context, and call cancel at the end, to ensure all our
	// goroutines exit without leaving leaks.
	// Particularly, if this function becomes part of a program with
	// a longer lifetime than this function.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// STAGE 1: start reader
	rowsBatch := [][]byte{}
	rowsCh := reader(ctx, &rowsBatch)

	// STAGE 2: create a slice of processed output channels with size of numWorkers
	// and assign each slot with the out channel from each worker.
	workersCh := make([]<-chan processed, numWorkers)
	for i := 0; i < numWorkers; i++ {
		workersCh[i] = worker(ctx, rowsCh)
	}

	// STAGE 3: read from the combined channel and calculate the final result.
	// this will end once all channels from workers are closed!
	for processed := range combiner(ctx, workersCh...) {
		i.logger.Debug("Processed ", log.Metadata{"processed": processed})
	}

	return true
}

func reader(f io.Reader, batchSize int) func(context.Context, *[][]byte) <-chan [][]byte {
	// reader creates and returns a channel that recieves
	// batches of rows (of length batchSize) from the file
	return func(ctx context.Context, rowsBatch *[][]byte) <-chan [][]byte {
		ch := make(chan [][]byte)
		contents := bufio.NewScanner(f)

		scanner := bufio.NewScanner(f)

		for scanner.Scan() {
			go processRow(scanner.Bytes())
		}

		go func(ch chan [][]byte, contents *bufio.Scanner) {
			defer close(ch) // close channel when we are done sending all rows

			for {
				scanned := contents.Scan()

				select {
				case <-ctx.Done():
					return
				default:
					row := contents.Bytes()
					// if batch size is complete or end of file, send batch out
					if len(*rowsBatch) == batchSize || !scanned {
						ch <- *rowsBatch
						*rowsBatch = [][]byte{} // clear batch
					}
					*rowsBatch = append(*rowsBatch, row) // add row to current batch
				}

				// if nothing else to scan return
				if !scanned {
					return
				}
			}
		}(ch, contents)

		return ch
	}
}

func worker() func(context.Context, <-chan [][]byte) <-chan processed {
	// worker takes in a read-only channel to recieve batches of rows.
	// After it processes each row-batch it sends out the processed output
	// on its channel.
	return func(ctx context.Context, rowBatch <-chan [][]byte) <-chan processed {
		out := make(chan processed)

		go func() {
			defer close(out)

			p := processed{}
			for rowBatch := range rowBatch {
				for _, row := range rowBatch {
					consentType, id, idType, regulationType := processRow(row)
					p.consentTypes = append(p.consentTypes, consentType)
					p.ids = append(p.ids, id)
					p.idTypes = append(p.idTypes, idType)
					p.regulationTypes = append(p.regulationTypes, regulationType)
					p.numRows++
				}
			}
			out <- p
		}()

		return out
	}
}

func combiner() func(context.Context, ...<-chan processed) <-chan processed {
	// combiner takes in multiple read-only channels that receive processed output
	// (from workers) and sends it out on it's own channel via a multiplexer.
	return func(ctx context.Context, inputs ...<-chan processed) <-chan processed {
		out := make(chan processed)

		var wg sync.WaitGroup
		multiplexer := func(p <-chan processed) {
			defer wg.Done()

			for in := range p {
				select {
				case <-ctx.Done():
				case out <- in:
				}
			}
		}

		// add length of input channels to be consumed by mutiplexer
		wg.Add(len(inputs))
		for _, in := range inputs {
			go multiplexer(in)
		}

		// close channel after all inputs channels are closed
		go func() {
			wg.Wait()
			close(out)
		}()

		return out
	}
}

func processRow(data []byte) (consentType, id, idType, regulationType string) {
	c := gjson.GetBytes(data, "consent_type").String()
	i := gjson.GetBytes(data, "id").String()
	t := gjson.GetBytes(data, "id_type").String()
	r := gjson.GetBytes(data, "regulation_type").String()
	return c, i, t, r
}
