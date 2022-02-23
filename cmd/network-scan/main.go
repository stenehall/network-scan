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

type subNets []string

func (i *subNets) String() string {
	return "wrong subnet format"
}

func (i *subNets) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	dbKnownHosts := database.Database("known_hosts.db")
	if len(os.Args) < 1 {
		log.Fatalf("Missing command")
	}

	command := os.Args[1]
	switch command {
	case "scan":
		mainScan(os.Args[2:], dbKnownHosts)
	case "list":
		mainList(dbKnownHosts)
	}
}

func mainList(db database.DB) {
	hosts := db.GetAll()
	for _, host := range hosts {
		log.Printf("%s\t%s\n", host.IP, host.Hostname)
	}
}

func mainScan(args []string, dbKnownHosts database.DB) {
	var subNets subNets

	fs := flag.NewFlagSet("mainScan", flag.ContinueOnError)
	fs.Var(&subNets, "subnet", "IP subnet to scan")
	pushoverToken := fs.String("pushoverToken", os.Getenv("PUSHOVER_TOKEN"), "Pushover access token")
	pushoverRecipient := fs.String("pushoverRecipient", os.Getenv("PUSHOVER_RECIPIENT"), "Pushover recipient token")
	_ = fs.Parse(args)

	push, err := pushover.PushOver(*pushoverToken, *pushoverRecipient)
	if *pushoverToken == "" || *pushoverRecipient == "" {
		log.Println("No PushOver tokens provided, only outputting to log")
	}
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
	log.Println("Scanning subnets")
	for _, subNet := range subNets {
		log.Println("- ", subNet)
	}

	scanner, err := nmapw.NewScanner(subNets)
	if err != nil {
		log.Fatalf("couldn't create a new  nmap instance")
	}

	err = scan(scanner, dbKnownHosts, push)
	if err != nil {
		log.Fatalf("initial scan %v", err)
	}

	// Wait 2 minutes then rescan
	go func() {
		for range time.Tick(1200 * time.Second) {
			err := scan(scanner, dbKnownHosts, push)
			log.Println(err)
		}
	}()

	// @TODO: Understand this black magic
	// A select statement blocks until at least one of it’s cases can proceed. With zero cases this will never happen.
	select {}
}

func scan(scanner *nmapw.Scanner, database database.DB, push pushover.Push) error {
	// Run initial scan
	hosts, err := scanner.Scan()

	if err != nil {
		return fmt.Errorf("scan error %w", err)
	}

	for _, host := range hosts {
		result := database.AddIfNotExist(host.IP, host.Hostname)

		// If the IP is missing in the db it means we have a new host
		// We should save it and alert the user.
		if result.Error != nil {
			msg := fmt.Sprintf("- New\t%v\t(%v)", host.IP, host.Hostname)
			log.Println(msg)
			push.Message(msg)
		}

		// Existing IP. No need to send a push.
		if result.Error == nil {
			log.Printf("- Existing\t%v\t(%v)\n", host.IP, host.Hostname)
		}
	}
	return nil
}
