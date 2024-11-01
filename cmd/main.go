package main

import (
	"embed"
	"fmt"
	"os"
	"runtime"

	"github.com/hantbk/vbackup"
	"github.com/hantbk/vbackup/internal/route"
	"github.com/spf13/cobra"
)

var (
	configPath     string
	serverBindHost string
	serverBindPort int
)

//go:embed web/dashboard
var embedWebDashboard embed.FS

var rootCmd = &cobra.Command{
	Use:   "vbackup_server",
	Short: "vBackup is a file backup service",
	Long:  `vBackup is a simple, open-source, fast, and secure server file backup tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		server.EmbedWebDashboard = embedWebDashboard
		err := server.Listen(route.InitRoute, serverBindHost, serverBindPort, configPath)
		if err != nil {
			os.Exit(1)
			return
		}
	},
	Version: fmt.Sprintf("%s compiled with %v on %v/%v at %s",
		v.Version, runtime.Version(), runtime.GOOS, runtime.GOARCH, v.BuildTime),
}
var v = vbackup.GetVersion()

var resetOtpCmd = &cobra.Command{
	Use:   "resetOtp [username]",
	Short: "Disable user OTP authentication.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		username := ""
		if len(args) > 0 {
			username = args[0]
		} else {
			fmt.Println("Username cannot be empty.")
			return
		}
		fmt.Println(fmt.Sprintf("Disabling two-factor authentication for: %s.", username))
		cmdServer.Instance(configPath).ClearOtp(username, 0)
	},
}

var resetPwdCmd = &cobra.Command{
	Use:   "resetPwd [username]",
	Short: "Reset user password.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		username := ""
		if len(args) > 0 {
			username = args[0]
		} else {
			fmt.Println("Username cannot be empty.")
			return
		}
		fmt.Println(fmt.Sprintf("Resetting password for: %s.", username))
		cmdServer.Instance(configPath).ClearPwd(username, 0)
	},
}

func init() {
	rootCmd.Flags().StringVar(&serverBindHost, "server-bind-host", "", "bind address")
	rootCmd.Flags().IntVarP(&serverBindPort, "server-bind-port", "p", 0, "bind port")
	rootCmd.PersistentFlags().StringVarP(&configPath, "config-path", "c", "", "config file path")
	rootCmd.AddCommand(resetOtpCmd)
	rootCmd.AddCommand(resetPwdCmd)
}
func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
		return
	}
}
