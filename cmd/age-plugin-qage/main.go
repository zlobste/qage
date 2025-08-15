package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"filippo.io/age"

	"github.com/zlobste/qage/internal/version"
	"github.com/zlobste/qage/pkg/qage"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	// Read first line - should be command
	if !scanner.Scan() {
		fatal("no input")
	}

	parts := strings.Fields(scanner.Text())
	if len(parts) < 1 {
		fatal("empty command")
	}

	cmd := parts[0]
	switch cmd {
	case "recipient-v1":
		handleRecipient(scanner, parts[1:])
	case "identity-v1":
		handleIdentity(scanner, parts[1:])
	case "--version":
		fmt.Println(version.String())
	default:
		fatal("unsupported command: " + cmd)
	}
}

func handleRecipient(scanner *bufio.Scanner, args []string) {
	if len(args) < 1 {
		fatal("recipient missing argument")
	}

	recipientStr := args[0]
	recipient, err := qage.ParseRecipient(recipientStr)
	if err != nil {
		fatal("invalid recipient: " + err.Error())
	}

	// Read file key from stdin (base64)
	if !scanner.Scan() {
		fatal("no file key")
	}
	fileKeyB64 := scanner.Text()

	fileKey, err := base64.StdEncoding.DecodeString(fileKeyB64)
	if err != nil {
		fatal("invalid file key: " + err.Error())
	}

	// Wrap the file key
	stanzas, err := recipient.Wrap(fileKey)
	if err != nil {
		fatal("wrap failed: " + err.Error())
	}

	// Output stanza
	if len(stanzas) != 1 {
		fatal("expected exactly one stanza")
	}

	s := stanzas[0]
	fmt.Printf("-> %s", s.Type)
	for _, arg := range s.Args {
		fmt.Printf(" %s", arg)
	}
	fmt.Println()

	bodyB64 := base64.StdEncoding.EncodeToString(s.Body)
	fmt.Println(bodyB64)
}

func handleIdentity(scanner *bufio.Scanner, args []string) {
	if len(args) < 1 {
		fatal("identity missing argument")
	}

	identityStr := args[0]

	// Parse identity from bech32 (expect QAGE-SECRET-KEY-1 prefix stripped)
	identity, err := qage.ParseIdentity(identityStr)
	if err != nil {
		fatal("invalid identity: " + err.Error())
	}

	// Read stanza from stdin
	if !scanner.Scan() {
		fatal("no stanza type")
	}
	stanzaLine := scanner.Text()

	if !strings.HasPrefix(stanzaLine, "-> ") {
		fatal("invalid stanza format")
	}

	stanzaParts := strings.Fields(stanzaLine[3:])
	if len(stanzaParts) < 1 || stanzaParts[0] != "qage" {
		return // Not for us
	}

	if !scanner.Scan() {
		fatal("no stanza body")
	}
	bodyB64 := scanner.Text()

	body, err := base64.StdEncoding.DecodeString(bodyB64)
	if err != nil {
		fatal("invalid stanza body: " + err.Error())
	}

	// Create stanza and unwrap
	ageStanza := &age.Stanza{
		Type: stanzaParts[0],
		Args: stanzaParts[1:],
		Body: body,
	}

	fileKey, err := identity.UnwrapStanza(ageStanza)

	if err != nil {
		return // Not for us or failed
	}

	// Output file key
	fileKeyB64 := base64.StdEncoding.EncodeToString(fileKey)
	fmt.Println(fileKeyB64)
}

func fatal(msg string) {
	fmt.Fprintf(os.Stderr, "age-plugin-qage: %s\n", msg)
	os.Exit(1)
}
