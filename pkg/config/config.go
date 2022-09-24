package config

type Config struct {
	ServiceDomain         string
	ServiceScheme         string
	ServicePort           int
	UnlockShares          int
	UnlockThreshold       int
	VaultRootTokenSecret  string
	VaultUnlockKeysSecret string
	Namespace             string
}
