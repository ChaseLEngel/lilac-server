package transfer

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// Parses user's ssh keys and returns SSH public key authentication method.
func publicKey() (*ssh.AuthMethod, error) {
	authMethod := new(ssh.AuthMethod)

	// Get details of current user.
	curUser, err := user.Current()
	if err != nil {
		return authMethod, fmt.Errorf("Failed to get current user: %v", err)
	}

	// Read private key file. Assuming that private key is in default location.
	key, err := ioutil.ReadFile(curUser.HomeDir + "/.ssh/id_rsa")
	if err != nil {
		return authMethod, fmt.Errorf("Failed to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return authMethod, fmt.Errorf("Failed to parse private key: %v", err)
	}

	*authMethod = ssh.PublicKeys(signer)

	return authMethod, nil
}

// Connect to SSH server and initiate SFTP session.
func newSFTP(host, port string, config *ssh.ClientConfig) (*sftp.Client, error) {
	address := host + ":" + port

	sshClient, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect: %v", err)
	}

	sftp, err := sftp.NewClient(sshClient)
	if err != nil {
		return nil, fmt.Errorf("Failed to open SFTP session: %v", err)
	}

	return sftp, nil
}

// Transfers source file to destination on remote server.
func sendFile(src, dest string, client *sftp.Client) error {
	srcBytes, err := ioutil.ReadFile(src)
	if err != nil {
		return fmt.Errorf("Failed to read source file: %v", err)
	}

	// Create file on remote server.
	createdFile, err := client.Create(dest)
	if err != nil {
		return fmt.Errorf("Failed to create file %v on remote server: %v", dest, err)
	}

	// Write source file contents to remote file.
	if _, err := createdFile.Write(srcBytes); err != nil {
		return fmt.Errorf("Failed to write to remote file: %v", err)
	}

	return nil
}

// Recursively traverses directory and fills array with contents.
func directoryWalk(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		*files = append(*files, path)
		return err
	}
}

// Returns path with root being target.
// /home/user/files/file1 turns into files/file1
func relativePath(target, path string) string {
	split := strings.Split(path, "/")
	for i, s := range split {
		if s == target {
			return strings.Join(split[i:], "/")
		}
	}
	fmt.Println("WARN: Failed to ")
	return ""
}

// Recursively transfer directory to remote destination.
func sendDirectory(src, dest string, client *sftp.Client) error {
	var files []string

	// Get all files and directories from source directory.
	if err := filepath.Walk(src, directoryWalk(&files)); err != nil {
		return fmt.Errorf("Failed to read contents of source directory: %v", err)
	}

	// Write files and create directories on remote.
	for _, file := range files {
		// Define destination path on remote for file.
		path := filepath.Join(dest, relativePath(filepath.Base(src), file))
		info, err := os.Stat(file)
		if err != nil {
			return fmt.Errorf("Failed to stat file: %v", err)
		}
		if info.IsDir() {
			// Create nested directories in remote destination directory.
			if err = client.Mkdir(path); err != nil {
				return fmt.Errorf("Failed to make remote directory: %v", err)
			}
		} else {
			if err := sendFile(file, path, client); err != nil {
				return err
			}
		}
	}

	return nil
}

// Init authentication and transfer source file or directory to destination on remote server.
// Assumes SSH keys have been setup.
func Transfer(src, dest, host, port, user string) error {
	pubKey, err := publicKey()
	if err != nil {
		return err
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			*pubKey,
		},
	}

	client, err := newSFTP(host, port, config)
	if err != nil {
		return err
	}
	defer client.Close()

	info, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("Failed to get info of source: %v", err)
	}

	if info.IsDir() {
		err = sendDirectory(src, dest, client)
	} else {
		err = sendFile(src, filepath.Join(dest, filepath.Base(src)), client)
	}
	if err != nil {
		return err
	}

	return nil
}
