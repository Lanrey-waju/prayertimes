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

// refreshCmd represents the refresh command
var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "ignores cache and refreshes prayer times",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := config.InitConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// get an instance of Queries
		queries := cache.New(db)

		// if city is not set, get location info
		config.EnsureConfig(timings.GetLocationParams)
		// assign city
		city := viper.GetString("location.city")

		prayertimes, err := timings.GetPrayerTimes(queries, city)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Prayer times refreshed successfully")
		fmt.Println(prayertimes)
	},
}

func init() {
	rootCmd.AddCommand(refreshCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// refreshCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// refreshCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
