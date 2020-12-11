// Package api - ADMIN API 기능을 제공하는 패키지
package api

import (
	"net/http"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/errors"
)

// ===== [ Constants and Variables ] =====

const (
	// RemovedOperation - 설정 제거 작업
	RemovedOperation ConfigurationOperation = iota
	// UpdatedOperation - 설정 변경 작업
	UpdatedOperation
	// AddedOperation - 설정 등록 작업
	AddedOperation
	// RemovedGroupOperation - 소스 제거 작업
	RemovedGroupOperation
	// AddedGroupOperation - 소스 추가 작업
	AddedGroupOperation
	// ApplyGroupsOperation - 설정 변경사항 모두 저장 (File or ETCD, ...)
	ApplyGroupsOperation
)

const (
	// NONE - 로드한 후에 변경이 없는 상태
	NONE ConfigurationState = iota
	// ADDED - Group가 신규로 생성된 경우
	ADDED
	// REMOVED - Group가 삭제된 경우
	REMOVED
	// CHANGED - Group 내의 Definition이 추가/수정/삭제된 경우
	CHANGED
)

var (
	// ErrAPIDefinitionNotFound - 레파지토리에 API 정의가 존재하지 않는 경우 오류
	ErrAPIDefinitionNotFound = errors.NewWithCode(http.StatusNotFound, "api definition not found")
	// ErrAPIsNotChanged - 레파지토리에 저장할 API 정의 변경 사항이 존재하지 않는 경우 오류
	ErrAPIsNotChanged = errors.NewWithCode(http.StatusNotModified, "api definitions are not changed")
	// ErrAPINameExists - 레파지토리에 동일한 이름의 API 정의가 존재하는 경우 오류
	ErrAPINameExists = errors.NewWithCode(http.StatusConflict, "api name is already registered")
	// ErrAPIListenPathExists - 레파지토리에 동일한 수신 경로의 API 정의가 존재하는 경우 오류
	ErrAPIListenPathExists = errors.NewWithCode(http.StatusConflict, "api listen path is already registered")
	// ErrGroupExists - 리파지토리에 동일한 이름의 소스가 존재하는 경우 오류
	ErrGroupExists = errors.NewWithCode(http.StatusConflict, "api group is already registered")
	// ErrGroupNotExists - 리파지토리에 동일한 이름의 소스가 존재하지 않는 경우 오류
	ErrGroupNotExists = errors.NewWithCode(http.StatusNotFound, "api group not found")
	// ErrInvalidRequestData - 요청에서 데이터를 추출하지 못했을 경우 오류
	ErrInvalidRequestData = errors.NewWithCode(http.StatusBadRequest, "invalid requested data")
)

// ===== [ Types ] =====

type (
	// ConfigurationState - Configuration 변경 상태 형식
	ConfigurationState int
	// ConfigurationOperation - Configuration 변경에 연계되는 Operation 형식
	ConfigurationOperation int

	// ConfigModel - 클라이언트와 통신에 사용할 정보 구조
	ConfigModel struct {
		Name        string                   `json:"name"`
		Definitions []*config.EndpointConfig `json:"definitions"`
	}

	// Configuration - API Definitions 관리 구조
	Configuration struct {
		DefinitionMaps []*DefinitionMap
	}

	// RepoChangedMessage - Repository의 변경이 발생한 경우 전송되는 메시지 형식 (Repository to Server)
	RepoChangedMessage struct {
		Configurations *Configuration
	}

	// ConfigChangedMessage - Configuration의 변경이 발생한 경우 전송되는 메시지 형식 (Server to Repository)
	ConfigChangedMessage struct {
		Operation   ConfigurationOperation
		Name        string
		Definitions []*config.EndpointConfig
	}
)

// ===== [ Implementations ] =====

// // EqualsTo - 현재 동작 중인 설정과 지정된 설정이 동일한지 여부 검증
// func (c *Configuration) EqualsTo(tc *Configuration) bool {
// 	return reflect.DeepEqual(c, tc)
// }

// GetAllDefinitions - 관리하고 있는 API Definition들 반환
func (c *Configuration) GetAllDefinitions() []*config.EndpointConfig {
	defs := make([]*config.EndpointConfig, 0)
	for _, dm := range c.DefinitionMaps {
		if dm.State != REMOVED {
			for _, def := range dm.Definitions {
				defs = append(defs, def)
			}
		}
	}
	return defs
}

// GetDefinitionMaps - 관리하고 있는 API Definition Map들 반환
func (c *Configuration) GetDefinitionMaps() []*DefinitionMap {
	maps := make([]*DefinitionMap, 0)
	for _, dm := range c.DefinitionMaps {
		if dm.State != REMOVED {
			maps = append(maps, dm)
		}
	}

	return maps
}

// ExistGroup - 지정한 Group가 존재하는지 검증
func (c *Configuration) ExistGroup(name string) bool {
	for _, dm := range c.DefinitionMaps {
		if dm.Name == name {
			return true
		}
	}
	return false
}

