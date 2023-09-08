test-with-report:
	mkdir -p coverage
	go test -v -race ./... -coverprofile=./coverage/coverage.out -json > ./coverage/test-report.json

eunit:
	@echo "Running unit tests"
	@go test -v ./...

set-pipeline:
	@echo setting concourse pipeline
	@fly -t prod.tizen-ads-sdk set-pipeline -c ci/pipeline/concourse_config.yml  \
	   -p sps-header-bidder \
	   -l ci/pipeline/vars.yml
	   
validate-pipeline:
	@echo validating concourse pipeline
	@fly -t prod.tizen-ads-sdk validate-pipeline -c ci/pipeline/concourse_config.yml \
	   -l ci/pipeline/vars.yml --strict

swagger:
	swagger generate spec -o ./swagger/swagger.json --scan-models

install_swagger:
	go get -u github.com/go-swagger/go-swagger/cmd/swagger
	go install github.com/go-swagger/go-swagger/cmd/swagger

.PHONY: eunit set-pipeline validate-pipeline

