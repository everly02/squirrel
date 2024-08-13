package squirrel

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	*viper.Viper
}

// NewConfig 创建一个新的配置实例
func NewConfig(configName, configPath string) *Config {
	v := viper.New()

	// 设置配置文件名（不带扩展名）
	v.SetConfigName(configName)
	// 设置查找配置文件的路径
	v.AddConfigPath(configPath)
	// 设置配置文件类型（如果是从环境变量或远程配置中获取则不需要）
	v.SetConfigType("yaml")

	// 读取环境变量的前缀（可选）
	v.SetEnvPrefix("APP")
	// 自动读取环境变量
	v.AutomaticEnv()

	// 尝试读取配置文件
	err := v.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	return &Config{Viper: v}
}

// GetString 获取字符串类型的配置值
func (c *Config) GetString(key string) string {
	return c.Viper.GetString(key)
}

// GetInt 获取整型配置值
func (c *Config) GetInt(key string) int {
	return c.Viper.GetInt(key)
}

// GetBool 获取布尔型配置值
func (c *Config) GetBool(key string) bool {
	return c.Viper.GetBool(key)
}

// Set 设置配置项的值
func (c *Config) Set(key string, value interface{}) {
	c.Viper.Set(key, value)
}

// SaveConfig 保存当前配置到文件
func (c *Config) SaveConfig() error {
	return c.WriteConfig()
}

// SaveConfigAs 保存当前配置到指定文件
func (c *Config) SaveConfigAs(filename string) error {
	return c.WriteConfigAs(filename)
}
