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
	"bytes"
	"math"
	"math/big"

	"github.com/fxamacker/cbor/v2"
	"github.com/onflow/atree"

	"github.com/onflow/cadence/runtime/common"
	"github.com/onflow/cadence/runtime/errors"
)

const cborTagSize = 2

var bigOne = big.NewInt(1)

func getBigIntCBORSize(v *big.Int) uint32 {
	sign := v.Sign()
	if sign < 0 {
		v = new(big.Int).Abs(v)
		v.Sub(v, bigOne)
	}

	// tag number + bytes
	return 1 + getBytesCBORSize(v.Bytes())
}

func getIntCBORSize(v int64) uint32 {
	if v < 0 {
		return getUintCBORSize(uint64(-v - 1))
	}
	return getUintCBORSize(uint64(v))
}

func getUintCBORSize(v uint64) uint32 {
	if v <= 23 {
		return 1
	}
	if v <= math.MaxUint8 {
		return 2
	}
	if v <= math.MaxUint16 {
		return 3
	}
	if v <= math.MaxUint32 {
		return 5
	}
	return 9
}

func getBytesCBORSize(b []byte) uint32 {
	length := len(b)
	if length == 0 {
		return 1
	}
	return getUintCBORSize(uint64(length)) + uint32(length)
}

// Cadence needs to encode different kinds of objects in CBOR, for instance,
// dictionaries, structs, resources, etc.
//
// However, CBOR only provides one native map type, and no support
// for directly representing e.g. structs or resources.
//
// To be able to encode/decode such semantically different values,
// we define custom CBOR tags.

// !!! *WARNING* !!!
//
// Only add new fields to encoded structs by
// appending new fields with the next highest key.
//
// DO *NOT* REPLACE EXISTING FIELDS!

const CBORTagBase = 128

// !!! *WARNING* !!!
//
// Only add new types by:
// - replacing existing placeholders (`_`) with new types
// - appending new types
//
// Only remove types by:
// - replace existing types with a placeholder `_`
//
// DO *NOT* REPLACE EXISTING TYPES!
// DO *NOT* ADD NEW TYPES IN BETWEEN!

const (
	CBORTagVoidValue = CBORTagBase + iota
	_                // DO *NOT* REPLACE. Previously used for dictionary values
	CBORTagSomeValue
	CBORTagAddressValue
	CBORTagCompositeValue
	CBORTagTypeValue
	_ // DO *NOT* REPLACE. Previously used for array values
	CBORTagStringValue
	CBORTagCharacterValue
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_

	// Int*
	CBORTagIntValue
	CBORTagInt8Value
	CBORTagInt16Value
	CBORTagInt32Value
	CBORTagInt64Value
	CBORTagInt128Value
	CBORTagInt256Value
	_

	// UInt*
	CBORTagUIntValue
	CBORTagUInt8Value
	CBORTagUInt16Value
	CBORTagUInt32Value
	CBORTagUInt64Value
	CBORTagUInt128Value
	CBORTagUInt256Value
	_

	// Word*
	_
	CBORTagWord8Value
	CBORTagWord16Value
	CBORTagWord32Value
	CBORTagWord64Value
	CBORTagWord128Value
	CBORTagWord256Value
	_

	// Fix*
	_
	_ // future: Fix8
	_ // future: Fix16
	_ // future: Fix32
	CBORTagFix64Value
	_ // future: Fix128
	_ // future: Fix256
	_

	// UFix*
	_
	_ // future: UFix8
	_ // future: UFix16
	_ // future: UFix32
	CBORTagUFix64Value
	_ // future: UFix128
	_ // future: UFix256
	_

	// Locations
	CBORTagAddressLocation
	CBORTagStringLocation
	CBORTagIdentifierLocation
	CBORTagTransactionLocation
	CBORTagScriptLocation
	_
	_
	_

	// Storage

	CBORTagPathValue
	// Deprecated: CBORTagPathCapabilityValue
	CBORTagPathCapabilityValue
	_ // DO NOT REPLACE! used to be used for storage references
	// Deprecated: CBORTagPathLinkValue
	CBORTagPathLinkValue
	CBORTagPublishedValue
	// Deprecated: CBORTagAccountLinkValue
	CBORTagAccountLinkValue
	CBORTagStorageCapabilityControllerValue
	CBORTagAccountCapabilityControllerValue
	CBORTagCapabilityValue
	_
	_
	_

	// Static Types
	CBORTagPrimitiveStaticType
	CBORTagCompositeStaticType
	CBORTagInterfaceStaticType
	CBORTagVariableSizedStaticType
	CBORTagConstantSizedStaticType
	CBORTagDictionaryStaticType
	CBORTagOptionalStaticType
	CBORTagReferenceStaticType
	CBORTagIntersectionStaticType
	CBORTagCapabilityStaticType
	CBORTagUnauthorizedStaticAuthorization
	CBORTagEntitlementMapStaticAuthorization
	CBORTagEntitlementSetStaticAuthorization
	CBORTagInclusiveRangeStaticType

	// !!! *WARNING* !!!
	// ADD NEW TYPES *BEFORE* THIS WARNING.
	// DO *NOT* ADD NEW TYPES AFTER THIS LINE!
	CBORTag_Count
)

// CBOREncMode
//
// See https://github.com/fxamacker/cbor:
// "For best performance, reuse EncMode and DecMode after creating them."
var CBOREncMode = func() cbor.EncMode {
	options := cbor.CanonicalEncOptions()
	options.BigIntConvert = cbor.BigIntConvertNone
	encMode, err := options.EncMode()
	if err != nil {
		panic(err)
	}
	return encMode
}()

