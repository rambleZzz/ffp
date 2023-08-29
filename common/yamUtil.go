package common

import (
	"github.com/spf13/viper"
	"log"
)

var configYaml = `# AutoScan Yaml Config
thirdparty: # 路径必须配置在当前目录下，不可自定义当前执行文件以外的其他目录
  ObserverWardPath: /thirdparty/observerWard/observer_ward # observer_ward 可执行文件
  ObserverWardDir: /thirdparty/observerWard/ # observerWard 所在目录
  WebFingerprintPath: /thirdparty/observerWard/web_fingerprint_v3.json # fingerprintHub指纹库
  GeoLite2Path: /thirdparty/cdnCheck/GeoLite2-ASN.mmdb # ip相关
`

func ReadYaml() {
	YamlName := CurrentRunPath + "/config.yaml"
	if CheckFileIsExist(YamlName) == false {
		WriteFile(YamlName, configYaml)
	}
	v := viper.New()
	v.AddConfigPath(CurrentRunPath)
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		log.Println(err)
	}
	ObserverWardPath = v.GetString("thirdparty.ObserverWardPath")
	ObserverWardDir = v.GetString("thirdparty.ObserverWardDir")
	WebFingerprintPath = v.GetString("thirdparty.WebFingerprintPath")
	GeoLite2Path = v.GetString("thirdparty.GeoLite2Path")
}
