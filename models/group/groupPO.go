package group

import (
	"github.com/link1st/gowebsocket/v2/lib/database"
	"log"
)

// Group 分组信息
type Group struct {
	Id         int64  `json:"id" gorm:"primary_key" gorm:"column:id" form:"id"`       /** 自增主键 **/
	Name       string `json:"name" gorm:"column:name" form:"name"`                    /** 分组名称 **/
	Uuid       string `json:"uuid" gorm:"column:uuid" form:"uuid"`                    /** 唯一ID **/
	AppId      string `json:"appId" gorm:"column:app_id" form:"appId"`                /** 平台号（类似房间号） **/
	UserId     string `json:"userId" gorm:"column:user_id" form:"userId"`             /** 创建人 **/
	Status     int32  `json:"status" gorm:"column:status" form:"status"`              /** 用户状态 1 在用 2 锁定 **/
	CreateTime int64  `json:"createTime" gorm:"column:create_time" form:"createTime"` /** 更新时间 **/
}

func (d *Group) TableName() string {
	return "group"
}

func (d *Group) Save(entity *Group) (err error) {
	err = database.GetDB().Model(d).Save(entity).Error
	return err
}

func (d *Group) GetByUuid(uuid string) (*Group, error) {
	res := new(Group)
	err := database.GetDB().Model(d).Where("uuid = ?", uuid).First(&res).Error
	if err != nil {
		log.Println("record not found", err)
		return nil, err
	}
	return res, nil
}

func (d *Group) Update(entity Group) error {
	err := database.GetDB().Model(d).Omit("id").Where("id = ?", entity.Id).Updates(&entity).Error
	log.Println("更新信息", err)
	return err
}

func (d *Group) List(entity Group) (res []Group, err error) {
	err = database.GetDB().Model(d).Where(entity).Find(&res).Error
	return res, err
}

func (d *Group) Delete(id int64) (err error) {
	err = database.GetDB().Model(d).Where("id = ?", id).Delete(&Group{}).Error
	return err
}

func (d *Group) UserGroupList(appId, userId string) ([]Group, error) {
	var res []Group
	// 使用 GORM 进行多表关联查询，通过 group.uuid 和 group_element.group_uuid 关联
	err := database.GetDB().Model(d).
		Joins("JOIN group_element ON group.uuid = group_element.group_uuid").
		Where("group.app_id = ? AND group_element.user_id = ?", appId, userId).
		Find(&res).Error

	if err != nil {
		return nil, err
	}
	// 使用 map 去重，以 Uuid 作为唯一标识
	uniqueGroups := make(map[string]Group)
	for _, group := range res {
		uniqueGroups[group.Uuid] = group
	}

	// 转换回 slice
	result := make([]Group, 0, len(uniqueGroups))
	for _, group := range uniqueGroups {
		result = append(result, group)
	}

	return result, nil
}
