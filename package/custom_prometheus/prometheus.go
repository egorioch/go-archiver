package custom_prometheus

import "github.com/prometheus/client_golang/prometheus"

/*
	Код, созданный с целью оптимизации введения новых метрик-криэйтеров
*/

func CreateCustomCounter(name, help string, labelsForClassification []string) *prometheus.CounterVec {
	newCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: name,
			Help: help,
		},
		labelsForClassification, // Метки для классификации ошибок
	)

	return newCounter
}
