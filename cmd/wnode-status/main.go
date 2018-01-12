package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/status-im/status-go/geth/api"
	"github.com/status-im/status-go/geth/params"
)

func main() {
	flag.Parse()

	config, err := makeNodeConfig()
	if err != nil {
		log.Fatalf("Making config failed: %v", err)
	}

	printHeader(config)

	if *injectAccounts {
		if err := LoadTestAccounts(config.DataDir); err != nil {
			log.Fatalf("Failed to load test accounts: %v", err)
		}
	}

	backend := api.NewStatusBackend()
	started, err := backend.StartNode(config)
	if err != nil {
		log.Fatalf("Node start failed: %v", err)
		return
	}

	// wait till node is started
	<-started

	if *injectAccounts {
		if err := InjectTestAccounts(backend.NodeManager()); err != nil {
			log.Fatalf("Failed to inject accounts: %v", err)
		}
	}

	// wait till node has been stopped
	node, err := backend.NodeManager().Node()
	if err != nil {
		log.Fatalf("Getting node failed: %v", err)
		return
	}

	// start debug server and collecting metrics
	if *statsEnabled {
		log.Printf("Stats enabled on %s", *statsAddr)

		go func() {
			http.Handle("/metrics", promhttp.Handler())
			log.Fatal(http.ListenAndServe(*statsAddr, nil))
		}()
	}

	node.Wait()
}

// printHeader prints command header
func printHeader(config *params.NodeConfig) {
	fmt.Println("Starting Whisper V5 node...")
	fmt.Printf("Config: %s\n", config.WhisperConfig)
}
