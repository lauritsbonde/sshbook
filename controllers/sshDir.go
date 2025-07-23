package controllers

import (
	"fmt"
	"log"
	"os"
	"sshbook/models"
	"strings"
)

func SshDirContents() models.SSHDirContents {
	// This function should return the contents of the SSH directory.
	userHome, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error getting user home directory: %v", err)
	}
	sshDir := userHome + "/.ssh"
	entries, err := os.ReadDir(sshDir)
	if err != nil {
		log.Fatalf("Error reading SSH directory: %v", err)
	}

	keys := []string{}

	for _, entry := range entries {
		if entry.IsDir() || entry.Name() == "config" || entry.Name() == "known_hosts" {
			continue // Skip directories
		}
		keys = append(keys, entry.Name())
	}

	// filter keys to only include those that are SSH keys
	var filteredKeys []string
	for _, key := range keys {
		if strings.HasSuffix(key, ".pub") || strings.HasSuffix(key, "_rsa") || strings.HasSuffix(key, "_dsa") || strings.HasSuffix(key, "_ed25519") {
			filteredKeys = append(filteredKeys, key)
		}
	}

	// create a list of hold the name of pairs
	var uniquePairs = make(map[string]bool)
	for _, key := range filteredKeys {
		// Check if the key is a public key
		if strings.HasSuffix(key, ".pub") {
			// Extract the base name (without .pub)
			baseName := strings.TrimSuffix(key, ".pub")
			// Check if the private key exists
			if _, err := os.Stat(sshDir + "/" + baseName); err == nil {
				uniquePairs[baseName] = true // Add to unique pairs
			}
		} else {
			// If it's not a public key, just add it as a private key
			uniquePairs[key] = true
		}
	}

	// Convert the map keys to a slice
	var uniqueKeys []string
	for key := range uniquePairs {
		uniqueKeys = append(uniqueKeys, key)
	}

	knownhosts, err := readKnownHostsFile()
	if err != nil {
		log.Fatalf("Error reading known_hosts file: %v", err)
	}

	return models.SSHDirContents{
		Keys:       uniqueKeys,
		Config:     sshDir + "/config",
		KnownHosts: knownhosts,
	}

}

func readKnownHostsFile() ([]string, error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error getting user home directory: %v", err)
	}
	knownHostsPath := userHome + "/.ssh/known_hosts"
	data, err := os.ReadFile(knownHostsPath)
	if err != nil {
		return nil, fmt.Errorf("error reading known_hosts file: %v", err)
	}

	// Split into lines
	lines := strings.Split(string(data), "\n")

	// Clean up empty lines
	cleaned := []string{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			cleaned = append(cleaned, line)
		}
	}

	return cleaned, nil
}
