// Code generated by "stringer -type=MemoryKind -trimprefix=MemoryKind"; DO NOT EDIT.

package common

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[MemoryKindUnknown-0]
	_ = x[MemoryKindAddressValue-1]
	_ = x[MemoryKindStringValue-2]
	_ = x[MemoryKindCharacterValue-3]
	_ = x[MemoryKindNumberValue-4]
	_ = x[MemoryKindArrayValueBase-5]
	_ = x[MemoryKindDictionaryValueBase-6]
	_ = x[MemoryKindCompositeValueBase-7]
	_ = x[MemoryKindSimpleCompositeValueBase-8]
	_ = x[MemoryKindOptionalValue-9]
	_ = x[MemoryKindTypeValue-10]
	_ = x[MemoryKindPathValue-11]
	_ = x[MemoryKindIDCapabilityValue-12]
	_ = x[MemoryKindPathCapabilityValue-13]
	_ = x[MemoryKindPathLinkValue-14]
	_ = x[MemoryKindAccountLinkValue-15]
	_ = x[MemoryKindStorageReferenceValue-16]
	_ = x[MemoryKindAccountReferenceValue-17]
	_ = x[MemoryKindEphemeralReferenceValue-18]
	_ = x[MemoryKindInterpretedFunctionValue-19]
	_ = x[MemoryKindHostFunctionValue-20]
	_ = x[MemoryKindBoundFunctionValue-21]
	_ = x[MemoryKindBigInt-22]
	_ = x[MemoryKindSimpleCompositeValue-23]
	_ = x[MemoryKindPublishedValue-24]
	_ = x[MemoryKindStorageCapabilityControllerValue-25]
	_ = x[MemoryKindAccountCapabilityControllerValue-26]
	_ = x[MemoryKindAtreeArrayDataSlab-27]
	_ = x[MemoryKindAtreeArrayMetaDataSlab-28]
	_ = x[MemoryKindAtreeArrayElementOverhead-29]
	_ = x[MemoryKindAtreeMapDataSlab-30]
	_ = x[MemoryKindAtreeMapMetaDataSlab-31]
	_ = x[MemoryKindAtreeMapElementOverhead-32]
	_ = x[MemoryKindAtreeMapPreAllocatedElement-33]
	_ = x[MemoryKindAtreeEncodedSlab-34]
	_ = x[MemoryKindPrimitiveStaticType-35]
	_ = x[MemoryKindCompositeStaticType-36]
	_ = x[MemoryKindInterfaceStaticType-37]
	_ = x[MemoryKindVariableSizedStaticType-38]
	_ = x[MemoryKindConstantSizedStaticType-39]
	_ = x[MemoryKindDictionaryStaticType-40]
	_ = x[MemoryKindInclusiveRangeStaticType-41]
	_ = x[MemoryKindOptionalStaticType-42]
	_ = x[MemoryKindRestrictedStaticType-43]
	_ = x[MemoryKindReferenceStaticType-44]
	_ = x[MemoryKindCapabilityStaticType-45]
	_ = x[MemoryKindFunctionStaticType-46]
	_ = x[MemoryKindCadenceVoidValue-47]
	_ = x[MemoryKindCadenceOptionalValue-48]
	_ = x[MemoryKindCadenceBoolValue-49]
	_ = x[MemoryKindCadenceStringValue-50]
	_ = x[MemoryKindCadenceCharacterValue-51]
	_ = x[MemoryKindCadenceAddressValue-52]
	_ = x[MemoryKindCadenceIntValue-53]
	_ = x[MemoryKindCadenceNumberValue-54]
	_ = x[MemoryKindCadenceArrayValueBase-55]
	_ = x[MemoryKindCadenceArrayValueLength-56]
	_ = x[MemoryKindCadenceDictionaryValue-57]
	_ = x[MemoryKindCadenceKeyValuePair-58]
	_ = x[MemoryKindCadenceStructValueBase-59]
	_ = x[MemoryKindCadenceStructValueSize-60]
	_ = x[MemoryKindCadenceResourceValueBase-61]
	_ = x[MemoryKindCadenceAttachmentValueBase-62]
	_ = x[MemoryKindCadenceResourceValueSize-63]
	_ = x[MemoryKindCadenceAttachmentValueSize-64]
	_ = x[MemoryKindCadenceEventValueBase-65]
	_ = x[MemoryKindCadenceEventValueSize-66]
	_ = x[MemoryKindCadenceContractValueBase-67]
	_ = x[MemoryKindCadenceContractValueSize-68]
	_ = x[MemoryKindCadenceEnumValueBase-69]
	_ = x[MemoryKindCadenceEnumValueSize-70]
	_ = x[MemoryKindCadencePathLinkValue-71]
	_ = x[MemoryKindCadenceAccountLinkValue-72]
	_ = x[MemoryKindCadencePathValue-73]
	_ = x[MemoryKindCadenceTypeValue-74]
	_ = x[MemoryKindCadenceIDCapabilityValue-75]
	_ = x[MemoryKindCadencePathCapabilityValue-76]
	_ = x[MemoryKindCadenceFunctionValue-77]
	_ = x[MemoryKindCadenceOptionalType-78]
	_ = x[MemoryKindCadenceVariableSizedArrayType-79]
	_ = x[MemoryKindCadenceConstantSizedArrayType-80]
	_ = x[MemoryKindCadenceDictionaryType-81]
	_ = x[MemoryKindCadenceInclusiveRangeType-82]
	_ = x[MemoryKindCadenceField-83]
	_ = x[MemoryKindCadenceParameter-84]
	_ = x[MemoryKindCadenceTypeParameter-85]
	_ = x[MemoryKindCadenceStructType-86]
	_ = x[MemoryKindCadenceResourceType-87]
	_ = x[MemoryKindCadenceAttachmentType-88]
	_ = x[MemoryKindCadenceEventType-89]
	_ = x[MemoryKindCadenceContractType-90]
	_ = x[MemoryKindCadenceStructInterfaceType-91]
	_ = x[MemoryKindCadenceResourceInterfaceType-92]
	_ = x[MemoryKindCadenceContractInterfaceType-93]
	_ = x[MemoryKindCadenceFunctionType-94]
	_ = x[MemoryKindCadenceReferenceType-95]
	_ = x[MemoryKindCadenceRestrictedType-96]
	_ = x[MemoryKindCadenceCapabilityType-97]
	_ = x[MemoryKindCadenceEnumType-98]
	_ = x[MemoryKindRawString-99]
	_ = x[MemoryKindAddressLocation-100]
	_ = x[MemoryKindBytes-101]
	_ = x[MemoryKindVariable-102]
	_ = x[MemoryKindCompositeTypeInfo-103]
	_ = x[MemoryKindCompositeField-104]
	_ = x[MemoryKindInvocation-105]
	_ = x[MemoryKindStorageMap-106]
	_ = x[MemoryKindStorageKey-107]
	_ = x[MemoryKindTypeToken-108]
	_ = x[MemoryKindErrorToken-109]
	_ = x[MemoryKindSpaceToken-110]
	_ = x[MemoryKindProgram-111]
	_ = x[MemoryKindIdentifier-112]
	_ = x[MemoryKindArgument-113]
	_ = x[MemoryKindBlock-114]
	_ = x[MemoryKindFunctionBlock-115]
	_ = x[MemoryKindParameter-116]
	_ = x[MemoryKindParameterList-117]
	_ = x[MemoryKindTypeParameter-118]
	_ = x[MemoryKindTypeParameterList-119]
	_ = x[MemoryKindTransfer-120]
	_ = x[MemoryKindMembers-121]
	_ = x[MemoryKindTypeAnnotation-122]
	_ = x[MemoryKindDictionaryEntry-123]
	_ = x[MemoryKindFunctionDeclaration-124]
	_ = x[MemoryKindCompositeDeclaration-125]
	_ = x[MemoryKindAttachmentDeclaration-126]
	_ = x[MemoryKindInterfaceDeclaration-127]
	_ = x[MemoryKindEnumCaseDeclaration-128]
	_ = x[MemoryKindFieldDeclaration-129]
	_ = x[MemoryKindTransactionDeclaration-130]
	_ = x[MemoryKindImportDeclaration-131]
	_ = x[MemoryKindVariableDeclaration-132]
	_ = x[MemoryKindSpecialFunctionDeclaration-133]
	_ = x[MemoryKindPragmaDeclaration-134]
	_ = x[MemoryKindAssignmentStatement-135]
	_ = x[MemoryKindBreakStatement-136]
	_ = x[MemoryKindContinueStatement-137]
	_ = x[MemoryKindEmitStatement-138]
	_ = x[MemoryKindExpressionStatement-139]
	_ = x[MemoryKindForStatement-140]
	_ = x[MemoryKindIfStatement-141]
	_ = x[MemoryKindReturnStatement-142]
	_ = x[MemoryKindSwapStatement-143]
	_ = x[MemoryKindSwitchStatement-144]
	_ = x[MemoryKindWhileStatement-145]
	_ = x[MemoryKindRemoveStatement-146]
	_ = x[MemoryKindBooleanExpression-147]
	_ = x[MemoryKindVoidExpression-148]
	_ = x[MemoryKindNilExpression-149]
	_ = x[MemoryKindStringExpression-150]
	_ = x[MemoryKindIntegerExpression-151]
	_ = x[MemoryKindFixedPointExpression-152]
	_ = x[MemoryKindArrayExpression-153]
	_ = x[MemoryKindDictionaryExpression-154]
	_ = x[MemoryKindIdentifierExpression-155]
	_ = x[MemoryKindInvocationExpression-156]
	_ = x[MemoryKindMemberExpression-157]
	_ = x[MemoryKindIndexExpression-158]
	_ = x[MemoryKindConditionalExpression-159]
	_ = x[MemoryKindUnaryExpression-160]
	_ = x[MemoryKindBinaryExpression-161]
	_ = x[MemoryKindFunctionExpression-162]
	_ = x[MemoryKindCastingExpression-163]
	_ = x[MemoryKindCreateExpression-164]
	_ = x[MemoryKindDestroyExpression-165]
	_ = x[MemoryKindReferenceExpression-166]
	_ = x[MemoryKindForceExpression-167]
	_ = x[MemoryKindPathExpression-168]
	_ = x[MemoryKindAttachExpression-169]
	_ = x[MemoryKindConstantSizedType-170]
	_ = x[MemoryKindDictionaryType-171]
	_ = x[MemoryKindFunctionType-172]
	_ = x[MemoryKindInstantiationType-173]
	_ = x[MemoryKindNominalType-174]
	_ = x[MemoryKindOptionalType-175]
	_ = x[MemoryKindReferenceType-176]
	_ = x[MemoryKindRestrictedType-177]
	_ = x[MemoryKindVariableSizedType-178]
	_ = x[MemoryKindPosition-179]
	_ = x[MemoryKindRange-180]
	_ = x[MemoryKindElaboration-181]
	_ = x[MemoryKindActivation-182]
	_ = x[MemoryKindActivationEntries-183]
	_ = x[MemoryKindVariableSizedSemaType-184]
	_ = x[MemoryKindConstantSizedSemaType-185]
	_ = x[MemoryKindDictionarySemaType-186]
	_ = x[MemoryKindOptionalSemaType-187]
	_ = x[MemoryKindRestrictedSemaType-188]
	_ = x[MemoryKindReferenceSemaType-189]
	_ = x[MemoryKindCapabilitySemaType-190]
	_ = x[MemoryKindOrderedMap-191]
	_ = x[MemoryKindOrderedMapEntryList-192]
	_ = x[MemoryKindOrderedMapEntry-193]
	_ = x[MemoryKindLast-194]
}

