package constants

// Index
const (
	ManifestExtension = ".yaml"
	DefaultIndexName  = "default"
	DefaultIndexURI   = "https://github.com/alex-held/devctl-index.git"
)

// devctl
const (
	DevctlPluginName = "devctl"
)

// environment variable keys
const (
	DEVCTL_ROOT_KEY              = "DEVCTL_ROOT"
	DEVCTL_ENV_KEY               = "DEVCTL_ENV"
	DEVCTL_DEFAULT_INDEX_URI_KEY = "DEVCTL_DEFAULT_INDEX_URI"
)

// paths
const (
	DefaultDevctlDir = ".devctl"
	ConfigDir        = "configs"
	IndexDir         = "index"
	SDKsDir          = "sdks"
	StoreDir         = "store"
	ReceiptsDir      = "receipts"
	BinDir           = "bin"
)
