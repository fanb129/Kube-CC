package ssh

import (
	"bufio"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

var (
	TypePassword = "password"
	TypeKey      = "key"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Type     string // password或者key
	Password string
	KeyPath  string // ssh id_rsa.id路径
}

type Ssh struct {
	Config Config
	client *ssh.Client
}

func NewSsh(c Config) (*Ssh, error) {
	// 创建ssh登录配置
	config := &ssh.ClientConfig{
		Timeout:         time.Second, // ssh连接time out时间一秒钟,如果ssh验证错误会在一秒钟返回
		User:            c.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 这个可以,但是不够安全
		//HostKeyCallback: hostKeyCallBackFunc(s.Host),
	}
	if c.Type == TypePassword {
		config.Auth = []ssh.AuthMethod{ssh.Password(c.Password)}
	} else if c.Type == TypeKey {
		config.Auth = []ssh.AuthMethod{publicKeyAuthFunc(c.KeyPath)}
	}

	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		zap.S().Errorln("创建ssh client失败：", err)
		return nil, err
	}
	return &Ssh{
		Config: c,
		client: sshClient,
	}, nil
}

func (s *Ssh) SendCmd(cmd string) (string, error) {
	session, err := s.client.NewSession()
	if err != nil {
		zap.S().Errorln(err)
		return "", err
	}
	defer session.Close()
	combinedOutput, err := session.CombinedOutput(cmd)
	if err != nil {
		zap.S().Errorln(err)
	}
	return string(combinedOutput), err
}

func (s *Ssh) CloseClient() error {
	err := s.client.Close()
	return err
}

func publicKeyAuthFunc(kPath string) ssh.AuthMethod {
	keyPath, err := homedir.Expand(kPath)
	if err != nil {
		log.Fatal("find key's home dir failed", err)
	}

	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatal("ssh key file read failed", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatal("ssh key signer failed", err)
	}
	return ssh.PublicKeys(signer)
}

func hostKeyCallBackFunc(host string) ssh.HostKeyCallback {
	hostPath, err := homedir.Expand("~/.ssh/known_hosts")
	if err != nil {
		log.Fatal("find known_hosts's home dir failed", err)
	}
	file, err := os.Open(hostPath)
	if err != nil {
		log.Fatal("can't find known_host file:", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hostKey ssh.PublicKey
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 {
			continue
		}
		if strings.Contains(fields[0], host) {
			var err error
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				log.Fatalf("error parsing %q: %v", fields[2], err)
			}
			break
		}
	}
	if hostKey == nil {
		log.Fatalf("no hostkey for %s,%v", host, err)
	}
	return ssh.FixedHostKey(hostKey)
}
