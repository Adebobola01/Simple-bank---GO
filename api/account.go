package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type createAccountRequest struct{
	owner string `json:"owner" binding:"required`
	currency string `json:"currency" binding: "required oneof= EUR USD"`
}

func (server *Server) createAccount(ctx *gin.Context){
	var req createAccountRequest
	err := ctx.ShouldBindJSON(&req); err != nil{
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := createAccountParams{
		Owner: req.owner
		Currency: req.currency
		Balance: 0
	}
	account, err := Server.store.CreatAccount(ctx, args)
	if err := nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.statusOK, account)

}