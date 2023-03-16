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

package checker

import (
	"testing"

	"github.com/onflow/cadence/runtime/sema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckBasicEntitlementDeclaration(t *testing.T) {

	t.Parallel()

	t.Run("basic", func(t *testing.T) {
		t.Parallel()
		checker, err := ParseAndCheck(t, `
			entitlement E
		`)

		assert.NoError(t, err)
		entitlement := checker.Elaboration.EntitlementType("S.test.E")
		assert.Equal(t, "E", entitlement.String())
	})

	t.Run("priv access", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			priv entitlement E 
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidAccessModifierError{}, errs[0])
	})
}

func TestCheckBasicEntitlementMappingDeclaration(t *testing.T) {

	t.Parallel()

	t.Run("basic", func(t *testing.T) {
		t.Parallel()
		checker, err := ParseAndCheck(t, `
			entitlement mapping M {}
		`)

		assert.NoError(t, err)
		entitlement := checker.Elaboration.EntitlementMapType("S.test.M")
		assert.Equal(t, "M", entitlement.String())
	})

	t.Run("with mappings", func(t *testing.T) {
		t.Parallel()
		checker, err := ParseAndCheck(t, `
			entitlement A 
			entitlement B
			entitlement C
			entitlement mapping M {
				A -> B
				B -> C
			}
		`)

		assert.NoError(t, err)
		entitlement := checker.Elaboration.EntitlementMapType("S.test.M")
		assert.Equal(t, "M", entitlement.String())
		assert.Equal(t, 2, len(entitlement.Relations))
	})

	t.Run("priv access", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			priv entitlement mapping M {}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidAccessModifierError{}, errs[0])
	})
}

func TestCheckBasicEntitlementMappingNonEntitlements(t *testing.T) {

	t.Parallel()

	t.Run("resource", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement A 
			resource B {}
			entitlement mapping M {
				A -> B
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidNonEntitlementTypeInMapError{}, errs[0])
	})

	t.Run("struct", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement A 
			struct B {}
			entitlement mapping M {
				A -> B
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidNonEntitlementTypeInMapError{}, errs[0])
	})

	t.Run("attachment", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement A 
			attachment B for AnyStruct {}
			entitlement mapping M {
				A -> B
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidNonEntitlementTypeInMapError{}, errs[0])
	})

	t.Run("interface", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement A 
			resource interface B {}
			entitlement mapping M {
				A -> B
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidNonEntitlementTypeInMapError{}, errs[0])
	})

	t.Run("contract", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement B
			contract A {}
			entitlement mapping M {
				A -> B
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidNonEntitlementTypeInMapError{}, errs[0])
	})

	t.Run("event", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement B
			event A()
			entitlement mapping M {
				A -> B
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidNonEntitlementTypeInMapError{}, errs[0])
	})

	t.Run("enum", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement B
			enum A: UInt8 {}
			entitlement mapping M {
				A -> B
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidNonEntitlementTypeInMapError{}, errs[0])
	})

	t.Run("simple type", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement B
			entitlement mapping M {
				Int -> B
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidNonEntitlementTypeInMapError{}, errs[0])
	})
}