// Encode encodes the value as a CBOR nil
func (v NilValue) Encode(e *atree.Encoder) error {
	// NOTE: when updating, also update NilValue.ByteSize
	return e.CBOR.EncodeNil()
}

// Encode encodes the value as a CBOR bool
func (v BoolValue) Encode(e *atree.Encoder) error {
	// NOTE: when updating, also update BoolValue.ByteSize
	return e.CBOR.EncodeBool(bool(v))
}

// Encode encodes the value as a CBOR string
func (v CharacterValue) Encode(e *atree.Encoder) error {
	err := e.CBOR.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagCharacterValue,
	})
	if err != nil {
		return err
	}
	return e.CBOR.EncodeString(v.Str)
}

// Encode encodes the value as a CBOR string
func (v *StringValue) Encode(e *atree.Encoder) error {
	err := e.CBOR.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagStringValue,
	})
	if err != nil {
		return err
	}
	return e.CBOR.EncodeString(v.Str)
}

// Encode encodes the value as a CBOR string
func (v StringAtreeValue) Encode(e *atree.Encoder) error {
	return e.CBOR.EncodeString(string(v))
}

// Encode encodes the value as a CBOR unsigned integer
func (v Uint64AtreeValue) Encode(e *atree.Encoder) error {
	return e.CBOR.EncodeUint64(uint64(v))
}

// cborVoidValue represents the CBOR value:
//
//	cbor.Tag{
//		Number: CBORTagVoidValue,
//		Content: nil
//	}
var cborVoidValue = []byte{
	// tag
	0xd8, CBORTagVoidValue,
	// null
	0xf6,
}

// Encode writes a value of type Void to the encoder
func (VoidValue) Encode(e *atree.Encoder) error {

	// TODO: optimize: use 0xf7, but decoded by github.com/fxamacker/cbor/v2 as Go `nil`:
	//   https://github.com/fxamacker/cbor/blob/a6ed6ff68e99cbb076997a08d19f03c453851555/README.md#limitations

	return e.CBOR.EncodeRawBytes(cborVoidValue)
}

// Encode encodes the value as
//
//	cbor.Tag{
//			Number:  CBORTagIntValue,
//			Content: *big.Int(v.BigInt),
//	}
func (v IntValue) Encode(e *atree.Encoder) error {
	err := e.CBOR.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagIntValue,
	})
	if err != nil {
		return err
	}
	return e.CBOR.EncodeBigInt(v.BigInt)
}

// Encode encodes Int8Value as
//
//	cbor.Tag{
//			Number:  CBORTagInt8Value,
//			Content: int8(v),
//	}
func (v Int8Value) Encode(e *atree.Encoder) error {
	err := e.CBOR.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagInt8Value,
	})
	if err != nil {
		return err
	}
	return e.CBOR.EncodeInt8(int8(v))
}

// Encode encodes Int16Value as
//
//	cbor.Tag{
//			Number:  CBORTagInt16Value,
//			Content: int16(v),
//	}
func (v Int16Value) Encode(e *atree.Encoder) error {
	err := e.CBOR.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagInt16Value,
	})
	if err != nil {
		return err
	}
	return e.CBOR.EncodeInt16(int16(v))
}

// Encode encodes Int32Value as
//
//	cbor.Tag{
//			Number:  CBORTagInt32Value,
//			Content: int32(v),
//	}
func (v Int32Value) Encode(e *atree.Encoder) error {
	err := e.CBOR.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagInt32Value,
	})
	if err != nil {
		return err
	}
	return e.CBOR.EncodeInt32(int32(v))
}

// Encode encodes Int64Value as
//
//	cbor.Tag{
//			Number:  CBORTagInt64Value,
//			Content: int64(v),
//	}
func (v Int64Value) Encode(e *atree.Encoder) error {
	err := e.CBOR.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagInt64Value,
	})
	if err != nil {
		return err
	}
	return e.CBOR.EncodeInt64(int64(v))
}

// Encode encodes Int128Value as
//
//	cbor.Tag{
//			Number:  CBORTagInt128Value,
//			Content: *big.Int(v.BigInt),
//	}
func (v Int128Value) Encode(e *atree.Encoder) error {
	err := e.CBOR.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagInt128Value,
	})
	if err != nil {
		return err
	}
	return e.CBOR.EncodeBigInt(v.BigInt)
}

// Encode encodes Int256Value as
//
//	cbor.Tag{
//			Number:  CBORTagInt256Value,
//			Content: *big.Int(v.BigInt),
//	}
func (v Int256Value) Encode(e *atree.Encoder) error {
	err := e.CBOR.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagInt256Value,
	})
	if err != nil {
		return err
	}
	return e.CBOR.EncodeBigInt(v.BigInt)
}

// Encode encodes UIntValue as
//
//	cbor.Tag{
//			Number:  CBORTagUIntValue,
//			Content: *big.Int(v.BigInt),
//	}
func (v UIntValue) Encode(e *atree.Encoder) error {
	err := e.CBOR.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagUIntValue,
	})
	if err != nil {
		return err
	}
	return e.CBOR.EncodeBigInt(v.BigInt)
}

