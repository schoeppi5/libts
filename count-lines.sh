#!/bin/bash
for line in $(find -name "*.go")
do
    lines=$(wc -l $line | awk '{print $1}')
    let total=$total+$lines
    if [[ $line =~ ^.*_test\.go$ ]]; then
        let tests=$tests+$lines
    else
        let code=$code+$lines
    fi
done

echo "Total lines of code: $total
    Lines of code:  $code
    Lines of tests: $tests"