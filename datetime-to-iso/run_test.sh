#!/bin/bash

set -eu

actual="tmp_test_output.txt"
expected="src/testdata.out.txt"
cargo run <src/testdata.txt >$actual
d=$(diff -U0 $expected $actual)

if [[ -z "$d" ]]; then
  >&2 echo "No diff, all good"
  exit 0
else
  >&2 echo "FAILED, diff!"
  echo $d
  exit 1
fi
