#!/bin/bash
#
# This script is a pre-commit hook that checks if the mage command is
# installed and if not, prompts the user to install it. If mage is
# installed, the script navigates to the repository root and runs the
# `mage generatepackagedocs` command.
# This command generates documentation for all Go packages in the current
# directory and its subdirectories by traversing the file tree and creating
# a new README.md file in each directory containing a Go package.
# If the command fails, the commit is stopped, and an error message is shown.
set -e

cd "$(git rev-parse --show-toplevel)" || exit 1

# Check if mage is installed
if ! command -v mage > /dev/null 2>&1; then
    echo -e "mage is not installed\n"
    echo -e "Please install mage by running the following command:\n"
    echo -e "go install github.com/magefile/mage@latest\n"
    exit 1
fi

# Run the mage generatepackagedocs command
mage generatepackagedocs
# Catch the exit code
exit_status=$?

# If the exit code is not zero (i.e., the command failed),
# then stop the commit and show an error message
if [ $exit_status -ne 0 ]; then
    echo "failed to generate package docs"
    exit 1
fi
