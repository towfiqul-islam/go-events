package routes

import (
	"net/http"
	"strconv"

	"example.com/rest-api/models"
	"example.com/rest-api/utils"
	"github.com/gin-gonic/gin"
)

func signup(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse user data"})
		return
	}

	user.Password, err = utils.HashPassword(user.Password)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not hash password"})
		return
	}

	err = user.Save()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save user"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "User created"})
}

func login(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse data"})
		return
	}

	err = user.ValidateUser()

	if err != nil {
		context.JSON(http.StatusForbidden, err)
		return
	}

	token, err := utils.GenerateToken(user.Email, user.ID)

	if err != nil {
		context.JSON(http.StatusForbidden, gin.H{"message": "Could not authenticate user"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "login success", "token": token})
}

func getUserByID(context *gin.Context) {
	userId, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse userId"})
		return
	}

	user, err := models.GetUser(userId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not find user"})
	}

	context.JSON(http.StatusOK, user)

	
}