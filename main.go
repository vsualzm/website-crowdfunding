package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/vsualzm/website-crowfunding/auth"
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

	authService := auth.NewService()

	token, err := authService.ValidateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxNn0.Cw3D1z9hd5PEK87wc0dahk47C2no1GvWJtPziFS3lAk")

	if err != nil {
		fmt.Println("ERROR")
		fmt.Println("ERROR")
		fmt.Println("ERROR")
	}

	if token.Valid {
		fmt.Println("VALID")
		fmt.Println("VALID")
		fmt.Println("VALID")
	} else {
		fmt.Println("INVALID")
		fmt.Println("INVALID")
		fmt.Println("INVALID")
	}
	// manual
	//userService.SaveAvatar(6, "images/20-profile.png")

	userHandler := handler.NewUserHandler(userService, authService)

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
