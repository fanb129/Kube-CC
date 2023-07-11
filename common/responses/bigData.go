package responses

type Bigdata struct {
	Ns
	DeployList []AppDeploy `json:"deploy_list"`
}

type BigdataListResponse struct {
	Response
	Length      int       `json:"length"`
	BigdataList []Bigdata `json:"bigdata_list"`
}
