package common

import (
	"github.com/asdine/storm/v3"
	"github.com/hantbk/vbackup/internal/server"
)

type DBService interface {
	GetDB(options DBOptions) storm.Node
}

type DefaultDBService struct{}

type DBOptions struct {
	DB storm.Node
}

func (d *DefaultDBService) GetDB(options DBOptions) storm.Node {
	if options.DB != nil {
		return options.DB
	}
	return server.DB()
}
