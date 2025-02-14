package indexer

import (
	"github.com/cosmostation/cvms/internal/common"
	"github.com/cosmostation/cvms/internal/common/api"
	"github.com/prometheus/client_golang/prometheus"
)

func (idx *AxelarAmplifierVerifierIndexer) initLabelsAndMetrics() {
	idx.MetricsMap[common.IndexPointerBlockHeightMetricName] = idx.Factory.NewGauge(prometheus.GaugeOpts{
		Namespace:   common.Namespace,
		Subsystem:   subsystem,
		Name:        common.IndexPointerBlockHeightMetricName,
		ConstLabels: idx.PackageLabels,
	})

	idx.MetricsMap[common.IndexPointerBlockTimestampMetricName] = idx.Factory.NewGauge(prometheus.GaugeOpts{
		Namespace:   common.Namespace,
		Subsystem:   subsystem,
		Name:        common.IndexPointerBlockTimestampMetricName,
		ConstLabels: idx.PackageLabels,
	})

	latestBlockHeightMetric := idx.Factory.NewGauge(prometheus.GaugeOpts{
		Namespace:   common.Namespace,
		Subsystem:   subsystem,
		Name:        common.LatestBlockHeightMetricName,
		ConstLabels: idx.PackageLabels,
	})
	latestBlockHeightMetric.Set(0)
	idx.MetricsMap[common.LatestBlockHeightMetricName] = latestBlockHeightMetric
}

func (idx *AxelarAmplifierVerifierIndexer) updatePrometheusMetrics(indexPointer int64) {
	idx.MetricsMap[common.IndexPointerBlockHeightMetricName].Set(float64(indexPointer))
	_, timestamp, _, _, _, _, err := api.GetBlock(idx.CommonClient, indexPointer)
	if err != nil {
		idx.Errorf("failed to get block %d: %s", indexPointer, err)
		return
	}
	idx.MetricsMap[common.IndexPointerBlockTimestampMetricName].Set((float64(timestamp.Unix())))
	idx.Debugf("update prometheus metrics %d height", indexPointer)
}
