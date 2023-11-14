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

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()

	return server, nil
}

// setupRouter sets up the router for the server.
func (s *Server) setupRouter() {
	s.router = gin.Default()

	userGroup := s.router.Group("/users")
	{
		userGroup.POST("", s.createUser)
		userGroup.GET("/users/login", s.loginUser)
		userGroup.GET("/users/:id", s.getUser)
	}

	accountsGroup := s.router.Group("/accounts")
	{
		accountsGroup.POST("", s.createAccount)
		accountsGroup.GET("/:id", s.getAccount)
		accountsGroup.GET("", s.listAccounts)
		accountsGroup.DELETE("/:id", s.deleteAccount)
		accountsGroup.PUT("", s.updateAccount)
	}

	transferGroup := s.router.Group("/transfers")
	{
		transferGroup.POST("", s.createTransfer)
	}
}

// Start runs the HTTP server on a specific address.
// router is seted to be private, that's why we have func Start to start a server.
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
