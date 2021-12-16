package common

type TbRestGetAllConfigResponse struct {
	Name       string       `json:"name"`
	ConfigInfo TbConfigInfo `json:"configInfo"`
}
