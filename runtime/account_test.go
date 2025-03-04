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

package runtime_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/onflow/atree"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/onflow/cadence"
	"github.com/onflow/cadence/encoding/json"
	. "github.com/onflow/cadence/runtime"
	"github.com/onflow/cadence/runtime/common"
	"github.com/onflow/cadence/runtime/errors"
	"github.com/onflow/cadence/runtime/interpreter"
	"github.com/onflow/cadence/runtime/sema"
	"github.com/onflow/cadence/runtime/stdlib"
	"github.com/onflow/cadence/runtime/tests/checker"
	. "github.com/onflow/cadence/runtime/tests/runtime_utils"
	"github.com/onflow/cadence/runtime/tests/utils"
	. "github.com/onflow/cadence/runtime/tests/utils"
)

func TestRuntimeAccountKeyConstructor(t *testing.T) {

	t.Parallel()

	rt := NewTestInterpreterRuntime()

	script := []byte(`
        access(all) fun main(): AccountKey {
            let key = AccountKey(
                PublicKey(
                    publicKey: "0102".decodeHex(),
                    signatureAlgorithm: SignatureAlgorithm.ECDSA_P256
                ),
                hashAlgorithm: HashAlgorithm.SHA3_256,
                weight: 1.7
            )

            return key
          }
    `)

	runtimeInterface := &TestRuntimeInterface{}

	_, err := rt.ExecuteScript(
		Script{
			Source: script,
		},
		Context{
			Interface: runtimeInterface,
			Location:  common.ScriptLocation{},
		},
	)
	RequireError(t, err)

	assert.Contains(t, err.Error(), "cannot find variable in this scope: `AccountKey`")
}

func noopRuntimeUInt64Getter(_ common.Address) (uint64, error) {
	return 0, nil
}

func TestRuntimeReturnPublicAccount(t *testing.T) {

	t.Parallel()

	rt := NewTestInterpreterRuntime()

	script := []byte(`
        access(all) fun main(): &Account {
            return getAccount(0x02)
          }
    `)

	runtimeInterface := &TestRuntimeInterface{
		OnGetAccountBalance:          noopRuntimeUInt64Getter,
		OnGetAccountAvailableBalance: noopRuntimeUInt64Getter,
		OnGetStorageUsed:             noopRuntimeUInt64Getter,
		OnGetStorageCapacity:         noopRuntimeUInt64Getter,
		OnAccountKeysCount:           noopRuntimeUInt64Getter,
		Storage:                      NewTestLedger(nil, nil),
	}

	_, err := rt.ExecuteScript(
		Script{
			Source: script,
		},
		Context{
			Interface: runtimeInterface,
			Location:  common.ScriptLocation{},
		},
	)
	require.NoError(t, err)
}

func TestRuntimeReturnAuthAccount(t *testing.T) {

	t.Parallel()

	rt := NewTestInterpreterRuntime()

	script := []byte(`
        access(all) fun main(): auth(Storage) &Account {
            return getAuthAccount<auth(Storage) &Account>(0x02)
          }
    `)

	runtimeInterface := &TestRuntimeInterface{
		OnGetAccountBalance:          noopRuntimeUInt64Getter,
		OnGetAccountAvailableBalance: noopRuntimeUInt64Getter,
		OnGetStorageUsed:             noopRuntimeUInt64Getter,
		OnGetStorageCapacity:         noopRuntimeUInt64Getter,
		OnAccountKeysCount:           noopRuntimeUInt64Getter,
		Storage:                      NewTestLedger(nil, nil),
	}

	_, err := rt.ExecuteScript(
		Script{
			Source: script,
		},
		Context{
			Interface: runtimeInterface,
			Location:  common.ScriptLocation{},
		},
	)
	require.NoError(t, err)
}

func TestRuntimeStoreAccountAPITypes(t *testing.T) {

	t.Parallel()

	nextTransactionLocation := NewTransactionLocationGenerator()

	for _, ty := range []sema.Type{
		sema.AccountKeyType,
		sema.PublicKeyType,
	} {

		rt := NewTestInterpreterRuntime()

		script := []byte(fmt.Sprintf(`
            transaction {

                prepare(signer: auth(SaveValue) &Account) {
                    signer.storage.save<%s>(panic(""))
                }
            }
        `, ty.String()))

		runtimeInterface := &TestRuntimeInterface{}

		err := rt.ExecuteTransaction(
			Script{
				Source: script,
			},
			Context{
				Interface: runtimeInterface,
				Location:  nextTransactionLocation(),
			},
		)

		RequireError(t, err)

		assert.Contains(t, err.Error(), "expected `Storable`")
	}
}

var accountKeyA = &stdlib.AccountKey{
	KeyIndex: 0,
	PublicKey: &stdlib.PublicKey{
		PublicKey: []byte{1, 2, 3},
		SignAlgo:  sema.SignatureAlgorithmECDSA_P256,
	},
	HashAlgo:  sema.HashAlgorithmSHA3_256,
	Weight:    100,
	IsRevoked: false,
}

var accountKeyB = &stdlib.AccountKey{
	KeyIndex: 1,
	PublicKey: &stdlib.PublicKey{
		PublicKey: []byte{4, 5, 6},
		SignAlgo:  sema.SignatureAlgorithmECDSA_secp256k1,
	},
	HashAlgo:  sema.HashAlgorithmSHA3_256,
	Weight:    100,
	IsRevoked: false,
}

var revokedAccountKeyA = func() *stdlib.AccountKey {
	revokedKey := *accountKeyA
	revokedKey.IsRevoked = true
	return &revokedKey
}()

type accountTestEnvironment struct {
	storage          *testAccountKeyStorage
	runtime          Runtime
	runtimeInterface *TestRuntimeInterface
}

func newAccountTestEnv() accountTestEnvironment {
	storage := newTestAccountKeyStorage()
	rt := NewTestInterpreterRuntime()
	rtInterface := getAccountKeyTestRuntimeInterface(storage)

	addPublicKeyValidation(rtInterface, nil)

	return accountTestEnvironment{
		storage,
		rt,
		rtInterface,
	}
}

