// Package proxy - Backend의 결과들을 Merge 처리하는 Merging 패키지
package proxy

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/errors"
)

// ===== [ Constants and Variables ] =====

const (
	defaultCombinerName = "default"
	sequentialKey       = "sequential"
)

var (
	errNullResult     = errors.New("invalid response")
	responseCombiners = initResponseCombiners()
	reMergeKey        = regexp.MustCompile(`\{\{\.Resp(\d+)_([\d\w-_\.]+)\}\}`)
)

// ===== [ Types ] =====

// incrementalMergeAccumulator - 점진적인 Merging 처리를 위한 데이터 구조
type incrementalMergeAccumulator struct {
	pending  int
	data     *Response
	combiner ResponseCombiner
	errs     []error
}

// mergeError - Merging 과정에서 발생하는 오류들 관리 구조
type mergeError struct {
	errs []error
}

// ResponseCombiner - 여러 Response의 데이터를 Merging 처리해서 하나의 Response 데이터로 구성하는 함수 정의
type ResponseCombiner func(int, []*Response) *Response

// ===== [ Implementations ] =====

// Merge - 지정한 Response에 대한 점진적인 Merging 처리
func (ima *incrementalMergeAccumulator) Merge(res *Response, err error) {
	ima.pending--
	if nil != err {
		ima.errs = append(ima.errs, err)
		if nil != ima.data {
			ima.data.IsComplete = false
		}
		return
	}
	if nil == res {
		ima.errs = append(ima.errs, errNullResult)
		return
	}
	if nil == ima.data {
		ima.data = res
		return
	}
	ima.data = ima.combiner(2, []*Response{ima.data, res})
}

// Result - 처리된 Merging 결과 반환
func (ima *incrementalMergeAccumulator) Result() (*Response, error) {
	if nil == ima.data {
		return &Response{Data: make(map[string]interface{}, 0), IsComplete: false}, newMergeError(ima.errs)
	}

	if 0 != ima.pending || 0 != len(ima.errs) {
		ima.data.IsComplete = false
	}
	return ima.data, newMergeError(ima.errs)
}

// Error - Merging 작업 중에 발생한 오류 메시지 반환
func (me mergeError) Error() string {
	msg := make([]string, len(me.errs))
	for i, err := range me.errs {
		msg[i] = err.Error()
	}
	return strings.Join(msg, "\n")
}

// ===== [ Private Functions ] =====

// newMergeError - Merging 처리 중에 발생한 오류들을 하나의 오류로 반환
func newMergeError(errs []error) error {
	if 0 == len(errs) {
		return nil
	}
	return mergeError{errs}
}

// requestPart - 지정한 요청을 호출하고 오류나 null 반환에 대해 cancel 처리
func requestPart(ctx context.Context, next Proxy, req *Request, out chan<- *Response, failed chan<- error) {
	localCtx, cancel := context.WithCancel(ctx)

	res, err := next(localCtx, req)
	if nil != err {
		failed <- err
		cancel()
		return
	}
	if nil == res {
		failed <- errNullResult
		cancel()
		return
	}
	select {
	case out <- res:
	case <-ctx.Done():
		failed <- ctx.Err()
	}
	cancel()
}

// newIncrementalMergeAccumultor - 지정한 Backend count 와 ResponseCombiner를 설정한 점진적인 Merge 처리기 생성
func newIncrementalMergeAccumultor(backendCount int, rc ResponseCombiner) *incrementalMergeAccumulator {
	return &incrementalMergeAccumulator{
		pending:  backendCount,
		combiner: rc,
		errs:     []error{},
	}
}

// combineData - 지정한 Backend count와 Response들을 기준으로 Merging 처리된 Response 반환
func combineData(backendCount int, reses []*Response) *Response {
	isComplete := len(reses) == backendCount
	var mergedResponse *Response
	for _, res := range reses {
		if nil == res || nil == res.Data {
			isComplete = false
			continue
		}

		isComplete = isComplete && res.IsComplete
		if nil == mergedResponse {
			mergedResponse = res
			continue
		}

		for k, v := range res.Data {
			mergedResponse.Data[k] = v
		}
	}

	if nil == mergedResponse {
		// do not allow nil data to response
		return &Response{Data: make(map[string]interface{}, 0), IsComplete: isComplete}
	}
	mergedResponse.IsComplete = isComplete
	return mergedResponse
}

// initResponseCombiners - Response 데이터를 Merging 하는 ResponseCombiner 초기화
func initResponseCombiners() *combinerRegister {
	return newCombinerRegister(map[string]ResponseCombiner{defaultCombinerName: combineData}, combineData)
}

// getResponseCombiner - 기본적으로 사용되는 ResponseCombiner 반환
func getResponseCombiner() ResponseCombiner {
	combiner, _ := responseCombiners.GetResponseCombiner(defaultCombinerName)
	return combiner
}

