package pod

import (
	"fmt"
	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"k8s_deploy_gin/common"
	"k8s_deploy_gin/conf"
	"k8s_deploy_gin/dao"
	"k8s_deploy_gin/service"
	"k8s_deploy_gin/service/ws"
	"net/http"
)

func Index(c *gin.Context) {
	ns := c.DefaultQuery("ns", "")
	u_id := c.DefaultQuery("u_id", "")
	selector := ""
	if u_id != "" {
		label := map[string]string{
			"u_id": u_id,
		}
		// 将map标签转换为string
		selector = labels.SelectorFromSet(label).String()
	}
	podListResponse, err := service.GetPod(ns, selector)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, podListResponse)
	}
}

func Delete(c *gin.Context) {
	ns := c.Query("ns")
	name := c.Query("name")
	response, err := service.DeletePod(name, ns)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{StatusCode: -1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, response)
	}
}

func Ssh(c *gin.Context) {
	fmt.Println("pod ssh")
	var (
		wsConn        *ws.WsConnection
		restConf      *rest.Config
		sshReq        *rest.Request
		podName       string
		podNs         string
		containerName string
		executor      remotecommand.Executor
		handler       *service.StreamHandler
		err           error
	)

	podNs = c.Query("podNs")
	podName = c.Query("podName")
	containerName = c.Query("containerName")

	// 得到websocket长连接
	if wsConn, err = ws.InitWebsocket(c.Writer, c.Request); err != nil {
		fmt.Println(err)
		return
	}

	// 获取k8s rest client配置
	if restConf, err = clientcmd.BuildConfigFromFlags("", conf.KubeConfig); err != nil {
		fmt.Println(err)
		return
	}

	// URL长相:
	// https://172.18.11.25:6443/api/v1/namespaces/default/pods/
	//nginx-deployment-5cbd8757f-d5qvx/exec?
	//command=sh&container=nginx&stderr=true&stdin=true&stdout=true&tty=true

	fmt.Println("pod")
	sshReq = dao.ClientSet.CoreV1().RESTClient().Post().
		Resource("pods").Name(podName).Namespace(podNs).SubResource("exec").Param("container", containerName).
		VersionedParams(&corev1.PodExecOptions{
			Container: containerName,
			Command:   []string{"bash"},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
		}, scheme.ParameterCodec)

	// 创建到容器的连接
	if executor, err = remotecommand.NewSPDYExecutor(restConf, "POST", sshReq.URL()); err != nil {
		fmt.Println(1, err)
		goto END
	}

	// 配置与容器之间的数据流处理回调

	handler = &service.StreamHandler{WsConn: wsConn, ResizeEvent: make(chan remotecommand.TerminalSize)}
	if err = executor.Stream(remotecommand.StreamOptions{
		Stdin:             handler,
		Stdout:            handler,
		Stderr:            handler,
		TerminalSizeQueue: handler,
		Tty:               true,
	}); err != nil {
		fmt.Println(2, err)
		goto END
	}
	return

END:
	fmt.Println(err)
	wsConn.WsClose()
}
