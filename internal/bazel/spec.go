package bazel

// See: https://github.com/EngFlow/credential-helper-spec/blob/main/schemas/get-credentials-request.schema.json
type GetCredentialsRequest struct {
	Uri string `json:"uri"`
}

// See: https://github.com/EngFlow/credential-helper-spec/blob/main/schemas/get-credentials-response.schema.json
type GetCredentialsReponse struct {
	Expires string              `json:"expires,omitempty"`
	Headers map[string][]string `json:"headers,omitempty"`
}

type BazelCredentialProvider interface {
	Get(request GetCredentialsRequest) (*GetCredentialsReponse, error)
}
