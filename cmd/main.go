package main

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/sirupsen/logrus"
)

var rootCmd = &cobra.Command{
	Use:   "secroute",
	Short: "SecRoute Webhook - enforce creation of secure routes",
}

var webhookServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Start HTTP server for processing OCP Routes object mutation",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("Starting up webhook server...")

	},
}

func init() {
	rootCmd.AddCommand(webhookServerCmd)
	// Init log
	logrus.SetOutput(os.Stdout)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

}

func main() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
