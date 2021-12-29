package main

import (
	"flag"
	"log"
	"netconf-go/vyos"
	"os"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:22", "The address of the system to connect to")
	user := flag.String("user", "root", "The username for the system")
	flag.Parse()

	if flag.NArg() != 1 {
		log.Fatal("Please provide the config file to diff against")
	}

	b, err := vyos.FetchConfig(*user, *addr)
	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(flag.Arg(0), b, 0755); err != nil {
		log.Fatal(err)
	}
}
