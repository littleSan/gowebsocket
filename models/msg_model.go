// Package models 数据模型
package models

import (
	"github.com/link1st/gowebsocket/v2/common"
	"time"
)

const (
	// MessageTypeText 文本类型消息
	MessageTypeText = "text"
	// MessageCmdMsg 文本类型消息
	MessageCmdMsg = "msg"
	// MessageCmdEnter 用户进入类型消息
	MessageCmdEnter = "enter"
	// MessageCmdExit 用户退出类型消息
	MessageCmdExit = "exit"
	// 群组消息
	MessageCmdGroup = "group"
	//广播消息
	MessageCmdBroadcast = "broadcast"
)

// Message 消息的定义
type Message struct {
	Target    string `json:"target"` // 目标
	Type      string `json:"type"`   // 消息类型 text/img/
	Msg       string `json:"msg"`    // 消息内容
	From      string `json:"from"`   // 发送者,
	Time      int64  `json:"time"`
	GroupUuid string `json:"groupUuid"` // 群组uuid
	Nickname  string `json:"nickname"`  //  昵称
	AvatarUrl string `json:"avatarUrl"` //  头像
}

type MessageTrans struct {
	Target    string `json:"target"` // 目标
	Type      string `json:"type"`   // 消息类型 text/img/
	Msg       string `json:"msg"`    // 消息内容
	From      string `json:"from"`   // 发送者,
	Time      int64  `json:"time"`
	GroupUuid string `json:"groupUuid"` // 群组uuid
	Nickname  string `json:"nickname"`  //  昵称
	AvatarUrl string `json:"avatarUrl"` //  头像
	Cmd       string `json:"cmd"`       // 命令
	MsgId     string `json:"msgId"`     // 消息id
}

// NewMsg 创建新的消息
func NewMsg(temp MessageTrans) (message *Message) {
	message = &Message{
		Type:      MessageTypeText,
		From:      temp.From,
		Msg:       temp.Msg,
		Target:    temp.Target,
		Time:      time.Now().UnixMilli(),
		GroupUuid: temp.GroupUuid,
		Nickname:  temp.Nickname,
		AvatarUrl: temp.AvatarUrl,
	}
	return
}

func getTextMsgData(temp MessageTrans) string {

	textMsg := NewMsg(temp)
	head := NewResponseHead(temp.MsgId, temp.Cmd, common.OK, "Ok", textMsg)

	return head.String()
}

// GetMsgData 文本消息
func GetMsgData(temp MessageTrans) string {
	return getTextMsgData(temp)
}

// GetTextMsgData 文本消息
func GetTextMsgData(temp MessageTrans) string {
	temp.Cmd = MessageCmdMsg
	return getTextMsgData(temp)
}

// GetTextMsgDataEnter 用户进入消息
func GetTextMsgDataEnter(temp MessageTrans) string {
	temp.Cmd = MessageCmdEnter
	return getTextMsgData(temp)
}

// GetTextMsgDataExit 用户退出消息
func GetTextMsgDataExit(temp MessageTrans) string {
	temp.Cmd = MessageCmdExit
	return getTextMsgData(temp)
}

// GetMsgDataGroup 用户群组消息
func GetMsgDataGroup(temp MessageTrans) string {
	temp.Cmd = MessageCmdGroup
	return getTextMsgData(temp)
}
