package service

import (
	"Kube-CC/common/responses"
	"Kube-CC/service/ssh"
	"go.uber.org/zap"
)

// CreateNodeBySealos  添加node
func CreateNodeBySealos(passwd string, hosts []string) (*responses.Response, error) {
	newSsh, err := ssh.NewSsh(ssh.MasterConfig)
	defer newSsh.CloseClient()
	if err != nil {
		zap.S().Errorln(err)
		return nil, err
	}
	nodes := ""
	i := 0
	for i = 0; i < len(hosts)-1; i++ {
		nodes += hosts[i] + ","
	}
	nodes += hosts[i]

	// 新版sealos 支持输入ssh密码
	//cmd := "echo 'y' | sealos add --nodes " + nodes + " --passwd " + passwd
	// 旧版sealos
	cmd := "echo 'y' | sealos add --nodes " + nodes

	zap.S().Info(cmd)
	r, err := newSsh.SendCmd(cmd)
	if err != nil {
		zap.S().Errorln(err)
		return nil, err
	}
	zap.S().Info(r)
	return &responses.OK, nil
}

// DeleteNodeBysealos  删除node节点
func DeleteNodeBysealos(name string) (*responses.Response, error) {
	newSsh, err := ssh.NewSsh(ssh.MasterConfig)
	defer newSsh.CloseClient()
	if err != nil {
		zap.S().Errorln(err)
		return nil, err
	}
	cmd := "echo 'y' | sealos delete --nodes " + name
	zap.S().Info(cmd)
	r, err := newSsh.SendCmd(cmd)
	if err != nil {
		zap.S().Errorln(err)
		return nil, err
	}
	zap.S().Info(r)
	return &responses.OK, nil
}
