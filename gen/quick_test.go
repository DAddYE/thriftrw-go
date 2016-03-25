// Copyright (c) 2015 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package gen

import (
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
	"time"

	tc "github.com/thriftrw/thriftrw-go/gen/testdata/containers"
	ts "github.com/thriftrw/thriftrw-go/gen/testdata/structs"

	"github.com/stretchr/testify/assert"
)

func TestQuickRoundTrip(t *testing.T) {
	tests := []reflect.Type{
		reflect.TypeOf(tc.PrimitiveContainers{}),
		reflect.TypeOf(tc.PrimitiveContainersRequired{}),
		reflect.TypeOf(tc.EnumContainers{}),
		reflect.TypeOf(ts.PrimitiveRequiredStruct{}),
		reflect.TypeOf(ts.PrimitiveOptionalStruct{}),
		reflect.TypeOf(ts.Point{}),
		reflect.TypeOf(ts.Size{}),
		// TODO Uncomment once we validate required fields
		// reflect.TypeOf(testdata.Frame{}),
		// reflect.TypeOf(testdata.Edge{}),
		// reflect.TypeOf(testdata.Graph{}),
		reflect.TypeOf(ts.ContactInfo{}),
		reflect.TypeOf(ts.User{}),
	}

	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	attempts := 1000
	for _, tt := range tests {
		for i := 0; i < attempts; i++ {
			structValue, ok := quick.Value(tt, rand)
			if !ok {
				t.Fatalf("failed to generate a value for %v", tt)
			}

			wireValue := structValue.Addr().MethodByName("ToWire").Call(nil)

			parsedValue := reflect.New(tt)
			result := parsedValue.MethodByName("FromWire").Call(wireValue)
			if result[0].Interface() != nil {
				t.Fatal("failed to parse", tt, "from", wireValue[0].Interface())
			}

			assert.Equal(t, structValue.Addr().Interface(), parsedValue.Interface())
		}
	}
}
