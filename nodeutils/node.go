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
	//resp, err := http.Get(url)
	//if err != nil {
	//	panic(err)
	//}
	//defer resp.Body.Close()
	//
	//s, err := ioutil.ReadAll(resp.Body)
	//return string(s)

	client := &http.Client{}
	reqest, _ := http.NewRequest("GET", url, nil)
	reqest.Header.Add("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 13_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.4 Mobile/15E148 Safari/604.1")

	resp, err := client.Do(reqest) //提交
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

//获取clash节点
func GetClashProxy(url string) string {

	var allNode string

	urlSplit := strings.Split(url, "+")
	for _, url := range urlSplit {
		data := GetUrlData(url)
		s := strings.Split(strings.Split(data, "Proxy:")[1], "Proxy Group:")[0]
		allNode += s
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
