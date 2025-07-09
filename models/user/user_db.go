package user

import (
	"github.com/link1st/gowebsocket/v2/lib/database"
	"log"
)

const (
	IsLogoffYes = 1
	IsLogoffNo  = 0
)

// User 用户
type UserPO struct {
	Id            int64  `json:"id" gorm:"primary_key" gorm:"column:id" form:"id"`                /** 自增主键 **/
	AccIp         string `json:"accIp" gorm:"column:acc_ip" form:"accIp"`                         /** acc Ip **/
	AccPort       string `json:"accPort" gorm:"column:acc_port" form:"accPort"`                   /** acc 端口 **/
	AppId         string `json:"appId" gorm:"column:app_id" form:"appId"`                         /** 平台号（类似房间号） **/
	UserId        string `json:"userId" gorm:"column:user_id" form:"userId"`                      /** 用户ID **/
	ClientIp      string `json:"clientIp" gorm:"column:client_ip" form:"clientIp"`                /** 客户端IP **/
	ClientPort    string `json:"clientPort" gorm:"column:client_port" form:"clientPort"`          /** 客户端端口 **/
	Qua           string `json:"qua" gorm:"column:qua" form:"qua"`                                /** qua **/
	DeviceInfo    string `json:"deviceInfo" gorm:"column:device_info" form:"deviceInfo"`          /** 设备信息 **/
	Status        int32  `json:"status" gorm:"column:status" form:"status"`                       /** 用户状态 1 在用 2 锁定 **/
	IsLogoff      int32  `json:"isLogoff" gorm:"column:is_logoff" form:"isLogoff"`                /** 是否登陆 0 否 1 是 **/
	LogoutTime    uint64 `json:"logoutTime" gorm:"column:logout_time" form:"logoutTime"`          /** 登出时间 **/
	AvatarUrl     string `json:"avatarUrl" gorm:"column:avatar_url" form:"avatarUrl"`             /** 用户头像链接 **/
	LoginTime     uint64 `json:"loginTime" gorm:"column:login_time" form:"loginTime"`             /** 登陆时间 **/
	HeartbeatTime uint64 `json:"heartbeatTime" gorm:"column:heartbeat_time" form:"heartbeatTime"` /** 上次心跳时间 **/
	CreateTime    uint64 `json:"createTime" gorm:"column:create_time" form:"createTime"`          /** 更新时间 **/
	Nickname      string `json:"nickname" gorm:"column:nickname" form:"nickname"`
}

func (d *UserPO) TableName() string {
	return "user"
}

func (d *UserPO) UserSave() (err error) {
	//判断是否存在，不存在则新增，存在则更新
	po, err := d.UserByUidAndAppId(d.UserId, d.AppId)
	if err != nil {
		err = database.GetDB().Model(d).Table(d.TableName()).Save(d).Error
		return err
	}
	// 更新信息
	d.Id = po.Id
	err = d.Update(*d)
	return err
}

func (d *UserPO) UserByUidAndAppId(userId, appId string) (*UserPO, error) {
	res := &UserPO{}
	err := database.GetDB().Model(d).Table(d.TableName()).Where("user_id = ? and app_id = ?", userId, appId).First(res).Error
	if err != nil {
		log.Println("record not found", err)
		return nil, err
	}
	return res, nil
}

func (d *UserPO) Update(entity UserPO) error {
	err := database.GetDB().Model(d).Omit("id", "app_id", "user_id", "create_time").Where("id = ?", entity.Id).Updates(&entity).Error
	log.Println("更新信息", err)
	return err
}

func (d *UserPO) List(entity UserPO) (res []UserPO, err error) {
	err = database.GetDB().Model(d).Where(entity).Find(&res).Error
	return res, err
}

func (d *UserPO) Delete(id int64) (err error) {
	err = database.GetDB().Model(d).Where("id = ?", id).Delete(&UserPO{}).Error
	return err
}

func (d *UserPO) UserByUserId(userId string) (*UserPO, error) {
	res := &UserPO{}
	err := database.GetDB().Model(d).Table(d.TableName()).Where("user_id = ?", userId).First(res).Error
	if err != nil {
		log.Println("record not found", err)
		return nil, err
	}
	return res, nil
}
