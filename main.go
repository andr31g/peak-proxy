package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"net/http"
	"os"
	"peakproxy/common"
	"peakproxy/util"
	"strconv"
)

func main() {
	logLevel := &LogLevel
	localPort := &LocalPort
	targetURI := &TargetURI
	iterationCount := &IterationCount
	app := &cli.App{
		Name:  "peakproxy",
		Usage: "Peak Proxy",
		Commands: []*cli.Command{
			{
				Name:    "run",
				Aliases: []string{"r"},
				Usage:   "Runs the proxy",
				Flags: []cli.Flag{
					logLevel,
					localPort,
					targetURI,
					iterationCount,
				},
				Action: func(c *cli.Context) error {
					iterations := c.Uint(iterationCount.Name)
					target := c.String(targetURI.Name)
					level := c.String(logLevel.Name)
					port := c.Int(localPort.Name)
					run(port, target, iterations, level)
					return nil
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(localPort int, targetURI string, iterationCount uint, logLevel string) {
	logger, level := util.ConfigureLogger(logLevel)
	logger.Info("Log level is set to", "log-level", level.Level())

	if proxy, err := NewPeakProxy(targetURI, iterationCount, logger); err == nil {
		http.Handle("/", proxy)
	} else {
		common.LogFatalFailedToCreateInstanceErrorWithMessage(err, GetPeakProxyStructName(), logger)
	}
	if err := http.ListenAndServe(":"+strconv.Itoa(localPort), nil); err != nil {
		common.LogFatalFailedToStartHTTPServerError(err, logger)
	}
}
