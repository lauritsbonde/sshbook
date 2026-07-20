package models

import (
	"fmt"
)

type SSHDirContents struct {
	Keys        []string
	Config      string
	Connections []Connection
}

func (s SSHDirContents) String() string {
	result := "SSH Directory Contents:\n"
	result += "Keys:\n"
	for _, key := range s.Keys {
		result += fmt.Sprintf("- %s\n", key)
	}
	result += fmt.Sprintf("Config file: %s\n", s.Config)
	result += fmt.Sprintf("Connections: %d\n", len(s.Connections))
	return result
}
