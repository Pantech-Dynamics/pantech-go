// Command webhook-receiver accepts Pantech webhook deliveries and verifies
// their signatures before acting on them.
//
//	PANTECH_WEBHOOK_SECRET=whsec_… go run ./examples/webhook-receiver -addr :8080
package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"

	pantech "github.com/Pantech-Dynamics/pantech-go"
)

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()
	secret := os.Getenv("PANTECH_WEBHOOK_SECRET")
	if secret == "" {
		log.Fatal("set PANTECH_WEBHOOK_SECRET to the secret shown when the webhook was created")
	}

	http.HandleFunc("POST /pantech/webhook", func(w http.ResponseWriter, r *http.Request) {
		// Verify the raw bytes, before any JSON parsing.
		body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
		if err != nil {
			http.Error(w, "unreadable body", http.StatusBadRequest)
			return
		}
		event, err := pantech.ParseWebhook(secret, r.Header.Get(pantech.WebhookSignatureHeader), body, pantech.DefaultWebhookTolerance)
		if err != nil {
			log.Printf("rejected a delivery: %v", err)
			http.Error(w, "invalid signature", http.StatusBadRequest)
			return
		}
		// Events can arrive more than once and out of order: deduplicate on
		// event.ID before acting.
		log.Printf("event %s %s (delivery %s)", event.ID, event.Type, r.Header.Get(pantech.WebhookDeliveryIDHeader))
		w.WriteHeader(http.StatusNoContent) // any 2xx within 10 seconds is success
	})
	log.Printf("listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
