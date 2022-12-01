package yaml

import (
	"Kube-CC/common"
	"Kube-CC/service/yamlApply"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"net/http"
)

func Apply(c *gin.Context) {
	//body := c.Request.Body
	//bytes, _ := ioutil.ReadAll(body)
	//fmt.Println((bytes))
	form := common.YamlApplyForm{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusOK, common.ValidatorResponse(err))
		return
	}
	yaml := form.Yaml.(map[string]interface{})
	kind := form.Kind
	name := form.Name
	ns := form.Ns
	metadata := yaml["metadata"].(map[string]interface{})
	var err error
	// 转为json
	jsonYaml, err := json.Marshal(yaml)
	if err != nil {
		goto END
	}
	if kind != "" && kind != yaml["kind"].(string) {
		err = errors.New("请勿修改yaml中kind")
		goto END
	}
	if name != "" && kind != metadata["name"].(string) {
		err = errors.New("请勿修改yaml中name")
		goto END
	}
	if ns != "" && kind != metadata["namespace"].(string) {
		err = errors.New("请勿修改yaml中namespace")
		goto END
	}
	switch yaml["kind"].(string) {
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
	form := common.YamlCreateForm{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusOK, common.ValidatorResponse(err))
		return
	}
	yaml := form.Yaml.(map[string]interface{})
	kind := form.Kind
	ns := form.Ns
	yamlKind := yaml["kind"].(string)
	//fmt.Printf("%v", yaml)
	// 转为json
	jsonYaml, err := json.Marshal(yaml)
	if err != nil {
		goto END
	}
	// 针对指定类型资源
	if kind != "" && ns == "" {
		err = errors.New("请选择namespace,并填入yaml")
		goto END
	}

	switch kind {
	case "Deploy":
		if yamlKind != "Deployment" && yamlKind != "deployment" {
			err = errors.New("yaml中kind必须为\"Deployment\"或\"Deploy\"")
			goto END
		}
	case "Service":
		if yamlKind != "Service" && yamlKind != "service" {
			err = errors.New("yaml中kind必须为\"Service\"或\"service\"")
			goto END
		}
	case "Pod":
		if yamlKind != "Pod" && yamlKind != "pod" {
			err = errors.New("yaml中kind必须为\"Pod\"或\"pod\"")
			goto END
		}
	}

	if ns != "" {
		metadata := yaml["metadata"].(map[string]interface{})
		if ns != metadata["namespace"].(string) {
			err = errors.New("yaml中namespace必须为所选" + ns)
			goto END
		}
	}
	// 所有类型资源
	switch yamlKind {
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
	case "Deployment", "deployment", "Deploy", "deploy":
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
