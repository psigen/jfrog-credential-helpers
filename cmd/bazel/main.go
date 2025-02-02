package main

import "github.com/psigen/jfrog-credential-helpers/internal/bazel"

func main() {
	bazel.Serve(ArtifactoryCredentialProvider{})
}
