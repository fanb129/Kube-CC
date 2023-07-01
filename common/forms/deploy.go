package forms

type DeployAddForm struct {
	Name      string `json:"name" form:"name"`
	Namespace string `json:"namespace" form:"namespace"`
	//Uid       uint   `json:"u_id" form:"u_id"`
}
