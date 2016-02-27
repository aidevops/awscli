// Package main -
package main

// Imports -
import (
	"os"

	"github.com/mitchellh/cli"

	"github.com/johnt337/awscli/command"
)

// Commands - map for sub commands
var Commands map[string]cli.CommandFactory

func init() {
	ui := &cli.BasicUi{Writer: os.Stdout}

	Commands = map[string]cli.CommandFactory{
		"ec2": func() (cli.Command, error) {
			return &command.EC2Command{
				UI: ui,
			}, nil
		},

		"ecr": func() (cli.Command, error) {
			return &command.ECRCommand{
				UI: ui,
			}, nil
		},

		"ecs": func() (cli.Command, error) {
			return &command.ECSCommand{
				UI: ui,
			}, nil
		},

		"version": func() (cli.Command, error) {
			return &command.VersionCommand{
				Revision:          GitCommit,
				Version:           Version,
				VersionPrerelease: VersionPrerelease,
				UI:                ui,
			}, nil
		},
	}
}