func TestCheckEntitlementDeclarationNesting(t *testing.T) {
	t.Parallel()
	t.Run("in contract", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			contract C {
				entitlement E
			}
		`)

		assert.NoError(t, err)
	})

	t.Run("in contract interface", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			contract interface C {
				entitlement E
			}
		`)

		assert.NoError(t, err)
	})

	t.Run("in resource", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			resource R {
				entitlement E
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidNestedDeclarationError{}, errs[0])
	})

	t.Run("in resource interface", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			resource interface R {
				entitlement E
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidNestedDeclarationError{}, errs[0])
	})

	t.Run("in attachment", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			attachment A for AnyStruct {
				entitlement E
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidNestedDeclarationError{}, errs[0])
	})

	t.Run("in struct", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			struct S {
				entitlement E
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidNestedDeclarationError{}, errs[0])
	})

	t.Run("in struct", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			struct interface S {
				entitlement E
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidNestedDeclarationError{}, errs[0])
	})

	t.Run("in enum", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			enum X: UInt8 {
				entitlement E
			}
		`)

		errs := RequireCheckerErrors(t, err, 2)

		require.IsType(t, &sema.InvalidNestedDeclarationError{}, errs[0])
		require.IsType(t, &sema.InvalidNonEnumCaseError{}, errs[1])
	})
}

func TestCheckEntitlementMappingDeclarationNesting(t *testing.T) {
	t.Parallel()
	t.Run("in contract", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			contract C {
				entitlement mapping M {}
			}
		`)

		assert.NoError(t, err)
	})

	t.Run("in contract interface", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			contract interface C {
				entitlement mapping M {}
			}
		`)

		assert.NoError(t, err)
	})

	t.Run("in resource", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			resource R {
				entitlement mapping M {}
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidNestedDeclarationError{}, errs[0])
	})

	t.Run("in resource interface", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			resource interface R {
				entitlement mapping M {}
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidNestedDeclarationError{}, errs[0])
	})

	t.Run("in attachment", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			attachment A for AnyStruct {
				entitlement mapping M {}
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidNestedDeclarationError{}, errs[0])
	})

	t.Run("in struct", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			struct S {
				entitlement mapping M {}
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidNestedDeclarationError{}, errs[0])
	})

	t.Run("in struct", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			struct interface S {
				entitlement mapping M {}
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidNestedDeclarationError{}, errs[0])
	})

	t.Run("in enum", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			enum X: UInt8 {
				entitlement mapping M {}
			}
		`)

		errs := RequireCheckerErrors(t, err, 2)

		require.IsType(t, &sema.InvalidNestedDeclarationError{}, errs[0])
		require.IsType(t, &sema.InvalidNonEnumCaseError{}, errs[1])
	})
}

func TestCheckBasicEntitlementAccess(t *testing.T) {

	t.Parallel()
	t.Run("valid", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement E
			struct interface S {
				access(E) let foo: String
			}
		`)

		assert.NoError(t, err)
	})

	t.Run("multiple entitlements conjunction", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement A
			entitlement B
			entitlement C
			resource interface R {
				access(A, B) let foo: String
				access(B, C) fun bar()
			}
		`)

		assert.NoError(t, err)
	})

	t.Run("multiple entitlements disjunction", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement A
			entitlement B
			entitlement C
			resource interface R {
				access(A | B) let foo: String
				access(B | C) fun bar()
			}
		`)

		assert.NoError(t, err)
	})

	t.Run("valid in contract", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			contract C {
				entitlement E
				struct interface S {
					access(E) let foo: String
				}
			}
		`)

		assert.NoError(t, err)
	})

	t.Run("valid in contract interface", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			contract interface C {
				entitlement E
				struct interface S {
					access(E) let foo: String
				}
			}
		`)

		assert.NoError(t, err)
	})

	t.Run("qualified", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			contract C {
				entitlement E
				struct interface S {
					access(E) let foo: String
				}
			}
			resource R {
				access(C.E) fun bar() {}
			}
		`)

		assert.NoError(t, err)
	})
}

func TestCheckBasicEntitlementMappingAccess(t *testing.T) {

	t.Parallel()
	t.Run("valid", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement mapping M {}
			struct interface S {
				access(M) let foo: String
			}
		`)

		assert.NoError(t, err)
	})

	t.Run("multiple mappings conjunction", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement mapping M {} 
			entitlement mapping N {}
			resource interface R {
				access(M, N) let foo: String
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidMultipleMappedEntitlementError{}, errs[0])
	})

	t.Run("multiple mappings conjunction with regular", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement mapping M {} 
			entitlement N
			resource interface R {
				access(M, N) let foo: String
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidMultipleMappedEntitlementError{}, errs[0])
	})

	t.Run("multiple mappings disjunction", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement mapping M {} 
			entitlement mapping N {}
			resource interface R {
				access(M | N) let foo: String
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidMultipleMappedEntitlementError{}, errs[0])
	})

	t.Run("multiple mappings disjunction with regular", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement M 
			entitlement mapping N {}
			resource interface R {
				access(M | N) let foo: String
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidNonEntitlementAccessError{}, errs[0])
	})

	t.Run("valid in contract", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			contract C {
				entitlement mapping M {} 
				struct interface S {
					access(M) let foo: String
				}
			}
		`)

		assert.NoError(t, err)
	})

	t.Run("valid in contract interface", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			contract interface C {
				entitlement mapping M {} 
				struct interface S {
					access(M) let foo: String
				}
			}
		`)

		assert.NoError(t, err)
	})

	t.Run("qualified", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			contract C {
				entitlement mapping M {} 
				struct interface S {
					access(M) let foo: String
				}
			}
			resource R {
				access(C.M) fun bar() {}
			}
		`)

		assert.NoError(t, err)
	})
}

