package conf

import (
	"encoding/json"
	"log"
	"os"
)

// GlobalConfig project config info
// maybe saved in os.Env or k8s service config yaml file.
var GlobalConfig *ProjectConfig

// envConfigName env config name for project
const envConfigName = "USER_GROWTH_CONFIG"

// ProjectConfig project's config
type ProjectConfig struct {
	Db struct {
		Engine          string // mysql
		Username        string // root
		Password        string // 12345678
		Host            string // localhost
		Port            int    // 3306
		Database        string // user_growth
		Charset         string // utf8
		ShowSql         bool   // true
		MaxIdleConns    int    // 2  最多保持2个空闲连接，避免频繁创建连接
		MaxOpenConns    int    // 10 最多10个连接.
		ConnMaxLifetime int    // 30 minute
	}
	Cache struct{}
}

// LoadConfigs load global config info
func init() {
	//LoadEnvConfig()
	LoadConfigFromFile("D:\\GolandProjects\\user_growth\\conf\\conf.json")
}

// LoadEnvConfig load configs from env config with name of envConfigName, json format
func LoadEnvConfig() {
	pc := &ProjectConfig{}

	// load from os env
	if strConfigs := os.Getenv(envConfigName); len(strConfigs) > 0 {
		if err := json.Unmarshal([]byte(strConfigs), pc); err != nil {
			log.Fatalf("conf.LoadEnvConfig(%s) error=%s\n",
				envConfigName, err.Error())
			return
		}
	}
	if pc == nil || pc.Db.Username == "" { // no config info
		log.Fatalln("empty os.Getenv config ", envConfigName)
		return
	}
	GlobalConfig = pc
}

func LoadConfigFromFile(filePath string) {
	pc := &ProjectConfig{}
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("conf.LoadConfigFromFile(%s) error=%s\n", filePath, err.Error())
		return
	}
	if err := json.Unmarshal(file, pc); err != nil {
		log.Fatalf("conf.LoadConfigFromFile(%s) json.Unmarshal error=%s\n", filePath, err.Error())
		return
	}
	if pc == nil || pc.Db.Username == "" { // no config info
		log.Fatalln("empty config file ", filePath)
		return
	}
	GlobalConfig = pc
}
