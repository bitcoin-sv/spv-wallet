package config

// Validate checks the configuration for specific rules
func (a *AppConfig) Validate() error {
	var err error

	if err = a.Authentication.Validate(); err != nil {
		return err
	}

	if err = a.Cache.Validate(); err != nil {
		return err
	}

	if err = a.Db.Validate(); err != nil {
		return err
	}

	if err = a.Paymail.Validate(); err != nil {
		return err
	}

	if err = a.BHS.Validate(); err != nil {
		return err
	}

	if err = a.Server.Validate(); err != nil {
		return err
	}

	if err = a.ARC.Validate(); err != nil {
		return err
	}

	return nil
}
