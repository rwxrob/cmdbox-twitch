package cmd

import (
	"github.com/rwxrob/cmdtab"
	_ "github.com/rwxrob/cmdtab-config"
)

func init() {
	cmdtab.New("twitch", "mark", "config")
	//x.Default = "update"
}
