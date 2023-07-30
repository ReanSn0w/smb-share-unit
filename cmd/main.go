package main

import (
	"github.com/ReanSn0w/smb-share-unit/pkg/server"
	"github.com/ReanSn0w/smb-share-unit/pkg/utils"
	"github.com/go-pkgz/lgr"
)

var (
	logger utils.Logger = lgr.Default()
)

func main() {
	opts := utils.LoadParameters(logger)
	logger = opts.Logger()

	smb := utils.NewSMB(
		logger,
		opts.SMB.Host,
		opts.SMB.Port,
		opts.SMB.Sharename,
		opts.SMB.Username,
		opts.SMB.Password)

	cache := utils.NewCache(logger, opts.CacheTimeout)

	srv := server.New(logger, opts.HTTPPort, cache, smb)
	srv.Run()
}
