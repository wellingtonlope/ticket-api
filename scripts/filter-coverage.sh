#!/bin/bash
# Script to filter coverage file based on .coverageignore patterns

COVERAGE_FILE="${1:-coverage.out}"
IGNORE_FILE="${2:-.coverageignore}"

if [ ! -f "$COVERAGE_FILE" ]; then
	echo "Error: Coverage file '$COVERAGE_FILE' not found"
	exit 1
fi

if [ ! -f "$IGNORE_FILE" ]; then
	echo "No .coverageignore file found, skipping filter"
	exit 0
fi

# Create temporary file
TEMP_FILE=$(mktemp)

# Copy coverage file to temp
cp "$COVERAGE_FILE" "$TEMP_FILE"

# Process each non-empty, non-comment line from ignore file
while IFS= read -r pattern || [ -n "$pattern" ]; do
	# Skip empty lines and comments
	pattern=$(echo "$pattern" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
	if [ -z "$pattern" ] || [[ "$pattern" == \#* ]]; then
		continue
	fi

	# Filter out lines matching the pattern
	grep -v "$pattern" "$TEMP_FILE" >"${TEMP_FILE}.tmp"
	mv "${TEMP_FILE}.tmp" "$TEMP_FILE"
done <"$IGNORE_FILE"

# Move filtered file back
mv "$TEMP_FILE" "$COVERAGE_FILE"

echo "Coverage filtered successfully"
