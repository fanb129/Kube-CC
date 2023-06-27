package forms

type YamlApplyForm struct {
	Kind string      `json:"kind" forms:"kind"`
	Name string      `json:"name" forms:"name"`
	Ns   string      `json:"ns" forms:"ns"`
	Yaml interface{} `json:"yaml" forms:"yaml" binding:"required"`
}

type YamlCreateForm struct {
	Kind string      `json:"kind" forms:"kind"`
	Ns   string      `json:"ns" forms:"ns"`
	Yaml interface{} `json:"yaml" forms:"yaml" binding:"required"`
}
