package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "simplebank/db/sqlc"
)

// Server is responsible for handling all http requests to this app
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new http server
func NewServer(s db.Store) Server {
	server := Server{store: s}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	router.POST("/accounts", server.createAccount)
	router.PUT("/accounts", server.updateAccount)
	router.GET("/accounts", server.getAccounts)
	router.GET("/accounts/:id", server.getAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)

	router.POST("/transfers", server.createTransfer)

	server.router = router
	return server
}

// Start runs the server on a specific http address
func (s *Server) Start(adrr string) error {
	return s.router.Run(adrr)
}

func errResponse(e error) gin.H {
	return gin.H{"error": e.Error()}
}
