package handler

import "github.com/gin-gonic/gin"

func NewCreateAccount() *CreateAccount {
	return &CreateAccount{}
}

type CreateAccount struct {
	// deps
}

type CreateAccountReqJSON struct {
	FirstName string `json:"first_name"`
	// ...
}

func (h *CreateAccount) Handle(c *gin.Context) {
	panic("not implemented")
}
