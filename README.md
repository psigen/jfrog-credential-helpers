# jfrog-credential-helpers

Credential shims using the JFrog CLI to access Artifactory.

Currently supports **bazel** and **docker** via the following helpers:

- `bazel-credential-jfrog`
- `docker-credential-jfrog`

This internally calls the JFrog go client library to access the same stored
credentials that are used by the JFrog CLI.

## Setting up the docker credential helper

1.  Install the `docker-credential-jfrog` binary somewhere on your `$PATH`.

    _(Note that you may need to `chmod +x` the file depending on how you fetched it)_

2.  Add the following entry to your `~/.docker/config.json`
    ```json
    {
      [...]
      "credHelpers": {
        "<server-id>.jfrog.io": "jfrog"
      }
    }
    ```

3.  Use `jf login` to get credentials to log into artifactory.

4.  You should now be able to run `docker pull <server.id>.jfrog.io/<repo>/<image>:<tag>`

## Setting up the bazel credential helper

1.  Install the `bazel-credential-jfrog` binary to a known location (or your `$PATH`)

    _(Note that you may need to `chmod +x` the file depending on how you fetched it)_

2.  Add the following to your `.bazelrc`:
    ```
    common --credential_helper=*.jfrog.io=/path/to/bazel-credential-jfrog
    ```

3.  Use `jf login` to get credentials to log into artifactory.

4.  Rules such as `http_archive` and `http_file` can now access artifactory URIs.

## Building manually

To manually build the executables in this package, use the following commands:

```bash
go build -o docker-credential-jfrog ./cmd/docker
go build -o bazel-credential-jfrog ./cmd/bazel
```

## References

- [EngFlow Credential Helper Spec](https://github.com/EngFlow/credential-helper-spec/blob/main/spec.md)
- [Bazel Credential Helper Proposal](https://github.com/bazelbuild/proposals/blob/main/designs/2022-06-07-bazel-credential-helpers.md)
- [`docker-credential-helpers` golang library](github.com/docker/docker-credential-helpers)
- [`jfrog-cli` golang library](https://github.com/jfrog/jfrog-cli)
