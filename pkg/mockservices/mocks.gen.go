// Code generated by MockGen. DO NOT EDIT.
// Source: dependencies.go

// Package mockservices is a generated GoMock package.
package mockservices

import (
	io "io"
	reflect "reflect"
	time "time"

	httpclient "github.com/adgear/go-commons/pkg/httpclient"
	log "github.com/adgear/go-commons/pkg/log"
	demand "github.com/adgear/sps-header-bidder/pkg/demand"
	types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	gomock "github.com/golang/mock/gomock"
	kgo "github.com/twmb/franz-go/pkg/kgo"
)

// Mocklogger is a mock of logger interface.
type Mocklogger struct {
	ctrl     *gomock.Controller
	recorder *MockloggerMockRecorder
}

// MockloggerMockRecorder is the mock recorder for Mocklogger.
type MockloggerMockRecorder struct {
	mock *Mocklogger
}

// NewMocklogger creates a new mock instance.
func NewMocklogger(ctrl *gomock.Controller) *Mocklogger {
	mock := &Mocklogger{ctrl: ctrl}
	mock.recorder = &MockloggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mocklogger) EXPECT() *MockloggerMockRecorder {
	return m.recorder
}

// Debug mocks base method.
func (m *Mocklogger) Debug(message string, metadata ...log.Metadata) {
	m.ctrl.T.Helper()
	varargs := []interface{}{message}
	for _, a := range metadata {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Debug", varargs...)
}

// Debug indicates an expected call of Debug.
func (mr *MockloggerMockRecorder) Debug(message interface{}, metadata ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{message}, metadata...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Debug", reflect.TypeOf((*Mocklogger)(nil).Debug), varargs...)
}

// Error mocks base method.
func (m *Mocklogger) Error(message string, metadata ...log.Metadata) {
	m.ctrl.T.Helper()
	varargs := []interface{}{message}
	for _, a := range metadata {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Error", varargs...)
}

// Error indicates an expected call of Error.
func (mr *MockloggerMockRecorder) Error(message interface{}, metadata ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{message}, metadata...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*Mocklogger)(nil).Error), varargs...)
}

// Fatal mocks base method.
func (m *Mocklogger) Fatal(message string, metadata ...log.Metadata) {
	m.ctrl.T.Helper()
	varargs := []interface{}{message}
	for _, a := range metadata {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Fatal", varargs...)
}

// Fatal indicates an expected call of Fatal.
func (mr *MockloggerMockRecorder) Fatal(message interface{}, metadata ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{message}, metadata...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fatal", reflect.TypeOf((*Mocklogger)(nil).Fatal), varargs...)
}

// GetDefaultLogFormat mocks base method.
func (m *Mocklogger) GetDefaultLogFormat() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDefaultLogFormat")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetDefaultLogFormat indicates an expected call of GetDefaultLogFormat.
func (mr *MockloggerMockRecorder) GetDefaultLogFormat() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDefaultLogFormat", reflect.TypeOf((*Mocklogger)(nil).GetDefaultLogFormat))
}

// GetLogFormat mocks base method.
func (m *Mocklogger) GetLogFormat() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLogFormat")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetLogFormat indicates an expected call of GetLogFormat.
func (mr *MockloggerMockRecorder) GetLogFormat() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLogFormat", reflect.TypeOf((*Mocklogger)(nil).GetLogFormat))
}

// GetLogLevel mocks base method.
func (m *Mocklogger) GetLogLevel() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLogLevel")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetLogLevel indicates an expected call of GetLogLevel.
func (mr *MockloggerMockRecorder) GetLogLevel() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLogLevel", reflect.TypeOf((*Mocklogger)(nil).GetLogLevel))
}

// Info mocks base method.
func (m *Mocklogger) Info(message string, metadata ...log.Metadata) {
	m.ctrl.T.Helper()
	varargs := []interface{}{message}
	for _, a := range metadata {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Info", varargs...)
}

// Info indicates an expected call of Info.
func (mr *MockloggerMockRecorder) Info(message interface{}, metadata ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{message}, metadata...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*Mocklogger)(nil).Info), varargs...)
}

