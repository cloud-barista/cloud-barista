// Package admin -
package admin

import (
	"errors"
	"net/http"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/admin/response"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/api"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/core"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/core/adapters/gin"
	"go.opencensus.io/trace"
)

// ===== [ Constants and Variables ] =====
// ===== [ Types ] =====

type (
	// APIHandler - Admin API 운영을 위한 REST Handler 정보 형식
	APIHandler struct {
		configurationChan chan<- api.ConfigChangedMessage

		Configs *api.Configuration
	}
)

// ===== [ Implementations ] =====

// GetDefinitions - 전체 Definition 정보 반환
func (ah *APIHandler) GetDefinitions() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		_, span := trace.StartSpan(req.Context(), "definitions.GetAll")
		defer span.End()

		if nil == ah.Configs.DefinitionMaps {
			// API Definition이 없을 경우는 빈 JSON Array 처리 (ID 기준)
			response.Write(rw, req, []int{})
			return
		}

		response.Write(rw, req, ah.Configs.GetDefinitionMaps())
	}
}

// UpdateDefinition - Request의 정보를 기준으로 Definition 정보 갱신
func (ah *APIHandler) UpdateDefinition() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		cm := &api.ConfigModel{}

		err := core.JSONDecode(req.Body, cm)
		if nil != err {
			response.Errorf(rw, req, -1, err)
			return
		}

		_, span := trace.StartSpan(req.Context(), "definition.FindByName")
		def := ah.Configs.FindByName(cm.Name, cm.Definitions[0].Name)
		span.End()

		if nil == def {
			response.Errorf(rw, req, -1, api.ErrAPIDefinitionNotFound)
			return
		}

		err = cm.Definitions[0].Validate()
		if nil != err {
			response.Errorf(rw, req, -1, err)
			return
		}

		// 동일한 경로가 다른 Definition 이름으로 등록되어있는 경우 검증 (전체 대상)
		_, span = trace.StartSpan(req.Context(), "repo.FindByListenPath")
		existingDef := ah.Configs.FindByListenPath(cm.Definitions[0].Endpoint)
		span.End()

		if nil != existingDef && existingDef.Name != cm.Definitions[0].Name {
			response.Errorf(rw, req, -1, api.ErrAPIListenPathExists)
			return
		}

		// Definition 갱신 채널 처리
		_, span = trace.StartSpan(req.Context(), "repo.Update")
		ah.configurationChan <- api.ConfigChangedMessage{
			Name:        cm.Name,
			Operation:   api.UpdatedOperation,
			Definitions: cm.Definitions,
		}
		span.End()

		response.Write(rw, req, nil)
	}
}

// AddDefinition - Request 정보를 기준으로 Definition 정보 추가
func (ah *APIHandler) AddDefinition() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		cm := &api.ConfigModel{}

		err := core.JSONDecode(req.Body, cm)
		if nil != err {
			response.Errorf(rw, req, -1, err)
			return
		}

		err = cm.Definitions[0].Validate()
		if nil != err {
			response.Errorf(rw, req, -1, err)
			return
		}

		// 기존 정보가 존재하는지 검증 (Name 및 Endpoint Path, ...)
		_, span := trace.StartSpan(req.Context(), "definition.Exists")
		exists, err := ah.Configs.Exists(cm.Name, cm.Definitions[0])
		span.End()

		if nil != err {
			response.Errorf(rw, req, -1, err)
			return
		}
		if exists {
			response.Errorf(rw, req, -1, api.ErrAPINameExists)
			return
		}

		// Definition 갱신 채널 처리
		_, span = trace.StartSpan(req.Context(), "repo.Add")
		ah.configurationChan <- api.ConfigChangedMessage{
			Name:        cm.Name,
			Operation:   api.AddedOperation,
			Definitions: cm.Definitions,
		}
		span.End()

		response.Write(rw, req, nil)
	}
}

// RemoveDefinition - Request 정보를 기준으로 Definition 정보 삭제
func (ah *APIHandler) RemoveDefinition() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		cm := &api.ConfigModel{}

		cm.Name = gin.URLParam(req, "gid")
		cm.Definitions = make([]*config.EndpointConfig, 0)

		def := config.NewDefinition()
		def.Name = gin.URLParam(req, "id")
		cm.Definitions = append(cm.Definitions, def)

		_, span := trace.StartSpan(req.Context(), "definition.Exists")
		exists, err := ah.Configs.Exists(cm.Name, cm.Definitions[0])
		span.End()

		if nil != err && err != api.ErrAPINameExists {
			response.Errorf(rw, req, -1, err)
			return
		}
		if !exists {
			response.Errorf(rw, req, -1, api.ErrAPIDefinitionNotFound)
			return
		}

		_, span = trace.StartSpan(req.Context(), "repo.Remove")
		defer span.End()

		// Definition 갱신 채널 처리
		ah.configurationChan <- api.ConfigChangedMessage{
			Name:        cm.Name,
			Operation:   api.RemovedOperation,
			Definitions: cm.Definitions,
		}

		response.Write(rw, req, nil)
	}
}

