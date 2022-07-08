package common

type TbGeoLocation struct {
	BriefAddr     string `json:"briefAddr"`
	CloudType     string `json:"cloudType"`
	Latitude      string `json:"latitude"`
	Longitude     string `json:"longitude"`
	NativeRegion  string `json:"nativeRegion"`
}