// Encode encodes UInt8Value as
//
//	cbor.Tag{
//			Number:  CBORTagUInt8Value,
//			Content: uint8(v),
//	}
func (v UInt8Value) Encode(e *atree.Encoder) error {
	err := e.CBOR.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagUInt8Value,
	})
	if err != nil {
		return err
	}
	return e.CBOR.EncodeUint8(uint8(v))
}

// Encode encodes UInt16Value as
//
//	cbor.Tag{
//			Number:  CBORTagUInt16Value,
//			Content: uint16(v),
//	}
func (v UInt16Value) Encode(e *atree.Encoder) error {
	err := e.CBOR.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagUInt16Value,
	})
	if err != nil {
		return err
	}
	return e.CBOR.EncodeUint16(uint16(v))
}

// Encode encodes UInt32Value as
//
//	cbor.Tag{
//			Number:  CBORTagUInt32Value,
//			Content: uint32(v),
//	}
func (v UInt32Value) Encode(e *atree.Encoder) error {
	err := e.CBOR.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagUInt32Value,
	})
	if err != nil {
		return err
	}
	return e.CBOR.EncodeUint32(uint32(v))
}

// Encode encodes UInt64Value as
//
//	cbor.Tag{
//			Number:  CBORTagUInt64Value,
//			Content: uint64(v),
//	}
func (v UInt64Value) Encode(e *atree.Encoder) error {
	err := e.CBOR.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagUInt64Value,
	})
	if err != nil {
		return err
	}
	return e.CBOR.EncodeUint64(uint64(v))
}

// Encode encodes UInt128Value as
//
//	cbor.Tag{
//			Number:  CBORTagUInt128Value,
//			Content: *big.Int(v.BigInt),
//	}
func (v UInt128Value) Encode(e *atree.Encoder) error {
	err := e.CBOR.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagUInt128Value,
	})
	if err != nil {
		return err
	}
	return e.CBOR.EncodeBigInt(v.BigInt)
}

// Encode encodes UInt256Value as
//
//	cbor.Tag{
//			Number:  CBORTagUInt256Value,
//			Content: *big.Int(v.BigInt),
//	}
func (v UInt256Value) Encode(e *atree.Encoder) error {
	err := e.CBOR.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagUInt256Value,
	})
	if err != nil {
		return err
	}
	return e.CBOR.EncodeBigInt(v.BigInt)
}

// Encode encodes Word8Value as
//
//	cbor.Tag{
//			Number:  CBORTagWord8Value,
//			Content: uint8(v),
//	}
func (v Word8Value) Encode(e *atree.Encoder) error {
	err := e.CBOR.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagWord8Value,
	})
	if err != nil {
		return err
	}
	return e.CBOR.EncodeUint8(uint8(v))
}

// Encode encodes Word16Value as
//
//	cbor.Tag{
//			Number:  CBORTagWord16Value,
//			Content: uint16(v),
//	}
func (v Word16Value) Encode(e *atree.Encoder) error {
	err := e.CBOR.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagWord16Value,
	})
	if err != nil {
		return err
	}
	return e.CBOR.EncodeUint16(uint16(v))
}

// Encode encodes Word32Value as
//
//	cbor.Tag{
//			Number:  CBORTagWord32Value,
//			Content: uint32(v),
//	}
func (v Word32Value) Encode(e *atree.Encoder) error {
	err := e.CBOR.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagWord32Value,
	})
	if err != nil {
		return err
	}
	return e.CBOR.EncodeUint32(uint32(v))
}

// Encode encodes Word64Value as
//
//	cbor.Tag{
//			Number:  CBORTagWord64Value,
//			Content: uint64(v),
//	}
func (v Word64Value) Encode(e *atree.Encoder) error {
	err := e.CBOR.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagWord64Value,
	})
	if err != nil {
		return err
	}
	return e.CBOR.EncodeUint64(uint64(v))
}

// Encode encodes Word128Value as
//
//	cbor.Tag{
//			Number:  CBORTagWord128Value,
//			Content: *big.Int(v.BigInt),
//	}
func (v Word128Value) Encode(e *atree.Encoder) error {
	err := e.CBOR.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagWord128Value,
	})
	if err != nil {
		return err
	}
	return e.CBOR.EncodeBigInt(v.BigInt)
}

// Encode encodes Word256Value as
//
//	cbor.Tag{
//			Number:  CBORTagWord256Value,
//			Content: *big.Int(v.BigInt),
//	}
func (v Word256Value) Encode(e *atree.Encoder) error {
	err := e.CBOR.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagWord256Value,
	})
	if err != nil {
		return err
	}
	return e.CBOR.EncodeBigInt(v.BigInt)
}

// Encode encodes Fix64Value as
//
//	cbor.Tag{
//			Number:  CBORTagFix64Value,
//			Content: int64(v),
//	}
func (v Fix64Value) Encode(e *atree.Encoder) error {
	err := e.CBOR.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagFix64Value,
	})
	if err != nil {
		return err
	}
	return e.CBOR.EncodeInt64(int64(v))
}

// Encode encodes UFix64Value as
//
//	cbor.Tag{
//			Number:  CBORTagUFix64Value,
//			Content: uint64(v),
//	}
func (v UFix64Value) Encode(e *atree.Encoder) error {
	err := e.CBOR.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagUFix64Value,
	})
	if err != nil {
		return err
	}
	return e.CBOR.EncodeUint64(uint64(v))
}

// Encode encodes SomeStorable as
//
//	cbor.Tag{
//			Number: CBORTagSomeValue,
//			Content: Value(v.Value),
//	}
func (s SomeStorable) Encode(e *atree.Encoder) error {
	// NOTE: when updating, also update SomeStorable.ByteSize
	err := e.CBOR.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagSomeValue,
	})
	if err != nil {
		return err
	}
	return s.Storable.Encode(e)
}

