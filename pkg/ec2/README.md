# Running the Tests

You'll need a bunch of environment variables to run these tests.
The easiest way to do this is to copy and fill in the provided template:

```bash
cp test_env_template test_env
```

Fill in all of the `export` commands in `test_env` with the appropriate values,
then source the file to add the environment variables to your shell:

```bash
source test_env
```

**Note:** `test_env` is captured in the repo's `.gitignore`,
so you don't need to worry about accidentally committing secrets.
