package api

import (
	"github.com/gin-gonic/gin"
	"https://github.com/Adebobola01/Simple-bank---GO/tree/master/db/sqlc"
)

type Server struct {
	store  *sqlc.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *server {
	server := &Server{store: store}
	router: gin.default()

	router.POST("/accounts", server.createAccout)
	server.router = router
	return server
}

func errorResponse (err error) gin.H{
	return gin.H{"error": err.Error};
}