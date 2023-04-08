package args

import (
	"os"

	"github.com/jessevdk/go-flags"
)

type Options struct {
	Help         bool   `long:"help" description:"Show Usage Information"`
	Program_Name string `long:"program_name" description:"Name of the program where you got a duplicate on" required:"true"`
	H1_Session   string `long:"h1_session" description:"HackerOne Account Session Cookie" required:"true"`
	Report_ID    string `long:"report_id" description:"HackerOne Report ID" required:"true"`
}

func GetArgs() Options {
	var opts Options
	parser := flags.NewParser(&opts, flags.Default)

	// Parse command-line flags
	_, err := parser.Parse()
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}

	return opts
}
