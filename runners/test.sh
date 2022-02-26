#!/usr/bin/env bash

command="./bin/favicheck"

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

test_stderr "Errors on missing files" \
  "panic: Could not open favicon file: fixtures/NOTEXIST.ico" \
  "fixtures/NOTEXIST.ico"

test_stdout "Detects a matching favicon from file" \
  "Web framework: cgiirc (0.5.9)" \
  "fixtures/cgiirc.ico"

# Non-matching framework from file
test_stdout "Shows when there is no match" \
  "No matching web framework for this favicon" \
  "fixtures/noframework.ico"

# Matching framework from URL
test_stdout "Detects a matching favicon from file" \
  "Web framework: cgiirc (0.5.9)" \
  "https://static-labs.tryhackme.cloud/sites/favicon/images/favicon.ico"

# Non-matching framework from URL

# Neither a file nor a URL

# URL is not a favicon

# Handles the case of Zero byte favicons
