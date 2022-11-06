package main

import (
	"fmt"
	"os"
	"time"

	log "github.com/go-pkgz/lgr"
	"github.com/nazandr/secret-pass/server"
	"github.com/umputun/go-flags"
)

var revision string

var opts struct {
	Port     string        `short:"p" long:"port" env:"PORT" default:":8080" description:"Service port"`
	Lifespan time.Duration `short:"ls" long:"lifespan" env:"LIFESPAN" description:"Secret lifespan"`
	Debug    bool          `long:"dbg" env:"DEBUG" description:"Enable debug mode"`
}

func main() {
	fmt.Printf("Secret-pass %s\n", revision)

	p := flags.NewParser(&opts, flags.PassDoubleDash|flags.HelpFlag)
	if _, err := p.Parse(); err != nil {
		if err.(*flags.Error).Type != flags.ErrHelp {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		p.WriteHelp(os.Stderr)
		os.Exit(2)
	}

	setupLog(opts.Debug)

	server := server.NewServer(opts.Port, opts.Lifespan)
	if err := server.Run(); err != nil {
		log.Fatalf("[ERROR] cannot start server: %s", err)
	}
}

func setupLog(dbg bool) {
	logOpts := []log.Option{log.Msec, log.LevelBraces, log.StackTraceOnError}
	if dbg {
		logOpts = append(logOpts, log.Debug, log.CallerFunc)
	}
	log.SetupStdLogger(logOpts...)
	log.Setup(logOpts...)
}
