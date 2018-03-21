package main

import (
	"flag"
	"fmt"
	"os"
	"xqb"
	"github.com/jonnywang/go-kits/redis"
)

var optionConfigFile= flag.String("config", "./config.xml", "configure xml file")

func usage() {
	fmt.Printf("Version: %s\nUsage: %s [options]Options:", xqb.VERSION, os.Args[0])
	flag.PrintDefaults()
	os.Exit(0)
}

func main()  {
	flag.Usage = usage
	flag.Parse()

	if len(os.Args) < 2 {
		usage()
	}

	_, err := xqb.ParseXmlConfig(*optionConfigFile)
	if err != nil {
		redis.Logger.Print(err)
		os.Exit(1)
	}

	xqb.Run()
}
