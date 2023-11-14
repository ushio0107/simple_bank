package api

import (
	"fmt"
	db "simple_bank/db/sqlc"
	"simple_bank/token"
	"simple_bank/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("Failed to create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	userGroup := router.Group("/users")
	{
		userGroup.POST("", server.createUser)
		userGroup.GET("/users/:id", server.getUser)
	}

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
	return server, nil
}

// Start runs the HTTP server on a specific address.
// router is seted to be private, that's why we have func Start to start a server.
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
