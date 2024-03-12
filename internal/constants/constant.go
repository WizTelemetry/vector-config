package constants

const (
	SecretLabelKey      = "app.kubernetes.io/component"
	Aggregator          = "Aggregator"
	Agent               = "Agent"
	FilePath            = "/etc/vector/custom"
	SecretFinalizer     = "finalizers.kubesphere.io/secret-finalizer"
	SecretLabel         = "logging.whizard.io/vector-role"
	CertificationLabel  = "logging.whizard.io/certification"
	ConfigReloadEnabled = "logging.whizard.io/enable"
	GeneratedFiles      = "generated-files"
	Certification       = "certification"
)

var (
	FileDir    = ""
	VectorRole = ""
)
