package app

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type conf struct {
	CurrentContext string                    `yaml:"current-context"`
	Contexts       map[string]*ConfigContext `yaml:"contexts"`
}

type ConfigContext struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
	Mckscli   *CliConfig
	Spidercli *CliConfig
}

type CliConfig struct {
	ServerAddr string `yaml:"server_addr"`
	Timeout    string `yaml:"timeout"`
	Tls        struct {
		TlsCa string `yaml:"tls_ca"`
	} `yaml:"tls"`
	Interceptors struct {
		AuthJwt struct {
			JwtToken string `yaml:"jwt_token"`
		} `yaml:"auth_jwt"`
		Opentracing struct {
			Jaeger struct {
				Endpoint    string `yaml:"endpoint"`
				ServiceName string `yaml:"service_name"`
				SampleRate  string `yaml:"sample_rate"`
			} `yaml:"jaeger"`
		} `yaml:"opentracing"`
	} `yaml:"interceptors"`
}

var (
	Config *conf
)

func (self *conf) WriteConfig() error {

	if b, err := yaml.Marshal(self); err != nil {
		return err
	} else {
		os.WriteFile(viper.ConfigFileUsed(), b, os.ModePerm)
	}
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

func OnConfigInitialize(cfgFile string) error {

	dir := filepath.Join(HomeDir(), ".cbadm")
	viper.AddConfigPath(dir)
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	// set default
	viper.SetDefault("current-context", "")
	viper.SetDefault("contexts.local.name", "local")
	viper.SetDefault("contexts.local.mckscli.server_addr", "127.0.0.1:50254")
	viper.SetDefault("contexts.local.mckscli.timeout", "1000s")
	viper.SetDefault("contexts.local.mckscli.interceptors.opentracing.jaeger", map[string]string{
		"endpoint":     "localhost:6834",
		"service_name": "mcks grpc client",
		"sample_rate":  "1",
	})
	viper.SetDefault("contexts.local.spidercli.server_addr", "127.0.0.1:2048")
	viper.SetDefault("contexts.local.spidercli.timeout", "1000s")
	viper.SetDefault("contexts.local.spidercli.interceptors.opentracing.jaeger", map[string]string{
		"endpoint":     "localhost:6832",
		"service_name": "spider grpc client",
		"sample_rate":  "1",
	})
	// read a config file
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		//return err

		// the default config save to "${HOME}/.cbctl/config"
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.MkdirAll(dir, os.ModePerm)
		}
		if _, err := os.Stat(filepath.Join(dir, "config")); os.IsNotExist(err) {
			os.WriteFile(filepath.Join(dir, "config"), []byte{}, os.ModePerm)
		}
		if err := viper.WriteConfig(); err != nil {
			fmt.Println(err)
		}
	}

	// unmarshal
	if err := viper.Unmarshal(&Config,
		viper.DecoderConfigOption(func(decoderConfig *mapstructure.DecoderConfig) {
			decoderConfig.TagName = "yaml"
		})); err != nil {
		return fmt.Errorf("unable to decode into config struct, %v", err)
	}

	// current-context
	if Config.Contexts[Config.CurrentContext] == nil {
		Config.CurrentContext = func() string {
			if len(Config.Contexts) > 0 {
				for k := range Config.Contexts {
					return k
				}
			}
			return ""
		}()
	}
	if Config.CurrentContext == "" {
		return fmt.Errorf("unable to find current context")
	}

	return nil

}

func (self *conf) GetCurrentContext() *ConfigContext {
	return self.Contexts[self.CurrentContext]
}

func HomeDir() string {

	if runtime.GOOS == "windows" {
		home := os.Getenv("HOME")
		homeDriveHomePath := ""
		if homeDrive, homePath := os.Getenv("HOMEDRIVE"), os.Getenv("HOMEPATH"); len(homeDrive) > 0 && len(homePath) > 0 {
			homeDriveHomePath = homeDrive + homePath
		}
		userProfile := os.Getenv("USERPROFILE")

		// Return first of %HOME%, %HOMEDRIVE%/%HOMEPATH%, %USERPROFILE% that contains a `.cbctl\config` file.
		// %HOMEDRIVE%/%HOMEPATH% is preferred over %USERPROFILE% for backwards-compatibility.
		for _, p := range []string{home, homeDriveHomePath, userProfile} {
			if len(p) == 0 {
				continue
			}
			if _, err := os.Stat(filepath.Join(p, ".cbadm", "config")); err != nil {
				continue
			}
			return p
		}

		firstSetPath := ""
		firstExistingPath := ""

		// Prefer %USERPROFILE% over %HOMEDRIVE%/%HOMEPATH% for compatibility with other auth-writing tools
		for _, p := range []string{home, userProfile, homeDriveHomePath} {
			if len(p) == 0 {
				continue
			}
			if len(firstSetPath) == 0 {
				// remember the first path that is set
				firstSetPath = p
			}
			info, err := os.Stat(p)
			if err != nil {
				continue
			}
			if len(firstExistingPath) == 0 {
				// remember the first path that exists
				firstExistingPath = p
			}
			if info.IsDir() && info.Mode().Perm()&(1<<(uint(7))) != 0 {
				// return first path that is writeable
				return p
			}
		}

		// If none are writeable, return first location that exists
		if len(firstExistingPath) > 0 {
			return firstExistingPath
		}

		// If none exist, return first location that is set
		if len(firstSetPath) > 0 {
			return firstSetPath
		}

		// We've got nothing
		return ""
	}
	return os.Getenv("HOME")
}
