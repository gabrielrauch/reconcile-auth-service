package ports

type TokenProvider interface {
	Generate(email string) (string, error)
	Validate(token string) bool
}
