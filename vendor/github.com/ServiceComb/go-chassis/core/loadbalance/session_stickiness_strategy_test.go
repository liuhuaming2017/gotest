package loadbalance_test

import (
	"testing"

	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/loadbalance"
	"github.com/ServiceComb/go-chassis/core/registry"
	"github.com/ServiceComb/go-chassis/session"
	"github.com/ServiceComb/go-chassis/third_party/forked/go-micro/selector"
	"github.com/stretchr/testify/assert"
	"time"
)

func TestGetSuccessiveFailureCount(t *testing.T) {
	c := loadbalance.GetSuccessiveFailureCount("0807040b-0f08-4609-4608-010c00050e03")
	assert.Equal(t, 0, c)
	loadbalance.IncreaseSuccessiveFailureCount("0807040b-0f08-4609-4608-010c00050e03")
	c = loadbalance.GetSuccessiveFailureCount("0807040b-0f08-4609-4608-010c00050e03")
	assert.Equal(t, 1, c)
	loadbalance.IncreaseSuccessiveFailureCount("0807040b-0f08-4609-4608-010c00050e03")
	c = loadbalance.GetSuccessiveFailureCount("0807040b-0f08-4609-4608-010c00050e03")
	assert.Equal(t, 2, c)
	loadbalance.DeleteSuccessiveFailureCount("0807040b-0f08-4609-4608-010c00050e03")
	c = loadbalance.GetSuccessiveFailureCount("0807040b-0f08-4609-4608-010c00050e03")
	assert.Equal(t, 0, c)
	loadbalance.ResetSuccessiveFailureMap()
	c = loadbalance.GetSuccessiveFailureCount("0807040b-0f08-4609-4608-010c00050e03")
	assert.Equal(t, 0, c)
}
func TestSessionStickyStrategies(t *testing.T) {
	config.Init()
	testData := []*registry.MicroServiceInstance{
		{
			EndpointsMap: map[string]string{"rest": "127.0.0.1:8080", "highway": "10.0.0.3:8080"},
		},
		{
			EndpointsMap: map[string]string{"rest": "127.0.0.1:8080", "highway": "10.0.0.3:8080"},
		},
	}

	for name, strategy := range map[string]selector.Strategy{"sessionstickiness": loadbalance.SessionStickiness} {

		next := strategy(testData, nil)
		counts := make(map[string]int)

		for i := 0; i < 100; i++ {
			node, err := next()
			if err != nil {
				t.Fatal(err)
			}
			counts[node.InstanceID]++
		}

		t.Logf("%s: %+v", name, counts)
	}
}
func TestStickySessionStrategy(t *testing.T) {
	config.Init()

	testData := []*registry.MicroServiceInstance{
		{
			EndpointsMap: map[string]string{"rest": "127.0.0.1:8080", "highway": "10.0.0.3:8080"},
		},
		{
			EndpointsMap: map[string]string{"rest": "127.0.0.1:8080", "highway": "10.0.0.3:8080"},
		},
	}

	for name, strategy := range map[string]selector.Strategy{"sessionstickiness": loadbalance.SessionStickiness} {
		session.Save("sticky1", "sdhgfa", time.Second*10)
		next := strategy(testData, "sticky1")

		for i := 0; i < 100; i++ {
			_, err := next()
			if err != nil {
				t.Fatal(err)
			}
		}

		t.Logf("%s", name)
		next1 := strategy(testData, "sticky1")

		for i := 0; i < 100; i++ {
			_, err := next1()
			if err != nil {
				t.Fatal(err)
			}
		}

		t.Logf("%s", name)
	}
	for name, strategy := range map[string]selector.Strategy{"sessionstickiness": loadbalance.SessionStickiness} {
		LBstr := make(map[string]string)
		LBstr["name"] = "SessionStickiness"
		LBstr["sessionTimeoutInSeconds"] = "30"
		config.GetLoadBalancing().Strategy = LBstr
		next := strategy(testData, "sticky3")

		for i := 0; i < 100; i++ {
			_, err := next()
			if err != nil {
				t.Fatal(err)
			}
		}

		t.Logf("%s", name)

		next1 := strategy(testData, "sticky3")

		for i := 0; i < 100; i++ {
			_, err := next1()
			if err != nil {
				t.Fatal(err)
			}
		}

		t.Logf("%s", name)
	}
}
