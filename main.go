package main

import (
	"flag"
	"github.com/anliksim/cmd-wrapper/hub"
	"github.com/spf13/viper"
	"log"
	"os"
)

func main() {

	configFilePtr := flag.String("configFile", "config", "Name of the YAML file without the extension")
	configPathPtr := flag.String("configPath", ".", "Directory where config.yml is located")
	boolPtr := flag.Bool("d", false, "Run as daemon with schedule sync")
	flag.Parse()
	if *boolPtr {
		log.Println("Running in daemon mode")
		log.Fatal("Unsupported operation")
	}

	viper.SetConfigName(*configFilePtr)
	viper.AddConfigPath(*configPathPtr)
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file, %s", err)
	}

	clone := viper.GetBool("clone")
	verbose := viper.GetBool("verbose")
	workDir := viper.GetString("workdir")

	orgs := viper.GetStringSlice("orgs")
	repoMap := viper.GetStringMapStringSlice("repos")

	hubCmd := hub.Hub(verbose)
	log.Printf("Starting work in %s\n", workDir)
	for e := range orgs {
		org := orgs[e]
		log.Printf("Processing organisation %s\n", org)
		repos := repoMap[org]

		for e := range repos {
			repo := repos[e]
			dir := workDir + repo
			_, err := os.Stat(dir)
			if err != nil {
				if clone {
					cloneRepo(org, repo, workDir, hubCmd)
					syncAndClean(dir, hubCmd)
				} else {
					log.Printf("Ignoring %s\n", dir)
				}
			} else {
				syncAndClean(dir, hubCmd)
			}
		}
	}
}

func syncAndClean(dir string, hubCmd *hub.Cmd) {
	log.Printf("Processing %s\n", dir)
	hubCmd.WorkDir = dir
	hubCmd.Sync()
	hubCmd.Exec("gc")
}

func cloneRepo(org string, repo string, workDir string, hubCmd *hub.Cmd) {
	cloneUrl := "git@github.com:" + org + "/" + repo + ".git"
	log.Printf("Cloning %s\n", cloneUrl)
	hubCmd.WorkDir = workDir
	hubCmd.Exec("clone", cloneUrl)
}
