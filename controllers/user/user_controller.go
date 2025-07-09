// Package user 用户调用接口
package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/link1st/gowebsocket/v2/common"
	"github.com/link1st/gowebsocket/v2/controllers"
	"github.com/link1st/gowebsocket/v2/lib/cache"
	"github.com/link1st/gowebsocket/v2/models"
	"github.com/link1st/gowebsocket/v2/models/friend"
	"github.com/link1st/gowebsocket/v2/models/message"
	"github.com/link1st/gowebsocket/v2/servers/websocket"
	"strconv"
	"strings"
)

// List 查看全部在线用户
func List(c *gin.Context) {
	appIDStr := c.Query("appID")
	appIDUint64, _ := strconv.ParseInt(appIDStr, 10, 32)
	appID := uint32(appIDUint64)
	fmt.Println("http_request 查看全部在线用户", appID)
	data := make(map[string]interface{})
	userList := websocket.UserList(appID)
	data["userList"] = userList
	data["userCount"] = len(userList)
	controllers.Response(c, common.OK, "", data)
}

// Online 查看用户是否在线
func Online(c *gin.Context) {
	userID := c.Query("userID")
	appIDStr := c.Query("appID")
	appIDUint64, _ := strconv.ParseInt(appIDStr, 10, 32)
	appID := uint32(appIDUint64)
	fmt.Println("http_request 查看用户是否在线", userID, appIDStr)
	data := make(map[string]interface{})
	online := websocket.CheckUserOnline(appID, userID)
	data["userID"] = userID
	data["online"] = online
	controllers.Response(c, common.OK, "", data)
}

type SendMessageRequest struct {
	AppID     int    `json:"appID"`
	From      string `json:"from"`
	MsgID     string `json:"msgID"`
	Msg       string `json:"msg"`
	Target    string `json:"target"`
	GroupUuid string `json:"groupUuid"`
}

// SendMessage 给用户发送消息
func SendMessage(c *gin.Context) {
	// 获取参数
	reqParam := SendMessageRequest{}
	err := c.ShouldBindJSON(&reqParam)
	if err != nil {
		controllers.Response(c, common.ParameterIllegal, err.Error(), nil)
		return
	}
	fromUserId := strings.ToLower(reqParam.From)
	msgID := reqParam.MsgID
	//msgID := strconv.FormatInt(time.Now().UnixMilli(), 10)
	msg := reqParam.Msg
	appID := reqParam.AppID
	targetUserId := strings.ToLower(reqParam.Target)
	fmt.Println("http_request 给用户发送消息", appID, fromUserId, msgID, msg)

	// TODO::进行用户权限认证，一般是客户端传入TOKEN，然后检验TOKEN是否合法，通过TOKEN解析出来用户ID
	// 本项目只是演示，所以直接过去客户端传入的用户ID(userID)
	data := make(map[string]interface{})
	if cache.SeqDuplicates(msgID) {
		fmt.Println("给用户发送消息 重复提交:", msgID)
		controllers.Response(c, common.OK, "", data)
		return
	}
	sendResults, err := websocket.SendUserMessage(uint32(appID), targetUserId, msgID, msg, fromUserId)
	if err != nil {
		data["sendResultsErr"] = err.Error()
	}
	data["sendResults"] = sendResults
	controllers.Response(c, common.OK, "", data)
}

// SendMessageAll 给全员发送消息
func SendMessageAll(c *gin.Context) {
	// 获取参数
	// 获取参数
	reqParam := SendMessageRequest{}
	err := c.ShouldBindJSON(&reqParam)
	if err != nil {
		controllers.Response(c, common.ParameterIllegal, err.Error(), nil)
		return
	}

	fmt.Println("http_request 给全体用户发送消息", reqParam.AppID, reqParam.From, reqParam.MsgID, reqParam.Msg)
	data := make(map[string]interface{})
	if cache.SeqDuplicates(reqParam.MsgID) {
		fmt.Println("给用户发送消息 重复提交:", reqParam.MsgID)
		controllers.Response(c, common.OK, "", data)
		return
	}
	sendResults, err := websocket.SendUserMessageAll(uint32(reqParam.AppID), reqParam.From, reqParam.MsgID, models.MessageCmdBroadcast, reqParam.Msg)
	if err != nil {
		data["sendResultsErr"] = err.Error()
	}
	data["sendResults"] = sendResults
	controllers.Response(c, common.OK, "", data)

}

