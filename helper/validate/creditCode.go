/*
*

	@author:
	@date : 2023/10/9
*/
package validate

import (
	"fmt"
	"github.com/bluesky335/IDCheck/IDNumber"
	"github.com/bluesky335/IDCheck/USCI"
)

// ValidateCreditCode 校验社会统一信用代码
func ValidateCreditCode(code string) bool {

	var usci = USCI.New(code)
	return usci.IsValid()
}

// ValidateIDCard 校验身份证号码
func ValidateIDCard(code string) bool {
	var id = IDNumber.New(code)

	//var birthday = id.GetBirthday()
	//if birthday != nil {
	//	fmt.Printf("生日：%s-%s-%s\n", birthday.Year, birthday.Month, birthday.Day)
	//} else {
	//	// 不合法的身份证
	//}

	// 产生一个随机的符合校验规则的身份证号码。注意！它虽然符合校验规则，但不一定真实存在。
	randomIDCard := IDNumber.Random()
	fmt.Printf("随机身份证：%s\n", randomIDCard)
	return id.IsValid()
}
