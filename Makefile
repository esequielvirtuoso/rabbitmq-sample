version ?= latest
PROJECT  = rabbitmq-sample
HUB_USER = esequielvirtuoso
BUILD    = $(shell git rev-parse --short HEAD)
IMAGE    = $(HUB_USER)/$(PROJECT):$(BUILD)
