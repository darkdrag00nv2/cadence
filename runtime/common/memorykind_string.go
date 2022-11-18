// Code generated by "stringer -type=MemoryKind -trimprefix=MemoryKind"; DO NOT EDIT.

package common

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[MemoryKindUnknown-0]
	_ = x[MemoryKindBoolValue-1]
	_ = x[MemoryKindAddressValue-2]
	_ = x[MemoryKindStringValue-3]
	_ = x[MemoryKindCharacterValue-4]
	_ = x[MemoryKindNumberValue-5]
	_ = x[MemoryKindArrayValueBase-6]
	_ = x[MemoryKindDictionaryValueBase-7]
	_ = x[MemoryKindCompositeValueBase-8]
	_ = x[MemoryKindSimpleCompositeValueBase-9]
	_ = x[MemoryKindOptionalValue-10]
	_ = x[MemoryKindNilValue-11]
	_ = x[MemoryKindVoidValue-12]
	_ = x[MemoryKindTypeValue-13]
	_ = x[MemoryKindPathValue-14]
	_ = x[MemoryKindStorageCapabilityValue-15]
	_ = x[MemoryKindLinkValue-16]
	_ = x[MemoryKindStorageReferenceValue-17]
	_ = x[MemoryKindEphemeralReferenceValue-18]
	_ = x[MemoryKindInterpretedFunctionValue-19]
	_ = x[MemoryKindHostFunctionValue-20]
	_ = x[MemoryKindBoundFunctionValue-21]
	_ = x[MemoryKindBigInt-22]
	_ = x[MemoryKindSimpleCompositeValue-23]
	_ = x[MemoryKindPublishedValue-24]
	_ = x[MemoryKindAtreeArrayDataSlab-25]
	_ = x[MemoryKindAtreeArrayMetaDataSlab-26]
	_ = x[MemoryKindAtreeArrayElementOverhead-27]
	_ = x[MemoryKindAtreeMapDataSlab-28]
	_ = x[MemoryKindAtreeMapMetaDataSlab-29]
	_ = x[MemoryKindAtreeMapElementOverhead-30]
	_ = x[MemoryKindAtreeMapPreAllocatedElement-31]
	_ = x[MemoryKindAtreeEncodedSlab-32]
	_ = x[MemoryKindPrimitiveStaticType-33]
	_ = x[MemoryKindCompositeStaticType-34]
	_ = x[MemoryKindInterfaceStaticType-35]
	_ = x[MemoryKindVariableSizedStaticType-36]
	_ = x[MemoryKindConstantSizedStaticType-37]
	_ = x[MemoryKindDictionaryStaticType-38]
	_ = x[MemoryKindOptionalStaticType-39]
	_ = x[MemoryKindRestrictedStaticType-40]
	_ = x[MemoryKindReferenceStaticType-41]
	_ = x[MemoryKindCapabilityStaticType-42]
	_ = x[MemoryKindFunctionStaticType-43]
	_ = x[MemoryKindCadenceVoidValue-44]
	_ = x[MemoryKindCadenceOptionalValue-45]
	_ = x[MemoryKindCadenceBoolValue-46]
	_ = x[MemoryKindCadenceStringValue-47]
	_ = x[MemoryKindCadenceCharacterValue-48]
	_ = x[MemoryKindCadenceAddressValue-49]
	_ = x[MemoryKindCadenceIntValue-50]
	_ = x[MemoryKindCadenceNumberValue-51]
	_ = x[MemoryKindCadenceArrayValueBase-52]
	_ = x[MemoryKindCadenceArrayValueLength-53]
	_ = x[MemoryKindCadenceDictionaryValue-54]
	_ = x[MemoryKindCadenceKeyValuePair-55]
	_ = x[MemoryKindCadenceStructValueBase-56]
	_ = x[MemoryKindCadenceStructValueSize-57]
	_ = x[MemoryKindCadenceResourceValueBase-58]
	_ = x[MemoryKindCadenceResourceValueSize-59]
	_ = x[MemoryKindCadenceEventValueBase-60]
	_ = x[MemoryKindCadenceEventValueSize-61]
	_ = x[MemoryKindCadenceContractValueBase-62]
	_ = x[MemoryKindCadenceContractValueSize-63]
	_ = x[MemoryKindCadenceEnumValueBase-64]
	_ = x[MemoryKindCadenceEnumValueSize-65]
	_ = x[MemoryKindCadenceLinkValue-66]
	_ = x[MemoryKindCadencePathValue-67]
	_ = x[MemoryKindCadenceTypeValue-68]
	_ = x[MemoryKindCadenceStorageCapabilityValue-69]
	_ = x[MemoryKindCadenceFunctionValue-70]
	_ = x[MemoryKindCadenceSimpleType-71]
	_ = x[MemoryKindCadenceOptionalType-72]
	_ = x[MemoryKindCadenceVariableSizedArrayType-73]
	_ = x[MemoryKindCadenceConstantSizedArrayType-74]
	_ = x[MemoryKindCadenceDictionaryType-75]
	_ = x[MemoryKindCadenceField-76]
	_ = x[MemoryKindCadenceParameter-77]
	_ = x[MemoryKindCadenceStructType-78]
	_ = x[MemoryKindCadenceResourceType-79]
	_ = x[MemoryKindCadenceEventType-80]
	_ = x[MemoryKindCadenceContractType-81]
	_ = x[MemoryKindCadenceStructInterfaceType-82]
	_ = x[MemoryKindCadenceResourceInterfaceType-83]
	_ = x[MemoryKindCadenceContractInterfaceType-84]
	_ = x[MemoryKindCadenceFunctionType-85]
	_ = x[MemoryKindCadenceReferenceType-86]
	_ = x[MemoryKindCadenceRestrictedType-87]
	_ = x[MemoryKindCadenceCapabilityType-88]
	_ = x[MemoryKindCadenceEnumType-89]
	_ = x[MemoryKindRawString-90]
	_ = x[MemoryKindAddressLocation-91]
	_ = x[MemoryKindBytes-92]
	_ = x[MemoryKindVariable-93]
	_ = x[MemoryKindCompositeTypeInfo-94]
	_ = x[MemoryKindCompositeField-95]
	_ = x[MemoryKindInvocation-96]
	_ = x[MemoryKindStorageMap-97]
	_ = x[MemoryKindStorageKey-98]
	_ = x[MemoryKindTypeToken-99]
	_ = x[MemoryKindErrorToken-100]
	_ = x[MemoryKindSpaceToken-101]
	_ = x[MemoryKindProgram-102]
	_ = x[MemoryKindIdentifier-103]
	_ = x[MemoryKindArgument-104]
	_ = x[MemoryKindBlock-105]
	_ = x[MemoryKindFunctionBlock-106]
	_ = x[MemoryKindParameter-107]
	_ = x[MemoryKindParameterList-108]
	_ = x[MemoryKindTransfer-109]
	_ = x[MemoryKindMembers-110]
	_ = x[MemoryKindTypeAnnotation-111]
	_ = x[MemoryKindDictionaryEntry-112]
	_ = x[MemoryKindFunctionDeclaration-113]
	_ = x[MemoryKindCompositeDeclaration-114]
	_ = x[MemoryKindInterfaceDeclaration-115]
	_ = x[MemoryKindEnumCaseDeclaration-116]
	_ = x[MemoryKindFieldDeclaration-117]
	_ = x[MemoryKindTransactionDeclaration-118]
	_ = x[MemoryKindImportDeclaration-119]
	_ = x[MemoryKindVariableDeclaration-120]
	_ = x[MemoryKindSpecialFunctionDeclaration-121]
	_ = x[MemoryKindPragmaDeclaration-122]
	_ = x[MemoryKindAssignmentStatement-123]
	_ = x[MemoryKindBreakStatement-124]
	_ = x[MemoryKindContinueStatement-125]
	_ = x[MemoryKindEmitStatement-126]
	_ = x[MemoryKindExpressionStatement-127]
	_ = x[MemoryKindForStatement-128]
	_ = x[MemoryKindIfStatement-129]
	_ = x[MemoryKindReturnStatement-130]
	_ = x[MemoryKindSwapStatement-131]
	_ = x[MemoryKindSwitchStatement-132]
	_ = x[MemoryKindWhileStatement-133]
	_ = x[MemoryKindBooleanExpression-134]
	_ = x[MemoryKindNilExpression-135]
	_ = x[MemoryKindStringExpression-136]
	_ = x[MemoryKindIntegerExpression-137]
	_ = x[MemoryKindFixedPointExpression-138]
	_ = x[MemoryKindArrayExpression-139]
	_ = x[MemoryKindDictionaryExpression-140]
	_ = x[MemoryKindIdentifierExpression-141]
	_ = x[MemoryKindInvocationExpression-142]
	_ = x[MemoryKindMemberExpression-143]
	_ = x[MemoryKindIndexExpression-144]
	_ = x[MemoryKindConditionalExpression-145]
	_ = x[MemoryKindUnaryExpression-146]
	_ = x[MemoryKindBinaryExpression-147]
	_ = x[MemoryKindFunctionExpression-148]
	_ = x[MemoryKindCastingExpression-149]
	_ = x[MemoryKindCreateExpression-150]
	_ = x[MemoryKindDestroyExpression-151]
	_ = x[MemoryKindReferenceExpression-152]
	_ = x[MemoryKindForceExpression-153]
	_ = x[MemoryKindPathExpression-154]
	_ = x[MemoryKindConstantSizedType-155]
	_ = x[MemoryKindDictionaryType-156]
	_ = x[MemoryKindFunctionType-157]
	_ = x[MemoryKindInstantiationType-158]
	_ = x[MemoryKindNominalType-159]
	_ = x[MemoryKindOptionalType-160]
	_ = x[MemoryKindReferenceType-161]
	_ = x[MemoryKindRestrictedType-162]
	_ = x[MemoryKindVariableSizedType-163]
	_ = x[MemoryKindPosition-164]
	_ = x[MemoryKindRange-165]
	_ = x[MemoryKindElaboration-166]
	_ = x[MemoryKindActivation-167]
	_ = x[MemoryKindActivationEntries-168]
	_ = x[MemoryKindVariableSizedSemaType-169]
	_ = x[MemoryKindConstantSizedSemaType-170]
	_ = x[MemoryKindDictionarySemaType-171]
	_ = x[MemoryKindOptionalSemaType-172]
	_ = x[MemoryKindRestrictedSemaType-173]
	_ = x[MemoryKindReferenceSemaType-174]
	_ = x[MemoryKindCapabilitySemaType-175]
	_ = x[MemoryKindOrderedMap-176]
	_ = x[MemoryKindOrderedMapEntryList-177]
	_ = x[MemoryKindOrderedMapEntry-178]
	_ = x[MemoryKindLast-179]
}

