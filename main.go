package main

import (
	"dagger/dagger-setup-php/internal/dagger"
	"fmt"
)

type Php struct {
	// The PHP extensions to install
	Extensions []string

	// The version of PHP to install
	Version string
}

func New(
	// The version of PHP to install
	// +optional
	// +default="8.4"
	version string,
	// The PHP extensions to install
	// +optional
	// +default=["bcmath","cli","common","curl","intl","mbstring","mysql","opcache","readline","xml","zip"]
	extensions []string,
) *Php {
	return &Php{
		Version:    version,
		Extensions: extensions,
	}
}

// Create a container with PHP and the specified extensions installed
func (m *Php) Setup() *dagger.Container {
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
		WithExec(append([]string{"apt", "install", "-y", "php" + m.Version}, ext...)).
		WithWorkdir("/app").
		WithEntrypoint([]string{"php"})
}

// Run a PHP command with the specified arguments
func (m *Php) Run(
	// The arguments to pass to the PHP command
	args []string,
) *dagger.Container {
	return m.Setup().WithExec(append([]string{"php"}, args...))
}