func TestRuntimeAuthAccountKeys(t *testing.T) {

	t.Parallel()

	initTestEnvironment := func(t *testing.T, location Location) accountTestEnvironment {
		testEnv := newAccountTestEnv()
		addAuthAccountKey(t, testEnv.runtime, testEnv.runtimeInterface, location)
		return testEnv
	}

	t.Run("add key", func(t *testing.T) {

		t.Parallel()

		nextTransactionLocation := NewTransactionLocationGenerator()
		testEnv := initTestEnvironment(t, nextTransactionLocation())

		assert.Equal(t, []*stdlib.AccountKey{accountKeyA}, testEnv.storage.keys)
		assert.Equal(t, accountKeyA, testEnv.storage.returnedKey)
	})

	t.Run("get existing key", func(t *testing.T) {

		t.Parallel()

		nextTransactionLocation := NewTransactionLocationGenerator()
		testEnv := initTestEnvironment(t, nextTransactionLocation())

		test := accountKeyTestCase{
			code: `
                transaction {
                    prepare(signer: &Account) {
                        let key = signer.keys.get(keyIndex: 0) ?? panic("unexpectedly nil")
                        log(key)
                        assert(!key.isRevoked)
                    }
                }`,
			args: []cadence.Value{},
		}

		err := test.executeTransaction(
			testEnv.runtime,
			testEnv.runtimeInterface,
			nextTransactionLocation(),
		)
		require.NoError(t, err)

		assert.Equal(t, []*stdlib.AccountKey{accountKeyA}, testEnv.storage.keys)
		assert.Equal(t, accountKeyA, testEnv.storage.returnedKey)
		assert.Equal(
			t,
			[]string{
				"AccountKey(keyIndex: 0, publicKey: PublicKey(publicKey: [1, 2, 3], signatureAlgorithm: SignatureAlgorithm(rawValue: 1)), hashAlgorithm: HashAlgorithm(rawValue: 3), weight: 100.00000000, isRevoked: false)",
			},
			testEnv.storage.logs,
		)
	})

	t.Run("get non-existing key", func(t *testing.T) {

		t.Parallel()

		nextTransactionLocation := NewTransactionLocationGenerator()
		testEnv := initTestEnvironment(t, nextTransactionLocation())

		test := accountKeyTestCase{
			code: `
                transaction {
                    prepare(signer: &Account) {
                        let key: AccountKey? = signer.keys.get(keyIndex: 5)
                        assert(key == nil)
                    }
                }`,
			args: []cadence.Value{},
		}

		err := test.executeTransaction(
			testEnv.runtime,
			testEnv.runtimeInterface,
			nextTransactionLocation(),
		)
		require.NoError(t, err)
		assert.Nil(t, testEnv.storage.returnedKey)
	})

	t.Run("revoke existing key", func(t *testing.T) {

		t.Parallel()

		nextTransactionLocation := NewTransactionLocationGenerator()
		testEnv := initTestEnvironment(t, nextTransactionLocation())

		test := accountKeyTestCase{
			code: `
                transaction {
                    prepare(signer: auth(RevokeKey) &Account) {
                        let key = signer.keys.revoke(keyIndex: 0) ?? panic("unexpectedly nil")
                        assert(key.isRevoked)
                    }
                }`,
			args: []cadence.Value{},
		}

		err := test.executeTransaction(
			testEnv.runtime,
			testEnv.runtimeInterface,
			nextTransactionLocation(),
		)
		require.NoError(t, err)

		assert.Equal(t, []*stdlib.AccountKey{revokedAccountKeyA}, testEnv.storage.keys)
		assert.Equal(t, revokedAccountKeyA, testEnv.storage.returnedKey)
	})

	t.Run("revoke non-existing key", func(t *testing.T) {

		t.Parallel()

		nextTransactionLocation := NewTransactionLocationGenerator()
		testEnv := initTestEnvironment(t, nextTransactionLocation())

		test := accountKeyTestCase{
			code: `
                transaction {
                    prepare(signer: auth(RevokeKey) &Account) {
                        let key: AccountKey? = signer.keys.revoke(keyIndex: 5)
                        assert(key == nil)
                    }
                }`,
			args: []cadence.Value{},
		}

		err := test.executeTransaction(
			testEnv.runtime,
			testEnv.runtimeInterface,
			nextTransactionLocation(),
		)
		require.NoError(t, err)
		assert.Nil(t, testEnv.storage.returnedKey)
	})

	t.Run("get key count after revocation", func(t *testing.T) {
		t.Parallel()

		nextTransactionLocation := NewTransactionLocationGenerator()
		testEnv := initTestEnvironment(t, nextTransactionLocation())

		test := accountKeyTestCase{
			code: `
                transaction {
                    prepare(signer: auth(RevokeKey) &Account) {
                        assert(signer.keys.count == 1)

                        let key = signer.keys.revoke(keyIndex: 0) ?? panic("unexpectedly nil")
                        assert(key.isRevoked)

                        assert(signer.keys.count == 0)
                    }
                }
            `,
			args: []cadence.Value{},
		}

		err := test.executeTransaction(
			testEnv.runtime,
			testEnv.runtimeInterface,
			nextTransactionLocation(),
		)
		require.NoError(t, err)
		assert.Equal(t, []*stdlib.AccountKey{revokedAccountKeyA}, testEnv.storage.keys)
		assert.Equal(t, revokedAccountKeyA, testEnv.storage.returnedKey)
	})

	t.Run("test keys forEach, after add and revoke", func(t *testing.T) {
		t.Parallel()

		nextTransactionLocation := NewTransactionLocationGenerator()
		testEnv := initTestEnvironment(t, nextTransactionLocation())

		test := accountKeyTestCase{
			code: `
                transaction {
                    prepare(signer: auth(Keys) &Account) {
                        signer.keys.add(
                            publicKey: PublicKey(
                                publicKey: [1, 2, 3],
                                signatureAlgorithm: SignatureAlgorithm.ECDSA_P256
                            ),
                            hashAlgorithm: HashAlgorithm.SHA3_256,
                            weight: 100.0
                        )

                        signer.keys.revoke(keyIndex: 0) ?? panic("unexpectedly nil")

                        signer.keys.forEach(fun(key: AccountKey): Bool {
                            log(key.keyIndex)
                            return true
                        })
                    }
                }
            `,
			args: []cadence.Value{},
		}

		err := test.executeTransaction(
			testEnv.runtime,
			testEnv.runtimeInterface,
			nextTransactionLocation(),
		)
		require.NoError(t, err)

		keys := make(map[int]*AccountKey, len(testEnv.storage.keys))
		for _, key := range testEnv.storage.keys {
			keys[key.KeyIndex] = key
		}
		for _, loggedIndex := range testEnv.storage.logs {
			keyIdx, err := strconv.Atoi(loggedIndex)
			require.NoError(t, err)

			key, ok := keys[keyIdx]

			require.NotNil(t, key)

			assert.True(t, ok) // no key should be passed to the callback twice
			keys[keyIdx] = nil
		}
	})
}

func TestRuntimeAuthAccountKeysAdd(t *testing.T) {

	t.Parallel()

	rt := NewTestInterpreterRuntime()

	pubKey := newBytesValue([]byte{1, 2, 3})

	const code = `
       transaction(publicKey: [UInt8]) {
           prepare(signer: auth(BorrowValue) &Account) {
               let acct = Account(payer: signer)
               acct.keys.add(
                   publicKey: PublicKey(
                       publicKey: publicKey,
                       signatureAlgorithm: SignatureAlgorithm.ECDSA_P256
                   ),
                   hashAlgorithm: HashAlgorithm.SHA3_256,
                   weight: 100.0
               )
           }
       }
   `

	storage := newTestAccountKeyStorage()
	runtimeInterface := getAccountKeyTestRuntimeInterface(storage)
	addPublicKeyValidation(runtimeInterface, nil)

	nextTransactionLocation := NewTransactionLocationGenerator()

	err := rt.ExecuteTransaction(
		Script{
			Source:    []byte(code),
			Arguments: encodeArgs([]cadence.Value{pubKey}),
		},
		Context{
			Location:  nextTransactionLocation(),
			Interface: runtimeInterface,
		},
	)

	require.NoError(t, err)
	assert.Len(t, storage.keys, 1)

	require.Len(t, storage.events, 2)

	assert.EqualValues(t,
		stdlib.AccountCreatedEventType.ID(),
		storage.events[0].Type().ID(),
	)

	assert.EqualValues(t,
		stdlib.AccountKeyAddedFromPublicKeyEventType.ID(),
		storage.events[1].Type().ID(),
	)
}

