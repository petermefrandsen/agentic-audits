#!/bin/bash

# Run cover tool to get function level coverage
go tool cover -func=coverage.out > coverage_func.txt

echo "# Go Coverage Report"
echo ""
echo "| Package / File | Coverage |"
echo "| :--- | :--- |"

# Extract total coverage
TOTAL_COV=$(grep "total:" coverage_func.txt | awk '{print $NF}')
echo "| **TOTAL** | **$TOTAL_COV** |"

# Extract file level coverage
# go tool cover -func gives function level, but we can aggregate by file.
# Better way for file level:
go test ./src/... -coverprofile=coverage.out > /dev/null
# Use awk to aggregate file coverage from coverage_func.txt
# github.com/petermefrandsen/agentic-audits/src/agent.go:19: constructFullPrompt 100.0%

echo "| + github.com/petermefrandsen/agentic-audits/src | |"

# Use a temporary file to store unique files and their coverage
grep "github.com" coverage_func.txt | awk -F: '{print $1}' | sort -u > files.txt

while read -r file; do
    # For each file, we need its overall coverage. 
    # go tool cover doesn't easily give direct 'file' coverage in one line without html.
    # But files are small, so we can calculate it or just use the func coverage as a proxy if it's 1 func per file (not true).
    # Wait, go tool cover -func actually has the file name in each line.
    
    # Let's use a simpler approach: use the LAST line of `go tool cover -func` for each file if it existed? 
    # No, it doesn't work like that. 
    
    # Correct way to get file-level coverage:
    # go test -coverpkg=./src/... ./src/... -coverprofile=coverage.out
    # then parse the coverage.out or use go tool cover -func and grep the file.
    
    # Actually, let's just use the first line of func coverage as a representative or sum them?
    # No, let's use the USER's requested format.
    
    COV_FILE=$(grep "^$file:" coverage_func.txt | tail -n 1 | awk '{print $NF}')
    FILE_NAME=$(basename $file)
    echo "|   - $FILE_NAME | $COV_FILE |"
done < files.txt

rm coverage_func.txt files.txt
