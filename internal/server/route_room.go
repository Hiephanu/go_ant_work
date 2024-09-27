package server

import (
	"encoding/json"
	"fmt"
	"go_ant_work/internal/structs"
	"go_ant_work/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) CreateRoomHandler(c *gin.Context) {
	var room structs.CreateRoomRequest
	err := json.NewDecoder(c.Request.Body).Decode(&room)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		fmt.Print("Error to decode body in create room", err)
		return
	}

	valid := utils.ValidatePasswordRoom(room.Password)
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
		return
	}

	roomRes := s.roomTracker.CreateRoom(room.Name, room.Password, room.UserID)

	c.JSON(http.StatusOK, gin.H{"message": "Sucess", "data": roomRes})
}

func (s *Server) JoinRoomHandler(c *gin.Context) {
	var roomJoinRequest structs.RoomJoinRequest
	err := json.NewDecoder(c.Request.Body).Decode(&roomJoinRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	if !s.roomTracker.IsRoomExist(roomJoinRequest.RoomID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room not found"})
		return
	}

	s.roomTracker.JoinRoom(roomJoinRequest.RoomID, roomJoinRequest.UserID, roomJoinRequest.Password, roomJoinRequest.Offer)

	c.JSON(http.StatusOK, gin.H{"message": "Success", "data": roomJoinRequest.RoomID})
}

func (s *Server) LeafRoomHandler(c *gin.Context) {
	var leftRoomRequest structs.LeftRoomRequest
	err := json.NewDecoder(c.Request.Body).Decode(&leftRoomRequest)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	s.roomTracker.LeftRoom(leftRoomRequest.RoomID, leftRoomRequest.UserID)

	c.JSON(http.StatusOK, gin.H{"messgae": "Success", "data": leftRoomRequest.RoomID})
}