// Exists - 지정한 Group내에 지정한 Definition이 존재하는지 검증
func (c *Configuration) Exists(name string, ec *config.EndpointConfig) (bool, error) {
	for _, dm := range c.DefinitionMaps {
		if dm.Name == name {
			for _, def := range dm.Definitions {
				if def.Name == ec.Name {
					return true, ErrAPINameExists
				}

				if def.Endpoint == ec.Endpoint {
					return true, ErrAPIListenPathExists
				}
			}
		}
	}

	return false, nil
}

// ExistsDefinition - 지정한 Definition이 존재하는지 검증 (Group에 상관없음)
func (c *Configuration) ExistsDefinition(ec *config.EndpointConfig) error {
	for _, dm := range c.DefinitionMaps {
		for _, def := range dm.Definitions {
			if def.Name == ec.Name {
				return ErrAPINameExists
			}

			if def.Endpoint == ec.Endpoint {
				return ErrAPIListenPathExists
			}
		}
	}

	return nil
}

// FindByName - 지정한 이름의 Endpoint Definition이 존재하는지 검증 (동일 소스 대상)
func (c *Configuration) FindByName(gname, dname string) *config.EndpointConfig {
	for _, dm := range c.DefinitionMaps {
		if dm.Name == gname {
			for _, def := range dm.Definitions {
				if def.Name == dname {
					return def
				}
			}
		}
	}

	return nil
}

// FindByListenPath - 지정한 Path의 Endpoint Definition이 존재하는 검증 (전체 대상)
func (c *Configuration) FindByListenPath(listenPath string) *config.EndpointConfig {
	for _, dm := range c.DefinitionMaps {
		for _, def := range dm.Definitions {
			if def.Endpoint == listenPath {
				return def
			}
		}
	}

	return nil
}

// AddDefinition - 지정한 정보를 기준으로 관리 중인 API Defintion 추가
func (c *Configuration) AddDefinition(name string, ec *config.EndpointConfig) error {
	for _, dm := range c.DefinitionMaps {
		if dm.Name == name {
			dm.Definitions = append(dm.Definitions, ec)
			if dm.State != ADDED {
				dm.State = CHANGED
			}
			return nil
		}
	}
	return errors.New("Specific group path not exist [" + name + "]")
}

// UpdateDefinition - 지정한 정보를 기준으로 관리 중인 API Definition 갱신
func (c *Configuration) UpdateDefinition(name string, ec *config.EndpointConfig) error {
	for _, dm := range c.DefinitionMaps {
		if dm.Name == name {
			for i, def := range dm.Definitions {
				if def.Name == ec.Name {
					dm.Definitions[i] = ec
					if dm.State != ADDED {
						dm.State = CHANGED
					}
					return nil
				}
			}
		}
	}

	return errors.New("No definition to update in group path [" + name + "]")
}

// RemoveDefinition - 지정한 정보를 기준으로 관리 중인 API Definition 삭제
func (c *Configuration) RemoveDefinition(name string, ec *config.EndpointConfig) error {
	for _, dm := range c.DefinitionMaps {
		// find group
		if dm.Name == name {
			// remove definitions
			for i, def := range dm.Definitions {
				if def.Name == ec.Name {
					dm.Definitions = append(dm.Definitions[:i], dm.Definitions[i+1:]...)
					if dm.State != ADDED {
						dm.State = CHANGED
					}
					return nil
				}
			}
		}
	}
	return errors.New("No defintion to remove in group path [" + name + "]")
}

// GetGroup - 지정한 Group 정보 반환
func (c *Configuration) GetGroup(name string) *DefinitionMap {
	for _, dm := range c.DefinitionMaps {
		if dm.Name == name {
			return dm
		}
	}
	return nil
}

// AddGroupAndDefinitions - 지정한 정보를 기준으로 API Group을 생성하고 Definition들 추가
func (c *Configuration) AddGroupAndDefinitions(name string, ecs []*config.EndpointConfig) error {
	c.DefinitionMaps = append(c.DefinitionMaps, &DefinitionMap{Name: name, State: ADDED, Definitions: ecs})
	return nil
}

// AddGroup - 지정한 정보를 기준으로 API Group 생성
func (c *Configuration) AddGroup(name string) error {
	c.DefinitionMaps = append(c.DefinitionMaps, &DefinitionMap{Name: name, State: ADDED, Definitions: make([]*config.EndpointConfig, 0)})
	return nil
}

// RemoveGroup - 지정한 정보를 기준으로 API Group 삭제
func (c *Configuration) RemoveGroup(name string) error {
	for i, dm := range c.DefinitionMaps {
		if dm.Name == name {
			if dm.State == ADDED {
				c.DefinitionMaps = append(c.DefinitionMaps[:i], c.DefinitionMaps[i+1:]...)
				return nil
			}

			dm.State = REMOVED
			return nil
		}
	}
	return errors.New("No group to remove in groups [" + name + "]")
}

// ClearRemoved - 현재 관리 중인 API Defintion Soruce들 중에서 삭제된 내용을 제거
func (c *Configuration) ClearRemoved() {
	for i, dm := range c.DefinitionMaps {
		if dm.State == REMOVED {
			c.DefinitionMaps = append(c.DefinitionMaps[:i], c.DefinitionMaps[i+1:]...)
		}
	}
}

// ===== [ Private Functions ] =====
// ===== [ Public Functions ] =====
