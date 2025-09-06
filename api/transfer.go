package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/Adebobola01/Simple-bank---GO/db/sqlc"
	"github.com/gin-gonic/gin"
)


type createTransferRequest struct{
	FromAccountID int64 `json:"from_account_id" binding:"required"`
	ToAccountID int64 `json:"to_account_id" binding:"required"`
	Amount int64 `json:"amount" binding:"required"`
	Currency string `json:"currency" binding:"required"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req createTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil{
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !server.isValidAccount(ctx, req.FromAccountID, req.Currency){
		return
	}

	if !server.isValidAccount(ctx, req.ToAccountID, req.Currency){
		return
	}

	args := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID: req.ToAccountID,
		Amount: req.Amount,
	}

	transfer, err := server.store.TransferTx(ctx, args)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, transfer)

}

func (server *Server) isValidAccount(ctx *gin.Context, accountId int64, currency string) bool {
	account, err := server.store.GetAccount(ctx, accountId)
	if err != nil{
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false

	}

	if account.Currency != currency {
		err = fmt.Errorf("Error: Account %v currency mismatch %s instead of %s", account.ID, currency, account.Currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return  false
	}

	return true
}