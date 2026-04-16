package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mostlygeek/llama-swap/proxy"
)

const version = "0.1.0"

func main() {
	var (
		configFile  = flag.String("config", "config.yaml)
		listenAddr  = flag.String("listen", ":8080", "address to listen on (host:port)")
		showVersion = flag.Bool("version", false, "print version and exit")
		logLevel    = flag.String("log-level", "info", "log level: debug, info, warn, error")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("llama-swap version %s\n", version)
		os.Exit(0)
	}

	// Load configuration
	cfg, err := proxy.LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("failed to load config from %q: %v", *configFile, err)
	}

	// Override listen address if set in config and not overridden by flag
	if cfg.ListenAddress != "" && !isFlagSet("listen") {
		*listenAddr = cfg.ListenAddress
	}

	log.Printf("llama-swap v%s starting", version)
	log.Printf("config: %s", *configFile)
	log.Printf("listen: %s", *listenAddr)
	log.Printf("log-level: %s", *logLevel)

	// Create and start the proxy server
	server, err := proxy.New(cfg, *listenAddr)
	if err != nil {
		log.Fatalf("failed to create proxy server: %v", err)
	}

	// Handle graceful shutdown on SIGINT / SIGTERM
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		log.Printf("received signal %s, shutting down...", sig)
 err := server.Start(); err !=: %v", err)
	}

	log.Println("llama-swap stopped")
}

// isFlagSet returns true if the named flag was explicitly set on the command line.
func isFlagSet(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
