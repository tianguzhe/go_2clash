package main

import (
	. "awesomeProject/model"
	. "awesomeProject/nodedao"
	. "awesomeProject/nodeutils"
	"bufio"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var emojiData = make(chan string)
var ppp string

//提取节点信息
func getAllNodes(nodes string) []NodeBean {

	addNodes := make([]NodeBean, 0)

	split := strings.Split(nodes, "\n")
	for _, node := range split {
		if node == "" {
			continue
		}

		i := strings.Split(node, "//")

		if i[0] == "vmess:" {
			ppp := getVmessNode(i[1])
			addNodes = append(addNodes, ppp)
		} else if i[0] == "ss:" {
			ppp := getSSNode(i[1])
			addNodes = append(addNodes, ppp)
		}

	}
	return addNodes
}

func getClashProxyName(nodes string) string {
	nodeItem := strings.Split(nodes, "\n")

	var allNodeName string

	for _, value := range nodeItem {

		item := strings.Split(strings.Split(value, ",")[0], ":")

		if len(item) == 2 {
			allNodeName += item[1] + ","
		}

	}

	return allNodeName[:len(allNodeName)-1]
}

func getSSNode(s string) NodeBean {

	info, err := b64.StdEncoding.DecodeString(DecodeInfoByByte(strings.Split(s, "@")[0]))
	if err != nil {
		panic(err)
	}

	var nodeBean NodeBean

	nodeBean.Name, _ = url.QueryUnescape(strings.Split(s, "#")[1])
	nodeBean.Name = strings.ReplaceAll(nodeBean.Name, "#", "")
	nodeBean.NodeType = "ss"
	nodeBean.Server = strings.Split(strings.Split(s, "@")[1], ":")[0]
	nodeBean.Port = strings.Split(strings.Split(s, ":")[1], "#")[0]
	nodeBean.Cipher = strings.Split(string(info), ":")[0]
	nodeBean.Password = strings.Split(string(info), ":")[1]

	return nodeBean

}

func getVmessNode(s string) NodeBean {

	decodeString, err := b64.StdEncoding.DecodeString(DecodeInfoByByte(s))
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s \n", decodeString)

	var nodeBean NodeBean

	err = json.Unmarshal(decodeString, &nodeBean)
	if err != nil {
		fmt.Println(err)
	}

	nodeBean.NodeType = "vmess"
	fmt.Printf("%#v \n", nodeBean)

	return nodeBean
}

func setNodes(infos []NodeBean) string {

	proxies := make([]string, 0)

	//ppp := getUrlData("https://raw.githubusercontent.com/bddddd/ConfConvertor/master/Emoji/flag_emoji.json")
	//ppp := <-emojiData

	m := make(map[string][]string)

	err := json.Unmarshal([]byte(ppp), &m)
	if err != nil {
		panic(err)
	}

	for _, node := range infos {
		name := strings.ReplaceAll(node.Name, "\r", "")

		for k, v := range m {
			for _, subV := range v {
				if strings.Contains(name, subV) {
					name = k + name
					break
				}
			}
		}

		nodeType := node.NodeType
		server := node.Server
		port := node.Port
		uuid := node.Uuid
		alterId := node.AlterId
		network := node.Network
		host := node.Host
		cipher := node.Cipher
		password := node.Password
		path := node.Path
		tls := node.Tls

		var proxy string

		if nodeType == "vmess" {
			if tls != "" {
				if host == "" {
					proxy = fmt.Sprintf("- { name: \"%s\", type: %s, server: %s, port: %v, uuid: %s, alterId: %v, cipher: auto, network: %s, ws-path: \"%s\" , tls: true }", name, nodeType, server, port, uuid, alterId, network, path)
				} else {
					proxy = fmt.Sprintf("- { name: \"%s\", type: %s, server: %s, port: %v, uuid: %s, alterId: %v, cipher: auto, network: %s, ws-path: \"%s\", ws-headers: {Host: %s}, tls: true }", name, nodeType, server, port, uuid, alterId, network, path, host)
				}
			} else {
				if host == "" {
					proxy = fmt.Sprintf("- { name: \"%s\", type: %s, server: %s, port: %v, uuid: %s, alterId: %v, cipher: auto, network: %s, ws-path: \"%s\" }", name, nodeType, server, port, uuid, alterId, network, path)
				} else {
					proxy = fmt.Sprintf("- { name: \"%s\", type: %s, server: %s, port: %v, uuid: %s, alterId: %v, cipher: auto, network: %s, ws-path: \"%s\", ws-headers: {Host: %s} }", name, nodeType, server, port, uuid, alterId, network, path, host)

				}
			}

		} else if nodeType == "ss" {
			proxy = fmt.Sprintf("- { name: %s, type: %s, server: %s, port: %v, cipher: %s, password: %s }", name, nodeType, server, port, cipher, password)
		}

		proxies = append(proxies, proxy)

	}

	return "\nProxy:\n" + strings.Join(proxies, "\n")
}

func setPG(infos []NodeBean) string {

	var allName string

	//ppp := getUrlData("https://raw.githubusercontent.com/bddddd/ConfConvertor/master/Emoji/flag_emoji.json")

	//ppp := <-emojiData

	m := make(map[string][]string)

	err := json.Unmarshal([]byte(ppp), &m)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(infos); i++ {

		name := infos[i].Name

		for k, v := range m {
			for _, subV := range v {
				if strings.Contains(name, subV) {
					name = k + name
					break
				}
			}
		}

		if i == 0 {
			allName = fmt.Sprintf("\"%s\"", strings.ReplaceAll(name, "\r", ""))
		} else {
			allName = fmt.Sprintf("%s, \"%s\"", allName, strings.ReplaceAll(name, "\r", ""))
		}
	}

	Proxy0 := fmt.Sprintf("- { name: \"PROXY\", type: select, proxies: [%s] }\n", allName)
	Proxy1 := fmt.Sprintf("- { name: \"Domestic\", type: select, proxies:  [\"DIRECT\",%s] }\n", allName)
	Proxy2 := fmt.Sprintf("- { name: \"GlobalTV\", type: select, proxies:  [\"PROXY\",%s] }\n", allName)
	Proxy3 := fmt.Sprintf("- { name: \"Prevent\", type: select, proxies:  [\"REJECT\",\"DIRECT\"] }\n")
	Proxy4 := fmt.Sprintf("- { name: \"Others\", type: select, proxies:  [\"PROXY\",\"DIRECT\"] }\n")

	sha1Str := "\n" + "#" + urlSha1 + "\n"

	Rule := "\n"

	return "\n\nProxy Group:\n" + Proxy0 + Proxy1 + Proxy2 + Proxy3 + Proxy4 + sha1Str + Rule
}

func setClashPg(s string) string {
	var allName string

	split := strings.Split(s, ",")

	for i := 0; i < len(split); i++ {

		name := strings.TrimSpace(strings.ReplaceAll(split[i], "\"", ""))

		if i == 0 {
			allName = fmt.Sprintf("\"%s\"", strings.ReplaceAll(name, "\r", ""))
		} else {
			allName = fmt.Sprintf("%s, \"%s\"", allName, strings.ReplaceAll(name, "\r", ""))
		}
	}

	Proxy0 := fmt.Sprintf("- { name: \"PROXY\", type: select, proxies: [%s] }\n", allName)
	Proxy1 := fmt.Sprintf("- { name: \"Domestic\", type: select, proxies:  [\"DIRECT\",%s] }\n", allName)
	Proxy2 := fmt.Sprintf("- { name: \"GlobalTV\", type: select, proxies:  [\"PROXY\",%s] }\n", allName)
	Proxy3 := fmt.Sprintf("- { name: \"Prevent\", type: select, proxies:  [\"REJECT\",\"DIRECT\"] }\n")
	Proxy4 := fmt.Sprintf("- { name: \"Others\", type: select, proxies:  [\"PROXY\",\"DIRECT\"] }\n")

	sha1Str := "\n" + "#" + urlSha1 + "\n"

	Rule := "\n"

	return "\n\nProxy Group:\n" + Proxy0 + Proxy1 + Proxy2 + Proxy3 + Proxy4 + sha1Str + Rule
}

func getEmojiData(ch chan string) {
	ppp := GetUrlData("https://raw.githubusercontent.com/bddddd/ConfConvertor/master/Emoji/flag_emoji.json")
	ch <- ppp
}

var urlSha1 string

func main() {

	header := GetUrlData("https://raw.githubusercontent.com/tianguzhe/ClashRConf/master/General.yml")

	foot := GetUrlData("https://raw.githubusercontent.com/tianguzhe/ClashRConf/master/Rule.yml")

	r := gin.Default()
	r.GET("/2clash", func(c *gin.Context) {
		urlAddress := c.DefaultQuery("url", "")
		urlSha1 = c.DefaultQuery("code", "")

		file, err := os.OpenFile("./newUrlList.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
		if err != nil {
			panic(err)
		}

		writer := bufio.NewWriter(file)
		writer.WriteString(urlAddress + "\n")
		writer.Flush()

		if urlAddress != "" {
			checkUrl, err := CheckUrl(urlAddress)
			if err != nil {
				err = CreateMD5Url(urlAddress)
				if err != nil {
					fmt.Printf("[Error] Insert Data Error %s", err)
					return
				}
			} else {
				_, err := UpdateUrl(urlAddress)
				if err != nil {
					fmt.Printf("[Error] Insert Data Error %s", err)
					return
				}

				urlSha1 = checkUrl.Sha1Str
			}
		}

		if urlSha1 != "" {
			checkUrl, err := CheckUrlSha1(urlSha1)
			if err != nil {
				return
			}

			urlAddress = checkUrl.Urls

		}

		go getEmojiData(emojiData)

		nodeInfo := getAllNodes(Base64Encode(urlAddress))

		ppp = <-emojiData

		proxy := setNodes(nodeInfo)
		proxyGroup := setPG(nodeInfo)

		c.String(http.StatusOK, "%s", header+proxy+proxyGroup+foot)
	})

	r.GET("/clash2clash", func(c *gin.Context) {
		urlAddress := c.DefaultQuery("url", "")
		urlSha1 = c.DefaultQuery("code", "")

		file, err := os.OpenFile("./newUrlList.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
		if err != nil {
			panic(err)
		}

		writer := bufio.NewWriter(file)
		writer.WriteString(urlAddress + "\n")
		writer.Flush()

		if urlAddress != "" {
			checkUrl, err := CheckUrl(urlAddress)
			if err != nil {
				err = CreateMD5Url(urlAddress)
				if err != nil {
					fmt.Printf("[Error] Insert Data Error %s", err)
					return
				}
			} else {
				_, err := UpdateUrl(urlAddress)
				if err != nil {
					fmt.Printf("[Error] Insert Data Error %s", err)
					return
				}

				urlSha1 = checkUrl.Sha1Str
			}
		}

		if urlSha1 != "" {
			checkUrl, err := CheckUrlSha1(urlSha1)
			if err != nil {
				return
			}

			urlAddress = checkUrl.Urls

		}

		proxy := GetClashProxy(urlAddress)
		proxyGroup := setClashPg(getClashProxyName(proxy))

		c.String(http.StatusOK, "%s", header+"\nProxy:\n"+proxy+proxyGroup+foot)
	})
	r.Run(":9001")

}
