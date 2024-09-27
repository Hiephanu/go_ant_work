package server

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Thay đổi thành URL frontend của bạn
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           86400, // Thời gian cache preflight request
	}
	// Sử dụng middleware CORS
	r.Use(cors.New(corsConfig))

	r.GET("ws", s.UpgradeToWS)
	//register websocket connection
	api := r.Group("/api/v1")
	{
		api.GET("/", s.HelloWorldHandler)
		api.GET("/health", s.healthHandler)
		api.GET("/health/redis", s.healthHandlerRedis)
		api.POST("/register", s.RegisterHandler)
		api.POST("/login", s.LoginHandler)
		api.GET("/users/:id", s.GetUserByIdHandler)
		api.PUT("/users/:id", s.UpdateUserHandler)
		api.POST("/rooms/create", s.CreateRoomHandler)
		api.POST("/rooms/join", s.JoinRoomHandler)
		api.POST("/rooms/left", s.LeafRoomHandler)
		api.GET("/rooms/all", s.GetAllRoom)
	}

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}

func (s *Server) healthHandlerRedis(c *gin.Context) {
	c.JSON(http.StatusOK, s.redis.Health())
}
