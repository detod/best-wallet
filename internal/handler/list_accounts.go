package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewListAccounts(
	db *pgxpool.Pool,
) *ListAccounts {
	return &ListAccounts{
		db: db,
	}
}

type ListAccounts struct {
	db *pgxpool.Pool
}

type ListAccountsResponse struct {
	Accounts []dbGetAccountForCustomer `json:"accounts"`
}

func (h *ListAccounts) Handle(c *gin.Context) {
	// Parse customer id from request.
	customerIDRaw := c.GetHeader("BestWallet-Customer-ID")
	if customerIDRaw == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, "missing customer id")
		return
	}
	customerID, err := uuid.Parse(customerIDRaw)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "malformed customer id")
		return
	}

	// Check if customer exists in db.
	exists, err := h.dbCustomerExists(c, customerID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !exists {
		c.AbortWithStatusJSON(http.StatusBadRequest, "customer not found")
		return
	}

	// Read customer accounts from db.
	accounts, err := h.dbGetAccountsForCustomer(c, customerID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, ListAccountsResponse{Accounts: accounts})
}

func (h *ListAccounts) dbCustomerExists(ctx context.Context, customerID uuid.UUID) (ok bool, err error) {
	sql := `SELECT true FROM customers WHERE id = $1`

	err = h.db.QueryRow(ctx, sql, customerID).Scan(&ok)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return false, nil
	case err != nil:
		return false, err
	default:
		return ok, nil
	}
}

type dbGetAccountForCustomer struct {
	Number  string `json:"number"`
	Balance int    `json:"balance"`
}

func (h *ListAccounts) dbGetAccountsForCustomer(ctx context.Context, customerID uuid.UUID) ([]dbGetAccountForCustomer, error) {
	sql := `SELECT number, balance FROM accounts WHERE customer_id = $1 ORDER BY created_at desc`
	// TODO index on customer_id.

	rows, _ := h.db.Query(ctx, sql, customerID)
	res, err := pgx.CollectRows[dbGetAccountForCustomer](rows, pgx.RowToStructByName[dbGetAccountForCustomer])
	switch {
	case err != nil:
		return nil, err
	default:
		return res, nil
	}
}
