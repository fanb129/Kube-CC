package responses

type Hadoop struct {
	Ns
	DeployList []AppDeploy `json:"deploy_list"`
}

type HadoopListResponse struct {
	Response
	Length     int      `json:"length"`
	HadoopList []Hadoop `json:"hadoop_list"`
}
