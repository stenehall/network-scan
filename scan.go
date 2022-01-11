package main

import (
	"fmt"
	"github.com/Ullaakut/nmap/v2"
	"log"
)

func scan(scanner *nmap.Scanner, push Push) {
	fmt.Println("Starting a new scan")

	result, warnings, err := scanner.Run()
	if err != nil {
		fmt.Printf("%v", warnings)
		log.Fatalf("nmap scan failed: %v", err)
	}

	database := Database("known_hosts.db")

	for _, host := range result.Hosts {
		if len(host.Addresses) == 0 {
			continue
		}

		ip, hostname := extract(host)
		result := database.Check(ip, hostname)

		// If the IP is missing in the db it means we have a new host
		// We should save it and alert the user.
		if result.Error != nil {
			msg := fmt.Sprintf("Found a new IP on your network %v (%v)", ip, hostname)
			fmt.Println(msg)
			push.Message(msg)
		}

		// Existing IP. No need to send a push.
		if result.Error == nil {
			fmt.Printf("Exists in db %v (%v)", ip, hostname)
		}
	}

	fmt.Printf("Nmap done: %d hosts up scanned in %3f seconds\n", len(result.Hosts), result.Stats.Finished.Elapsed)
}

func extract(host nmap.Host) (string, string) {
	hostname := "unknown"
	if len(host.Hostnames) > 0 {
		hostname = host.Hostnames[0].Name
	}

	return host.Addresses[0].Addr, hostname
}
