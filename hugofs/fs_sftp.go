package hugofs

import (
	"errors"
	"github.com/pkg/sftp"
	"github.com/spf13/afero"
	"github.com/spf13/afero/sftpfs"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	"io/ioutil"
	"net/url"
	"os/user"
	"strings"
)

type SftpFsContext struct {
	sshc   *ssh.Client
	sshcfg *ssh.ClientConfig
	sftpc  *sftp.Client
}

func readPrivateKey(dir string) (auth ssh.AuthMethod, err error) {
	pemBytes, err := ioutil.ReadFile(dir + "/.ssh/id_rsa")
	if err != nil {
		err = errors.New("no password provided and couln't read private ssh key file: " + err.Error())
		return
	}

	signer, err := ssh.ParsePrivateKey(pemBytes)

	if err != nil {
		err = errors.New("no password provided and couln't parse private ssh key file: " + err.Error())
		return
	}
	auth = ssh.PublicKeys(signer)
	return
}

func SftpConnect(username, password, host string) (ctx *SftpFsContext, err error) {

	_user, err := user.Current()

	if username == "" {
		username = _user.Username
	}

	// initialize key database from "~/.ssl/knowh_hosts"
	hostkeyCB, err :=
		knownhosts.New(_user.HomeDir + "/.ssh/known_hosts")

	if err != nil {
		return
	}

	var auth ssh.AuthMethod

	if password == "" {
		auth, err = readPrivateKey(_user.HomeDir)
		if err != nil {
			return
		}
	} else {
		auth = ssh.Password(password)
	}

	sshcfg := &ssh.ClientConfig{
		User:            username,
		Auth:            []ssh.AuthMethod{auth},
		HostKeyCallback: hostkeyCB,
	}

	if !strings.Contains(host, ":") {
		host = host + ":22"
	}

	sshc, err := ssh.Dial("tcp", host, sshcfg)

	if err != nil {
		return nil, err
	}

	sftpc, err := sftp.NewClient(sshc)

	if err != nil {
		return nil, err
	}

	ctx = &SftpFsContext{
		sshc:   sshc,
		sshcfg: sshcfg,
		sftpc:  sftpc,
	}

	return ctx, nil
}

func (ctx *SftpFsContext) Disconnect() error {
	ctx.sftpc.Close()
	ctx.sshc.Close()
	return nil
}

func NewSftpFs(url *url.URL) (fs afero.Fs, err error) {
	var user, pwd string

	if url.User != nil {
		user = url.User.Username()
		pwd, _ = url.User.Password()
	} else {
		user = ""
		pwd = ""
	}

	ctx, err := SftpConnect(user, pwd, url.Host)
	if err != nil {
		return
	}
	// TODO: fix afero.sftp.Fs to
	// be able to call client Disconnect by
	// adding a Method to Fs or making
	// client a public field
	return sftpfs.New(ctx.sftpc), err
}
