package srv

import (
	"tikapp/api"
	"tikapp/common/db"
	"tikapp/common/log"
	"tikapp/common/model"
	"tikapp/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var logger = log.NameSpace("UserService")

type User struct{}

// 根据 Uber 的指导原则 这里是检查 User 是否实现了 api 中的所有方法
// 即检查项目是否缺少必要的接口
var _ api.UserHandler = (*User)(nil)

type UserLoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginResp struct {
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

func (u User) Login(c *gin.Context) (interface{}, error) {
	var req UserLoginReq
	err := c.ShouldBindWith(&req, binding.JSON)
	if err != nil {
		logger.Error("parse json error")
		return nil, err
	}
	var user model.User
	var count int64
	err = db.MySQL.Debug().Model(&model.User{}).Where("username = ? and password = ?", req.Username, req.Password).Select("id").First(&user).Count(&count).Error
	if err != nil {
		logger.Error("mysql happen error")
		return nil, err
	}
	if count != 1 {
		logger.Error("no user or user repeat")
		return nil, err
	}
	token, err := util.CreateAccessToken(user.Id)
	if err != nil {
		logger.Error("create access token error")
		return nil, err
	}
	refreshToken, err := util.CreateRefreshToken(user.Id)
	if err != nil {
		logger.Error("create refresh token error")
		return nil, err
	}
	c.Header("token", token)
	db.Redis.Set(token, refreshToken, 30*24*time.Hour)
	return UserLoginResp{
		UserId: user.Id,
		Token:  token,
	}, nil
}
