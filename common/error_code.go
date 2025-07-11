// Package common 通用函数
package common

const (
	OK                  = 200  // Success
	NotLoggedIn         = 1000 // 未登录
	ParameterIllegal    = 1001 // 参数不合法
	UnauthorizedUserID  = 1002 // 非法的用户 ID
	Unauthorized        = 1003 // 未授权
	ServerError         = 1004 // 系统错误
	NotData             = 1005 // 没有数据
	ModelAddError       = 1006 // 添加错误
	ModelDeleteError    = 1007 // 删除错误
	ModelStoreError     = 1008 // 存储错误
	OperationFailure    = 1009 // 操作失败
	RoutingNotExist     = 1010 // 路由不存在
	GroupNotExist       = 1011 //  群组不存在
	GroupMemberExist    = 1012 // 群组成员已存在
	GroupMemberNotExist = 1013 // 群组成员不存在
	FileFormatError     = 1014
)

// GetErrorMessage 根据错误码 获取错误信息
func GetErrorMessage(code uint32, message string) string {
	var codeMessage string
	codeMap := map[uint32]string{
		OK:                  "Success",
		NotLoggedIn:         "未登录",
		ParameterIllegal:    "参数不合法",
		UnauthorizedUserID:  "非法的用户ID",
		Unauthorized:        "未授权",
		NotData:             "没有数据",
		ServerError:         "系统错误",
		ModelAddError:       "添加错误",
		ModelDeleteError:    "删除错误",
		ModelStoreError:     "存储错误",
		OperationFailure:    "操作失败",
		RoutingNotExist:     "路由不存在",
		GroupNotExist:       "群组不存在",
		GroupMemberExist:    "群组成员已存在",
		GroupMemberNotExist: "群组成员不存在",
	}

	if message == "" {
		if value, ok := codeMap[code]; ok {
			// 存在
			codeMessage = value
		} else {
			codeMessage = "未定义错误类型!"
		}
	} else {
		codeMessage = message
	}

	return codeMessage
}
