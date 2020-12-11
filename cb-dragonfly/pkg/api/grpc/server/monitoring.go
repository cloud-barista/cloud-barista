package server

import (
	"context"
	"net/http"
	"time"

	"github.com/cloud-barista/cb-dragonfly/pkg/api/grpc/common"
	pb "github.com/cloud-barista/cb-dragonfly/pkg/api/grpc/protobuf/cbdragonfly"
	"github.com/cloud-barista/cb-dragonfly/pkg/config"
	coreconfig "github.com/cloud-barista/cb-dragonfly/pkg/core/config"
	coremetric "github.com/cloud-barista/cb-dragonfly/pkg/core/metric"
)

type MonitoringService struct{}

func (c MonitoringService) GetMCISMonInfo(ctx context.Context, request *pb.VMMCISMonQryRequest) (*pb.MCISMonInfoResponse, error) {
	mcisMetric, statusCode, err := coremetric.GetMCISCommonMonInfo(request.NsId, request.McisId, request.VmId, request.AgentIp, request.MetricName)
	if statusCode != http.StatusOK {
		return nil, common.ConvGrpcStatusErr(err, "", "MonitoringService.GetMCISMonInfo()")
	}

	var resp pb.MCISMonInfoResponse
	err = common.CopySrcToDest(mcisMetric, &resp)
	if err != nil {
		return nil, common.ConvGrpcStatusErr(err, "", "MonitoringService.GetMCISMonInfo()")
	}
	return &resp, nil
}

func (c MonitoringService) GetVMOnDemandMonCpuInfo(ctx context.Context, request *pb.VMOnDemandMonQryRequest) (*pb.CpuOnDemandInfoResponse, error) {
	cpuMetric, statusCode, err := coremetric.GetVMOnDemandMonInfo(request.NsId, request.McisId, request.VmId, coremetric.Cpu, request.AgentIp)
	if statusCode != http.StatusOK {
		return nil, common.ConvGrpcStatusErr(err, "", "MonitoringService.GetVMMonCpuInfo()")
	}

	cpuMetricMap := cpuMetric.(map[string]interface{})

	// convert to grpc object
	var tagInfo *pb.Tags
	var metricInfo *pb.CpuOnDemandInfo
	err = common.CopySrcToDest(cpuMetricMap["tags"], &tagInfo)
	if err != nil {
		return nil, common.ConvGrpcStatusErr(err, "", "MonitoringService.GetVMMonCpuInfo()")
	}
	err = common.CopySrcToDest(cpuMetricMap["values"], &metricInfo)
	if err != nil {
		return nil, common.ConvGrpcStatusErr(err, "", "MonitoringService.GetVMMonCpuInfo()")
	}

	resp := &pb.CpuOnDemandInfoResponse{
		Name:   cpuMetricMap["name"].(string),
		Tags:   tagInfo,
		Time:   cpuMetricMap["time"].(time.Time).String(),
		Values: metricInfo,
	}
	return resp, nil
}

func (c MonitoringService) GetVMOnDemandMonCpuFreqInfo(ctx context.Context, request *pb.VMOnDemandMonQryRequest) (*pb.CpuFreqOnDemandInfoResponse, error) {
	cpuFreqMetric, statusCode, err := coremetric.GetVMOnDemandMonInfo(request.NsId, request.McisId, request.VmId, coremetric.CpuFreqency, request.AgentIp)
	if statusCode != http.StatusOK {
		return nil, common.ConvGrpcStatusErr(err, "", "MonitoringService.GetVMMonCpuInfo()")
	}

	cpuMetricMap := cpuFreqMetric.(map[string]interface{})

	// convert to grpc object
	var tagInfo *pb.Tags
	var metricInfo *pb.CpuFreqOnDemandInfo
	err = common.CopySrcToDest(cpuMetricMap["tags"], &tagInfo)
	if err != nil {
		return nil, common.ConvGrpcStatusErr(err, "", "MonitoringService.GetVMMonCpuInfo()")
	}
	err = common.CopySrcToDest(cpuMetricMap["values"], &metricInfo)
	if err != nil {
		return nil, common.ConvGrpcStatusErr(err, "", "MonitoringService.GetVMMonCpuInfo()")
	}

	resp := &pb.CpuFreqOnDemandInfoResponse{
		Name:   cpuMetricMap["name"].(string),
		Tags:   tagInfo,
		Time:   cpuMetricMap["time"].(time.Time).String(),
		Values: metricInfo,
	}
	return resp, nil
}

