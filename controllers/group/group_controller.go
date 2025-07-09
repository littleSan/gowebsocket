/*
*

	@author:
	@date : 2025/5/27
*/
package group

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/google/uuid"
	"github.com/link1st/gowebsocket/v2/common"
	"github.com/link1st/gowebsocket/v2/controllers"
	"github.com/link1st/gowebsocket/v2/models/group"
	"github.com/link1st/gowebsocket/v2/models/groupElement"
	user2 "github.com/link1st/gowebsocket/v2/models/user"
	"strings"
	"time"
)

// 创建分组
func CreateGroup(c *gin.Context) {
	param := new(group.Group)
	err := c.BindJSON(param)
	if err != nil {
		controllers.Response(c, common.ParameterIllegal, err.Error(), nil)
		return
	}
	//产生UUID
	param.Uuid = uuid.NewString()
	param.Status = 1
	param.CreateTime = time.Now().UnixMilli()
	//查询创建者账户是否存在
	user := new(user2.UserPO)
	param.UserId = strings.ToLower(param.UserId)
	_, err = user.UserByUidAndAppId(param.UserId, param.AppId)
	if err != nil {
		controllers.Response(c, common.UnauthorizedUserID, common.GetErrorMessage(common.UnauthorizedUserID, ""), nil)
		return
	}
	err = param.Save(param)
	if err != nil {
		controllers.Response(c, common.ModelAddError, err.Error(), nil)
		return
	}
	//同步将个人添加至群成员
	ele := new(groupElement.GroupElement)
	ele.UserId = param.UserId
	ele.GroupUuid = param.Uuid
	ele.CreateTime = time.Now().UnixMilli()
	ele.Status = 1
	ele.AppId = param.AppId
	err = ele.Save(ele)
	data := make(map[string]interface{})
	data["group"] = param
	controllers.Response(c, common.OK, "", data)
}

// 添加分组成员
func AddGroupMember(c *gin.Context) {
	param := new(groupElement.GroupElement)
	err := c.BindJSON(param)
	if err != nil {
		controllers.Response(c, common.ParameterIllegal, err.Error(), nil)
		return
	}
	//查看用户是否存在
	userCv := new(user2.UserPO)
	param.UserId = strings.ToLower(param.UserId)
	_, err = userCv.UserByUidAndAppId(param.UserId, param.AppId)
	if err != nil {
		controllers.Response(c, common.UnauthorizedUserID, common.GetErrorMessage(common.UnauthorizedUserID, ""), nil)
		return
	}
	//查看分组属否存在
	group := new(group.Group)
	_, err = group.GetByUuid(param.GroupUuid)
	if err != nil {
		controllers.Response(c, common.GroupNotExist, common.GetErrorMessage(common.GroupNotExist, ""), nil)
		return
	}
	//查询群组成员是否存在
	el, err := param.GetByUuidAndUserId(param.GroupUuid, param.UserId)
	if err == nil && el != nil {
		controllers.Response(c, common.GroupMemberExist, common.GetErrorMessage(common.GroupMemberExist, ""), nil)
		return
	}
	//添加成员
	err = param.Save(param)
	if err != nil {
		controllers.Response(c, common.ModelAddError, err.Error(), nil)
		return
	}
	data := make(map[string]interface{})
	data["groupElement"] = param
	controllers.Response(c, common.OK, "", data)
}

// 删除分组成员
func DeleteGroupMember(c *gin.Context) {
	param := new(groupElement.GroupElement)
	err := c.BindJSON(param)
	if err != nil {
		controllers.Response(c, common.ParameterIllegal, err.Error(), nil)
		return
	}
	//查询是否存在
	param.UserId = strings.ToLower(param.UserId)
	el, err := param.GetByUuidAndUserId(param.GroupUuid, param.UserId)
	if err != nil || el == nil {
		controllers.Response(c, common.GroupMemberNotExist, common.GetErrorMessage(common.GroupMemberNotExist, ""), nil)
		return
	}

	err = param.Delete(el.Id)
	if err != nil {
		controllers.Response(c, common.ModelDeleteError, common.GetErrorMessage(common.ModelDeleteError, ""), nil)
		return
	}
	controllers.Response(c, common.OK, "", nil)
}

// 获取用户所有的所属分组信息
func GroupList(c *gin.Context) {
	userId := c.Query("userId")
	appId := c.Query("appId")
	if userId == "" || appId == "" {
		controllers.Response(c, common.ParameterIllegal, common.GetErrorMessage(common.ParameterIllegal, ""), nil)
		return
	}
	userId = strings.ToLower(userId)
	data := make(map[string]interface{})
	groupCv := new(group.Group)
	data["groupList"], _ = groupCv.UserGroupList(appId, userId)
	controllers.Response(c, common.OK, "", data)
}

// 获取群组成员信息
func GroupMemberList(c *gin.Context) {
	uuid := c.Query("uuid")
	//校验群组 是否存在
	group := new(group.Group)
	_, err := group.GetByUuid(uuid)
	if err != nil {
		controllers.Response(c, common.GroupNotExist, common.GetErrorMessage(common.GroupNotExist, ""), nil)
		return
	}
	elementCv := new(groupElement.GroupElement)
	data := make(map[string]interface{})
	data["groupMemberList"], _ = elementCv.ListVo(groupElement.GroupElement{GroupUuid: uuid})

	controllers.Response(c, common.OK, "", data)
}
