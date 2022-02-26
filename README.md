# Favicheck

![Favicheck logo](logo.svg)

Find the web framework a site uses by checking its favicon against the [OWASP database of common favicons](https://wiki.owasp.org/index.php/OWASP_favicon_database)

```sh
$ favicheck https://static-labs.tryhackme.cloud/sites/favicon/images/favicon.ico
Web framework: cgiirc (0.5.9)
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
./runners/build.sh
```

## Running the test suite

```sh
./runners/test.sh
```
