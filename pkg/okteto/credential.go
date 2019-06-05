package okteto

import (
	"fmt"
)

// Credentials top body answer
type Credentials struct {
	Credentials Credential
}

//Credential represents an Okteto Space k8s credentials
type Credential struct {
	Server      string `json:"server" yaml:"server"`
	Certificate string `json:"certificate" yaml:"certificate"`
	Token       string `json:"token" yaml:"token"`
	Namespace   string `json:"namespace" yaml:"namespace"`
}

// GetCredentials returns the space config credentials
func GetCredentials(namespace string) (*Credential, error) {
	q := fmt.Sprintf(`query{
		credentials(space: "%s"){
			server, certificate, token, namespace
		},
	}`, namespace)

	var cred Credentials
	if err := query(q, &cred); err != nil {
		return nil, err
	}

	return &cred.Credentials, nil
}