func TestRuntimePublicAccountKeys(t *testing.T) {

	t.Parallel()

	initTestEnv := func(keys ...*AccountKey) accountTestEnvironment {
		testEnv := newAccountTestEnv()
		testEnv.storage.keys = append(testEnv.storage.keys, keys...)
		for _, key := range keys {
			if !key.IsRevoked {
				testEnv.storage.unrevokedKeyCount++
			}
		}
		return testEnv
	}

	t.Run("get key", func(t *testing.T) {

		t.Parallel()

		testEnv := initTestEnv(accountKeyA, accountKeyB)
		test := accountKeyTestCase{
			code: `
              access(all) fun main(): AccountKey? {
                  let acc = getAccount(0x02)
                  return acc.keys.get(keyIndex: 0)
              }
            `,
			args: []cadence.Value{},
		}

		value, err := test.executeScript(testEnv.runtime, testEnv.runtimeInterface)
		require.NoError(t, err)
		require.NotNil(t, value)

		require.IsType(t, cadence.Optional{}, value)
		optionalValue := value.(cadence.Optional)

		expectedValue := accountKeyExportedValue(
			0,
			[]byte{1, 2, 3},
			sema.SignatureAlgorithmECDSA_P256,
			sema.HashAlgorithmSHA3_256,
			"100.0",
			false,
		)

		assert.Equal(t, expectedValue, optionalValue.Value)
		assert.Equal(t, accountKeyA, testEnv.storage.returnedKey)

	})

	t.Run("get another key", func(t *testing.T) {

		t.Parallel()

		testEnv := initTestEnv(accountKeyA, accountKeyB)

		test := accountKeyTestCase{
			code: `
              access(all) fun main(): AccountKey? {
                  let acc = getAccount(0x02)
                  return acc.keys.get(keyIndex: 1)
              }
            `,
			args: []cadence.Value{},
		}

		value, err := test.executeScript(testEnv.runtime, testEnv.runtimeInterface)
		require.NoError(t, err)
		require.NotNil(t, value)

		require.IsType(t, cadence.Optional{}, value)
		optionalValue := value.(cadence.Optional)

		expectedValue := accountKeyExportedValue(
			1,
			[]byte{4, 5, 6},
			sema.SignatureAlgorithmECDSA_secp256k1,
			sema.HashAlgorithmSHA3_256,
			"100.0",
			false,
		)

		assert.Equal(t, expectedValue, optionalValue.Value)
		assert.Equal(t, accountKeyB, testEnv.storage.returnedKey)
	})

	t.Run("get non-existing key", func(t *testing.T) {

		t.Parallel()

		testEnv := initTestEnv(accountKeyA, accountKeyB)

		test := accountKeyTestCase{
			code: `
                access(all) fun main(): AccountKey? {
                    let acc = getAccount(0x02)
                    return acc.keys.get(keyIndex: 4)
                }
            `,
			args: []cadence.Value{},
		}

		value, err := test.executeScript(testEnv.runtime, testEnv.runtimeInterface)
		require.NoError(t, err)
		require.NotNil(t, value)

		require.IsType(t, cadence.Optional{}, value)
		optionalValue := value.(cadence.Optional)

		assert.Nil(t, optionalValue.Value)
	})

	t.Run("get revoked key", func(t *testing.T) {

		t.Parallel()

		testEnv := initTestEnv(revokedAccountKeyA, accountKeyB)

		test := accountKeyTestCase{
			code: `
              access(all) fun main(): AccountKey? {
                  let acc = getAccount(0x02)
                  var keys: &Account.Keys = acc.keys
                  return keys.get(keyIndex: 0)
              }
            `,
			args: []cadence.Value{},
		}

		value, err := test.executeScript(testEnv.runtime, testEnv.runtimeInterface)
		require.NoError(t, err)
		require.NotNil(t, value)

		require.IsType(t, cadence.Optional{}, value)
		optionalValue := value.(cadence.Optional)

		expectedValue := accountKeyExportedValue(
			0,
			[]byte{1, 2, 3},
			sema.SignatureAlgorithmECDSA_P256,
			sema.HashAlgorithmSHA3_256,
			"100.0",
			true,
		)

		assert.Equal(t, expectedValue, optionalValue.Value)
		assert.Equal(t, revokedAccountKeyA, testEnv.storage.returnedKey)
	})

	t.Run("get key count", func(t *testing.T) {
		t.Parallel()

		testEnv := initTestEnv(revokedAccountKeyA, accountKeyB)

		test := accountKeyTestCase{
			code: `
            access(all) fun main(): UInt64 {
                return getAccount(0x02).keys.count
            }
            `,
		}

		value, err := test.executeScript(testEnv.runtime, testEnv.runtimeInterface)
		require.NoError(t, err)
		require.NotNil(t, value)

		expected := cadence.UInt64(1)
		assert.Equal(t, expected, value)
	})

	t.Run("test keys.forEach", func(t *testing.T) {
		t.Parallel()

		testEnv := initTestEnv(revokedAccountKeyA, accountKeyB)
		test := accountKeyTestCase{
			code: `
                access(all) fun main() {
                        getAccount(0x02).keys.forEach(fun(key: AccountKey): Bool {
                            log(key.keyIndex)
                            return true
                        })
                    }
            `,
			args: []cadence.Value{},
		}

		value, err := test.executeScript(testEnv.runtime, testEnv.runtimeInterface)
		require.NoError(t, err)
		utils.AssertEqualWithDiff(t, cadence.Void{}, value)

		keys := make(map[int]*AccountKey, len(testEnv.storage.keys))
		for _, key := range testEnv.storage.keys {
			keys[key.KeyIndex] = key
		}
		for _, loggedIndex := range testEnv.storage.logs {
			keyIdx, err := strconv.Atoi(loggedIndex)
			require.NoError(t, err)

			key, ok := keys[keyIdx]

			assert.True(t, ok)

			require.NotNil(t, key)
			keys[keyIdx] = nil // no key should be passed to the callback twice
		}
	})
}

func TestRuntimeHashAlgorithm(t *testing.T) {

	t.Parallel()

	rt := NewTestInterpreterRuntime()

	script := []byte(`
        access(all) fun main(): [HashAlgorithm?] {
            var key1: HashAlgorithm? = HashAlgorithm.SHA3_256

            var key2: HashAlgorithm? = HashAlgorithm(rawValue: 3)

            var key3: HashAlgorithm? = HashAlgorithm(rawValue: 100)
            return [key1, key2, key3]
          }
    `)

	storage := NewTestLedger(nil, nil)

	runtimeInterface := &TestRuntimeInterface{
		Storage: storage,
	}

	result, err := rt.ExecuteScript(
		Script{
			Source: script,
		},
		Context{
			Interface: runtimeInterface,
			Location:  common.ScriptLocation{},
		},
	)
	require.NoError(t, err)

	require.IsType(t, cadence.Array{}, result)
	array := result.(cadence.Array)

	require.Len(t, array.Values, 3)

	// Check key1
	require.IsType(t, cadence.Optional{}, array.Values[0])
	optionalValue := array.Values[0].(cadence.Optional)

	require.IsType(t, cadence.Enum{}, optionalValue.Value)
	builtinStruct := optionalValue.Value.(cadence.Enum)

	require.Len(t, builtinStruct.Fields, 1)
	assert.Equal(t,
		cadence.NewUInt8(HashAlgorithmSHA3_256.RawValue()),
		builtinStruct.Fields[0],
	)

	// Check key2
	require.IsType(t, cadence.Optional{}, array.Values[1])
	optionalValue = array.Values[1].(cadence.Optional)

	require.IsType(t, cadence.Enum{}, optionalValue.Value)
	builtinStruct = optionalValue.Value.(cadence.Enum)

	require.Len(t, builtinStruct.Fields, 1)
	assert.Equal(t,
		cadence.NewUInt8(HashAlgorithmSHA3_256.RawValue()),
		builtinStruct.Fields[0],
	)

	// Check key3
	require.IsType(t, cadence.Optional{}, array.Values[2])
	optionalValue = array.Values[2].(cadence.Optional)

	require.Nil(t, optionalValue.Value)
}

