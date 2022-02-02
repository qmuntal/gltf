// Code generated by "stringer -linecomment -type ComponentType,AccessorType,PrimitiveMode,AlphaMode,MagFilter,MinFilter,WrappingMode,Interpolation,Target,TRSProperty -output const_strings.go"; DO NOT EDIT.

package gltf

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ComponentFloat-0]
	_ = x[ComponentByte-1]
	_ = x[ComponentUbyte-2]
	_ = x[ComponentShort-3]
	_ = x[ComponentUshort-4]
	_ = x[ComponentUint-5]
}

const _ComponentType_name = "FLOATBYTEUNSIGNED_BYTESHORTUNSIGNED_SHORTUNSIGNED_INT"

var _ComponentType_index = [...]uint8{0, 5, 9, 22, 27, 41, 53}

func (i ComponentType) String() string {
	if i >= ComponentType(len(_ComponentType_index)-1) {
		return "ComponentType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ComponentType_name[_ComponentType_index[i]:_ComponentType_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[AccessorScalar-0]
	_ = x[AccessorVec2-1]
	_ = x[AccessorVec3-2]
	_ = x[AccessorVec4-3]
	_ = x[AccessorMat2-4]
	_ = x[AccessorMat3-5]
	_ = x[AccessorMat4-6]
}

const _AccessorType_name = "SCALARVEC2VEC3VEC4MAT2MAT3MAT4"

var _AccessorType_index = [...]uint8{0, 6, 10, 14, 18, 22, 26, 30}

func (i AccessorType) String() string {
	if i >= AccessorType(len(_AccessorType_index)-1) {
		return "AccessorType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _AccessorType_name[_AccessorType_index[i]:_AccessorType_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[PrimitiveTriangles-0]
	_ = x[PrimitivePoints-1]
	_ = x[PrimitiveLines-2]
	_ = x[PrimitiveLineLoop-3]
	_ = x[PrimitiveLineStrip-4]
	_ = x[PrimitiveTriangleStrip-5]
	_ = x[PrimitiveTriangleFan-6]
}

const _PrimitiveMode_name = "TRIANGLESPOINTSLINESLINE_LOOPLINE_STRIPTRIANGLE_STRIPTRIANGLE_FAN"

var _PrimitiveMode_index = [...]uint8{0, 9, 15, 20, 29, 39, 53, 65}

func (i PrimitiveMode) String() string {
	if i >= PrimitiveMode(len(_PrimitiveMode_index)-1) {
		return "PrimitiveMode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _PrimitiveMode_name[_PrimitiveMode_index[i]:_PrimitiveMode_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[AlphaOpaque-0]
	_ = x[AlphaMask-1]
	_ = x[AlphaBlend-2]
}

const _AlphaMode_name = "OPAQUEMASKBLEND"

var _AlphaMode_index = [...]uint8{0, 6, 10, 15}

func (i AlphaMode) String() string {
	if i >= AlphaMode(len(_AlphaMode_index)-1) {
		return "AlphaMode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _AlphaMode_name[_AlphaMode_index[i]:_AlphaMode_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[MagUndefined-0]
	_ = x[MagLinear-1]
	_ = x[MagNearest-2]
}

const _MagFilter_name = "UNDEFINEDLINEARNEAREST"

var _MagFilter_index = [...]uint8{0, 9, 15, 22}

func (i MagFilter) String() string {
	if i >= MagFilter(len(_MagFilter_index)-1) {
		return "MagFilter(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _MagFilter_name[_MagFilter_index[i]:_MagFilter_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[MinUndefined-0]
	_ = x[MinLinear-1]
	_ = x[MinNearestMipMapLinear-2]
	_ = x[MinNearest-3]
	_ = x[MinNearestMipMapNearest-4]
	_ = x[MinLinearMipMapNearest-5]
	_ = x[MinLinearMipMapLinear-6]
}

const _MinFilter_name = "UNDEFINEDLINEARNEAREST_MIPMAP_LINEARNEARESTNEAREST_MIPMAP_NEARESTLINEAR_MIPMAP_NEARESTLINEAR_MIPMAP_LINEAR"

var _MinFilter_index = [...]uint8{0, 9, 15, 36, 43, 65, 86, 106}

func (i MinFilter) String() string {
	if i >= MinFilter(len(_MinFilter_index)-1) {
		return "MinFilter(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _MinFilter_name[_MinFilter_index[i]:_MinFilter_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[WrapRepeat-0]
	_ = x[WrapClampToEdge-1]
	_ = x[WrapMirroredRepeat-2]
}

const _WrappingMode_name = "REPEATCLAMP_TO_EDGEMIRRORED_REPEAT"

var _WrappingMode_index = [...]uint8{0, 6, 19, 34}

func (i WrappingMode) String() string {
	if i >= WrappingMode(len(_WrappingMode_index)-1) {
		return "WrappingMode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _WrappingMode_name[_WrappingMode_index[i]:_WrappingMode_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[InterpolationLinear-0]
	_ = x[InterpolationStep-1]
	_ = x[InterpolationCubicSpline-2]
}

const _Interpolation_name = "LINEARSTEPCUBICSPLINE"

var _Interpolation_index = [...]uint8{0, 6, 10, 21}

func (i Interpolation) String() string {
	if i >= Interpolation(len(_Interpolation_index)-1) {
		return "Interpolation(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Interpolation_name[_Interpolation_index[i]:_Interpolation_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TargetNone-0]
	_ = x[TargetArrayBuffer-34962]
	_ = x[TargetElementArrayBuffer-34963]
}

const (
	_Target_name_0 = "NONE"
	_Target_name_1 = "ARRAY_BUFFERELEMENT_ARRAY_BUFFER"
)

var (
	_Target_index_1 = [...]uint8{0, 12, 32}
)

func (i Target) String() string {
	switch {
	case i == 0:
		return _Target_name_0
	case 34962 <= i && i <= 34963:
		i -= 34962
		return _Target_name_1[_Target_index_1[i]:_Target_index_1[i+1]]
	default:
		return "Target(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TRSTranslation-0]
	_ = x[TRSRotation-1]
	_ = x[TRSScale-2]
	_ = x[TRSWeights-3]
}

const _TRSProperty_name = "translationrotationscaleweights"

var _TRSProperty_index = [...]uint8{0, 11, 19, 24, 31}

func (i TRSProperty) String() string {
	if i >= TRSProperty(len(_TRSProperty_index)-1) {
		return "TRSProperty(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TRSProperty_name[_TRSProperty_index[i]:_TRSProperty_index[i+1]]
}
