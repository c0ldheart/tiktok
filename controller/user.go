package ctrl

import (
	"github.com/gin-gonic/gin"
	"strconv"
	res "tikapp/common/result"
	srv "tikapp/service"
	"tikapp/util"
)

// Register 新用户注册
func Register(c *gin.Context) {
	var u srv.User
	register, err := u.Register(c)
	if err != nil {
		// 用户名重复
		if err == srv.ErrUsernameExits {
			res.Error(c, res.Status{
				StatusCode: res.UsernameExitErrorStatus.StatusCode,
				StatusMsg:  res.UsernameExitErrorStatus.StatusMsg,
			})
		} else if err == srv.ErrEmpty {
			res.Error(c, res.Status{
				StatusCode: res.EmptyErrorStatus.StatusCode,
				StatusMsg:  res.EmptyErrorStatus.StatusMsg,
			})
		} else {
			res.Error(c, res.Status{
				StatusCode: res.RegisterErrorStatus.StatusCode,
				StatusMsg:  res.RegisterErrorStatus.StatusMsg,
			})
		}
		return
	}
	data := register.(srv.UserRegisterResp)
	res.Success(c, res.R{
		"userid": data.UserId,
		"token":  data.Token,
	})
}

// Info 获取用户信息
func Info(c *gin.Context) {
	var u srv.User
	var myUserID int64
	var err error
	targetUserID, _ := strconv.Atoi(c.Query("user_id"))
	token := c.Query("token")

	// 通过token获取当前用户ID，如果是游客（token为空），则当前用户ID为0
	if token != "" {
		myUserID, err = util.GetUsernameFormToken(token)
		if err != nil {
			res.Error(c, res.Status{
				StatusCode: res.TokenErrorStatus.StatusCode,
				StatusMsg:  "token error",
			})
			return
		}
	}

	// 调用服务层
	user, err := u.Info(myUserID, int64(targetUserID))
	if err != nil {
		res.Error(c, res.Status{
			StatusCode: res.InfoErrorStatus.StatusCode,
			StatusMsg:  "info error",
		})
		return
	}

	// 因为看文档返回时user要打包，所以这里也要打包
	res.Success(c, res.R{
		"user": user,
	})
}

// Login 用户登录
func Login(c *gin.Context) {
	var u srv.User
	login, err := u.Login(c)
	if err != nil {
		res.Error(c, res.Status{
			StatusCode: res.LoginErrorStatus.StatusCode,
			StatusMsg:  res.LoginErrorStatus.StatusMsg,
		})
		return
	}
	data := login.(srv.UserLoginResp)
	res.Success(c, res.R{
		"user_id": data.UserId,
		"token":   data.Token,
	})
}
