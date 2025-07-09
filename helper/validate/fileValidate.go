/*
*

	@author:
	@date : 2023/11/15
*/
package validate

import (
	"mime/multipart"
	"strings"
)

// / 文件校验
var IMGAllow = "jpg,jpeg,png,gif,bmp,webp,pdf,tiff,psd,dds,pdt,webp,xmp,svg,pdd,raw,cr2"
var VideoAllow = "avi,wmv,mpeg,mp4,flv,mov,mkv,webm,f4v,avchd,dat,mod"
var FileAllow = "zip,doc,docx,json,txt,xlsx,pdf,xls,rar,pptx,csv"

var TotalAllow = IMGAllow + "," + FileAllow + "," + VideoAllow
var ImageMaxSize = 30 << 20
var FileMaxSize = 30 << 20
var VideoMaxSize = 100 << 20

func FileValidate(fi multipart.FileHeader, typeAllow string, maxSize int64) bool {

	names := fi.Filename
	if !prefix(names, typeAllow) {
		return false
	}
	size := fi.Size
	return size != 0 && size <= maxSize
}
func FilesValidate(fi []*multipart.FileHeader, typeAllow string, maxSize int64) bool {
	var sizeTotal int64
	for _, f := range fi {
		names := f.Filename
		if !prefix(names, typeAllow) {
			return false
		}
		sizeTotal += f.Size
	}
	return sizeTotal != 0 && sizeTotal <= maxSize

}
func prefix(filename, allow string) bool {
	filename = strings.ToLower(filename)
	index := strings.Index(filename, "?")
	if index > 0 {
		filename = filename[:index]
	}
	allowType := strings.Split(allow, ",")
	types := strings.Split(filename, ".")
	t := types[len(types)-1]
	return itemExistSlice(t, allowType)
}

func itemExistSlice(a string, source []string) bool {
	for _, s := range source {
		if strings.EqualFold(s, a) {
			return true
		}
	}
	return false
}
