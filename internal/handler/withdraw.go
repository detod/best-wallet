package handler

import "github.com/gin-gonic/gin"

func NewWithdraw() *Withdraw {
	return &Withdraw{}
}

type Withdraw struct {
	// deps
}

func (h *Withdraw) Handle(c *gin.Context) {
	panic("not implemented")
}
