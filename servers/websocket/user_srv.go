// Package websocket 处理
package websocket

import (
	"errors"
	"fmt"
	friend2 "github.com/link1st/gowebsocket/v2/models/friend"
	"github.com/link1st/gowebsocket/v2/models/message"
	user2 "github.com/link1st/gowebsocket/v2/models/user"
	"strconv"
	"time"

	"github.com/link1st/gowebsocket/v2/lib/cache"
	"github.com/link1st/gowebsocket/v2/models"
	"github.com/link1st/gowebsocket/v2/servers/grpcclient"

	"github.com/redis/go-redis/v9"
)

// UserList 查询所有用户
func UserList(appID uint32) (userList []string) {
	userList = make([]string, 0)
	currentTime := uint64(time.Now().Unix())
	servers, err := cache.GetServerAll(currentTime)
	if err != nil {
		fmt.Println("给全体用户发消息", err)
		return
	}
	for _, server := range servers {
		var (
			list []string
		)
		if IsLocal(server) {
			list = GetUserList(appID)
		} else {
			list, _ = grpcclient.GetUserList(server, appID)
		}
		userList = append(userList, list...)
	}
	return
}

// CheckUserOnline 查询用户是否在线
func CheckUserOnline(appID uint32, userID string) (online bool) {
	// 全平台查询
	if appID == 0 {
		for _, appID := range GetAppIDs() {
			online, _ = checkUserOnline(appID, userID)
			if online == true {
				break
			}
		}
	} else {
		online, _ = checkUserOnline(appID, userID)
	}
	return
}

// checkUserOnline 查询用户 是否在线
func checkUserOnline(appID uint32, userID string) (online bool, err error) {
	key := GetUserKey(appID, userID)
	userOnline, err := cache.GetUserOnlineInfo(key)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			fmt.Println("GetUserOnlineInfo", appID, userID, err)
			return false, nil
		}
		fmt.Println("GetUserOnlineInfo", appID, userID, err)
		return
	}
	online = userOnline.IsOnline()
	return
}

// SendUserMessage 给用户发送消息
func SendUserMessage(appID uint32, targetUserId string, msgID, content string, fromUserId string) (sendResults bool, err error) {
	tempMsg := models.MessageTrans{MsgId: msgID, Msg: content, From: fromUserId, Target: targetUserId}
	data := models.GetTextMsgData(tempMsg)
	client := GetUserClient(appID, targetUserId)
	// 保存消息发送本地记录
	SaveMessageToDb(fromUserId, targetUserId, strconv.FormatInt(int64(appID), 10), msgID, content, "")
	SaveFriendRelation(fromUserId, targetUserId, strconv.FormatInt(int64(appID), 10), content)
	if client != nil {
		sendResults, err = SendUserMessageLocal(appID, targetUserId, data)
		if err != nil {
			fmt.Println("给用户发送消息", appID, targetUserId, err)
		}
		if sendResults == true {
			// 发送成功
			UpdateMessage(msgID)
		}
		return
	}
	key := GetUserKey(appID, targetUserId)
	info, err := cache.GetUserOnlineInfo(key)
	if err != nil {
		cache.PushOfflineMsg(key, data)
		fmt.Println("给用户发送消息失败", key, err)
		return true, nil
	}
	if !info.IsOnline() {
		cache.PushOfflineMsg(key, data)
		fmt.Println("用户不在线", key)
		return true, nil
	}
	server := models.NewServer(info.AccIp, info.AccPort)
	msg, err := grpcclient.SendMsg(server, msgID, appID, fromUserId, targetUserId, models.MessageCmdMsg, models.MessageCmdMsg, content)
	if err != nil {
		fmt.Println("给用户发送消息失败", key, err)
		cache.PushOfflineMsg(key, data)
		return false, err
	}
	fmt.Println("给用户发送消息成功-rpc", msg)
	sendResults = true
	return
}

