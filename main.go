package main

import (
	"os"

	"github.com/dedgarsites/dedgar/routers"
)

func main() {
	e := routers.Routers
	if localPort := os.Getenv("LOCAL_TESTING"); localPort != "" {
		e.Logger.Info(e.Start(":" + localPort))
	} else {
		e.Logger.Info(e.StartTLS(":8443", "/cert/lego/certificates/dashboard.crt", "/cert/lego/certificates/dashboard.key"))
	}
}
