package metrics

import (
	"time"
)

const (
	bucketSize = 10000
)

type DataPoint struct {
	recordedAt time.Time
	value      float64
}

type Metric struct {
	cursor          int
	oneSecondCursor int
	oneSecondSum    float64
	data            [bucketSize]DataPoint
}

type MetricsRegistry struct {
	metrics map[string]*Metric
}

func New() *MetricsRegistry {
	return &MetricsRegistry{
		metrics: map[string]*Metric{},
	}
}

func (m *MetricsRegistry) Inc(name string, value float64) {
	if _, ok := m.metrics[name]; !ok {
		m.metrics[name] = &Metric{}
	}

	metric := m.metrics[name]
	dataPoint := DataPoint{recordedAt: time.Now(), value: value}
	metric.data[metric.cursor] = dataPoint
	metric.oneSecondSum += dataPoint.value

	start := metric.data[metric.oneSecondCursor]
	finish := dataPoint
	for ((finish.recordedAt.Sub(start.recordedAt)) > time.Second) && metric.oneSecondCursor != metric.cursor {
		metric.oneSecondSum -= start.value
		metric.oneSecondCursor = (metric.oneSecondCursor + 1) % bucketSize
		start = metric.data[metric.oneSecondCursor]
	}

	metric.cursor = (metric.cursor + 1) % bucketSize
}

func (m *MetricsRegistry) GetOneSecondSum(name string) float64 {
	metric := m.metrics[name]
	if _, ok := m.metrics[name]; !ok {
		return 0
	}

	start := metric.data[metric.oneSecondCursor]
	for ((time.Since(start.recordedAt)) > time.Second) && metric.oneSecondCursor != metric.cursor {
		metric.oneSecondSum -= start.value
		metric.oneSecondCursor = (metric.oneSecondCursor + 1) % bucketSize
		start = metric.data[metric.oneSecondCursor]
	}
	return metric.oneSecondSum
}
