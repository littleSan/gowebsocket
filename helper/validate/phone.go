/*
*

	@author:
	@date : 2023/9/27
*/
package validate

import "regexp"

func ValidatePhoneNumber(phone string) bool {
	pattern := `^1[3-9]\d{9}$`
	match, _ := regexp.MatchString(pattern, phone)
	return match
}
