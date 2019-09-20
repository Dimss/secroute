package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var runWebhookServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Start HTTP server for processing OCP Routes object mutation",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("Starting up webhook server...")

	},
}


func StartHttpRouter() {
	//cert := viper.GetString("http.crt")
	//key := viper.GetString("http.key")
	//pair, err := tls.LoadX509KeyPair(cert, key)
	//if err != nil {
	//	logrus.Error("Failed to load key pair: %v", err)
	//}
	//// Buffered channel for AD users
	//adUsersChan := make(chan string, 100)
	//// Watch and process ad users when they pushed to adUsersChan channel
	//
	//// Handel admission webhook request
	//http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	oauthtokenwebhook.WebHookHandler(w, r, adUsersChan)
	//})
	//// Handle health check request
	//http.HandleFunc("/healthz", oauthtokenwebhook.LivenessHandler)
	//// Create HTTPS server configuration
	//s := &http.Server{
	//	Addr:      ":8080",
	//	TLSConfig: &tls.Config{Certificates: []tls.Certificate{pair}},
	//}
	//// Start HTTPS server
	//log.Fatal(s.ListenAndServeTLS("", ""))
}
