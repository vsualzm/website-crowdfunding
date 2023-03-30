package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/vsualzm/website-crowfunding/auth"
	"github.com/vsualzm/website-crowfunding/campaign"
	"github.com/vsualzm/website-crowfunding/handler"
	"github.com/vsualzm/website-crowfunding/helper"
	"github.com/vsualzm/website-crowfunding/transaction"
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
	db.AutoMigrate(&campaign.Campaign{})
	db.AutoMigrate(&campaign.CampaignImage{})
	db.AutoMigrate(&transaction.Transaction{})

	// inisiasi Repository, service , handler
	userRepository := user.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)
	transactionRepository := transaction.NewRepository(db)

	userService := user.NewService(userRepository)
	campaignService := campaign.NewService(campaignRepository)
	authService := auth.NewService()
	transactionService := transaction.NewService(transactionRepository, campaignRepository)

	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	// pemanggilan router API
	router := gin.Default()
	router.Static("/images", "./images")
	api := router.Group("/api/v1")

	// metode POST, GET, Update dan Delete
	// REST-API
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/session", userHandler.Login)
	api.POST("/email_chekers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)

	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaign)
	api.POST("/campaigns", authMiddleware(authService, userService), campaignHandler.CreateCampaign)
	api.PUT("/campaigns/:id", authMiddleware(authService, userService), campaignHandler.UpdateCampaign)
	api.POST("/campaign-images", authMiddleware(authService, userService), campaignHandler.UploadImage)

	api.GET("/campaigns/:id/transactions", authMiddleware(authService, userService), transactionHandler.GetCampaignTransactions)

	// router menjalankan GIN
	router.Run()

}

// middleware

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

// // campaigns, err := campaignRepository.FindAll()
// cekId, err := campaignRepository.FindByUserID(1)
// // fmt.Println(cekId)
// fmt.Println("debug")
// fmt.Println("debug")
// fmt.Println("debug")
// fmt.Println("debug")
// // panjangnya isinya campaign
// fmt.Println(len(cekId))
// for _, kempen := range cekId {
// 	fmt.Println(kempen.Name)
// 	if len(kempen.CampaignImages) > 0 {
// 		fmt.Println("jumlah gambar")
// 		fmt.Println(len(kempen.CampaignImages))
// 		fmt.Println(kempen.CampaignImages[0].FileName)
// 	}
// }
