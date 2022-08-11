// Copyright  OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mongodbatlasreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/mongodbatlasreceiver"

import (
	"errors"
	"fmt"
	"net"

	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/configtls"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.opentelemetry.io/collector/receiver/scraperhelper"
	"go.uber.org/multierr"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/mongodbatlasreceiver/internal/metadata"
)

var _ config.Receiver = (*Config)(nil)

type Config struct {
	scraperhelper.ScraperControllerSettings `mapstructure:",squash"`
	PublicKey                               string                       `mapstructure:"public_key"`
	PrivateKey                              string                       `mapstructure:"private_key"`
	Granularity                             string                       `mapstructure:"granularity"`
	Metrics                                 metadata.MetricsSettings     `mapstructure:"metrics"`
	Alerts                                  AlertConfig                  `mapstructure:"alerts"`
	RetrySettings                           exporterhelper.RetrySettings `mapstructure:"retry_on_failure"`
}

type AlertConfig struct {
	Enabled  bool                        `mapstructure:"enabled"`
	Endpoint string                      `mapstructure:"endpoint"`
	Secret   string                      `mapstructure:"secret"`
	TLS      *configtls.TLSServerSetting `mapstructure:"tls"`
}

var (
	errNoEndpoint = errors.New("an endpoint must be specified")
	errNoSecret   = errors.New("a webhook secret must be specified")
	errNoCert     = errors.New("tls was configured, but no cert file was specified")
	errNoKey      = errors.New("tls was configured, but no key file was specified")
)

func (c *Config) Validate() error {
	var errs error

	errs = multierr.Append(errs, c.ScraperControllerSettings.Validate())
	errs = multierr.Append(errs, c.Alerts.validate())

	return errs
}

func (a *AlertConfig) validate() error {
	var errs error

	if !a.Enabled {
		// No need to further validate, receiving alerts is disabled.
		return nil
	}

	if a.Endpoint == "" {
		errs = multierr.Append(errs, errNoEndpoint)
	}

	_, _, err := net.SplitHostPort(a.Endpoint)
	if err != nil {
		errs = multierr.Append(errs, fmt.Errorf("failed to split endpoint into 'host:port' pair: %w", err))
	}

	if a.Secret == "" {
		errs = multierr.Append(errs, errNoSecret)
	}

	if a.TLS != nil {
		if a.TLS.CertFile == "" {
			errs = multierr.Append(errs, errNoCert)
		}

		if a.TLS.KeyFile == "" {
			errs = multierr.Append(errs, errNoKey)
		}
	}

	return errs
}