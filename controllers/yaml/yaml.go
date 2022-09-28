package yaml

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s_deploy_gin/common"
	"k8s_deploy_gin/service/yamlApply"
	"net/http"
)

func Apply(c *gin.Context) {
	//body := c.Request.Body
	//bytes, _ := ioutil.ReadAll(body)
	//fmt.Println((bytes))
	form := common.YamlForm{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusOK, common.ValidatorResponse(err))
		return
	}
	yaml := form.Yaml.(map[string]interface{})
	fmt.Printf("%v", yaml)
	// 转为json
	jsonYaml, err := json.Marshal(yaml)
	if err != nil {
		goto END
	}
	switch yaml["kind"] {
	case "Namespace", "namespace":
		// json转为struct
		ns := corev1.Namespace{}
		if err = json.Unmarshal(jsonYaml, &ns); err != nil {
			goto END
		}
		if _, err = yamlApply.NamespaceApply(&ns); err != nil {
			goto END
		}
		goto SUCCESS
	case "Deployment", "deployment", "Deploy", "deploy":
		deploy := appsv1.Deployment{}
		if err = json.Unmarshal(jsonYaml, &deploy); err != nil {
			goto END
		}
		if _, err = yamlApply.DeployApply(&deploy); err != nil {
			goto END
		}
		goto SUCCESS
	case "Service", "service":
		svc := corev1.Service{}
		if err = json.Unmarshal(jsonYaml, &svc); err != nil {
			goto END
		}
		if _, err = yamlApply.ServiceApply(&svc); err != nil {
			goto END
		}
		goto SUCCESS
	case "Pod", "pod":
		pod := corev1.Pod{}
		if err = json.Unmarshal(jsonYaml, &pod); err != nil {
			goto END
		}
		if _, err = yamlApply.PodApply(&pod); err != nil {
			goto END
		}
		goto SUCCESS
	default:
		c.JSON(http.StatusOK, common.Response{
			StatusCode: -1,
			StatusMsg:  "类型不匹配",
		})
		return
	}
SUCCESS:
	c.JSON(http.StatusOK, common.OK)
	return
END:
	c.JSON(http.StatusOK, common.Response{StatusCode: -1, StatusMsg: err.Error()})
	return
}

func Create(c *gin.Context) {
	form := common.YamlForm{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusOK, common.ValidatorResponse(err))
		return
	}
	yaml := form.Yaml.(map[string]interface{})
	//fmt.Printf("%v", yaml)
	// 转为json
	jsonYaml, err := json.Marshal(yaml)
	if err != nil {
		goto END
	}
	switch yaml["kind"] {
	case "Namespace", "namespace":
		// json转为struct
		ns := corev1.Namespace{}
		if err = json.Unmarshal(jsonYaml, &ns); err != nil {
			goto END
		}
		if _, err = yamlApply.NamespaceCreate(&ns); err != nil {
			goto END
		}
		goto SUCCESS
	case "Deployment", "deployment":
		deploy := appsv1.Deployment{}
		if err = json.Unmarshal(jsonYaml, &deploy); err != nil {
			goto END
		}
		if _, err = yamlApply.DeployCreate(&deploy); err != nil {
			goto END
		}
		goto SUCCESS
	case "Service", "service":
		svc := corev1.Service{}
		if err = json.Unmarshal(jsonYaml, &svc); err != nil {
			goto END
		}
		if _, err = yamlApply.ServiceCreate(&svc); err != nil {
			goto END
		}
		goto SUCCESS
	case "Pod", "pod":
		pod := corev1.Pod{}
		if err = json.Unmarshal(jsonYaml, &pod); err != nil {
			goto END
		}
		if _, err = yamlApply.PodCreate(&pod); err != nil {
			goto END
		}
		goto SUCCESS
	default:
		c.JSON(http.StatusOK, common.Response{
			StatusCode: -1,
			StatusMsg:  "类型不匹配",
		})
		return
	}

SUCCESS:
	c.JSON(http.StatusOK, common.OK)
	return
END:
	c.JSON(http.StatusOK, common.Response{StatusCode: -1, StatusMsg: err.Error()})
	return
}

func Delete(c *gin.Context) {
	c.JSON(http.StatusOK, common.Response{1, "delete"})
}