func (c MonitoringService) GetVMOnDemandMonMemoryInfo(ctx context.Context, request *pb.VMOnDemandMonQryRequest) (*pb.MemoryOnDemandInfoResponse, error) {
	memMetric, statusCode, err := coremetric.GetVMOnDemandMonInfo(request.NsId, request.McisId, request.VmId, coremetric.Memory, request.AgentIp)
	if statusCode != http.StatusOK {
		return nil, common.ConvGrpcStatusErr(err, "", "MonitoringService.GetVMMonCpuInfo()")
	}

	memMetricMap := memMetric.(map[string]interface{})

	// convert to grpc object
	var tagInfo *pb.Tags
	var metricInfo *pb.MemoryOnDemandInfo
	err = common.CopySrcToDest(memMetricMap["tags"], &tagInfo)
	if err != nil {
		return nil, common.ConvGrpcStatusErr(err, "", "MonitoringService.GetVMMonCpuInfo()")
	}
	err = common.CopySrcToDest(memMetricMap["values"], &metricInfo)
	if err != nil {
		return nil, common.ConvGrpcStatusErr(err, "", "MonitoringService.GetVMMonCpuInfo()")
	}

	resp := &pb.MemoryOnDemandInfoResponse{
		Name:   memMetricMap["name"].(string),
		Tags:   tagInfo,
		Time:   memMetricMap["time"].(time.Time).String(),
		Values: metricInfo,
	}
	return resp, nil
}

func (c MonitoringService) GetVMOnDemandMonDiskInfo(ctx context.Context, request *pb.VMOnDemandMonQryRequest) (*pb.DiskOnDemandInfoResponse, error) {
	diskMetric, statusCode, err := coremetric.GetVMOnDemandMonInfo(request.NsId, request.McisId, request.VmId, coremetric.Disk, request.AgentIp)
	if statusCode != http.StatusOK {
		return nil, common.ConvGrpcStatusErr(err, "", "MonitoringService.GetVMMonCpuInfo()")
	}

	diskMetricMap := diskMetric.(map[string]interface{})

	// convert to grpc object
	var tagInfo *pb.Tags
	var metricInfo *pb.DiskOnDemandInfo
	err = common.CopySrcToDest(diskMetricMap["tags"], &tagInfo)
	if err != nil {
		return nil, common.ConvGrpcStatusErr(err, "", "MonitoringService.GetVMMonCpuInfo()")
	}
	err = common.CopySrcToDest(diskMetricMap["values"], &metricInfo)
	if err != nil {
		return nil, common.ConvGrpcStatusErr(err, "", "MonitoringService.GetVMMonCpuInfo()")
	}

	resp := &pb.DiskOnDemandInfoResponse{
		Name:   diskMetricMap["name"].(string),
		Tags:   tagInfo,
		Time:   diskMetricMap["time"].(time.Time).String(),
		Values: metricInfo,
	}
	return resp, nil
}

func (c MonitoringService) GetVMOnDemandMonNetworkInfo(ctx context.Context, request *pb.VMOnDemandMonQryRequest) (*pb.NetworkOnDemandInfoResponse, error) {
	netMetric, statusCode, err := coremetric.GetVMOnDemandMonInfo(request.NsId, request.McisId, request.VmId, coremetric.Disk, request.AgentIp)
	if statusCode != http.StatusOK {
		return nil, common.ConvGrpcStatusErr(err, "", "MonitoringService.GetVMMonCpuInfo()")
	}

	netMetricMap := netMetric.(map[string]interface{})

	// convert to grpc object
	var tagInfo *pb.Tags
	var metricInfo *pb.NetworkOnDemandInfo
	err = common.CopySrcToDest(netMetricMap["tags"], &tagInfo)
	if err != nil {
		return nil, common.ConvGrpcStatusErr(err, "", "MonitoringService.GetVMMonCpuInfo()")
	}
	err = common.CopySrcToDest(netMetricMap["values"], &metricInfo)
	if err != nil {
		return nil, common.ConvGrpcStatusErr(err, "", "MonitoringService.GetVMMonCpuInfo()")
	}

	resp := &pb.NetworkOnDemandInfoResponse{
		Name:   netMetricMap["name"].(string),
		Tags:   tagInfo,
		Time:   netMetricMap["time"].(time.Time).String(),
		Values: metricInfo,
	}
	return resp, nil
}

