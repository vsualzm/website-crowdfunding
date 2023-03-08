package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vsualzm/website-crowfunding/helper"
	"github.com/vsualzm/website-crowfunding/user"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	// tangkap input dari user
	// map input dari user ke struct RegisterUserInput
	// struct di atas kita passing sebagai parameter service

	// input dariregister
	var input user.RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Register Account Failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)

		return
	}

	newUser, err := h.userService.RegisterUserInput(input)

	if err != nil {
		response := helper.APIResponse("Register Account Failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// token, err := h.jwtService.GenerateToken()

	formatter := user.FormatUser(newUser, "TOKEN DI SINI ")

	// helper output
	response := helper.APIResponse("Account has been registered", http.StatusOK, "succes", formatter)
	c.JSON(http.StatusOK, response)
}
