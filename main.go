package main

import (
	"flag"
	"github.com/anliksim/cmd-wrapper/hub"
	"github.com/spf13/viper"
	"log"
	"os"
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
	gitUser := viper.GetString("git_user")
	repos := viper.GetStringSlice("repos")

	hubCmd := hub.Hub(verbose)
	log.Printf("Starting work in %s\n", workDir)
	for e := range repos {
		repo := repos[e]
		dir := workDir + repo
		log.Printf("Processing %s\n", dir)
		_, err := os.Stat(dir)
		if err != nil {
			cloneUrl := "git@github.com:" + gitUser + "/" + repo + ".git"
			log.Printf("Cloning %s\n", cloneUrl)
			hubCmd.WorkDir = workDir
			hubCmd.Exec("clone", cloneUrl)
		}
		hubCmd.WorkDir = dir
		hubCmd.Sync()
		hubCmd.Exec("gc")
	}

}
