/*
*

	@author:
	@date : 2023/10/9
*/
package validate

import "testing"

func TestValidateCreditCode(t *testing.T) {
	ok := ValidateCreditCode("91320102MA1WCY0L8N")
	if !ok {
		t.Fatal(ok)
	}
	t.Logf("执行成功")
}

func TestValidateIDCard(t *testing.T) {
	ok := ValidateIDCard("41272219930125613X")
	if ok {
		t.Logf("success")
	}
	t.Fatal(ok)
}
