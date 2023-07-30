package utils

import (
	"os"

	"github.com/go-pkgz/lgr"
	"github.com/umputun/go-flags"
)

func LoadParameters(log Logger) *Parameters {
	opts := &Parameters{}

	p := flags.NewParser(opts, flags.PrintErrors|flags.PassDoubleDash|flags.HelpFlag)
	p.SubcommandsOptional = true
	if _, err := p.Parse(); err != nil {
		if err.(*flags.Error).Type != flags.ErrHelp {
			log.Logf("[ERROR] cli error: %v", err)
		}
		os.Exit(2)
	}

	log.Logf("[DEBUG] Verbose: %v", opts.Verbose)
	log.Logf("[DEBUG] HTTPPort: %v", opts.HTTPPort)
	log.Logf("[DEBUG] Cache Timeout: %v", opts.CacheTimeout)
	log.Logf("[DEBUG] SMB Host: %v", opts.SMB.Host)
	log.Logf("[DEBUG] SMB Port: %v", opts.SMB.Port)
	log.Logf("[DEBUG] SMB Sharename: %v", opts.SMB.Sharename)

	return opts
}

type Parameters struct {
	Verbose      bool `short:"v" long:"verbose" env:"VERBOSE" description:"verbose mode"`
	HTTPPort     int  `short:"p" long:"port" env:"PORT" default:"8080" description:"http port for website"`
	CacheTimeout int  `short:"t" long:"timeout" env:"TIMEOUT" default:"3600" description:"cache timeout in seconds"`

	SMB struct {
		Host      string `long:"host" env:"HOST" default:"localhost" description:"smb host"`
		Port      int    `long:"port" env:"PORT" default:"445" description:"smb port"`
		Sharename string `long:"sharename" env:"SHARENAME" description:"smb sharename"`
		Username  string `long:"username" env:"USERNAME" default:"anonymous" description:"smb username"`
		Password  string `long:"password" env:"PASSWORD" default:"" description:"smb password"`
	} `group:"smb" namespace:"smb" env-namespace:"SMB"`
}

func (p *Parameters) Logger() Logger {
	if p.Verbose {
		return lgr.New(lgr.Msec, lgr.Debug, lgr.CallerFile, lgr.CallerFunc)
	}

	return lgr.Default()
}
