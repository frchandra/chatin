package controller

import (
	"github.com/frchandra/chatin/app/model"
	"github.com/frchandra/chatin/app/service"
	"github.com/frchandra/chatin/app/util"
	"github.com/frchandra/chatin/app/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
	userService *service.UserService
	log         *util.LogUtil
}

func NewUserController(userService *service.UserService, log *util.LogUtil) *UserController {
	return &UserController{userService: userService, log: log}
}

func (u *UserController) Register(c *gin.Context) {
	//validate the input data
	var inputData validation.RegisterValidation
	if err := c.ShouldBindJSON(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}
	//get or insert the new user data
	newUser := model.User{
		Name:     inputData.Name,
		Email:    inputData.Email,
		Password: inputData.Password,
	}
	userResult, err := u.userService.GetOrInsertOne(&newUser)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"user":   userResult,
	})
	return
}
