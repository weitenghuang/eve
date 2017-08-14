package aws

// AuthType represents an integer value for
// AuthType constants
type AuthType int

const (
	IAM        AuthType = iota
	AssumeRole AuthType = iota
)

// Account represents and AWS Account
type Account struct {
	Id       int64
	Name     string
	Roles    []string
	Regions  []string
	AuthType AuthType
	Meta     map[string]string
}
