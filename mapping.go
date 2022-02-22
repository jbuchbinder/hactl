package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Define all application mappings

var (
	haMap      = map[string]haMapping{}
	clusterMap = map[string][]haServer{}
	serverMap  = map[string][]haMapping{}
)

type haConfig struct {
	Map      map[string]haMapping  `yaml:"map"`
	Clusters map[string][]haServer `yaml:"clusters"`
}

type haMapping struct {
	Cluster      string   `yaml:"cluster"`
	BackendSlugs []string `yaml:"backends,flow"`
	Servers      []string `yaml:"servers,flow"`
}

type haServer struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func parseConfig() {
	data, err := ioutil.ReadFile(*configFile)
	if err != nil {
		panic(err)
	}
	cfg := haConfig{}
	err = yaml.Unmarshal([]byte(data), &cfg)
	if err != nil {
		panic(err)
	}
	haMap = cfg.Map
	clusterMap = cfg.Clusters

	// Interpolate ServerMap
	for _, v := range haMap {
		for _, s := range v.Servers {
			if _, ok := serverMap[s]; !ok {
				serverMap[s] = make([]haMapping, 0)
			}
			serverMap[s] = append(serverMap[s], v)
		}
	}
}
