package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterHandler(c *gin.Context) {
	var registerRequest RegisterRequest
	err := json.NewDecoder(c.Request.Body).Decode(&registerRequest)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	res, err := s.Register(c.Request.Context(), registerRequest)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration success", "data": res})
}

func (s *Server) LoginHandler(c *gin.Context) {
	var loginRequest LoginRequest
	err := json.NewDecoder(c.Request.Body).Decode(&loginRequest)

	fmt.Println(loginRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	res, err := s.Login(c.Request.Context(), &loginRequest)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login sucess", "data": res})
}
