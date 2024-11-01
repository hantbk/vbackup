package config

import (
	"fmt"

	"io/ioutil"
	"os"
	pathu "path"
	"path/filepath"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func ReadConfig(path string) (c *config.Config, err error) {
	path = fileutil.FixPath(path)
	if "" == path {
		c = defaultConfig()
		path = fileutil.FixPath(pathu.Join(fileutil.HomeDir(), string(filepath.Separator), ".vbackup"))
		c.Data.CacheDir = pathu.Join(path, string(filepath.Separator), "cache")
		c.Logger.LogPath = pathu.Join(path, string(filepath.Separator), "log")
		c.Data.DbDir = pathu.Join(path, string(filepath.Separator), "db")
		confpath := pathu.Join(path, string(filepath.Separator), "conf")
		if !fileutil.Exist(confpath) {
			err = os.MkdirAll(confpath, 0777)
			if err != nil {
				return nil, err
			}
		}
		fpath := pathu.Join(confpath, string(filepath.Separator), "app.yml")
		if !fileutil.Exist(fpath) {
			fmt.Printf("Loading default configuration: %s\n", path)
			bytes, _ := yaml.Marshal(c)
			err = ioutil.WriteFile(fpath, bytes, 0666)
			if err != nil {
				return nil, err
			}
			return c, nil
		} else {
			c, err = ReadConfig(path)
			if err != nil {
				return nil, err
			}
			return c, nil
		}
	} else {
		// Read configuration file from the path
		v := viper.New()
		v.SetConfigName("app")
		v.SetConfigType("yaml")
		realDir := fileutil.ReplaceHomeDir(pathu.Join(path, string(filepath.Separator), "conf"))
		if exists := fileutil.Exist(realDir); !exists {
			return nil, fmt.Errorf("Failed to read configuration file %s: the conf directory does not exist; do not include 'conf' when configuring.", path)
		}
		v.AddConfigPath(realDir)
		if err = v.ReadInConfig(); err != nil {
			fmt.Println(fmt.Errorf("Failed to read configuration file %s: %s, %s.", path, realDir, err.Error()))
			return nil, err
		}
		if err = v.Unmarshal(&c); err != nil {
			fmt.Println(fmt.Errorf("Failed to read configuration file %s: %s.", path, err.Error()))
			return nil, err
		}
		if c.Logger.LogPath == "" {
			c.Logger.LogPath = pathu.Join(path, string(filepath.Separator), "log")
		}
		if c.Data.DbDir == "" {
			c.Data.DbDir = pathu.Join(path, string(filepath.Separator), "db")
		}
		if "" == c.Data.CacheDir {
			c.Data.CacheDir = pathu.Join(path, string(filepath.Separator), "cache")
		}
		fmt.Println(fmt.Sprintf("Configuration loading completed: %s", path))
	}
	return c, nil
}

// defaultConfig retrieves the default configuration.
func defaultConfig() *config.Config {
	c := &config.Config{}
	c.Server.Name = "vbackup"
	c.Data.NoCache = false
	c.Server.Debug = false
	c.Logger.Level = "info"
	c.Jwt.Key = "vbackup-secret"
	c.Jwt.MaxAge = 1800
	c.Server.Bind.Port = 8012
	c.Server.Bind.Host = ""
	c.Prometheus.Enable = false
	return c
}
