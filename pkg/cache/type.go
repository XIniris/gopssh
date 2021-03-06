package cache

import (
	"fmt"
	"path/filepath"

	"gopssh/pkg/base64"
	"gopssh/pkg/file"
	"gopssh/pkg/rand"
	"gopssh/pkg/ssh"
)

/*
Cache logic:

          ┌───────────────┐
          │  if use cache ├─────────────────────┐
          └───────┬───────┘     not use cache   │
                  │ use cache                   │
                  ▼                             │
     ┌──────────────────────────┐               │
     │ if config cache in index ├───────────────┤
     └────────────┬─────────────┘   cache       │
                  │ cache found     not found   │
                  │ in index                    │
                  ▼                             │
    ┌────────────────────────────┐              │
    │ if config mod time changed ├──────────────┤
    └─────────────┬──────────────┘   mod time   │
                  │  mod time        changed    │
                  │  not changed                │
                  ▼                             ▼
         ┌──────────────────┐          ┌─────────────────┐
         │ use config cache │          │ analysis config │
         └──────────────────┘          └─────────────────┘

ASCII Workflow build with: https://asciiflow.com/
*/

const (
	cacheDirName    = ".gopssh/.cache"
	cacheFilePrefix = "cache"
	cacheIndexName  = "cache_index.yaml"

	randStringLength = 5
)

var (
	cacheDir       = filepath.Join(file.MustGetUserHome(), cacheDirName)
	cacheIndexPath = filepath.Join(cacheDir, cacheIndexName)
)

type CacheIndex []*Cache

type Cache struct {
	ConfigPath    string    `json:"config_path" yaml:"config_path"`
	InstancesPath string    `json:"instances_path" yaml:"instances_path"`
	ModTime       int64     `json:"mod_time" yaml:"mod_time"`
	Instances     Instances `json:"-" yaml:"-"`
}

type Instances []*Instance

type Instance struct {
	SSH    *ssh.SSH          `json:"ssh" yaml:"ssh"`
	Labels map[string]string `json:"labels" yaml:"labels"`
}

func GetCacheDir() string {
	return cacheDir
}

func GetCacheIndexPath() string {
	return cacheIndexPath
}

func (i CacheIndex) GetCacheByConfigPath(configPath string) *Cache {
	for _, cache := range i {
		if cache.ConfigPath == configPath {
			return cache
		}
	}

	return nil
}

func GenerateCacheFilePathByConfig(configPath string) string {
	return filepath.Join(GetCacheDir(), GenerateCacheFileNameByConfig(configPath))
}

func GenerateCacheFileNameByConfig(configPath string) string {
	// Suffix .json or .yaml, length 5
	cfgName := filepath.Base(configPath)
	// Sample: cache_j9sxt_cfgName.yaml
	return fmt.Sprintf("%s_%s_%s.yaml", cacheFilePrefix, cfgName[:len(cfgName)-5], rand.String(randStringLength))
}

func (i *Instance) HasLabels(labels map[string]string) bool {
	for k, v := range labels {
		if i.Labels[k] != v {
			return false
		}
	}

	return true
}

func (i *Instance) String() string {
	return fmt.Sprintf("Ip: %s, Port: %d, Username: %s, Password: %s, Labels: %v", i.SSH.Address.Ip, i.SSH.Address.Port, i.SSH.Username, i.SSH.Password, PrettyMap(i.Labels))
}

func (i *Instance) StringDecodePassword() (string, error) {
	pwd := i.SSH.Password
	decPwd, err := base64.Decode(pwd)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Ip: %s, Port: %d, Username: %s, Password: %s, Labels: %#v", i.SSH.Address.Ip, i.SSH.Address.Port, i.SSH.Username, decPwd, PrettyMap(i.Labels)), nil
}

func PrettyMap(input map[string]string) string {
	output := "{ "
	for k, v := range input {
		output += fmt.Sprintf("%s: %s, ", k, v)
	}

	return output[:len(output)-2] + " }"
}