const _MemoryKind_name = "UnknownAddressValueStringValueCharacterValueNumberValueArrayValueBaseDictionaryValueBaseCompositeValueBaseSimpleCompositeValueBaseOptionalValueTypeValuePathValueIDCapabilityValuePathCapabilityValuePathLinkValueAccountLinkValueStorageReferenceValueAccountReferenceValueEphemeralReferenceValueInterpretedFunctionValueHostFunctionValueBoundFunctionValueBigIntSimpleCompositeValuePublishedValueStorageCapabilityControllerValueAccountCapabilityControllerValueAtreeArrayDataSlabAtreeArrayMetaDataSlabAtreeArrayElementOverheadAtreeMapDataSlabAtreeMapMetaDataSlabAtreeMapElementOverheadAtreeMapPreAllocatedElementAtreeEncodedSlabPrimitiveStaticTypeCompositeStaticTypeInterfaceStaticTypeVariableSizedStaticTypeConstantSizedStaticTypeDictionaryStaticTypeInclusiveRangeStaticTypeOptionalStaticTypeRestrictedStaticTypeReferenceStaticTypeCapabilityStaticTypeFunctionStaticTypeCadenceVoidValueCadenceOptionalValueCadenceBoolValueCadenceStringValueCadenceCharacterValueCadenceAddressValueCadenceIntValueCadenceNumberValueCadenceArrayValueBaseCadenceArrayValueLengthCadenceDictionaryValueCadenceKeyValuePairCadenceStructValueBaseCadenceStructValueSizeCadenceResourceValueBaseCadenceAttachmentValueBaseCadenceResourceValueSizeCadenceAttachmentValueSizeCadenceEventValueBaseCadenceEventValueSizeCadenceContractValueBaseCadenceContractValueSizeCadenceEnumValueBaseCadenceEnumValueSizeCadencePathLinkValueCadenceAccountLinkValueCadencePathValueCadenceTypeValueCadenceIDCapabilityValueCadencePathCapabilityValueCadenceFunctionValueCadenceOptionalTypeCadenceVariableSizedArrayTypeCadenceConstantSizedArrayTypeCadenceDictionaryTypeCadenceInclusiveRangeTypeCadenceFieldCadenceParameterCadenceTypeParameterCadenceStructTypeCadenceResourceTypeCadenceAttachmentTypeCadenceEventTypeCadenceContractTypeCadenceStructInterfaceTypeCadenceResourceInterfaceTypeCadenceContractInterfaceTypeCadenceFunctionTypeCadenceReferenceTypeCadenceRestrictedTypeCadenceCapabilityTypeCadenceEnumTypeRawStringAddressLocationBytesVariableCompositeTypeInfoCompositeFieldInvocationStorageMapStorageKeyTypeTokenErrorTokenSpaceTokenProgramIdentifierArgumentBlockFunctionBlockParameterParameterListTypeParameterTypeParameterListTransferMembersTypeAnnotationDictionaryEntryFunctionDeclarationCompositeDeclarationAttachmentDeclarationInterfaceDeclarationEnumCaseDeclarationFieldDeclarationTransactionDeclarationImportDeclarationVariableDeclarationSpecialFunctionDeclarationPragmaDeclarationAssignmentStatementBreakStatementContinueStatementEmitStatementExpressionStatementForStatementIfStatementReturnStatementSwapStatementSwitchStatementWhileStatementRemoveStatementBooleanExpressionVoidExpressionNilExpressionStringExpressionIntegerExpressionFixedPointExpressionArrayExpressionDictionaryExpressionIdentifierExpressionInvocationExpressionMemberExpressionIndexExpressionConditionalExpressionUnaryExpressionBinaryExpressionFunctionExpressionCastingExpressionCreateExpressionDestroyExpressionReferenceExpressionForceExpressionPathExpressionAttachExpressionConstantSizedTypeDictionaryTypeFunctionTypeInstantiationTypeNominalTypeOptionalTypeReferenceTypeRestrictedTypeVariableSizedTypePositionRangeElaborationActivationActivationEntriesVariableSizedSemaTypeConstantSizedSemaTypeDictionarySemaTypeOptionalSemaTypeRestrictedSemaTypeReferenceSemaTypeCapabilitySemaTypeOrderedMapOrderedMapEntryListOrderedMapEntryLast"

