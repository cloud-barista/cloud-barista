package tumblebug

type VMStatus string

const (
	VM_USER_ACCOUNT = "cb-user"

	VMSTATUS_CREATING VMStatus = "Creating" // from launch to running
	VMSTATUS_RUNNING  VMStatus = "Running"
	VMSTATUS_FAILED   VMStatus = "Failed"
	//VMSTATUS_SUSPENDING  VMStatus = "Suspending" // from running to suspended
	//VMSTATUS_SUSPENDED   VMStatus = "Suspended"
	//VMSTATUS_RESUMING    VMStatus = "Resuming"    // from suspended to running
	//VMSTATUS_REBOOTING   VMStatus = "Rebooting"   // from running to running
	//VMSTATUS_TERMINATING VMStatus = "Terminating" // from running, suspended to terminated
	//VMSTATUS_TERMINATED  VMStatus = "Terminated"
	//VMSTATUS_NOTEXIST    VMStatus = "NotExist" // VM does not exist
)

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Model struct {
	Name      string `json:"name"`
	Namespace string
}

// Namespace
type NS struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Connection info.
type Connection struct {
	Model
	ConfigName     string `json:"ConfigName"`
	ProviderName   string `json:"ProviderName"`
	DriverName     string `json:"DriverName"`
	CredentialName string `json:"CredentialName"`
	RegionName     string `json:"RegionName"`
}

type Region struct {
	Model
	RegionName       string     `json:"RegionName"`
	ProviderName     string     `json:"ProviderName"`
	KeyValueInfoList []KeyValue `json:"KeyValueInfoList"`
}

// MCIR
type VPC struct {
	Model
	Config       string     `json:"connectionName"`
	CidrBlock    string     `json:"cidrBlock"`
	Subnets      []Subnet   `json:"subnetInfoList"`
	Description  string     `json:"description"`
	CspVNetId    string     `json:"cspVNetId"`    // output
	CspVNetName  string     `json:"cspVNetName"`  // output
	Status       string     `json:"status"`       // output
	KeyValueList []KeyValue `json:"keyValueList"` // output
}

type Subnet struct {
	Name      string `json:"Name"`
	CidrBlock string `json:"IPv4_CIDR"`
}

type Firewall struct {
	Model
	Config               string          `json:"connectionName"`
	VPCId                string          `json:"vNetId"`
	Description          string          `json:"description"`
	FirewallRules        []FirewallRules `json:"firewallRules"`
	CspSecurityGroupId   string          `json:"cspSecurityGroupId"`   // output
	CspSecurityGroupName string          `json:"cspSecurityGroupName"` // output
	KeyValueList         []KeyValue      `json:"keyValueList"`         // output
}
type FirewallRules struct {
	From      string `json:"fromPort"`
	To        string `json:"toPort"`
	Protocol  string `json:"ipProtocol"`
	Direction string `json:"direction"`
}

type Image struct {
	Model
	Config       string     `json:"connectionName"`
	CspImageId   string     `json:"cspImageId"`
	CspImageName string     `json:"cspImageName"` // output
	CreationDate string     `json:"creationDate"` // output
	Description  string     `json:"description"`  //
	GuestOS      string     `json:"guestOS"`      //
	Status       string     `json:"status"`       // output
	KeyValueList []KeyValue `json:"keyValueList"` // output
}

type Spec struct {
	Model
	Config      string `json:"connectionName"`
	CspSpecName string `json:"cspSpecName"`
}
type SSHKey struct {
	Model
	Config     string `json:"connectionName"`
	Username   string `json:"username"`
	PrivateKey string `json:"privateKey"` // output
}

type LookupSpec struct {
	Model
	Config string `json:"connectionName"`
	Spec   string `json:"cspSpecName"`
	Region string `json:"region"` // output
	Memory string `json:"mem"`    // output
	CPU    struct {
		Count string `json:"count"` // output
		Clock string `json:"clock"` // output - GHz
	} `json:"vcpu"`
}

type LookupSpecs struct {
	Model
	Config  string `json:"connectionName"`
	Vmspecs []struct {
		Name   string `json:"name"` // output
		Memory string `json:"mem"`  // output
		CPU    struct {
			Count string `json:"count"` // output
			Clock string `json:"clock"` // output - GHz
		} `json:"vcpu"`
	} `json:"vmspec"`
}

type LookupImages struct {
	Model
	ConnectionName string `json:"connectionName"`
	Images         []struct {
		Name string `json:"name"`
		IId  struct {
			NameId   string `json:"nameId"`   // output - NameID by user
			SystemId string `json:"systemId"` // output - SystemID by CloudOS
		} `json:"iid"`
		GuestOS      string     `json:"guestOS"`      // output - Windows7, Ubuntu etc.
		Status       string     `json:"status"`       // output - available, unavailable
		KeyValueList []KeyValue `json:"keyValueList"` //output -
	} `json:"image"`
}

// MCIS
type MCIS struct {
	Model
	Description     string `json:"description"`
	Label           string `json:"label"`
	SystemLabel     string `json:"systemLabel"`
	InstallMonAgent string `json:"installMonAgent"`
	Status          string `json:"status"`
	VMs             []VM   `json:"vm"` // output
}

type VM struct {
	Model
	mcisName      string   //private
	Config        string   `json:"connectionName"`
	VPC           string   `json:"vNetId"`
	Subnet        string   `json:"subnetId"`
	Firewalls     []string `json:"securityGroupIds"`
	SSHKey        string   `json:"sshKeyId"`
	Image         string   `json:"imageId"`
	Spec          string   `json:"specId"`
	UserAccount   string   `json:"vmUserAccount"`
	UserPassword  string   `json:"vmUserPassword"`
	Description   string   `json:"description"`
	PublicIP      string   `json:"publicIP"`      // output
	PrivateIP     string   `json:"privateIP"`     // output
	Status        VMStatus `json:"status"`        // output
	SystemMessage string   `json:"systemMessage"` // output
	Region        struct {
		Region string `json:"region"`
		Zone   string `json:"zone"`
	} `json:"region"` // output
	CspViewVmDetail struct {
		VMSpecName string `json:"vmspecName"`
	} `json:"cspViewVmDetail"` // output

}
