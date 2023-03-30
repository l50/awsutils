# awsutils

[![License](https://img.shields.io/github/license/l50/awsutils?label=License&style=flat&color=blue&logo=github)](https://github.com/l50/awsutils/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/l50/awsutils)](https://goreportcard.com/report/github.com/l50/awsutils)
[![Tests](https://github.com/l50/awsutils/actions/workflows/tests.yaml/badge.svg)](https://github.com/l50/awsutils/actions/workflows/tests.yaml)
[![🚨 CodeQL Analysis](https://github.com/l50/awsutils/actions/workflows/codeql-analysis.yaml/badge.svg)](https://github.com/l50/awsutils/actions/workflows/codeql-analysis.yaml)
[![🚨 Semgrep Analysis](https://github.com/l50/awsutils/actions/workflows/semgrep.yaml/badge.svg)](https://github.com/l50/awsutils/actions/workflows/semgrep.yaml)
[![Coverage Status](https://coveralls.io/repos/github/l50/awsutils/badge.svg?branch=main)](https://coveralls.io/github/l50/awsutils?branch=main)
[![Renovate](https://github.com/l50/awsutils/actions/workflows/renovate.yaml/badge.svg)](https://github.com/l50/awsutils/actions/workflows/renovate.yaml)

This repo is comprised of aws utilities that I use across various go projects.

## Dependencies

- [Install gvm](https://github.com/moovweb/gvm):

  ```bash
  bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
  ```

- [Install golang](https://go.dev/):

  ```bash
  source .gvm
  ```

- [Install pre-commit](https://pre-commit.com/):

  ```bash
  brew install pre-commit
  ```

- [Install Mage](https://magefile.org/):

  ```bash
  go install github.com/magefile/mage@latest
  ```

- [Install AWS CLI](https://aws.amazon.com/cli/)

---

## Developer Environment Setup

1. [Fork this project](https://docs.github.com/en/get-started/quickstart/fork-a-repo)

1. (Optional) If you installed gvm:

   ```bash
   source "${HOME}/.gvm"
   ```

1. Install pre-commit hooks and dependencies:

   ```bash
   mage installPreCommitHooks
   ```

1. Update and run pre-commit hooks locally:

   ```bash
   mage runPreCommit
   ```

1. Fill out all TODO values in `test_env` files
   for code that you want to modify and test:

   ```bash
   vi pkg/ec2/test_env # fill in TODO values
   source pkg/ec2/test_env # export the env vars
   ```

   Alternatively, you can debug with vscode by
   removing the export statements in front of each
   variable.
