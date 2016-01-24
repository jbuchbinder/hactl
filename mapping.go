package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Define all application mappings

var (
	HaMap      = map[string]HaMapping{}
	ClusterMap = map[string][]HaServer{}
	ServerMap  = map[string][]HaMapping{}
)

type HaConfig struct {
	Map      map[string]HaMapping  `yaml:"map"`
	Clusters map[string][]HaServer `yaml:"clusters"`
}

type HaMapping struct {
	Cluster      string   `yaml:"cluster"`
	BackendSlugs []string `yaml:"backends,flow"`
	Servers      []string `yaml:"servers,flow"`
}

type HaServer struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func parseConfig() {
	data, err := ioutil.ReadFile(*ConfigFile)
	if err != nil {
		panic(err)
	}
	cfg := HaConfig{}
	err = yaml.Unmarshal([]byte(data), &cfg)
	if err != nil {
		panic(err)
	}
	HaMap = cfg.Map
	ClusterMap = cfg.Clusters

	// Interpolate ServerMap
	for _, v := range HaMap {
		for _, s := range v.Servers {
			if _, ok := ServerMap[s]; !ok {
				ServerMap[s] = make([]HaMapping, 0)
			}
			ServerMap[s] = append(ServerMap[s], v)
		}
	}
}
