package config

type ServiceEndpoints struct {
	UserGraphQL string
}

type Conf struct {
	Address                    string
	FirebaseCredentialFilePath string
	Port                       int
	ProjectID                  string
	PubSubSubscribeMember      string
	PubSubTopicMember          string
	ServiceEndpoints           ServiceEndpoints
	TokenSecretName            string
	V0RESTfulSrvTargetURL      string
}

func (c *Conf) Valid() bool {
	return true
}
