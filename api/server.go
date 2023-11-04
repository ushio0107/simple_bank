package api

import (
	db "simple_bank/db/sqlc"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	accountsGroup := router.Group("/accounts")
	{
		accountsGroup.POST("", server.createAccount)
		accountsGroup.GET("/:id", server.getAccount)
		accountsGroup.GET("", server.listAccounts)
		accountsGroup.DELETE("/:id", server.deleteAccount)
		accountsGroup.PUT("", server.updateAccount)
	}

	transferGroup := router.Group("/transfers")
	{
		transferGroup.POST("", server.createTransfer)

	}

	server.router = router
	return server
}

// Start runs the HTTP server on a specific address.
// router is seted to be private, that's why we have func Start to start a server.
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
