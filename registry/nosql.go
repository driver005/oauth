package registry

func (m *RegistrySQL) CanHandle(dsn string) bool {
	return m.alwaysCanHandle(dsn)
}
