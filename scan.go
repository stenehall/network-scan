package main

import (
	"fmt"
	"github.com/Ullaakut/nmap/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

func scan(scanner *nmap.Scanner, push Push) {
	// @TODO Were do we want this? Outside the function?
	type Host struct {
		gorm.Model
		IP       string
		hostname string
	}

	// @TODO Were do we want this? Capitalized?
	var hosts []Host

	db, err := gorm.Open(sqlite.Open("known_hosts.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&hosts)

	fmt.Println("Starting a new scan")

	result, warnings, err := scanner.Run()
	if err != nil {
		fmt.Println("%#v", warnings)
		log.Fatalf("nmap scan failed: %v", err)
	}

	for _, host := range result.Hosts {
		if len(host.Addresses) == 0 {
			continue
		}

		ip := host.Addresses[0].Addr
		hostname := "unknown"
		if len(host.Hostnames) > 0 {
			hostname = host.Hostnames[0].Name
		}

		result := db.Where("IP = ?", ip).First(&hosts)

		// If the IP is missing in the db it means we have a new host
		// We should save it and alert the user.
		if result.Error != nil {
			msg := fmt.Sprintf("Found a new IP on your network %v (%v)", ip, hostname)
			fmt.Println(msg)
			push.Message(msg)

			// Save the new IP to DB.
			db.Create(&Host{IP: ip, hostname: hostname})
		}

		// Existing IP. No need to send a push.
		if result.Error == nil {
			fmt.Println("Exists in db %v (%v)", ip, hostname)
		}

		fmt.Printf("Host %q %v:\n", host.Addresses[0], host.Hostnames)
	}

	fmt.Printf("Nmap done: %d hosts up scanned in %3f seconds\n", len(result.Hosts), result.Stats.Finished.Elapsed)
}