func TestCheckInvalidEntitlementAccess(t *testing.T) {

	t.Parallel()

	t.Run("invalid variable decl", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement E
			access(E) var x: String = ""
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidEntitlementAccessError{}, errs[0])
	})

	t.Run("invalid fun decl", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement E
			access(E) fun foo() {}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidEntitlementAccessError{}, errs[0])
	})

	t.Run("invalid contract field", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement E
			contract C {
				access(E) fun foo() {}
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidEntitlementAccessError{}, errs[0])
	})

	t.Run("invalid contract interface field", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement E
			contract interface C {
				access(E) fun foo()
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidEntitlementAccessError{}, errs[0])
	})

	t.Run("invalid event", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement E
			resource I {
				access(E) event Foo()
			}
		`)

		errs := RequireCheckerErrors(t, err, 2)

		require.IsType(t, &sema.InvalidNestedDeclarationError{}, errs[0])
		require.IsType(t, &sema.InvalidEntitlementAccessError{}, errs[1])
	})

	t.Run("invalid enum case", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement E
			enum X: UInt8 {
				access(E) case red
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidAccessModifierError{}, errs[0])
	})

	t.Run("missing entitlement declaration fun", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			resource R {
				access(E) fun foo() {}
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.NotDeclaredError{}, errs[0])
	})

	t.Run("missing entitlement declaration field", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			struct interface S {
				access(E) let foo: String
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.NotDeclaredError{}, errs[0])
	})
}

func TestCheckInvalidEntitlementMappingAccess(t *testing.T) {

	t.Parallel()

	t.Run("invalid variable decl", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement mapping M {}
			access(M) var x: String = ""
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidEntitlementAccessError{}, errs[0])
	})

	t.Run("invalid fun decl", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement mapping M {}
			access(M) fun foo() {}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidEntitlementAccessError{}, errs[0])
	})

	t.Run("invalid contract field", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement mapping M {}
			contract C {
				access(M) fun foo() {}
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidEntitlementAccessError{}, errs[0])
	})

	t.Run("invalid contract interface field", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement mapping M {}
			contract interface C {
				access(M) fun foo()
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidEntitlementAccessError{}, errs[0])
	})

	t.Run("invalid event", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement mapping M {}
			resource I {
				access(M) event Foo()
			}
		`)

		errs := RequireCheckerErrors(t, err, 2)

		require.IsType(t, &sema.InvalidNestedDeclarationError{}, errs[0])
		require.IsType(t, &sema.InvalidEntitlementAccessError{}, errs[1])
	})

	t.Run("invalid enum case", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement mapping M {}
			enum X: UInt8 {
				access(M) case red
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidAccessModifierError{}, errs[0])
	})

	t.Run("missing entitlement mapping declaration fun", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			resource R {
				access(M) fun foo() {}
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.NotDeclaredError{}, errs[0])
	})

	t.Run("missing entitlement mapping declaration field", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			struct interface S {
				access(M) let foo: String
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.NotDeclaredError{}, errs[0])
	})
}

func TestCheckNonEntitlementAccess(t *testing.T) {

	t.Parallel()

	t.Run("resource", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			resource E {}
			resource R {
				access(E) fun foo() {}
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidNonEntitlementAccessError{}, errs[0])
	})

	t.Run("resource interface", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			resource interface E {}
			resource R {
				access(E) fun foo() {}
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidNonEntitlementAccessError{}, errs[0])
	})

	t.Run("attachment", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			attachment E for AnyStruct {}
			resource R {
				access(E) fun foo() {}
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidNonEntitlementAccessError{}, errs[0])
	})

	t.Run("struct", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			struct E {}
			resource R {
				access(E) fun foo() {}
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidNonEntitlementAccessError{}, errs[0])
	})

	t.Run("struct interface", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			resource E {}
			resource R {
				access(E) fun foo() {}
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidNonEntitlementAccessError{}, errs[0])
	})

	t.Run("event", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			event E()
			resource R {
				access(E) fun foo() {}
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidNonEntitlementAccessError{}, errs[0])
	})

	t.Run("contract", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			contract E {}
			resource R {
				access(E) fun foo() {}
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidNonEntitlementAccessError{}, errs[0])
	})

	t.Run("contract interface", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			contract interface E {}
			resource R {
				access(E) fun foo() {}
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidNonEntitlementAccessError{}, errs[0])
	})

	t.Run("enum", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			enum E: UInt8 {}
			resource R {
				access(E) fun foo() {}
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidNonEntitlementAccessError{}, errs[0])
	})
}

