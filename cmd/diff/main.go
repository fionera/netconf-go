package main

import (
	"flag"
	"fmt"
	"github.com/r3labs/diff/v2"
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

	wb, err := os.ReadFile(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	wantedCfg, err := vyos.Parse(wb)
	if err != nil {
		log.Fatal(err)
	}

	b, err := vyos.FetchConfig(*user, *addr)
	if err != nil {
		log.Fatal(err)
	}

	currentCfg, err := vyos.Parse(b)
	if err != nil {
		log.Fatal(err)
	}

	changelog, err := diff.Diff(currentCfg, wantedCfg)
	if err != nil {
		log.Fatal(err)
	}

	for _, s := range vyos.CommandsFromDiff(changelog) {
		fmt.Println(s)
	}
}
