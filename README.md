# docker-credential-artifactory

Docker credential shim using the JFrog API to access Artifactory.

This internally calls the JFrog go client library to access the same stored
credentials that are used by the JFrog CLI.

## Instructions

1.  Install the `docker-credential-artifactory` binary somewhere on your `PATH`.
2.  Add the following entry to your `~/.docker/config.json`
    ```json
    {
      [...]
      "credHelpers": {
        "<server-id>.jfrog.io": "artifactory"
      }
    }
    ```
3.  Use `jf login` to get credentials to log into artifactory.
4.  You should now be able to run `docker pull <server.id>.jfrog.io/<repo>/<image>:<tag>`
