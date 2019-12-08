package main

import (
	"flag"
	"github.com/anliksim/cmd-wrapper/hub"
	"github.com/spf13/viper"
	"io/ioutil"
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
	remove := viper.GetBool("remove")
	verbose := viper.GetBool("verbose")
	workDir := viper.GetString("workdir")
	remoteUrlPrefix := viper.GetString("remote_url_prefix")

	orgs := viper.GetStringSlice("orgs")
	repoMap := viper.GetStringMapStringSlice("repos")

	dirs := listAllDirectories(workDir)

	hubCmd := hub.Hub(verbose)
	log.Printf("Starting work in %s\n", workDir)
	for e := range orgs {
		org := orgs[e]
		log.Printf("Processing organisation %s\n", org)
		repos := repoMap[org]

		for e := range repos {
			repoDir := repos[e]
			delete(dirs, repoDir)
			dirPath := workDir + repoDir
			_, err := os.Stat(dirPath)
			if err != nil {
				if clone {
					cloneRepo(remoteUrlPrefix, org, repoDir, workDir, hubCmd)
					syncAndClean(dirPath, hubCmd)
				} else {
					log.Printf("Ignoring %s\n", dirPath)
				}
			} else {
				syncAndClean(dirPath, hubCmd)
			}
		}
	}

	if remove {
		log.Printf("Cleaning up %s\n", workDir)
		removedName := "removed"
		removedPath := workDir + removedName + "/"
		delete(dirs, removedName)

		if isNotEmpty(dirs) {
			createIfNotExists(removedPath)
			for dir, _ := range dirs {
				moveToRemovedDir(dir, workDir, removedPath)
			}
		}
	}
}

func isNotEmpty(dirs map[string]int) bool {
	return len(dirs) > 0
}

func moveToRemovedDir(dir string, oldPath string, newPath string) {
	log.Printf("Removing %s\n", dir)
	err := os.Rename(oldPath+dir, newPath+dir)
	if err != nil {
		log.Fatal(err)
	}
}

func createIfNotExists(dirPath string) {
	_, err := os.Stat(dirPath)
	if err != nil {
		err := os.Mkdir(dirPath, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func listAllDirectories(workDir string) map[string]int {
	files, err := ioutil.ReadDir(workDir)
	if err != nil {
		log.Fatal(err)
	}
	dirs := map[string]int{}
	for _, f := range files {
		if f.IsDir() {
			dirs[f.Name()] = 0
		}
	}
	return dirs
}

func syncAndClean(dir string, hubCmd *hub.Cmd) {
	log.Printf("Processing %s\n", dir)
	hubCmd.WorkDir = dir
	hubCmd.Sync()
	hubCmd.Exec("gc")
}

func cloneRepo(remoteUrlPrefix string, org string, repo string, workDir string, hubCmd *hub.Cmd) {
	cloneUrl := remoteUrlPrefix + org + "/" + repo + ".git"
	log.Printf("Cloning %s\n", cloneUrl)
	hubCmd.WorkDir = workDir
	hubCmd.Exec("clone", cloneUrl)
}
