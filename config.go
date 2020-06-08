package staticfile

type Config struct {
	Nginx *Nginx `yaml:"nginx,omitempty"`
}

type Nginx struct {
	RootDir               string            `yaml:"root"`
	HostDotFiles          bool              `yaml:"host_dot_files"`
	LocationInclude       string            `yaml:"location_include"`
	DirectoryIndex        bool              `yaml:"directory"`
	SSI                   bool              `yaml:"ssi"`
	PushState             bool              `yaml:"pushstate"`
	HSTS                  bool              `yaml:"http_strict_transport_security"`
	HSTSIncludeSubDomains bool              `yaml:"http_strict_transport_security_include_subdomains"`
	HSTSPreload           bool              `yaml:"http_strict_transport_security_preload"`
	ForceHTTPS            bool              `yaml:"force_https"`
	BasicAuth             bool              `yaml:"basic_auth"`
	StatusCodes           map[string]string `yaml:"status_codes"`
}
