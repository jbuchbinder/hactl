# HACTL

[![Build Status](https://secure.travis-ci.org/jbuchbinder/hactl.png)](http://travis-ci.org/jbuchbinder/hactl)
[![Report Card](https://goreportcard.com/badge/github.com/jbuchbinder/hactl)](https://goreportcard.com/report/github.com/jbuchbinder/hactl)

haproxy control via stats/admin socket for 1.4/1.5+

## Building

        go get -d ; go build

## Configuration

Enable admin sockets in your haproxy configuration in the ``global`` section:

	nbproc 2
	stats socket 10.0.1.3:6101 uid 99 gid 99 level admin process 1
	stats socket 10.0.1.3:6102 uid 99 gid 99 level admin process 2

*This shows support for binding on 10.0.1.3 for specific defaults -- your
configuration will most likely vary.*

Edit the ``hactl.yml`` configuration file to support your configuration of
haproxy server sockets (including support for nbproc > 1, which requires
multiple sockets). **hactl** uses the concept of "clusters" of haproxy
servers, for those who use N+1 layouts.

From there, the ``map`` yaml element contains your list of web applications,
along with one or more backend, one or more server (mapping to the ``server``
directive in haproxy's ``backend`` configuration), and an arbitrary cluster
name (which maps to the aforementioned ``cluster`` yaml map).

