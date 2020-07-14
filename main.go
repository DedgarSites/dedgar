package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/dedgarsites/dedgar/downloader"
	"github.com/dedgarsites/dedgar/routers"
)

var (
	certFile    = os.Getenv("CERT_FILE")
	keyFile     = os.Getenv("KEY_FILE")
	downloadURL = os.Getenv("DOWNLOAD_URL")
	filePath    = os.Getenv("TLS_FILE_PATH")
)

func main() {
	e := routers.Routers

	err := downloader.FileFromURL(downloadURL, filePath, certFile, keyFile)
	if err != nil {
		fmt.Println(err)
	}

	if localPort := os.Getenv("LOCAL_TESTING"); localPort != "" {
		e.Logger.Info(e.Start(":" + localPort))
	} else {
		go func() {
			time.Sleep(24 * 60 * time.Hour)
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			if err := e.Shutdown(ctx); err != nil {
				e.Logger.Info(err)
			}
		}()

		if _, err := os.Stat(filePath + certFile); os.IsNotExist(err) {
			fmt.Println("Cert file does not exist:", err)
		}
		e.Logger.Info(e.StartTLS(":8443", filePath+certFile, filePath+keyFile))
	}
}
