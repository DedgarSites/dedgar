package main

import (
	"context"
	"os"
	"time"

	"github.com/dedgarsites/dedgar/routers"
)

func main() {
	e := routers.Routers
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

		e.Logger.Info(e.StartTLS(":8443", "/cert/lego/certificates/dedgar.crt", "/cert/lego/certificates/dedgar.key"))
	}
}
