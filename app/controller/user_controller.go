package controller

import (
	"github.com/frchandra/chatin/app/model"
	"github.com/frchandra/chatin/app/service"
	"github.com/frchandra/chatin/app/util"
	"github.com/frchandra/chatin/app/validation"
	"github.com/frchandra/chatin/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type UserController struct {
	userService *service.UserService
	tokenUtil   *util.TokenUtil
	config      *config.AppConfig
	log         *util.LogUtil
}

func NewUserController(userService *service.UserService, tokenUtil *util.TokenUtil, config *config.AppConfig, log *util.LogUtil) *UserController {
	return &UserController{userService: userService, tokenUtil: tokenUtil, config: config, log: log}
}

func (u *UserController) Register(c *gin.Context) {

	var inputData validation.RegisterValidation //validate the input data
	if err := c.ShouldBindJSON(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}

	newUser := model.User{ //transform the data into user object
		Name:     inputData.Name,
		Email:    inputData.Email,
		Password: inputData.Password,
		Role:     "user",
	}

	userResult, err := u.userService.GetOrInsertOne(&newUser) //update or get the new data
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}

	token, err := u.userService.GenerateToken(&newUser) //generate token for this user
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}

	c.SetCookie("access_token", token.AccessToken, int(u.config.AccessMinute/time.Second), "/", "localhost", false, true)
	c.SetCookie("refresh_token", token.RefreshToken, int(u.config.AccessMinute/time.Second), "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"user":   userResult,
		"token":  token,
	})
	return
}

func (u *UserController) SignIn(c *gin.Context) {
	var inputData validation.RegisterValidation //validate the input data
	if err := c.ShouldBindJSON(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}

	newUser := model.User{ //transform the data into user object
		Name:     inputData.Name,
		Email:    inputData.Email,
		Password: inputData.Password,
		Role:     "user",
	}

	userResult, err := u.userService.InsertOne(&newUser) //insert the new user data
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

func (u *UserController) Login(c *gin.Context) {
	var inputData validation.LoginValidation //validate the input data
	if err := c.ShouldBindJSON(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}

	user := model.User{ //transform the input data to user object
		Email:    inputData.Email,
		Name:     inputData.Name,
		Password: inputData.Password,
	}

	resultUser, err := u.userService.ValidateLogin(&user) //validate if user exist and credential is correct
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}

	token, err := u.userService.GenerateToken(&resultUser) //generate token for this user
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}

	c.SetCookie("access_token", token.AccessToken, int(u.config.AccessMinute/time.Second), "/", "localhost", false, true)
	c.SetCookie("refresh_token", token.RefreshToken, int(u.config.AccessMinute/time.Second), "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"token":  token,
	})
	return

}

func (u *UserController) CurrentUser(c *gin.Context) {
	//get the details about the current user that make request from the context passed by user middleware
	contextData, isExist := c.Get("accessDetails")
	if isExist == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  "cannot get access details",
		})
		return
	}

	accessDetails, _ := contextData.(*util.AccessDetails) //get the user data given the user id from the token
	user, err := u.userService.GetOneById(accessDetails.UserId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   user,
	})
	return
}

func (u *UserController) Logout(c *gin.Context) {
	accessDetails, err := u.tokenUtil.GetValidatedAccess(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	deleted, err := u.tokenUtil.DeleteAuthn(accessDetails.AccessUuid)
	if err != nil || deleted == 0 { //if any goes wrong
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	c.SetCookie("access_token", "", -1, "", "", false, true)
	c.SetCookie("refresh_token", "", -1, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
	return
}

func (u *UserController) RefreshToken(c *gin.Context) {
	token, err := u.tokenUtil.Refresh(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"token":  token,
	})
	return
}
