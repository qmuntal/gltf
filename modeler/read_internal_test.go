package modeler

import (
	"errors"
	"io"
	"testing"

	"github.com/qmuntal/gltf"
)

func TestValidateZeroInitConversionCapacity(t *testing.T) {
	acr := &gltf.Accessor{ComponentType: gltf.ComponentUbyte, Type: gltf.AccessorScalar, Count: 1 << 28}
	if err := validateZeroInitConversion(acr, 4, 0); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("no caller storage: expected io.ErrShortBuffer, got %v", err)
	}
	if err := validateZeroInitConversion(acr, 4, acr.Count-1); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("partial caller storage: expected io.ErrShortBuffer, got %v", err)
	}
	if err := validateZeroInitConversion(acr, 4, acr.Count); err != nil {
		t.Fatalf("full caller storage: expected success, got %v", err)
	}

	weights := &gltf.Accessor{ComponentType: gltf.ComponentUbyte, Type: gltf.AccessorVec4, Count: 60_000_000, Normalized: true}
	if err := validateZeroInitConversion(weights, 16, 0); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("weights: expected io.ErrShortBuffer, got %v", err)
	}
	if err := validateZeroInitConversion(weights, 16, weights.Count); err != nil {
		t.Fatalf("weights with caller storage: expected success, got %v", err)
	}
}

func TestValidateZeroInitAccessorCountUnpadded(t *testing.T) {
	acr := &gltf.Accessor{ComponentType: gltf.ComponentUbyte, Type: gltf.AccessorMat2, Count: 140_000_000}
	if err := validateZeroInitAccessorCount(acr, nil); err != nil {
		t.Fatalf("140M MAT2 ubyte: expected success, got %v", err)
	}
	over := &gltf.Accessor{ComponentType: gltf.ComponentUbyte, Type: gltf.AccessorMat2, Count: 270_000_000}
	if err := validateZeroInitAccessorCount(over, nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("270M MAT2 ubyte: expected io.ErrShortBuffer, got %v", err)
	}
}

func TestValidateZeroInitAccessorCountCallerCapacity(t *testing.T) {
	acr := &gltf.Accessor{ComponentType: gltf.ComponentFloat, Type: gltf.AccessorVec3, Count: 90_000_000}
	caller := make([]byte, 0, acr.Count*12)
	if err := validateZeroInitAccessorCount(acr, caller); err != nil {
		t.Fatalf("caller capacity covers result: expected success, got %v", err)
	}
	if err := validateZeroInitAccessorCount(acr, nil); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("no caller storage: expected io.ErrShortBuffer, got %v", err)
	}
}

func TestValidateZeroInitConversionColorUbyteVec3(t *testing.T) {
	// Regression: ReadColor's only widening source is UBYTE/VEC3 (3
	// bytes/elem raw -> 4 bytes/elem [4]uint8 output). 360M elements are
	// 1.08 GB of raw data alone, within the raw-only cap on its own, but
	// combined with the widened output exceeds it.
	acr := &gltf.Accessor{ComponentType: gltf.ComponentUbyte, Type: gltf.AccessorVec3, Count: 360_000_000, Normalized: true}
	if err := validateZeroInitConversion(acr, 4, 0); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("no caller storage: expected io.ErrShortBuffer, got %v", err)
	}
	if err := validateZeroInitConversion(acr, 4, acr.Count); err != nil {
		t.Fatalf("full caller storage: expected success, got %v", err)
	}
}

func TestValidateZeroInitConversionIdentityReaders(t *testing.T) {
	// The non-converting float readers (ReadPosition, ReadNormal, ...)
	// pass their own element size as outSize: data plus a same-size
	// output must stay within the cap. For FLOAT/VEC3 that is 24
	// bytes/elem combined, so the boundary is 1 GiB / 24 = 44,739,242.
	acr := &gltf.Accessor{ComponentType: gltf.ComponentFloat, Type: gltf.AccessorVec3, Count: 44_739_243}
	if err := validateZeroInitConversion(acr, 12, 0); !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("combined over cap: expected io.ErrShortBuffer, got %v", err)
	}
	acr.Count = 44_739_242
	if err := validateZeroInitConversion(acr, 12, 0); err != nil {
		t.Fatalf("combined at cap: expected success, got %v", err)
	}
	if err := validateZeroInitConversion(acr, 12, acr.Count); err != nil {
		t.Fatalf("full caller storage: expected success, got %v", err)
	}
}