// GetGroup - 지정된 Group 정보 반환
func (ah *APIHandler) GetGroup() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		cm := &api.ConfigModel{}

		cm.Name = gin.URLParam(req, "gid")

		//cm.Name = core.GetURLVariable(req.URL.Path, GroupBasePath)

		if "" == cm.Name {
			response.Errorf(rw, req, -1, errors.New("cannot found group data from request"))
			return
		}

		_, span := trace.StartSpan(req.Context(), "repo.ExistsGroup")
		exists := ah.Configs.ExistGroup(cm.Name)
		span.End()

		if !exists {
			response.Errorf(rw, req, -1, api.ErrGroupNotExists)
			return
		}

		_, span = trace.StartSpan(req.Context(), "repo.GetGroup")
		defer span.End()

		if nil == ah.Configs.DefinitionMaps {
			// API Definition이 없을 경우는 빈 JSON Array 처리 (ID 기준)
			response.Write(rw, req, []int{})
			return
		}

		response.Write(rw, req, ah.Configs.GetGroup(cm.Name))
	}
}

// AddGroup - Request 정보를 기준으로 Definition을 관리하기 위한 신규 Group 생성
func (ah *APIHandler) AddGroup() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		cm := &api.ConfigModel{}

		err := core.JSONDecode(req.Body, cm)
		if nil != err {
			response.Errorf(rw, req, -1, err)
			return
		}

		_, span := trace.StartSpan(req.Context(), "repo.ExistsGroup")
		exists := ah.Configs.ExistGroup(cm.Name)
		span.End()

		if exists {
			response.Errorf(rw, req, -1, api.ErrGroupExists)
			return
		}

		if len(cm.Definitions) > 0 {
			_, span := trace.StartSpan(req.Context(), "repo.Exists")
			for _, ec := range cm.Definitions {
				err := ah.Configs.ExistsDefinition(ec)
				if nil != err {
					response.Errorf(rw, req, -1, err)
					span.End()
					return
				}
			}
			span.End()
		}

		// Group 등록
		_, span = trace.StartSpan(req.Context(), "repo.AddGroup")
		// Group 추가 채널 처리
		ah.configurationChan <- api.ConfigChangedMessage{
			Name:        cm.Name,
			Operation:   api.AddedGroupOperation,
			Definitions: cm.Definitions,
		}
		span.End()

		response.Write(rw, req, nil)
	}
}

// RemoveGroup - Request 정보를 기준으로 관리 중인 Group 삭제
func (ah *APIHandler) RemoveGroup() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		cm := &api.ConfigModel{}

		cm.Name = gin.URLParam(req, "gid")

		if "" == cm.Name {
			response.Errorf(rw, req, -1, errors.New("cannot found group data from request"))
			return
		}

		_, span := trace.StartSpan(req.Context(), "repo.Exists")
		exists := ah.Configs.ExistGroup(cm.Name)
		span.End()

		if !exists {
			response.Errorf(rw, req, -1, api.ErrGroupNotExists)
			return
		}

		// Group 삭제
		_, span = trace.StartSpan(req.Context(), "repo.RemoveGroup")
		// Group 삭제 채널 처리
		ah.configurationChan <- api.ConfigChangedMessage{
			Name:      cm.Name,
			Operation: api.RemovedGroupOperation,
		}
		span.End()

		response.Write(rw, req, nil)
	}
}

// ApplyGroups - 관리 중인 Group들의 변경사항 적용 (Persistence)
func (ah *APIHandler) ApplyGroups() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		_, span := trace.StartSpan(req.Context(), "repo.ApplyGroups")
		span.End()

		// 리파지토리에 전체 변경내역 저장
		ah.configurationChan <- api.ConfigChangedMessage{
			Operation: api.ApplyGroupsOperation,
		}

		response.Write(rw, req, nil)
	}
}

// ===== [ Private Functions ] =====
// ===== [ Public Functions ] =====

// NewAPIHandler - 지정한 Configuration 변경을 설정한 Admin API Handler 인스턴스 생성
func NewAPIHandler(configurationChan chan<- api.ConfigChangedMessage) *APIHandler {
	return &APIHandler{
		configurationChan: configurationChan,
	}
}
