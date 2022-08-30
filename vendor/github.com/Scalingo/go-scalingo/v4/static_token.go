package scalingo

// StaticTokenGenerator is an implementation of TokenGenerator which always return the same token.
// This token is provided to the constructor. The TokenGenerator is used by the Client to
// authenticate to the Scalingo API.
//
// Usage:
//		t := GetStaticTokenGenerator("my-token")
//		client := NewClient(ClientConfig{
//			TokenGenerator: t,
//		})
//
// Any subsequent calls to the Scalingo API will use this token to authenticate.
type StaticTokenGenerator struct {
	token  string
	client *Client
}

// NewStaticTokenGenerator returns a new StaticTokenGenerator. The only argument is the token that
// will always be returned by this generator.
func NewStaticTokenGenerator(token string) *StaticTokenGenerator {
	return &StaticTokenGenerator{
		token: token,
	}
}

// GetAccessToken always returns the configured token.
func (t *StaticTokenGenerator) GetAccessToken() (string, error) {
	return t.token, nil
}

// SetClient sets the client attribute of this token generator.
func (t *StaticTokenGenerator) SetClient(c *Client) {
	t.client = c
}
