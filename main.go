package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/solerf/artix-mirrors/conf"
	"github.com/solerf/artix-mirrors/source"
)

type Artix struct{}
type Arch struct{}

type cmd struct {
	Artix          Artix `cmd:"" help:"generates Artix mirror file"`
	Arch           Arch  `cmd:"" help:"generates Arch mirror file"`
	AddHttp        bool  `help:"consider http mirrors" default:"false"`
	MirrorTimeout  int   `short:"x" help:"connection timeout (seconds) when checking mirrors" default:"3"`
	MaxConcurrency int   `short:"c" help:"max mirrors to check in parallel" default:"5"`
	TopMirrors     int   `short:"t" help:"top fastest mirrors to consider for output" default:"15"`
}

func (x *Artix) Run(c *conf.Settings) error {
	return source.Artix(c, os.Stdout)
}

func (h *Arch) Run(c *conf.Settings) error {
	return source.Arch(c, os.Stdout)
}

var cli = &cmd{}

func main() {
	kong.UsageOnError()
	kCtx := kong.Parse(cli, kong.Description("Generator of Artix and Arch linux mirrors for pacman"))

	c := conf.Default().
		WithAddHttp(cli.AddHttp).
		WithMaxConcurrency(cli.MaxConcurrency).
		WithMirrorTimeout(cli.MirrorTimeout).
		WithTopMirrors(cli.TopMirrors)

	if err := kCtx.Run(c); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
