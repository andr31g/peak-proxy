package main

import "github.com/urfave/cli/v2"

var (
	LocalPort = cli.IntFlag{
		Name:     "local-port",
		Aliases:  []string{"p"},
		Usage:    "Local port the proxy will listen on",
		Required: true,
	}
	TargetURI = cli.StringFlag{
		Name:     "target-uri",
		Aliases:  []string{"t"},
		Usage:    "URI of the target Prometheus host",
		Required: true,
	}
	IterationCount = cli.UintFlag{
		Name:     "iterations",
		Value:    1,
		Aliases:  []string{"i"},
		Usage:    "Number of iterations to run the peak detector",
		Required: false,
	}
	LogLevel = cli.StringFlag{
		Name:     "log-level",
		Value:    "info",
		Aliases:  []string{"l"},
		Usage:    "Logging level. Case insensitive. Must be one of: DEBUG, INFO, WARN, ERROR",
		Required: false,
	}
)
