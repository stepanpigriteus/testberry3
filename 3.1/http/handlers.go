package http

import (
	"treeOne/domain"

	"github.com/rs/zerolog"
)

type HandleNotify struct {
	logger  zerolog.Logger
	service domain.Service
}

func NewHandleNotify(logger zerolog.Logger, service domain.Service) *HandleNotify {
	return &HandleNotify{
		logger:  logger,
		service: service,
	}
}