// Encode encodes AddressValue as
//
//	cbor.Tag{
//			Number:  CBORTagAddressValue,
//			Content: []byte(v.ToAddress().Bytes()),
//	}
func (v AddressValue) Encode(e *atree.Encoder) error {
	err := e.CBOR.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagAddressValue,
	})
	if err != nil {
		return err
	}
	return e.CBOR.EncodeBytes(v.ToAddress().Bytes())
}

// NOTE: NEVER change, only add/increment; ensure uint64
const (
	// encodedPathValueDomainFieldKey     uint64 = 0
	// encodedPathValueIdentifierFieldKey uint64 = 1

	// !!! *WARNING* !!!
	//
	// encodedPathValueLength MUST be updated when new element is added.
	// It is used to verify encoded path length during decoding.
	encodedPathValueLength = 2
)

// Encode encodes PathValue as
//
//	cbor.Tag{
//				Number: CBORTagPathValue,
//				Content: []any{
//					encodedPathValueDomainFieldKey:     uint(v.Domain),
//					encodedPathValueIdentifierFieldKey: string(v.Identifier),
//				},
//	}
func (v PathValue) Encode(e *atree.Encoder) error {
	// Encode tag number and array head
	err := e.CBOR.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagPathValue,
		// array, 2 items follow
		0x82,
	})
	if err != nil {
		return err
	}

	// Encode domain at array index encodedPathValueDomainFieldKey
	err = e.CBOR.EncodeUint(uint(v.Domain))
	if err != nil {
		return err
	}

	// Encode identifier at array index encodedPathValueIdentifierFieldKey
	return e.CBOR.EncodeString(v.Identifier)
}

// NOTE: NEVER change, only add/increment; ensure uint64
const (
	// encodedCapabilityValueAddressFieldKey    uint64 = 0
	// encodedCapabilityValueIDFieldKey         uint64 = 1
	// encodedCapabilityValueBorrowTypeFieldKey uint64 = 2

	// !!! *WARNING* !!!
	//
	// encodedCapabilityValueLength MUST be updated when new element is added.
	// It is used to verify encoded capability length during decoding.
	encodedCapabilityValueLength = 3
)

// Encode encodes CapabilityValue as
//
//	cbor.Tag{
//				Number: CBORTagCapabilityValue,
//				Content: []any{
//						encodedCapabilityValueAddressFieldKey:    AddressValue(v.Address),
//						encodedCapabilityValueIDFieldKey:         v.ID,
//						encodedCapabilityValueBorrowTypeFieldKey: StaticType(v.BorrowType),
//					},
//	}
func (v *CapabilityValue) Encode(e *atree.Encoder) error {
	// Encode tag number and array head
	err := e.CBOR.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagCapabilityValue,
		// array, 3 items follow
		0x83,
	})
	if err != nil {
		return err
	}

	// Encode address at array index encodedCapabilityValueAddressFieldKey
	err = v.Address.Encode(e)
	if err != nil {
		return err
	}

	// Encode ID at array index encodedCapabilityValueIDFieldKey
	err = e.CBOR.EncodeUint64(uint64(v.ID))
	if err != nil {
		return err
	}

	// Encode borrow type at array index encodedCapabilityValueBorrowTypeFieldKey
	return v.BorrowType.Encode(e.CBOR)
}

// NOTE: NEVER change, only add/increment; ensure uint64
const (
	// encodedAddressLocationAddressFieldKey uint64 = 0
	// encodedAddressLocationNameFieldKey    uint64 = 1

	// !!! *WARNING* !!!
	//
	// encodedAddressLocationLength MUST be updated when new element is added.
	// It is used to verify encoded address location length during decoding.
	encodedAddressLocationLength = 2
)

func encodeLocation(e *cbor.StreamEncoder, l common.Location) error {
	if l == nil {
		return e.EncodeNil()
	}

	switch l := l.(type) {

	case common.StringLocation:
		// common.StringLocation is encoded as
		// cbor.Tag{
		//		Number:  CBORTagStringLocation,
		//		Content: string(l),
		// }
		err := e.EncodeRawBytes([]byte{
			// tag number
			0xd8, CBORTagStringLocation,
		})
		if err != nil {
			return err
		}

		return e.EncodeString(string(l))

	case common.IdentifierLocation:
		// common.IdentifierLocation is encoded as
		// cbor.Tag{
		//		Number:  CBORTagIdentifierLocation,
		//		Content: string(l),
		// }
		err := e.EncodeRawBytes([]byte{
			// tag number
			0xd8, CBORTagIdentifierLocation,
		})
		if err != nil {
			return err
		}

		return e.EncodeString(string(l))

	case common.AddressLocation:
		// common.AddressLocation is encoded as
		// cbor.Tag{
		//		Number: CBORTagAddressLocation,
		//		Content: []any{
		//			encodedAddressLocationAddressFieldKey: []byte{l.Address.Bytes()},
		//			encodedAddressLocationNameFieldKey:    string(l.Name),
		//		},
		// }
		// Encode tag number and array head
		err := e.EncodeRawBytes([]byte{
			// tag number
			0xd8, CBORTagAddressLocation,
			// array, 2 items follow
			0x82,
		})
		if err != nil {
			return err
		}

		// Encode address at array index encodedAddressLocationAddressFieldKey
		err = e.EncodeBytes(l.Address.Bytes())
		if err != nil {
			return err
		}

		// Encode name at array index encodedAddressLocationNameFieldKey
		return e.EncodeString(l.Name)

	case common.TransactionLocation:
		// common.TransactionLocation is encoded as
		// cbor.Tag{
		//		Number: CBORTagTransactionLocation,
		//		Content: []byte(l),
		// }
		// Encode tag number and array head
		err := e.EncodeRawBytes([]byte{
			// tag number
			0xd8, CBORTagTransactionLocation,
		})
		if err != nil {
			return err
		}

		return e.EncodeBytes(l[:])

	case common.ScriptLocation:
		// common.ScriptLocation is encoded as
		// cbor.Tag{
		//		Number: CBORTagScriptLocation,
		//		Content: []byte(l),
		// }
		// Encode tag number and array head
		err := e.EncodeRawBytes([]byte{
			// tag number
			0xd8, CBORTagScriptLocation,
		})
		if err != nil {
			return err
		}

		return e.EncodeBytes(l[:])

	default:
		return errors.NewUnexpectedError("unsupported location: %T", l)
	}
}

