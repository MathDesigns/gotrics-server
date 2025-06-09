package storage

import (
	"context"
	"fmt"
	"gotrics-server/internal/config"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type InfluxStore struct {
	client   influxdb2.Client
	writeAPI api.WriteAPIBlocking
	queryAPI api.QueryAPI
}

type Metric struct {
	Hostname        string    `json:"hostname"`
	Timestamp       time.Time `json:"timestamp"`
	CPUUsagePercent float64   `json:"cpu_usage_percent"`
	MemoryUsedMB    uint64    `json:"memory_used_mb"`
	DiskUsedPercent float64   `json:"disk_used_percent"`
}

func NewInfluxStore(cfg config.InfluxDBConfig) *InfluxStore {
	client := influxdb2.NewClient(cfg.URL, cfg.Token)
	writeAPI := client.WriteAPIBlocking(cfg.Org, cfg.Bucket)
	queryAPI := client.QueryAPI(cfg.Org)
	return &InfluxStore{
		client:   client,
		writeAPI: writeAPI,
		queryAPI: queryAPI,
	}
}

func (s *InfluxStore) WriteMetric(ctx context.Context, m *Metric) error {
	p := influxdb2.NewPointWithMeasurement("system").
		AddTag("hostname", m.Hostname).
		AddField("cpu_usage_percent", m.CPUUsagePercent).
		AddField("memory_used_mb", m.MemoryUsedMB).
		AddField("disk_used_percent", m.DiskUsedPercent).
		SetTime(m.Timestamp)

	return s.writeAPI.WritePoint(ctx, p)
}

func (s *InfluxStore) GetHostMetrics(ctx context.Context, hostname string, duration time.Duration) ([]*Metric, error) {
	query := fmt.Sprintf(
		`from(bucket: "gotrics")
		|> range(start: -%s)
		|> filter(fn: (r) => r._measurement == "system" and r.hostname == "%s")
		|> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")`,
		duration.String(), hostname,
	)

	result, err := s.queryAPI.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var metrics []*Metric
	for result.Next() {
		record := result.Record()
		m := &Metric{
			Hostname:        record.ValueByKey("hostname").(string),
			Timestamp:       record.Time(),
			CPUUsagePercent: record.ValueByKey("cpu_usage_percent").(float64),
			MemoryUsedMB:    record.ValueByKey("memory_used_mb").(uint64),
			DiskUsedPercent: record.ValueByKey("disk_used_percent").(float64),
		}
		metrics = append(metrics, m)
	}
	return metrics, result.Err()
}

func (s *InfluxStore) GetKnownHosts(ctx context.Context) ([]string, error) {
	query := `import "influxdata/influxdb/schema"
		schema.tagValues(bucket: "gotrics", tag: "hostname")`

	result, err := s.queryAPI.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var hosts []string
	for result.Next() {
		hosts = append(hosts, result.Record().Value().(string))
	}
	return hosts, result.Err()
}

func (s *InfluxStore) Close() {
	s.client.Close()
}
