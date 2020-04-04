package config

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

const testFile = "linux/config.yaml"
const noConfig = ""

func TestUnmarshalString_DefaultVerbose(t *testing.T) {
	config := UnmarshalString(noConfig)
	assert.False(t, config.Verbose)
}

func TestUnmarshalString_Verbose(t *testing.T) {
	config := UnmarshalString("verbose: true")
	assert.True(t, config.Verbose)
}

func TestUnmarshalString_DefaultRemove(t *testing.T) {
	config := UnmarshalString(noConfig)
	assert.False(t, config.Remove)
}

func TestUnmarshalString_Remove(t *testing.T) {
	config := UnmarshalString("remove: true")
	assert.True(t, config.Remove)
}

func TestUnmarshalString_DefaultClone(t *testing.T) {
	config := UnmarshalString(noConfig)
	assert.False(t, config.Clone)
}

func TestUnmarshalString_Clone(t *testing.T) {
	config := UnmarshalString("clone: true")
	assert.True(t, config.Clone)
}

func TestUnmarshalString_WorkDir(t *testing.T) {
	config := UnmarshalString("workdir: test")
	assert.Equal(t, config.WorkDir, "test")
}

func TestUnmarshalString_DefaultUrlPrefix(t *testing.T) {
	config := UnmarshalString(noConfig)
	assert.Equal(t, config.RemoteUrlPrefix, githubSSH)
}

func TestUnmarshalString_UrlPrefix(t *testing.T) {
	config := UnmarshalString("remote_url_prefix: test")
	assert.Equal(t, config.RemoteUrlPrefix, "test")
}

func TestUnmarshal_ContainsRepos(t *testing.T) {
	config := UnmarshalFile(testFile)
	repos := config.Repositories()
	assert.Contains(t, repos, Repository{"myuser", "repo-a"})
	assert.Contains(t, repos, Repository{"myuser", "repo-b"})
	assert.Contains(t, repos, Repository{"myuser", "repo-c"})
	assert.Contains(t, repos, Repository{"myorg", "work"})
}

func ExampleConfiguration() {
	config := UnmarshalFile(testFile)
	log.Printf("%v\n", config)
	log.Printf("%v\n", config.Repositories())
}
