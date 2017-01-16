package utils

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/toolkits/str"
)

func HashIt(passwd string) (hashed string) {
	salt := viper.GetString("salt")
	log.Debugf("salf is %v", salt)
	if salt == "" {
		log.Error("salt is empty, please check your conf")
	}
	hashed = str.Md5Encode(salt + passwd)
	return
}
