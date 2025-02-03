GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)

.PHONY: mock
mock:
	mockgen -source v1/config/config.go -destination v1/config/mock/config.go -package config_mock
	mockgen -source v1/nats/queues.go -destination v1/nats/mock/queues.go -package nats_mock
	mockgen -source v2/nats/queues.go -destination v2/nats/mock/queues.go -package nats_mock
	mockgen -source v3/dialer/dialer.go -destination v3/dialer/mock/dialer.go -package dialer_mock
	mockgen -source v3/dialer/dialer_manager.go -destination v3/dialer/mock/dialer_manager.go -package dialer_mock
	mockgen -source v3/jwt/secrets.go -destination v3/jwt/mock/secrets.go -package jwt_mock
