package utils

import (
	"fmt"
	"net"

	"github.com/hirochachacha/go-smb2"
)

// NewSMB returns a new SMB client.
func NewSMB(log Logger, host string, port int, share, username, password string) SMB {
	return &smb{
		log:      log,
		host:     host,
		port:     port,
		share:    share,
		username: username,
		password: password,
	}
}

type smb struct {
	log      Logger
	host     string
	port     int
	share    string
	username string
	password string
}

// Get returns the file contents from the SMB share.
func (s *smb) Get(filename string) ([]byte, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", s.host, s.port))
	if err != nil {
		s.log.Logf("[ERROR] tcp connection error: %v", err)
		return nil, err
	}
	defer conn.Close()

	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     s.username,
			Password: s.password,
		},
	}

	smb, err := d.Dial(conn)
	if err != nil {
		s.log.Logf("[ERROR] smb connection error: %v", err)
		return nil, err
	}
	defer smb.Logoff()

	fs, err := smb.Mount(s.share)
	if err != nil {
		s.log.Logf("[ERROR] smb mount error: %v", err)
		return nil, err
	}
	defer fs.Umount()

	return fs.ReadFile(filename)
}
