package models

import "fmt"

// Connection is a single Host block from ~/.ssh/config.
type Connection struct {
	Name         string // the Host alias
	HostName     string
	User         string
	Port         string
	IdentityFile string
}

// Summary is a one-line description for the list view.
func (c Connection) Summary() string {
	target := c.HostName
	if target == "" {
		target = c.Name
	}
	if c.User != "" {
		target = c.User + "@" + target
	}
	if c.Port != "" && c.Port != "22" {
		target += ":" + c.Port
	}
	return fmt.Sprintf("%s  (%s)", c.Name, target)
}
