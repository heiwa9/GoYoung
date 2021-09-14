package main

import (
	"encoding/json"
	"fmt"
	"github.com/tinyhubs/tinydom"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/qifengzhang007/goCurl"
)

const (
	baiduURL         = "baidu.com:443"
	feiYoungRedirect = "http://www.msftconnecttest.com/redirect"
	filePath         = "./data.json"
)

var httpClient = goCurl.CreateHttpClient()

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
		return "INTERNAL_ERROR_CODE:-101"
	}
	if res.StatusCode != http.StatusFound {
		return "INTERNAL_ERROR_CODE:-102"
	}
	u, _ := url.Parse(res.Header.Get("Location"))
	query := u.Query()
	nextURL := fmt.Sprintf("http://58.53.199.144:8001/?userip=%s&wlanacname=&nasip=%s&usermac=%s&aidcauthtype=0",
		query.Get("userip"), query.Get("nasip"), query.Get("usermac"))

	resp, err := httpClient.Get(nextURL, goCurl.Options{
		Headers: map[string]interface{}{
			"User-Agent": "CDMA+WLAN(Mios)",
		},
		SetResCharset: "utf-8",
		Timeout:       1,
	})
	if err != nil {
		return "INTERNAL_ERROR_CODE:-103"
	}
	body, err := resp.GetContents()
	if err != nil {
		return "INTERNAL_ERROR_CODE:-104"
	}
	user.LastLoginURL = user.ParseXML(body, "WISPAccessGatewayParam", "Redirect", "LoginURL")
	token := "UserName=" + user.UserHard + user.UserAccount + "&Password=" + user.PassWord + "&AidcAuthAttr1=" + time.Now().Format("20060102150405") +
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
		return "INTERNAL_ERROR_CODE:-105"
	}
	body, err = resp.GetContents()
	if err != nil {
		return "INTERNAL_ERROR_CODE:-106"
	}
	body = user.ParseXML(body, "WISPAccessGatewayParam", "AuthenticationReply", "ReplyMessage")
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
		return "INTERNAL_ERROR_CODE:-201"
	}
	body, err := resp.GetContents()
	if err != nil {
		return "INTERNAL_ERROR_CODE:-202"
	}
	body = user.ParseXML(body, "WISPAccessGatewayParam", "LogoffReply", "ResponseCode")
	if body == "150" {
		return "150:Successfully offline"
	} else {
		return "255:Offline failure"
	}
}

func (user *User) SaveUserInfoJson() string {
	data, err := json.Marshal(user)
	if err != nil {
		return "INTERNAL_ERROR_CODE:-401"
	}
	err = ioutil.WriteFile(filePath, data, os.ModeAppend)
	if err != nil {
		return "INTERNAL_ERROR_CODE:-402"
	}
	return ""
}

func (user *User) ReadUserInfoJson() string {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "INTERNAL_ERROR_CODE:-501"
	}
	err = json.Unmarshal(bytes, &user)
	if err != nil {
		return "INTERNAL_ERROR_CODE:-502"
	}
	return ""
}

func (user *User) ParseXML(str, str1, str2, str3 string) string {
	doc, err := tinydom.LoadDocument(strings.NewReader(str))
	if err != nil {
		return "INTERNAL_ERROR_CODE:-701"
	}
	elem := doc.FirstChildElement(str1).FirstChildElement(str2).FirstChildElement(str3)
	return elem.Text()
}

func (user *User) CheckServer(url string) string {
	_, err := net.DialTimeout("tcp", url, 5*time.Second)
	if err == nil {
		return "Has been connected to the Internet and can only be offline"
	} else {
		return "Not connected to the internet"
	}
}