func TestRuntimeSignatureAlgorithm(t *testing.T) {

	t.Parallel()

	rt := NewTestInterpreterRuntime()

	script := []byte(`
        access(all) fun main(): [SignatureAlgorithm?] {
            var key1: SignatureAlgorithm? = SignatureAlgorithm.ECDSA_secp256k1

            var key2: SignatureAlgorithm? = SignatureAlgorithm(rawValue: 2)

            var key3: SignatureAlgorithm? = SignatureAlgorithm(rawValue: 100)
            return [key1, key2, key3]
        }
    `)

	storage := NewTestLedger(nil, nil)

	runtimeInterface := &TestRuntimeInterface{
		Storage: storage,
	}

	result, err := rt.ExecuteScript(
		Script{
			Source: script,
		},
		Context{
			Interface: runtimeInterface,
			Location:  common.ScriptLocation{},
		},
	)
	require.NoError(t, err)

	require.IsType(t, cadence.Array{}, result)
	array := result.(cadence.Array)

	require.Len(t, array.Values, 3)

	// Check key1
	require.IsType(t, cadence.Optional{}, array.Values[0])
	optionalValue := array.Values[0].(cadence.Optional)

	require.IsType(t, cadence.Enum{}, optionalValue.Value)
	builtinStruct := optionalValue.Value.(cadence.Enum)

	require.Len(t, builtinStruct.Fields, 1)
	assert.Equal(t,
		cadence.NewUInt8(SignatureAlgorithmECDSA_secp256k1.RawValue()),
		builtinStruct.Fields[0],
	)

	// Check key2
	require.IsType(t, cadence.Optional{}, array.Values[1])
	optionalValue = array.Values[1].(cadence.Optional)

	require.IsType(t, cadence.Enum{}, optionalValue.Value)
	builtinStruct = optionalValue.Value.(cadence.Enum)

	require.Len(t, builtinStruct.Fields, 1)
	assert.Equal(t,
		cadence.NewUInt8(SignatureAlgorithmECDSA_secp256k1.RawValue()),
		builtinStruct.Fields[0],
	)

	// Check key3
	require.IsType(t, cadence.Optional{}, array.Values[2])
	optionalValue = array.Values[2].(cadence.Optional)

	require.Nil(t, optionalValue.Value)
}

// Utility methods and types

var AccountKeyType = ExportedBuiltinType(sema.AccountKeyType).(*cadence.StructType)
var PublicKeyType = ExportedBuiltinType(sema.PublicKeyType).(*cadence.StructType)
var SignAlgoType = ExportedBuiltinType(sema.SignatureAlgorithmType).(*cadence.EnumType)
var HashAlgoType = ExportedBuiltinType(sema.HashAlgorithmType).(*cadence.EnumType)

func ExportedBuiltinType(internalType sema.Type) cadence.Type {
	return ExportType(internalType, map[sema.TypeID]cadence.Type{})
}

func newBytesValue(bytes []byte) cadence.Array {
	result := make([]cadence.Value, len(bytes))
	for index, value := range bytes {
		result[index] = cadence.NewUInt8(value)
	}
	return cadence.NewArray(result).
		WithType(&cadence.VariableSizedArrayType{
			ElementType: cadence.UInt8Type,
		})
}

func newSignAlgoValue(signAlgo sema.SignatureAlgorithm) cadence.Enum {
	return cadence.NewEnum([]cadence.Value{
		cadence.NewUInt8(signAlgo.RawValue()),
	}).WithType(SignAlgoType)
}

func accountKeyExportedValue(
	index int,
	publicKeyBytes []byte,
	signAlgo sema.SignatureAlgorithm,
	hashAlgo sema.HashAlgorithm,
	weight string,
	isRevoked bool,
) cadence.Struct {

	weightUFix64, err := cadence.NewUFix64(weight)
	if err != nil {
		panic(err)
	}

	return cadence.Struct{
		StructType: AccountKeyType,
		Fields: []cadence.Value{
			// Key index
			cadence.NewInt(index),

			// Public Key (struct)
			cadence.Struct{
				StructType: PublicKeyType,
				Fields: []cadence.Value{
					// Public key (bytes)
					newBytesValue(publicKeyBytes),

					// Signature Algo
					newSignAlgoValue(signAlgo),
				},
			},

			// Hash algo
			cadence.NewEnum([]cadence.Value{
				cadence.NewUInt8(hashAlgo.RawValue()),
			}).WithType(HashAlgoType),

			// Weight
			weightUFix64,

			// IsRevoked
			cadence.NewBool(isRevoked),
		},
	}
}

func getAccountKeyTestRuntimeInterface(storage *testAccountKeyStorage) *TestRuntimeInterface {
	return &TestRuntimeInterface{
		Storage: NewTestLedger(nil, nil),
		OnGetSigningAccounts: func() ([]Address, error) {
			return []Address{{42}}, nil
		},
		OnCreateAccount: func(payer Address) (address Address, err error) {
			return Address{42}, nil
		},
		OnAddAccountKey: func(address Address, publicKey *stdlib.PublicKey, hashAlgo HashAlgorithm, weight int) (*stdlib.AccountKey, error) {
			index := len(storage.keys)
			accountKey := &stdlib.AccountKey{
				KeyIndex:  index,
				PublicKey: publicKey,
				HashAlgo:  hashAlgo,
				Weight:    weight,
				IsRevoked: false,
			}

			storage.keys = append(storage.keys, accountKey)
			storage.unrevokedKeyCount += 1
			storage.returnedKey = accountKey
			return accountKey, nil
		},
		OnGetAccountKey: func(address Address, index int) (*stdlib.AccountKey, error) {
			if index >= len(storage.keys) {
				storage.returnedKey = nil
				return nil, nil
			}

			accountKey := storage.keys[index]
			storage.returnedKey = accountKey
			return accountKey, nil
		},
		OnRemoveAccountKey: func(address Address, index int) (*stdlib.AccountKey, error) {
			if index >= len(storage.keys) {
				storage.returnedKey = nil
				return nil, nil
			}

			accountKey := storage.keys[index]

			if !accountKey.IsRevoked {
				storage.unrevokedKeyCount -= 1
			}

			accountKey.IsRevoked = true

			storage.keys[index] = accountKey
			storage.returnedKey = accountKey

			return accountKey, nil
		},
		OnAccountKeysCount: func(address Address) (uint64, error) {
			return uint64(storage.unrevokedKeyCount), nil
		},
		OnProgramLog: func(message string) {
			storage.logs = append(storage.logs, message)
		},
		OnEmitEvent: func(event cadence.Event) error {
			storage.events = append(storage.events, event)
			return nil
		},
		OnDecodeArgument: func(b []byte, t cadence.Type) (value cadence.Value, err error) {
			return json.Decode(nil, b)
		},
	}
}

func addAuthAccountKey(t *testing.T, runtime Runtime, runtimeInterface *TestRuntimeInterface, location Location) {
	test := accountKeyTestCase{
		name: "Add key",
		code: `
                transaction {
                    prepare(signer: auth(AddKey) &Account) {
                        let key = PublicKey(
                            publicKey: "010203".decodeHex(),
                            signatureAlgorithm: SignatureAlgorithm.ECDSA_P256
                        )

                        var addedKey: AccountKey = signer.keys.add(
                            publicKey: key,
                            hashAlgorithm: HashAlgorithm.SHA3_256,
                            weight: 100.0
                        )
                    }
                }`,
		args: []cadence.Value{},
	}

	err := test.executeTransaction(runtime, runtimeInterface, location)
	require.NoError(t, err)
}

func addPublicKeyValidation(runtimeInterface *TestRuntimeInterface, returnError error) {
	runtimeInterface.OnValidatePublicKey = func(_ *stdlib.PublicKey) error {
		return returnError
	}
}

func encodeArgs(argValues []cadence.Value) [][]byte {
	args := make([][]byte, len(argValues))
	for i, arg := range argValues {
		var err error
		args[i], err = json.Encode(arg)
		if err != nil {
			panic(errors.NewUnexpectedError("broken test: invalid argument: %w", err))
		}
	}
	return args
}

