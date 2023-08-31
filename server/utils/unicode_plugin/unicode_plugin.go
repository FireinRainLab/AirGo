package unicode_plugin

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// 字符串转为unicode码
func StrToUnicode(str string) string {
	DD := []rune(str) //需要分割的字符串内容，将它转为字符，然后取长度。
	finallStr := ""
	for i := 0; i < len(DD); i++ {
		if unicode.Is(unicode.Scripts["Han"], DD[i]) {
			textQuoted := strconv.QuoteToASCII(string(DD[i]))
			finallStr += textQuoted[1 : len(textQuoted)-1]
		} else {
			h := fmt.Sprintf("%x", DD[i])
			finallStr += "\\u" + isFullFour(h)
		}
	}
	return finallStr
}

func isFullFour(str string) string {
	if len(str) == 1 {
		str = "000" + str
	} else if len(str) == 2 {
		str = "00" + str
	} else if len(str) == 3 {
		str = "0" + str
	}
	return str
}

// unicode 转字符
func UnicodeToStr(raw []byte) ([]byte, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1))
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}
