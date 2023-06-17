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

package interpreter

import (
	"github.com/onflow/cadence/runtime/common"
	"github.com/onflow/cadence/runtime/errors"
	"github.com/onflow/cadence/runtime/sema"
)

func NewRangeValue(
	interpreter *Interpreter,
	locationRange LocationRange,
	start IntegerValue,
	endInclusive IntegerValue,
	rangeType RangeStaticType,
) *CompositeValue {
	startComparable, startOk := start.(ComparableValue)
	endInclusiveComparable, endInclusiveOk := endInclusive.(ComparableValue)
	if !startOk || !endInclusiveOk {
		panic(errors.NewUnreachableError())
	}

	step := getValueForIntegerType(1, rangeType.ElementType)
	if startComparable.Greater(interpreter, endInclusiveComparable, locationRange) {
		negatedStep, ok := step.Negate(interpreter, locationRange).(IntegerValue)
		if !ok {
			panic(errors.NewUnreachableError())
		}

		step = negatedStep
	}

	return NewRangeValueWithStep(interpreter, locationRange, start, endInclusive, step, rangeType)
}

// NewRangeValue constructs a Range value.
func NewRangeValueWithStep(
	interpreter *Interpreter,
	locationRange LocationRange,
	start IntegerValue,
	endInclusive IntegerValue,
	step IntegerValue,
	rangeType RangeStaticType,
) *CompositeValue {

	// TODO: Validate if the sequence is moving away from the endInclusive value.
	// Also validate that step is non-zero.

	fields := []CompositeField{
		{
			Name:  sema.RangeTypeStartFieldName,
			Value: start,
		},
		{
			Name:  sema.RangeTypeEndInclusiveFieldName,
			Value: endInclusive,
		},
		{
			Name:  sema.RangeTypeStepFieldName,
			Value: step,
		},
	}

	rangeSemaType := getRangeSemaType(interpreter, rangeType)

	rangeValue := NewCompositeValueWithStaticType(
		interpreter,
		locationRange,
		sema.PublicKeyType.Location, // TODO:
		rangeSemaType.QualifiedString(),
		common.CompositeKindStructure,
		fields,
		common.ZeroAddress,
		rangeType,
	)

	rangeValue.ComputedFields = map[string]ComputedField{
		sema.RangeTypeCountFieldName: func(interpreter *Interpreter, locationRange LocationRange) Value {
			start := getFieldAsIntegerValue(rangeValue, interpreter, locationRange, sema.RangeTypeStartFieldName)
			endInclusive := getFieldAsIntegerValue(rangeValue, interpreter, locationRange, sema.RangeTypeEndInclusiveFieldName)
			step := getFieldAsIntegerValue(rangeValue, interpreter, locationRange, sema.RangeTypeStepFieldName)

			diff := convertAndAssertIntegerValue(endInclusive.Minus(interpreter, start, locationRange))

			// Perform integer division & drop the decimal part.
			// Note that step is guaranteed to be non-zero.
			return diff.Div(interpreter, step, locationRange)
		},
	}
	rangeValue.Functions = map[string]FunctionValue{
		sema.RangeTypeContainsFunctionName: NewHostFunctionValue(
			interpreter,
			sema.RangeContainsFunctionType(
				rangeSemaType.MemberType,
			),
			func(invocation Invocation) Value {
				return rangeContains(
					rangeValue,
					rangeType,
					invocation.Interpreter,
					invocation.LocationRange,
					invocation.Arguments[0],
				)
			},
		),
	}

	return rangeValue
}

func getRangeSemaType(interpreter *Interpreter, rangeType RangeStaticType) *sema.RangeType {
	return interpreter.MustConvertStaticToSemaType(rangeType).(*sema.RangeType)
}

func rangeContains(
	rangeValue *CompositeValue,
	rangeType RangeStaticType,
	interpreter *Interpreter,
	locationRange LocationRange,
	needleValue Value,
) BoolValue {
	start := getFieldAsIntegerValue(rangeValue, interpreter, locationRange, sema.RangeTypeStartFieldName)
	endInclusive := getFieldAsIntegerValue(rangeValue, interpreter, locationRange, sema.RangeTypeEndInclusiveFieldName)
	step := getFieldAsIntegerValue(rangeValue, interpreter, locationRange, sema.RangeTypeStepFieldName)

	needleInteger := convertAndAssertIntegerValue(needleValue)

	var result bool
	result = start.Equal(interpreter, locationRange, needleInteger) ||
		endInclusive.Equal(interpreter, locationRange, needleInteger)

	if !result {
		greaterThanStart := needleInteger.Greater(interpreter, start, locationRange)
		greaterThanEndInclusive := needleInteger.Greater(interpreter, endInclusive, locationRange)

		if greaterThanStart == greaterThanEndInclusive {
			// If needle is greater or smaller than both start & endInclusive, then it is outside the range.
			result = false
		} else {
			// needle is in between start and endInclusive.
			// start + k * step should be equal to needle i.e. (needle - start) mod step == 0.
			diff, ok := needleInteger.Minus(interpreter, start, locationRange).(IntegerValue)
			if !ok {
				panic(errors.NewUnreachableError())
			}

			result = diff.Mod(interpreter, step, locationRange).Equal(interpreter, locationRange, getValueForIntegerType(0, rangeType.ElementType))
		}
	}

	return AsBoolValue(result)
}

func getFieldAsIntegerValue(
	rangeValue *CompositeValue,
	interpreter *Interpreter,
	locationRange LocationRange,
	name string,
) IntegerValue {
	return convertAndAssertIntegerValue(
		rangeValue.GetField(
			interpreter,
			locationRange,
			sema.RangeTypeStartFieldName,
		),
	)
}

func convertAndAssertIntegerValue(value Value) IntegerValue {
	integerValue, ok := value.(IntegerValue)
	if !ok {
		panic(errors.NewUnreachableError())
	}
	return integerValue
}
