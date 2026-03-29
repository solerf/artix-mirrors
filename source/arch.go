package source

import (
	"fmt"
	"io"

	"github.com/solerf/artix-mirrors/conf"
	"github.com/solerf/artix-mirrors/mirror"
)

// https://archlinux.org/mirrors/status/tier/1/json/
// https://packages.artixlinux.org/mirrorlist/all/

const archJsonUrl = "https://archlinux.org/mirrors/status/tier/1/json/"

func Arch(c *conf.Settings, writeTo io.Writer) error {
	bb, err := fetchMirrorList(archJsonUrl)
	if err != nil {
		return fmt.Errorf("arch : %w", err)
	}

	m, err := mirror.FromJson(c.AddHttp, bb)
	if err != nil {
		return fmt.Errorf("arch: %w", err)
	}

	m = m.Rank(c)
	return write(archJsonUrl, m.Content(), writeTo)
}
