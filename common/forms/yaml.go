package forms

type YamlApplyForm struct {
	Kind string      `json:"kind" form:"kind"`
	Name string      `json:"name" form:"name"`
	Ns   string      `json:"ns" form:"ns"`
	Yaml interface{} `json:"yaml" form:"yaml" binding:"required"`
}

type YamlCreateForm struct {
	Kind string      `json:"kind" form:"kind"`
	Ns   string      `json:"ns" form:"ns"`
	Yaml interface{} `json:"yaml" form:"yaml" binding:"required"`
}
