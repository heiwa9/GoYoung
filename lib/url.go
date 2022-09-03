/**
 * @Author Kokutas
 * @Description //TODO
 * @Date 2021/2/11 0:13
 **/
package lib

import (
	"log"
	"net/url"
)

func ParseURL(urlStr string)*url.URL{
	link,err:=url.Parse(urlStr)
	if err!=nil{
		log.Println("Could not parse URL:",err)
	}
	return link
}
