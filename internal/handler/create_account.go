package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/detod/best-wallet/internal/domain"
	"github.com/detod/best-wallet/internal/util"
)

func NewCreateAccount(
	db *pgxpool.Pool,
) *CreateAccount {
	return &CreateAccount{
		db: db,
	}
}

type CreateAccount struct {
	db *pgxpool.Pool
}

type CreateAccountResponse struct {
	ID     uuid.UUID `json:"id"`
	Number string    `json:"number"`
}

func (h *CreateAccount) Handle(c *gin.Context) {
	// Read customer id.
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

	// Only customers with approved KYC can open accounts.
	kycStatus, exists, err := h.dbGetCustomerKYCStatus(c, customerID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !exists {
		c.AbortWithStatusJSON(http.StatusBadRequest, "customer not found")
		return
	}
	if kycStatus != domain.KYCStatusApproved {
		c.AbortWithStatusJSON(http.StatusBadRequest, "customer not verified, try again later")
		return
	}

	// Create new account.
	id := uuid.New()
	number := uuid.New().String()
	if err = h.dbInsertAccount(c, dbInsertAccountArgs{
		id:         id,
		customerID: customerID,
		number:     number,
		balance:    0,
	}); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Kickstart background notification.
	go util.Recover(func() {
		// TODO notify customer.
	})

	// Return new account identifiers.
	c.JSON(http.StatusCreated, CreateAccountResponse{
		ID:     id,
		Number: number,
	})
}

func (h *CreateAccount) dbGetCustomerKYCStatus(ctx context.Context, customerID uuid.UUID) (domain.KYCStatus, bool, error) {
	sql := `SELECT kyc_status FROM customers WHERE id = $1`
	var res domain.KYCStatus

	err := h.db.QueryRow(ctx, sql, customerID).Scan(&res)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return "", false, nil
	case err != nil:
		return "", false, err
	default:
		return res, true, nil
	}
}

type dbInsertAccountArgs struct {
	id         uuid.UUID
	customerID uuid.UUID
	number     string
	balance    int
}

func (h *CreateAccount) dbInsertAccount(ctx context.Context, args dbInsertAccountArgs) error {
	sql := `
		INSERT INTO accounts (id, customer_id, number, balance)
		VALUES ($1, $2, $3, $4)`

	_, err := h.db.Exec(ctx, sql,
		args.id,
		args.customerID,
		args.number,
		args.balance,
	)

	return err
}