func TestCheckEntitlementInheritance(t *testing.T) {

	t.Parallel()
	t.Run("valid", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement E
			struct interface I {
				access(E) fun foo() 
			}
			struct S {
				access(E) fun foo() {}
			}
		`)

		assert.NoError(t, err)
	})

	t.Run("pub subtyping invalid", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement E
			struct interface I {
				pub fun foo() 
			}
			struct S: I {
				access(E) fun foo() {}
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.ConformanceError{}, errs[0])
	})

	t.Run("pub(set) subtyping invalid", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement E
			struct interface I {
				pub(set) var x: String
			}
			struct S: I {
				access(E) var x: String
				init() {
					self.x = ""
				}
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.ConformanceError{}, errs[0])
	})

	t.Run("pub supertying invalid", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement E
			struct interface I {
				access(E) fun foo() 
			}
			struct S: I {
				pub fun foo() {}
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.ConformanceError{}, errs[0])
	})

	t.Run("pub(set) supertyping invalid", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement E
			struct interface I {
				access(E) var x: String
			}
			struct S: I {
				pub(set) var x: String
				init() {
					self.x = ""
				}
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.ConformanceError{}, errs[0])
	})

	t.Run("access contract subtyping invalid", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement E
			struct interface I {
				access(contract) fun foo() 
			}
			struct S: I {
				access(E) fun foo() {}
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.ConformanceError{}, errs[0])
	})

	t.Run("access account subtyping invalid", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement E
			struct interface I {
				access(account) fun foo() 
			}
			struct S: I {
				access(E) fun foo() {}
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.ConformanceError{}, errs[0])
	})

	t.Run("access account supertying invalid", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement E
			struct interface I {
				access(E) fun foo() 
			}
			struct S: I {
				access(account) fun foo() {}
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.ConformanceError{}, errs[0])
	})

	t.Run("access contract supertying invalid", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement E
			struct interface I {
				access(E) fun foo() 
			}
			struct S: I {
				access(contract) fun foo() {}
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.ConformanceError{}, errs[0])
	})

	t.Run("priv supertying invalid", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement E
			struct interface I {
				access(E) fun foo() 
			}
			struct S: I {
				priv fun foo() {}
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.ConformanceError{}, errs[0])
	})

	t.Run("expanded entitlements valid in disjunction", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement E
			entitlement F
			struct interface I {
				access(E) fun foo() 
			}
			struct interface J {
				access(F) fun foo() 
			}
			struct S: I, J {
				access(E | F) fun foo() {}
			}
		`)

		assert.NoError(t, err)
	})

	t.Run("reduced entitlements valid with conjunction", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement E
			entitlement F 
			struct interface I {
				access(E) fun foo() 
			}
			struct interface J {
				access(E, F) fun foo() 
			}
			struct S: I, J {
				access(E) fun foo() {}
			}
		`)

		assert.NoError(t, err)
	})

	t.Run("expanded entitlements invalid in conjunction", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement E
			entitlement F
			struct interface I {
				access(E) fun foo() 
			}
			struct interface J {
				access(F) fun foo() 
			}
			struct S: I, J {
				access(E, F) fun foo() {}
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.ConformanceError{}, errs[0])
	})

	t.Run("expanded entitlements invalid", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement E
			entitlement F
			struct interface I {
				access(E) fun foo() 
			}
			struct interface J {
				access(F) fun foo() 
			}
			struct S: I, J {
				access(E) fun foo() {}
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.ConformanceError{}, errs[0])
	})

	t.Run("reduced entitlements invalid with disjunction", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement E
			entitlement F 
			struct interface I {
				access(E) fun foo() 
			}
			struct interface J {
				access(E | F) fun foo() 
			}
			struct S: I, J {
				access(E) fun foo() {}
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.ConformanceError{}, errs[0])
	})

	t.Run("different entitlements invalid", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement E
			entitlement F
			entitlement G 
			struct interface I {
				access(E) fun foo() 
			}
			struct interface J {
				access(F) fun foo() 
			}
			struct S: I, J {
				access(E | G) fun foo() {}
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.ConformanceError{}, errs[0])
	})
}

func TestCheckEntitlementTypeAnnotation(t *testing.T) {

	t.Parallel()

	t.Run("invalid local annot", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement E
			let x: E = ""
		`)

		errs := RequireCheckerErrors(t, err, 2)

		require.IsType(t, &sema.DirectEntitlementAnnotationError{}, errs[0])
		require.IsType(t, &sema.TypeMismatchError{}, errs[1])
	})

	t.Run("invalid param annot", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement E
			pub fun foo(e: E) {}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.DirectEntitlementAnnotationError{}, errs[0])
	})

	t.Run("invalid return annot", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement E
			resource interface I {
				pub fun foo(): E 
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.DirectEntitlementAnnotationError{}, errs[0])
	})

	t.Run("invalid field annot", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement E
			resource interface I {
				let e: E
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.DirectEntitlementAnnotationError{}, errs[0])
	})

	t.Run("invalid conformance annotation", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement E
			resource R: E {}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidConformanceError{}, errs[0])
	})

	t.Run("invalid array annot", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement E
			resource interface I {
				let e: [E]
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.DirectEntitlementAnnotationError{}, errs[0])
	})

	t.Run("invalid fun annot", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement E
			resource interface I {
				let e: (fun (E): Void)
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.DirectEntitlementAnnotationError{}, errs[0])
	})

	t.Run("invalid enum conformance", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement E
			enum X: E {}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidEnumRawTypeError{}, errs[0])
	})

	t.Run("invalid dict annot", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement E
			resource interface I {
				let e: {E: E}
			}
		`)

		errs := RequireCheckerErrors(t, err, 2)

		// key
		require.IsType(t, &sema.InvalidDictionaryKeyTypeError{}, errs[0])
		// value
		require.IsType(t, &sema.DirectEntitlementAnnotationError{}, errs[1])
	})

	t.Run("invalid fun annot", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement E
			resource interface I {
				let e: (fun (E): Void)
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.DirectEntitlementAnnotationError{}, errs[0])
	})

	t.Run("runtype type", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement E
			let e = Type<E>()
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.DirectEntitlementAnnotationError{}, errs[0])
	})

	t.Run("type arg", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheckAccount(t, `
			entitlement E
			let e = authAccount.load<E>(from: /storage/foo)
		`)

		errs := RequireCheckerErrors(t, err, 2)

		require.IsType(t, &sema.DirectEntitlementAnnotationError{}, errs[0])
		// entitlements are not storable either
		require.IsType(t, &sema.TypeMismatchError{}, errs[1])
	})

	t.Run("restricted", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheckAccount(t, `
			entitlement E
			resource interface I {
				let e: E{E}
			}
		`)

		errs := RequireCheckerErrors(t, err, 2)

		require.IsType(t, &sema.InvalidRestrictionTypeError{}, errs[0])
		require.IsType(t, &sema.InvalidRestrictedTypeError{}, errs[1])
	})

	t.Run("reference", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheckAccount(t, `
			entitlement E
			resource interface I {
				let e: &E
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.DirectEntitlementAnnotationError{}, errs[0])
	})

	t.Run("capability", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheckAccount(t, `
			entitlement E
			resource interface I {
				let e: Capability<&E>
			}
		`)

		errs := RequireCheckerErrors(t, err, 2)

		require.IsType(t, &sema.DirectEntitlementAnnotationError{}, errs[0])
		require.IsType(t, &sema.DirectEntitlementAnnotationError{}, errs[1])
	})

	t.Run("optional", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheckAccount(t, `
			entitlement E
			resource interface I {
				let e: E?
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.DirectEntitlementAnnotationError{}, errs[0])
	})
}

func TestCheckEntitlementMappingTypeAnnotation(t *testing.T) {

	t.Parallel()

	t.Run("invalid local annot", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement mapping E {}
			let x: E = ""
		`)

		errs := RequireCheckerErrors(t, err, 2)

		require.IsType(t, &sema.DirectEntitlementAnnotationError{}, errs[0])
		require.IsType(t, &sema.TypeMismatchError{}, errs[1])
	})

	t.Run("invalid param annot", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement mapping E {}
			pub fun foo(e: E) {}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.DirectEntitlementAnnotationError{}, errs[0])
	})

	t.Run("invalid return annot", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement mapping E {}
			resource interface I {
				pub fun foo(): E 
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.DirectEntitlementAnnotationError{}, errs[0])
	})

	t.Run("invalid field annot", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement mapping E {}
			resource interface I {
				let e: E
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.DirectEntitlementAnnotationError{}, errs[0])
	})

	t.Run("invalid conformance annotation", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement mapping E {}
			resource R: E {}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidConformanceError{}, errs[0])
	})

	t.Run("invalid array annot", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement mapping E {}
			resource interface I {
				let e: [E]
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.DirectEntitlementAnnotationError{}, errs[0])
	})

	t.Run("invalid fun annot", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement mapping E {}
			resource interface I {
				let e: (fun (E): Void)
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.DirectEntitlementAnnotationError{}, errs[0])
	})

	t.Run("invalid enum conformance", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement mapping E {}
			enum X: E {}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidEnumRawTypeError{}, errs[0])
	})

	t.Run("invalid dict annot", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement mapping E {}
			resource interface I {
				let e: {E: E}
			}
		`)

		errs := RequireCheckerErrors(t, err, 2)

		// key
		require.IsType(t, &sema.InvalidDictionaryKeyTypeError{}, errs[0])
		// value
		require.IsType(t, &sema.DirectEntitlementAnnotationError{}, errs[1])
	})

	t.Run("invalid fun annot", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement mapping E {}
			resource interface I {
				let e: (fun (E): Void)
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.DirectEntitlementAnnotationError{}, errs[0])
	})

	t.Run("runtype type", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement mapping E {}
			let e = Type<E>()
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.DirectEntitlementAnnotationError{}, errs[0])
	})

	t.Run("type arg", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheckAccount(t, `
			entitlement mapping E {}
			let e = authAccount.load<E>(from: /storage/foo)
		`)

		errs := RequireCheckerErrors(t, err, 2)

		require.IsType(t, &sema.DirectEntitlementAnnotationError{}, errs[0])
		// entitlements are not storable either
		require.IsType(t, &sema.TypeMismatchError{}, errs[1])
	})

	t.Run("restricted", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheckAccount(t, `
			entitlement mapping E {}
			resource interface I {
				let e: E{E}
			}
		`)

		errs := RequireCheckerErrors(t, err, 2)

		require.IsType(t, &sema.InvalidRestrictionTypeError{}, errs[0])
		require.IsType(t, &sema.InvalidRestrictedTypeError{}, errs[1])
	})

	t.Run("reference", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheckAccount(t, `
			entitlement mapping E {}
			resource interface I {
				let e: &E
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.DirectEntitlementAnnotationError{}, errs[0])
	})

	t.Run("capability", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheckAccount(t, `
			entitlement mapping E {}
			resource interface I {
				let e: Capability<&E>
			}
		`)

		errs := RequireCheckerErrors(t, err, 2)

		require.IsType(t, &sema.DirectEntitlementAnnotationError{}, errs[0])
		require.IsType(t, &sema.DirectEntitlementAnnotationError{}, errs[1])
	})

	t.Run("optional", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheckAccount(t, `
			entitlement mapping E {}
			resource interface I {
				let e: E?
			}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.DirectEntitlementAnnotationError{}, errs[0])
	})
}

func TestChecAttachmentEntitlementAccessAnnotation(t *testing.T) {

	t.Parallel()
	t.Run("mapping allowed", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement mapping E {}
			access(E) attachment A for AnyStruct {}
		`)

		assert.NoError(t, err)
	})

	t.Run("entitlement set not allowed", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
			entitlement E
			entitlement F
			access(E, F) attachment A for AnyStruct {}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidEntitlementAccessError{}, errs[0])
	})

	t.Run("mapping allowed in contract", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
		contract C {
			entitlement mapping E {}
			access(E) attachment A for AnyStruct {}
		}
		`)

		assert.NoError(t, err)
	})

	t.Run("entitlement set not allowed in contract", func(t *testing.T) {
		t.Parallel()
		_, err := ParseAndCheck(t, `
		contract C {
			entitlement E
			access(E) attachment A for AnyStruct {}
		}
		`)

		errs := RequireCheckerErrors(t, err, 1)

		require.IsType(t, &sema.InvalidEntitlementAccessError{}, errs[0])
	})

}
