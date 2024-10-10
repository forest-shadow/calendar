#!/bin/bash
set -e

# Create pre-commit hook
cat > .git/hooks/pre-commit << EOL
#!/bin/bash
set -e

echo "Running golangci-lint"
golangci-lint run --fix

# Check if there are any changes after linting
if ! git diff --exit-code; then
    echo "Linter made changes. Adding them to the commit."
    git add .
else
    echo "No changes made by the linter."
fi
EOL

# Make the pre-commit hook executable
chmod +x .git/hooks/pre-commit

echo "Pre-commit hook has been set up successfully."