type accountKeyTestCase struct {
	name string
	code string
	args []cadence.Value
}

func (test accountKeyTestCase) executeTransaction(
	runtime Runtime,
	runtimeInterface *TestRuntimeInterface,
	location Location,
) error {
	args := encodeArgs(test.args)

	err := runtime.ExecuteTransaction(
		Script{
			Source:    []byte(test.code),
			Arguments: args,
		},
		Context{
			Interface: runtimeInterface,
			Location:  location,
		},
	)
	return err
}

func (test accountKeyTestCase) executeScript(
	runtime Runtime,
	runtimeInterface *TestRuntimeInterface,
) (cadence.Value, error) {

	args := encodeArgs(test.args)

	value, err := runtime.ExecuteScript(
		Script{
			Source:    []byte(test.code),
			Arguments: args,
		},
		Context{
			Interface: runtimeInterface,
			Location:  common.ScriptLocation{},
		},
	)

	return value, err
}

func newTestAccountKeyStorage() *testAccountKeyStorage {
	return &testAccountKeyStorage{
		events:            make([]cadence.Event, 0),
		keys:              make([]*stdlib.AccountKey, 0),
		unrevokedKeyCount: 0,
	}
}

type testAccountKeyStorage struct {
	returnedKey       *stdlib.AccountKey
	events            []cadence.Event
	keys              []*stdlib.AccountKey
	logs              []string
	unrevokedKeyCount int
}

func TestRuntimePublicKey(t *testing.T) {

	t.Parallel()

	executeScript := func(code string, runtimeInterface Interface) (cadence.Value, error) {
		rt := NewTestInterpreterRuntime()

		value, err := rt.ExecuteScript(
			Script{
				Source: []byte(code),
			},
			Context{
				Interface: runtimeInterface,
				Location:  common.ScriptLocation{},
			},
		)

		return value, err
	}

	t.Run("Constructor", func(t *testing.T) {
		t.Parallel()

		script := `
            access(all) fun main(): PublicKey {
                let publicKey = PublicKey(
                    publicKey: "0102".decodeHex(),
                    signatureAlgorithm: SignatureAlgorithm.ECDSA_P256
                )

                return publicKey
            }
        `

		storage := NewTestLedger(nil, nil)

		runtimeInterface := &TestRuntimeInterface{
			Storage: storage,
		}
		addPublicKeyValidation(runtimeInterface, nil)

		value, err := executeScript(script, runtimeInterface)
		require.NoError(t, err)

		expected := cadence.Struct{
			StructType: PublicKeyType,
			Fields: []cadence.Value{
				// Public key (bytes)
				newBytesValue([]byte{1, 2}),

				// Signature Algo
				newSignAlgoValue(sema.SignatureAlgorithmECDSA_P256),
			},
		}

		assert.Equal(t, expected, value)
	})

	t.Run("Validate func", func(t *testing.T) {
		t.Parallel()

		script := `
            access(all) fun main(): Bool {
                let publicKey =  PublicKey(
                    publicKey: "0102".decodeHex(),
                    signatureAlgorithm: SignatureAlgorithm.ECDSA_P256
                )

                return publicKey.validate()
            }
        `

		runtimeInterface := &TestRuntimeInterface{}

		_, err := executeScript(script, runtimeInterface)
		RequireError(t, err)

		assert.Contains(t, err.Error(), "value of type `PublicKey` has no member `validate`")
	})

	t.Run("Construct PublicKey in Cadence code", func(t *testing.T) {
		t.Parallel()

		script := `
          access(all) fun main(): PublicKey {
              let publicKey = PublicKey(
                  publicKey: "0102".decodeHex(),
                  signatureAlgorithm: SignatureAlgorithm.ECDSA_P256
              )

              return publicKey
          }
        `

		fakeError := &fakeError{}
		for _, errorToReturn := range []error{fakeError, nil} {
			var invoked bool

			storage := NewTestLedger(nil, nil)

			runtimeInterface := &TestRuntimeInterface{
				Storage: storage,
				OnValidatePublicKey: func(publicKey *stdlib.PublicKey) error {
					invoked = true
					return errorToReturn
				},
			}

			value, err := executeScript(script, runtimeInterface)

			assert.True(t, invoked, "validatePublicKey was not invoked")

			if errorToReturn == nil {
				assert.NotNil(t, value)
				assert.NoError(t, err)
			} else {
				assert.Nil(t, value)
				RequireError(t, err)

				assert.ErrorAs(t, err, &errorToReturn)
				assert.ErrorAs(t, err, &interpreter.InvalidPublicKeyError{})
			}
		}
	})

	t.Run("PublicKey from host env", func(t *testing.T) {
		t.Parallel()

		storage := newTestAccountKeyStorage()
		storage.keys = append(storage.keys, accountKeyA, accountKeyB)

		for index := range storage.keys {
			script := fmt.Sprintf(
				`
                  access(all) fun main(): PublicKey {
                      // Get a public key from host env
                      let acc = getAccount(0x02)
                      let publicKey = acc.keys.get(keyIndex: %d)!.publicKey
                      return publicKey
                  }
                `,
				index,
			)

			var invoked bool

			runtimeInterface := getAccountKeyTestRuntimeInterface(storage)
			runtimeInterface.OnValidatePublicKey = func(publicKey *stdlib.PublicKey) error {
				invoked = true
				return nil
			}

			value, err := executeScript(script, runtimeInterface)

			// skip validation when key comes from host env aka FVM
			assert.False(t, invoked, "validatePublicKey was not invoked")

			assert.IsType(t, cadence.Struct{}, value)
			assert.Nil(t, err)
		}
	})

	t.Run("Verify", func(t *testing.T) {
		t.Parallel()

		script := `
            access(all) fun main(): Bool {
                let publicKey =  PublicKey(
                    publicKey: "0102".decodeHex(),
                    signatureAlgorithm: SignatureAlgorithm.ECDSA_P256
                )

                return publicKey.verify(
                    signature: [],
                    signedData: [],
                    domainSeparationTag: "something",
                    hashAlgorithm: HashAlgorithm.SHA2_256
                )
            }
        `

		var invoked bool

		storage := NewTestLedger(nil, nil)

		runtimeInterface := &TestRuntimeInterface{
			Storage: storage,
			OnVerifySignature: func(
				_ []byte,
				_ string,
				_ []byte,
				_ []byte,
				_ SignatureAlgorithm,
				_ HashAlgorithm,
			) (bool, error) {
				invoked = true
				return true, nil
			},
		}
		addPublicKeyValidation(runtimeInterface, nil)

		value, err := executeScript(script, runtimeInterface)
		require.NoError(t, err)

		assert.True(t, invoked)
		assert.Equal(t, cadence.Bool(true), value)
	})

	t.Run("Verify - publicKey from host env", func(t *testing.T) {
		t.Parallel()

		storage := newTestAccountKeyStorage()
		storage.keys = append(storage.keys, accountKeyA, accountKeyB)

		script := `
            access(all) fun main(): Bool {
                // Get a public key from host env
                let acc = getAccount(0x02)
                let publicKey = acc.keys.get(keyIndex: 0)!.publicKey

                return publicKey.verify(
                    signature: [],
                    signedData: [],
                    domainSeparationTag: "something",
                    hashAlgorithm: HashAlgorithm.SHA2_256
                )
            }
        `

		var invoked bool

		runtimeInterface := getAccountKeyTestRuntimeInterface(storage)
		runtimeInterface.OnVerifySignature = func(
			_ []byte,
			_ string,
			_ []byte,
			_ []byte,
			_ SignatureAlgorithm,
			_ HashAlgorithm,
		) (bool, error) {
			invoked = true
			return true, nil
		}

		value, err := executeScript(script, runtimeInterface)
		require.NoError(t, err)

		assert.True(t, invoked)
		assert.Equal(t, cadence.Bool(true), value)
	})

	t.Run("field mutability", func(t *testing.T) {
		t.Parallel()

		script := `
            access(all) fun main(): PublicKey {
                let publicKey =  PublicKey(
                    publicKey: "0102".decodeHex(),
                    signatureAlgorithm: SignatureAlgorithm.ECDSA_P256
                )

                publicKey.publicKey = []
                publicKey.signatureAlgorithm = SignatureAlgorithm.ECDSA_secp256k1

                return publicKey
            }
        `

		storage := NewTestLedger(nil, nil)

		runtimeInterface := &TestRuntimeInterface{
			Storage: storage,
		}

		_, err := executeScript(script, runtimeInterface)
		errs := checker.RequireCheckerErrors(t, err, 4)

		assert.IsType(t, &sema.InvalidAssignmentAccessError{}, errs[0])
		assert.IsType(t, &sema.AssignmentToConstantMemberError{}, errs[1])
		assert.IsType(t, &sema.InvalidAssignmentAccessError{}, errs[2])
		assert.IsType(t, &sema.AssignmentToConstantMemberError{}, errs[3])
	})

	t.Run("raw-key mutability", func(t *testing.T) {
		t.Parallel()

		script := `
            access(all) fun main(): PublicKey {
                let publicKey =  PublicKey(
                    publicKey: "0102".decodeHex(),
                    signatureAlgorithm: SignatureAlgorithm.ECDSA_P256
                )

                publicKey.publicKey[0] = 5

                return publicKey
            }
        `

		storage := NewTestLedger(nil, nil)

		runtimeInterface := &TestRuntimeInterface{
			Storage: storage,
		}
		addPublicKeyValidation(runtimeInterface, nil)

		value, err := executeScript(script, runtimeInterface)
		require.NoError(t, err)

		expected := cadence.Struct{
			StructType: PublicKeyType,
			Fields: []cadence.Value{
				// Public key (bytes)
				newBytesValue([]byte{1, 2}),

				// Signature Algo
				newSignAlgoValue(sema.SignatureAlgorithmECDSA_P256),
			},
		}

		assert.Equal(t, expected, value)
	})

	t.Run("raw-key reference mutability", func(t *testing.T) {
		t.Parallel()

		script := `
          access(all) fun main(): PublicKey {
            let publicKey =  PublicKey(
                publicKey: "0102".decodeHex(),
                signatureAlgorithm: SignatureAlgorithm.ECDSA_P256
            )

            var publickeyRef = &publicKey.publicKey as auth(Mutate) &[UInt8]
            publickeyRef[0] = 3

            return publicKey
          }
        `

		storage := NewTestLedger(nil, nil)

		runtimeInterface := &TestRuntimeInterface{
			Storage: storage,
		}
		addPublicKeyValidation(runtimeInterface, nil)

		value, err := executeScript(script, runtimeInterface)
		require.NoError(t, err)

		expected := cadence.Struct{
			StructType: PublicKeyType,
			Fields: []cadence.Value{
				// Public key (bytes)
				newBytesValue([]byte{1, 2}),
				// Signature Algo
				newSignAlgoValue(sema.SignatureAlgorithmECDSA_P256),
			},
		}

		assert.Equal(t, expected, value)
	})

}

