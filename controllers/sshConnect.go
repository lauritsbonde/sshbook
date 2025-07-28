package controllers

import (
	"os"
	"os/exec"
	"sshbook/models"
	"strings"

	ui "github.com/gizak/termui/v3"
)

// FindSSHConfigEntry finds the first matching SSH config entry for a given alias
func FindSSHConfigEntry(host string, config []models.SSHConfigEntry) *models.SSHConfigEntry {
	for _, entry := range config {
		if entry.Host == host {
			return &entry
		}
	}
	return nil
}

func OpenSSHSession(app *models.AppState) error {
	selected := app.CurrentSelectedHost()
	if selected == "" {
		return nil // no-op
	}

	cleanHost := strings.Split(selected, " ")[0] // remove known_hosts fingerprint etc.

	entry := FindSSHConfigEntry(cleanHost, app.SSHDirContents.SSHConfigEntries)

	sshArgs := []string{}
	if entry != nil {
		if entry.User != "" {
			sshArgs = append(sshArgs, "-l", entry.User)
		}
		if entry.IdentityFile != "" {
			sshArgs = append(sshArgs, "-i", entry.IdentityFile)
		}
		if entry.HostName != "" {
			sshArgs = append(sshArgs, entry.HostName)
		} else {
			sshArgs = append(sshArgs, cleanHost)
		}
	} else {
		// fallback: use current OS user and the host as-is
		sshArgs = append(sshArgs, cleanHost)
	}

	cmd := exec.Command("ssh", sshArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Close TermUI before attaching terminal
	ui.Close()
	return cmd.Run()
}
