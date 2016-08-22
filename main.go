package main

import (
	"log"
	"os"
	"runtime"

	"github.com/codegangsta/cli"
	"github.com/containerops/vessel/cmd"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	runtime.GOMAXPROCS(runtime.NumCPU())
	cmdWeb := cmd.GetCmdWeb()

	app := cli.NewApp()

	app.Name = "Vessel"
	app.Usage = "ContainerOps CI Service"
	app.Version = "v0.1.0-alpha.0"

	app.Commands = []cli.Command{
		cmdWeb,
	}

	app.Flags = append(app.Flags, []cli.Flag{}...)
	app.Run(os.Args)
}
