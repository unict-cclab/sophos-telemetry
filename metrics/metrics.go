package metrics

import (
	"context"
	"fmt"
	"os"
	"time"

	prometheusapi "github.com/prometheus/client_golang/api"
	prometheus "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

var prometheusAddress = os.Getenv("PROMETHEUS_ADDRESS")

func newPrometheusClient(serverAddress string) (prometheus.API, error) {
	client, err := prometheusapi.NewClient(prometheusapi.Config{
		Address: serverAddress,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating metrics client: %v", err)
	}
	return prometheus.NewAPI(client), nil
}

func GetAppRequestsPerSecond(appGroupName, appName, rangeWidth string) (model.Vector, prometheus.Warnings, error) {
	prometheusClient, err := newPrometheusClient(prometheusAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create metrics client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := prometheusClient.Query(ctx, `
		sum(rate(istio_requests_total{app_group="`+appGroupName+`", app="`+appName+`", source_app!="unknown", destination_app!="unknown"}[`+rangeWidth+`])) by (source_app, destination_app)
	`, time.Now())

	if err != nil {
		return nil, nil, fmt.Errorf("error during query execution: %v", err)
	}

	vector, ok := result.(model.Vector)

	if !ok {
		return nil, nil, fmt.Errorf("query result is not a vector: %v", err)
	}

	return vector, warnings, err
}

func GetAppsRequestsPerSecond(appGroupName, rangeWidth string) (model.Vector, prometheus.Warnings, error) {
	prometheusClient, err := newPrometheusClient(prometheusAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create metrics client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := prometheusClient.Query(ctx, `
		sum(rate(istio_requests_total{reporter="source", app_group="`+appGroupName+`", source_app!="unknown", destination_app!="unknown"}[`+rangeWidth+`])) by (source_app, destination_app)
	`, time.Now())

	if err != nil {
		return nil, nil, fmt.Errorf("error during query execution: %v", err)
	}

	vector, ok := result.(model.Vector)

	if !ok {
		return nil, nil, fmt.Errorf("query result is not a vector: %v", err)
	}

	return vector, warnings, err
}

func GetAppTraffic(appGroupName, appName, rangeWidth string) (model.Vector, prometheus.Warnings, error) {
	prometheusClient, err := newPrometheusClient(prometheusAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create metrics client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := prometheusClient.Query(ctx, `
		sum(
			rate(istio_request_bytes_sum{app_group="`+appGroupName+`", app="`+appName+`", source_app!="unknown", destination_app!="unknown"}[`+rangeWidth+`])
			+
			rate(istio_response_bytes_sum{app_group="`+appGroupName+`", app="`+appName+`", source_app!="unknown", destination_app!="unknown"}[`+rangeWidth+`])
		) by (source_app, destination_app)
		or 
		sum(
			rate(istio_tcp_sent_bytes_total{app_group="`+appGroupName+`", app="`+appName+`", source_app!="unknown", destination_app!="unknown"}[`+rangeWidth+`]) 
			+ 
			rate(istio_tcp_received_bytes_total{app_group="`+appGroupName+`", app="`+appName+`", source_app!="unknown", destination_app!="unknown"}[`+rangeWidth+`])
		) by (source_app, destination_app)
	`, time.Now())

	if err != nil {
		return nil, nil, fmt.Errorf("error during query execution: %v", err)
	}

	vector, ok := result.(model.Vector)

	if !ok {
		return nil, nil, fmt.Errorf("query result is not a vector: %v", err)
	}

	return vector, warnings, err
}

func GetAppsTraffic(appGroupName, rangeWidth string) (model.Vector, prometheus.Warnings, error) {
	prometheusClient, err := newPrometheusClient(prometheusAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create metrics client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := prometheusClient.Query(ctx, `
		sum(
			rate(istio_request_bytes_sum{reporter="source", app_group="`+appGroupName+`", source_app!="unknown", destination_app!="unknown"}[`+rangeWidth+`])
			+
			rate(istio_response_bytes_sum{reporter="source", app_group="`+appGroupName+`", source_app!="unknown", destination_app!="unknown"}[`+rangeWidth+`])
		) by (source_app, destination_app)
		or 
		sum(
			rate(istio_tcp_sent_bytes_total{reporter="source", app_group="`+appGroupName+`", source_app!="unknown", destination_app!="unknown"}[`+rangeWidth+`]) 
			+ 
			rate(istio_tcp_received_bytes_total{reporter="source", app_group="`+appGroupName+`", source_app!="unknown", destination_app!="unknown"}[`+rangeWidth+`])
		) by (source_app, destination_app)
	`, time.Now())

	if err != nil {
		return nil, nil, fmt.Errorf("error during query execution: %v", err)
	}

	vector, ok := result.(model.Vector)

	if !ok {
		return nil, nil, fmt.Errorf("query result is not a vector: %v", err)
	}

	return vector, warnings, err
}
func GetAppCpuUsage(appGroupName, appName, rangeWidth string) (model.Vector, prometheus.Warnings, error) {
	prometheusClient, err := newPrometheusClient(prometheusAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create metrics client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	containerName := appGroupName + "-" + appName
	result, warnings, err := prometheusClient.Query(ctx, `
		avg by(container) (rate(container_cpu_usage_seconds_total{container="`+containerName+`"}[`+rangeWidth+`])) * 1000
	`, time.Now())

	if err != nil {
		return nil, nil, fmt.Errorf("error during query execution: %v", err)
	}

	vector, ok := result.(model.Vector)

	if !ok {
		return nil, nil, fmt.Errorf("query result is not a vector: %v", err)
	}

	return vector, warnings, err
}

func GetAppsCpuUsage(appGroupName, rangeWidth string) (model.Vector, prometheus.Warnings, error) {
	prometheusClient, err := newPrometheusClient(prometheusAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create metrics client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := prometheusClient.Query(ctx, `
		avg by(container) (rate(container_cpu_usage_seconds_total{container=~"`+appGroupName+`-.*"}[`+rangeWidth+`])) * 1000
	`, time.Now())

	if err != nil {
		return nil, nil, fmt.Errorf("error during query execution: %v", err)
	}

	vector, ok := result.(model.Vector)

	if !ok {
		return nil, nil, fmt.Errorf("query result is not a vector: %v", err)
	}

	return vector, warnings, err
}

func GetAppMemoryUsage(appGroupName, appName, rangeWidth string) (model.Vector, prometheus.Warnings, error) {
	prometheusClient, err := newPrometheusClient(prometheusAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create metrics client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	containerName := appGroupName + "-" + appName
	result, warnings, err := prometheusClient.Query(ctx, `
		avg by(container) (avg_over_time(container_memory_working_set_bytes{container="`+containerName+`"}[`+rangeWidth+`]) / (1024 * 1024))
	`, time.Now())

	if err != nil {
		return nil, nil, fmt.Errorf("error during query execution: %v", err)
	}

	vector, ok := result.(model.Vector)

	if !ok {
		return nil, nil, fmt.Errorf("query result is not a vector: %v", err)
	}

	return vector, warnings, err
}

func GetAppsMemoryUsage(appGroupName, rangeWidth string) (model.Vector, prometheus.Warnings, error) {
	prometheusClient, err := newPrometheusClient(prometheusAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create metrics client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := prometheusClient.Query(ctx, `
		avg by(container) (avg_over_time(container_memory_working_set_bytes{container=~"`+appGroupName+`-.*"}[`+rangeWidth+`]) / (1024 * 1024))
	`, time.Now())

	if err != nil {
		return nil, nil, fmt.Errorf("error during query execution: %v", err)
	}

	vector, ok := result.(model.Vector)

	if !ok {
		return nil, nil, fmt.Errorf("query result is not a vector: %v", err)
	}

	return vector, warnings, err
}

func GetAppNetworkBandwidthUsage(appGroupName, appName, rangeWidth string) (model.Vector, prometheus.Warnings, error) {
	prometheusClient, err := newPrometheusClient(prometheusAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create metrics client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := prometheusClient.Query(ctx, `
		avg by(interface) (rate(container_network_transmit_bytes_total{pod=~"`+appName+`-.*",interface="eth0"}[`+rangeWidth+`]) + rate(container_network_receive_bytes_total{pod=~"`+appName+`-.*",interface="eth0"}[`+rangeWidth+`])) / (1024 * 1024)
	`, time.Now())

	if err != nil {
		return nil, nil, fmt.Errorf("error during query execution: %v", err)
	}

	vector, ok := result.(model.Vector)

	if !ok {
		return nil, nil, fmt.Errorf("query result is not a vector: %v", err)
	}

	return vector, warnings, err
}

func GetAppDiskBandwidthUsage(appGroupName, appName, rangeWidth string) (model.Vector, prometheus.Warnings, error) {
	prometheusClient, err := newPrometheusClient(prometheusAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create metrics client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	containerName := appGroupName + "-" + appName
	result, warnings, err := prometheusClient.Query(ctx, `
		avg by(container) (rate(container_fs_writes_bytes_total{container="`+containerName+`"}[`+rangeWidth+`]) + rate(container_fs_reads_bytes_total{container="`+containerName+`"}[`+rangeWidth+`])) / (1024 * 1024)
	`, time.Now())

	if err != nil {
		return nil, nil, fmt.Errorf("error during query execution: %v", err)
	}

	vector, ok := result.(model.Vector)

	if !ok {
		return nil, nil, fmt.Errorf("query result is not a vector: %v", err)
	}

	return vector, warnings, err
}

func GetAppsDiskBandwidthUsage(appGroupName, rangeWidth string) (model.Vector, prometheus.Warnings, error) {
	prometheusClient, err := newPrometheusClient(prometheusAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create metrics client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := prometheusClient.Query(ctx, `
		avg by(container) (rate(container_fs_writes_bytes_total{container=~"`+appGroupName+`-.*"}[`+rangeWidth+`]) + rate(container_fs_reads_bytes_total{container=~"`+appGroupName+`-.*"}[`+rangeWidth+`])) / (1024 * 1024)
	`, time.Now())

	if err != nil {
		return nil, nil, fmt.Errorf("error during query execution: %v", err)
	}

	vector, ok := result.(model.Vector)

	if !ok {
		return nil, nil, fmt.Errorf("query result is not a vector: %v", err)
	}

	return vector, warnings, err
}

func GetNodeLatencies(nodeName, rangeWidth string) (model.Vector, prometheus.Warnings, error) {
	prometheusClient, err := newPrometheusClient(prometheusAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create metrics client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := prometheusClient.Query(ctx, `
		(rate(node_latency_sum{origin_node="`+nodeName+`"}[`+rangeWidth+`]) / rate(node_latency_count{origin_node="`+nodeName+`"}[`+rangeWidth+`])) * 1000
	`, time.Now())

	if err != nil {
		return nil, nil, fmt.Errorf("error during query execution: %v", err)
	}

	vector, ok := result.(model.Vector)

	if !ok {
		return nil, nil, fmt.Errorf("query result is not a vector: %v", err)
	}

	return vector, warnings, err
}

func GetNodesLatencies(rangeWidth string) (model.Vector, prometheus.Warnings, error) {
	prometheusClient, err := newPrometheusClient(prometheusAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create metrics client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := prometheusClient.Query(ctx, `
		(rate(node_latency_sum[`+rangeWidth+`]) / rate(node_latency_count[`+rangeWidth+`])) * 1000
	`, time.Now())

	if err != nil {
		return nil, nil, fmt.Errorf("error during query execution: %v", err)
	}

	vector, ok := result.(model.Vector)

	if !ok {
		return nil, nil, fmt.Errorf("query result is not a vector: %v", err)
	}

	return vector, warnings, err
}

func GetNodeMemoryUsage(nodeName, rangeWidth string) (model.Vector, prometheus.Warnings, error) {
	prometheusClient, err := newPrometheusClient(prometheusAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create metrics client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := prometheusClient.Query(ctx, `
		(avg_over_time(node_memory_MemTotal_bytes{node_name="`+nodeName+`"}[`+rangeWidth+`]) - avg_over_time(node_memory_MemFree_bytes{node_name="`+nodeName+`"}[`+rangeWidth+`]) - avg_over_time(node_memory_Buffers_bytes{node_name="`+nodeName+`"}[`+rangeWidth+`]) - avg_over_time(node_memory_Cached_bytes{node_name="`+nodeName+`"}[`+rangeWidth+`])) / (1024 * 1024)
	`, time.Now())

	if err != nil {
		return nil, nil, fmt.Errorf("error during query execution: %v", err)
	}

	vector, ok := result.(model.Vector)

	if !ok {
		return nil, nil, fmt.Errorf("query result is not a vector: %v", err)
	}

	return vector, warnings, err
}

func GetNodesMemoryUsage(rangeWidth string) (model.Vector, prometheus.Warnings, error) {
	prometheusClient, err := newPrometheusClient(prometheusAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create metrics client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := prometheusClient.Query(ctx, `
		(avg_over_time(node_memory_MemTotal_bytes[`+rangeWidth+`]) - avg_over_time(node_memory_MemFree_bytes[`+rangeWidth+`]) - avg_over_time(node_memory_Buffers_bytes[`+rangeWidth+`]) - avg_over_time(node_memory_Cached_bytes[`+rangeWidth+`])) / (1024 * 1024)
	`, time.Now())

	if err != nil {
		return nil, nil, fmt.Errorf("error during query execution: %v", err)
	}

	vector, ok := result.(model.Vector)

	if !ok {
		return nil, nil, fmt.Errorf("query result is not a vector: %v", err)
	}

	return vector, warnings, err
}

func GetNodeCpuUsage(nodeName, rangeWidth string) (model.Vector, prometheus.Warnings, error) {
	prometheusClient, err := newPrometheusClient(prometheusAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create metrics client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := prometheusClient.Query(ctx, `
		sum by (node_name) (rate(node_cpu_seconds_total{node_name="`+nodeName+`",mode!="idle"}[`+rangeWidth+`]))
	`, time.Now())

	if err != nil {
		return nil, nil, fmt.Errorf("error during query execution: %v", err)
	}

	vector, ok := result.(model.Vector)

	if !ok {
		return nil, nil, fmt.Errorf("query result is not a vector: %v", err)
	}

	return vector, warnings, err
}

func GetNodesCpuUsage(rangeWidth string) (model.Vector, prometheus.Warnings, error) {
	prometheusClient, err := newPrometheusClient(prometheusAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create metrics client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := prometheusClient.Query(ctx, `
		sum by (node_name) (rate(node_cpu_seconds_total{mode!="idle"}[`+rangeWidth+`]))
	`, time.Now())

	if err != nil {
		return nil, nil, fmt.Errorf("error during query execution: %v", err)
	}

	vector, ok := result.(model.Vector)

	if !ok {
		return nil, nil, fmt.Errorf("query result is not a vector: %v", err)
	}

	return vector, warnings, err
}

func GetNodeNetworkBandwidthUsage(nodeName, rangeWidth string) (model.Vector, prometheus.Warnings, error) {
	prometheusClient, err := newPrometheusClient(prometheusAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create metrics client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := prometheusClient.Query(ctx, `
		sum by (node_name) (rate(node_network_transmit_bytes_total{node_name="`+nodeName+`"}[`+rangeWidth+`]) + rate(node_network_receive_bytes_total{node_name="`+nodeName+`"}[`+rangeWidth+`])) / (1024 * 1024)
	`, time.Now())

	if err != nil {
		return nil, nil, fmt.Errorf("error during query execution: %v", err)
	}

	vector, ok := result.(model.Vector)

	if !ok {
		return nil, nil, fmt.Errorf("query result is not a vector: %v", err)
	}

	return vector, warnings, err
}

func GetNodesNetworkBandwidthUsage(rangeWidth string) (model.Vector, prometheus.Warnings, error) {
	prometheusClient, err := newPrometheusClient(prometheusAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create metrics client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := prometheusClient.Query(ctx, `
		sum by (node_name) (rate(node_network_transmit_bytes_total[`+rangeWidth+`]) + rate(node_network_receive_bytes_total[`+rangeWidth+`])) / (1024 * 1024)
	`, time.Now())

	if err != nil {
		return nil, nil, fmt.Errorf("error during query execution: %v", err)
	}

	vector, ok := result.(model.Vector)

	if !ok {
		return nil, nil, fmt.Errorf("query result is not a vector: %v", err)
	}

	return vector, warnings, err
}

func GetNodeDiskBandwidthUsage(nodeName, rangeWidth string) (model.Vector, prometheus.Warnings, error) {
	prometheusClient, err := newPrometheusClient(prometheusAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create metrics client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := prometheusClient.Query(ctx, `
		sum by (node_name) (rate(node_disk_written_bytes_total{node_name="`+nodeName+`"}[`+rangeWidth+`]) + rate(node_disk_read_bytes_total{node_name="`+nodeName+`"}[`+rangeWidth+`])) / (1024 * 1024)
	`, time.Now())

	if err != nil {
		return nil, nil, fmt.Errorf("error during query execution: %v", err)
	}

	vector, ok := result.(model.Vector)

	if !ok {
		return nil, nil, fmt.Errorf("query result is not a vector: %v", err)
	}

	return vector, warnings, err
}

func GetNodesDiskBandwidthUsage(rangeWidth string) (model.Vector, prometheus.Warnings, error) {
	prometheusClient, err := newPrometheusClient(prometheusAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create metrics client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := prometheusClient.Query(ctx, `
		sum by (node_name) (rate(node_disk_written_bytes_total[`+rangeWidth+`]) + rate(node_disk_read_bytes_total[`+rangeWidth+`])) / (1024 * 1024)
	`, time.Now())

	if err != nil {
		return nil, nil, fmt.Errorf("error during query execution: %v", err)
	}

	vector, ok := result.(model.Vector)

	if !ok {
		return nil, nil, fmt.Errorf("query result is not a vector: %v", err)
	}

	return vector, warnings, err
}
