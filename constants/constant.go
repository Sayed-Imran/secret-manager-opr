package constants

import v1 "k8s.io/api/core/v1"


var SecretTypes = map[string]v1.SecretType{
	"generic":         v1.SecretTypeOpaque,
	"docker":          v1.SecretTypeDockerConfigJson,
	"docker-registry": v1.SecretTypeDockerConfigJson,
	"service-account": v1.SecretTypeServiceAccountToken,
	"tls":             v1.SecretTypeTLS,
	"basic-auth":      v1.SecretTypeBasicAuth,
	"ssh-auth":        v1.SecretTypeSSHAuth,
}
