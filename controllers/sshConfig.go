package controllers

import (
	"bufio"
	"fmt"
	"os"
	"sshbook/models"
	"strings"
)

// ParseConfig reads an ssh config file and returns its Host blocks as
// connections. A missing file is not an error — it returns an empty slice.
func ParseConfig(path string) ([]models.Connection, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("open ssh config: %w", err)
	}
	defer f.Close()

	var connections []models.Connection
	var cur *models.Connection

	flush := func() {
		if cur != nil {
			connections = append(connections, *cur)
			cur = nil
		}
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value := splitField(line)
		switch strings.ToLower(key) {
		case "host":
			flush()
			// Skip wildcard/pattern hosts — not connectable targets.
			if strings.ContainsAny(value, "*?!") {
				continue
			}
			cur = &models.Connection{Name: value}
		case "hostname":
			if cur != nil {
				cur.HostName = value
			}
		case "user":
			if cur != nil {
				cur.User = value
			}
		case "port":
			if cur != nil {
				cur.Port = value
			}
		case "identityfile":
			if cur != nil {
				cur.IdentityFile = value
			}
		}
	}
	flush()

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan ssh config: %w", err)
	}
	return connections, nil
}

// AppendConnection writes a new Host block to the ssh config file, creating
// the file if it does not exist.
func AppendConnection(path string, c models.Connection) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("open ssh config for write: %w", err)
	}
	defer f.Close()

	var b strings.Builder
	fmt.Fprintf(&b, "\nHost %s\n", c.Name)
	if c.HostName != "" {
		fmt.Fprintf(&b, "    HostName %s\n", c.HostName)
	}
	if c.User != "" {
		fmt.Fprintf(&b, "    User %s\n", c.User)
	}
	if c.Port != "" {
		fmt.Fprintf(&b, "    Port %s\n", c.Port)
	}
	if c.IdentityFile != "" {
		fmt.Fprintf(&b, "    IdentityFile %s\n", c.IdentityFile)
	}

	if _, err := f.WriteString(b.String()); err != nil {
		return fmt.Errorf("write ssh config: %w", err)
	}
	return nil
}

// splitField splits an ssh config line into its keyword and value, handling
// both "Key value" and "Key = value" forms.
func splitField(line string) (string, string) {
	if i := strings.IndexAny(line, " \t="); i >= 0 {
		key := strings.TrimSpace(line[:i])
		value := strings.TrimSpace(strings.TrimLeft(line[i:], " \t="))
		return key, value
	}
	return line, ""
}
