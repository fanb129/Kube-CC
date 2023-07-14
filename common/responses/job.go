package responses

// AppJob job的信息,用于查看详细信息
type AppJob struct {
	Name        string `json:"name"`
	Namespace   string `json:"namespace"`
	Succeeded   int32  `json:"succeeded"`   // 完成数
	Completions int32  `json:"completions"` // 任务数
	Duration    string `json:"duration"`    // 用时
	Image       string `json:"image" form:"image" binding:"required"`
	CreatedAt   string `json:"created_at"`
	PodList     []JobPod
}

type AppJobList struct {
	Response
	Length     int      `json:"length"`
	AppJobList []AppJob `json:"app_job_list"`
}
