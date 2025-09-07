package api

import (
	"net/http"
	"time"

	db "github.com/Adebobola01/Simple-bank---GO/db/sqlc"
	"github.com/Adebobola01/Simple-bank---GO/util"
	"github.com/gin-gonic/gin"
)


type createUserRequest struct{
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type UserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func (server *Server) createUser(ctx *gin.Context){
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	args := db.CreateUsersParams{
		Username: req.Username,
		HashedPassword: hashedPassword,
		FullName: req.FullName,
		Email: req.Email,
	}
	user, err := server.store.CreateUsers(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	

resp := UserResponse{
	Username: user.Username,
	FullName: user.FullName,
	Email: user.Email,
	PasswordChangedAt: user.PasswordChangedAt,
	CreatedAt: user.CreatedAt,
}
	ctx.JSON(http.StatusOK, resp)

}