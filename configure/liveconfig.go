package configure

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/kr/pretty"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

/*
{
  "server": [
    {
      "appname": "live",
      "live": true,
	  "hls": true,
	  "static_push": []
    }
  ]
}
*/

// Application 客户端对象
type Application struct {
	AppName    string   `mapstructure:"appName"`
	Live       bool     `mapstructure:"live"`
	Hls        bool     `mapstructure:"hls"`
	Flv        bool     `mapstructure:"flv"`
	Api        bool     `mapstructure:"api"`
	StaticPush []string `mapstructure:"static_push"`
}

type Applications []Application

type JWT struct {
	Secret    string `mapstructure:"secret"`
	Algorithm string `mapstructure:"algorithm"`
}
type ServerCfg struct {
	Level           string       `mapstructure:"level"`
	ConfigFile      string       `mapstructure:"config_file"`
	FLVArchive      bool         `mapstructure:"flv_archive"`
	FLVDir          string       `mapstructure:"flv_dir"`
	RTMPNoAuth      bool         `mapstructure:"rtmp_noauth"`
	RTMPAddr        string       `mapstructure:"rtmp_addr"`
	HTTPFLVAddr     string       `mapstructure:"httpflv_addr"`
	HLSAddr         string       `mapstructure:"hls_addr"`
	HLSKeepAfterEnd bool         `mapstructure:"hls_keep_after_end"`
	APIAddr         string       `mapstructure:"api_addr"`
	RedisAddr       string       `mapstructure:"redis_addr"`
	RedisPwd        string       `mapstructure:"redis_pwd"`
	ReadTimeout     int          `mapstructure:"read_timeout"`
	WriteTimeout    int          `mapstructure:"write_timeout"`
	EnableTLSVerify bool         `mapstructure:"enable_tls_verify"`
	GopNum          int          `mapstructure:"gop_num"`
	JWT             JWT          `mapstructure:"jwt"`
	Server          Applications `mapstructure:"server"`
}

// default config
var defaultConf = ServerCfg{
	ConfigFile:      "livego.yaml",
	FLVArchive:      false,
	RTMPNoAuth:      false,
	RTMPAddr:        ":1935",
	HTTPFLVAddr:     ":7001",
	HLSAddr:         ":7002",
	HLSKeepAfterEnd: false,
	APIAddr:         ":8090",
	WriteTimeout:    10,
	ReadTimeout:     10,
	EnableTLSVerify: true,
	GopNum:          1,
	Server: Applications{{
		AppName:    "live",
		Live:       true,
		Hls:        true,
		Flv:        true,
		Api:        true,
		StaticPush: nil,
	}},
}

var (
	Config = viper.New()

	// BypassInit can be used to bypass the init() function by setting this
	// value to True at compile time.
	//
	// go build -ldflags "-X 'livego/configure.BypassInit=true'" -o livego main.go
	BypassInit = "" // 是否跳过init()
)

func initLog() {
	if l, err := log.ParseLevel(Config.GetString("level")); err == nil {
		log.SetLevel(l)
		log.SetReportCaller(l == log.DebugLevel)
	}
}

func init() {
	if BypassInit == "" {
		initDefault()
	}
}

func initDefault() {
	defer Init()

	// Default config 目的是在没有配置文件的情况下，使用默认配置
	b, _ := json.Marshal(defaultConf)              // 将defaultConf结构体转换为json格式的字节数组b
	defaultConfig := bytes.NewReader(b)            // 把b封装为bytes.Reader类型的defaultConfig
	viper.SetConfigType("json")                    // 设置配置类型为json格式
	_ = viper.ReadConfig(defaultConfig)            // 读取配置文件
	_ = Config.MergeConfigMap(viper.AllSettings()) // 将读取到的配置合并到全局配置中

	// Flags
	pflag.String("rtmp_addr", ":1935", "RTMP server listen address")
	pflag.Bool("enable_rtmps", false, "enable server session RTMPS")
	pflag.String("rtmps_cert", "server.crt", "cert file path required for RTMPS")
	pflag.String("rtmps_key", "server.key", "key file path required for RTMPS")
	pflag.String("httpflv_addr", ":7001", "HTTP-FLV server listen address")
	pflag.String("hls_addr", ":7002", "HLS server listen address")
	pflag.String("api_addr", ":8090", "HTTP manage interface server listen address")
	pflag.String("config_file", "livego.yaml", "configure filename")
	pflag.String("level", "info", "Log level")
	pflag.Bool("hls_keep_after_end", false, "Maintains the HLS after the stream ends")
	pflag.String("flv_dir", "tmp", "output flv file at flvDir/APP/KEY_TIME.flv")
	pflag.Int("read_timeout", 10, "read time out")
	pflag.Int("write_timeout", 10, "write time out")
	pflag.Int("gop_num", 1, "gop num")
	pflag.Bool("enable_tls_verify", true, "Use system root CA to verify RTMPS connection, set this flag to false on Windows")
	pflag.Parse()
	_ = Config.BindPFlags(pflag.CommandLine)

	// File
	Config.SetConfigFile(Config.GetString("config_file"))
	Config.AddConfigPath(".")
	err := Config.ReadInConfig()
	if err != nil {
		log.Warning(err)
		log.Info("Using default config")
	} else {
		Config.MergeInConfig()
	}

	// Environment
	replacer := strings.NewReplacer(".", "_")
	Config.SetEnvKeyReplacer(replacer)
	Config.AllowEmptyEnv(true)
	Config.AutomaticEnv()

	// Log
	initLog()

	// Print final config
	c := ServerCfg{}
	Config.Unmarshal(&c)
	log.Debugf("Current configurations: \n%# v", pretty.Formatter(c))
}

func CheckAppName(appname string) bool {
	apps := Applications{}
	Config.UnmarshalKey("server", &apps)
	for _, app := range apps {
		if app.AppName == appname {
			return app.Live
		}
	}
	return false
}

func GetStaticPushUrlList(appname string) ([]string, bool) {
	apps := Applications{}
	Config.UnmarshalKey("server", &apps)
	for _, app := range apps {
		if (app.AppName == appname) && app.Live {
			if len(app.StaticPush) > 0 {
				return app.StaticPush, true
			} else {
				return nil, false
			}
		}
	}
	return nil, false
}
