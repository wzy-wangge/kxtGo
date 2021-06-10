package tool

import "github.com/Unknwon/goconfig"

var GlobalConfig *goconfig.ConfigFile

type Config struct {
}

//加载配置
func (c *Config) loadConfigFile() error{
	cfg, err := goconfig.LoadConfigFile("conf.ini")
	if err != nil{
		return err
	}
	GlobalConfig = cfg
	return nil
}


func (c *Config) getCfg() (*goconfig.ConfigFile,error) {
	if GlobalConfig == nil {
		err := c.loadConfigFile()
		if err != nil {
			return nil,err
		}
	}
	return GlobalConfig,nil
}

// 获取块配置
func (c *Config) GetSection (section string) (map[string]string , error) {
	cfg,err := c.getCfg()
	if err != nil {
		return nil,err
	}

	sectionConfMap,err := cfg.GetSection(section)

	if err != nil {
		return nil, err
	}
	return sectionConfMap, nil
}