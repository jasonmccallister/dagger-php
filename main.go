package main

import (
	"context"
	"dagger/dagger-setup-php/internal/dagger"
	"fmt"
)

type Php struct {
	EnableXdebug bool

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
	// Enable Xdebug for the PHP container
	// +optional
	// +default=false
	enableXdebug bool,
) *Php {
	return &Php{
		EnableXdebug: enableXdebug,
		Version:      version,
		Extensions:   extensions,
	}
}

// Create a container with PHP and the specified extensions installed
func (m *Php) Setup() *dagger.Container {
	var ext []string
	for _, e := range m.Extensions {
		ext = append(ext, fmt.Sprintf("php%s-%s", m.Version, e))
	}

	if m.EnableXdebug {
		ext = append(ext, fmt.Sprintf("php%s-xdebug", m.Version))
	}

	c := dag.Container().
		From("ubuntu:24.04").
		WithEnvVariable("DEBIAN_FRONTEND", "noninteractive").
		WithExec([]string{"apt", "update", "-y"}).
		WithExec([]string{"apt", "install", "-y", "-q", "software-properties-common"}).
		WithExec([]string{"add-apt-repository", "ppa:ondrej/php"}).
		WithExec([]string{"apt", "update", "-y"}).
		WithExec(append([]string{"apt", "install", "-y", "php" + m.Version}, ext...)).
		WithWorkdir("/app").
		WithEntrypoint([]string{"php"})

	if m.EnableXdebug {
		c = c.WithEnvVariable("XDEBUG_MODE", "coverage")
	}

	return c
}

// Run a PHP command with the specified arguments
func (m *Php) Run(
	ctx context.Context,
	// The arguments to pass to the PHP command
	args []string,
) (string, error) {
	return m.Setup().WithExec(append([]string{"php"}, args...)).Stdout(ctx)
}