func (c MonitoringService) GetVMMonCpuInfo(ctx context.Context, request *pb.VMMonQryRequest) (*pb.CpuInfoResponse, error) {
	cpuMetric, statusCode, err := coremetric.GetVMMonInfo(request.NsId, request.McisId, request.VmId, coremetric.Cpu, request.PeriodType, request.StatisticsCriteria, request.Duration)
	if statusCode != http.StatusOK {
		return nil, common.ConvGrpcStatusErr(err, "", "MonitoringService.GetVMMonCpuInfo()")
	}

	cpuMetricMap := cpuMetric.(map[string]interface{})

	// convert to grpc object
	var tagInfo *pb.Tags
	var metricInfo []*pb.CpuInfo
	err = common.CopySrcToDest(cpuMetricMap["tags"], &tagInfo)
	if err != nil {
		return nil, common.ConvGrpcStatusErr(err, "", "MonitoringService.GetVMMonCpuInfo()")
	}
	err = common.CopySrcToDest(cpuMetricMap["values"], &metricInfo)
	if err != nil {
		return nil, common.ConvGrpcStatusErr(err, "", "MonitoringService.GetVMMonCpuInfo()")
	}

	resp := &pb.CpuInfoResponse{
		Name:   cpuMetricMap["name"].(string),
		Tags:   tagInfo,
		Values: metricInfo,
	}
	return resp, nil
}

func (c MonitoringService) GetVMMonCpuFreqInfo(ctx context.Context, request *pb.VMMonQryRequest) (*pb.CpuFreqInfoResponse, error) {
	cpuFreqMetric, statusCode, err := coremetric.GetVMMonInfo(request.NsId, request.McisId, request.VmId, coremetric.CpuFreqency, request.PeriodType, request.StatisticsCriteria, request.Duration)
	cpuFreqMetricMap := cpuFreqMetric.(map[string]interface{})
	if statusCode != http.StatusOK {
		return nil, common.ConvGrpcStatusErr(err, "", "MonitoringService.GetVMMonCpuInfo()")
	}

	// convert to grpc object
	var tagInfo *pb.Tags
	var metricInfo []*pb.CpuFreqInfo
	err = common.CopySrcToDest(cpuFreqMetricMap["tags"], &tagInfo)
	if err != nil {
		return nil, common.ConvGrpcStatusErr(err, "", "MonitoringService.GetVMMonCpuInfo()")
	}
	err = common.CopySrcToDest(cpuFreqMetricMap["values"], &metricInfo)
	if err != nil {
		return nil, common.ConvGrpcStatusErr(err, "", "MonitoringService.GetVMMonCpuInfo()")
	}

	resp := &pb.CpuFreqInfoResponse{
		Name:   cpuFreqMetricMap["name"].(string),
		Tags:   tagInfo,
		Values: metricInfo,
	}
	return resp, nil
}

func (c MonitoringService) GetVMMonMemoryInfo(ctx context.Context, request *pb.VMMonQryRequest) (*pb.MemoryInfoResponse, error) {
	memMetric, statusCode, err := coremetric.GetVMMonInfo(request.NsId, request.McisId, request.VmId, coremetric.Memory, request.PeriodType, request.StatisticsCriteria, request.Duration)
	memMetricMap := memMetric.(map[string]interface{})
	if statusCode != http.StatusOK {
		return nil, common.ConvGrpcStatusErr(err, "", "MonitoringService.GetVMMonCpuInfo()")
	}

	// convert to grpc object
	var tagInfo *pb.Tags
	var metricInfo []*pb.MemoryInfo
	err = common.CopySrcToDest(memMetricMap["tags"], &tagInfo)
	if err != nil {
		return nil, common.ConvGrpcStatusErr(err, "", "MonitoringService.GetVMMonCpuInfo()")
	}
	err = common.CopySrcToDest(memMetricMap["values"], &metricInfo)
	if err != nil {
		return nil, common.ConvGrpcStatusErr(err, "", "MonitoringService.GetVMMonCpuInfo()")
	}

	resp := &pb.MemoryInfoResponse{
		Name:   memMetricMap["name"].(string),
		Tags:   tagInfo,
		Values: metricInfo,
	}
	return resp, nil
}

func (c MonitoringService) GetVMMonDiskInfo(ctx context.Context, request *pb.VMMonQryRequest) (*pb.DiskInfoResponse, error) {
	diskMetric, statusCode, err := coremetric.GetVMMonInfo(request.NsId, request.McisId, request.VmId, coremetric.Disk, request.PeriodType, request.StatisticsCriteria, request.Duration)
	diskMetricMap := diskMetric.(map[string]interface{})
	if statusCode != http.StatusOK {
		return nil, common.ConvGrpcStatusErr(err, "", "MonitoringService.GetVMMonCpuInfo()")
	}

	// convert to grpc object
	var tagInfo *pb.Tags
	var metricInfo []*pb.DiskInfo
	err = common.CopySrcToDest(diskMetricMap["tags"], &tagInfo)
	if err != nil {
		return nil, common.ConvGrpcStatusErr(err, "", "MonitoringService.GetVMMonCpuInfo()")
	}
	err = common.CopySrcToDest(diskMetricMap["values"], &metricInfo)
	if err != nil {
		return nil, common.ConvGrpcStatusErr(err, "", "MonitoringService.GetVMMonCpuInfo()")
	}

	resp := &pb.DiskInfoResponse{
		Name:   diskMetricMap["name"].(string),
		Tags:   tagInfo,
		Values: metricInfo,
	}
	return resp, nil
}

