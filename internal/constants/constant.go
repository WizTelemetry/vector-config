package constants

const (
	SecretLabelKey  = "app.kubernetes.io/component"
	Aggregator      = "Aggregator"
	Agent           = "Agent"
	FilePath        = "/etc/vector"
	SecretFinalizer = "finalizers.kubesphere.io/secret-finalizer"
	SecretLabel     = "vector.kubesphere.io/config"
)
