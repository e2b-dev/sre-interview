package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

// durationAround returns target ± spread (uniform random).
func durationAround(target, spread time.Duration) time.Duration {
	jitter := time.Duration(rand.Int63n(int64(2*spread+1))) - spread
	return target + jitter
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func main() {
	start := time.Now()
	healthAfter := durationAround(5*time.Second, 1*time.Second)  // ~4–6s
	readyAfter := durationAround(15*time.Second, 3*time.Second)    // ~12–18s

	mainMux := http.NewServeMux()
	mainMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "ok")
	})

	internalMux := http.NewServeMux()
	internalMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if time.Since(start) < healthAfter {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprint(w, "starting")
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "ok")
	})
	internalMux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		if time.Since(start) < readyAfter {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprint(w, "starting")
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "ok")
	})

	mainPort := envOrDefault("PORT", "8080")
	internalPort := envOrDefault("INTERNAL_PORT", "9090")

	go func() {
		addr := ":" + internalPort
		log.Printf("internal server listening on %s", addr)
		if err := http.ListenAndServe(addr, internalMux); err != nil {
			log.Fatalf("internal server: %v", err)
		}
	}()

	addr := ":" + mainPort
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mainMux); err != nil {
		log.Fatal(err)
	}
}
