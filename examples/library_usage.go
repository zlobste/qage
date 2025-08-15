package main

import (
	"fmt"
	"log"

	"github.com/zlobste/qage/pkg/qage"
)

func main() {
	// Generate a new identity
	identity, err := qage.NewIdentity()
	if err != nil {
		log.Fatal(err)
	}

	// Get the recipient
	recipient := identity.Recipient()

	// Get string representations
	identityStr, err := identity.String()
	if err != nil {
		log.Fatal(err)
	}

	recipientStr, err := recipient.String()
	if err != nil {
		log.Fatal(err)
	}

	// Parse the recipient string back
	parsed, err := qage.ParseRecipient(recipientStr)
	if err != nil {
		log.Fatal(err)
	}

	parsedRecipientStr, err := parsed.String()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Generated identity: %s\n", identityStr)
	fmt.Printf("Public recipient: %s\n", recipientStr)
	fmt.Printf("Parsed recipient: %s\n", parsedRecipientStr)
	fmt.Printf("Recipients match: %t\n", recipientStr == parsedRecipientStr)

	// Test parsing the identity string back
	parsedIdentity, err := qage.ParseIdentity(identityStr)
	if err != nil {
		log.Fatal(err)
	}

	parsedIdentityStr, err := parsedIdentity.String()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Identities match: %t\n", identityStr == parsedIdentityStr)
}
