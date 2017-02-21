package main

import (
	"fmt"
	"log"
	"os"

	"github.com/coreos/go-semver/semver"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	De "github.com/visionmedia/go-debug"
)

var debug = De.Debug("meshblu-connector-eradicator:main")

func main() {
	app := cli.NewApp()
	app.Name = "meshblu-connector-eradicator"
	app.Version = version()
	app.Action = run
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "dry-run",
			EnvVar: "MESHBLU_CONNECTOR_ERADICATOR_DRY_RUN",
			Usage:  "Print out what the connector would have done, without doing it",
		},
		cli.StringFlag{
			Name:   "local-app-data, l",
			Usage:  "Local AppData directory of the user.",
			EnvVar: "LOCALAPPDATA",
		},
	}
	app.Run(os.Args)
}

func run(context *cli.Context) {
	dryRun, localAppData := getOpts(context)
	fmt.Println("dryRun:", dryRun, "localAppData:", localAppData)
}

func getOpts(context *cli.Context) (bool, string) {
	dryRun := context.Bool("dry-run")
	localAppData := context.String("local-app-data")

	if localAppData == "" {
		cli.ShowAppHelp(context)

		if localAppData == "" {
			color.Red("  Missing required flag --local-app-data or LOCALAPPDATA")
		}
		os.Exit(1)
	}

	return dryRun, localAppData
}

func version() string {
	version, err := semver.NewVersion(VERSION)
	if err != nil {
		errorMessage := fmt.Sprintf("Error with version number: %v", VERSION)
		log.Panicln(errorMessage, err.Error())
	}
	return version.String()
}
