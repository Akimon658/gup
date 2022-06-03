package config

func (c *Config) GetFlags(pkgName string) *BuildFlags {
	flags := c.Global

	for _, v := range c.Packages {
		if v.Name == pkgName {
			if v.Ldflags != "" {
				flags.Ldflags = v.Ldflags
			}
			if v.Tags != "" {
				flags.Tags = v.Tags
			}

			break
		}
	}

	return &flags
}
