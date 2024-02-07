package constants

const (
	SecretLabelKey      = "app.kubernetes.io/component"
	Aggregator          = "Aggregator"
	Agent               = "Agent"
	FilePath            = "/etc/vector/custom"
	SecretFinalizer     = "finalizers.kubesphere.io/secret-finalizer"
	SecretLabel         = "logging.whizard.io/vector-role"
	ConfigReloadEnabled = "logging.whizard.io/enable"
	GeneratedFiles      = "generated-files"
	Config              = "config"
	Ccertification      = "certification"
)

var (
	FileDir    = ""
	VectorRole = ""
)
