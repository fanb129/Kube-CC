package responses

type Sc struct {
	Name                 string `json:"name"`
	CreatedAt            string `json:"created_at"`
	ReclaimPolicy        string `json:"reclaim_policy"`
	Provisioner          string `json:"provisioner"`
	VolumeBindingMode    string `json:"volume_binding_mode"`
	AllowVolumeExpansion bool   `json:"allow_volume_expansion"`
}

type ScListResponse struct {
	Response
	Length  int  `json:"length"`
	PvcList []Sc `json:"sc_list"`
}
