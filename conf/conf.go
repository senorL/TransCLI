package conf

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/validator.v2"
	"gopkg.in/yaml.v2"
)

var (
	conf *Config
	once sync.Once
)

type Config struct {
	TModle TModle `yaml:"TModle"`
}

type TModle struct {
	APPID  string `yaml:"appid" validate:"nonzero"`
	APIKey string
}

// GetConf returns a singleton configuration instance.
func GetConf() *Config {
	once.Do(func() {
		err := initConf()
		if err != nil {
			panic(fmt.Sprintf("初始化配置文件失败: %v", err))
		}
	})
	return conf
}

// initConf initializes the configuration by loading from a YAML file.
func initConf() error {
	welcome()
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前工作目录失败: %w", err)
	}

	// 构建配置文件路径
	confFileRelPath := filepath.Join(wd, "trans.yaml")
	content, err := os.ReadFile(confFileRelPath)
	if err != nil {
		if os.IsNotExist(err) {
			// 如果配置文件不存在，则创建默认配置文件
			fmt.Println("加载配置文件失败，创建默认配置文件...")
			err = createDefaultConfig(confFileRelPath)
			if err != nil {
				return fmt.Errorf("创建默认配置文件失败: %w", err)
			}
			fmt.Println("默认配置文件创建成功！")
			// 重新读取配置文件
			content, err = os.ReadFile(confFileRelPath)
			if err != nil {
				return fmt.Errorf("读取配置文件失败: %w", err)
			}
		} else {
			return fmt.Errorf("读取配置文件失败: %w", err)
		}
	}

	conf = new(Config)
	err = yaml.Unmarshal(content, conf)
	if err != nil {
		return fmt.Errorf("解析 YAML 失败: %w", err)
	}

	if err := validator.Validate(conf); err != nil {
		return fmt.Errorf("验证配置文件失败: %w", err)
	}

	return nil
}

// createDefaultConfig 创建默认配置文件
func createDefaultConfig(filePath string) error {
	defaultConfig := Config{
		TModle: TModle{
			APPID:  "20241123002209596",
			APIKey: "FfVqlwuCY7TIUXzImag8",
		},
	}

	content, err := yaml.Marshal(&defaultConfig)
	if err != nil {
		return fmt.Errorf("序列化默认配置失败: %w", err)
	}

	err = os.WriteFile(filePath, content, 0644)
	if err != nil {
		return fmt.Errorf("写入默认配置文件失败: %w", err)
	}

	return nil
}

func welcome() {
	asciiArt := `

	████████╗██████╗  █████╗ ███╗   ██╗███████╗       ██████╗██╗     ██╗
	╚══██╔══╝██╔══██╗██╔══██╗████╗  ██║██╔════╝      ██╔════╝██║     ██║
	   ██║   ██████╔╝███████║██╔██╗ ██║███████╗█████╗██║     ██║     ██║
	   ██║   ██╔══██╗██╔══██║██║╚██╗██║╚════██║╚════╝██║     ██║     ██║
	   ██║   ██║  ██║██║  ██║██║ ╚████║███████║      ╚██████╗███████╗██║
	   ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═══╝╚══════╝       ╚═════╝╚══════╝╚═╝
`
	fmt.Print(asciiArt)
}
