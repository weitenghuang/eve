package eve

type HealthService interface {
	GetHealth() *HealthInfo
}

type QuoinService interface {
	GetQuoin(name string) (*Quoin, error)
	GetQuoinArchive(id string) (*QuoinArchive, error)
	GetQuoinArchiveIds(quoinName string) ([]string, error)
	GetQuoinArchiveIdFromUri(archiveUri string) string
	CreateQuoin(quoin *Quoin) (*Quoin, error)
	CreateQuoinArchive(quoinArchive *QuoinArchive) error
	DeleteQuoin(name string) error
	DeleteQuoinArchive(id string) error
}

type InfrastructureService interface {
	QuoinService
	GetInfrastructure(name string) (*Infrastructure, error)
	GetInfrastructureState(name string) (map[string]interface{}, error)
	CreateInfrastructure(infra *Infrastructure) error
	DeleteInfrastructure(name string) error
	DeleteInfrastructureState(name string) error
	UpdateInfrastructureState(name string, state map[string]interface{}) error
	UpdateInfrastructureStatus(name string, status Status) error
	SubscribeAsyncProc(subject Subject, handler InfrastructureAsyncHandler) error
	PublishMessageToQueue(subject Subject, infra *Infrastructure) error
}

type InfrastructureAsyncHandler func(infra *Infrastructure)

type HealthInfo struct {
	Hostname string            `json:"hostname"`
	Errors   []Error           `json:"errors,omitempty"`
	Metadata map[string]string `json:"metadata"`
	Uptime   string            `json:"uptime"`
}

type Error struct {
	Description string            `json:"description"`
	Error       string            `json:"error"`
	Metadata    map[string]string `json:"metadata"`
	Type        string            `json:"type"`
}

type Quoin struct {
	Id            string        `json:"id,omitempty"`            // UUID for each entry. Generated by rethinkdb uuid() based on quoin.Name
	Name          string        `json:"name"`                    // quoin unique name as db index field
	ArchiveUri    string        `json:"archiveUri,omitempty"`    // quoin scheme://host:port/quoin/name/upload/
	Variables     []QuoinVar    `json:"variables,omitempty"`     // quoin module input variables
	Authorization Authorization `json:"authorization,omitempty"` // quoin authorization setting
}

// Quoin Archive content is a collection of terraform modules in tarball format
type QuoinArchive struct {
	Id            string
	QuoinName     string // Archive will be linked to specific quoin instance
	Modules       []byte
	Authorization Authorization
}

type QuoinVar struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Infrastructure struct {
	Id            string                 `json:"id,omitempty"`            // UUID for each entry. Generated by rethinkdb uuid() based on name
	Name          string                 `json:"name"`                    // infrastructure unique name as db index field
	Quoin         *Quoin                 `json:"quoin"`                   // infrastructure quoin type
	Variables     []QuoinVar             `json:"variables,omitempty"`     // infrastructure environment variables
	State         map[string]interface{} `json:"state,omitempty"`         // Terraform state output
	Status        Status                 `json:"status,omitempty"`        // infrastructure environment lifecycle status
	Authorization Authorization          `json:"authorization,omitempty"` // infrastructure authorization setting
	ProviderSlug  string                 `json:"providerSlug"`            // infrastructure provider in slug format <provider:schema-type> aws:account
}

// Team's permission on resource
type PolicyMode int

const (
	POLICY_NONE PolicyMode = iota
	POLICY_ALL  PolicyMode = 1 << iota
	POLICY_READ
	POLICY_WRITE
	POLICY_EXECUTE
	POLICY_READ_WRITE
	POLICY_READ_EXECUTE
	POLICY_WRITE_EXECUTE
)

type Authorization struct {
	// Owner could be a user, or an organization
	Owner UserId

	// Grant other team/organization the access to this resource
	GroupAccess map[Group]PolicyMode
}

type Group string

type Organization Group

type Team Group

type UserId string

type User struct {
	Id UserId

	Organization

	Teams []Team
}

type Status int

// Resource lifecycle status in iota int
const (
	DEFAULT   Status = iota
	VALIDATED Status = 1 << iota
	RUNNING
	DEPLOYED
	DESTROYED
	OBSOLETED
	FAILED
)

type Subject string

// NATS.io Message's "subject"
const (
	CREATE_INFRA Subject = "create-infra"
	DELETE_INFRA Subject = "delete-infra"
)

type ProviderService interface {
	GetProvider(name string) (*Provider, error)
}

type Provider struct {
	Id            string `json:id,omitempty`
	Name          string `json:name`
	Schema        Schema
	Authorization Authorization
}

type Schema struct {
	Type string      `json:type`
	Data interface{} `json:data`
}

const (
	AGENT_USER = "terraform"
)
