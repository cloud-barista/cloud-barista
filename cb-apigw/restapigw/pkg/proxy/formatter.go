// Package proxy - Backend 결과에 대한 Mapping, Whitelist, Blacklist 등의 처리를 수행하는 Formatter 패키지
package proxy

import (
	"strings"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/core"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/core/flatmap/tree"
)

// ===== [ Constants and Variables ] =====

const (
	flatmapFilter = "flatmap_filter"
)

// ===== [ Types ] =====

// flatmapFormatter - Flatmap을 사용해서 Response를 Format 처리하는 구조 정의
type flatmapFormatter struct {
	Target string
	Prefix string
	Ops    []flatmapOp
}

// flatmapOp - Flatmap 운영을 위한 설정 구조 정의
type flatmapOp struct {
	Type string     `yaml:"type"`
	Args [][]string `yarml:"args"`
}

// propertyFilter - Reponse Filtering에 사용할 함수 정의
type propertyFilter func(*Response)

// entityFormatter - PropertyFilter를 이용해서 Response를 Format처리하기 위한 설정 구조 정의
type entityFormatter struct {
	Target         string
	Prefix         string
	PropertyFilter propertyFilter
	Mapping        map[string]string
}

// EntityFormatter - Response를 Format 처리하는 Entity 기반 인터페이스 정의
type EntityFormatter interface {
	Format(Response) Response
}

// ===== [ Implementations ] =====

// Format - PropertyFilter를 활용하는 EntityFormatter 구현
func (ef entityFormatter) Format(entity Response) Response {
	// Target 처리
	if ef.Target != "" {
		extractTarget(ef.Target, &entity)
	}
	if 0 < len(entity.Data) {
		ef.PropertyFilter(&entity)
		// Mapping 처리
		for formerKey, newKey := range ef.Mapping {
			if v, ok := entity.Data[formerKey]; ok {
				// Collection Wrapping에 대한 정보 재 처리
				if _, ok := entity.Data[core.CollectionTag]; ok {
					if _, ok := entity.Data[core.WrappingTag]; ok {
						entity.Data[core.WrappingTag] = newKey
					}
				}
				entity.Data[newKey] = v
				delete(entity.Data, formerKey)
			}
		}
	}
	if "" != ef.Prefix {
		entity.Data = map[string]interface{}{ef.Prefix: entity.Data}
	}
	return entity
}

// Format - Flatmap을 활용하는 FlatmapFormatter 구현
func (ff flatmapFormatter) Format(entity Response) Response {
	// Target 처리
	if "" != ff.Target {
		extractTarget(ff.Target, &entity)
	}

	// Flatmap 처리
	ff.processOps(&entity)

	if "" != ff.Prefix {
		entity.Data = map[string]interface{}{ff.Prefix: entity.Data}
	}
	return entity
}

// processOps - Flatmap 설정에 대한 처리
func (ff flatmapFormatter) processOps(entity *Response) {
	flatten, err := tree.New(entity.Data)
	if nil != err {
		return
	}
	for _, op := range ff.Ops {
		switch op.Type {
		case "move":
			// move - like whitelist and mapping
			flatten.Move(op.Args[0], op.Args[1])
		case "del":
			// del - like blacklist
			for _, val := range op.Args {
				flatten.Del(val)
			}

		default:
		}
	}

	// 처리된 데이터 설정
	entity.Data, _ = flatten.Get([]string{}).(map[string]interface{})
}

// ===== [ Private Functions ] =====

// whitelistPrune - 지정한 Whitelist 맵과 Dictionary 맵을 통해서 데이터 추출 (비 대상은 모두 제거)
func whitelistPrune(wlDict map[string]interface{}, inDict map[string]interface{}) bool {
	canDelete := true
	var deleteSibling bool
	for k, v := range inDict {
		deleteSibling = true
		if subWl, ok := wlDict[k]; ok {
			if subWlDict, okk := subWl.(map[string]interface{}); okk {
				if subInDict, isDict := v.(map[string]interface{}); isDict && !whitelistPrune(subWlDict, subInDict) {
					deleteSibling = false
				}
			} else {
				// whitelist leaf, maintain this branch
				deleteSibling = false
			}
		}
		if deleteSibling {
			delete(inDict, k)
		} else {
			canDelete = false
		}
	}
	return canDelete
}

// buildDictPath - 지정한 맵과 필드들의 정보를 이용해서 필드명 기준의 맵 생성
func buildDictPath(accumulator map[string]interface{}, fields []string) map[string]interface{} {
	ok := true
	var c map[string]interface{}
	var fIdx int
	fEnd := len(fields)
	p := accumulator
	for fIdx = 0; fIdx < fEnd; fIdx++ {
		if c, ok = p[fields[fIdx]].(map[string]interface{}); !ok {
			break
		}
		p = c
	}
	for ; fIdx < fEnd; fIdx++ {
		c = make(map[string]interface{})
		p[fields[fIdx]] = c
		p = c
	}
	return p
}

