package main

import (
	"flag"
	"github.com/anliksim/cmd-wrapper/hub"
	"github.com/spf13/viper"
	"log"
)

func main() {

	boolPtr := flag.Bool("d", false, "Run as daemon with schedule sync")
	flag.Parse()
	if *boolPtr {
		log.Println("Running in daemon mode")
		log.Fatal("Unsupported operation")
	}

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file, %s", err)
	}

	verbose := viper.GetBool("verbose")
	workDir := viper.GetString("workdir")
	repos := viper.GetStringSlice("repos")

	hubCmd := hub.Hub(verbose)
	log.Printf("Starting work in %s\n", workDir)
	for e := range repos {
		dir := workDir + repos[e]
		hubCmd.WorkDir = dir
		log.Printf("Processing %s\n", dir)
		hubCmd.Sync()
		hubCmd.Exec("gc")
	}
}
