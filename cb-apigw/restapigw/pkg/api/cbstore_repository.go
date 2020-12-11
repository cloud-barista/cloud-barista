// Package api - CB-Store 기반 Repository
package api

import (
	"context"
	"time"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/core"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/logging"
	cbstore "github.com/cloud-barista/cb-store"
	icbs "github.com/cloud-barista/cb-store/interfaces"
)

// ===== [ Constants and Variables ] =====
// ===== [ Types ] =====

type (
	// CbStoreRepository - CB-Store 기반 Repository 관리 정보 형식
	CbStoreRepository struct {
		sConf *config.ServiceConfig
		*InMemoryRepository
		store       icbs.Store
		storeKey    string
		refreshTime time.Duration
	}
)

// ===== [ Implementations ] =====

// getStorePath - API Definition 저장을 위한 Store key path 반환
func (csr *CbStoreRepository) getStorePath(path string) string {
	return csr.storeKey + "/" + path
}

// Write - 변경된 리파지토리 내용을 대상 파일로 출력
func (csr *CbStoreRepository) Write(definitionMaps []*DefinitionMap) error {
	csr.Groups = definitionMaps

	for _, dm := range csr.Groups {
		if dm.State == REMOVED {
			err := csr.store.Delete(csr.getStorePath(dm.Name))
			if nil != err {
				return err
			}
		} else if dm.State != NONE {
			data, err := groupDefinitions(dm)
			if nil != err {
				return err
			}

			err = csr.store.Put(csr.getStorePath(dm.Name), string(data))
			if nil != err {
				return err
			}
		}
		dm.State = NONE
	}

	return nil
}

// Close - 사용 중인 Repository 세션 종료
func (csr *CbStoreRepository) Close() error {
	csr.store.Close()
	logging.GetLogger().Debug("[REPOSITORY] CB-STORE > Repository closed")
	return nil
}

// Watch - CB-STORE 리파지토리의 변경 감시 및 처리 (Timer Reading)
func (csr *CbStoreRepository) Watch(ctx context.Context, repoChan chan<- RepoChangedMessage) {
	log := logging.GetLogger()
	ticker := time.NewTicker(csr.refreshTime)

	go func(refreshTicker *time.Ticker) {
		defer refreshTicker.Stop()
		log.Debug("[REPOSITORY] CB-STORE > Watching CB-Store repository...")

		for {
			select {
			case <-refreshTicker.C:
				definitionMaps, err := csr.FindAll()
				if nil != err {
					log.WithError(err).Error("[REPOSITORY] CB-STORE > Failed to get API definitions on watch")
					continue
				}

				repoChan <- RepoChangedMessage{Configurations: &Configuration{DefinitionMaps: definitionMaps}}
			case <-ctx.Done():
				return
			}
		}
	}(ticker)
}

// ===== [ Private Functions ] =====
// ===== [ Public Functions ] =====

// NewCbStoreRepository - CB-Store 기반의 Repository 인스턴스 생성
func NewCbStoreRepository(sConf *config.ServiceConfig, key string, refreshTime time.Duration) (*CbStoreRepository, error) {
	log := logging.GetLogger()
	repo := CbStoreRepository{sConf: sConf, InMemoryRepository: NewInMemoryRepository(), store: cbstore.GetStore(), storeKey: key, refreshTime: refreshTime}

	// Grab configuration from CB-STORE
	keyValues, err := repo.store.GetList(key, true)
	if nil != err {
		return nil, err
	}

	for _, kv := range keyValues {
		// Skip Root
		if kv.Key == key {
			continue
		}

		apiDef, err := parseEndpoint(sConf, []byte(kv.Value))
		if nil != err {
			log.WithError(err).Error("[REPOSITORY] CB-STORE > Failed during parsing definitions")
			return nil, err
		}

		for _, def := range apiDef.Definitions {
			if err := repo.add(core.GetLastPart(kv.Key, "/"), def); nil != err {
				log.WithField("endpoint", def.Endpoint).WithError(err).Error("[REPOSITORY] CB-STORE > Failed during add endpoint to the repository")
				return nil, err
			}
		}
	}

	return &repo, nil
}
