package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	pgxgoogleuuid "github.com/vgarvardt/pgx-google-uuid/v5"

	"github.com/detod/best-wallet/internal/handler"
	"github.com/detod/best-wallet/internal/middleware"
)

func main() {
	var ctx = context.Background()
	// TODO structured logging.

	// Postgres.
	pgxConf, err := pgxpool.ParseConfig(os.Getenv("POSTGRES_CONN_STRING"))
	if err != nil {
		log.Fatal("Can't parse pgx config: ", err)
	}
	pgxConf.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxgoogleuuid.Register(conn.TypeMap()) // So we can use google/uuid type with pgx.
		return nil
	}
	db, err := pgxpool.NewWithConfig(ctx, pgxConf) // TODO configure postgres.
	if err != nil {
		log.Fatal("Can't create pgx pool: ", err)
	}

	// Redis.
	redis := redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS_ADDR")})
	if _, err := redis.Ping(ctx).Result(); err != nil { // TODO configure redis.
		log.Fatal("Can't talk to redis: ", err)
	}

	// TODO tracing.
	// TODO catch signals and shutdown gracefully.

	// Handlers.
	createCustomer := handler.NewCreateCustomer(db)
	createAccount := handler.NewCreateAccount(db)
	listAccounts := handler.NewListAccounts(db)
	deposit := handler.NewDeposit()
	withdraw := handler.NewWithdraw()
	transfer := handler.NewTransfer()

	// Middleware.
	hmacVerifier := middleware.HMACVerifier(nil)

	// Routing.
	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		v1.POST("/customers", hmacVerifier, createCustomer.Handle) // Create customer.

		v1.POST("/accounts", hmacVerifier, createAccount.Handle)             // Open a new personal account for a customer.
		v1.GET("/accounts", hmacVerifier, listAccounts.Handle)               // List all accounts for a customer.
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
	r.Run(os.Getenv("LISTEN_ADDR")) // TODO configure server.
}
