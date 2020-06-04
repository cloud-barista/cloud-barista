// CB-Store is a common repository for managing Meta Info of Cloud-Barista.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
//
// by powerkim@etri.re.kr, 2019.08.

package etcddriver

import (
	"go.etcd.io/etcd/clientv3"
	"context"
	"time"
	"strconv"

	"github.com/cloud-barista/cb-store/config"
	icbs "github.com/cloud-barista/cb-store/interfaces"
)

var cli *clientv3.Client
var ctx context.Context

type ETCDDriver struct{}

func init() {
	etcdServerPort := config.GetConfigInfos().ETCD.ETCDSERVERPORT  
//	config.Cblogger.Info("serverPort:" + etcdServerPort)

        etcdcli, err := clientv3.New(clientv3.Config{
                Endpoints:   []string{"http://" + etcdServerPort}, // @TODO set multiple Server
                DialTimeout: 5 * time.Second,
        })
	
	//config.Cblogger.Infof("etcdcli: %#v",  etcdcli)

	if err != nil {
                config.Cblogger.Error(err)
        }

        cli = etcdcli
	ctx = context.Background()
}

func (etcdDriver *ETCDDriver) InitDB() error {
        config.Cblogger.Info("init db")

        _, err := cli.Delete(ctx, "/", clientv3.WithPrefix()) // @todo

        if err != nil {
                config.Cblogger.Error(err)
        }

        return err
        return nil
}

func (etcdDriver *ETCDDriver) InitData() error {
        config.Cblogger.Info("init data")

        _, err := cli.Delete(ctx, "/", clientv3.WithPrefix())

        if err != nil {
                config.Cblogger.Error(err)
        }

        return err
	return nil
}

func (etcdDriver *ETCDDriver) Put(key string, value string) error {
	config.Cblogger.Info("Key:" + key  + ", value:" + value)

        _, err := cli.Put(ctx, key, value)
	if err != nil {
                config.Cblogger.Error(err)
		return err
        }

	return nil
}

func (etcdDriver *ETCDDriver) Get(key string) (*icbs.KeyValue, error) {
	config.Cblogger.Info("Key:" + key)

        resp, err := cli.Get(ctx, key)
        if err != nil {
                config.Cblogger.Error(err)
		return nil, err
        }

        for _, ev := range resp.Kvs {
		keyValue := icbs.KeyValue{string(ev.Key), string(ev.Value)}
                return &keyValue, nil
        }

        //return nil, fmt.Errorf("No Results with %s Key!!", key)
        return nil, nil
}

func (etcdDriver *ETCDDriver) GetList(key string, sortAscend bool) ([]*icbs.KeyValue, error) {
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
                tmpOne := icbs.KeyValue{string(ev.Key), string(ev.Value)}
                keyValueList[k] = &tmpOne
        }

        return keyValueList, nil
}

func (etcdDriver *ETCDDriver)Delete(key string) error {
	config.Cblogger.Info("Key:" + key)

	_, err := cli.Delete(ctx, key)
	
        if err != nil {
                config.Cblogger.Error(err)
        }

	return err
}

func Close() {
        cli.Close()
	cli = nil
}