// Panic mocks base method.
func (m *Mocklogger) Panic(msg string, metadata ...log.Metadata) {
	m.ctrl.T.Helper()
	varargs := []interface{}{msg}
	for _, a := range metadata {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Panic", varargs...)
}

// Panic indicates an expected call of Panic.
func (mr *MockloggerMockRecorder) Panic(msg interface{}, metadata ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{msg}, metadata...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Panic", reflect.TypeOf((*Mocklogger)(nil).Panic), varargs...)
}

// SetEnvKey mocks base method.
func (m *Mocklogger) SetEnvKey(environment string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetEnvKey", environment)
}

// SetEnvKey indicates an expected call of SetEnvKey.
func (mr *MockloggerMockRecorder) SetEnvKey(environment interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetEnvKey", reflect.TypeOf((*Mocklogger)(nil).SetEnvKey), environment)
}

// SetLevel mocks base method.
func (m *Mocklogger) SetLevel(level log.Level) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetLevel", level)
}

// SetLevel indicates an expected call of SetLevel.
func (mr *MockloggerMockRecorder) SetLevel(level interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLevel", reflect.TypeOf((*Mocklogger)(nil).SetLevel), level)
}

// SetLogFormat mocks base method.
func (m *Mocklogger) SetLogFormat(logFormat string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetLogFormat", logFormat)
}

// SetLogFormat indicates an expected call of SetLogFormat.
func (mr *MockloggerMockRecorder) SetLogFormat(logFormat interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLogFormat", reflect.TypeOf((*Mocklogger)(nil).SetLogFormat), logFormat)
}

// SetLogLevel mocks base method.
func (m *Mocklogger) SetLogLevel(logLevel string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetLogLevel", logLevel)
}

// SetLogLevel indicates an expected call of SetLogLevel.
func (mr *MockloggerMockRecorder) SetLogLevel(logLevel interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLogLevel", reflect.TypeOf((*Mocklogger)(nil).SetLogLevel), logLevel)
}

// SetLogParams mocks base method.
func (m *Mocklogger) SetLogParams(logFormat, environment, logLevel, appName, fileName string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetLogParams", logFormat, environment, logLevel, appName, fileName)
}

// SetLogParams indicates an expected call of SetLogParams.
func (mr *MockloggerMockRecorder) SetLogParams(logFormat, environment, logLevel, appName, fileName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLogParams", reflect.TypeOf((*Mocklogger)(nil).SetLogParams), logFormat, environment, logLevel, appName, fileName)
}

// SetWriteSyncer mocks base method.
func (m *Mocklogger) SetWriteSyncer(out io.Writer) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetWriteSyncer", out)
}

// SetWriteSyncer indicates an expected call of SetWriteSyncer.
func (mr *MockloggerMockRecorder) SetWriteSyncer(out interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetWriteSyncer", reflect.TypeOf((*Mocklogger)(nil).SetWriteSyncer), out)
}

// Sync mocks base method.
func (m *Mocklogger) Sync() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sync")
	ret0, _ := ret[0].(error)
	return ret0
}

// Sync indicates an expected call of Sync.
func (mr *MockloggerMockRecorder) Sync() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sync", reflect.TypeOf((*Mocklogger)(nil).Sync))
}

// Warn mocks base method.
func (m *Mocklogger) Warn(message string, metadata ...log.Metadata) {
	m.ctrl.T.Helper()
	varargs := []interface{}{message}
	for _, a := range metadata {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Warn", varargs...)
}

// Warn indicates an expected call of Warn.
func (mr *MockloggerMockRecorder) Warn(message interface{}, metadata ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{message}, metadata...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Warn", reflect.TypeOf((*Mocklogger)(nil).Warn), varargs...)
}

// MockdemandClient is a mock of demandClient interface.
type MockdemandClient struct {
	ctrl     *gomock.Controller
	recorder *MockdemandClientMockRecorder
}

// MockdemandClientMockRecorder is the mock recorder for MockdemandClient.
type MockdemandClientMockRecorder struct {
	mock *MockdemandClient
}

// NewMockdemandClient creates a new mock instance.
func NewMockdemandClient(ctrl *gomock.Controller) *MockdemandClient {
	mock := &MockdemandClient{ctrl: ctrl}
	mock.recorder = &MockdemandClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockdemandClient) EXPECT() *MockdemandClientMockRecorder {
	return m.recorder
}

// BidOrtbReq mocks base method.
func (m *MockdemandClient) BidOrtbReq(demandParams demand.DemandExtParams) (int, []byte) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BidOrtbReq", demandParams)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].([]byte)
	return ret0, ret1
}

