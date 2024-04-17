package configloader

import (
	"bytes"
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func FromFS(embedFS *embed.FS, patterns []string, conf interface{}) (err error) {
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetConfigType("toml")

	for i := range patterns {
		matches, er := fs.Glob(embedFS, patterns[i])
		if er != nil {
			err = er
			return
		}

		for i2 := range matches {
			fileData, er := embedFS.ReadFile(matches[i2])
			if er != nil {
				err = er
				return
			}

			if err = viper.MergeConfig(bytes.NewReader(fileData)); err != nil {
				return
			}
		}
	}

	err = viper.Unmarshal(conf)
	return
}

func FromFiles(patterns []string, conf interface{}) (err error) {
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetConfigType("toml")

	for i := range patterns {
		matches, er := filepath.Glob(patterns[i])
		if er != nil {
			err = er
			return
		}

		for i2 := range matches {
			if _, err := os.Stat(matches[i2]); err == nil {
				viper.SetConfigFile(matches[i2])
				if err := viper.MergeInConfig(); err != nil {
					return err
				}
			}
		}
	}

	err = viper.Unmarshal(conf)
	return
}
