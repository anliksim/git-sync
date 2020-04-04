package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

const githubSSH = "git@github.com:"

type Configuration struct {
	// work directory where the listed repositories are located
	WorkDir string `yaml:"workdir"`

	// the remote url prefix lets you change the protocol or use a different git server
	// - for SSL use: 'git@github.com:' (default)
	// - for HTTPS use: 'https://github.com/'
	RemoteUrlPrefix string `yaml:"remote_url_prefix"`

	// clone projects if they do not exist
	Clone bool `yaml:"clone"`

	// remove projects if they are not listed (i.e. moves them into a folder 'removed')
	Remove bool `yaml:"remove"`

	// if command output should be redirected to stdout
	Verbose bool `yaml:"verbose"`

	// list of repositories to be processes if the org is listed
	Repos map[string][]string `yaml:"repos"`
}

type Repository struct {
	//
	Organisation string
	//
	Name string
}

func (c *Configuration) Repositories() []Repository {
	var repositories []Repository
	for k, v := range c.Repos {
		for _, vv := range v {
			repositories = append(repositories, Repository{k, vv})
		}
	}
	return repositories
}

func UnmarshalString(yamlString string) Configuration {
	return Unmarshal([]byte(yamlString))
}

func UnmarshalFile(filePath string) Configuration {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return Unmarshal(content)
}

func Unmarshal(yamlContent []byte) Configuration {
	config := Configuration{
		RemoteUrlPrefix: githubSSH,
		Clone:           false,
		Remove:          false,
		Verbose:         false,
	}
	err := yaml.Unmarshal(yamlContent, &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}
