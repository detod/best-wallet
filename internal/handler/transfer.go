package handler

import "github.com/gin-gonic/gin"

func NewTransfer() *Transfer {
	return &Transfer{}
}

type Transfer struct {
	// deps
}

func (h *Transfer) Handle(c *gin.Context) {
	panic("not implemented")
}
