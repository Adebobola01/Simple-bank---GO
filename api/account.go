package api

import (
	"fmt"
	"net/http"

	db "github.com/Adebobola01/Simple-bank---GO/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createAccountRequest struct{
	Owner string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=EUR USD"`
}

func (server *Server) createAccount(ctx *gin.Context){
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.CreateAccountParams{
		Owner: req.Owner,
		Currency: req.Currency,
		Balance: 0,
	}
	fmt.Println(args, req)
	account, err := server.store.CreateAccount(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)

}