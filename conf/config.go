package conf

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/capyflow/allspark-go/ds"
)

// CenterConfig 中心配置
type CenterConfig struct {
	Port     int          `json:"port" toml:"port"`
	DBConfig *ds.DsConfig `json:"db_config" toml:"dbConfig"`
}

// LoadConfig 加载配置
func LoadConfig(confPath string) *CenterConfig {
	open, err := os.Open(confPath)
	if nil != err {
		panic(err)
	}
	defer open.Close()
	config := &CenterConfig{}
	toml.NewDecoder(open).Decode(config)
	return config
}
