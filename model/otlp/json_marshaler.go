// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package otlp // import "go.opentelemetry.io/collector/model/otlp"

import (
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

// NewJSONTracesMarshaler returns a model.TracesMarshaler. Marshals to OTLP json bytes.
// Deprecated: [v0.49.0] Use ptrace.NewJSONMarshaler instead.
var NewJSONTracesMarshaler = ptrace.NewJSONMarshaler

// NewJSONMetricsMarshaler returns a model.MetricsMarshaler. Marshals to OTLP json bytes.
// Deprecated: [v0.49.0] Use pmetric.NewJSONMarshaler instead.
var NewJSONMetricsMarshaler = pmetric.NewJSONMarshaler

// NewJSONLogsMarshaler returns a model.LogsMarshaler. Marshals to OTLP json bytes.
// Deprecated: [v0.49.0] Use plog.NewJSONMarshaler instead.
var NewJSONLogsMarshaler = plog.NewJSONMarshaler