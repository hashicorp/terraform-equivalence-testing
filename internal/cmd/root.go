package cmd

import "github.com/spf13/cobra"

// Root returns the root command for the equivalence test binary.
//
// The Root command is made up of two sub commands (update and diff), and
// contains three global flags that are specified in the Args struct.
func Root() (*cobra.Command, error) {
	command := &cobra.Command{
		Use:   "equivalence-test",
		Short: "Compare or update Terraform equivalence test golden files.",
		Long: `Compare or update Terraform equivalence test golden files.

The Terraform equivalence test framework can be used to find differences between Terraform or Provider versions, or even between different Terraform configurations.

There are two sub commands (diff and update) and three flags (goldens, tests, and binary).`,
		Example: "equivalence-test --help",
	}

	args := Args{}
	command.PersistentFlags().StringVarP(&args.GoldenFilesDirectory, "goldens", "g", "", "Absolute or relative path to the directory containing the golden files.")
	command.PersistentFlags().StringVarP(&args.TestingFilesDirectory, "tests", "t", "", "Absolute or relative path to the directory containing the tests and specifications.")
	command.PersistentFlags().StringVarP(&args.TerraformBinaryPath, "binary", "b", "terraform", "Absolute or relative path to the target Terraform binary.")

	if err := command.MarkPersistentFlagRequired("goldens"); err != nil {
		return nil, err
	}
	if err := command.MarkPersistentFlagRequired("tests"); err != nil {
		return nil, err
	}

	command.AddCommand(diff(&args))
	command.AddCommand(update(&args))
	return command, nil
}
