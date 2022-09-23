package cmd

import (
	"github.com/spf13/cobra"

	"github.com/hashicorp/terraform-equivalence-testing/internal/terraform"
	"github.com/hashicorp/terraform-equivalence-testing/internal/tests"
)

func update(args *Args) *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Update the equivalence test golden files.",
		Long: `Update the equivalence test golden files.

This command will execute all the test cases within the tests directory, and write the outputs into the specified golden files directory. This will overwrite any existing golden files. 

Note, that this command won't report any differences it finds. It will only update the golden files.`,
		Example: "equivalence-test update --goldens=examples/example_golden_files --tests=examples/example_test_cases --binary=terraform",
		RunE: func(cmd *cobra.Command, _ []string) error {
			testCases, err := tests.ReadFrom(args.TestingFilesDirectory)
			if err != nil {
				return err
			}
			cmd.Printf("Found %d test cases in %s\n", len(testCases), args.TestingFilesDirectory)

			successfulTests := 0
			failedTests := 0

			tf := terraform.New(args.TerraformBinaryPath)
			for _, test := range testCases {
				cmd.Printf("\n[%s]: starting...\n", test.Name)

				output, err := test.RunWith(tf)
				if err != nil {
					failedTests++
					if tfErr, ok := err.(terraform.Error); ok {
						cmd.Printf("[%s]: %s\n", test.Name, tfErr.Error())
						continue
					}
					cmd.Printf("[%s]: unknown error (%v)\n", test.Name, err)
					continue
				}

				cmd.Printf("[%s]: updating golden files...\n", test.Name)

				if err := output.UpdateGoldenFiles(args.GoldenFilesDirectory); err != nil {
					failedTests++
					cmd.Printf("[%s]: unknown error (%v)\n", test.Name, err)
					continue
				}

				successfulTests++
				cmd.Printf("[%s]: complete\n", test.Name)
			}

			cmd.Printf("\nEquivalence testing complete.\n")
			cmd.Printf("\tAttempted %d test(s).\n", len(testCases))

			if successfulTests > 0 {
				cmd.Printf("\t%d test(s) were successfully updated.\n", successfulTests)
			}
			if failedTests > 0 {
				cmd.Printf("\t%d tests failed to update.\n", failedTests)
			}

			return nil
		},
	}
}