// BidOrtbReq indicates an expected call of BidOrtbReq.
func (mr *MockdemandClientMockRecorder) BidOrtbReq(demandParams interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BidOrtbReq", reflect.TypeOf((*MockdemandClient)(nil).BidOrtbReq), demandParams)
}

// MockhttpClient is a mock of httpClient interface.
type MockhttpClient struct {
	ctrl     *gomock.Controller
	recorder *MockhttpClientMockRecorder
}

// MockhttpClientMockRecorder is the mock recorder for MockhttpClient.
type MockhttpClientMockRecorder struct {
	mock *MockhttpClient
}

// NewMockhttpClient creates a new mock instance.
func NewMockhttpClient(ctrl *gomock.Controller) *MockhttpClient {
	mock := &MockhttpClient{ctrl: ctrl}
	mock.recorder = &MockhttpClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockhttpClient) EXPECT() *MockhttpClientMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockhttpClient) Delete(request *httpclient.Request) (*httpclient.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", request)
	ret0, _ := ret[0].(*httpclient.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockhttpClientMockRecorder) Delete(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockhttpClient)(nil).Delete), request)
}

// Get mocks base method.
func (m *MockhttpClient) Get(request *httpclient.Request) (*httpclient.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", request)
	ret0, _ := ret[0].(*httpclient.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockhttpClientMockRecorder) Get(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockhttpClient)(nil).Get), request)
}

// Post mocks base method.
func (m *MockhttpClient) Post(request *httpclient.Request) (*httpclient.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Post", request)
	ret0, _ := ret[0].(*httpclient.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Post indicates an expected call of Post.
func (mr *MockhttpClientMockRecorder) Post(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Post", reflect.TypeOf((*MockhttpClient)(nil).Post), request)
}

// Put mocks base method.
func (m *MockhttpClient) Put(request *httpclient.Request) (*httpclient.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Put", request)
	ret0, _ := ret[0].(*httpclient.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Put indicates an expected call of Put.
func (mr *MockhttpClientMockRecorder) Put(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockhttpClient)(nil).Put), request)
}

// MockkafkaProducer is a mock of kafkaProducer interface.
type MockkafkaProducer struct {
	ctrl     *gomock.Controller
	recorder *MockkafkaProducerMockRecorder
}

// MockkafkaProducerMockRecorder is the mock recorder for MockkafkaProducer.
type MockkafkaProducerMockRecorder struct {
	mock *MockkafkaProducer
}

// NewMockkafkaProducer creates a new mock instance.
func NewMockkafkaProducer(ctrl *gomock.Controller) *MockkafkaProducer {
	mock := &MockkafkaProducer{ctrl: ctrl}
	mock.recorder = &MockkafkaProducerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockkafkaProducer) EXPECT() *MockkafkaProducerMockRecorder {
	return m.recorder
}

// Produce mocks base method.
func (m *MockkafkaProducer) Produce(record *kgo.Record, promise func(*kgo.Record, error)) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Produce", record, promise)
}

// Produce indicates an expected call of Produce.
func (mr *MockkafkaProducerMockRecorder) Produce(record, promise interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Produce", reflect.TypeOf((*MockkafkaProducer)(nil).Produce), record, promise)
}

// ProduceSync mocks base method.
func (m *MockkafkaProducer) ProduceSync(record *kgo.Record) (*kgo.Record, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProduceSync", record)
	ret0, _ := ret[0].(*kgo.Record)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProduceSync indicates an expected call of ProduceSync.
func (mr *MockkafkaProducerMockRecorder) ProduceSync(record interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProduceSync", reflect.TypeOf((*MockkafkaProducer)(nil).ProduceSync), record)
}

// PublishKafkaEvent mocks base method.
func (m *MockkafkaProducer) PublishKafkaEvent(topic string, message []byte, schemaKey, schemaFile string, promise func(*kgo.Record, error)) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "PublishKafkaEvent", topic, message, schemaKey, schemaFile, promise)
}

// PublishKafkaEvent indicates an expected call of PublishKafkaEvent.
func (mr *MockkafkaProducerMockRecorder) PublishKafkaEvent(topic, message, schemaKey, schemaFile, promise interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishKafkaEvent", reflect.TypeOf((*MockkafkaProducer)(nil).PublishKafkaEvent), topic, message, schemaKey, schemaFile, promise)
}

