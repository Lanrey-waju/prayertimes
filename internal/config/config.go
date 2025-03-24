package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

func InitConfig() {
	viper.SetDefault("location.country", "")
	viper.SetDefault("location.city", "")
	viper.SetDefault("location.latitude", 0.0)
	viper.SetDefault("location.longitude", 0.0)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(filepath.Join(os.Getenv("HOME"), ".config", "prayertimes"))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// config file does not exist
		} else {
			fmt.Printf("error reading config: %v", err)
			os.Exit(1)
		}
	}
}

// EnsureConfig ensures config file exists with usable values
func EnsureConfig(locationProvider func() (string, float64, float64)) {
	if viper.GetString("location.city") == "" {
		city, lat, lon := locationProvider()

		if city == "" {
			fmt.Println("error: location provider returned an empty city name")
			os.Exit(1)
		}

		// Set Values
		viper.Set("location.city", city)
		viper.Set("location.latitude", lat)
		viper.Set("location.longitude", lon)

		configDir := filepath.Clean(filepath.Join(os.Getenv("HOME"), ".config", "prayertimes"))
		if err := os.MkdirAll(configDir, 0o755); err != nil {
			fmt.Printf("could not create necessary dirs: %v", err)
			os.Exit(1)
		}

		configFile := filepath.Join(configDir, "config.yaml")
		if err := viper.SafeWriteConfigAs(configFile); err != nil {
			if os.IsExist(err) {
				viper.WatchConfig()
			} else {
				fmt.Printf("error writing config %v\n", err)
				os.Exit(1)
			}
		}
		fmt.Println("config updated successfully")
	}
}

func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
