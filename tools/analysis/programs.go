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

package analysis

import (
	"github.com/onflow/cadence/runtime/ast"
	"github.com/onflow/cadence/runtime/common"
	"github.com/onflow/cadence/runtime/parser"
	"github.com/onflow/cadence/runtime/sema"
	"github.com/onflow/cadence/runtime/stdlib"
)

type Programs map[common.Location]*Program

type importResolutionResults map[common.Location]bool

func (programs Programs) Load(config *Config, location common.Location) error {
	return programs.load(
		config,
		location,
		nil,
		ast.Range{},
		importResolutionResults{
			// Entry point program is also currently in check.
			location: true,
		},
	)
}

func (programs Programs) load(
	config *Config,
	location common.Location,
	importingLocation common.Location,
	importRange ast.Range,
	seenImports importResolutionResults,
) error {

	if programs[location] != nil {
		return nil
	}

	wrapError := func(err error) ParsingCheckingError {
		return ParsingCheckingError{
			error:    err,
			location: location,
		}
	}

	code, err := config.ResolveCode(location, importingLocation, importRange)
	if err != nil {
		return err
	}

	program, err := parser.ParseProgram(nil, code, parser.Config{})
	if err != nil {
		return wrapError(err)
	}

	var checker *sema.Checker
	if config.Mode&NeedTypes != 0 {
		checker, err = programs.check(config, program, location, seenImports)
		if err != nil {
			return wrapError(err)
		}
	}

	programs[location] = &Program{
		Location: location,
		Code:     code,
		Program:  program,
		Checker:  checker,
	}

	return nil
}

func (programs Programs) check(
	config *Config,
	program *ast.Program,
	location common.Location,
	seenImports importResolutionResults,
) (
	*sema.Checker,
	error,
) {
	baseValueActivation := sema.NewVariableActivation(sema.BaseValueActivation)
	for _, value := range stdlib.DefaultScriptStandardLibraryValues(nil) {
		baseValueActivation.DeclareValue(value)
	}

	checker, err := sema.NewChecker(
		program,
		location,
		nil,
		&sema.Config{
			BaseValueActivationHandler: func(_ common.Location) *sema.VariableActivation {
				return baseValueActivation
			},
			AccessCheckMode: sema.AccessCheckModeStrict,
			LocationHandler: sema.AddressLocationHandlerFunc(
				config.ResolveAddressContractNames,
			),
			PositionInfoEnabled:        config.Mode&NeedPositionInfo != 0,
			ExtendedElaborationEnabled: config.Mode&NeedExtendedElaboration != 0,
			ImportHandler: func(
				checker *sema.Checker,
				importedLocation common.Location,
				importRange ast.Range,
			) (sema.Import, error) {

				var elaboration *sema.Elaboration
				switch importedLocation {
				case stdlib.CryptoCheckerLocation:
					cryptoChecker := stdlib.CryptoChecker()
					elaboration = cryptoChecker.Elaboration

				default:
					if seenImports[importedLocation] {
						return nil, &sema.CyclicImportsError{
							Location: importedLocation,
							Range:    importRange,
						}
					}
					seenImports[importedLocation] = true
					defer delete(seenImports, importedLocation)

					err := programs.load(config, importedLocation, location, importRange, seenImports)
					if err != nil {
						return nil, err
					}

					elaboration = programs[importedLocation].Checker.Elaboration
				}

				return sema.ElaborationImport{
					Elaboration: elaboration,
				}, nil
			},
		},
	)
	if err != nil {
		return nil, err
	}

	err = checker.Check()
	if err != nil {
		return nil, err
	}

	return checker, nil
}
