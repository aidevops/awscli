package awscli

type AwsCli struct {
  version int `json:"version"`
  key_id string
  access_key string
  token string
  region string
  account string
}

func New(version int, access_key, token, region, account string) *AwsCli {
  return &AwsCli{
    version: version,
    access_key: acess_key,
    token: token,
    region: region,
    account: account
  }
}
