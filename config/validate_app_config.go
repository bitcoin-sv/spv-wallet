package config

// Validate checks the configuration for specific rules
func (c *AppConfig) Validate() error {
	var err error

	if err = c.Authentication.Validate(); err != nil {
		return err
	}

	if err = c.Cache.Validate(); err != nil {
		return err
	}

	if err = c.Db.Validate(); err != nil {
		return err
	}

	if err = c.Paymail.Validate(); err != nil {
		return err
	}

	if err = c.BHS.Validate(); err != nil {
		return err
	}

	if err = c.Server.Validate(); err != nil {
		return err
	}

	if err = c.ARC.Validate(); err != nil {
		return err
	}

	return nil
}
