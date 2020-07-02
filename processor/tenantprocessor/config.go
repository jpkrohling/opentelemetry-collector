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
	"go.opentelemetry.io/collector/config/configmodels"
)

// Config is the configuration for the span processor.
// Prior to any actions being applied, each span is compared against
// the include properties and then the exclude properties if they are specified.
// This determines if a span is to be processed or not.
type Config struct {
	configmodels.ProcessorSettings `mapstructure:",squash"`

	FromBearerToken BearerToken `mapstructure:"from_bearer_token"`
}

// BearerToken defines that the tenant information originates from a bearer token
type BearerToken struct {
	// Header specifies which context header to use to determine the tenant information
	Header string `mapstructure:"header"`

	// FromAttributes determines the JWT attributes to use to determine the tenant information.
	// If the attribute contains a dot and no such attribute exists, the processor attempts to parse it as
	// a json path
	FromAttributes []string `mapstructure:"from_attributes"`

	// Format determines how the tenant value should be formatted, following the semantics from the "fmt" package.
	Format string `mapstructure:"format"`
}
