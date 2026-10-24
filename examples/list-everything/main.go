// Command list-everything prints what the API key's organization runs:
// instances, volumes, networks, databases and Kubernetes clusters.
//
//	PANTECH_API_KEY=PAN_… go run ./examples/list-everything
package main

import (
	"context"
	"fmt"
	"log"

	pantech "github.com/Pantech-Dynamics/pantech-go"
)

func main() {
	c, err := pantech.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	me, err := c.Me.Get(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("organization %s, key %q (%v)\n\n", me.OrganizationID, me.APIKey.Name, me.APIKey.Scopes)

	// All iterates over every page; the loop ends at the last one.
	fmt.Println("instances:")
	for vm, err := range c.Instances.All(ctx, &pantech.ListParams{Limit: 100}) {
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  %s  %-24s %s\n", vm.ID, vm.Name, vm.ObservedState)
	}

	fmt.Println("volumes:")
	for vol, err := range c.Volumes.All(ctx, nil) {
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  %s  %-24s %d GB\n", vol.ID, vol.Name, vol.SizeGB)
	}

	fmt.Println("networks:")
	for n, err := range c.Networks.All(ctx, nil) {
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  %s  %-24s %s\n", n.ID, n.Name, n.CIDR)
	}

	fmt.Println("databases:")
	for db, err := range c.Databases.All(ctx, nil) {
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  %s  %-24s %s %s\n", db.ID, db.Name, db.Engine, db.Version)
	}

	// Some lists are one page; their List returns the whole slice.
	clusters, err := c.Kubernetes.ListClusters(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("kubernetes clusters:")
	for _, k := range clusters {
		fmt.Printf("  %s  %-24s %s\n", k.ID, k.Name, k.KubernetesVersion)
	}
}