// NOTE: NEVER change, only add/increment; ensure uint64
const (
	// encodedPublishedValueRecipientFieldKey uint64 = 0
	// encodedPublishedValueValueFieldKey     uint64 = 1

	// !!! *WARNING* !!!
	//
	// encodedPublishedValueLength MUST be updated when new element is added.
	// It is used to verify encoded link length during decoding.
	encodedPublishedValueLength = 2
)

// Encode encodes PublishedValue as
//
//	cbor.Tag{
//				Number: CBORTagPublishedValue,
//				Content: []any{
//					encodedPublishedValueRecipientFieldKey: AddressValue(v.Recipient),
//					encodedPublishedValueValueFieldKey:     v.Value,
//				},
//	}
func (v *PublishedValue) Encode(e *atree.Encoder) error {
	// Encode tag number and array head
	err := e.CBOR.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagPublishedValue,
		// array, 2 items follow
		0x82,
	})
	if err != nil {
		return err
	}
	// Encode path at array index encodedPublishedValueRecipientFieldKey
	err = v.Recipient.Encode(e)
	if err != nil {
		return err
	}
	// Encode type at array index encodedPublishedValueValueFieldKey
	return v.Value.Encode(e)
}

// NOTE: NEVER change, only add/increment; ensure uint64
const (
	// encodedTypeValueTypeFieldKey uint64 = 0

	// !!! *WARNING* !!!
	//
	// encodedTypeValueTypeLength MUST be updated when new element is added.
	// It is used to verify encoded type length during decoding.
	encodedTypeValueTypeLength = 1
)

// Encode encodes TypeValue as
//
//	cbor.Tag{
//				Number: CBORTagTypeValue,
//				Content: cborArray{
//					encodedTypeValueTypeFieldKey: StaticType(v.Type),
//				},
//		}
func (v TypeValue) Encode(e *atree.Encoder) error {
	// Encode tag number and array head
	err := e.CBOR.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagTypeValue,
		// array, 1 item follow
		0x81,
	})
	if err != nil {
		return err
	}

	// Encode type at array index encodedTypeValueTypeFieldKey
	if v.Type == nil {
		return e.CBOR.EncodeNil()
	} else {
		return v.Type.Encode(e.CBOR)
	}
}

// NOTE: NEVER change, only add/increment; ensure uint64
const (
	// encodedStorageCapabilityControllerValueBorrowTypeFieldKey   uint64 = 0
	// encodedStorageCapabilityControllerValueCapabilityIDFieldKey uint64 = 1
	// encodedStorageCapabilityControllerValueTargetPathFieldKey   uint64 = 2

	// !!! *WARNING* !!!
	//
	// encodedStorageCapabilityControllerValueLength MUST be updated when new element is added.
	// It is used to verify encoded storage capability controller length during decoding.
	encodedStorageCapabilityControllerValueLength = 3
)

// Encode encodes StorageCapabilityControllerValue as
//
//	cbor.Tag{
//				Number: CBORTagStorageCapabilityControllerValue,
//				Content: []any{
//					encodedStorageCapabilityControllerValueBorrowTypeFieldKey:   StaticType(v.BorrowType),
//					encodedStorageCapabilityControllerValueCapabilityIDFieldKey: UInt64Value(v.CapabilityID),
//					encodedStorageCapabilityControllerValueTargetPathFieldKey:   PathValue(v.TargetPath),
//				},
//	}
func (v *StorageCapabilityControllerValue) Encode(e *atree.Encoder) error {
	// Encode tag number and array head
	err := e.CBOR.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagStorageCapabilityControllerValue,
		// array, 3 items follow
		0x83,
	})
	if err != nil {
		return err
	}

	// Encode borrow type at array index encodedStorageCapabilityControllerValueBorrowTypeFieldKey
	err = v.BorrowType.Encode(e.CBOR)
	if err != nil {
		return err
	}

	// Encode ID at array index encodedStorageCapabilityControllerValueCapabilityIDFieldKey
	err = e.CBOR.EncodeUint64(uint64(v.CapabilityID))
	if err != nil {
		return err
	}

	// Encode target path at array index encodedStorageCapabilityControllerValueTargetPathFieldKey
	return v.TargetPath.Encode(e)
}

