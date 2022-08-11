// Copyright  The OpenTelemetry Authors
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

package transformprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor"

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor/processorhelper"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor/internal/logs"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor/internal/metrics"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor/internal/traces"
)

const (
	typeStr   = "transform"
	stability = component.StabilityLevelAlpha
)

var processorCapabilities = consumer.Capabilities{MutatesData: true}

func NewFactory() component.ProcessorFactory {
	return component.NewProcessorFactory(
		typeStr,
		createDefaultConfig,
		component.WithLogsProcessor(createLogsProcessor, stability),
		component.WithTracesProcessor(createTracesProcessor, stability),
		component.WithMetricsProcessor(createMetricsProcessor, stability),
	)
}

func createDefaultConfig() config.Processor {
	return &Config{
		ProcessorSettings: config.NewProcessorSettings(config.NewComponentID(typeStr)),
		Logs: SignalConfig{
			Queries: []string{},

			functions: logs.DefaultFunctions(),
		},
		Traces: SignalConfig{
			Queries: []string{},

			functions: traces.DefaultFunctions(),
		},
		Metrics: SignalConfig{
			Queries: []string{},

			functions: metrics.DefaultFunctions(),
		},
	}
}

func createLogsProcessor(
	ctx context.Context,
	set component.ProcessorCreateSettings,
	cfg config.Processor,
	nextConsumer consumer.Logs,
) (component.LogsProcessor, error) {
	oCfg := cfg.(*Config)

	proc, err := logs.NewProcessor(oCfg.Logs.Queries, oCfg.Logs.functions, set)
	if err != nil {
		return nil, fmt.Errorf("invalid config for \"transform\" processor %w", err)
	}
	return processorhelper.NewLogsProcessorWithCreateSettings(
		ctx,
		set,
		cfg,
		nextConsumer,
		proc.ProcessLogs,
		processorhelper.WithCapabilities(processorCapabilities))
}

func createTracesProcessor(
	ctx context.Context,
	set component.ProcessorCreateSettings,
	cfg config.Processor,
	nextConsumer consumer.Traces,
) (component.TracesProcessor, error) {
	oCfg := cfg.(*Config)

	proc, err := traces.NewProcessor(oCfg.Traces.Queries, oCfg.Traces.functions, set)
	if err != nil {
		return nil, fmt.Errorf("invalid config for \"transform\" processor %w", err)
	}
	return processorhelper.NewTracesProcessorWithCreateSettings(
		ctx,
		set,
		cfg,
		nextConsumer,
		proc.ProcessTraces,
		processorhelper.WithCapabilities(processorCapabilities))
}

func createMetricsProcessor(
	ctx context.Context,
	set component.ProcessorCreateSettings,
	cfg config.Processor,
	nextConsumer consumer.Metrics,
) (component.MetricsProcessor, error) {
	oCfg := cfg.(*Config)

	proc, err := metrics.NewProcessor(oCfg.Metrics.Queries, oCfg.Metrics.functions, set)
	if err != nil {
		return nil, fmt.Errorf("invalid config for \"transform\" processor %w", err)
	}
	return processorhelper.NewMetricsProcessorWithCreateSettings(
		ctx,
		set,
		cfg,
		nextConsumer,
		proc.ProcessMetrics,
		processorhelper.WithCapabilities(processorCapabilities))
}