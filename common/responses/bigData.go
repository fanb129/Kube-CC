package responses

import "Kube-CC/common/forms"

type Bigdata struct {
	Ns
	DeployList []AppDeploy `json:"deploy_list"`
}

type BigdataListResponse struct {
	Response
	Length      int       `json:"length"`
	BigdataList []Bigdata `json:"bigdata_list"`
}

type InfoHadoop struct {
	Response
	Form forms.HadoopUpdateForm
}

type InfoSpark struct {
	Response
	Form forms.SparkUpdateForm
}
