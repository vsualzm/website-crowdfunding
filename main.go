package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/vsualzm/website-crowfunding/handler"
	"github.com/vsualzm/website-crowfunding/user"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	// connection database
	dsn := "root:@tcp(127.0.0.1:3306)/startup_bwa?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// cheking database connection
	if err != nil {
		log.Fatal(err.Error())
	}

	// migrate database
	db.AutoMigrate(&user.User{})

	// inisiasi Repository, service , handler
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	// manual
	//userService.SaveAvatar(6, "images/20-profile.png")

	userHandler := handler.NewUserHandler(userService)

	// pemanggilan router API
	router := gin.Default()
	api := router.Group("/api/v1")

	// metode POST, GET, Update dan Delete
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/session", userHandler.Login)
	api.POST("/email_chekers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", userHandler.UploadAvatar)

	// router menjalankan GIN
	router.Run()

}
