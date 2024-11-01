package server

import (
	"github.com/asdine/storm/v3"
	"github.com/fanjindong/go-cache"
	"github.com/hantbk/vbackup/internal/config"
	"github.com/kataras/iris/v12"
	"github.com/sirupsen/logrus"
)

type BackupServer struct {
	app          *iris.Application
	logger       *logrus.Logger
	rootRoute    iris.Party
	config       *config.Config
	db           *storm.DB
	cache        cache.ICache
	systemStatus string
	isDocker     bool
}

var EmbedWebDashboard embed.FS

var bs *BackupServer

