package main

import (
	"flag"
	"fmt"
	"strings"
)

var (
	ConfigFile   = flag.String("config", "/etc/hactl.yml", "Config file location")
	Debug        = flag.Bool("debug", false, "Debug (don't execute commands)")
	ListServices = flag.Bool("list", false, "List services")
	App          = flag.String("app", "", "Application")
	Server       = flag.String("server", "", "Server")
	Enable       = flag.Bool("enable", false, "Enable service")
	Disable      = flag.Bool("disable", false, "Disable service")
	EnableAll    = flag.Bool("enableall", false, "Enable service")
	DisableAll   = flag.Bool("disableall", false, "Disable service")
	Failover     = flag.Bool("failover", false, "Disable all services")
	Recover      = flag.Bool("recover", false, "Enable all services")
)

func main() {
	flag.Parse()

	parseConfig()

	if *ListServices {
		listServices()
		return
	}

	if !*Enable && !*Disable && !*EnableAll && !*DisableAll && !*Failover && !*Recover {
		usage()
		return
	}

	// Handle single service enable/disable
	if *Enable || *Disable {
		m, ok := HaMap[*App]
		if !ok {
			fmt.Printf("Service '%s' not found.\n", *App)
			return
		}

		c := ClusterMap[m.Cluster]
		for _, h := range c {
			for _, b := range m.BackendSlugs {
				action := "enable"
				if *Disable {
					action = "disable"
				}
				fmt.Printf("%s:%d[%s] %s %s\n", h.Host, h.Port, b, action, *Server)
				err := HaproxySetStatus(h.Host, h.Port, b, *Server, *Enable)
				if err != nil {
					fmt.Printf("%v\n", err)
				}
			}
		}
	}

	// Handle enable/disable for entire server
	if *EnableAll || *DisableAll {
		serverAction(*Server, *EnableAll)
		return
	}

	if *Failover || *Recover {
		for s, _ := range ServerMap {
			fmt.Printf("Processing server %s\n", s)
			serverAction(s, *Recover)
		}
	}
}

func serverAction(s string, enable bool) {
	_, ok := ServerMap[s]
	if !ok {
		fmt.Printf("Server '%s' not found.\n", s)
		return
	}

	for _, v := range ServerMap[s] {
		c := ClusterMap[v.Cluster]
		for _, h := range c {
			for _, b := range v.BackendSlugs {
				action := "enable"
				if *DisableAll {
					action = "disable"
				}
				fmt.Printf("%s:%d[%s] %s %s\n", h.Host, h.Port, b, action, s)
				err := HaproxySetStatus(h.Host, h.Port, b, s, enable)
				if err != nil {
					fmt.Printf("%v\n", err)
				}
			}
		}
	}
}

func listServices() {
	fmt.Printf("Supported services:\n")
	for k, v := range HaMap {
		fmt.Printf("\t%s [%s]\n", k, strings.Join(v.Servers, " "))
	}
}

func usage() {
	flag.PrintDefaults()
}
