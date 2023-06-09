/*
 * Cadence - The resource-oriented smart contract programming language
 *
 * Copyright Dapper Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package stdlib

import (
	"fmt"

	"github.com/onflow/cadence/runtime/errors"
	"github.com/onflow/cadence/runtime/interpreter"
	"github.com/onflow/cadence/runtime/sema"
)

// RangeConstructorFunction

const rangeConstructorFunctionDocString = `
 Constructs a Range covering from start to endInclusive.
 
 The step argument is optional and determines the step size. 
 If not provided, the value of +1 or -1 is used based on the values of start and endInclusive. 
 `

var rangeConstructorFunctionType = func() *sema.FunctionType {
	typeParameter := &sema.TypeParameter{
		Name:      "T",
		TypeBound: sema.IntegerType,
	}

	typeAnnotation := sema.NewTypeAnnotation(
		&sema.GenericType{
			TypeParameter: typeParameter,
		},
	)

	return &sema.FunctionType{
		TypeParameters: []*sema.TypeParameter{
			typeParameter,
		},
		Parameters: []sema.Parameter{
			{
				Label:          sema.ArgumentLabelNotRequired,
				Identifier:     "start",
				TypeAnnotation: typeAnnotation,
			},
			{
				Label:          sema.ArgumentLabelNotRequired,
				Identifier:     "endInclusive",
				TypeAnnotation: typeAnnotation,
			},
			{
				Identifier:     "step",
				TypeAnnotation: typeAnnotation,
			},
		},
		ReturnTypeAnnotation: sema.NewTypeAnnotation(
			&sema.RangeType{
				MemberType: typeAnnotation.Type,
			},
		),
		RequiredArgumentCount: sema.RequiredArgumentCount(2),
	}
}()

var RangeConstructorFunction = NewStandardLibraryFunction(
	"Range",
	rangeConstructorFunctionType,
	rangeConstructorFunctionDocString,
	func(invocation interpreter.Invocation) interpreter.Value {
		panic("TODO")
	},
)

// RangeConstructionError

type RangeConstructionError struct {
	interpreter.LocationRange
	Message string
}

var _ errors.UserError = RangeConstructionError{}

func (RangeConstructionError) IsUserError() {}

func (e RangeConstructionError) Error() string {
	const message = "Range construction failed"
	if e.Message == "" {
		return message
	}
	return fmt.Sprintf("%s: %s", message, e.Message)
}
