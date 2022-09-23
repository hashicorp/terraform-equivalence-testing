package main

import (
	"fmt"

	"github.com/hashicorp/terraform-equivalence-testing/internal/cmd"
)

func main() {
	command, err := cmd.Root()
	if err != nil {
		fmt.Printf("failed to run equivalence tests: %v", err)
		return
	}

	if err := command.Execute(); err != nil {
		fmt.Printf("failed to run equivalence tests: %v", err)
	}
}
