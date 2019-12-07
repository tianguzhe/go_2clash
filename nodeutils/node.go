package nodeutils

import (
	b64 "encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

//获取订阅链接文本
func GetUrlData(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	s, err := ioutil.ReadAll(resp.Body)
	return string(s)
}

//base64解码
func Base64Encode(url string) string {
	var allNode string

	urlSplit := strings.Split(url, "+")
	for _, url := range urlSplit {
		data := GetUrlData(url)
		decodeString, err := b64.URLEncoding.DecodeString(DecodeInfoByByte(data))
		if err != nil {
			fmt.Printf("[base64Encode Error] %s \n", err)
			continue
		}
		allNode += string(decodeString)
	}
	return allNode
}

func DecodeInfoByByte(info string) string {
	lens := len(info)
	if lens%4 == 1 {
		info = info + "==="
	} else if lens%4 == 2 {
		info = info + "=="
	} else if lens%4 == 3 {
		info = info + "="
	}
	return info
}