func TestRuntimeAuthAccountContracts(t *testing.T) {

	t.Parallel()

	t.Run("get existing contract", func(t *testing.T) {
		t.Parallel()

		rt := NewTestInterpreterRuntime()

		script := []byte(`
            transaction {
                prepare(acc: &Account) {
                    let deployedContract = acc.contracts.get(name: "foo")
                    assert(deployedContract!.name == "foo")
                }
            }
        `)

		var invoked bool

		runtimeInterface := &TestRuntimeInterface{
			OnGetSigningAccounts: func() ([]Address, error) {
				return []Address{{42}}, nil
			},
			OnGetAccountContractCode: func(_ common.AddressLocation) ([]byte, error) {
				invoked = true
				return []byte{1, 2}, nil
			},
		}

		nextTransactionLocation := NewTransactionLocationGenerator()

		err := rt.ExecuteTransaction(
			Script{
				Source: script,
			},
			Context{
				Interface: runtimeInterface,
				Location:  nextTransactionLocation(),
			},
		)

		require.NoError(t, err)
		assert.True(t, invoked)
	})

	t.Run("get non-existing contract", func(t *testing.T) {
		t.Parallel()

		rt := NewTestInterpreterRuntime()

		script := []byte(`
            transaction {
                prepare(acc: &Account) {
                    let deployedContract = acc.contracts.get(name: "foo")
                    assert(deployedContract == nil)
                }
            }
        `)

		var invoked bool

		runtimeInterface := &TestRuntimeInterface{
			OnGetSigningAccounts: func() ([]Address, error) {
				return []Address{{42}}, nil
			},
			OnGetAccountContractCode: func(_ common.AddressLocation) ([]byte, error) {
				invoked = true
				return nil, nil
			},
		}

		nextTransactionLocation := NewTransactionLocationGenerator()

		err := rt.ExecuteTransaction(
			Script{
				Source: script,
			},
			Context{
				Interface: runtimeInterface,
				Location:  nextTransactionLocation(),
			},
		)

		require.NoError(t, err)
		assert.True(t, invoked)
	})

	t.Run("borrow existing contract", func(t *testing.T) {
		t.Parallel()

		rt := NewTestInterpreterRuntime()

		accountCodes := map[Location][]byte{}
		var events []cadence.Event

		runtimeInterface := &TestRuntimeInterface{
			OnGetCode: func(location Location) (bytes []byte, err error) {
				return accountCodes[location], nil
			},
			Storage: NewTestLedger(nil, nil),
			OnGetSigningAccounts: func() ([]Address, error) {
				return []Address{{0, 0, 0, 0, 0, 0, 0, 0x42}}, nil
			},
			OnResolveLocation: NewSingleIdentifierLocationResolver(t),
			OnGetAccountContractCode: func(location common.AddressLocation) (code []byte, err error) {
				return accountCodes[location], nil
			},
			OnUpdateAccountContractCode: func(location common.AddressLocation, code []byte) error {
				accountCodes[location] = code
				return nil
			},
			OnEmitEvent: func(event cadence.Event) error {
				events = append(events, event)
				return nil
			},
		}

		nextTransactionLocation := NewTransactionLocationGenerator()

		// Deploy contract interface
		err := rt.ExecuteTransaction(
			Script{
				Source: DeploymentTransaction("HelloInterface", []byte(`
                  access(all) contract interface HelloInterface {

                      access(all) fun hello(): String
                  }
                `)),
			},
			Context{
				Interface: runtimeInterface,
				Location:  nextTransactionLocation(),
			},
		)
		require.NoError(t, err)

		// Deploy concrete contract
		err = rt.ExecuteTransaction(
			Script{
				Source: DeploymentTransaction("Hello", []byte(`
                  import HelloInterface from 0x42

                  access(all) contract Hello: HelloInterface {

                      access(all) fun hello(): String {
                          return "Hello!"
                      }
                  }
                `)),
			},
			Context{
				Interface: runtimeInterface,
				Location:  nextTransactionLocation(),
			},
		)
		require.NoError(t, err)

		// Test usage

		err = rt.ExecuteTransaction(
			Script{
				Source: []byte(`
                  import HelloInterface from 0x42

                  transaction {
                      prepare(acc: &Account) {
                          let hello = acc.contracts.borrow<&HelloInterface>(name: "Hello")
                          assert(hello?.hello() == "Hello!")
                      }
                  }
              `),
			},
			Context{
				Interface: runtimeInterface,
				Location:  nextTransactionLocation(),
			},
		)
		require.NoError(t, err)
	})

	t.Run("borrow existing contract with incorrect type", func(t *testing.T) {
		t.Parallel()

		rt := NewTestInterpreterRuntime()

		accountCodes := map[Location][]byte{}
		var events []cadence.Event

		runtimeInterface := &TestRuntimeInterface{
			OnGetCode: func(location Location) (bytes []byte, err error) {
				return accountCodes[location], nil
			},
			Storage: NewTestLedger(nil, nil),
			OnGetSigningAccounts: func() ([]Address, error) {
				return []Address{{0, 0, 0, 0, 0, 0, 0, 0x42}}, nil
			},
			OnResolveLocation: NewSingleIdentifierLocationResolver(t),
			OnGetAccountContractCode: func(location common.AddressLocation) (code []byte, err error) {
				return accountCodes[location], nil
			},
			OnUpdateAccountContractCode: func(location common.AddressLocation, code []byte) error {
				accountCodes[location] = code
				return nil
			},
			OnEmitEvent: func(event cadence.Event) error {
				events = append(events, event)
				return nil
			},
		}

		nextTransactionLocation := NewTransactionLocationGenerator()

		// Deploy contract interface
		err := rt.ExecuteTransaction(
			Script{
				Source: DeploymentTransaction("HelloInterface", []byte(`
                  access(all) contract interface HelloInterface {

                      access(all) fun hello(): String
                  }
                `)),
			},
			Context{
				Interface: runtimeInterface,
				Location:  nextTransactionLocation(),
			},
		)
		require.NoError(t, err)

		// Deploy concrete contract
		err = rt.ExecuteTransaction(
			Script{
				Source: DeploymentTransaction("Hello", []byte(`
                  access(all) contract Hello {

                      access(all) fun hello(): String {
                          return "Hello!"
                      }
                  }
                `)),
			},
			Context{
				Interface: runtimeInterface,
				Location:  nextTransactionLocation(),
			},
		)
		require.NoError(t, err)

		// Test usage

		err = rt.ExecuteTransaction(
			Script{
				Source: []byte(`
                  import HelloInterface from 0x42

                  transaction {
                      prepare(acc: &Account) {
                          let hello = acc.contracts.borrow<&HelloInterface>(name: "Hello")
                          assert(hello == nil)
                      }
                  }
              `),
			},
			Context{
				Interface: runtimeInterface,
				Location:  nextTransactionLocation(),
			},
		)
		require.NoError(t, err)
	})

	t.Run("borrow non-existing contract", func(t *testing.T) {
		t.Parallel()

		rt := NewTestInterpreterRuntime()

		runtimeInterface := &TestRuntimeInterface{
			Storage: NewTestLedger(nil, nil),
			OnGetSigningAccounts: func() ([]Address, error) {
				return []Address{{0, 0, 0, 0, 0, 0, 0, 0x42}}, nil
			},
			OnGetAccountContractCode: func(_ common.AddressLocation) ([]byte, error) {
				return nil, nil
			},
		}

		nextTransactionLocation := NewTransactionLocationGenerator()

		err := rt.ExecuteTransaction(
			Script{
				Source: []byte(`
                  transaction {
                      prepare(acc: &Account) {
                          let hello = acc.contracts.borrow<&AnyStruct>(name: "Hello")
                          assert(hello == nil)
                      }
                  }
              `),
			},
			Context{
				Interface: runtimeInterface,
				Location:  nextTransactionLocation(),
			},
		)
		require.NoError(t, err)
	})

	t.Run("get names", func(t *testing.T) {
		t.Parallel()

		rt := NewTestInterpreterRuntime()

		script := []byte(`
            transaction {
                prepare(signer: &Account) {
                    let names = signer.contracts.names

                    assert(names.isInstance(Type<[String]>()))
                    assert(names.length == 2)
                }
            }
        `)

		var invoked bool

		runtimeInterface := &TestRuntimeInterface{
			OnGetSigningAccounts: func() ([]Address, error) {
				return []Address{{42}}, nil
			},
			OnGetAccountContractNames: func(_ Address) ([]string, error) {
				invoked = true
				return []string{"foo", "bar"}, nil
			},
		}

		nextTransactionLocation := NewTransactionLocationGenerator()

		err := rt.ExecuteTransaction(
			Script{
				Source: script,
			},
			Context{
				Interface: runtimeInterface,
				Location:  nextTransactionLocation(),
			},
		)

		require.NoError(t, err)
		assert.True(t, invoked)
	})

	t.Run("update names", func(t *testing.T) {
		t.Parallel()

		rt := NewTestInterpreterRuntime()

		script := []byte(`
            transaction {
                prepare(signer: auth(Mutate) &Account) {
                    signer.contracts.names[0] = "baz"
                }
            }
        `)

		runtimeInterface := &TestRuntimeInterface{
			OnGetSigningAccounts: func() ([]Address, error) {
				return []Address{{42}}, nil
			},
			OnGetAccountContractNames: func(_ Address) ([]string, error) {
				return []string{"foo", "bar"}, nil
			},
		}

		nextTransactionLocation := NewTransactionLocationGenerator()

		err := rt.ExecuteTransaction(
			Script{
				Source: script,
			},
			Context{
				Interface: runtimeInterface,
				Location:  nextTransactionLocation(),
			},
		)

		errs := checker.RequireCheckerErrors(t, err, 1)

		assert.IsType(t, &sema.UnauthorizedReferenceAssignmentError{}, errs[0])
	})

	t.Run("update names through reference", func(t *testing.T) {
		t.Parallel()

		rt := NewTestInterpreterRuntime()

		script := []byte(`
            transaction {
                prepare(signer: auth(Mutate) &Account) {
                    let namesRef = signer.contracts.names
                    namesRef[0] = "baz"

                    assert(signer.contracts.names[0] == "foo")
                }
            }
        `)

		runtimeInterface := &TestRuntimeInterface{
			OnGetSigningAccounts: func() ([]Address, error) {
				return []Address{{42}}, nil
			},
			OnGetAccountContractNames: func(_ Address) ([]string, error) {
				return []string{"foo", "bar"}, nil
			},
		}

		nextTransactionLocation := NewTransactionLocationGenerator()

		err := rt.ExecuteTransaction(
			Script{
				Source: script,
			},
			Context{
				Interface: runtimeInterface,
				Location:  nextTransactionLocation(),
			},
		)

		errs := checker.RequireCheckerErrors(t, err, 1)

		assert.IsType(t, &sema.UnauthorizedReferenceAssignmentError{}, errs[0])
	})
}

