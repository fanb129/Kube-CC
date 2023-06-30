package forms

type DeployAddForm struct {
	Name      string `json:"name" forms:"name"`
	Namespace string `json:"namespace" forms:"namespace"`
	//Uid       uint   `json:"u_id" forms:"u_id"`
}
