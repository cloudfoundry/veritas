package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"

	"code.cloudfoundry.org/bbs"
	"code.cloudfoundry.org/cfdot/commands"
	"code.cloudfoundry.org/lager"
)

func main() {
	bbsParser := flags.NewParser(&commands.BBSOptions, flags.IgnoreUnknown|flags.PassDoubleDash)
	// ignoring error since we catch on the main parser below
	bbsParser.Parse()
	var bbsClient bbs.Client
	if commands.BBSOptions.BBSURL != "" {
		bbsClient = bbs.NewClient(commands.BBSOptions.BBSURL)
	}

	parser := flags.NewParser(&commands.CFdot, flags.HelpFlag|flags.PassDoubleDash)
	logger := lager.NewLogger("cfdot")
	commands.Configure(logger, os.Stdout, bbsClient)

	retargs, err := parser.Parse()
	if err != nil {
		if err == commands.ErrShowHelpMessage || (len(retargs) == 1 && retargs[0] == "") {
			command := os.Args[0]
			writeSynopsis(command)

			helpParser := flags.NewParser(&commands.CFdot, flags.IgnoreUnknown|flags.HelpFlag)
			helpParser.NamespaceDelimiter = "-"
			helpParser.ParseArgs([]string{"-h"})
			helpParser.Command.Name = command

			helpParser.WriteHelp(os.Stdout)
			os.Exit(0)
		} else {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			os.Exit(1)
		}
	}
}

func writeSynopsis(command string) {
	fmt.Fprintln(os.Stdout, "SYNOPSIS:")
	fmt.Fprintf(os.Stdout, "  %s: A command-line tool to interact with a Cloud Foundry Diego deployment.\n\n", command)
}
