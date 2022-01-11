package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Ullaakut/nmap/v2"
	"log"
	"time"
)

type arraySubNets []string

func (i *arraySubNets) String() string {
	return "wrong subnet format"
}

func (i *arraySubNets) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	var subNets arraySubNets

	flag.Var(&subNets, "subnet", "IP subnet to scan")
	pushoverToken := flag.String("pushoverToken", "", "Pushover access token")
	pushoverRecipient := flag.String("pushoverRecipient", "", "Pushover recipient token")
	flag.Parse()

	if *pushoverToken == "" {
		log.Fatalf("no pushover access token provided")
	}

	if *pushoverRecipient == "" {
		log.Fatalf("no pushover recipient token provided")
	}

	push := PushOver(*pushoverToken, *pushoverRecipient)

	// Print out the subnets selected for scaning
	fmt.Println("Scanning the following IP subnets")
	for _, subNet := range subNets {
		fmt.Println("- ", subNet)
	}

	// Create a context with a 5 minut timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	scanner, err := nmap.NewScanner(
		// Fancy unpack of an []string
		nmap.WithTargets(subNets...),
		nmap.WithFastMode(),
		nmap.WithContext(ctx),
	)
	if err != nil {
		log.Fatalf("unable to create nmap scanner: %v", err)
	}

	// Run initial scan
	scan(scanner, push)

	// Wait 2 minutes then rescan
	go func() {
		c := time.Tick(120 * time.Second)
		for range c {
			scan(scanner, push)
		}
	}()

	// @TODO: Understand this black magic
	// A select statement blocks until at least one of it’s cases can proceed. With zero cases this will never happen.
	select {}

	fmt.Println("Stopping?")
}
