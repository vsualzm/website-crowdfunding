package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/vsualzm/website-crowfunding/auth"
	"github.com/vsualzm/website-crowfunding/handler"
	"github.com/vsualzm/website-crowfunding/helper"
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

	// cek token dengan manual
	// token, err := authService.ValidateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxNn0.Cw3D1z9hd5PEK87wc0dahk47C2no1GvWJtPziFS3lAk")

	// if err != nil {
	// 	fmt.Println("ERROR")
	// 	fmt.Println("ERROR")
	// 	fmt.Println("ERROR")
	// }

	// if token.Valid {
	// 	fmt.Println("VALID")
	// 	fmt.Println("VALID")
	// 	fmt.Println("VALID")
	// } else {
	// 	fmt.Println("INVALID")
	// 	fmt.Println("INVALID")
	// 	fmt.Println("INVALID")
	// }
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
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)

	// router menjalankan GIN
	router.Run()

}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {

		// ambil Authorization
		authHeader := c.GetHeader("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)

			// semacam menghentikan operasi
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// tokennnn!!
		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_id"].(float64))
		user, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)
	}
}

// dalam middlware
// ambil nilai header Authorization : Bearer tokentokentoken
// dari header authorization, kita nilai token nya saja (maniuppulasi token nya saja)
// ktia validasi token
// ambil user_id
// ambil user dari db berdasakan user_id lewat service
// kita set context isinya user
