package constants

const (
	SecretLabelKey      = "app.kubernetes.io/component"
	Aggregator          = "Aggregator"
	Agent               = "Agent"
	FilePath            = "/etc/vector/custom"
	SecretFinalizer     = "finalizers.kubesphere.io/secret-finalizer"
	SecretLabel         = "logging.whizard.io/vector-role"
	CALabel            =  "logging.whizard.io"
	ConfigReloadEnabled = "logging.whizard.io/enable"
	GeneratedFiles      = "generated-files"
	Config              = "config"
	Certification      = "certification"
)

var (
	FileDir    = ""
	VectorRole = ""
)
