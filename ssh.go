package botil

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"time"
)
type SSHConnecter struct{
	user string
	pass string
	host string
	port int
	Client *ssh.Client
}

func NewSSHConnecter(user, password, host string, port int)*SSHConnecter{
	connecter:=&SSHConnecter{}
	connecter.user=user
	connecter.pass=password
	connecter.host=host
	connecter.port=port
	return connecter;
}

func (self *SSHConnecter)Connect() (error) {
	if(self.Client!=nil) {
		return nil
	}
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(self.pass))

	clientConfig = &ssh.ClientConfig{
		User:    self.user,
		Auth:    auth,
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", self.host, self.port)

	if self.Client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return err
	}

	return nil
}
type SFTPConnecter struct{
	Client *sftp.Client
	SSHconnecter *SSHConnecter
}
func NewSFTPConnecter(user, password, host string, port int)*SFTPConnecter{
	connecter:=&SFTPConnecter{}
	connecter.SSHconnecter=NewSSHConnecter(user, password, host, port)
	return connecter;
}
func NewSFTPConnecterWithSSH(SSHconnecter *SSHConnecter)*SFTPConnecter{
	connecter:=&SFTPConnecter{}
	connecter.SSHconnecter=SSHconnecter
	return connecter;
}
func (self *SFTPConnecter)Connect() (error) {
	if(self.Client!=nil){
		return nil
	}
	err:=self.SSHconnecter.Connect()
	if  err!=nil{
		return err
	}
	// create sftp client
	if self.Client, err= sftp.NewClient(self.SSHconnecter.Client); err != nil {
		return err
	}
	return nil
}

func UploadFile(client *sftp.Client,srcFilePath,desDir string)error{
	srcFilePath = filepath.Clean(srcFilePath)
	desDir = filepath.ToSlash(filepath.Clean(desDir))

	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	_, err = client.Stat(desDir)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if(err != nil&&os.IsNotExist(err)) {
		err = client.MkdirAll(desDir)
		if err != nil {
			return err
		}
	}

	remoteFileName:= filepath.Base(srcFilePath)
	remoteFilePath:=filepath.ToSlash(filepath.Join(desDir, remoteFileName))
	dstFile, err := client.Create(remoteFilePath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	buf := make([]byte, 1024*1024)
	readsize:=int64(0)
	state, err := srcFile.Stat()
	if (err != nil) {
		return err
	}
	end:=make(chan bool,1)
	go func(){
		for  {
			select {
			case <-time.After(time.Second):
				percent:=float64(readsize)/float64(state.Size())*100
				fmt.Printf("\r%s %d/%d  %.2f%%", state.Name(), readsize, state.Size(),percent )
				if int64(readsize) == state.Size() {
					fmt.Println("")
					end<-true
					return
				}
			}
		}
	}()
	for {
		n, _ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		readsize+=int64(n)
		dstFile.Write(buf[0:n])
	}
	readsize=state.Size()
	<-end
	return nil
}

func UploadDir(client *sftp.Client,srcDir,desDir string) (err error){
	srcDir = filepath.Clean(srcDir)
	desDir = filepath.ToSlash(filepath.Clean(desDir))
	fmt.Printf("folder %s to %s \n",srcDir,desDir)

	si, err := os.Stat(srcDir)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}


	_, err = client.Stat(desDir)
	if err != nil && !os.IsNotExist(err) {
		return
	}

	if(err != nil&&os.IsNotExist(err)) {
		err = client.MkdirAll(desDir)
		if err != nil {
			return
		}
	}

	entries, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		dstPath := filepath.Join(desDir, entry.Name())

		if entry.IsDir() {
			err = UploadDir(client,srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = UploadFile(client,srcPath, desDir)
			if err != nil {
				return
			}
		}
	}
	return nil
}