package client

// ApplyOptions are options for applying objects to a cluster.
type ApplyOptions struct {
	Create bool
	GcTag  string
	SkipGc bool
	DryRun bool
	Client *Config
}

// DeleteOptions are options for deleting from a cluster.
type DeleteOptions struct {
	GracePeriod int64
	Client      *Config
}