// newWhitelistFilter - 지정한 Whitelist를 Response에서 추출하기 위한 Filter 생성
func newWhitelistFilter(whitelist []string) propertyFilter {
	wlDict := make(map[string]interface{})
	for _, k := range whitelist {
		wlFields := strings.Split(k, ".")
		d := buildDictPath(wlDict, wlFields[:len(wlFields)-1])
		d[wlFields[len(wlFields)-1]] = true
	}

	return func(entity *Response) {
		if whitelistPrune(wlDict, entity.Data) {
			for k := range entity.Data {
				delete(entity.Data, k)
			}
		}
	}
}

// newBlacklistFilter - 지정한 Blacklist 를 Response에서 제거하기 위한 Filter 생성
func newBlacklistFilter(blacklist []string) propertyFilter {
	bl := make(map[string][]string, len(blacklist))
	for _, key := range blacklist {
		keys := strings.Split(key, ".")
		if 1 < len(keys) {
			if sub, ok := bl[keys[0]]; ok {
				bl[keys[0]] = append(sub, keys[1])
			} else {
				bl[keys[0]] = []string{keys[1]}
			}
		} else {
			bl[keys[0]] = []string{}
		}
	}

	return func(entity *Response) {
		for k, sub := range bl {
			if 0 == len(sub) {
				delete(entity.Data, k)
			} else {
				if tmp := blacklistFilterSub(entity.Data[k], sub); 0 < len(tmp) {
					entity.Data[k] = tmp
				}
			}
		}
	}
}

// blacklistFilterSub - 지정한 Value Map에서 지정한 Blacklist 지정 데이터를 제거
func blacklistFilterSub(vMap interface{}, blacklist []string) map[string]interface{} {
	tmp, ok := vMap.(map[string]interface{})
	if !ok {
		return map[string]interface{}{}
	}
	for _, key := range blacklist {
		delete(tmp, key)
	}
	return tmp
}

// newFlatmapFormatter - 지정된 BackendConfig 기준으로 Flatmap을 활용하는 Formatter 생성
func newFlatmapFormatter(bConf *config.BackendConfig) EntityFormatter {
	if v, ok := bConf.Middleware[MWNamespace]; ok {
		if e, ok := v.(map[string]interface{}); ok {
			if vs, ok := e[flatmapFilter].([]interface{}); ok {
				if 0 == len(vs) {
					return nil
				}
				ops := []flatmapOp{}
				for _, v := range vs {
					m, ok := v.(map[interface{}]interface{})
					if !ok {
						continue
					}
					op := flatmapOp{}
					if t, ok := m["type"].(string); ok {
						op.Type = t
					} else {
						continue
					}
					if args, ok := m["args"].([]interface{}); ok {
						op.Args = make([][]string, len(args))
						for k, arg := range args {
							if t, ok := arg.(string); ok {
								op.Args[k] = strings.Split(t, ".")
							}
						}
					}
					ops = append(ops, op)
				}
				if 0 == len(ops) {
					return nil
				}

				return &flatmapFormatter{
					Target: bConf.Target,
					Prefix: bConf.Group,
					Ops:    ops,
				}
			}
		}
	}
	return nil
}

// extractTarget - 지정한 Response에 대해 지정한 Target이 존재하는지를 검증하고 반환 (단, Map 형식이어야 하며, 만일 없거나, 변환 불가이면 빈 데이터로 처리)
func extractTarget(target string, entity *Response) {
	for _, part := range strings.Split(target, ".") {
		if tmp, ok := entity.Data[part]; ok {

			entity.Data, ok = tmp.(map[string]interface{})
			if !ok {
				entity.Data = map[string]interface{}{}
				return
			}
		} else {
			entity.Data = map[string]interface{}{}
			return
		}
	}
}

// ===== [ Public Functions ] =====

// NewEntityFormatter - 지정된 Backend 설정을 기준으로 Response 처리에 사용할 EntityFormatter 생성
func NewEntityFormatter(bConf *config.BackendConfig) EntityFormatter {
	if ff := newFlatmapFormatter(bConf); nil != ff {
		return ff
	}

	var pf propertyFilter
	if 0 < len(bConf.Whitelist) {
		// Response를 대상으로 whitelist 필터링
		pf = newWhitelistFilter(bConf.Whitelist)
	} else {
		// Response를 대상으로 blacklist 필터링
		pf = newBlacklistFilter(bConf.Blacklist)
	}
	sanitizedMappings := make(map[string]string, len(bConf.Mapping))
	for i, m := range bConf.Mapping {
		v := strings.Split(m, ".")
		sanitizedMappings[i] = v[0]
	}
	return entityFormatter{
		Target:         bConf.Target,
		Prefix:         bConf.Group,
		PropertyFilter: pf,
		Mapping:        sanitizedMappings,
	}
}
