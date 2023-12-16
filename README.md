# awsutils

[![License](https://img.shields.io/github/license/l50/awsutils?label=License&style=flat&color=blue&logo=github)](https://github.com/l50/awsutils/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/l50/awsutils)](https://goreportcard.com/report/github.com/l50/awsutils)
[![Tests](https://github.com/l50/awsutils/actions/workflows/tests.yaml/badge.svg)](https://github.com/l50/awsutils/actions/workflows/tests.yaml)
[![🚨 CodeQL Analysis](https://github.com/l50/awsutils/actions/workflows/codeql-analysis.yaml/badge.svg)](https://github.com/l50/awsutils/actions/workflows/codeql-analysis.yaml)
[![🚨 Semgrep Analysis](https://github.com/l50/awsutils/actions/workflows/semgrep.yaml/badge.svg)](https://github.com/l50/awsutils/actions/workflows/semgrep.yaml)
[![Coverage Status](https://coveralls.io/repos/github/l50/awsutils/badge.svg?branch=main)](https://coveralls.io/github/l50/awsutils?branch=main)
[![Renovate](https://github.com/l50/awsutils/actions/workflows/renovate.yaml/badge.svg)](https://github.com/l50/awsutils/actions/workflows/renovate.yaml)
[![Pre-Commit](https://github.com/l50/awsutils/actions/workflows/pre-commit.yaml/badge.svg)](https://github.com/l50/awsutils/actions/workflows/pre-commit.yaml)

This repo is comprised of aws utilities that I use across various go projects.

## Dependencies

- [Install pre-commit](https://pre-commit.com/):

  ```bash
  python3 -m pip install --upgrade pip
  python3 -m pip install pre-commit
  ```

- [Install Mage](https://magefile.org/):

  ```bash
  go install github.com/magefile/mage@latest
  ```

- [Install AWS CLI](https://aws.amazon.com/cli/)

---

## For Contributors and Developers

1. [Fork this project](https://docs.github.com/en/get-started/quickstart/fork-a-repo)

1. Install dependencies:

   ```bash
   mage installDeps
   ```

1. Update and run pre-commit hooks locally:

   ```bash
   mage runPreCommit
   ```

---

## Create New Release

- Download and install the [gh cli tool](https://cli.github.com/):

  - [macOS](https://github.com/cli/cli#macos)
  - [Linux](https://github.com/cli/cli/blob/trunk/docs/install_linux.md)
  - [Windows](https://github.com/cli/cli#windows)

- Install changelog extension:

  ```bash
  gh extension install chelnak/gh-changelog
  ```

- Generate changelog:

  ```bash
  NEXT_VERSION=v1.1.3
  gh changelog new --next-version "${NEXT_VERSION}"
  ```

- Create release:

  ```bash
  gh release create "${NEXT_VERSION}" -F CHANGELOG.md
  ```

---

## Developer Environment Setup

1. [Fork this project](https://docs.github.com/en/get-started/quickstart/fork-a-repo)

1. Update and run pre-commit hooks locally:

   ```bash
   mage runPreCommit
   ```

1. If running tests, create a `test_env` file from `test_env_template`
   and fill out all TODO values.

   Once all vars are filled out, export them:

   ```bash
   source test_env
   ```

   Alternatively, you can debug with vscode by
   removing the export statements in front of each
   variable.
