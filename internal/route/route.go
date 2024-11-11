package route

import (
	"fmt"

	"github.com/hantbk/vbackup/internal/api"
	v1 "github.com/hantbk/vbackup/internal/api/v1"
	"github.com/hantbk/vbackup/internal/cron"
	"github.com/hantbk/vbackup/internal/server"
	"github.com/hantbk/vbackup/internal/service/v1/user"
	"github.com/hantbk/vbackup/pkg/utils"
	resticProxy "github.com/hantbk/vbackup/restic_proxy"
	"github.com/kataras/iris/v12"
)

func InitRoute(party iris.Party) {
	initOthers()
	apiParty := party.Party("/api")
	api.AddPingRoute(apiParty)
	v1.AddV1Route(apiParty)
	ininPrint()
}
func initOthers() {
	go resticProxy.InitRepository()
	initAdmin()
	utils.InitJwt()
	go cron.InitCron()
}

func ininPrint() {
	if server.Config().Prometheus.Enable {
		fmt.Printf("Prometheus is deploy to: http://%s:%d/%s\n",
			server.Config().Server.Bind.Host,
			server.Config().Server.Bind.Port,
			"metrics")
	}
	fmt.Printf("Health Check is deploy to: http://%s:%d/%s\n",
		server.Config().Server.Bind.Host,
		server.Config().Server.Bind.Port,
		"api/ping")
}

// initAdmin initializes the admin account.
func initAdmin() {
	userServer := user.GetService()
	err := userServer.InitAdmin()
	if err != nil {
		fmt.Println("Failed to initialize admin account: ", err.Error())
		return
	}
}
