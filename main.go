/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"time"

	"github.com/lanrey-waju/prayer-times/cmd"
	"github.com/lanrey-waju/prayer-times/internal/config"
)

func main() {
	defer config.TimeTrack(time.Now(), "main")
	cmd.Execute()
}
