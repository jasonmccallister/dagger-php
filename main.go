// A generated module for SetupPhp functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"dagger/dagger-setup-php/internal/dagger"
	"fmt"
)

type SetupPhp struct {
	// The PHP extensions to install
	Extensions []string

	// The version of PHP to install
	Version string
}

func New(
	// +optional
	// +default="8.4"
	version string,
	// +optional
	// +default=["bcmath","cli","common","curl","intl","mbstring","mysql","opcache","readline","xml","zip"]
	extensions []string,
) *SetupPhp {
	return &SetupPhp{
		Version:    version,
		Extensions: extensions,
	}
}

// Create a container with PHP and the specified extensions installed
func (m *SetupPhp) Build() *dagger.Container {
	var ext []string
	for _, e := range m.Extensions {
		ext = append(ext, fmt.Sprintf("php%s-%s", m.Version, e))
	}

	return dag.Container().
		From("ubuntu:24.04").
		WithEnvVariable("DEBIAN_FRONTEND", "noninteractive").
		WithExec([]string{"apt", "update", "-y"}).
		WithExec([]string{"apt", "install", "-y", "-q", "software-properties-common"}).
		WithExec([]string{"add-apt-repository", "ppa:ondrej/php"}).
		WithExec([]string{"apt", "update", "-y"}).
		WithExec(append([]string{"apt", "install", "-y", "php" + m.Version}, ext...))
}

// Run a PHP command with the specified arguments
func (m *SetupPhp) Run(
	// The arguments to pass to the PHP command
	args []string,
) *dagger.Container {
	return m.Build().WithExec(append([]string{"php"}, args...))
}