const _MemoryKind_name = "UnknownBoolValueAddressValueStringValueCharacterValueNumberValueArrayValueBaseDictionaryValueBaseCompositeValueBaseSimpleCompositeValueBaseOptionalValueNilValueVoidValueTypeValuePathValueStorageCapabilityValueLinkValueStorageReferenceValueEphemeralReferenceValueInterpretedFunctionValueHostFunctionValueBoundFunctionValueBigIntSimpleCompositeValuePublishedValueAtreeArrayDataSlabAtreeArrayMetaDataSlabAtreeArrayElementOverheadAtreeMapDataSlabAtreeMapMetaDataSlabAtreeMapElementOverheadAtreeMapPreAllocatedElementAtreeEncodedSlabPrimitiveStaticTypeCompositeStaticTypeInterfaceStaticTypeVariableSizedStaticTypeConstantSizedStaticTypeDictionaryStaticTypeOptionalStaticTypeRestrictedStaticTypeReferenceStaticTypeCapabilityStaticTypeFunctionStaticTypeCadenceVoidValueCadenceOptionalValueCadenceBoolValueCadenceStringValueCadenceCharacterValueCadenceAddressValueCadenceIntValueCadenceNumberValueCadenceArrayValueBaseCadenceArrayValueLengthCadenceDictionaryValueCadenceKeyValuePairCadenceStructValueBaseCadenceStructValueSizeCadenceResourceValueBaseCadenceResourceValueSizeCadenceEventValueBaseCadenceEventValueSizeCadenceContractValueBaseCadenceContractValueSizeCadenceEnumValueBaseCadenceEnumValueSizeCadenceLinkValueCadencePathValueCadenceTypeValueCadenceStorageCapabilityValueCadenceFunctionValueCadenceSimpleTypeCadenceOptionalTypeCadenceVariableSizedArrayTypeCadenceConstantSizedArrayTypeCadenceDictionaryTypeCadenceFieldCadenceParameterCadenceStructTypeCadenceResourceTypeCadenceEventTypeCadenceContractTypeCadenceStructInterfaceTypeCadenceResourceInterfaceTypeCadenceContractInterfaceTypeCadenceFunctionTypeCadenceReferenceTypeCadenceRestrictedTypeCadenceCapabilityTypeCadenceEnumTypeRawStringAddressLocationBytesVariableCompositeTypeInfoCompositeFieldInvocationStorageMapStorageKeyTypeTokenErrorTokenSpaceTokenProgramIdentifierArgumentBlockFunctionBlockParameterParameterListTransferMembersTypeAnnotationDictionaryEntryFunctionDeclarationCompositeDeclarationInterfaceDeclarationEnumCaseDeclarationFieldDeclarationTransactionDeclarationImportDeclarationVariableDeclarationSpecialFunctionDeclarationPragmaDeclarationAssignmentStatementBreakStatementContinueStatementEmitStatementExpressionStatementForStatementIfStatementReturnStatementSwapStatementSwitchStatementWhileStatementBooleanExpressionNilExpressionStringExpressionIntegerExpressionFixedPointExpressionArrayExpressionDictionaryExpressionIdentifierExpressionInvocationExpressionMemberExpressionIndexExpressionConditionalExpressionUnaryExpressionBinaryExpressionFunctionExpressionCastingExpressionCreateExpressionDestroyExpressionReferenceExpressionForceExpressionPathExpressionConstantSizedTypeDictionaryTypeFunctionTypeInstantiationTypeNominalTypeOptionalTypeReferenceTypeRestrictedTypeVariableSizedTypePositionRangeElaborationActivationActivationEntriesVariableSizedSemaTypeConstantSizedSemaTypeDictionarySemaTypeOptionalSemaTypeRestrictedSemaTypeReferenceSemaTypeCapabilitySemaTypeOrderedMapOrderedMapEntryListOrderedMapEntryLast"

