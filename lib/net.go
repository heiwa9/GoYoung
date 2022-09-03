package lib

import (
	"net"
	"time"
)

func CheckServer(url string) string {
	timeout := 5 * time.Second
	// t1 := time.Now()
	_, err := net.DialTimeout("tcp", url, timeout)
	// massage += "\n网络测试时长 :" + time.Since(t1).String()

	if err == nil {
		return "已接入互联网，只能进行下线操作"
	} else {
		return "未接入互联网"
	}
}
