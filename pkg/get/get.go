// Package get provides functionality similar to the 'go get' command.
package get

import (
	"go/build"
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
