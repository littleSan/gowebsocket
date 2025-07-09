// Package routers 路由
package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/link1st/gowebsocket/v2/controllers/fileControl"
	"github.com/link1st/gowebsocket/v2/controllers/group"
	"github.com/link1st/gowebsocket/v2/controllers/home"
	"github.com/link1st/gowebsocket/v2/controllers/systems"
	"github.com/link1st/gowebsocket/v2/controllers/user"
)

// Init http 接口路由
func Init(router *gin.Engine) {
	router.LoadHTMLGlob("views/**/*")

	// 用户组
	userRouter := router.Group("/user")
	{
		userRouter.GET("/list", user.List)
		userRouter.GET("/online", user.Online)
		userRouter.POST("/sendMessage", user.SendMessage)
		userRouter.POST("/sendMessageAll", user.SendMessageAll)
		userRouter.POST("/sendMessageGroup", user.SendMessageGroup)
		userRouter.GET("/friends", user.FriendList)
		userRouter.GET("/messages", user.MessagesList)
	}

	// 系统
	systemRouter := router.Group("/system")
	{
		systemRouter.GET("/state", systems.Status)
	}

	// home
	homeRouter := router.Group("/home")
	{
		homeRouter.GET("/index", home.Index)
	}
	fileRouter := router.Group("/file")
	{
		fileRouter.POST("/upload", fileControl.UploadPic)
	}

	//群组
	groupRouter := router.Group("/group")
	{
		groupRouter.POST("/create", group.CreateGroup)
		groupRouter.POST("/addMember", group.AddGroupMember)
		groupRouter.POST("/deleteMember", group.DeleteGroupMember)
		groupRouter.GET("/list", group.GroupList)
		groupRouter.GET("/memberList", group.GroupMemberList)
		groupRouter.GET("/messages", user.GroupMessagesList)
	}
}
