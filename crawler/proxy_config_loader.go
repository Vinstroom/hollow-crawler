// proxy_config_loader.go

package crawler

import (
	"io/ioutil"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

// ProxyConfig represents the proxy configuration data
type ProxyConfig struct {
	Proxies []string `yaml:"proxies"`
}

// ProxyConfigLoader loads the proxy configuration from a file
type ProxyConfigLoader struct {
	filePath string
	config   *ProxyConfig
	mutex   sync.Mutex
}

// NewProxyConfigLoader initializes a new ProxyConfigLoader
func NewProxyConfigLoader(filePath string) *ProxyConfigLoader {
	return &ProxyConfigLoader{filePath: filePath}
}

// Load loads the proxy configuration
func (p *ProxyConfigLoader) Load() error {
	data, err := ioutil.ReadFile(p.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var config ProxyConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	p.mutex.Lock()
	p.config = &config
	p.mutex.Unlock()

	return nil
}

// GetProxies returns a copy of the loaded proxy configuration
func (p *ProxyConfigLoader) GetProxies() []string {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.config == nil {
		return nil
	}

	return p.config.Proxies
}
