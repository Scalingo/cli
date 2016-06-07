package gitremote

func (config *Config) Update(r *Remote) error {
	err := config.Delete(r)
	if err != nil {
		return err
	}
	return config.Add(r)
}
