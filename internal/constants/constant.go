package constants

const (
	SecretLabelKey      = "app.kubernetes.io/component"
	Aggregator          = "Aggregator"
	Agent               = "Agent"
	FilePath            = "/etc/vector"
	SecretFinalizer     = "finalizers.kubesphere.io/secret-finalizer"
	SecretLabel         = "logging.whizard.io/vector-role"
	ConfigReloadEnabled = "logging.whizard.io/enable"
	GeneratedFiles      = "generated-files"
)

var (
	FileDir    = ""
	VectorRole = ""
)
