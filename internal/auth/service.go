package auth

// Config holds the secrets and settings the auth package needs.
//
// Nothing outside cmd/ reads process configuration directly — cmd/server
// loads it once and constructs a Service with the values auth actually
// needs, instead of the auth package reaching into a global config.
type Config struct {
	// JWTSecretKey signs/verifies access tokens.
	JWTSecretKey []byte
	// JWTRefreshSecretKey signs/verifies refresh tokens.
	JWTRefreshSecretKey []byte
	// CookieSecure sets the Secure flag on auth cookies. Should be true
	// whenever the app is served over HTTPS (i.e. in production).
	CookieSecure bool
}

// Service issues and validates authentication sessions: JWTs, the cookies
// that carry them, and the auth.Principal derived from them.
type Service struct {
	cfg Config
}

// New constructs an auth Service from explicit configuration.
func New(cfg Config) *Service {
	return &Service{cfg: cfg}
}
