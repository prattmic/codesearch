// Package get provides functionality similar to the 'go get' command.
package get

import (
	"bytes"
	"fmt"
	"go/build"
	"os"
	"os/exec"
)

// findDeps adds all of the (recursive) dependencies of pkg to deps.
func findDeps(pkg string, deps map[string]bool) error {
	p, err := build.Import(pkg, ".", build.ImportComment)
	if err != nil {
		return err
	}

	for _, i := range p.Imports {
		if _, ok := deps[i]; !ok {
			deps[i] = true
			if err := findDeps(i, deps); err != nil {
				return err
			}
		}
	}

	return nil
}

// PackageDependencies returns all of the (recursive) dependencies of the
// specified package.
func PackageDependencies(pkg string) ([]string, error) {
	deps := make(map[string]bool)

	if err := findDeps(pkg, deps); err != nil {
		return nil, err
	}

	var out []string
	for dep := range deps {
		out = append(out, dep)
	}
	return out, nil
}

// execError wraps an error from exec.Command, including the command stderr.
type execError struct {
	Stderr string
	Err    error
}

func (e execError) Error() string {
	return fmt.Sprintf("%s; stderr: %s", e.Err, e.Stderr)
}

// Get runs 'go get', downloading the package and its dependencies into gopath.
func Get(pkg string, gopath string) error {
	cmd := exec.Command("go", "get", "-d", pkg)
	cmd.Env = append(cmd.Env, fmt.Sprintf("PATH=%s", os.Getenv("PATH")))
	cmd.Env = append(cmd.Env, fmt.Sprintf("GOPATH=%s", gopath))

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return execError{Err: err, Stderr: stderr.String()}
	}

	return nil
}
