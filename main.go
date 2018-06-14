package main

import (
	"flag"
	"fmt"
	"strings"
)

var (
	configFile   = flag.String("config", "/etc/hactl.yml", "Config file location")
	debug        = flag.Bool("debug", false, "Debug (don't execute commands)")
	listServices = flag.Bool("list", false, "List services")
	app          = flag.String("app", "", "Application")
	server       = flag.String("server", "", "Server")
	enable       = flag.Bool("enable", false, "Enable service")
	disable      = flag.Bool("disable", false, "Disable service")
	enableAll    = flag.Bool("enableall", false, "Enable service")
	disableAll   = flag.Bool("disableall", false, "Disable service")
	failoverAll  = flag.Bool("failover", false, "Disable all services")
	recoverAll   = flag.Bool("recover", false, "Enable all services")
)

func main() {
	flag.Parse()

	parseConfig()

	if *listServices {
		listServicesAction()
		return
	}

	if !*enable && !*disable && !*enableAll && !*disableAll && !*failoverAll && !*recoverAll {
		usage()
		return
	}

	// Handle single service enable/disable
	if *enable || *disable {
		m, ok := HaMap[*app]
		if !ok {
			fmt.Printf("Service '%s' not found.\n", *app)
			return
		}

		c := ClusterMap[m.Cluster]
		for _, h := range c {
			for _, b := range m.BackendSlugs {
				action := "enable"
				if *disable {
					action = "disable"
				}
				fmt.Printf("%s:%d[%s] %s %s\n", h.Host, h.Port, b, action, *server)
				err := haproxySetStatus(h.Host, h.Port, b, *server, *enable)
				if err != nil {
					fmt.Printf("%v\n", err)
				}
			}
		}
	}

	// Handle enable/disable for entire server
	if *enableAll || *disableAll {
		serverAction(*server, *enableAll)
		return
	}

	if *failoverAll || *recoverAll {
		for s := range ServerMap {
			fmt.Printf("Processing server %s\n", s)
			serverAction(s, *recoverAll)
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
				if *disableAll {
					action = "disable"
				}
				fmt.Printf("%s:%d[%s] %s %s\n", h.Host, h.Port, b, action, s)
				err := haproxySetStatus(h.Host, h.Port, b, s, enable)
				if err != nil {
					fmt.Printf("%v\n", err)
				}
			}
		}
	}
}

func listServicesAction() {
	fmt.Printf("Supported services:\n")
	for k, v := range HaMap {
		fmt.Printf("\t%s [%s]\n", k, strings.Join(v.Servers, " "))
	}
}

func usage() {
	flag.PrintDefaults()
}
