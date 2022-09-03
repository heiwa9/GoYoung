package main

import (
	"GoYoung/lib"
	"crypto/md5"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/qifengzhang007/goCurl"
)

const (
	baiduURL         = "baidu.com:443"
	feiYoungRedirect = "http://www.msftconnecttest.com/redirect"
)

type User struct {
	UserHard     string `json:"user_hard"`
	UserAccount  string `json:"user_account"`
	PassWord     string `json:"pass_word"`
	LastLoginURL string `json:"last_login_url"`
}

func (user *User) login() string {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: 1 * time.Second,
	}
	res, err := client.Get(feiYoungRedirect)
	if err != nil {
		return "内部错误,代码-101"
	}
	if res.StatusCode != http.StatusFound {
		return "内部错误,代码-102"
	}
	resp, err := httpClient.Get(res.Header.Get("Location"), goCurl.Options{
		Headers: map[string]interface{}{
			"User-Agent": "CDMA+WLAN(Mios)",
		},
		SetResCharset: "utf-8",
		Timeout:       1,
	})
	if err != nil {
		return "内部错误,代码-103"
	}
	body, err := resp.GetContents()
	if err != nil {
		return "内部错误,代码-104"
	}
	user.LastLoginURL = lib.ParseXML(body, "WISPAccessGatewayParam", "Redirect", "LoginURL")
	encryptPass := func() string {
		if len(user.PassWord) != 16 {
			return user.EncryptPassword()
		}
		return user.PassWord
	}()
	token := "UserName=" + user.UserHard + user.UserAccount + "&Password=" + encryptPass + "&AidcAuthAttr1=" + time.Now().Format("20060102150405") +
		"&AidcAuthAttr3=keuyGQlK&AidcAuthAttr4=zrDgXllCChyJHjwkcRwhygP0&AidcAuthAttr5=kfe1GQhXdGqOFDteego5zwP9IsNoxX7djTWspPrYm1A%3D%3D&" +
		"AidcAuthAttr6=5Ia4cQhDfXSFbTtUDGY1yx8%3D&AidcAuthAttr7=6ZWiVlwdNiHMXCpOagQv2w2MQs0ohTWJnTu8qK5OibhCydTpTxkI88wadKPWby%2F2PKCVaZ" +
		"UxglbBs96%2FtmLE89M8AJ6y28o7qolpFep%2FcYFFRLd7H4MAMrDUMRO0F%2B93jh14fiAZYmtk9hdp%2BZ5w%2BjMQUoV4TCtM9VJ07XQwxlMVg%2F0YKrS1s3hXA" +
		"stdQ1fvdSn3nAVGgdxc%2BJQDrQ%3D%3D&AidcAuthAttr8=jPSyBQxVaXWTQWUaakluj06scJ98nyqCyX7y%2FLUk1OkXiNjkXhVGvJhyTuLDaCPhK%2FOFJttlxxi" +
		"VqNKupnDXkp9%2BR9D9j8p2j5h8FOxoatMaGu0oRdk%3D&createAuthorFlag=0"
	resp, err = httpClient.Post(user.LastLoginURL, goCurl.Options{
		Headers: map[string]interface{}{
			"User-Agent":   "CDMA+WLAN(Mios)",
			"Content-Type": "application/x-www-form-urlencoded",
		},
		XML:           token,
		SetResCharset: "utf-8",
		Timeout:       1,
	})
	if err != nil {
		return "内部错误,代码-105"
	}
	body, err = resp.GetContents()
	if err != nil {
		return "内部错误,代码-106"
	}
	body = lib.ParseXML(body, "WISPAccessGatewayParam", "AuthenticationReply", "ReplyMessage")
	return body
}

func (user *User) Logout() string {
	u := "http://58.53.199.144:8001/wispr_logout.jsp?" + strings.Split(user.LastLoginURL, "?")[1]
	var httpClient = goCurl.CreateHttpClient()
	resp, err := httpClient.Get(u, goCurl.Options{
		Headers: map[string]interface{}{
			"User-Agent": "CDMA+WLAN(Mios)",
		},
		SetResCharset: "utf-8",
		Timeout:       2,
	})
	if err != nil {
		return "内部错误,代码-201"
	}
	body, err := resp.GetContents()
	if err != nil {
		return "内部错误,代码-202"
	}
	body = lib.ParseXML(body, "WISPAccessGatewayParam", "LogoffReply", "ResponseCode")
	if body == "150" {
		return "150:下线成功"
	} else {
		return "255:下线失败"
	}
}

func (user *User) EncryptPassword() string {
	wordData := map[int]string{
		1:  "5084972163",
		2:  "9801567243",
		3:  "7286059143",
		4:  "1850394726",
		5:  "1462578093",
		6:  "5042936178",
		7:  "0145937682",
		8:  "0964238571",
		9:  "3497651802",
		10: "9125780643",
		11: "8634972150",
		12: "5924673801",
		13: "8274053169",
		14: "5841792063",
		15: "2469385701",
		16: "8205349671",
		17: "7429516038",
		18: "3769458021",
		19: "5862370914",
		20: "8529364170",
		21: "7936082154",
		22: "5786241930",
		23: "0728643951",
		24: "9418360257",
		25: "5093287146",
		26: "5647830192",
		27: "3986145207",
		28: "0942587136",
		29: "4357069128",
		30: "0956723814",
		31: "1502796384",
	}
	word := wordData[time.Now().Day()]
	wordLen := len(word)
	wordByte := []byte(word)

	var token []int
	for i := 0; i < 256; i++ {
		if i < 128 {
			token = append(token, i)
		} else {
			token = append(token, i-256)
		}
	}

	for index, temp, i := 0, 0, 0; i < 256; i++ {
		temp, _ = strconv.Atoi(string(wordByte[i%wordLen]))
		index += token[i] + temp&255
		if index < 0 {
			index = (index%256 + 256) % 256
		} else {
			index %= 256
		}
		token[index], token[i] = token[i], token[index]
	}

	passwdByte := []byte(user.PassWord)
	passwdLen := len(user.PassWord)
	var passwdToken []int
	for i := range user.PassWord {
		passwdToken = append(passwdToken, i)
	}

	for i, index, index1, index2 := 0, 0, 0, 0; i < passwdLen; i++ {
		index1 += 1 & 255
		index1 %= 256
		index2 += token[index1] & 255
		index2 %= 256
		token[index1], token[index2] = token[index2], token[index1]
		index = token[index1] + token[index2]&255
		if index < 0 {
			index = (index%256 + 256) % 256
		} else {
			index %= 256
		}
		passwdToken[i] = 256 + token[index] ^ int(passwdByte[i])
		passwdToken[i] %= 256
	}
	var temp []byte
	for _, i2 := range passwdToken {
		temp = append(temp, byte(i2))
	}
	return fmt.Sprintf("%x", md5.Sum(temp))[8:24]
}
