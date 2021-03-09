package main

import (
	"github.com/rwxrob/cmdtab"
	_ "github.com/rwxrob/cmdtab-twitch"
)

func main() {
	cmdtab.Execute("twitch")
}
