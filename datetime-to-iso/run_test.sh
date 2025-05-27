#!/bin/bash

set -eu

input="src/testdata.txt"
actual="tmp_test_output.txt"
rm -fv $actual
expected="src/testdata.txt.expected"

cargo run <$input >$actual
if diff "$expected" "$actual"; then
  >&2 echo "all good"
else
  >&2 echo "FAILED. If this is the new status quo, then run"
  >&2 echo "cargo run <$input >$expected"
fi
