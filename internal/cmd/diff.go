package cmd

import (
	"github.com/spf13/cobra"

	"github.com/hashicorp/terraform-equivalence-testing/internal/terraform"
	"github.com/hashicorp/terraform-equivalence-testing/internal/tests"
)

func diff(args *Args) *cobra.Command {
	return &cobra.Command{
		Use:   "diff",
		Short: "Compare and report the diff between a fresh run of the equivalence tests and the golden files.",
		Long: `Compare and report the diff between a fresh run of the equivalence tests and the golden files.

This command will execute all the test cases within the tests directory, and report any differences between the output and the existing golden files.`,
		Example: "equivalence-test diff --goldens=examples/example_golden_files --tests=examples/example_test_cases --binary=terraform",
		RunE: func(cmd *cobra.Command, _ []string) error {
			testCases, err := tests.ReadFrom(args.TestingFilesDirectory)
			if err != nil {
				return err
			}
			cmd.Printf("Found %d test cases in %s\n", len(testCases), args.TestingFilesDirectory)

			successfulTests := 0
			testsWithDiffs := 0
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

				cmd.Printf("[%s]: computing diffs...\n", test.Name)

				files, err := output.ComputeDiff(args.GoldenFilesDirectory)
				if err != nil {
					failedTests++
					cmd.Printf("[%s]: unknown error (%v)\n", test.Name, err)
					continue
				}

				newFileCount := 0
				noChangeCount := 0
				changeCount := 0

				for file, diff := range files {
					switch diff {
					case tests.NewFile:
						newFileCount++
						cmd.Printf("[%s]: %s was a new file\n", test.Name, file)
					case tests.NoChange:
						noChangeCount++
						cmd.Printf("[%s]: %s had no diffs\n", test.Name, file)
					default:
						changeCount++
						cmd.Printf("[%s]: %s had diffs:\n%s\n", test.Name, file, diff)
					}
				}

				successfulTests++
				if newFileCount+changeCount > 0 {
					testsWithDiffs++
				}

				cmd.Printf("[%s]: complete\n", test.Name)
			}

			cmd.Printf("\nEquivalence testing complete.\n")
			cmd.Printf("\tAttempted %d test(s).\n", len(testCases))

			if successfulTests > 0 {
				cmd.Printf("\t%d test(s) were successful.\n", successfulTests)
			}

			if testsWithDiffs > 0 {
				cmd.Printf("\t%d test(s) had diffs.\n", testsWithDiffs)
			}

			if failedTests > 0 {
				cmd.Printf("\t%d tests failed.\n", failedTests)
			}

			return nil
		},
	}
}
