// Package api -
package api

import (
	"sync"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/logging"
)

// ===== [ Constants and Variables ] =====
// ===== [ Types ] =====
type (
	// InMemoryRepository - Memory 기반의 Repository 관리 형식
	InMemoryRepository struct {
		sync.RWMutex
		Groups []*DefinitionMap
	}
)

// ===== [ Implementations ] =====

// getGroup - 지정한 소스 경로에 맞는 GroupMap 반환
func (imr *InMemoryRepository) getGroup(path string) *DefinitionMap {
	for _, sm := range imr.Groups {
		if sm.Name == path {
			return sm
		}
	}

	return nil
}

// Add adds an api definition to the repository
func (imr *InMemoryRepository) add(group string, eConf *config.EndpointConfig) error {
	imr.Lock()
	defer imr.Unlock()

	log := logging.GetLogger()

	err := eConf.Validate()
	if nil != err {
		log.WithError(err).Error("[REPOSITORY] MEMORY > Validation errors")
		return err
	}

	sm := imr.getGroup(group)
	if nil != sm {
		sm.Definitions = append(sm.Definitions, eConf)
		log.Debug("[REPOSITORY] MEMORY > " + eConf.Name + " definition added to " + group + " group")
	} else {
		sm := &DefinitionMap{Name: group, State: NONE, Definitions: make([]*config.EndpointConfig, 0)}
		sm.Definitions = append(sm.Definitions, eConf)
		imr.Groups = append(imr.Groups, sm)
	}
	return nil
}

// Close - 사용 중인 Memory Repository 세션 종료
func (imr *InMemoryRepository) Close() error {
	return nil
}

// FindGroups - 리포지토리에서 관리하는 API Group 경로들을 반환
func (imr *InMemoryRepository) FindGroups() ([]string, error) {
	imr.RLock()
	defer imr.RUnlock()

	groups := make([]string, 0)
	for _, sm := range imr.Groups {
		groups = append(groups, sm.Name)
	}

	return groups, nil
}

// FindAllByGroup - 지정한 Group에서 API Routing 설정 정보 반환
func (imr *InMemoryRepository) FindAllByGroup(group string) ([]*config.EndpointConfig, error) {
	imr.RLock()
	defer imr.RUnlock()

	endpoints := make([]*config.EndpointConfig, 0)
	sm := imr.getGroup(group)
	if nil != sm {
		for _, v := range sm.Definitions {
			endpoints = append(endpoints, v)
		}
	}
	return endpoints, nil
}

// FindAll - 사용 가능한 모든 API Routing 설정 검증 및 반환
func (imr *InMemoryRepository) FindAll() ([]*DefinitionMap, error) {
	imr.RLock()
	defer imr.RUnlock()

	return imr.Groups, nil
}

// ===== [ Private Functions ] =====
// ===== [ Public Functions ] =====

// NewInMemoryRepository creates a in memory repository
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{Groups: make([]*DefinitionMap, 0)}
}
