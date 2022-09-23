package cmd

// Args is a helpful struct that contains the global flags for the equivalence
// test binary.
type Args struct {
	// The relative or absolute path to the directory that contains the golden
	// files.
	GoldenFilesDirectory string

	// The relative or absolute path to the directory that contains the test
	// files and specifications.
	TestingFilesDirectory string

	// The relative or absolute path to the target Terraform binary.
	TerraformBinaryPath string
}
