package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/detod/best-wallet/internal/domain"
	"github.com/detod/best-wallet/internal/util"
)

func NewCreateCustomer(
	db *pgxpool.Pool,
) *CreateCustomer {
	return &CreateCustomer{
		db: db,
	}
}

type CreateCustomer struct {
	db *pgxpool.Pool
}

type CreateCustomerRequest struct {
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	Email            string    `json:"email"`
	ResidenceAddress string    `json:"residence_address"`
	BirthDate        time.Time `json:"birth_date"`
}

type CreateCustomerResponse struct {
	ID uuid.UUID `json:"id"`
}

func (h *CreateCustomer) Handle(c *gin.Context) {
	var req CreateCustomerRequest
	if err := c.BindJSON(&req); err != nil {
		return
	}

	// TODO validate request.

	id := uuid.New()
	if err := h.dbInsertCustomer(c, dbInsertCustomerArgs{
		id:               id,
		firstName:        req.FirstName,
		lastName:         req.LastName,
		email:            req.Email,
		residenceAddress: req.ResidenceAddress,
		birthDate:        req.BirthDate,
		kycStatus:        domain.KYCStatusPending,
	}); err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to insert customer in db: %w", err))
		return
	}

	go util.Recover(func() {
		if err := h.dbUpdateKYCStatus(c, dbUpdateKYCStatusArgs{
			id:        id,
			oldStatus: domain.KYCStatusPending,
			newStatus: domain.KYCStatusInProgress,
		}); err != nil {
			log.Println("failed to set KYC status in_progress", err)
			return
		}

		// Simulate KYC process.
		time.Sleep(time.Minute)

		if err := h.dbUpdateKYCStatus(c, dbUpdateKYCStatusArgs{
			id:        id,
			oldStatus: domain.KYCStatusInProgress,
			newStatus: domain.KYCStatusApproved,
		}); err != nil {
			log.Println("failed to set KYC status approved", err)
			return
		}

		// TODO notify customer.
	})

	c.JSON(http.StatusCreated, CreateCustomerResponse{ID: id})
}

type dbInsertCustomerArgs struct {
	id               uuid.UUID
	firstName        string
	lastName         string
	email            string
	residenceAddress string
	birthDate        time.Time
	kycStatus        domain.KYCStatus
}

func (h *CreateCustomer) dbInsertCustomer(ctx context.Context, args dbInsertCustomerArgs) error {
	sql := `
		INSERT INTO customers (id, first_name, last_name, email, residence_address, birth_date, kyc_status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := h.db.Exec(ctx, sql,
		args.id,
		args.firstName,
		args.lastName,
		args.email,
		args.residenceAddress,
		args.birthDate,
		args.kycStatus,
	)

	return err
}

type dbUpdateKYCStatusArgs struct {
	id        uuid.UUID
	oldStatus domain.KYCStatus
	newStatus domain.KYCStatus
}

func (h *CreateCustomer) dbUpdateKYCStatus(ctx context.Context, args dbUpdateKYCStatusArgs) error {
	sql := `
		UPDATE customers SET kyc_status = $1, updated_at = now()
		WHERE id = $2 AND kyc_status = $3`

	res, err := h.db.Exec(ctx, sql, args.newStatus, args.id, args.oldStatus)
	if err != nil {
		return err
	}
	if ra := res.RowsAffected(); ra != 1 {
		return fmt.Errorf("expected rows affected to be 1, got %d", ra)
	}

	return nil
}
