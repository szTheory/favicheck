#!/usr/bin/env bash

command="./bin/favicheck"
success=1

function test_stdout() {
  test "$1" "$2" "$3" "1"
}

function test_stderr() {
  test "$1" "$2" "$3" "2"
}

function test() {
  description=$1
  expected=$2
  param=$3
  output_stream=$4

  original_command="$command $param"
  eval_command="$original_command $output_stream>&1 | grep -q -e '$expected'"

  echo -n "===> $description: "
  if eval "$eval_command"; then
    echo "$(tput setaf 42)Pass$(tput sgr0)"
    return 1
  else
    echo "$(tput setaf 1)FAIL$(tput sgr0)"
    echo "======> Failed command was: $(tput setaf 1)$original_command$(tput sgr0)"
    success=0
    return 0
  fi
}

# Rebuild
./runners/build.sh

# Run tests
echo "Running tests..."

test_stdout "Shows usage when not enough parameters" \
  "Usage: favicheck <filepath|url>" \
  ""

test_stdout "Detects a matching favicon from file" \
  "Web framework: cgiirc (0.5.9)" \
  "fixtures/cgiirc.ico"

test_stdout "Shows when there is no match" \
  "No matching web framework for this favicon" \
  "fixtures/noframework.ico"

test_stdout "Detects a matching favicon from URL" \
  "Web framework: cgiirc (0.5.9)" \
  "https://static-labs.tryhackme.cloud/sites/favicon/images/favicon.ico"

test_stdout "Detects a non-matching favicon from URL" \
  "No matching web framework for this favicon" \
  "https://www.google.com/favicon.ico"

test_stderr "Errors when the URL is not a favicon" \
  "The URL is not a favicon" \
  "https://www.google.com/"

test_stderr "Errors when the file is not a favicon" \
  "The file is not a favicon" \
  "README.md"

test_stderr "Errors when the file doesn't exist on the filesystem" \
  "Could not open favicon file: fixtures/doesntexist.ico" \
  "fixtures/doesntexist.ico"

test_stderr "Errors when favicon at URL not found" \
  "Error while downloading favicon: HTTP status code 404" \
  "https://www.google.com/doesntexist.ico"

if [ $success -ne 1 ]; then
  # Test suite failure
  exit 1
fi
