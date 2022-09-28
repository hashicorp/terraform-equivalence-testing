package cmd

import (
	"fmt"
	"github.com/mitchellh/cli"
	"strings"

	"github.com/hashicorp/terraform-equivalence-testing/internal/terraform"
	"github.com/hashicorp/terraform-equivalence-testing/internal/tests"
)

func UpdateCommandFactory(ui cli.Ui) cli.CommandFactory {
	return func() (cli.Command, error) {
		return &updateCommand{
			ui: ui,
		}, nil
	}
}

type updateCommand struct {
	ui cli.Ui
}

func (cmd *updateCommand) Help() string {
	return strings.TrimSpace(`
Usage: terraform-equivalence-testing update --goldens=examples/example_golden_files --tests=examples/example_test_cases [--binary=terraform]

Update the equivalence test golden files.

This command will execute all the test cases within the tests directory, and write the outputs into the specified golden files directory. This will overwrite any existing golden files.

Note, that this command won't report any differences it finds. It will only update the golden files.`)
}

func (cmd *updateCommand) Run(args []string) int {
	flags, err := ParseFlags("update", args)
	if err != nil {
		cmd.ui.Error(err.Error())
		return 1
	}

	testCases, err := tests.ReadFrom(flags.TestingFilesDirectory)
	if err != nil {
		cmd.ui.Error(err.Error())
		return 1
	}
	cmd.ui.Output(fmt.Sprintf("Found %d test cases in %s\n", len(testCases), flags.TestingFilesDirectory))

	successfulTests := 0
	failedTests := 0

	tf := terraform.New(flags.TerraformBinaryPath)
	for _, test := range testCases {
		cmd.ui.Output(fmt.Sprintf("[%s]: starting...", test.Name))

		output, err := test.RunWith(tf)
		if err != nil {
			failedTests++
			if tfErr, ok := err.(terraform.Error); ok {
				cmd.ui.Output(fmt.Sprintf("[%s]: %s", test.Name, tfErr.Error()))
				continue
			}
			cmd.ui.Output(fmt.Sprintf("[%s]: unknown error (%v)", test.Name, err))
			continue
		}

		cmd.ui.Output(fmt.Sprintf("[%s]: updating golden files...", test.Name))

		if err := output.UpdateGoldenFiles(flags.GoldenFilesDirectory); err != nil {
			failedTests++
			cmd.ui.Output(fmt.Sprintf("[%s]: unknown error (%v)", test.Name, err))
			continue
		}

		successfulTests++
		cmd.ui.Output(fmt.Sprintf("[%s]: complete\n", test.Name))
	}

	cmd.ui.Output(fmt.Sprintf("Equivalence testing complete."))
	cmd.ui.Output(fmt.Sprintf("\tAttempted %d test(s).", len(testCases)))

	if successfulTests > 0 {
		cmd.ui.Output(fmt.Sprintf("\t%d test(s) were successfully updated.", successfulTests))
	}
	if failedTests > 0 {
		cmd.ui.Output(fmt.Sprintf("\t%d tests failed to update.", failedTests))
	}

	return 0
}

func (cmd *updateCommand) Synopsis() string {
	return "Update the equivalence test golden files."
}
