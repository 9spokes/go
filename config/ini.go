package config

import "github.com/go-ini/ini"

//Read reads an .ini-style file and returns a *ini.File object
func Read(filename string) (*ini.File, error) {
	cfg, err := ini.Load(filename)

	if err != nil {
		return nil, err
	}

	return cfg, nil
}