// SendUserMessageLocal 给本机用户发送消息
func SendUserMessageLocal(appID uint32, userID string, data string) (sendResults bool, err error) {
	client := GetUserClient(appID, userID)
	if client == nil {
		err = errors.New("用户不在线")
		return
	}
	// 发送消息
	client.SendMsg([]byte(data))
	sendResults = true
	return
}

// SendUserMessageAll 给全体用户发消息
func SendUserMessageAll(appID uint32, fromUserId string, msgID, cmd, message string) (sendResults bool, err error) {
	sendResults = true
	currentTime := uint64(time.Now().Unix())
	servers, err := cache.GetServerAll(currentTime)
	if err != nil {
		fmt.Println("给全体用户发消息", err)
		return
	}
	//存群发消息
	SaveMessageToDb(fromUserId, "", strconv.FormatInt(int64(appID), 10), msgID, message, "")
	for _, server := range servers {
		if IsLocal(server) {
			tempMsg := models.MessageTrans{MsgId: msgID, Msg: message, From: fromUserId, Cmd: cmd, Target: ""}
			data := models.GetMsgData(tempMsg)
			AllSendMessages(appID, fromUserId, data)
		} else {
			_, _ = grpcclient.SendMsgAll(server, msgID, appID, fromUserId, cmd, message)
		}
	}
	return
}

func SendUserMessageGroup(appID uint32, fromUserId, groupUuid, cmd, msgID, message string) (sendResults bool, err error) {
	sendResults = true
	currentTime := uint64(time.Now().Unix())
	servers, err := cache.GetServerAll(currentTime)
	if err != nil {
		fmt.Println("给用户群发消息", err)
		return
	}
	//存群发消息
	SaveMessageToDb(fromUserId, "", strconv.FormatInt(int64(appID), 10), msgID, message, groupUuid)
	for _, server := range servers {
		if IsLocal(server) {
			//获取 nickname  和 avaterUsr
			tempMsg := models.MessageTrans{MsgId: msgID, Msg: message, From: fromUserId, Target: "", GroupUuid: groupUuid}
			keys := fmt.Sprintf("%d_%s", appID, fromUserId)
			onlineUser, er := cache.GetUserOnlineInfo(keys)
			if er != nil && onlineUser == nil {
				//从数据库查询
				userCv := user2.UserPO{}
				u1, er := userCv.UserByUidAndAppId(fromUserId, strconv.FormatInt(int64(appID), 10))
				if er != nil {
					fmt.Println("查询用户信息失败", err)
					sendResults = false
					return
				}
				tempMsg.Nickname = u1.Nickname
				tempMsg.AvatarUrl = u1.AvatarUrl
			} else {
				tempMsg.Nickname = onlineUser.Nickname
				tempMsg.AvatarUrl = onlineUser.AvatarUrl
			}

			data := models.GetMsgDataGroup(tempMsg)
			groupSendMessages(data, appID, fromUserId, groupUuid)
		} else {
			_, _ = grpcclient.SendMsgAll(server, msgID, appID, fromUserId, cmd, message)
		}
	}
	return
}

func SaveMessageToDb(from, targetUserId, appId, msgID, content, groupUuid string) {
	msg := message.Message{}
	msg.SaveMsg(from, targetUserId, appId, msgID, "text", content, message.MsgSendStatusSending, message.MsgReadNo, groupUuid)
}

func SaveFriendRelation(from, targetUserId, appId, content string) {
	//维护列表关系
	fr := friend2.Friend{}
	fr.UserId = from
	fr.FriendUserId = targetUserId
	fr.AppId = appId
	fr.LastWords = content
	fr.FriendSave(&fr)
	fr.FriendUserId, fr.UserId = fr.UserId, fr.FriendUserId
	fr.Id = 0
	fr.FriendSave(&fr)
}

func UpdateMessage(msgID string) {
	msg := message.Message{}
	msg.UpdateStatus(msgID, message.MsgSendStatusSuccess)
}
