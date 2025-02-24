# Dagger PHP Module

This module is a PHP module for Dagger. It is primarily used to setup a PHP environment for your Dagger modules and pipelines.

## Usage

If you need to get into a php container to run commands, you can use the following command:

```bash
dagger call -m github.com/jasonmccallister call setup terminal
```

You can specify the version of PHP you would like to use by passing the `--version` flag:

```bash
dagger call -m github.com/jasonmccallister call --version=8.3 setup
```

Optionally, you can mount a source directory to the container by passing the `--source` flag:

```bash
dagger call -m github.com/jasonmccallister call --source=/path/to/source setup terminal
```

To use this in another Dagger module, you can run the `dagger install github.com/jasonmccallister/dagger-php` command in your module directory and it will add the module as a dependency in your `dagger.json` file.
