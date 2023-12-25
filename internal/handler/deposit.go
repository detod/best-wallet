package handler

import "github.com/gin-gonic/gin"

func NewDeposit() *Deposit {
	return &Deposit{}
}

type Deposit struct {
	// deps
}

func (h *Deposit) Handle(c *gin.Context) {
	panic("not implemented")
}
