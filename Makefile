IMG ?= ghcr.io/unict-cclab/sophos-telemetry:latest
PROMETHEUS_ADDRESS ?= http://localhost:9090

run:
	PROMETHEUS_ADDRESS=${PROMETHEUS_ADDRESS} go run main.go

build:
	go build

build-image:
	docker build -t ${IMG} .

push-image:
	docker push ${IMG}