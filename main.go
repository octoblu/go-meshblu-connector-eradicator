package main

import (
	"fmt"
	"log"
	"os"

	"github.com/coreos/go-semver/semver"
	"github.com/fatih/color"
	"github.com/octoblu/go-meshblu-connector-service/manage"
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

	uuids, err := manage.ListUserLogin(localAppData)
	if err != nil {
		log.Fatalln(err.Error())
	}

	if dryRun {
		printDryRun(uuids)
		os.Exit(0)
	}

	errChans := make([]chan error, len(uuids))

	for i, uuid := range uuids {
		errChans[i] = uninstall(localAppData, uuid)
	}

	allSuccessful := true

	for _, errChan := range errChans {
		err := <-errChan
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			allSuccessful = false
		}
	}

	if allSuccessful {
		os.Exit(0)
	}
	os.Exit(1)
}

func printDryRun(uuids []string) {
	for _, uuid := range uuids {
		fmt.Println("UNINSTALL: ", uuid)
	}
}

func uninstall(localAppData, uuid string) chan error {
	debug("UNINSTALL: %v", uuid)

	err := make(chan error)

	go func() {
		err <- manage.UninstallUserLogin(&manage.UninstallUserLoginOptions{
			LocalAppData: localAppData,
			UUID:         uuid,
		})
		debug("DONE: %v", uuid)
	}()

	return err
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
