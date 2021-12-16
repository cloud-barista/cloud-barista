// CB-Store is a common repository for managing Meta Info of Cloud-Barista.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
//
// by powerkim@etri.re.kr, 2019.08.

package etcddriver

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	// GetList 에서 ETCD 환경이 구성되지 않았을 경우에 No Response 상태 방지용.
	"github.com/etcd-io/etcd/clientv3"
	"google.golang.org/grpc"

	// "go.etcd.io/etcd/clientv3"

	"github.com/cloud-barista/cb-store/config"
	icbs "github.com/cloud-barista/cb-store/interfaces"
)

var cli *clientv3.Client
var ctx context.Context

// ETCD 환경이 제대로 구성되지 않았을 경우 처리용
var errNoETCD = errors.New("no etcd environment available to connect")

// ETCDDriver - ETCD 처리 정보 구조
type ETCDDriver struct{}

func init() {
	// 패키지 로드 검증용
	// etcdServerPort := config.GetConfigInfos().ETCD.ETCDSERVERPORT
	// //	config.Cblogger.Info("serverPort:" + etcdServerPort)

	// etcdcli, err := clientv3.New(clientv3.Config{
	// 	// Original
	// 	// Endpoints:   []string{"http://" + etcdServerPort}, // @TODO set multiple Server
	// 	// Modified by ccambomorris
	// 	Endpoints:   strings.Split(etcdServerPort, ","),
	// 	DialTimeout: 5 * time.Second,
	// })

	// //config.Cblogger.Infof("etcdcli: %#v",  etcdcli)

	// if err != nil {
	// 	config.Cblogger.Error(err)
	// }

	// cli = etcdcli
	// ctx = context.Background()
}

// InitDB - ETCD의 데이터 초기화
func (etcdDriver *ETCDDriver) InitDB() error {
	// ETCD 환경 확인
	if nil == cli {
		return errNoETCD
	}

	config.Cblogger.Info("init db")

	_, err := cli.Delete(ctx, "/", clientv3.WithPrefix()) // @todo

	if err != nil {
		config.Cblogger.Error(err)
	}

	return err
}

// InitData - ETCD의 데이터 초기화
func (etcdDriver *ETCDDriver) InitData() error {
	// ETCD 환경 확인
	if nil == cli {
		return errNoETCD
	}

	config.Cblogger.Info("init data")

	_, err := cli.Delete(ctx, "/", clientv3.WithPrefix())

	if err != nil {
		config.Cblogger.Error(err)
	}

	return err
}

// Put - 지정한 키/값을 ETCD의 데이터로 추가
func (etcdDriver *ETCDDriver) Put(key string, value string) error {
	// ETCD 환경 확인
	if nil == cli {
		return errNoETCD
	}

	config.Cblogger.Info("Key:" + key + ", value:" + value)

	_, err := cli.Put(ctx, key, value)
	if err != nil {
		config.Cblogger.Error(err)
		return err
	}

	return nil
}

// Get - 지정한 키의 값을 ETCD의 데이터에서 추출
func (etcdDriver *ETCDDriver) Get(key string) (*icbs.KeyValue, error) {
	// ETCD 환경 확인
	if nil == cli {
		return nil, errNoETCD
	}

	config.Cblogger.Info("Key:" + key)

	resp, err := cli.Get(ctx, key)
	if err != nil {
		config.Cblogger.Error(err)
		return nil, err
	}

	for _, ev := range resp.Kvs {
		keyValue := icbs.KeyValue{Key: string(ev.Key), Value: string(ev.Value)}
		return &keyValue, nil
	}

	//return nil, fmt.Errorf("No Results with %s Key!!", key)
	return nil, nil
}

// GetList - 지정한 키와 정렬 조건을 기준으로 ETCD의 데이터에서 추출
func (etcdDriver *ETCDDriver) GetList(key string, sortAscend bool) ([]*icbs.KeyValue, error) {
	// ETCD 환경 확인
	if nil == cli {
		return nil, errNoETCD
	}

	config.Cblogger.Info("Key:" + key + ", sortAscend:" + strconv.FormatBool(sortAscend))

	order := clientv3.SortAscend

	if !sortAscend {
		order = clientv3.SortDescend
	}
	resp, err := cli.Get(ctx, key, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, order))

	if err != nil {
		config.Cblogger.Error(err)
		return nil, err
	}

	keyValueList := make([]*icbs.KeyValue, len(resp.Kvs))
	for k, ev := range resp.Kvs {
		tmpOne := icbs.KeyValue{Key: string(ev.Key), Value: string(ev.Value)}
		keyValueList[k] = &tmpOne
	}

	return keyValueList, nil
}

// Delete - 지정한 키의 데이터를 ETCD의 데이터에서 삭제
func (etcdDriver *ETCDDriver) Delete(key string) error {
	// ETCD 환경 확인
	if nil == cli {
		return errNoETCD
	}

	config.Cblogger.Info("Key:" + key)

	_, err := cli.Delete(ctx, key)

	if err != nil {
		config.Cblogger.Error(err)
	}

	return err
}

// Close - ETCD 클라이언트 종료
func (etcdDriver *ETCDDriver) Close() error {
	// ETCD 환경 확인
	if nil != cli {
		cli.Close()
		cli = nil
	}

	return nil
}

// InitializeDriver - ETCD Driver 초기화
func InitializeDriver() {
	etcdServerPort := config.GetConfigInfos().ETCD.ETCDSERVERPORT
	//	config.Cblogger.Info("serverPort:" + etcdServerPort)

	etcdcli, err := clientv3.New(clientv3.Config{
		// ETCD Endpoints에 지정된 Server 들을 CLI 로 지정하는 방식 적용 (Modified by ccambomorris)
		Endpoints:   strings.Split(etcdServerPort, ","),
		DialTimeout: 120 * time.Second,
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	})

	//config.Cblogger.Infof("etcdcli: %#v",  etcdcli)

	if err != nil {
		config.Cblogger.Error(err)
	}

	cli = etcdcli
	ctx = context.Background()
}
