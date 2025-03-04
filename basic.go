package main

import (
	"flag"
	"fmt"
	"github.com/erbekin/webgo/internal/server"
	"log"
	"os"
	"regexp"
)

func main() {
	parsedArgs := parseArgs()
	log.Printf("Parsed args: %#v\n", parsedArgs)
	if parsedArgs.run {
		// Start server
		var addr string = ""
		if parsedArgs.allowAllOrigins {
			fmt.Println("Allowing all origins")
			addr = "0.0.0.0"
		}

		portRegex := regexp.MustCompile(":[0-9]+")
		if parsedArgs.port == "" || !portRegex.MatchString(parsedArgs.port) {
			log.Println("Invalid port or you did not specify a port. Defaulting to :8080.")
			parsedArgs.port = ":8080"
		}
		err := server.Serve(addr + parsedArgs.port)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	flag.Usage()

}

// Parse args and return args type
func parseArgs() args {
	parsedArgs := args{}
	flag.BoolVar(&parsedArgs.run, "serve", false, "Starts server if set if this was not set, other options will not have any effect")
	flag.StringVar(&parsedArgs.port, "port", "", "Port to listen on when serving. It must be a valid port number. This can be used with -all flag.")
	flag.BoolVar(&parsedArgs.allowAllOrigins, "all", false, "Allow all origins of the server when serving. If you also want to specify a port number use this with -port")
	flag.Parse()
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Hello dev! Basic web server is stands in front of you!"+"\nHere is the usage:\n"+
			"Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	return parsedArgs
}

type args struct {
	run             bool
	port            string
	allowAllOrigins bool
}
