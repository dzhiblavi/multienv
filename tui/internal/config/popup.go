package config

import (
	"fmt"
	"strconv"
	"strings"
)

type Popup struct {
	Enabled bool   `yaml:"enabled"`
	Width   string `yaml:"width"`
	Height  string `yaml:"height"`
}

func (p Popup) validate() error {
	for name, value := range map[string]string{"width": p.Width, "height": p.Height} {
		percent := strings.HasSuffix(value, "%")
		number := strings.TrimSuffix(value, "%")
		if strings.Trim(number, "0123456789") != "" {
			return fmt.Errorf("popup.%s must be a positive cell count or a percentage from 1%% to 100%%", name)
		}

		size, err := strconv.Atoi(number)
		if err != nil || size < 1 || (percent && size > 100) {
			return fmt.Errorf("popup.%s must be a positive cell count or a percentage from 1%% to 100%%", name)
		}
	}

	return nil
}
