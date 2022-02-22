package main

import (
	"fmt"
	"net"
	"time"
)

func haproxySetStatus(hahost string, haport int, backend, server string, enabled bool) error {
	cmd := "enable"
	if !enabled {
		cmd = "disable"
	}

	if *debug {
		fmt.Printf("%s:%d => %s server %s/%s\n", hahost, haport, cmd, backend, server)
		return nil
	}

	conn, err := net.DialTimeout("tcp",
		fmt.Sprintf("%s:%d", hahost, haport),
		time.Second*time.Duration(*timeout))
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = fmt.Fprintf(conn, "%s server %s/%s\n", cmd, backend, server)
	return err
}
