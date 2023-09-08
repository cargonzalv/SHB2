package apitest

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"runtime"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

const SAMPLE_FREEWHEEL_PATH = "testing/freewheel/"
const SAMPLE_SPRINGSERVE_PATH = "testing/springserve/"

var sample *Sample

type SpringserveApiRequest struct {
	Valid       []byte
	InvalidJson []byte
}

type FreewheelApiRequest struct {
	ValidDealId      []byte
	DealsWithoutId   []byte
	DealsWithEmptyId []byte
	NoDeals          []byte
	Unexpected       []byte
	InvalidJson      []byte
}

type DemandApiResponse struct {
	ValidWithDealId    []byte
	ValidWithoutDealId []byte
	Unexpected         []byte
}

type Sample struct {
	Url                    string
	SpringserveApiRequests SpringserveApiRequest
	FreewheelApiRequests   FreewheelApiRequest
	DemandApiResponse      DemandApiResponse
}

func Samples() *Sample {
	if sample == nil {
		loadSample()
	}

	return sample
}

func Request(body []byte, url string) (*http.Request, error) {
	payload := bytes.NewBuffer(body)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return req, err
	}
	req.Header.Add("-Encoding", "gzip")

	return req, nil
}

// serve() serves http request using provided rest handler
func ProcessRequest(listener *fasthttputil.InmemoryListener, r *http.Request, disableCompression bool) (*http.Response, error) {
	defer listener.Close()

	tr := &http.Transport{
		DisableCompression: disableCompression,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return listener.Dial()
		},
	}

	client := http.Client{Transport: tr}
	return client.Do(r)
}

func InitServer(h fasthttp.RequestHandler) *fasthttputil.InmemoryListener {
	listener := fasthttputil.NewInmemoryListener()

	go func() {
		err := fasthttp.Serve(listener, h)
		if err != nil {
			panic(fmt.Errorf("failed to serve: %v", err))
		}
	}()

	return listener
}

func loadSample() {
	springserveApiRequest := SpringserveApiRequest{
		Valid:       loadSpringserveSampleFile("springserve_request_valid.json"),
		InvalidJson: []byte(`invalid json`),
	}
	freewheelRequestsSample := FreewheelApiRequest{
		ValidDealId:      loadFreeWheelSampleFile("freewheel_request_with_deal_id.json"),
		DealsWithoutId:   loadFreeWheelSampleFile("freewheel_request_without_deal_id.json"),
		DealsWithEmptyId: loadFreeWheelSampleFile("freewheel_request_with_empty_deal_id.json"),
		NoDeals:          loadFreeWheelSampleFile("freewheel_request_with_no_deals.json"),
		Unexpected:       loadFreeWheelSampleFile("freewheel_request_unexpected.json"),
		InvalidJson:      []byte(`invalid json`),
	}

	demandResponseSample := DemandApiResponse{
		ValidWithDealId:    loadFreeWheelSampleFile("demand_response_with_deal_id.json"),
		ValidWithoutDealId: loadFreeWheelSampleFile("demand_response_without_deal_id.json"),
		Unexpected:         loadFreeWheelSampleFile("demand_unexpected_response.json"),
	}

	sample = &Sample{
		Url:                    "http://localhost:8085/",
		SpringserveApiRequests: springserveApiRequest,
		FreewheelApiRequests:   freewheelRequestsSample,
		DemandApiResponse:      demandResponseSample,
	}
}

func loadSpringserveSampleFile(file string) []byte {
	return loadFile(SAMPLE_SPRINGSERVE_PATH + file)
}

func loadFreeWheelSampleFile(file string) []byte {
	return loadFile(SAMPLE_FREEWHEEL_PATH + file)
}

func loadFile(filepath string) []byte {
	dir := ""
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		dir = path.Join(path.Dir(filename), "..", "..", "..")
	}
	fullPath := path.Join(dir, filepath)
	body, err := os.ReadFile(fullPath)
	if err != nil {
		log.Println("Error apitest.loadFile: unable to read file", err)
		panic(err)
	}
	return body
}
