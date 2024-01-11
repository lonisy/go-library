package viper

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	etcdtools "go-library/config/etcd"
	clientv3 "go.etcd.io/etcd/client/v3"
	"os"
	"strings"
	"time"
)

var AppConfig *viper.Viper

func InitConfig() {
	key := "/youzu-ad/services/apipool/click/config-develop.json"
	if os.Getenv("APP_ENV") == "release" || os.Getenv("APP_ENV") == "pord" {
		key = "/youzu-ad/services/apipool/click/config.json"
	}
	log.Println(key)
	AppConfig = viper.New()
	AppConfig.SetConfigType("json")
	Etcd := etcdtools.NewEtcdTools().Init(clientv3.Config{
		Endpoints:   strings.Split("ad-etcd-s1.youzu.com:2379,ad-etcd-s2.youzu.com:2379,ad-etcd-s3.youzu.com:2379,ad-etcd-develop.adm5.uuzu.com:2379", ","),
		DialTimeout: 5 * time.Second,
	})
	Etcd.LoadData(key, func(x, y interface{}) {
		if err := AppConfig.ReadConfig(bytes.NewBufferString(cast.ToString(y))); err == nil {
			log.Println("LoadData...")
			if err := AppConfig.WriteConfigAs(".config.json"); err != nil {
				log.Println(err.Error())
			}
		} else {
			log.Panicln(err.Error())
		}
	})
	Etcd.WatchData(key, func(x, y interface{}) {
		if err := AppConfig.ReadConfig(bytes.NewBufferString(cast.ToString(y))); err == nil {
			log.Println("WatchData...")
			if err := AppConfig.WriteConfigAs(".config.json"); err != nil {
				log.Println(err.Error())
			}
		} else {
			log.Panicln(err.Error())
		}
	})
}
