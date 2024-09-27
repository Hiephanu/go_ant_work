package server

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) GetUserByIdHandler(c *gin.Context) {
	userId, _ := c.Params.Get("id")

	res, err := s.GetUserById(c.Request.Context(), userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration success", "data": res})
}

func (s *Server) UpdateUserHandler(c *gin.Context) {
	var userUpdateRequest UserUpdateRequest
	err := json.NewDecoder(c.Request.Body).Decode(&userUpdateRequest)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	res, err := s.UpdateUser(c.Request.Context(), userUpdateRequest)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login sucess", "data": res})
}
