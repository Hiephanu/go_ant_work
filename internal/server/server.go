package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"go_ant_work/internal/database"
)

type Server struct {
	port int

	db database.Service

	redis database.RedisService

	roomTracker *RoomTracker
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port:        port,
		db:          database.New(),
		redis:       database.NewRedisService(),
		roomTracker: NewTrackerRoom(),
	}

	handler := NewServer.RegisterRoutes()
	// handlerWithMiddelware := middleware.JwtAuthMiddleWare(handler)
	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      handler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
