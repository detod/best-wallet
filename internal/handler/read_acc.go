package handler

import "github.com/gin-gonic/gin"

func NewReadAccount() *ReadAccount {
	return &ReadAccount{}
}

type ReadAccount struct {
	// deps
}

func (h *ReadAccount) Handle(c *gin.Context) {
	panic("not implemented")
}