// NOTE: NEVER change, only add/increment; ensure uint64
const (
	// encodedAccountCapabilityControllerValueBorrowTypeFieldKey   uint64 = 0
	// encodedAccountCapabilityControllerValueCapabilityIDFieldKey uint64 = 1

	// !!! *WARNING* !!!
	//
	// encodedAccountCapabilityControllerValueLength MUST be updated when new element is added.
	// It is used to verify encoded account capability controller length during decoding.
	encodedAccountCapabilityControllerValueLength = 2
)

// Encode encodes AccountCapabilityControllerValue as
//
//	cbor.Tag{
//				Number: CBORTagAccountCapabilityControllerValue,
//				Content: []any{
//					encodedAccountCapabilityControllerValueBorrowTypeFieldKey:   StaticType(v.BorrowType),
//					encodedAccountCapabilityControllerValueCapabilityIDFieldKey: UInt64Value(v.CapabilityID),
//				},
//	}
func (v *AccountCapabilityControllerValue) Encode(e *atree.Encoder) error {
	// Encode tag number and array head
	err := e.CBOR.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagAccountCapabilityControllerValue,
		// array, 2 items follow
		0x82,
	})
	if err != nil {
		return err
	}

	// Encode borrow type at array index encodedAccountCapabilityControllerValueBorrowTypeFieldKey
	err = v.BorrowType.Encode(e.CBOR)
	if err != nil {
		return err
	}

	// Encode ID at array index encodedAccountCapabilityControllerValueCapabilityIDFieldKey
	return e.CBOR.EncodeUint64(uint64(v.CapabilityID))
}

func StaticTypeToBytes(t StaticType) (cbor.RawMessage, error) {
	var buf bytes.Buffer
	enc := CBOREncMode.NewStreamEncoder(&buf)

	err := t.Encode(enc)
	if err != nil {
		return nil, err
	}

	err = enc.Flush()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Encode encodes PrimitiveStaticType as
//
//	cbor.Tag{
//			Number:  CBORTagPrimitiveStaticType,
//			Content: uint(v),
//	}
func (t PrimitiveStaticType) Encode(e *cbor.StreamEncoder) error {
	err := e.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagPrimitiveStaticType,
	})
	if err != nil {
		return err
	}
	return e.EncodeUint(uint(t))
}

// Encode encodes OptionalStaticType as
//
//	cbor.Tag{
//			Number:  CBORTagOptionalStaticType,
//			Content: StaticType(v.Type),
//	}
func (t *OptionalStaticType) Encode(e *cbor.StreamEncoder) error {
	err := e.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagOptionalStaticType,
	})
	if err != nil {
		return err
	}

	return t.Type.Encode(e)
}

// NOTE: NEVER change, only add/increment; ensure uint64
const (
	// encodedCompositeStaticTypeLocationFieldKey            uint64 = 0
	// encodedCompositeStaticTypeQualifiedIdentifierFieldKey uint64 = 1

	// !!! *WARNING* !!!
	//
	// encodedCompositeStaticTypeLength MUST be updated when new element is added.
	// It is used to verify encoded composite static type length during decoding.
	encodedCompositeStaticTypeLength = 2
)

// Encode encodes CompositeStaticType as
//
//	cbor.Tag{
//				Number: CBORTagCompositeStaticType,
//				Content: cborArray{
//					encodedCompositeStaticTypeLocationFieldKey:            Location(v.Location),
//					encodedCompositeStaticTypeQualifiedIdentifierFieldKey: string(v.QualifiedIdentifier),
//			},
//	}
func (t *CompositeStaticType) Encode(e *cbor.StreamEncoder) error {
	// Encode tag number and array head
	err := e.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagCompositeStaticType,
		// array, 2 items follow
		0x82,
	})
	if err != nil {
		return err
	}

	// Encode location at array index encodedCompositeStaticTypeLocationFieldKey
	err = encodeLocation(e, t.Location)
	if err != nil {
		return err
	}

	// Encode qualified identifier at array index encodedCompositeStaticTypeQualifiedIdentifierFieldKey
	return e.EncodeString(t.QualifiedIdentifier)
}

// NOTE: NEVER change, only add/increment; ensure uint64
const (
	// encodedInterfaceStaticTypeLocationFieldKey            uint64 = 0
	// encodedInterfaceStaticTypeQualifiedIdentifierFieldKey uint64 = 1

	// !!! *WARNING* !!!
	//
	// encodedInterfaceStaticTypeLength MUST be updated when new element is added.
	// It is used to verify encoded interface static type length during decoding.
	encodedInterfaceStaticTypeLength = 2
)

// Encode encodes InterfaceStaticType as
//
//	cbor.Tag{
//			Number: CBORTagInterfaceStaticType,
//			Content: cborArray{
//					encodedInterfaceStaticTypeLocationFieldKey:            Location(v.Location),
//					encodedInterfaceStaticTypeQualifiedIdentifierFieldKey: string(v.QualifiedIdentifier),
//			},
//	}
func (t *InterfaceStaticType) Encode(e *cbor.StreamEncoder) error {
	// Encode tag number and array head
	err := e.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagInterfaceStaticType,
		// array, 2 items follow
		0x82,
	})
	if err != nil {
		return err
	}

	// Encode location at array index encodedInterfaceStaticTypeLocationFieldKey
	err = encodeLocation(e, t.Location)
	if err != nil {
		return err
	}

	// Encode qualified identifier at array index encodedInterfaceStaticTypeQualifiedIdentifierFieldKey
	return e.EncodeString(t.QualifiedIdentifier)
}

