package main

import (
	"flag"
	"github.com/pckhoi/tent/internal/app/loaders/postgres"
	"github.com/pckhoi/tent/internal/app/settings"
)

func main() {
	rootDir := flag.String("rootdir", "", "Specify where the data should be stored")
	flag.Parse()
	settings.Set(settings.Settings{
		RootDir: *rootDir,
	})
	postgres.Load(flag.Arg(0))
}
