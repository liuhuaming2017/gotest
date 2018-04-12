package selector_test

// Forked from github.com/micro/go-micro
// Some parts of this file have been modified to make it functional in this package
import (
	"github.com/ServiceComb/go-chassis/core/registry"
	"github.com/ServiceComb/go-chassis/third_party/forked/go-micro/selector"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilterMD(t *testing.T) {
	t.Log("testing filter md with specified labels")
	testData := []*registry.MicroServiceInstance{
		{
			EndpointsMap: map[string]string{"w": "rest:127.0.0.1", "q": "highway:10.0.0.3"},
			Metadata:     map[string]string{"key": "1"},
		},
		{
			EndpointsMap: map[string]string{"w": "rest:10.0.0.1", "q": "highway:10.0.0.3"},
		},
	}
	f := selector.FilterMD("key", "1")
	instances := f(testData)
	assert.Equal(t, 1, len(instances))
	ins := instances[0]
	assert.Equal(t, "1", ins.Metadata["key"])
}
