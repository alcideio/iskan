package util

import (
	"io/ioutil"
	"os"
	"os/exec"
)

var (
	DefaultCmdRunner = &cmdRunner{}
)

type File interface {
	Name() string
	Read([]byte) (int, error)
}

// CmdRunner the cmdRunner to the outside "world".
// Wraps methods that modify global state and hence make the code that use them very hard to test.
type CmdRunner interface {
	Environ() []string
	LookPath(string) (string, error)
	RunCmd(cmd *exec.Cmd) ([]byte, error)
	TempFile(dir, pattern string) (File, error)
	Remove(name string) error
}

type cmdRunner struct {
}

func (a *cmdRunner) Environ() []string {
	return os.Environ()
}

func (a *cmdRunner) RunCmd(cmd *exec.Cmd) ([]byte, error) {
	return cmd.CombinedOutput()
}

func (a *cmdRunner) TempFile(dir, pattern string) (File, error) {
	return ioutil.TempFile(dir, pattern)
}

func (a *cmdRunner) Remove(name string) error {
	return os.Remove(name)
}

func (a *cmdRunner) LookPath(file string) (string, error) {
	return exec.LookPath(file)
}