var _MemoryKind_index = [...]uint16{0, 7, 19, 30, 44, 55, 69, 88, 106, 130, 143, 152, 161, 178, 197, 210, 226, 247, 268, 291, 315, 332, 350, 356, 376, 390, 422, 454, 472, 494, 519, 535, 555, 578, 605, 621, 640, 659, 678, 701, 724, 744, 768, 786, 806, 825, 845, 863, 879, 899, 915, 933, 954, 973, 988, 1006, 1027, 1050, 1072, 1091, 1113, 1135, 1159, 1185, 1209, 1235, 1256, 1277, 1301, 1325, 1345, 1365, 1385, 1408, 1424, 1440, 1464, 1490, 1510, 1529, 1558, 1587, 1608, 1633, 1645, 1661, 1681, 1698, 1717, 1738, 1754, 1773, 1799, 1827, 1855, 1874, 1894, 1915, 1936, 1951, 1960, 1975, 1980, 1988, 2005, 2019, 2029, 2039, 2049, 2058, 2068, 2078, 2085, 2095, 2103, 2108, 2121, 2130, 2143, 2156, 2173, 2181, 2188, 2202, 2217, 2236, 2256, 2277, 2297, 2316, 2332, 2354, 2371, 2390, 2416, 2433, 2452, 2466, 2483, 2496, 2515, 2527, 2538, 2553, 2566, 2581, 2595, 2610, 2627, 2641, 2654, 2670, 2687, 2707, 2722, 2742, 2762, 2782, 2798, 2813, 2834, 2849, 2865, 2883, 2900, 2916, 2933, 2952, 2967, 2981, 2997, 3014, 3028, 3040, 3057, 3068, 3080, 3093, 3107, 3124, 3132, 3137, 3148, 3158, 3175, 3196, 3217, 3235, 3251, 3269, 3286, 3304, 3314, 3333, 3348, 3352}

func (i MemoryKind) String() string {
	if i >= MemoryKind(len(_MemoryKind_index)-1) {
		return "MemoryKind(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _MemoryKind_name[_MemoryKind_index[i]:_MemoryKind_index[i+1]]
}
