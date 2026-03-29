package source

import (
	"fmt"
	"io"

	"github.com/solerf/artix-mirrors/conf"
	"github.com/solerf/artix-mirrors/mirror"
)

// https://packages.artixlinux.org/mirrorlist/all/
// CHAOTIC AUR ==> https://gitlab.com/chaotic-aur/pkgbuilds/-/raw/main/chaotic-mirrorlist/mirrorlist

const artixJsonUrl = "https://status.artixlinux.org/mirrors/status/json/"

func Artix(c *conf.Settings, writeTo io.Writer) error {
	bb, err := fetchMirrorList(artixJsonUrl)
	if err != nil {
		return fmt.Errorf("artix : %w", err)
	}

	m, err := mirror.FromJson(c.AddHttp, bb)
	if err != nil {
		return fmt.Errorf("artix: %w", err)
	}

	m = m.Rank(c)
	return write(artixJsonUrl, m.Content(), writeTo)
}
