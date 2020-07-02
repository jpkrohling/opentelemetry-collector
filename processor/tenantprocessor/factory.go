// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tenantprocessor

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configerror"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/consumer"
)

const (
	// typeStr is the value of "type" Span processor in the configuration.
	typeStr = "tenant"

	defaultHeader    = "authorization" // HTTP/2 has lower-case headers
	defaultAttribute = "kubernetes.io/serviceaccount/namespace"
	defaultFormat    = "%s"
)

// Factory is the factory for the Span processor.
type Factory struct {
}

var _ component.ProcessorFactory = (*Factory)(nil)

// Type gets the type of the config created by this factory.
func (f *Factory) Type() configmodels.Type {
	return typeStr
}

// CreateDefaultConfig creates the default configuration for processor.
func (f *Factory) CreateDefaultConfig() configmodels.Processor {
	return &Config{
		ProcessorSettings: configmodels.ProcessorSettings{
			TypeVal: typeStr,
			NameVal: typeStr,
		},
	}
}

// CreateTraceProcessor creates a trace processor based on this config.
func (f *Factory) CreateTraceProcessor(
	_ context.Context,
	params component.ProcessorCreateParams,
	nextConsumer consumer.TraceConsumer,
	cfg configmodels.Processor) (component.TraceProcessor, error) {

	oCfg := cfg.(*Config)

	if len(oCfg.FromBearerToken.Header) == 0 {
		oCfg.FromBearerToken.Header = defaultHeader
	}
	if len(oCfg.FromBearerToken.FromAttributes) == 0 {
		oCfg.FromBearerToken.FromAttributes = []string{defaultAttribute}
	}
	if len(oCfg.FromBearerToken.Format) == 0 {
		oCfg.FromBearerToken.Format = defaultFormat
	}

	return newTenantProcessor(nextConsumer, params, *oCfg)
}

// CreateMetricsProcessor creates a metric processor based on this config.
func (f *Factory) CreateMetricsProcessor(
	_ context.Context,
	_ component.ProcessorCreateParams,
	_ consumer.MetricsConsumer,
	_ configmodels.Processor) (component.MetricsProcessor, error) {
	// Span Processor does not support Metrics.
	return nil, configerror.ErrDataTypeIsNotSupported
}
