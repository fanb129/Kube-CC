package forms

type NodeAddForm struct {
	Nodes []struct {
		Host string `json:"host" forms:"host" binding:"required"`
	} `json:"nodes" forms:"nodes" binding:"required"`

	Port     int    `json:"port" forms:"port" binding:"required"`
	User     string `json:"user" forms:"user" binding:"required"`
	Password string `json:"password" forms:"password" binding:"required"`
}
