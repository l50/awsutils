# awsutils

[![Go Report Card](https://goreportcard.com/badge/github.com/l50/awsutils)](https://goreportcard.com/report/github.com/l50/awsutils)
[![License](http://img.shields.io/:license-mit-blue.svg)](https://github.com/l50/awsutils/blob/master/LICENSE)
[![Tests](https://github.com/l50/awsutils/actions/workflows/tests.yaml/badge.svg)](https://github.com/l50/awsutils/actions/workflows/tests.yaml)
[![🚨 CodeQL Analysis](https://github.com/l50/awsutils/actions/workflows/codeql-analysis.yaml/badge.svg)](https://github.com/l50/awsutils/actions/workflows/codeql-analysis.yaml)
[![🚨 Semgrep Analysis](https://github.com/l50/awsutils/actions/workflows/semgrep.yaml/badge.svg)](https://github.com/l50/awsutils/actions/workflows/semgrep.yaml)
[![Coverage Status](https://coveralls.io/repos/github/l50/awsutils/badge.svg?branch=main)](https://coveralls.io/github/l50/awsutils?branch=main)

This repo is comprised of aws utilities that I use across various go projects.

## Dependencies

- [Install golang](https://go.dev/):

  ```bash
  gvm install go1.18
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

- [Optional - install gvm](https://github.com/moovweb/gvm):

  ```bash
  bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
  source "${GVM_BIN}"
  ```

---

## Developer Environment Setup

1. [Fork this project](https://docs.github.com/en/get-started/quickstart/fork-a-repo)

2. (Optional) If you installed gvm, create golang pkgset specifically for this project:

   ```bash
   mkdir "${HOME}/go"
   GVM_BIN="${HOME}/.gvm/scripts/gvm"
   export GOPATH="${HOME}/go"
   VERSION='1.18'
   PROJECT=awsutils

   bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
   source $GVM_BIN
   gvm install "go${VERSION}"
   gvm use "go${VERSION}"
   gvm pkgset create "${PROJECT}"
   gvm pkgset use "${PROJECT}"
   ```

3. Generate the `magefile` binary:

   ```bash
   mage -d .mage/ -compile ../magefile
   ```

4. Install pre-commit hooks and dependencies:

   ```bash
   ./magefile installPreCommitHooks
   ```

5. Update and run pre-commit hooks locally:

   ```bash
   ./magefile runPreCommit
   ```

6. Set up `go.mod` for development:

   ```bash
   ./magefile localGoMod
   ```

7. Fill out all TODO values in `test_env` files
   for code that you want to modify and test:

   ```bash
   vi pkg/ec2/test_env # fill in TODO values
   source pkg/ec2/test_env # export the env vars
   ```
