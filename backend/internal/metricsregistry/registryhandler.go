package metricsregistry

import (
	"github.com/gin-gonic/gin"
)

func (mr *MetricsRegistry) toSerializable() interface{} {
	entries := make([]interface{}, len(mr.entries))
	for i, entry := range mr.entries {
		entries[i] = entry.toSerializable()
	}
	return entries
}

func (mre *MetricsRegistryEntry) toSerializable() interface{} {

	params := make([]interface{}, len(mre.params))
	for i, param := range mre.params {
		params[i] = map[string]interface{}{
			"id":          param.id,
			"name":        param.name,
			"description": param.description,
		}
	}

	serializable := map[string]interface{}{
		"name":        mre.name,
		"description": mre.description,
		"path":        mre.path,
		"metric_type": mre.metricType,
		"params":      params,
	}
	return serializable

}

func (mr *MetricsRegistry) GetRegistry(c *gin.Context) {
	res := mr.toSerializable()
	c.JSON(200, res)
}

func (mr *MetricsRegistry) GetParamOptions(c *gin.Context) {
	paramName := c.Param("param")
	if optsFunc, exists := mr.paramOptsFuncs[paramName]; exists {
		ctx := c.Request.Context()
		opts, err := optsFunc(ctx)
		if err != nil {
			c.JSON(500, map[string]string{"error": "Failed to retrieve parameter options"})
			return
		}
		c.JSON(200, opts)
	} else {
		c.JSON(404, map[string]string{"error": "Parameter not found"})
	}
}
