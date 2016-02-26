// Package awscli - smaller aws toolkit
package awscli

// AwsCli - base struct
type AwsCli struct {
	version   int `json:"version"`
	keyID     string
	accessKey string
	token     string
	region    string
	account   string
}

// New - returns a new pointer to AwsCli
func New(version int, keyID, accessKey, token, region, account string) *AwsCli {
	return &AwsCli{
		version:   version,
		accessKey: accessKey,
		token:     token,
		region:    region,
		account:   account,
	}
}
