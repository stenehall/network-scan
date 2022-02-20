package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"network-scan/internal/database"
	"network-scan/internal/nmapw"
	"network-scan/internal/pushover"
)

type SubNets []string

func (i *SubNets) String() string {
	return "wrong subnet format"
}

func (i *SubNets) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	var subNets SubNets

	flag.Var(&subNets, "subnet", "IP subnet to scan")
	pushoverToken := flag.String("pushoverToken", os.Getenv("PUSHOVER_TOKEN"), "Pushover access token")
	pushoverRecipient := flag.String("pushoverRecipient", os.Getenv("PUSHOVER_RECIPIENT"), "Pushover recipient token")
	flag.Parse()

	push := pushover.PushOver(*pushoverToken, *pushoverRecipient)
	if *pushoverToken == "" || *pushoverRecipient == "" {
		fmt.Println("No PushOver tokens provided, only outputting to log")
	}

	err := push.Validate()
	if err != nil {
		log.Fatalf("The PushOver tokens provided seems invalid, %v\n", err)
	}

	if len(subNets) == 0 {
		subNets = strings.Fields(os.Getenv("SUBNETS"))
	}
	if len(subNets) == 0 {
		log.Fatalf("No subnet provided. We have nothing to do\n")
	}

	// Print out the subnets selected for scan
	fmt.Println("Scanning the following IP subnets")
	for _, subNet := range subNets {
		fmt.Println("- ", subNet)
	}

	if err != nil {
		log.Fatalf("unable to create nmapw scanner: %v", err)
	}

	db := database.Database("known_hosts.db")

	s, err := nmapw.NewScanner(subNets)
	if err != nil {
		log.Fatalf("couldn't instanciate nmap")
	}

	scan(s, db)

	// Wait 2 minutes then rescan
	go func() {
		c := time.Tick(1200 * time.Second)
		for range c {
			scan(s, db)
		}
	}()

	// @TODO: Understand this black magic
	// A select statement blocks until at least one of it’s cases can proceed. With zero cases this will never happen.
	select {}
}

func scan(scanner *nmapw.Scanner, database database.DB) {
	// Run initial scan
	hosts := scanner.Scan()

	for _, host := range hosts {
		result := database.AddIfNotExist(host.Ip, host.Hostname)

		// If the IP is missing in the db it means we have a new host
		// We should save it and alert the user.
		if result.Error != nil {
			msg := fmt.Sprintf("- New\t%v\t(%v)", host.Ip, host.Hostname)
			fmt.Println(msg)
			// push.Message(msg)
		}

		// Existing IP. No need to send a push.
		if result.Error == nil {
			fmt.Printf("- Existing\t%v\t(%v)\n", host.Ip, host.Hostname)
		}
	}
}
