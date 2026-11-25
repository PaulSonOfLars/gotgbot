package main

import (
	"log"
	"net/http"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/gin-gonic/gin"
)

func StartServer(
	updater *ext.Updater,
	listenAddr string,
) {
	gin.SetMode(gin.ReleaseMode)
	server := gin.New()

	server.POST("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	// This will be the handler for each bot.
	// As stated in main.go, we will set webhook to
	// "http://localhost:8080/bots/<botToken>", so we will
	// use "/bots/:token" as the path.
	server.POST("/bots/:token", func(c *gin.Context) {
		handler := updater.GetHandlerFunc("/bots/")
		handler(c.Writer, c.Request)
	})

	log.Printf("Webhook server started at %s\n", listenAddr)
	go func() {
		err := http.ListenAndServe(listenAddr, server)
		if err != nil {
			panic("failed to start server: " + err.Error())
		}
	}()
}
