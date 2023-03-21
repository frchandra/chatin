package middleware

import (
	"github.com/frchandra/chatin/app/service"
	"github.com/frchandra/chatin/app/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type UserMiddleware struct {
	tokenUtil   *util.TokenUtil
	log         *logrus.Logger
	userService *service.UserService
}

func NewUserMiddleware(tokenUtil *util.TokenUtil, log *logrus.Logger, userService *service.UserService) *UserMiddleware {
	return &UserMiddleware{tokenUtil: tokenUtil, log: log, userService: userService}
}

func (u *UserMiddleware) HandleUserAccess(c *gin.Context) {
	accessDetails, err := u.tokenUtil.GetValidatedAccess(c) //get the user data from the token in the request header
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  "your credentials are invalid",
		})
		u.log.
			WithField("occurrence", "UserMiddelware@HandleUserAcccess").
			WithField("source_messages", err.Error()).
			WithField("client_ip", c.ClientIP()).
			WithField("endpoint", c.FullPath()).
			Info("cannot find token in the http request")
		c.Abort()
		return
	}
	err = u.tokenUtil.FetchAuthn(accessDetails.AccessUuid) //check if token exist in the token storage (Check if the token is expired)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "fail",
			"error":  "your credentials are invalid. try to refresh your credentials",
		})
		u.log.
			WithField("occurrence", "UserMiddleware@HandleUserAccess").
			WithField("client_ip", c.ClientIP()).
			WithField("endpoint", c.FullPath()).
			WithField("source_messages", err.Error()).
			Info("cannot find access token in the storage")
		c.Abort()
		return
	}
	c.Set("accessDetails", accessDetails)
	c.Next()
}

func (u *UserMiddleware) HandleAdminAccess(c *gin.Context) {
	accessDetails, err := u.tokenUtil.GetValidatedAccess(c) //get the user data from the token in the request header
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  "your credentials are invalid",
		})
		u.log.
			WithField("occurrence", "UserMiddelware@HandleUserAcccess").
			WithField("source_messages", err.Error()).
			WithField("client_ip", c.ClientIP()).
			WithField("endpoint", c.FullPath()).
			Info("cannot find token in the http request")
		c.Abort()
		return
	}

	err = u.tokenUtil.FetchAuthn(accessDetails.AccessUuid) //check if token exist in the token storage (Check if the token is expired)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "fail",
			"error":  "your credentials are invalid. try to refresh your credentials",
		})
		u.log.
			WithField("occurrence", "AdminMiddelware@HandleAdminAcccess").
			WithField("client_ip", c.ClientIP()).
			WithField("endpoint", c.FullPath()).
			WithField("source_messages", err.Error()).
			Info("cannot find access token in the storage")
		c.Abort()
		return
	}

	userResult, _ := u.userService.GetOneById(accessDetails.UserId)
	if userResult.Role == "admin" { //check is this user is an admin
		c.Set("accessDetails", accessDetails)
		c.Next()
		return
	}
	c.Abort()
	c.JSON(http.StatusUnauthorized, gin.H{
		"status": "fail",
		"error":  "you are not authorized",
	})
	return
}