var _MemoryKind_index = [...]uint16{0, 7, 16, 28, 39, 53, 64, 78, 97, 115, 139, 152, 160, 169, 178, 187, 209, 218, 239, 262, 286, 303, 321, 327, 347, 361, 379, 401, 426, 442, 462, 485, 512, 528, 547, 566, 585, 608, 631, 651, 669, 689, 708, 728, 746, 762, 782, 798, 816, 837, 856, 871, 889, 910, 933, 955, 974, 996, 1018, 1042, 1066, 1087, 1108, 1132, 1156, 1176, 1196, 1212, 1228, 1244, 1273, 1293, 1310, 1329, 1358, 1387, 1408, 1420, 1436, 1453, 1472, 1488, 1507, 1533, 1561, 1589, 1608, 1628, 1649, 1670, 1685, 1694, 1709, 1714, 1722, 1739, 1753, 1763, 1773, 1783, 1792, 1802, 1812, 1819, 1829, 1837, 1842, 1855, 1864, 1877, 1885, 1892, 1906, 1921, 1940, 1960, 1980, 1999, 2015, 2037, 2054, 2073, 2099, 2116, 2135, 2149, 2166, 2179, 2198, 2210, 2221, 2236, 2249, 2264, 2278, 2295, 2308, 2324, 2341, 2361, 2376, 2396, 2416, 2436, 2452, 2467, 2488, 2503, 2519, 2537, 2554, 2570, 2587, 2606, 2621, 2635, 2652, 2666, 2678, 2695, 2706, 2718, 2731, 2745, 2762, 2770, 2775, 2786, 2796, 2813, 2834, 2855, 2873, 2889, 2907, 2924, 2942, 2952, 2971, 2986, 2990}

func (i MemoryKind) String() string {
	if i >= MemoryKind(len(_MemoryKind_index)-1) {
		return "MemoryKind(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _MemoryKind_name[_MemoryKind_index[i]:_MemoryKind_index[i+1]]
}
