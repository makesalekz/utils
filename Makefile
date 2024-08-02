GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)

.PHONY: mock
mock:
	mockgen -source v1/nats/queues.go -destination v1/nats/mock/queues.go -package nats_mock
	mockgen -source v2/dialer/dialer.go -destination v2/dialer/mock/dialer.go -package dialer_mock
	mockgen -source v2/dialer/dialer_manager.go -destination v2/dialer/mock/dialer_manager.go -package dialer_mock
	mockgen -source v2/jwt/claims.go -destination v2/jwt/mock/claims.go -package jwt_mock
	mockgen -source v2/jwt/processor.go -destination v2/jwt/mock/processor.go -package jwt_mock
