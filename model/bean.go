package model

type NodeBean struct {
	Server   string      `json:"add"`
	NodeType string      `json:"node_type"`
	Host     string      `json:"host"`
	Uuid     string      `json:"id"`
	Network  string      `json:"net"`
	Path     string      `json:"path"`
	Port     interface{} `json:"port"`
	Name     string      `json:"ps"`
	Tls      string      `json:"tls"`
	//V        int    `json:"v"`
	AlterId       interface{} `json:"aid"`
	Type          string      `json:"type"`
	Cipher        string      `json:"cipher"`
	Password      string      `json:"password"`
	Protocol      string
	Protocolparam string
	Obfs          string
	Obfsparam     string
}

type NodeOrm struct {
	Id         int64
	Urls       string `xrom:"unique"`
	CreateTime string
	UpdateTime string
	Sha1Str    string `xrom:"unique"`
	Hide       bool
}
