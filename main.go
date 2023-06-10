package main

import (
	"os"

	"github.com/Nekodigi/gpt-playground-backend/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	//handler.Firestore()
	if len(os.Args) == 2 && os.Args[1] == "setup" {
		//config.Setup()
	} else {
		engine := gin.Default()
		handler.Router(engine)
		engine.Run(":8080")
	}
}
