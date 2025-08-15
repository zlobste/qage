package main

import (
	"bytes"
	"testing"

	"github.com/zlobste/qage/cmd/qage/cmd"
)

func TestRootHelp(t *testing.T) {
	b := &bytes.Buffer{}
	rootCmd := cmd.NewRootCmd()
	rootCmd.SetOut(b)
	rootCmd.SetErr(b)
	rootCmd.SetArgs([]string{"--help"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("help: %v", err)
	}
	output := b.String()
	if !bytes.Contains([]byte(output), []byte("qage provides post-quantum")) {
		t.Fatalf("missing help text in: %s", output)
	}
}

func TestKeygenCommandDry(t *testing.T) {
	b := &bytes.Buffer{}
	rootCmd := cmd.NewRootCmd()
	rootCmd.SetOut(b)
	rootCmd.SetErr(b)
	rootCmd.SetArgs([]string{"keygen"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("keygen: %v", err)
	}
	output := b.String()
	if !bytes.Contains([]byte(output), []byte("QAGE-SECRET-KEY-1")) {
		t.Fatalf("expected key output, got: %s", output)
	}
}
