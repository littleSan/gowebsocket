package friend

import (
	"github.com/link1st/gowebsocket/v2/lib/database"
	user2 "github.com/link1st/gowebsocket/v2/models/user"
	"time"
)

// Friend 聊天好友列表
type Friend struct {
	Id           int64  `json:"id" gorm:"primary_key" gorm:"column:id" form:"id"`              /** 自增主键 **/
	UserId       string `json:"userId" gorm:"column:user_id" form:"userId"`                    /** 用户ID **/
	FriendUserId string `json:"friendUserId" gorm:"column:friend_user_id" form:"friendUserId"` /** 朋友userID **/
	AppId        string `json:"appId" gorm:"column:app_id" form:"appId"`                       /** 平台号（类似房间号） **/
	ClientIp     string `json:"clientIp" gorm:"column:client_ip" form:"clientIp"`              /** 客户端IP **/
	Qua          string `json:"qua" gorm:"column:qua" form:"qua"`                              /** qua **/
	Nickname     string `json:"nickname" gorm:"column:nickname" form:"nickname"`               /** 昵称 **/
	AvatarUrl    string `json:"avatarUrl" gorm:"column:avatar_url" form:"avatarUrl"`           /** 用户头像链接 **/
	LastWords    string `json:"lastWords" gorm:"column:last_words" form:"lastWords"`           /** 最近一次聊天内容 **/
	LastTime     int64  `json:"lastTime" gorm:"column:last_time" form:"lastTime"`              /** 最近一次聊天时间 **/
	CreateTime   int64  `json:"createTime" gorm:"column:create_time" form:"createTime"`        /** 创建时间 **/
}

func (d *Friend) TableName() string {
	return "friend"
}

func (d *Friend) FriendSave(entity *Friend) (err error) {
	//存在数据，不进行添加

	res, err := d.GetByUid(entity.UserId, entity.FriendUserId)
	if err == nil && res.Id > 0 {
		return d.UpdateFriendInfo(entity.UserId, entity.FriendUserId, entity.LastWords)
	} else {
		user := user2.UserPO{}
		user.UserId = entity.FriendUserId
		user.AppId = entity.AppId
		res2, err := user.UserByUidAndAppId(user.UserId, user.AppId)
		if err == nil && res2.Id > 0 {
			entity.Nickname = res2.Nickname
			entity.AvatarUrl = res2.AvatarUrl
		}
		user.UserId = entity.UserId
		res3, err := user.UserByUidAndAppId(user.UserId, user.AppId)
		if err == nil && res3.Id > 0 {
			entity.ClientIp = res3.ClientIp
		}
		err = database.GetDB().Model(d).Table(d.TableName()).Save(entity).Error
	}
	return err
}

func (d *Friend) GetByUid(userId, friendUserId string) (*Friend, error) {
	res := &Friend{}
	err := database.GetDB().Model(d).Where("user_id = ? and friend_user_id = ?", userId, friendUserId).First(&res).Error
	return res, err
}

func (d *Friend) UpdateFriendInfo(userId, friendUserId, lastWords string) error {
	err := database.GetDB().Model(d).Omit("id", "user_id", "friend_user_id", "app_id").
		Where("user_id = ? and friend_user_id = ?", userId, friendUserId).UpdateColumns(map[string]interface{}{
		"last_words": lastWords,
		"last_time":  time.Now().UnixMilli(),
	}).Error
	return err
}

func (d *Friend) List(entity Friend) (res []Friend, err error) {
	err = database.GetDB().Model(d).Where(entity).Order("last_time desc").Find(&res).Error
	return res, err
}

func (d *Friend) Delete(id int64) (err error) {
	err = database.GetDB().Model(d).Where("id = ?", id).Delete(&Friend{}).Error
	return err
}