// Encode encodes VariableSizedStaticType as
//
//	cbor.Tag{
//			Number:  CBORTagVariableSizedStaticType,
//			Content: StaticType(v.Type),
//	}
func (t *VariableSizedStaticType) Encode(e *cbor.StreamEncoder) error {
	err := e.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagVariableSizedStaticType,
	})
	if err != nil {
		return err
	}
	return t.Type.Encode(e)
}

// NOTE: NEVER change, only add/increment; ensure uint64
const (
	// encodedConstantSizedStaticTypeSizeFieldKey uint64 = 0
	// encodedConstantSizedStaticTypeTypeFieldKey uint64 = 1

	// !!! *WARNING* !!!
	//
	// encodedConstantSizedStaticTypeLength MUST be updated when new element is added.
	// It is used to verify encoded constant sized static type length during decoding.
	encodedConstantSizedStaticTypeLength = 2
)

// Encode encodes ConstantSizedStaticType as
//
//	cbor.Tag{
//			Number: CBORTagConstantSizedStaticType,
//			Content: cborArray{
//					encodedConstantSizedStaticTypeSizeFieldKey: int64(v.Size),
//					encodedConstantSizedStaticTypeTypeFieldKey: StaticType(v.Type),
//			},
//	}
func (t *ConstantSizedStaticType) Encode(e *cbor.StreamEncoder) error {
	// Encode tag number and array head
	err := e.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagConstantSizedStaticType,
		// array, 2 items follow
		0x82,
	})
	if err != nil {
		return err
	}
	// Encode size at array index encodedConstantSizedStaticTypeSizeFieldKey
	err = e.EncodeInt64(t.Size)
	if err != nil {
		return err
	}
	// Encode type at array index encodedConstantSizedStaticTypeTypeFieldKey
	return t.Type.Encode(e)
}

func (t Unauthorized) Encode(e *cbor.StreamEncoder) error {
	err := e.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagUnauthorizedStaticAuthorization,
	})
	if err != nil {
		return err
	}
	return e.EncodeNil()
}

func (a EntitlementMapAuthorization) Encode(e *cbor.StreamEncoder) error {
	err := e.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagEntitlementMapStaticAuthorization,
	})
	if err != nil {
		return err
	}
	return e.EncodeString(string(a.TypeID))
}

// NOTE: NEVER change, only add/increment; ensure uint64
const (
	// encodedSetAuthorizationStaticTypeKindKey                uint64 = 0
	// encodedSetAuthorizationStaticTypeEntitlementsKey        uint64 = 1

	// !!! *WARNING* !!!
	//
	// encodedSetAuthorizationStaticTypeLength MUST be updated when new element is added.
	// It is used to verify encoded reference static type length during decoding.
	encodedSetAuthorizationStaticTypeLength = 2
)

func (a EntitlementSetAuthorization) Encode(e *cbor.StreamEncoder) error {
	err := e.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagEntitlementSetStaticAuthorization,
		// array, 2 items follow
		0x82,
	})
	if err != nil {
		return err
	}

	err = e.EncodeUint8(uint8(a.SetKind))
	if err != nil {
		return err
	}

	err = e.EncodeArrayHead(uint64(a.Entitlements.Len()))
	if err != nil {
		return err
	}
	return a.Entitlements.ForeachWithError(func(entitlement common.TypeID, value struct{}) error {
		// Encode entitlement as array entitlements element
		return e.EncodeString(string(entitlement))
	})
}

// NOTE: NEVER change, only add/increment; ensure uint64
const (
	// encodedReferenceStaticTypeAuthorizationFieldKey uint64 = 0
	// encodedReferenceStaticTypeTypeFieldKey          uint64 = 1

	// !!! *WARNING* !!!
	//
	// encodedReferenceStaticTypeLength MUST be updated when new element is added.
	// It is used to verify encoded reference static type length during decoding.
	encodedReferenceStaticTypeLength = 2
)

// Encode encodes ReferenceStaticType as
//
//	cbor.Tag{
//			Number: CBORTagReferenceStaticType,
//			Content: cborArray{
//					encodedReferenceStaticTypeAuthorizationFieldKey: v.Authorization,
//					encodedReferenceStaticTypeTypeFieldKey:          StaticType(v.Type),
//			},
//		}
func (t *ReferenceStaticType) Encode(e *cbor.StreamEncoder) error {
	// Encode tag number and array head
	err := e.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagReferenceStaticType,
		// array, 2 items follow
		0x82,
	})
	if err != nil {
		return err
	}
	// Encode authorized at array index encodedReferenceStaticTypeAuthorizationFieldKey
	err = t.Authorization.Encode(e)
	if err != nil {
		return err
	}
	// Encode type at array index encodedReferenceStaticTypeTypeFieldKey
	return t.ReferencedType.Encode(e)
}

// NOTE: NEVER change, only add/increment; ensure uint64
const (
	// encodedDictionaryStaticTypeKeyTypeFieldKey   uint64 = 0
	// encodedDictionaryStaticTypeValueTypeFieldKey uint64 = 1

	// !!! *WARNING* !!!
	//
	// encodedDictionaryStaticTypeLength MUST be updated when new element is added.
	// It is used to verify encoded dictionary static type length during decoding.
	encodedDictionaryStaticTypeLength = 2
)