func (c MonitoringService) GetVMMonNetworkInfo(ctx context.Context, request *pb.VMMonQryRequest) (*pb.NetworkInfoResponse, error) {
	netMetric, statusCode, err := coremetric.GetVMMonInfo(request.NsId, request.McisId, request.VmId, coremetric.Network, request.PeriodType, request.StatisticsCriteria, request.Duration)
	netMetricMap := netMetric.(map[string]interface{})
	if statusCode != http.StatusOK {
		return nil, common.ConvGrpcStatusErr(err, "", "MonitoringService.GetVMMonCpuInfo()")
	}

	// convert to grpc object
	var tagInfo *pb.Tags
	var metricInfo []*pb.NetworkInfo
	err = common.CopySrcToDest(netMetricMap["tags"], &tagInfo)
	if err != nil {
		return nil, common.ConvGrpcStatusErr(err, "", "MonitoringService.GetVMMonCpuInfo()")
	}
	err = common.CopySrcToDest(netMetricMap["values"], &metricInfo)
	if err != nil {
		return nil, common.ConvGrpcStatusErr(err, "", "MonitoringService.GetVMMonCpuInfo()")
	}

	resp := &pb.NetworkInfoResponse{
		Name:   netMetricMap["name"].(string),
		Tags:   tagInfo,
		Values: metricInfo,
	}
	return resp, nil
}

func (c MonitoringService) SetMonConfig(ctx context.Context, request *pb.MonitoringConfigRequest) (*pb.MonitoringConfigResponse, error) {
	// convert grpc request to config struct
	reqParams := config.Monitoring{
		AgentInterval:     int(request.Item.AgentInterval),
		CollectorInterval: int(request.Item.CollectorInterval),
		MonitoringPolicy:  request.Item.MonitoringPolicy,
		MaxHostCount:      int(request.Item.MaxHostCount),
	}
	monConfig, statusCode, err := coreconfig.SetMonConfig(reqParams)
	if statusCode != http.StatusOK {
		return nil, common.ConvGrpcStatusErr(err, "", "MonitoringService.SetMonConfig()")
	}
	var grpcObj *pb.MonitoringConfigInfo
	err = common.CopySrcToDest(&monConfig, &grpcObj)
	if err != nil {
		return nil, common.ConvGrpcStatusErr(err, "", "MonitoringService.SetMonConfig()")
	}
	resp := &pb.MonitoringConfigResponse{Item: grpcObj}
	return resp, nil
}

func (c MonitoringService) GetMonConfig(ctx context.Context, empty *pb.Empty) (*pb.MonitoringConfigResponse, error) {
	monConfig, statusCode, err := coreconfig.GetMonConfig()
	if statusCode != http.StatusOK {
		return nil, common.ConvGrpcStatusErr(err, "", "MonitoringService.GetMonConfig()")
	}
	var grpcObj *pb.MonitoringConfigInfo
	err = common.CopySrcToDest(&monConfig, &grpcObj)
	if err != nil {
		return nil, common.ConvGrpcStatusErr(err, "", "MonitoringService.GetMonConfig()")
	}
	resp := &pb.MonitoringConfigResponse{Item: grpcObj}
	return resp, nil
}

func (c MonitoringService) ResetMonConfig(ctx context.Context, empty *pb.Empty) (*pb.MonitoringConfigResponse, error) {
	monConfig, statusCode, err := coreconfig.ResetMonConfig()
	if statusCode != http.StatusOK {
		return nil, common.ConvGrpcStatusErr(err, "", "MonitoringService.ResetMonConfig()")
	}
	var grpcObj *pb.MonitoringConfigInfo
	err = common.CopySrcToDest(&monConfig, &grpcObj)
	if err != nil {
		return nil, common.ConvGrpcStatusErr(err, "", "MonitoringService.ResetMonConfig()")
	}
	resp := &pb.MonitoringConfigResponse{Item: grpcObj}
	return resp, nil
}

func (c MonitoringService) InstallTelegraf(ctx context.Context, request *pb.InstallTelegrafRequest) (*pb.MessageResponse, error) {
	panic("implement me")
}
