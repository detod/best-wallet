package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/detod/best-wallet/internal/handler"
	"github.com/detod/best-wallet/internal/middleware"
)

var ctx = context.Background()

func main() {
	// Postgres.
	_, err := pgxpool.New(ctx, os.Getenv("POSTGRES_CONN_STRING"))
	if err != nil {
		log.Fatal("Can't connect to postgres: ", err)
	}

	// Redis.
	redis := redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS_ADDR")})
	if _, err := redis.Ping(ctx).Result(); err != nil {
		log.Fatal("Can't talk to redis: ", err)
	}

	// TODO logging.
	// TODO tracing.
	// TODO catch signals and shutdown gracefully.

	// Handlers.
	createAcc := handler.NewCreateAccount()
	readAcc := handler.NewReadAccount()
	deposit := handler.NewDeposit()
	withdraw := handler.NewWithdraw()
	transfer := handler.NewTransfer()

	// Middleware.
	hmacVerifier := middleware.HMACVerifier(nil)

	// Routing.
	r := gin.Default()
	v1 := r.Group("/api/v1")
	{

		v1.POST("/accounts", hmacVerifier, createAcc.Handle)      // Create account.
		v1.GET("/accounts/:number", hmacVerifier, readAcc.Handle) // Read account data.

		v1.POST("/accounts/:number/deposit", hmacVerifier, deposit.Handle)   // Money coming into the wallet.
		v1.POST("/accounts/:number/withdraw", hmacVerifier, withdraw.Handle) // Money leaving the wallet.
		v1.POST("/accounts/transfer", hmacVerifier, transfer.Handle)         // Money moving within the wallet.
	}

	// TODO deep healthchecks.
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Serve HTTP.
	r.Run(os.Getenv("LISTEN_ADDR"))
}