// Encode encodes DictionaryStaticType as
//
//	cbor.Tag{
//			Number: CBORTagDictionaryStaticType,
//			Content: []any{
//					encodedDictionaryStaticTypeKeyTypeFieldKey:   StaticType(v.KeyType),
//					encodedDictionaryStaticTypeValueTypeFieldKey: StaticType(v.ValueType),
//			},
//	}
func (t *DictionaryStaticType) Encode(e *cbor.StreamEncoder) error {
	// Encode tag number and array head
	err := e.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagDictionaryStaticType,
		// array, 2 items follow
		0x82,
	})
	if err != nil {
		return err
	}
	// Encode key type at array index encodedDictionaryStaticTypeKeyTypeFieldKey
	err = t.KeyType.Encode(e)
	if err != nil {
		return err
	}
	// Encode value type at array index encodedDictionaryStaticTypeValueTypeFieldKey
	return t.ValueType.Encode(e)
}

// Encode encodes InclusiveRangeStaticType as
//
//	cbor.Tag{
//			Number: CBORTagInclusiveRangeStaticType,
//			Content: StaticType(v.Type),
//	}
func (t InclusiveRangeStaticType) Encode(e *cbor.StreamEncoder) error {
	// Encode tag number and array head
	err := e.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagInclusiveRangeStaticType,
	})
	if err != nil {
		return err
	}

	return t.ElementType.Encode(e)
}

// NOTE: NEVER change, only add/increment; ensure uint64
const (
	// encodedIntersectionStaticTypeLegacyTypeFieldKey  uint64 = 0
	// encodedIntersectionStaticTypeTypesFieldKey       uint64 = 1

	// !!! *WARNING* !!!
	//
	// encodedIntersectionStaticTypeLength MUST be updated when new element is added.
	// It is used to verify encoded intersection static type length during decoding.
	encodedIntersectionStaticTypeLength = 2
)

// Encode encodes IntersectionStaticType as
//
//	cbor.Tag{
//			Number: CBORTagIntersectionStaticType,
//			Content: cborArray{
//					encodedIntersectionStaticTypeLegacyTypeFieldKey: StaticType(v.LegacyRestrictedType),
//					encodedIntersectionStaticTypeTypesFieldKey:		[]any(v.Types),
//			},
//	}
func (t *IntersectionStaticType) Encode(e *cbor.StreamEncoder) error {
	// Encode tag number and array head
	err := e.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagIntersectionStaticType,
		// array, 2 items follow
		0x82,
	})
	if err != nil {
		return err
	}

	if t.LegacyType != nil {
		// Encode type at array index encodedIntersectionStaticTypeTypeFieldKey
		err = t.LegacyType.Encode(e)
		if err != nil {
			return err
		}
	} else {
		err = e.EncodeNil()
		if err != nil {
			return err
		}
	}

	// Encode types (as array) at array index encodedIntersectionStaticTypeTypesFieldKey
	err = e.EncodeArrayHead(uint64(len(t.Types)))
	if err != nil {
		return err
	}

	for _, typ := range t.Types {
		// Encode typ as array types element
		err = typ.Encode(e)
		if err != nil {
			return err
		}
	}
	return nil
}

// Encode encodes CapabilityStaticType as
//
//	cbor.Tag{
//			Number:  CBORTagCapabilityStaticType,
//			Content: StaticType(v.BorrowType),
//	}
func (t *CapabilityStaticType) Encode(e *cbor.StreamEncoder) error {
	err := e.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagCapabilityStaticType,
	})
	if err != nil {
		return err
	}
	if t.BorrowType == nil {
		return e.EncodeNil()
	} else {
		return t.BorrowType.Encode(e)
	}
}

func (t FunctionStaticType) Encode(_ *cbor.StreamEncoder) error {
	return NonStorableStaticTypeError{
		Type: t.Type,
	}
}

// compositeTypeInfo
type compositeTypeInfo struct {
	location            common.Location
	qualifiedIdentifier string
	kind                common.CompositeKind
}

func NewCompositeTypeInfo(
	memoryGauge common.MemoryGauge,
	location common.Location,
	qualifiedIdentifier string,
	kind common.CompositeKind,
) compositeTypeInfo {
	common.UseMemory(memoryGauge, common.CompositeTypeInfoMemoryUsage)

	return compositeTypeInfo{
		location:            location,
		qualifiedIdentifier: qualifiedIdentifier,
		kind:                kind,
	}
}

var _ atree.TypeInfo = compositeTypeInfo{}

const encodedCompositeTypeInfoLength = 3

func (c compositeTypeInfo) Encode(e *cbor.StreamEncoder) error {
	err := e.EncodeRawBytes([]byte{
		// tag number
		0xd8, CBORTagCompositeValue,
		// array, 3 items follow
		0x83,
	})
	if err != nil {
		return err
	}

	err = encodeLocation(e, c.location)
	if err != nil {
		return err
	}

	err = e.EncodeString(c.qualifiedIdentifier)
	if err != nil {
		return err
	}

	err = e.EncodeUint64(uint64(c.kind))
	if err != nil {
		return err
	}

	return nil
}

func (c compositeTypeInfo) Equal(o atree.TypeInfo) bool {
	other, ok := o.(compositeTypeInfo)
	return ok &&
		c.location == other.location &&
		c.qualifiedIdentifier == other.qualifiedIdentifier &&
		c.kind == other.kind
}

// EmptyTypeInfo
type EmptyTypeInfo struct{}

func (e EmptyTypeInfo) Encode(encoder *cbor.StreamEncoder) error {
	return encoder.EncodeNil()
}

var emptyTypeInfo atree.TypeInfo = EmptyTypeInfo{}
