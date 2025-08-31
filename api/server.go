package api

import (
	"github.com/gin-gonic/gin"
	"https://github.com/Adebobola01/Simple-bank---GO/tree/master/db/sqlc"
)

type Server struct {
	store  *sqlc.Store
	router *gin.Engine
}