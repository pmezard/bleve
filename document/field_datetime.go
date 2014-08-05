//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.
package document

import (
	"fmt"
	"time"

	"github.com/couchbaselabs/bleve/analysis"
	"github.com/couchbaselabs/bleve/numeric_util"
)

const DEFAULT_DATETIME_INDEXING_OPTIONS = INDEX_FIELD

const DEFAULT_DATETIME_PRECISION_STEP uint = 4

type DateTimeField struct {
	name    string
	options IndexingOptions
	value   numeric_util.PrefixCoded
}

func (n *DateTimeField) Name() string {
	return n.name
}

func (n *DateTimeField) Options() IndexingOptions {
	return n.options
}

func (n *DateTimeField) Analyze() (int, analysis.TokenFrequencies) {
	tokens := make(analysis.TokenStream, 0)
	tokens = append(tokens, &analysis.Token{
		Start:    0,
		End:      len(n.value),
		Term:     n.value,
		Position: 1,
		Type:     analysis.DateTime,
	})

	original, err := n.value.Int64()
	if err == nil {

		shift := DEFAULT_PRECISION_STEP
		for shift < 64 {
			shiftEncoded, err := numeric_util.NewPrefixCodedInt64(original, shift)
			if err != nil {
				break
			}
			token := analysis.Token{
				Start:    0,
				End:      len(shiftEncoded),
				Term:     shiftEncoded,
				Position: 1,
				Type:     analysis.DateTime,
			}
			tokens = append(tokens, &token)
			shift += DEFAULT_PRECISION_STEP
		}
	}

	fieldLength := len(tokens)
	tokenFreqs := analysis.TokenFrequency(tokens)
	return fieldLength, tokenFreqs
}

func (n *DateTimeField) Value() []byte {
	return n.value
}

func (n *DateTimeField) GoString() string {
	return fmt.Sprintf("&document.DateField{Name:%s, Options: %s, Value: %s}", n.name, n.options, n.value)
}

func NewDateTimeField(name string, dt time.Time) *DateTimeField {
	return NewDateTimeFieldWithIndexingOptions(name, dt, DEFAULT_NUMERIC_INDEXING_OPTIONS)
}

func NewDateTimeFieldWithIndexingOptions(name string, dt time.Time, options IndexingOptions) *DateTimeField {
	dtInt64 := dt.UnixNano()
	prefixCoded := numeric_util.MustNewPrefixCodedInt64(dtInt64, 0)
	return &DateTimeField{
		name:    name,
		value:   prefixCoded,
		options: options,
	}
}