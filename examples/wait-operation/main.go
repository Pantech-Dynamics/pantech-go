// Command wait-operation stops or starts an instance and follows the change,
// either by polling or over the event stream.
//
//	PANTECH_API_KEY=PAN_… go run ./examples/wait-operation -instance vm_… -action stop [-stream]
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"time"

	pantech "github.com/Pantech-Dynamics/pantech-go"
)

func main() {
	id := flag.String("instance", "", "instance id (vm_…)")
	action := flag.String("action", "stop", "stop or start")
	stream := flag.Bool("stream", false, "follow the operation over Server-Sent Events instead of polling")
	flag.Parse()
	if *id == "" {
		log.Fatal("-instance is required")
	}

	c, err := pantech.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	var accepted *pantech.Accepted
	switch *action {
	case "stop":
		accepted, err = c.Instances.Stop(ctx, *id)
	case "start":
		accepted, err = c.Instances.Start(ctx, *id)
	default:
		log.Fatalf("unknown -action %q", *action)
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("operation", accepted.OperationID)

	if *stream {
		for op, err := range c.Operations.Stream(ctx, accepted.OperationID) {
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("  status:", op.Status)
		}
		return
	}

	op, err := c.Operations.Wait(ctx, accepted.OperationID, &pantech.WaitOptions{Interval: 2 * time.Second})
	var failed *pantech.OperationError
	if errors.As(err, &failed) {
		// failure.reason is safe to show to users; branch on failure.code.
		log.Fatalf("the change failed: %s (%s)", failed.Operation.Failure.Reason, failed.Operation.Failure.Code)
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("done:", op.Status)
}
