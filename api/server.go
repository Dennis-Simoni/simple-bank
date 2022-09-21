package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "simplebank/db/sqlc"
	"simplebank/token"
	"simplebank/util"
)

// Server is responsible for handling all http requests to this app
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new http server
func NewServer(config util.Config, s db.Store) (*Server, error) {

	fmt.Println(config.TokenSymmetricKey)
	maker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %v", err)
	}

	server := &Server{
		config:     config,
		store:      s,
		tokenMaker: maker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setUpRoutes()
	return server, nil
}

func (s *Server) setUpRoutes() {

	s.router = gin.Default()

	s.router.POST("/accounts", s.createAccount)
	s.router.PUT("/accounts", s.updateAccount)
	s.router.GET("/accounts", s.getAccounts)
	s.router.GET("/accounts/:id", s.getAccount)
	s.router.DELETE("/accounts/:id", s.deleteAccount)

	s.router.POST("/transfers", s.createTransfer)

	s.router.POST("/users", s.createUser)
	s.router.POST("/users/login", s.loginUser)
}

// Start runs the server on a specific http address
func (s *Server) Start(adrr string) error {
	return s.router.Run(adrr)
}

func errResponse(e error) gin.H {
	return gin.H{"error": e.Error()}
}
