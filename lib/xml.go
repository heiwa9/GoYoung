package lib

import (
	"strings"

	"github.com/tinyhubs/tinydom"
)

func ParseXML(str, str1, str2, str3 string) string {
	doc, err := tinydom.LoadDocument(strings.NewReader(str))
	if err != nil {
		return "内部错误,代码-701"
	}
	elem := doc.FirstChildElement(str1).FirstChildElement(str2).FirstChildElement(str3)
	return elem.Text()
}
