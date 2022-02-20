package nmapw

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Ullaakut/nmap/v2"
)

type NmapHost struct {
	Ip       string
	Hostname string
}

type Scanner struct {
	scanner *nmap.Scanner
}

func NewScanner(subNets []string) (*Scanner, error) {
	scanner, err := nmap.NewScanner(
		nmap.WithTargets(subNets...),
		nmap.WithUDPDiscovery(),
		nmap.WithSCTPDiscovery(),
		nmap.WithACKDiscovery(),
		nmap.WithSYNDiscovery(),
		nmap.WithICMPEchoDiscovery(),
		nmap.WithICMPNetMaskDiscovery(),
	)

	return &Scanner{scanner: scanner}, err
}

func track(msg string) (string, time.Time) {
	return msg, time.Now()
}

func duration(msg string, start time.Time) {
	log.Printf("%v: %v\n", msg, time.Since(start))
}

func (s *Scanner) Scan() []NmapHost {
	defer duration(track("nmap scan took"))

	fmt.Println("Starting a new scan")

	var hosts []NmapHost
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
	defer cancel()
	s.scanner.AddOptions(nmap.WithContext(ctx))

	progress := make(chan float32, 1)

	fmt.Printf("Progress: %v %%", 0)

	go func() {
		for p := range progress {
			fmt.Printf("\rProgress: %v %%", p)
		}
	}()
	result, warnings, err := s.scanner.RunWithProgress(progress)

	if err != nil {
		fmt.Printf("%v", warnings)
		log.Fatalf("nmapw scan failed: %v", err)
	}

	for _, host := range result.Hosts {
		if len(host.Addresses) == 0 {
			continue
		}

		ip, hostname := extract(host)
		hosts = append(hosts, NmapHost{
			Ip:       ip,
			Hostname: hostname,
		})

	}

	fmt.Printf("Nmap done: %d hosts up scanned in %3f seconds\n", len(result.Hosts), result.Stats.Finished.Elapsed)
	return hosts
}

func extract(host nmap.Host) (string, string) {
	hostname := "unknown"
	if len(host.Hostnames) > 0 {
		hostname = host.Hostnames[0].Name
	}

	return host.Addresses[0].Addr, hostname
}
