package crawler

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProxyScorer_GetBestProxyOrRetry(t *testing.T) {
	ps := NewProxyScorer()
	ps.proxies = map[string]ProxyScore{
		"proxy1": {Score: 0.8, LastUsed: time.Now().Add(-time.Hour)},
		"proxy2": {Score: 0.5, LastUsed: time.Now().Add(-time.Minute)},
		"proxy3": {Score: 0.9, LastUsed: time.Now().Add(-time.Second)},
	}

	t.Run("returns the best proxy under maxIdleTime and maxRetries", func(t *testing.T) {
		proxy, isValid, err := ps.GetBestProxyOrRetry(time.Minute, 2, 100*time.Millisecond)
		require.NoError(t, err)
		require.True(t, isValid)
		assert.Equal(t, "proxy3", proxy)
	})

	t.Run("retries and returns the best proxy with exponential backoff", func(t *testing.T) {
		ps.proxies["proxy1"].LastUsed = time.Now().Add(time.Hour)
		ps.proxies["proxy2"].LastUsed = time.Now().Add(time.Hour)
		ps.proxies["proxy3"].LastUsed = time.Now().Add(time.Hour)

		proxy, isValid, err := ps.GetBestProxyOrRetry(time.Minute, 2, 100*time.Millisecond)
		require.NoError(t, err)
		require.True(t, isValid)
		assert.Equal(t, "proxy3", proxy)
	})

	t.Run("returns ErrNoValidProxy if no valid proxy is found after retrying", func(t *testing.T) {
		for _, proxy := range ps.proxies {
			proxy.LastUsed = time.Now().Add(time.Hour)
		}

		proxy, isValid, err := ps.GetBestProxyOrRetry(time.Minute, 2, 100*time.Millisecond)
		require.EqualError(t, err, "No valid proxy found after retrying")
		require.False(t, isValid)
		require.Empty(t, proxy)
	})
}

func TestProxyScorer_retryWithBackoff(t *testing.T) {
	ps := NewProxyScorer()

	t.Run("retries with exponential backoff until succeeds or reaches maxRetries", func(t *testing.T) {
		var retries int
		var fn func() (string, bool, error)

		fn = func() (string, bool, error) {
			retries++
			if retries < 4 {
				return "", false, ErrNoValidProxy
			}
			return "test", true, nil
		}

		proxy, isValid, err := ps.retryWithBackoff(5, 100*time.Millisecond, fn)
		require.NoError(t, err)
		require.True(t, isValid)
		assert.Equal(t, 4, retries)
		assert.Equal(t, "test", proxy)
	})

	t.Run("returns ErrMaxRetriesExceeded if maxRetries is exceeded", func(t *testing.T) {
		var fn func() (string, bool, error)

		fn = func() (string, bool, error) {
			return "", false, ErrNoValidProxy
		}

		proxy, isValid, err := ps.retryWithBackoff(2, 100*time.Millisecond, fn)
		require.EqualError(t, err, "Maximum number of retries exceeded")
		require.False(t, isValid)
		require.Empty(t, proxy)
	})
}
