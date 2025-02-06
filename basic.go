package main

import (
	"flag"
	"fmt"
	"github.com/Emre-Erbekin/webgo/internal/server"
	"log"
)

func main() {
	parsedArgs := parseArgs()
	if parsedArgs.run != true {
		fmt.Println("Hello! All things are ready.")
		flag.Usage()
		return
	}
	if parsedArgs.port == "" {
		fmt.Println("You didn't specify a port, defaulting to :8080.")
		parsedArgs.port = ":8080"
	}
	// Start server
	err := server.Serve(parsedArgs.port)
	if err != nil {
		log.Fatal(err)
	}

}

// Parse args and return args type
func parseArgs() args {
	parsedArgs := args{}
	flag.BoolVar(&parsedArgs.run, "serve", false, "Starts server.")
	flag.StringVar(&parsedArgs.port, "port", "", "Port to listen on. It must be a valid port number.")
	flag.Parse()
	return parsedArgs
}

type args struct {
	run  bool
	port string
}
