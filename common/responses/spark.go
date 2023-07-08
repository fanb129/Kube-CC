package responses

type Spark struct {
	Ns
	DeployList []AppDeploy `json:"deploy_list"`
}

type SparkListResponse struct {
	Response
	Length    int     `json:"length"`
	SparkList []Spark `json:"spark_list"`
}
