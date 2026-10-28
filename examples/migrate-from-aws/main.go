// Command migrate-from-aws connects an AWS account to Pantech through a
// read-only role, scans it, plans the move, prints the price comparison and
// writes the plan as Terraform.
//
//	PANTECH_API_KEY=PAN_… go run ./examples/migrate-from-aws [-name "Production AWS"] [-region af-abj] [-out ./pantech-migration]
//
// It prints a link that creates the role in your AWS account. Create the
// stack, then paste its RoleArn output when asked, or pass it up front with
// -role-arn or PANTECH_ROLE_ARN. No AWS access keys are ever sent.
package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	pantech "github.com/Pantech-Dynamics/pantech-go"
)

func main() {
	name := flag.String("name", "AWS account", "a name for the connection and the plan")
	region := flag.String("region", "", "the Pantech region to plan for (default: the first region)")
	roleARN := flag.String("role-arn", os.Getenv("PANTECH_ROLE_ARN"), "the stack's RoleArn output (default: $PANTECH_ROLE_ARN, else asked for)")
	out := flag.String("out", "pantech-migration", "directory to write the Terraform files to")
	flag.Parse()

	c, err := pantech.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	// 1. Start a connection. Nothing is shared yet.
	conn, err := c.Migrate.CreateConnection(ctx, pantech.CreateCloudConnectionRequest{Name: *name})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connection", conn.ID)

	// 2. Create the read-only role in AWS, from the one-click link or the template.
	if err := os.MkdirAll(*out, 0o755); err != nil {
		log.Fatal(err)
	}
	if conn.Setup.QuickCreateURL != nil {
		fmt.Println("\nOpen this link, signed in to AWS, and create the stack:")
		fmt.Println("  " + *conn.Setup.QuickCreateURL)
	} else {
		path := filepath.Join(*out, "pantech-role.template.json")
		if err := os.WriteFile(path, []byte(conn.Setup.Template), 0o644); err != nil {
			log.Fatal(err)
		}
		fmt.Println("\nCreate a CloudFormation stack in AWS from", path)
	}

	// 3. Send the stack's RoleArn output.
	arn := strings.TrimSpace(*roleARN)
	if arn == "" {
		fmt.Print("\nPaste the stack's RoleArn output and press Enter: ")
		line, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil && line == "" {
			log.Fatal("no role ARN given")
		}
		arn = strings.TrimSpace(line)
	}
	if _, err := c.Migrate.UpdateConnection(ctx, conn.ID, pantech.UpdateCloudConnectionRequest{RoleARN: &arn}); err != nil {
		log.Fatal(err)
	}

	// 4. Wait for Pantech to assume the role.
	wait := &pantech.WaitOptions{Interval: 5 * time.Second, Timeout: 10 * time.Minute}
	conn, err = c.Migrate.WaitForConnection(ctx, conn.ID, wait)
	exitOnFailure(err)
	fmt.Println("connected to AWS account", deref(conn.AWSAccountID))

	// 5. Scan the account.
	scan, err := c.Migrate.StartScan(ctx, conn.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("scanning", scan.ID, "(most accounts take under five minutes)")
	scan, err = c.Migrate.WaitForScan(ctx, scan.ID, &pantech.WaitOptions{Interval: 10 * time.Second, Timeout: 30 * time.Minute})
	exitOnFailure(err)
	if n := scan.Counts; n != nil {
		fmt.Printf("found %d servers, %d disks, %d networks, %d databases, %d load balancers in %s\n",
			n.Instances, n.Volumes, n.Networks, n.Databases, n.LoadBalancers, strings.Join(scan.Regions, ", "))
	}
	for _, w := range scan.Warnings {
		fmt.Println("  not read:", w)
	}

	// 6. Plan the move. Nothing is created on Pantech.
	req := pantech.CreateMigrationRequest{ScanID: scan.ID, Name: *name}
	if *region != "" {
		req.Region = region
	}
	plan, err := c.Migrate.CreateMigration(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nplan %s for %s: %d items planned, %d excluded, %d not supported yet\n",
		plan.ID, plan.Region, plan.Summary.Planned, plan.Summary.Excluded, plan.Summary.Unsupported)
	for _, item := range plan.Items {
		for _, f := range item.Flags {
			if f.Blocking {
				fmt.Printf("  %s: %s\n", item.Key, f.Message)
			}
		}
	}

	// 7. The price comparison.
	if e := plan.Estimate; e != nil {
		if e.PantechPriced && e.PantechMonthlyMinor != nil {
			fmt.Printf("Pantech: %s a month (covers %s)\n", money(*e.PantechMonthlyMinor, e.Currency), strings.Join(e.Covers, ", "))
		} else {
			fmt.Println("Pantech: part of the plan could not be priced")
		}
		if len(e.NotCovered) > 0 {
			fmt.Println("  not in the price:", strings.Join(e.NotCovered, ", "))
		}
		if e.SourceKnown && e.SourceMonthlyMinor != nil {
			fmt.Printf("AWS: %s a month, the whole bill averaged over %d months\n", money(*e.SourceMonthlyMinor, deref(e.SourceCurrency)), e.SourceMonths)
		}
	}

	// 8. Export Terraform.
	tf, err := c.Migrate.ExportTerraform(ctx, plan.ID)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range tf.Files {
		path := filepath.Join(*out, filepath.Base(f.Name))
		if err := os.WriteFile(path, []byte(f.Content), 0o644); err != nil {
			log.Fatal(err)
		}
		fmt.Println("wrote", path)
	}
	fmt.Println("\nReview the files, fill in the variables, then run terraform init and terraform plan in", *out)
}

// exitOnFailure explains a failed connection or scan; its message says what to do.
func exitOnFailure(err error) {
	var failed *pantech.CloudFailureError
	if errors.As(err, &failed) && failed.Failure != nil {
		log.Fatalf("the %s failed: %s (%s)", failed.Kind, failed.Failure.Message, failed.Failure.Code)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func money(minor int64, currency string) string {
	return fmt.Sprintf("%s %d.%02d", currency, minor/100, minor%100)
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
