package terraform

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/hashicorp/terraform-equivalence-testing/internal/files"

	"github.com/hashicorp/terraform-exec/tfexec"
)

// Terraform is an interface that can execute a single equivalence test within a
// directory using the ExecuteTest method.
//
// We hold this in an interface, so we can mock it for testing purposes.
type Terraform interface {
	// ExecuteTest executes a series of terraform commands in order and returns the
	// output of the apply and plan steps, the Terraform state, and any additionally
	// requested files.
	ExecuteTest(directory string, includeFiles []string) (map[string]*files.File, error)

	// Version returns the version of the underlying Terraform binary.
	Version() string
}

// New returns a Terraform compatible struct that executes the tests using the
// Terraform binary provided in the argument.
func New(binary string) (Terraform, error) {

	// First, sanity check binary actually points to a Terraform binary file.
	//
	// We do this by fetching the version using tfexec. tfexec tries to be
	// clever and look up cached provider versions as well, but we're not
	// interested in this, so we just set the working directory to be the
	// current directory and tfexec just won't find any terraform or provider
	// files.
	//
	// Note, ideally we could actually just tfexec for everything. tfexec
	// doesn't (yet) support returning JSON files from the apply command so for
	// now we do the rest ourselves. Something to revisit in the future.
	tf, err := tfexec.NewTerraform(".", binary)
	if err != nil {
		return nil, err
	}

	version, _, err := tf.Version(context.Background(), true)
	if err != nil {
		return nil, err
	}

	return &terraform{
		binary:  binary,
		version: version.String(),
	}, nil
}

type terraform struct {
	binary  string
	version string
}

func (t *terraform) Version() string {
	return t.version
}

func (t *terraform) ExecuteTest(directory string, includeFiles []string) (map[string]*files.File, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	if err := os.Chdir(directory); err != nil {
		return nil, err
	}
	defer os.Chdir(wd)

	if err := t.init(); err != nil {
		return nil, err
	}

	savedFiles := map[string]*files.File{}
	if savedFiles["plan"], err = t.plan(); err != nil {
		return nil, err
	}
	if savedFiles["apply.json"], err = t.apply(); err != nil {
		return nil, err
	}
	if savedFiles["state.json"], err = t.showState(); err != nil {
		return nil, err
	}
	if savedFiles["plan.json"], err = t.showPlan(); err != nil {
		return nil, err
	}

	for _, includeFile := range includeFiles {
		raw, err := os.ReadFile(includeFile)
		if err != nil {
			return nil, fmt.Errorf("could not read additional file (%s): %v", includeFile, err)
		}
		if savedFiles[includeFile], err = files.NewFile(includeFile, raw); err != nil {
			return nil, fmt.Errorf("could not unmarshal additional file (%s): %v", includeFile, err)
		}
	}

	return savedFiles, nil
}

func (t *terraform) init() error {
	_, err := run(exec.Command(t.binary, "init"), "init")
	if err != nil {
		return err
	}
	return nil
}

func (t *terraform) plan() (*files.File, error) {
	capture, err := run(exec.Command(t.binary, "plan", "-out=equivalence_test_plan", "-no-color"), "plan")
	if err != nil {
		return nil, err
	}
	return files.NewRawFile(capture.ToString()), nil
}

func (t *terraform) apply() (*files.File, error) {
	capture, err := run(exec.Command(t.binary, "apply", "-json", "equivalence_test_plan"), "apply")
	if err != nil {
		return nil, err
	}

	json, err := capture.ToJson(true)
	if err != nil {
		return nil, err
	}
	return files.NewJsonFile(json), nil
}

func (t *terraform) showPlan() (*files.File, error) {
	capture, err := run(exec.Command(t.binary, "show", "-json", "equivalence_test_plan"), "show plan")
	if err != nil {
		return nil, err
	}

	json, err := capture.ToJson(false)
	if err != nil {
		return nil, err
	}
	return files.NewJsonFile(json), nil
}

func (t *terraform) showState() (*files.File, error) {
	capture, err := run(exec.Command(t.binary, "show", "-json"), "show state")
	if err != nil {
		return nil, err
	}

	json, err := capture.ToJson(false)
	if err != nil {
		return nil, err
	}
	return files.NewJsonFile(json), nil
}

func run(cmd *exec.Cmd, command string) (*capture, error) {
	capture := Capture(cmd)
	if err := cmd.Run(); err != nil {
		return capture, Error{
			Command:   command,
			Go:        err,
			Terraform: capture.ToError(),
		}
	}
	return capture, nil
}
