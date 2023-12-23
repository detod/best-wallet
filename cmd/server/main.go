package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	urlExample := fmt.Sprintf("postgres://postgres:asd9fwepub83lf@postgres:5432/bestwallet")
	conn, err := pgxpool.New(ctx, urlExample)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()
	log.Println("Successfuly connected to postgres")

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer func() {
		if err := rdb.Close(); err != nil {
			log.Println("failed to close redis client")
		}
		log.Println("Closed redis client")
	}()
	res, err := rdb.Ping(ctx).Result()
	log.Println("Redis ping returned: ", res, err)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
