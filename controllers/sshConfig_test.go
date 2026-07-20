package controllers

import (
	"os"
	"path/filepath"
	"sshbook/models"
	"testing"
)

func TestParseConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config")

	seed := "# comment\n" +
		"Host web\n" +
		"    HostName 1.2.3.4\n" +
		"    User root\n" +
		"    Port 2222\n" +
		"    IdentityFile ~/.ssh/id_ed25519\n" +
		"\n" +
		"Host bare\n" +
		"    HostName example.com\n" +
		"\n" +
		"Host *\n" + // wildcard — must be skipped
		"    User ignored\n"

	if err := os.WriteFile(path, []byte(seed), 0600); err != nil {
		t.Fatal(err)
	}

	conns, err := ParseConfig(path)
	if err != nil {
		t.Fatalf("ParseConfig: %v", err)
	}
	if len(conns) != 2 {
		t.Fatalf("want 2 connections, got %d: %+v", len(conns), conns)
	}

	web := conns[0]
	if web.Name != "web" || web.HostName != "1.2.3.4" || web.User != "root" ||
		web.Port != "2222" || web.IdentityFile != "~/.ssh/id_ed25519" {
		t.Errorf("web parsed wrong: %+v", web)
	}
	if conns[1].Name != "bare" || conns[1].HostName != "example.com" {
		t.Errorf("bare parsed wrong: %+v", conns[1])
	}
}

func TestParseConfigMissingFile(t *testing.T) {
	conns, err := ParseConfig(filepath.Join(t.TempDir(), "nope"))
	if err != nil {
		t.Fatalf("missing file should not error, got %v", err)
	}
	if len(conns) != 0 {
		t.Fatalf("want empty, got %+v", conns)
	}
}

func TestAppendConnectionRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config")
	want := models.Connection{
		Name:         "myserver",
		HostName:     "10.0.0.1",
		User:         "laurits",
		Port:         "22",
		IdentityFile: "~/.ssh/id_ed25519",
	}
	if err := AppendConnection(path, want); err != nil {
		t.Fatalf("AppendConnection: %v", err)
	}

	conns, err := ParseConfig(path)
	if err != nil {
		t.Fatalf("ParseConfig: %v", err)
	}
	if len(conns) != 1 {
		t.Fatalf("want 1, got %d", len(conns))
	}
	if conns[0] != want {
		t.Errorf("round-trip mismatch:\n got %+v\nwant %+v", conns[0], want)
	}
}
