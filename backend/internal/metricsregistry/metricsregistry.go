package metricsregistry

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
)

type MetricsRegistry struct {
	entries        []MetricsRegistryEntry
	paramOptsFuncs map[string]func(context.Context) ([]MetricsParamEntry, error)
	rg             *gin.RouterGroup
}

func NewMetricsRegistry(rg *gin.RouterGroup) *MetricsRegistry {

	registry := &MetricsRegistry{
		rg:             rg,
		entries:        make([]MetricsRegistryEntry, 0),
		paramOptsFuncs: make(map[string]func(context.Context) ([]MetricsParamEntry, error)),
	}

	registry.rg.GET("/registry", registry.GetRegistry)
	registry.rg.GET("/registry/param-options/:param", registry.GetParamOptions)

	return registry
}

func (mr *MetricsRegistry) RegisterEntry(entry ...MetricsRegistryEntry) error {

	mr.entries = append(mr.entries, entry...)
	for _, e := range entry {
		mr.rg.GET(e.path, e.metricsFunc)
	}

	return nil
}

func (mr *MetricsRegistry) RegisterParamOptions(param MetricsParam, optsFunc func(context.Context) ([]MetricsParamEntry, error)) error {

	if _, exists := mr.paramOptsFuncs[param.id]; exists {
		return fmt.Errorf("options function for parameter %s already registered", param.id)
	}

	mr.paramOptsFuncs[param.id] = optsFunc
	return nil
}

type MetricsRegistryEntry struct {
	name          string
	description   string
	path          string
	metricType    MetricType
	params        []MetricsParam
	displayFormat DisplayFormat
	metricsFunc   func(*gin.Context)
}

func NewMetricsRegistryEntry(name, description, path string, metricType MetricType, metricParams []MetricsParam, displayFormat DisplayFormat, metricsFunc func(*gin.Context)) (*MetricsRegistryEntry, error) {

	params := templateParams(path)

	if len(params) != len(metricParams) {
		return nil, fmt.Errorf("mismatch between path parameters and metric parameters")
	}

	paramMap := make(map[string]struct{})
	for _, param := range metricParams {
		paramMap[param.id] = struct{}{}
	}

	for _, param := range params {
		if _, exists := paramMap[param]; !exists {
			return nil, fmt.Errorf("parameter %s is in the path but not in metricParams", param)
		}
	}

	return &MetricsRegistryEntry{
		name:          name,
		description:   description,
		path:          path,
		metricType:    metricType,
		params:        metricParams,
		displayFormat: displayFormat,
		metricsFunc:   metricsFunc,
	}, nil
}

type MetricsParam struct {
	id          string
	name        string
	description string
}

func NewMetricsParam(id, name, description string) *MetricsParam {
	return &MetricsParam{
		id:          id,
		name:        name,
		description: description,
	}
}

type MetricsParamEntry struct {
	DisplayName string `json:"display_name"`
	Value       string `json:"value"`
}
