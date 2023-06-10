package main

import (
	"real-chat-backend/controllers"
	"real-chat-backend/middlewares"
	"real-chat-backend/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	models.ConnectDataBase()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"*"},
		AllowHeaders:  []string{"*"},
		AllowWildcard: true,
	}))

	go controllers.Manager.Start()

	public := r.Group("/api")

	public.POST("/register", controllers.Register)
	public.POST("/login", controllers.Login)

	protected := r.Group("/api")
	protected.Use(middlewares.JwtAuthMiddleware())
	protected.GET("/user", controllers.CurrentUser)
	protected.GET("/getUsers", controllers.GetUsers)
	protected.POST("/saveMessage", controllers.SaveMessage)
	protected.GET("/getMessage", controllers.GetMessage)

	public.GET("/ws", controllers.WsHandler)

	r.Run(":8080")
}
