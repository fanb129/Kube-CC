package forms

type NodeAddForm struct {
	Nodes []struct {
		Host string `json:"host" form:"host" binding:"required"`
	} `json:"nodes" form:"nodes" binding:"required"`

	//Port     int    `json:"port" form:"port" binding:"required"`
	//User     string `json:"user" form:"user" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}
