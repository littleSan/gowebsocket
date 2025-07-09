package message

import (
	"fmt"
	"github.com/link1st/gowebsocket/v2/lib/database"
	"time"
)

const (
	MsgSendStatusSuccess = 1
	MsgSendStatusFail    = 2
	MsgSendStatusSending = 3
	MsgReadYes           = 1
	MsgReadNo            = 2
)

// Message 消息列表
type Message struct {
	Id         int64  `json:"id" gorm:"primary_key" gorm:"column:id" form:"id"`       /** 自增主键 **/
	From       string `json:"from" gorm:"column:from" form:"from"`                    /** 消息发送方 **/
	Target     string `json:"target" gorm:"column:target" form:"target"`              /** 消息接受方 **/
	AppId      string `json:"appId" gorm:"column:app_id" form:"appId"`                /** 平台号（类似房间号） **/
	Seq        string `json:"seq" gorm:"column:seq" form:"seq"`                       /** 消息ID **/
	MsgType    string `json:"msgType" gorm:"column:msg_type" form:"msgType"`          /** 消息类型 **/
	Content    string `json:"content" gorm:"column:content" form:"content"`           /** 消息内容 **/
	Status     int32  `json:"status" gorm:"column:status" form:"status"`              /** 发送状态 1 成功 2 失败 3 发送中  **/
	IsRead     int32  `json:"isRead" gorm:"column:is_read" form:"isRead"`             /** 读取状态 1 已读 2 未读 **/
	CreateTime int64  `json:"createTime" gorm:"column:create_time" form:"createTime"` /** 创建时间 **/
	GroupUuid  string `json:"groupUuid" gorm:"column:group_uuid" form:"groupUuid"`    /** 群组ID **/
}

type MessageVo struct {
	Message
	Nickname  string `json:"nickname" gorm:"column:nickname" form:"nickname"`
	AvatarUrl string `json:"avatarUrl" gorm:"column:avatar_url" form:"avatarUrl"`
}

type Resp struct {
	Seq      string `json:"seq"`
	Cmd      string `json:"cmd"`
	Response struct {
		Code    int    `json:"code"`
		CodeMsg string `json:"codeMsg"`
		Data    struct {
			Target    string `json:"target"`
			Type      string `json:"type"`
			Msg       string `json:"msg"`
			From      string `json:"from"`
			Time      int64  `json:"time"`
			GroupUuid string `json:"groupUuid"`
			Nickname  string `json:"nickname"`
			AvatarUrl string `json:"avatarUrl"`
		} `json:"data"`
	} `json:"response"`
}

func (d *Message) TableName() string {
	return "message"
}

func (d *Message) SaveMsg(from, target, appid, seq, msgType, content string, status, isRead int, groupUuid string) (err error) {
	entity := &Message{
		From:       from,
		Target:     target,
		AppId:      appid,
		Seq:        seq,
		MsgType:    msgType,
		Content:    content,
		Status:     int32(status),
		IsRead:     int32(isRead),
		CreateTime: time.Now().UnixMilli(),
		GroupUuid:  groupUuid,
	}

	err = database.GetDB().Model(d).Save(entity).Error
	return err
}

func (d *Message) GetById(id int64) (res *Message, err error) {
	err = database.GetDB().Model(d).Where("id = ?", id).First(&res).Error
	if err != nil {
		fmt.Sprintf("record not found %s", err)
		return nil, err
	}
	return
}

func (d *Message) GetBySeq(seq string) (*Message, error) {
	res := &Message{}
	err := database.GetDB().Model(d).Where("seq = ?", seq).First(res).Error
	return res, err
}
func (d *Message) UpdateStatus(seq string, status int) error {
	err := database.GetDB().Model(d).Omit("id", "from", "target", "app_id", "seq").Where("seq = ?", seq).
		Update("status", status).Error
	return err
}

func (d *Message) List(entity Message) (res []Message, err error) {
	err = database.GetDB().Model(d).Where(entity).Find(&res).Error
	return res, err
}

// LastMessage 查询用户之间的消息
func (d *Message) LastMessage(from, to string, page, size int) (res []Message, cnt int64, err error) {
	db := database.GetDB().Model(d).Where("`from` = ? and `target` = ? ", from, to).
		Or("`from` = ? and `target` = ? ", to, from)
	if page <= 0 {
		page = 1
	}
	if size <= 0 || page > 50 {
		size = 10
	}
	err = db.Count(&cnt).Error
	err = db.Offset((page - 1) * size).Limit(size).Order("create_time desc").Find(&res).Error
	return
}

func (d *Message) Delete(seq string) (err error) {
	err = database.GetDB().Model(d).Where("seq = ?", seq).Delete(&Message{}).Error
	return err
}

// LastMessage 查询用户之间的消息
func (d *Message) GroupMessage(groupUuid string, page, size int) (res []MessageVo, cnt int64, err error) {
	db := database.GetDB().Model(d).Table("message m").Select("m.*,u.nickname,u.avatar_url").
		Joins("LEFT JOIN user u ON m.from = u.user_id").Where("m.`group_uuid` = ? ", groupUuid)

	if page <= 0 {
		page = 1
	}
	if size <= 0 || page > 50 {
		size = 10
	}
	err = db.Count(&cnt).Error
	err = db.Offset((page - 1) * size).Limit(size).Order("create_time desc").Find(&res).Error
	return
}
