// Command create-instance orders an instance, waits until it is provisioned
// and running, and prints its address. It spends real credit.
//
//	PANTECH_API_KEY=PAN_… go run ./examples/create-instance -name web-1 -plan individual -image ubuntu-24-04 -ssh-key sshk_…
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	pantech "github.com/Pantech-Dynamics/pantech-go"
)

func main() {
	name := flag.String("name", "web-1", "instance name")
	plan := flag.String("plan", "individual", "plan slug (see GET /plans)")
	image := flag.String("image", "ubuntu-24-04", "image slug (see GET /images)")
	sshKey := flag.String("ssh-key", "", "SSH key id (sshk_…)")
	flag.Parse()
	if *sshKey == "" {
		log.Fatal("-ssh-key is required: images have no password login")
	}

	c, err := pantech.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
	defer cancel()

	order, err := c.Instances.Create(ctx, pantech.CreateInstanceRequest{
		Name:      *name,
		PlanSlug:  *plan,
		ImageSlug: *image,
		SSHKeyID:  *sshKey,
		Tags:      pantech.Ptr(map[string]string{"env": "example"}),
	})
	if pantech.IsCode(err, "INSUFFICIENT_CREDIT") {
		log.Fatalf("not enough credit: %v", err)
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(os.Stderr, "ordered %s (order %s, %d %s upfront)\n", order.InstanceID, order.OrderID, order.AmountMinor, order.Currency)

	provisioned, err := c.Instances.WaitForOrder(ctx, order.OrderID, nil)
	if err != nil {
		log.Fatal(err)
	}
	if provisioned.OperationID != nil {
		if _, err := c.Operations.Wait(ctx, *provisioned.OperationID, nil); err != nil {
			log.Fatal(err)
		}
	}

	vm, err := c.Instances.Get(ctx, order.InstanceID)
	if err != nil {
		log.Fatal(err)
	}
	addr := "(no public address)"
	switch {
	case vm.PublicIPv4 != nil:
		addr = *vm.PublicIPv4
	case vm.SubnetID == nil && vm.PrivateIPv4 != nil:
		// A standard instance reports its one address in private_ipv4.
		addr = *vm.PrivateIPv4
	}
	fmt.Printf("%s %s %s\n", vm.ID, vm.ObservedState, addr)
}