// shouldRunSequentialMerger - 지정된 설정 정보를 기준으로 Merging이 순차 처리가 되어야할지 검증
func shouldRunSequentialMerger(eConf *config.EndpointConfig) bool {
	if v, ok := eConf.Middleware[MWNamespace]; ok {
		if e, ok := v.(map[string]interface{}); ok {
			if v, ok := e[sequentialKey]; ok {
				c, ok := v.(bool)
				return ok && c
			}
		}
	}
	return false
}

// parallelMerge - 지정한 시간내에 Timeout 발생하는 Context 기반으로 Request를 처리하고 도착하는 Response를 병렬로 처리
func parallelMerge(timeout time.Duration, rc ResponseCombiner, next ...Proxy) Proxy {
	return func(ctx context.Context, req *Request) (*Response, error) {
		localCtx, cancel := context.WithTimeout(ctx, timeout)

		parts := make(chan *Response, len(next))
		failed := make(chan error, len(next))

		for _, n := range next {
			go requestPart(localCtx, n, req, parts, failed)
		}

		acc := newIncrementalMergeAccumultor(len(next), rc)
		for i := 0; i < len(next); i++ {
			select {
			case err := <-failed:
				acc.Merge(nil, err)
			case response := <-parts:
				acc.Merge(response, nil)
			}
		}

		result, err := acc.Result()
		cancel()
		return result, err
	}
}

// sequentialMerge - 지정한 시간내에 Timeout 발생하는 Context 기반으로 Request를 순차적으로 처리하고 이전 Response의 결과를 파라미터로 처리해서 다음 Request를 처리하는 방식으로 순차 처리
func sequentialMerge(patterns []string, timeout time.Duration, rc ResponseCombiner, next ...Proxy) Proxy {
	return func(ctx context.Context, req *Request) (*Response, error) {
		localCtx, cancel := context.WithTimeout(ctx, timeout)

		parts := make([]*Response, len(next))
		out := make(chan *Response, 1)
		errCh := make(chan error, 1)

		acc := newIncrementalMergeAccumultor(len(next), rc)
	TxLoop:
		for i, n := range next {
			// 두번째 부터 전 호출의 결과에서 파라미터 검증
			if 0 < i {
				for _, match := range reMergeKey.FindAllStringSubmatch(patterns[i], -1) {
					if 1 < len(match) {
						rNum, err := strconv.Atoi(match[1])
						if nil != err || rNum >= i || nil == parts[rNum] {
							continue
						}
						key := "Resp" + match[1] + "_" + match[2]

						var v interface{}
						var ok bool

						data := parts[rNum].Data
						keys := strings.Split(match[2], ".")
						if 1 < len(keys) {
							for _, k := range keys[:len(keys)-1] {
								v, ok = data[k]
								if !ok {
									break
								}
								switch clean := v.(type) {
								case map[string]interface{}:
									data = clean
								default:
									break

								}
							}
						}

						v, ok = data[keys[len(keys)-1]]
						if !ok {
							continue
						}
						switch clean := v.(type) {
						case string:
							req.Params[key] = clean
						case int:
							req.Params[key] = strconv.Itoa(clean)
						case float64:
							req.Params[key] = strconv.FormatFloat(clean, 'E', -1, 32)
						case bool:
							req.Params[key] = strconv.FormatBool(clean)
						default:
							req.Params[key] = fmt.Sprintf("%v", v)
						}
					}
				}
			}
			requestPart(localCtx, n, req, out, errCh)
			select {
			case err := <-errCh:
				if 0 == i {
					cancel()
					return nil, err
				}
				acc.Merge(nil, err)
				break TxLoop
			case response := <-out:
				acc.Merge(response, nil)
				if !response.IsComplete {
					break TxLoop
				}
				parts[i] = response
			}
		}

		result, err := acc.Result()
		cancel()
		return result, err
	}
}

// ===== [ Public Functions ] =====

// NewMergeDataChain - 전달된 Endpoint 설정을 기준으로 Backend 갯수에 따라서 Response를 Merging 하는 Proxy Call chain 생성
func NewMergeDataChain(eConf *config.EndpointConfig) CallChain {
	totalBackends := len(eConf.Backend)
	if 0 == totalBackends {
		panic(ErrNoBackends)
	}
	if 1 == totalBackends {
		return EmptyChain
	}

	serviceTimeout := time.Duration(85*eConf.Timeout.Nanoseconds()/100) * time.Nanosecond
	combiner := getResponseCombiner()

	return func(next ...Proxy) Proxy {
		if len(next) != totalBackends {
			panic(ErrNotEnoughProxies)
		}
		if !shouldRunSequentialMerger(eConf) {
			return parallelMerge(serviceTimeout, combiner, next...)
		}
		patterns := make([]string, totalBackends)
		for i, b := range eConf.Backend {
			patterns[i] = b.URLPattern
		}
		return sequentialMerge(patterns, serviceTimeout, combiner, next...)
	}
}
