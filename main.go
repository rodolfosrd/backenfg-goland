package main

import (
	"github.com/202lp2/go2/routers"
	"os"
)

func main() {

	r := routers.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
	  port = "8080"
	}
	r.Run(":"+port)//"localhost:8081"
}
