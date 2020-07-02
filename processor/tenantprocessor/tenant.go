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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenterror"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/consumer/pdata"
)

var (
	// ErrTokenHeaderNotFound is returned when the context header containing the token couldn't be found
	ErrTokenHeaderNotFound = errors.New("failed to find the bearer token header")

	// ErrMultipleHeadersFound is returned when more than one header was found with the given name
	ErrMultipleHeadersFound = errors.New("more than one header was found with the given name")

	// ErrInvalidJWT is returned when the bearer token has a different number of parts than expected
	ErrInvalidJWT = errors.New("the JWT found in the request is invalid, expected to have 3 parts (header, body, checksum)")
)

type tenantProcessor struct {
	logger       *zap.Logger
	nextConsumer consumer.TraceConsumer
	config       Config
}

var _ component.TraceProcessor = (*tenantProcessor)(nil)

// newTenantProcessor returns the span processor.
func newTenantProcessor(nextConsumer consumer.TraceConsumer, params component.ProcessorCreateParams, config Config) (*tenantProcessor, error) {
	if nextConsumer == nil {
		return nil, componenterror.ErrNilNextConsumer
	}

	p := &tenantProcessor{
		logger:       params.Logger,
		nextConsumer: nextConsumer,
		config:       config,
	}

	return p, nil
}

func (p *tenantProcessor) ConsumeTraces(ctx context.Context, td pdata.Traces) error {
	// for now, we only support parsing the tenant information from a bearer token
	return p.addTenantFromBearerToken(ctx, td)
}

func (p *tenantProcessor) addTenantFromBearerToken(ctx context.Context, td pdata.Traces) error {
	m, _ := metadata.FromIncomingContext(ctx)

	sToken, ok := m[p.config.FromBearerToken.Header]
	if !ok {
		return ErrTokenHeaderNotFound
	}

	if len(sToken) > 1 {
		return ErrMultipleHeadersFound
	}

	// a JWT has three parts: header, body, checksum. They are individually base64 encoded, separated by dots
	// the header contains a prefix, "Bearer ". As we don't care at all about the prefix nor the first part, we just
	// split by the dot, meaning that the first part will look like 'Bearer eyABCDEF'.
	// Note too that this parsing logic is only valid for simple JWT, which is the most common out there. For instance,
	// Kubernetes' service account token will parse correctly with this, but encrypted JWT won't. This won't make any
	// verification on the token either, meaning that we trust an auth layer runs in front of this processor, validating
	// the token.
	parts := strings.Split(sToken[0], ".")
	if len(parts) != 3 {
		return ErrInvalidJWT
	}

	encodedBody := parts[1]
	if l := len(encodedBody) % 4; l > 0 {
		encodedBody += strings.Repeat("=", 4-l)
	}

	body, err := base64.URLEncoding.DecodeString(encodedBody)
	if err != nil {
		return errors.Wrap(err, "failed to decode the JWT body")
	}

	var bodyJSON map[string]interface{}
	if err := json.Unmarshal(body, &bodyJSON); err != nil {
		return errors.Wrap(err, "failed to parse the JWT body as JSON")
	}

	var values []interface{}
	for _, attribute := range p.config.FromBearerToken.FromAttributes {
		if v, ok := bodyJSON[attribute]; ok {
			values = append(values, v)
		}

		// TODO: handle JSONPath notation, as promised by the docs
	}

	tenant := fmt.Sprintf(p.config.FromBearerToken.Format, values...)
	rss := td.ResourceSpans()
	for i := 0; i < rss.Len(); i++ {
		rs := rss.At(i)
		if rs.IsNil() {
			continue
		}

		rs.Resource().Attributes().Insert("tenant", pdata.NewAttributeValueString(tenant))
	}

	return p.nextConsumer.ConsumeTraces(ctx, td)
}

func (p *tenantProcessor) GetCapabilities() component.ProcessorCapabilities {
	return component.ProcessorCapabilities{MutatesConsumedData: true}
}

// Start is invoked during service startup.
func (p *tenantProcessor) Start(_ context.Context, _ component.Host) error {
	return nil
}

// Shutdown is invoked during service shutdown.
func (p *tenantProcessor) Shutdown(context.Context) error {
	return nil
}
