// Main package for the mongoexport tool.
package main

import (
	"github.com/mongodb/mongo-tools/common/db"
	"github.com/mongodb/mongo-tools/common/log"
	"github.com/mongodb/mongo-tools/common/options"
	"github.com/mongodb/mongo-tools/common/signals"
	"github.com/mongodb/mongo-tools/common/util"
	"github.com/mongodb/mongo-tools/mongoexport"
	"os"
)

func main() {
	go signals.Handle()
	// initialize command-line opts
	opts := options.New("mongoexport", "<options>", options.EnabledOptions{Auth: true, Connection: true, Namespace: true})

	outputOpts := &mongoexport.OutputFormatOptions{}
	opts.AddOptions(outputOpts)
	inputOpts := &mongoexport.InputOptions{}
	opts.AddOptions(inputOpts)

	args, err := opts.Parse()
	if err != nil {
		log.Logf(log.Always, "error parsing command line options: %v", err)
		log.Logf(log.Always, "try 'mongoexport --help' for more information")
		os.Exit(util.ExitBadOptions)
	}
	if len(args) != 0 {
		log.Logf(log.Always, "too many positional arguments: %v", args)
		log.Logf(log.Always, "try 'mongoexport --help' for more information")
		os.Exit(util.ExitBadOptions)
	}

	log.SetVerbosity(opts.Verbosity)

	// print help, if specified
	if opts.PrintHelp(false) {
		return
	}

	// print version, if specified
	if opts.PrintVersion() {
		return
	}

	// connect directly, unless a replica set name is explicitly specified
	_, setName := util.ParseConnectionString(opts.Host)
	opts.Direct = (setName == "")
	opts.ReplicaSetName = setName

	provider, err := db.NewSessionProvider(*opts)
	if err != nil {
		log.Logf(log.Always, "error connecting to host: %v", err)
		os.Exit(util.ExitError)
	}
	exporter := mongoexport.MongoExport{
		ToolOptions:     *opts,
		OutputOpts:      outputOpts,
		InputOpts:       inputOpts,
		SessionProvider: provider,
	}

	err = exporter.ValidateSettings()
	if err != nil {
		log.Logf(log.Always, "error validating settings: %v", err)
		log.Logf(log.Always, "try 'mongoexport --help' for more information")
		os.Exit(util.ExitError)
	}

	writer, err := exporter.GetOutputWriter()
	if err != nil {
		log.Logf(log.Always, "error opening output stream: %v", err)
		os.Exit(util.ExitError)
	}
	if writer == nil {
		writer = os.Stdout
	} else {
		defer writer.Close()
	}

	numDocs, err := exporter.Export(writer)
	if err != nil {
		log.Logf(log.Always, "Failed: %v", err)
		os.Exit(util.ExitError)
	}

	if numDocs == 1 {
		log.Logf(log.Always, "exported %v record", numDocs)
	} else {
		log.Logf(log.Always, "exported %v records", numDocs)
	}

}
