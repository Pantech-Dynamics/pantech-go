// Command dns-records hosts a zone on Pantech DNS (or uses the one you
// already have), points a name at an address, waits until Pantech DNS serves
// it, and prints the zone as a BIND zone file.
//
//	PANTECH_API_KEY=PAN_… go run ./examples/dns-records -zone example.com -name www -ip 203.0.113.10
//
// A new zone answers once the domain's registrar delegates to the two name
// servers it prints. Records can be added before that; they stay pending
// until the zone is active.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	pantech "github.com/Pantech-Dynamics/pantech-go"
)

func main() {
	zoneName := flag.String("zone", "", "the zone's name, such as example.com")
	name := flag.String("name", "www", "the record's name, relative to the zone (@ for the apex)")
	ip := flag.String("ip", "", "the IPv4 address to point the name at")
	ttl := flag.Int("ttl", 300, "the record's TTL in seconds, 60 to 86400")
	flag.Parse()
	if *zoneName == "" || *ip == "" {
		log.Fatal("-zone and -ip are required")
	}

	c, err := pantech.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	zone, err := findOrCreateZone(ctx, c, *zoneName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("zone %s (%s) is %s\n", zone.Name, zone.ID, zone.Status)
	if zone.Status != pantech.DNSZoneActive {
		if zone.Status == pantech.DNSZonePending {
			// The name servers are known once the zone leaves pending.
			zone, err = waitUntilNotPending(ctx, c, zone.ID)
			if err != nil {
				log.Fatal(err)
			}
		}
		fmt.Println("Set these name servers at the domain's registrar:")
		for _, ns := range zone.NameServers {
			fmt.Println("  " + ns)
		}
	}

	// Change the record if it exists, else create it.
	page, err := c.DNS.ListRecords(ctx, zone.ID, &pantech.ListDNSRecordsParams{Name: *name, Type: pantech.DNSRecordA})
	if err != nil {
		log.Fatal(err)
	}
	var record *pantech.DNSRecord
	if len(page.Data) == 1 {
		existing := page.Data[0]
		// If-Match: fail rather than overwrite a change made since we read it.
		record, err = c.DNS.UpdateRecord(ctx, zone.ID, existing.ID, pantech.UpdateDNSRecordRequest{TTL: ttl, Values: []string{*ip}}, pantech.WithIfMatch(existing.Version))
	} else {
		record, err = c.DNS.CreateRecord(ctx, zone.ID, pantech.CreateDNSRecordRequest{Name: *name, Type: pantech.DNSRecordA, TTL: ttl, Values: []string{*ip}})
	}
	if pantech.IsCode(err, pantech.CodeVersionConflict) {
		log.Fatal("the record changed while we were updating it: run again")
	}
	if err != nil {
		log.Fatal(err)
	}

	record, err = c.DNS.WaitForRecord(ctx, zone.ID, record.ID, &pantech.WaitOptions{Interval: 3 * time.Second})
	var failed *pantech.DNSFailureError
	if errors.As(err, &failed) {
		log.Fatalf("Pantech DNS refused the change (%s); the record keeps its previous values", failed.FailureCode)
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s %d IN A %v is served\n\n", record.FQDN, record.TTL, record.Values)

	file, err := c.DNS.ExportZone(ctx, zone.ID)
	if err != nil {
		log.Fatal(err)
	}
	_, _ = os.Stdout.Write(file)
}

func findOrCreateZone(ctx context.Context, c *pantech.Client, name string) (*pantech.DNSZone, error) {
	for z, err := range c.DNS.AllZones(ctx, &pantech.ListDNSZonesParams{Q: name}) {
		if err != nil {
			return nil, err
		}
		if z.Name == name && z.Status != pantech.DNSZoneDeleted && z.Status != pantech.DNSZoneDeleting && z.Status != pantech.DNSZoneFailed {
			return &z, nil
		}
	}
	return c.DNS.CreateZone(ctx, pantech.CreateDNSZoneRequest{Name: name})
}

func waitUntilNotPending(ctx context.Context, c *pantech.Client, zoneID string) (*pantech.DNSZone, error) {
	for {
		zone, err := c.DNS.GetZone(ctx, zoneID)
		if err != nil || zone.Status != pantech.DNSZonePending {
			return zone, err
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(2 * time.Second):
		}
	}
}
