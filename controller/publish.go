package ctrl

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	res "tikapp/common/result"
	srv "tikapp/service"
	"tikapp/util"
)

// PublishAction 已登录的用户上传视频
func PublishAction(c *gin.Context) {
	fmt.Println("进入publish")
	userId, _ := c.Get("userId")
	if userId == "" {
		res.Error(c, res.Status{
			StatusCode: res.NoLoginErrorStatus.StatusCode,
			StatusMsg:  res.NoLoginErrorStatus.StatusMsg,
		})
		return
	}
	title := c.PostForm("title")
	data, err := c.FormFile("data")

	if err != nil {
		res.Error(c, res.Status{
			StatusCode: res.FileErrorStatus.StatusCode,
			StatusMsg:  res.FileErrorStatus.StatusMsg,
		})
		return
	}
	var v srv.Video
	err = v.PublishAction(data, title, userId.(int64))
	if err != nil {
		res.Error(c, res.Status{
			StatusCode: res.PublishErrorStatus.StatusCode,
			StatusMsg:  res.PublishErrorStatus.StatusMsg,
		})
		return
	}
	res.Success(c, res.R{})
}

// PublishList 列出当前用户所有的投稿视频
func PublishList(c *gin.Context) {
	var myUserID int64
	var targetUserID int
	var err error

	token := c.Query("token")
	if token != "" {
		myUserID, err = util.GetUserIDFormToken(token)
	}

	targetUserID, _ = strconv.Atoi(c.Query("user_id"))
	if err != nil {
		res.Error(c, res.Status{
			StatusCode: res.TokenErrorStatus.StatusCode,
			StatusMsg:  "token error",
		})
		return
	}

	var v srv.Video
	list, err := v.PublishList(myUserID, int64(targetUserID))
	if err != nil {
		res.Error(c, res.Status{
			StatusCode: res.PublishErrorStatus.StatusCode,
			StatusMsg:  "获取视频列表错误",
		})
		return
	}
	res.Success(c, res.R{
		"video_list": list,
	})
}
