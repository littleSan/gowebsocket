package groupElement

import (
	"github.com/link1st/gowebsocket/v2/lib/database"
	"log"
)

// GroupElement 组成员信息
type GroupElement struct {
	Id         int64  `json:"id" gorm:"primary_key" gorm:"column:id" form:"id"`       /** 自增主键 **/
	GroupUuid  string `json:"groupUuid" gorm:"column:group_uuid" form:"groupUuid"`    /** 唯一ID **/
	AppId      string `json:"appId" gorm:"column:app_id" form:"appId"`                /** 平台号（类似房间号） **/
	UserId     string `json:"userId" gorm:"column:user_id" form:"userId"`             /** 成员信息 **/
	Status     int32  `json:"status" gorm:"column:status" form:"status"`              /** 用户状态 1 在用 2 锁定 **/
	CreateTime int64  `json:"createTime" gorm:"column:create_time" form:"createTime"` /** 加入时间 **/
}

type GroupElementVO struct {
	GroupUuid string `json:"groupUuid" gorm:"column:group_uuid" form:"groupUuid"` /** 唯一ID **/
	AppId     string `json:"appId" gorm:"column:app_id" form:"appId"`             /** 平台号（类似房间号） **/
	UserId    string `json:"userId" gorm:"column:user_id" form:"userId"`          /** 成员信息 **/
	Status    int32  `json:"status" gorm:"column:status" form:"status"`           /** 用户状态 1 在用 2 锁定 **/
	AvatarUrl string `json:"avatarUrl" gorm:"column:avatar_url" form:"avatarUrl"` /** 用户头像链接 **/
	Nickname  string `json:"nickname" gorm:"column:nickname" form:"nickname"`
}

func (d *GroupElement) TableName() string {
	return "group_element"
}

func (d *GroupElement) Save(entity *GroupElement) (err error) {
	err = database.GetDB().Model(d).Save(entity).Error
	return err
}

func (d *GroupElement) GetByUuidAndUserId(uuid, userId string) (*GroupElement, error) {
	res := &GroupElement{}
	err := database.GetDB().Model(d).Where("group_uuid = ? and user_id = ?", uuid, userId).First(&res).Error
	if err != nil {
		log.Println("record not found", err)
		return nil, err
	}
	return res, nil
}

func (d *GroupElement) Update(entity GroupElement) error {
	err := database.GetDB().Model(d).Omit("id").Where("id = ?", entity.Id).Updates(&entity).Error
	log.Println("更新信息", err)
	return err
}

func (d *GroupElement) List(entity GroupElement) (res []GroupElement, err error) {
	err = database.GetDB().Model(d).Where(entity).Find(&res).Error
	return res, err
}

func (d *GroupElement) Delete(id int64) (err error) {
	err = database.GetDB().Model(d).Where("id = ?", id).Delete(&GroupElement{}).Error
	return err
}

func (d *GroupElement) ListVo(entity GroupElement) ([]GroupElementVO, error) {
	res := make([]GroupElementVO, 0)
	//avatar_url, nickname 是user表字段 与element表通过 user_id 关联 group_uuid 为查询条件
	err := database.GetDB().Table(d.TableName()).Select("group_element.group_uuid, group_element.app_id,group_element.user_id, group_element.status," +
		" user.avatar_url, user.nickname").
		Joins("LEFT JOIN user ON group_element.user_id = user.user_id").
		Where(entity).Find(&res).Error
	return res, err
}
