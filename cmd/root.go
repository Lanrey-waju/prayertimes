/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/lanrey-waju/prayertimes/internal/cache"
	"github.com/lanrey-waju/prayertimes/internal/config"
	"github.com/lanrey-waju/prayertimes/internal/timings"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var version = "dev"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "prayertimes",
	Short: "A cli app to get prayer times",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		// initialize config and get an instance of sql.DB
		var prayerTimes *timings.PrayerTimes
		var err error

		osProvider := cache.NewDefaultOSProvider()

		db, err := config.InitConfig(osProvider)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// check if city is set in config
		city := viper.GetString("location.city")

		// get an instance of Queries
		queries := cache.New(db)

		// if city is not set, get location info
		if city == "" {
			config.EnsureConfig(timings.GetLocationParams)
			// assign city
			city = viper.GetString("location.city")
		}

		if cmd.Flag("refresh").Changed {
			// if refresh flag is set, ignore cache and refresh prayer times
			prayerTimes, err = timings.GetPrayerTimes(queries, city)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("Prayer times refreshed successfully")
			fmt.Println(prayerTimes)
			os.Exit(0)
		}

		if cmd.Flag("version").Changed {
			fmt.Printf("prayertimes %s\n", version)
			os.Exit(0)
		}

		prayerTimes, err = timings.RetrievePrayerTimes(queries, city)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(prayerTimes)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.prayer-times.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().BoolP("version", "v", false, "Show version")
	rootCmd.Flags().BoolP("refresh", "r", false, "Refresh prayer times from API")
}