// PublishKafkaEventSync mocks base method.
func (m *MockkafkaProducer) PublishKafkaEventSync(topic string, message []byte, schemaKey, schemaFile string) (*kgo.Record, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublishKafkaEventSync", topic, message, schemaKey, schemaFile)
	ret0, _ := ret[0].(*kgo.Record)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PublishKafkaEventSync indicates an expected call of PublishKafkaEventSync.
func (mr *MockkafkaProducerMockRecorder) PublishKafkaEventSync(topic, message, schemaKey, schemaFile interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishKafkaEventSync", reflect.TypeOf((*MockkafkaProducer)(nil).PublishKafkaEventSync), topic, message, schemaKey, schemaFile)
}

// MockcacheClient is a mock of cacheClient interface.
type MockcacheClient struct {
	ctrl     *gomock.Controller
	recorder *MockcacheClientMockRecorder
}

// MockcacheClientMockRecorder is the mock recorder for MockcacheClient.
type MockcacheClientMockRecorder struct {
	mock *MockcacheClient
}

// NewMockcacheClient creates a new mock instance.
func NewMockcacheClient(ctrl *gomock.Controller) *MockcacheClient {
	mock := &MockcacheClient{ctrl: ctrl}
	mock.recorder = &MockcacheClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockcacheClient) EXPECT() *MockcacheClientMockRecorder {
	return m.recorder
}

// GetTifa mocks base method.
func (m *MockcacheClient) GetTifa(arg0 string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTifa", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// GetTifa indicates an expected call of GetTifa.
func (mr *MockcacheClientMockRecorder) GetTifa(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTifa", reflect.TypeOf((*MockcacheClient)(nil).GetTifa), arg0)
}

// IsLastLoadTsExpired mocks base method.
func (m *MockcacheClient) IsLastLoadTsExpired(arg0 string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsLastLoadTsExpired", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsLastLoadTsExpired indicates an expected call of IsLastLoadTsExpired.
func (mr *MockcacheClientMockRecorder) IsLastLoadTsExpired(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsLastLoadTsExpired", reflect.TypeOf((*MockcacheClient)(nil).IsLastLoadTsExpired), arg0)
}

// SetTifa mocks base method.
func (m *MockcacheClient) SetTifa(arg0, arg1 string, arg2 time.Duration) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTifa", arg0, arg1, arg2)
}

// SetTifa indicates an expected call of SetTifa.
func (mr *MockcacheClientMockRecorder) SetTifa(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTifa", reflect.TypeOf((*MockcacheClient)(nil).SetTifa), arg0, arg1, arg2)
}

// Mockudws3Client is a mock of udws3Client interface.
type Mockudws3Client struct {
	ctrl     *gomock.Controller
	recorder *Mockudws3ClientMockRecorder
}

// Mockudws3ClientMockRecorder is the mock recorder for Mockudws3Client.
type Mockudws3ClientMockRecorder struct {
	mock *Mockudws3Client
}

// NewMockudws3Client creates a new mock instance.
func NewMockudws3Client(ctrl *gomock.Controller) *Mockudws3Client {
	mock := &Mockudws3Client{ctrl: ctrl}
	mock.recorder = &Mockudws3ClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockudws3Client) EXPECT() *Mockudws3ClientMockRecorder {
	return m.recorder
}

// DownloadGzFile mocks base method.
func (m *Mockudws3Client) DownloadGzFile(arg0, arg1 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownloadGzFile", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DownloadGzFile indicates an expected call of DownloadGzFile.
func (mr *Mockudws3ClientMockRecorder) DownloadGzFile(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadGzFile", reflect.TypeOf((*Mockudws3Client)(nil).DownloadGzFile), arg0, arg1)
}

// FetchGzFiles mocks base method.
func (m *Mockudws3Client) FetchGzFiles() ([]types.Object, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchGzFiles")
	ret0, _ := ret[0].([]types.Object)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchGzFiles indicates an expected call of FetchGzFiles.
func (mr *Mockudws3ClientMockRecorder) FetchGzFiles() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchGzFiles", reflect.TypeOf((*Mockudws3Client)(nil).FetchGzFiles))
}