func TestRuntimePublicAccountContracts(t *testing.T) {

	t.Parallel()

	t.Run("get existing contract", func(t *testing.T) {
		t.Parallel()

		rt := NewTestInterpreterRuntime()

		script := []byte(`
            access(all) fun main(): [AnyStruct] {
                let acc = getAccount(0x02)
                let deployedContract = acc.contracts.get(name: "foo")

                return [deployedContract!.name, deployedContract!.code]
            }
        `)

		var invoked bool

		runtimeInterface := &TestRuntimeInterface{
			OnGetSigningAccounts: func() ([]Address, error) {
				return []Address{{42}}, nil
			},
			OnGetAccountContractCode: func(_ common.AddressLocation) ([]byte, error) {
				invoked = true
				return []byte{1, 2}, nil
			},
		}

		result, err := rt.ExecuteScript(
			Script{
				Source: script,
			},
			Context{
				Interface: runtimeInterface,
				Location:  common.ScriptLocation{},
			},
		)

		require.NoError(t, err)
		assert.True(t, invoked)

		require.IsType(t, cadence.Array{}, result)
		array := result.(cadence.Array)

		require.Len(t, array.Values, 2)

		assert.Equal(t, cadence.String("foo"), array.Values[0])
		assert.Equal(t,
			cadence.Array{
				Values: []cadence.Value{
					cadence.UInt8(1),
					cadence.UInt8(2),
				},
			}.WithType(&cadence.VariableSizedArrayType{
				ElementType: cadence.UInt8Type,
			}),
			array.Values[1],
		)
	})

	t.Run("get non-existing contract", func(t *testing.T) {
		t.Parallel()

		rt := NewTestInterpreterRuntime()

		script := []byte(`
            access(all) fun main() {
                let acc = getAccount(0x02)
                assert(acc.contracts.get(name: "foo") == nil)
            }
        `)

		var invoked bool

		runtimeInterface := &TestRuntimeInterface{
			OnGetSigningAccounts: func() ([]Address, error) {
				return []Address{{42}}, nil
			},
			OnGetAccountContractCode: func(_ common.AddressLocation) ([]byte, error) {
				invoked = true
				return nil, nil
			},
		}

		_, err := rt.ExecuteScript(
			Script{
				Source: script,
			},
			Context{
				Interface: runtimeInterface,
				Location:  common.ScriptLocation{},
			},
		)

		require.NoError(t, err)
		assert.True(t, invoked)
	})

	t.Run("get names", func(t *testing.T) {
		t.Parallel()

		rt := NewTestInterpreterRuntime()

		script := []byte(`
            access(all) fun main(): &[String] {
                let acc = getAccount(0x02)
                return acc.contracts.names
            }
        `)

		var invoked bool

		runtimeInterface := &TestRuntimeInterface{
			OnGetSigningAccounts: func() ([]Address, error) {
				return []Address{{42}}, nil
			},
			OnGetAccountContractNames: func(_ Address) ([]string, error) {
				invoked = true
				return []string{"foo", "bar"}, nil
			},
		}

		result, err := rt.ExecuteScript(
			Script{
				Source: script,
			},
			Context{
				Interface: runtimeInterface,
				Location:  common.ScriptLocation{},
			},
		)

		require.NoError(t, err)
		assert.True(t, invoked)

		require.IsType(t, cadence.Array{}, result)
		array := result.(cadence.Array)

		require.Len(t, array.Values, 2)
		assert.Equal(t, cadence.String("foo"), array.Values[0])
		assert.Equal(t, cadence.String("bar"), array.Values[1])
	})
}