type FriendListRequest struct {
	AppID  int    `json:"appID"`
	UserID string `json:"userID"`
}

// FriendList 聊天好友历史记录消息
func FriendList(c *gin.Context) {
	req := FriendListRequest{}
	err := c.BindQuery(&req)
	if err != nil {
		controllers.Response(c, common.ParameterIllegal, err.Error(), nil)
		return
	}
	fri := friend.Friend{}
	fri.UserId = req.UserID
	fri.AppId = strconv.FormatInt(int64(req.AppID), 10)
	friends, err := fri.List(fri)
	data := make(map[string]interface{})
	data["friends"] = friends
	controllers.Response(c, common.OK, "", data)
}

// MessagesList 查询两个用户之间的聊天记录
type MessagesListRequest struct {
	AppID        int    `json:"appID" form:"appID"`
	UserID       string `json:"userID" form:"userID"`
	FriendUserID string `json:"friendUserID" form:"friendUserID"`
	GroupUuid    string `json:"groupUuid" form:"groupUuid"`
	Page         int    `json:"page" form:"page"`
	Limit        int    `json:"limit" form:"limit"`
}

func MessagesList(c *gin.Context) {
	req := MessagesListRequest{}
	err := c.BindQuery(&req)
	if err != nil {
		controllers.Response(c, common.ParameterIllegal, err.Error(), nil)
		return
	}
	data := make(map[string]interface{})
	msg := message.Message{}
	msg.AppId = strconv.FormatInt(int64(req.AppID), 10)
	msg.From = req.UserID
	msg.Target = req.FriendUserID
	msgList, cnt, err := msg.LastMessage(msg.From, msg.Target, req.Page, req.Limit)
	data["msgList"] = msgList
	data["total"] = cnt
	controllers.Response(c, common.OK, "", data)

}

// 发送群组消息
func SendMessageGroup(c *gin.Context) {
	req := SendMessageRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		controllers.Response(c, common.ParameterIllegal, err.Error(), nil)
		return
	}
	fmt.Println("http_request 给群组用户发消息", req.AppID, req.From, req.MsgID, req.Msg, req.Target, req.GroupUuid)
	data := make(map[string]interface{})
	if cache.SeqDuplicates(req.MsgID) {
		fmt.Println("群消息 重复提交:", req.MsgID)
		controllers.Response(c, common.OK, "", data)
		return
	}
	sendResults, err := websocket.SendUserMessageGroup(uint32(req.AppID), req.From, req.GroupUuid, models.MessageCmdGroup, req.MsgID, req.Msg)
	if err != nil {
		data["sendResultsErr"] = err.Error()
	}
	if sendResults {
		msg := message.Message{}
		msg.UpdateStatus(req.MsgID, message.MsgSendStatusSuccess)
	}
	data["sendResults"] = sendResults
	controllers.Response(c, common.OK, "", data)
}

func GroupMessagesList(c *gin.Context) {
	req := MessagesListRequest{}
	err := c.BindQuery(&req)
	if err != nil {
		controllers.Response(c, common.ParameterIllegal, err.Error(), nil)
		return
	}
	data := make(map[string]interface{})
	msg := message.Message{}
	msg.AppId = strconv.FormatInt(int64(req.AppID), 10)
	msg.From = req.UserID
	msg.Target = req.FriendUserID
	msgList, cnt, err := msg.GroupMessage(req.GroupUuid, req.Page, req.Limit)
	data["msgList"] = msgList
	data["total"] = cnt
	page := req.Page
	if page == 0 {
		page = 1
	}
	limit := req.Limit
	if limit == 0 {
		limit = 10
	}
	data["page"] = page
	data["limit"] = limit
	controllers.Response(c, common.OK, "", data)

}
