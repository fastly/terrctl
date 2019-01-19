package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jedisct1/dlog"
)

func deploy(root string, language string) (*TerrariumInstance, error) {
	var err error
	var instance *TerrariumInstance
	maxAttempts := Config().MaxDeployAttempts
	for attempt := uint(1); attempt <= maxAttempts; attempt++ {
		if attempt > 1 {
			dlog.Infof("Retrying (attempt %v out of %v)", attempt, maxAttempts)
		}
		instance, err = NewTerrariumInstance(root, language)
		if err != nil {
			dlog.Warn(err)
			continue
		}
		err = instance.WaitForDeployment()
		if err != nil {
			if err.Error() == "Timeout" {
				dlog.Warn(err)
				continue
			}
			return nil, err
		}
		dlog.Info("Instance is deployed")
		err = instance.WaitForHealth()
		if err != nil {
			dlog.Warn(err)
			continue
		}
		dlog.Notice("Instance is running and reachable over HTTPS")
		break
	}
	return instance, err
}

func usage() {
	outFd := flag.CommandLine.Output()
	fmt.Fprintln(outFd, "Usage: terrctl [options] <source code path>")
	fmt.Fprintln(outFd, "")
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	dlog.Init("terrctl", dlog.SeverityInfo, "USER")
	flag.Usage = usage
	if err := UpdateConfigFromFlags(); err != nil {
		dlog.Fatal(err)
	}
	root := flag.Arg(0)
	if len(root) <= 0 {
		usage()
	}
	instance, err := deploy(root, config.Language)
	if err != nil {
		dlog.Fatal(err)
	}
	url, err := instance.URL()
	if err != nil {
		dlog.Fatal(err)
	}
	dlog.Noticef("New instance deployed at [%v]", url)
}
