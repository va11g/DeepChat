package config

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

var Configure *Config

var (
	defaultUrl            = "https://api.deepseek.com"
	defaultModel          = "deepseek-chat"
	defaultApiKey         = ""
	defaultLogDir         = "./"
	defaultLogLevel       = "info"
	defaultshardNum       = 1024
	defaultChanBufferSize = 10
	configFile            = "./deepseek.conf"
)

type Config struct {
	ConfFile       string
	Model          string
	Url            string
	ApiKey         string
	LogDir         string
	LogLevel       string
	ShardNum       int
	ChanBufferSize int
	Others         map[string]string
}

type CfgError struct {
	message string
}

func (cErr *CfgError) Error() string {
	return cErr.message
}

func flagInit(cfg *Config) {
	// 标志的值，标志的名称，标志的默认值，标志的使用说明
	flag.StringVar(&(cfg.ConfFile), "config", configFile, "Appoint a config file: such as /etc/redis.conf")
	flag.StringVar(&(cfg.Model), "model", defaultModel, "model type:deepseek-chat/deepseek-reasoner")
	flag.StringVar(&(cfg.Url), "Url", defaultUrl, "DeepSeek Api Url")
	flag.StringVar(&(cfg.ApiKey), "ApiKey", defaultApiKey, "API key that allows you to authenticate your identity")
	flag.StringVar(&(cfg.LogDir), "logdir", defaultLogDir, "Create log directory: default is /tmp")
	flag.StringVar(&(cfg.LogLevel), "loglevel", defaultLogLevel, "Create log level: default is info")
	flag.IntVar(&(cfg.ChanBufferSize), "chanBufSize", defaultChanBufferSize, "set the buffer size of channels in PUB/SUB commands. ")
}

// 初始化并检查config
func SetUp() (*Config, error) {
	cfg := &Config{
		ConfFile:       configFile,
		Model:          defaultModel,
		Url:            defaultUrl,
		ApiKey:         defaultApiKey,
		LogDir:         defaultLogDir,
		LogLevel:       defaultLogLevel,
		ShardNum:       defaultshardNum,
		ChanBufferSize: defaultChanBufferSize,
		Others:         make(map[string]string),
	}
	flagInit(cfg)
	flag.Parse()
	Configure = cfg
	return cfg, nil
}

// 解析config文件
func (cfg *Config) Parse(cfgFile string) error {
	fl, err := os.Open(cfgFile)
	if err != nil {
		return err
	}
	defer func() error {
		err := fl.Close()
		if err != nil {
			return err
		}
		return nil
	}()

	reader := bufio.NewReader(fl)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
		// 注释
		if len(line) > 0 && line[0] == '#' {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 2 {
			cfgName := strings.ToLower(fields[0])
			switch cfgName {
			case "model":
				cfg.Model = fields[1]
			case "url":
				if ip := net.ParseIP(fields[1]); ip == nil {
					ipErr := &CfgError{
						message: fmt.Sprintf("Given ip address %s is invalid", cfg.Url),
					}
					return ipErr
				}
				cfg.Url = fields[1]
			case "apikey":
				cfg.ApiKey = fields[1]
			case "logdir":
				cfg.LogDir = strings.ToLower(fields[1])
			case "loglevel":
				cfg.LogLevel = strings.ToLower(fields[1])
			case "shardnum":
				cfg.ShardNum, err = strconv.Atoi(fields[1])
				if err != nil {
					fmt.Println("ShardNum should be a number. Get: ", fields[1])
					panic(err)
				}
			default:
				cfg.Others[cfgName] = fields[1]
			}
		}

		if err == io.EOF {
			break
		}

	}
	return nil
}

// func (cfg *Config) ParseConfigJson(path string) error {
// 	data, err := os.ReadFile(path)
// 	if err != nil {
// 		return errors.New("json config not exist")
// 	}
// 	err = json.Unmarshal(data, cfg)
// 	if cfg.NodeID <= 0 {
// 		panic("NodeID not set")
// 	}
// 	if err != nil {
// 		return errors.New("Invalid config file fields. ")
// 	}
// 	if cfg.RaftAddr == "" {
// 		cfg.RaftAddr = strings.Split(cfg.PeerAddrs, ",")[cfg.NodeID-1]
// 	}
// 	log.Println("RaftAddr = ", cfg.RaftAddr)
// 	// we only support a single database in cluster mode
// 	cfg.Databases = 1
// 	return nil
// }

// func (cfg *Config) Parse(cfgFile string) error {
// 	// 打开 JSON 配置文件
// 	fl, err := os.Open(cfgFile)
// 	if err != nil {
// 		return err
// 	}
// 	defer fl.Close()

// 	reader := bufio.NewReader(fl)

// 	var jsonData map[string]interface{}

// 	decoder := json.NewDecoder(reader)
// 	err = decoder.Decode(&jsonData)
// 	if err != nil {
// 		return fmt.Errorf("error decoding JSON: %v", err)
// 	}

// 	// 逐个映射字段
// 	for key, value := range jsonData {
// 		switch strings.ToLower(key) {
// 		case "model":
// 			cfg.Model = fmt.Sprintf("%v", value)
// 		case "url":
// 			cfg.Url = fmt.Sprintf("%v", value)
// 		case "apikey":
// 			cfg.ApiKey = fmt.Sprintf("%v", value)
// 		case "logdir":
// 			cfg.LogDir = fmt.Sprintf("%v", value)
// 		case "loglevel":
// 			cfg.LogLevel = fmt.Sprintf("%v", value)
// 		case "shardnum":
// 			// 将字符串转换为整数
// 			shardNum, err := strconv.Atoi(fmt.Sprintf("%v", value))
// 			if err != nil {
// 				return fmt.Errorf("invalid value for shardnum: %v", value)
// 			}
// 			cfg.ShardNum = shardNum
// 		case "chanbuffersize":
// 			// 将字符串转换为整数
// 			chanBufferSize, err := strconv.Atoi(fmt.Sprintf("%v", value))
// 			if err != nil {
// 				return fmt.Errorf("invalid value for chanBufferSize: %v", value)
// 			}
// 			cfg.ChanBufferSize = chanBufferSize
// 		default:
// 			// 其他字段存储在 Others 字段中
// 			if cfg.Others == nil {
// 				cfg.Others = make(map[string]string)
// 			}
// 			cfg.Others[key] = fmt.Sprintf("%v", value)
// 		}
// 	}

// 	// 校验 URL 是否有效
// 	if net.ParseIP(cfg.Url) == nil {
// 		return fmt.Errorf("invalid URL: %s", cfg.Url)
// 	}

// 	// 解析成功后返回
// 	return nil
// }
