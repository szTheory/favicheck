# Favicheck

![Favicheck logo](logo.svg)

Find the web framework a site uses by checking its favicon against the [OWASP database of common favicons](https://wiki.owasp.org/index.php/OWASP_favicon_database)

```sh
$ favicheck https://static-labs.tryhackme.cloud/sites/favicon/images/favicon.ico
Web framework: cgiirc (0.5.9)
```

## Install

### Option 1: Download the binary

Go to the [releases page](https://github.com/szTheory/favicheck/releases) and get the archive for your OS and CPU combination. Extract it, then copy the `favicheck` binary to somewhere in your PATH.

### Option 2: Brew

```sh
brew tap szTheory/favicheck https://github.com/szTheory/favicheck
brew install favicheck
```

### Option 3: Go Get

```sh
go get github.com/szTheory/favicheck
```

## Usage

```sh
favicheck <filepath|url>
```

### Examples

```sh
# Check favicon from URL
favicheck https://static-labs.tryhackme.cloud/sites/favicon/images/favicon.ico
```

```sh
# Check favicon from file
favicheck ~/Downloads/favicon.ico
```

## Building the binary

```sh
make build
```

## Running the test suite

```sh
make test
```

## Linting the code

```sh
make lint
```

## Releasing

```sh
git tag -a v0.1.0 
git push origin v0.1.0
```
