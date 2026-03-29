package conf

import "time"

type Settings struct {
	AddHttp        bool
	MirrorTimeout  time.Duration
	MaxConcurrency int
	TopMirrors     int
}

func Default() *Settings {
	return &Settings{}
}

func (c *Settings) WithAddHttp(addHttp bool) *Settings {
	c.AddHttp = addHttp
	return c
}

func (c *Settings) WithMirrorTimeout(timeoutSeconds int) *Settings {
	c.MirrorTimeout = time.Duration(timeoutSeconds) * time.Second
	return c
}

func (c *Settings) WithMaxConcurrency(maxConcurrency int) *Settings {
	c.MaxConcurrency = maxConcurrency
	return c
}

func (c *Settings) WithTopMirrors(top int) *Settings {
	c.TopMirrors = top
	return c
}
