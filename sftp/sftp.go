package sftp

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// Context contains the details of the remote file to be fetched
type Context struct {
	Hostname   string
	Username   string
	PrivateKey string
}

// TransferOptions contains a set of options that customise a file operation
type TransferOptions struct {
	RemoteDirectory  string
	LocalDirectory   string
	ArchiveDirectory string
	RemoveOnSuccess  bool
	FilePattern      string
}

// List gets a list of remote files in a directory
func (sftp Context) List(opt TransferOptions) ([]string, error) {

	script := make([]string, 0)

	if opt.RemoteDirectory != "" {
		script = append(script, "cd "+opt.RemoteDirectory)
	}

	script = append(script, "ls -1 "+opt.FilePattern)
	script = append(script, "exit")

	return sftp.runScript(script)
}

// Get retrieves a file from an SFTP server
func (sftp Context) Get(file string, opt TransferOptions) ([]string, error) {

	script := make([]string, 0)

	if opt.RemoteDirectory != "" {
		script = append(script, "cd "+opt.RemoteDirectory)
	}

	if opt.LocalDirectory != "" {
		script = append(script, "lcd "+opt.LocalDirectory)
	}

	script = append(script, "get "+file)

	if opt.ArchiveDirectory != "" {
		script = append(script, "mv "+file+" "+opt.ArchiveDirectory)
	} else if opt.RemoveOnSuccess {
		script = append(script, "rm "+file)
	}

	script = append(script, "exit")

	return sftp.runScript(script)
}

// Put retrieves a file from an SFTP server
func (sftp Context) Put(file string, opt TransferOptions) ([]string, error) {

	script := make([]string, 0)

	if opt.RemoteDirectory != "" {
		script = append(script, "cd "+opt.RemoteDirectory)
	}

	script = append(script, "put "+file)
	script = append(script, "exit")

	return sftp.runScript(script)
}

func (sftp Context) runScript(script []string) ([]string, error) {

	output := make([]string, 0)

	content := []byte(strings.Join(script, "\n"))
	tmpfile, err := ioutil.TempFile("", "_sftp")
	if err != nil {
		return output, fmt.Errorf("while creating temporary SFTP file: %s", err.Error())
	}

	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(content); err != nil {
		return output, fmt.Errorf("while writing temporary SFTP script: %s", err.Error())
	}

	if err := tmpfile.Close(); err != nil {
		return output, fmt.Errorf("while close temporary SFTP script handle: %s", err.Error())
	}

	cmd := exec.Command("sftp", "-q", "-oStrictHostKeyChecking=no", "-oUserKnownHostsFile=/dev/null", "-i", sftp.PrivateKey, "-b", tmpfile.Name(), sftp.Username+"@"+sftp.Hostname)

	out, err := cmd.CombinedOutput()

	if err != nil {
		return strings.Split(string(out), "\n"), err
	}
	return sftp.sanitiseOutput(out)
}

func (sftp Context) sanitiseOutput(output []byte) ([]string, error) {

	raw := bytes.Split(output, []byte{'\n'})
	lines := make([]string, 0)

	for i := range raw {
		line := string(raw[i])
		if len(line) > 0 && !strings.Contains(line, "sftp> ") {
			lines = append(lines, strings.Replace(line, " ", "", -1))
		}
	}

	return lines, nil
}
