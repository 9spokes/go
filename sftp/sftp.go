package sftp

import (
	"bufio"
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

// Get retrieves a file from an SFTP server
func (sftp Context) Get(dir, file string) ([]string, error) {

	return sftp.runScript([]string{
		"cd " + dir,
		"get " + file,
		"rm " + file,
		"exit",
	})
}

// Put retrieves a file from an SFTP server
func (sftp Context) Put(dir, file string) ([]string, error) {

	return sftp.runScript([]string{
		"cd " + dir,
		"put " + file,
		"exit",
	})
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

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return output, fmt.Errorf("while reading from stderr: %s", err.Error())
	}

	if err := cmd.Start(); err != nil {
		return output, fmt.Errorf("while launching command: %s", err.Error())
	}

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		output = append(output, scanner.Text())
	}

	if err := cmd.Wait(); err != nil {
		return output, err
	}

	return output, nil
}
