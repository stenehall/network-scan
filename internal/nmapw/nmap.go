package nmapw

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Ullaakut/nmap/v2"
)

// NmapHost is the result object.
type NmapHost struct {
	IP       string
	Hostname string
}

// Scanner model.
type Scanner struct {
	scanner *nmap.Scanner
}

// NewScanner returns a new instance of the scanner.
func NewScanner(subNets []string) (*Scanner, error) {
	scanner, err := nmap.NewScanner(
		nmap.WithTargets(subNets...),
		nmap.WithUDPDiscovery(),
		nmap.WithSCTPDiscovery(),
		nmap.WithACKDiscovery(),
		nmap.WithSYNDiscovery(),
		// nmap.WithICMPEchoDiscovery(),
		nmap.WithICMPNetMaskDiscovery(),
		nmap.WithIPProtocolPingDiscovery(),
		nmap.WithFastMode(),
	)

	if err != nil {
		err = fmt.Errorf("newScanner %w", err)
	}

	return &Scanner{scanner: scanner}, err
}

func track(msg string) (string, time.Time) {
	return msg, time.Now()
}

func duration(msg string, start time.Time) {
	log.Printf("%v: %v\n", msg, time.Since(start))
}

// Scan runs a new scan.
func (s *Scanner) Scan() ([]NmapHost, error) {
	defer duration(track("nmap scan took"))

	log.Println("Starting a new scan")

	hosts := make([]NmapHost, 0)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
	defer cancel()
	s.scanner.AddOptions(nmap.WithContext(ctx))

	progress := make(chan float32, 1)

	log.Printf("Progress: %v %%", 0)

	go func() {
		for p := range progress {
			log.Printf("\rProgress: %v %%", p)
		}
	}()
	result, warnings, err := s.scanner.RunWithProgress(progress)

	if err != nil {
		log.Printf("%v", warnings)
		return nil, fmt.Errorf("nmapw scan failed: %w", err)
	}

	log.Println("----------------------------------")
	log.Println(result)
	log.Println("----------------------------------")

	for _, host := range result.Hosts {
		if len(host.Addresses) == 0 {
			continue
		}
		log.Println("----------------------------------")
		log.Println(host.Hostnames)
		log.Println(host.Addresses)
		log.Println(host.Comment)
		log.Println(host.Ports)
		log.Println(host.Status)

		ip, hostname := extract(host)
		hosts = append(hosts, NmapHost{
			IP:       ip,
			Hostname: hostname,
		})
	}

	log.Printf("Nmap done: %d hosts up scanned in %3f seconds\n", len(result.Hosts), result.Stats.Finished.Elapsed)
	return hosts, nil
}

func extract(host nmap.Host) (string, string) {
	hostname := "unknown"
	if len(host.Hostnames) > 0 {
		hostname = host.Hostnames[0].Name
	}

	return host.Addresses[0].Addr, hostname
}
