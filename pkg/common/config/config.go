package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"github.com/gookit/validate"
	"github.com/spf13/viper"
)

type params struct {
	InputDir      string `mapstructure:"INPUT_DIR"      validate:"required|minLen:5"`   // Input data directory for archiving
	OutputDir     string `mapstructure:"OUTPUT_DIR"     validate:"required|minLen:5"`   // Directory for holding output archive data
	DownloadDir   string `mapstructure:"DOWNLOAD_DIR"   validate:"required|minLen:5"`   // Directory for downloaded files from minIO
	Port          string `mapstructure:"PORT"           validate:"customPortValidator"` // The localhost port on which HTTP requests are listened
	DBService     string `mapstructure:"ORA_SERVICE"    validate:"required|minLen:3"`   //
	DBUser        string `mapstructure:"ORA_USERNAME"   validate:"required|minLen:3"`   //
	DBPassword    string `mapstructure:"ORA_PASSWORD"   validate:"required|minLen:3"`   //
	DBServer      string `mapstructure:"ORA_SERVER"     validate:"required|minLen:3"`   //
	DBPort        string `mapstructure:"ORA_PORT"       validate:"required|minLen:3"`   //
	DBWallet      string `mapstructure:"ORA_WALLET"`
	ConnString    string `mapstructure:"ORA_CONN_STRING"` // Oracle connection string with wallet
	SrcConnString string `mapstructure:"SSIRB_CONN_STRING"`
}

type Config struct {
	Params params
	Path   string
}

var onceCfg sync.Once

var singleInstanceCfg *Config

func GetAppName() string {
	path, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Base(path)
}

func GetConfig(cfgPath ...string) (Config, error) {
	var err error

	if singleInstanceCfg == nil {
		singleInstanceCfg = &Config{}

		f := func(cfg *Config, err error) {
			onceCfg.Do(
				func() {
					fmt.Println("Creating single instance of config now")
					//singleInstanceCfg = &Config{}

					// load and parse config
					*cfg, err = loadConfig(cfgPath[0])
					if err != nil {
						log.Fatalln("Failed at loading init config", err)
						return
					}
				})
		}
		f(singleInstanceCfg, err)
		log.Printf("%#v\n", *singleInstanceCfg) // config is now setup
	}

	if err != nil {
		return Config{}, err
	}
	return *singleInstanceCfg, nil
}

func loadConfig(cfgPath string) (cfg Config, err error) {
	cfg.Path = cfgPath
	log.Println(os.Getwd())
	viper.AddConfigPath("/etc/dr")
	viper.AddConfigPath("/tmp/dr")
	viper.AddConfigPath("./")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../../")
	viper.SetConfigName(GetConfigFilename())
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&cfg.Params)

	return
}

func GetConfigPaths() []string {
	var paths []string
	paths = append(paths, "/etc/dr")
	paths = append(paths, "/tmp/dr")

	return paths
}

func GetConfigFilename() string {
	cfgName := strings.Join([]string{GetAppName(), "cfg"}, ".")

	return cfgName
}

// CustomPortValidator validates port in config struct, valid syntax is `:3055`
func (p params) CustomPortValidator(port string) bool {
	return len(port) > 4 && strings.HasPrefix(port, ":")
}

// Equals compares 2 configs each other
func (c Config) Equals(other Config) bool {
	return reflect.DeepEqual(c, other)
}

func (c Config) Validate() error {
	v := validate.Struct(c.Params)
	if !v.Validate() {
		return errors.New(v.Errors.Error()) // all error messages
	}
	return nil
}
