package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// ctx := context.Background()

	// pgUrl := os.Getenv("SHORTY_POSTGRES_URL")
	// if pgUrl == "" {
	// 	panic("empty pg url")
	// }

	// baseUrl := os.Getenv("SHORTY_BASE_URL")
	// if baseUrl == "" {
	// 	panic("empty base url")
	// }

	// db, err := pgx.Connect(ctx, pgUrl)
	// if err != nil {
	// 	panic(err)
	// }

	// storage := NewStorage(db)

	server := gin.New()
	server.Use(gin.Recovery())
	server.Use(gin.Logger())

	server.GET("/", func(ctx *gin.Context) {
		file, err := os.Open("views/index.html")
		if err != nil {
			ctx.Status(404)
		}
		defer file.Close()

		ctx.DataFromReader(200, -1, "text/html", file, nil)
	})

	server.GET("/create", func(ctx *gin.Context) {
		//create link and return page with it
	})

	server.GET("/:id", func(ctx *gin.Context) {
		//resolve link
	})

	server.Run(":8081")
}
