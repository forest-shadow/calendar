#!/bin/bash
set -e

# Create pre-commit hook
cat > .git/hooks/pre-commit << 'EOL'
#!/bin/bash
set -e

# Initialize asdf if it's installed
if command -v asdf >/dev/null 2>&1; then
  . "$HOME/.asdf/asdf.sh"
else
  echo "asdf not found; skipping asdf setup."
fi

if command -v golangci-lint >/dev/null 2>&1; then
  echo "Running golangci-lint"
  golangci-lint run --fix

  # Check if there are any changes after linting
  if ! git diff --exit-code; then
      echo "Linter made changes.\nAdding staged files changes to the commit."
      git add -u $(git diff --name-only --cached)
  else
      echo "No changes made by the linter."
  fi
else
  echo "golangci-lint not found; skipping linting."
fi
EOL

# Make the pre-commit hook executable
chmod +x .git/hooks/pre-commit

echo "Pre-commit hook has been set up successfully."
