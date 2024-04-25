package conf

import (
	"fmt"

	log "github.com/cihub/seelog"
	"github.com/spf13/viper"
)

type Opts struct {
	AdminRole string //能够管理配置的角色
	CertPath  string
	//吊销证书路径

	RbacModelPath string
	VuePath       string
	IP            string
	Port          int
	//log的级别，DEBUG INFO ERROR FATAL
	Log string
	//pfx的默认密码
	DefaultPWD      string
	DefaultRoleName string
	Authn           string
	Oidc            struct {
		Issuer        string
		ClientID      string
		Secret        string
		Callback      string
		CallbackProto string
	}
	LDAP struct {
		URL          string
		BindDN       string
		BindPassword string
		BaseDN       string
	}
	DB struct {
		DBDriver     string
		DBConnection string
	}
	CA struct {
		DefaultDN string
		CRLURL    string
	}
}

var Default *Opts //系统全局配置信息
var CurrentHost string

func InitConfig(configFile string) *Opts {
	viper.SetConfigFile(configFile)
	viper.AddConfigPath(".")
	viper.AllowEmptyEnv(true)
	err := viper.ReadInConfig() // Find and read the config file
	viper.SetEnvPrefix("CA")
	viper.AutomaticEnv()
	//	viper.BindEnv("port")
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %v \r\n", err))
	}
	config := &Opts{}
	err = viper.Unmarshal(config)
	if err != nil {
		log.Error("读取配置错误:", err)
	}
	Default = config

	return config
}