func TestRuntimeGetAuthAccount(t *testing.T) {

	t.Parallel()

	t.Run("script", func(t *testing.T) {
		t.Parallel()

		rt := NewTestInterpreterRuntime()

		script := []byte(`
            access(all) fun main(): UInt64 {
                let acc = getAccount(0x02)
                return acc.storage.used
            }
        `)

		runtimeInterface := &TestRuntimeInterface{
			OnGetStorageUsed: func(_ Address) (uint64, error) {
				return 1, nil
			},
		}

		result, err := rt.ExecuteScript(
			Script{
				Source: script,
			},
			Context{
				Interface: runtimeInterface,
				Location:  common.ScriptLocation{0x1},
			},
		)

		require.NoError(t, err)
		assert.Equal(t, cadence.UInt64(0x1), result)
	})

	t.Run("incorrect arg type", func(t *testing.T) {
		t.Parallel()

		rt := NewTestInterpreterRuntime()

		script := []byte(`
            access(all) fun main() {
                let acc = getAuthAccount<&Account>("")
            }
        `)

		runtimeInterface := &TestRuntimeInterface{}

		_, err := rt.ExecuteScript(
			Script{
				Source: script,
			},
			Context{
				Interface: runtimeInterface,
				Location:  common.ScriptLocation{0x1},
			},
		)

		errs := checker.RequireCheckerErrors(t, err, 1)

		assert.IsType(t, &sema.TypeMismatchError{}, errs[0])
	})

	t.Run("no args", func(t *testing.T) {
		t.Parallel()

		rt := NewTestInterpreterRuntime()

		script := []byte(`
            access(all) fun main() {
                let acc = getAuthAccount<&Account>()
            }
        `)

		runtimeInterface := &TestRuntimeInterface{}

		_, err := rt.ExecuteScript(
			Script{
				Source: script,
			},
			Context{
				Interface: runtimeInterface,
				Location:  common.ScriptLocation{0x1},
			},
		)

		errs := checker.RequireCheckerErrors(t, err, 1)

		assert.IsType(t, &sema.InsufficientArgumentsError{}, errs[0])
	})

	t.Run("too many args", func(t *testing.T) {
		t.Parallel()

		rt := NewTestInterpreterRuntime()

		script := []byte(`
            access(all) fun main() {
                let acc = getAuthAccount<&Account>(0x1, 0x2)
            }
        `)

		runtimeInterface := &TestRuntimeInterface{}

		_, err := rt.ExecuteScript(
			Script{
				Source: script,
			},
			Context{
				Interface: runtimeInterface,
				Location:  common.ScriptLocation{0x1},
			},
		)
		errs := checker.RequireCheckerErrors(t, err, 1)

		assert.IsType(t, &sema.ExcessiveArgumentsError{}, errs[0])
	})

	t.Run("transaction", func(t *testing.T) {
		t.Parallel()

		rt := NewTestInterpreterRuntime()

		script := []byte(`
            transaction {
                prepare() {
                    let acc = getAuthAccount<&Account>(0x02)
                    log(acc.storage.used)
                }
            }
        `)

		runtimeInterface := &TestRuntimeInterface{
			OnGetStorageUsed: func(_ Address) (uint64, error) {
				return 1, nil
			},
		}

		err := rt.ExecuteTransaction(
			Script{
				Source: script,
			},
			Context{
				Interface: runtimeInterface,
				Location:  common.TransactionLocation{0x1},
			},
		)

		errs := checker.RequireCheckerErrors(t, err, 1)

		assert.IsType(t, &sema.NotDeclaredError{}, errs[0])
	})
}

type fakeError struct{}

func (fakeError) Error() string {
	return "fake error for testing"
}

// TestRuntimePublicKeyPublicKeyField tests PublicKey's `publicKey` field.
// It is a computed field, which always returns a copy of the stored raw public key ([UInt8]).
//
// This test ensures that the field can be accessed even after the PublicKey value has been transferred (copied),
// and the original PublicKey value has been removed.
func TestRuntimePublicKeyPublicKeyField(t *testing.T) {

	t.Parallel()

	inter := NewTestInterpreter(t)

	locationRange := interpreter.EmptyLocationRange

	publicKey := interpreter.NewCompositeValue(
		inter,
		locationRange,
		nil,
		sema.PublicKeyType.Identifier,
		common.CompositeKindStructure,
		[]interpreter.CompositeField{
			{
				Name: sema.PublicKeyTypePublicKeyFieldName,
				Value: interpreter.NewArrayValue(
					inter,
					locationRange,
					interpreter.ByteArrayStaticType,
					common.ZeroAddress,
					interpreter.NewUnmeteredUInt8Value(1),
				),
			},
		},
		common.ZeroAddress,
	)

	publicKeyArray1 := publicKey.GetMember(
		inter,
		locationRange,
		sema.PublicKeyTypePublicKeyFieldName,
	)

	publicKey2 := publicKey.Transfer(
		inter,
		locationRange,
		atree.Address{},
		false,
		nil,
		nil,
	).(*interpreter.CompositeValue)

	publicKey.DeepRemove(inter)

	publicKeyArray2 := publicKey2.GetMember(
		inter,
		locationRange,
		sema.PublicKeyTypePublicKeyFieldName,
	)

	require.True(t,
		publicKeyArray2.(interpreter.EquatableValue).
			Equal(inter, locationRange, publicKeyArray1),
	)
}
