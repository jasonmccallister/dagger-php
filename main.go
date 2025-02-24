package main

import (
	"dagger/dagger-setup-php/internal/dagger"
	"fmt"
	"slices"
)

var (
	defaultExtensions = []string{"bcmath", "cli", "common", "curl", "intl", "mbstring", "mysql", "opcache", "readline", "sqlite3", "xml", "zip"}
)

type Php struct {
	EnableXdebug  bool              // +private
	Extensions    []string          // +private
	Version       string            // +private
	UbuntuVersion string            // +private
	Source        *dagger.Directory // +private
}

func New(
	// The version of PHP to install
	// +optional
	// +default="8.4"
	version string,
	// Additional PHP extensions to install
	// +optional
	extensions []string,
	// Enable Xdebug for the PHP container
	// +optional
	// +default=false
	enableXdebug bool,
	// The version of Ubuntu to use
	// +optional
	// +default="24.04"
	ubuntuVersion string,
	// The directory that contains your PHP project.
	// +optional
	source *dagger.Directory,
) *Php {
	return &Php{
		EnableXdebug: enableXdebug,
		Version:      version,
		Extensions: slices.Compact(
			append(defaultExtensions, extensions...),
		),
		UbuntuVersion: ubuntuVersion,
		Source:        source,
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
		From("ubuntu:"+m.UbuntuVersion).
		WithEnvVariable("DEBIAN_FRONTEND", "noninteractive").
		WithExec([]string{"apt", "update", "-y"}).
		WithExec([]string{"apt", "install", "-y", "-q", "software-properties-common", "git"}).
		WithExec([]string{"add-apt-repository", "ppa:ondrej/php"}).
		WithExec([]string{"apt", "update", "-y"}).
		WithExec(append([]string{"apt", "install", "-y", "php" + m.Version}, ext...)).
		WithWorkdir("/app").
		WithDefaultTerminalCmd([]string{"bash"}).
		WithEntrypoint([]string{"php"})

	if m.Source != nil {
		c = c.WithMountedDirectory("/app", m.Source)
	}

	if m.EnableXdebug {
		c = c.WithEnvVariable("XDEBUG_MODE", "coverage")
	}

	return c
}
