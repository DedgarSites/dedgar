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
	fmt.Println("dedgar v0.0.1")

	e := routers.Routers

	err := downloader.FileFromURL(downloadURL, filePath, certFile, keyFile)
	if err != nil {
		fmt.Println(err)
	}

	if localPort := os.Getenv("LOCAL_TESTING"); localPort != "" {
		e.Logger.Info(e.Start(":" + localPort))
	} else {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 24*60*time.Hour)
			defer cancel()
			if err := e.Shutdown(ctx); err != nil {
				e.Logger.Fatal(err)
			}
		}()

		e.Logger.Info(e.StartTLS(":8443", filePath+certFile, filePath+keyFile))
	}
}
