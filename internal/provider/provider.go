package provider

func New(version string) *TrinoProvider {
	return &TrinoProvider{
		version: version,
	}
}
