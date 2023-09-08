helm_bin = $(shell command -v helm3 2> /dev/null || command -v helm)
srv=sps-header-bidder
namespace=sps-header-bidder
THIS_FILE := $(lastword $(MAKEFILE_LIST))

dry-run-dev: deps dry-run-chart-dev
dry-run-canary: deps dry-run-chart-canary
dry-run-prod: deps dry-run-chart-prod
install-dev: deps install-chart-dev
install-canary: deps install-chart-canary
install-prod: deps install-chart-prod

deps:
	$(helm_bin) dep update

dry-run-chart-dev:
	$(helm_bin) deploy --dry-run --name $(srv) -n $(namespace) --kube-context use1-rdev -f values.yaml -f dev.yaml .

dry-run-chart-canary:
	$(helm_bin) deploy --dry-run --name $(srv) -n $(namespace)-canary --kube-context use1-rprod -f values.yaml -f canary.yaml .

dry-run-chart-prod:
	$(helm_bin) deploy --dry-run --name $(srv) -n $(namespace) --kube-context use1-rprod -f values.yaml -f prod.yaml .

install-chart-dev:
	$(helm_bin) deploy --name $(srv)  -n $(namespace) --kube-context use1-rdev -f values.yaml -f dev.yaml . --wet

install-chart-canary:
	$(helm_bin) deploy --name $(srv) -n $(namespace)-canary --kube-context use1-rprod -f values.yaml -f canary.yaml . --wet

install-chart-prod:
	$(helm_bin) deploy --name $(srv) -n $(namespace) --kube-context use1-rprod -f values.yaml -f prod.yaml . --wet

install-chart-dev:
	$(helm_bin) deploy --name $(srv) -n $(namespace) --kube-context use1-rdev -f values.yaml -f dev.yaml . --wet

install-chart-canary:
	$(helm_bin) deploy --name $(srv) -n $(namespace)-canary --kube-context use1-rprod -f values.yaml -f canary.yaml . --wet

install-chart-prod:
	$(helm_bin) deploy --name $(srv) -n $(namespace) --kube-context use1-rprod -f values.yaml -f prod.yaml . --wet

delete-dev:
	$(helm_bin) delete $(srv) -n $(namespace) --kube-context use1-rdev

delete-canary:
	$(helm_bin) delete $(srv) -n $(namespace)-canary --kube-context use1-rprod

delete-prod:
	$(helm_bin) delete $(srv) -n $(namespace) --kube-context use1-rprod