package nodedao

import (
	. "awesomeProject/model"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"os"
	"time"
)

var x *xorm.Engine

func init() {
	var err error
	x, err = xorm.NewEngine("sqlite3", "./md5Url.db")
	if err != nil {
		fmt.Printf("[Error] Create Engine Error %s", err)
		return
	}

	if err = x.Sync2(new(NodeOrm)); err != nil {
		fmt.Printf("[Error] Create Engine Error %s", err)
		return
	}

	file, err := os.OpenFile("./sql.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		fmt.Printf("[Error] Create File Error %s", err)
		return
	}

	//开启日志
	x.SetLogger(xorm.NewSimpleLogger(file))
	x.ShowSQL(true)

	//设置缓存
	cache := xorm.NewLRUCacher(xorm.NewMemoryStore(), 1000)
	x.SetDefaultCacher(cache)

}

func Sha1Encode(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func CreateMD5Url(url string) error {

	create := time.Now().Format("2006-01-02 15:04:05")

	md5Str := Sha1Encode(url)

	_, err := x.Insert(&NodeOrm{
		Urls:       url,
		CreateTime: create,
		Sha1Str:    md5Str,
		Hide:       false,
	})

	return err
}

func CheckUrl(url string) (*NodeOrm, error) {

	a := &NodeOrm{Urls: url}

	has, err := x.Get(a)
	if err != nil {
		fmt.Printf("[Error]  Get Engine Error %s", err)
		return nil, err
	} else if !has {
		err = errors.New("[Error] Not found")
		return nil, err
	}

	return a, err
}

func CheckUrlSha1(sha1 string) (*NodeOrm, error) {

	a := &NodeOrm{Sha1Str: sha1}

	has, err := x.Get(a)
	if err != nil {
		fmt.Printf("[Error]  Get Engine Error %s", err)
		return nil, err
	} else if !has {
		err = errors.New("[Error] Not found")
		return nil, err
	}

	return a, err
}

func UpdateUrl(s string) (*NodeOrm, error) {

	updateTime := time.Now().Format("2006-01-02 15:04:05")

	orm, err := CheckUrl(s)
	if err != nil {
		fmt.Printf("[Error]  Get Engine Error %s", err)
		return orm, err
	}

	orm.UpdateTime = updateTime

	_, err = x.ID(orm.Id).Update(orm)

	return orm, err
}
