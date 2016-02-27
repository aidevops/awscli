// Package awscli - smaller aws toolkit
package awscli

// AwsCli - base struct
type AwsCli struct {
	Version   int    `json:"version"`
	KeyID     string `json:"key_id"`
	AccessKey string `json:"access_key"`
	Token     string `json:"token"`
	Region    string `json:"region"`
	Account   string `json:"account"`
}

// New - returns a new pointer to AwsCli
func New(version int, keyID, accessKey, token, region, account string) *AwsCli {
	return &AwsCli{
		Version:   version,
		AccessKey: accessKey,
		Token:     token,
		Region:    region,
		Account:   account,
	}
}
