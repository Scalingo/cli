package term

import (
	"os"
	"os/exec"
)

func makeRaw(f *os.File) error {
	return stty(f, "-icanon", "-echo").Run()
}

func restore(f *os.File) error {
	return stty(f, "icanon", "echo").Run()
}

// helpers
func stty(f *os.File, args ...string) *exec.Cmd {
	c := exec.Command("stty", args...)
	c.Stdin = f
	return c
}
