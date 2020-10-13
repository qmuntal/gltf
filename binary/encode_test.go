package binary

import (
	"encoding/binary"
	"image/color"
	"math"
	"reflect"
	"testing"

	"github.com/qmuntal/gltf"
)

func buildBuffer1(n int, empty ...int) []byte {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte(i + 1)
	}
	for _, e := range empty {
		b[e] = 0
	}
	return b
}

func buildBuffer2(n int, empty ...int) []byte {
	b := make([]byte, 2*n)
	for i := 0; i < n; i++ {
		binary.LittleEndian.PutUint16(b[i*2:], uint16(i+1))
	}
	for _, e := range empty {
		b[e] = 0
	}
	return b
}

func buildBuffer3(n int) []byte {
	b := make([]byte, 4*n)
	for i := 0; i < n; i++ {
		binary.LittleEndian.PutUint32(b[i*4:], uint32(i+1))
	}
	return b
}

func buildBufferF(n int) []byte {
	b := make([]byte, 4*n)
	for i := 0; i < n; i++ {
		binary.LittleEndian.PutUint32(b[i*4:], math.Float32bits(float32(i+1)))
	}
	return b
}

func TestRead(t *testing.T) {
	type args struct {
		b    []byte
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"not suported", args{make([]byte, 2), 1}, nil, true},
		{"small", args{[]byte{0, 0}, []int8{1, 2, 3}}, []int8{0, 0}, true},
		{"empty", args{make([]byte, 0), []int8{}}, []int8{}, false},
		{"int8", args{buildBuffer1(4), make([]int8, 4)}, []int8{1, 2, 3, 4}, false},
		{"int8-FE", args{[]byte{0xFE}, make([]int8, 1)}, []int8{-2}, false},
		{"[2]int8", args{[]byte{1, 2, 0, 0, 3, 4, 0, 0}, make([][2]int8, 2)}, [][2]int8{{1, 2}, {3, 4}}, false},
		{"[3]int8", args{[]byte{1, 2, 3, 0, 4, 5, 6, 0}, make([][3]int8, 2)}, [][3]int8{{1, 2, 3}, {4, 5, 6}}, false},
		{"[4]int8", args{buildBuffer1(2 * 4), make([][4]int8, 2)}, [][4]int8{{1, 2, 3, 4}, {5, 6, 7, 8}}, false},
		{"[2][2]int8", args{buildBuffer1(16, 2, 3, 6, 7, 10, 11, 14, 15), make([][2][2]int8, 2)}, [][2][2]int8{
			{{1, 5}, {2, 6}},
			{{9, 13}, {10, 14}},
		}, false},
		{"[3][3]int8", args{buildBuffer1(24, 3, 7, 11, 15, 19, 23), make([][3][3]int8, 2)}, [][3][3]int8{
			{{1, 5, 9}, {2, 6, 10}, {3, 7, 11}},
			{{13, 17, 21}, {14, 18, 22}, {15, 19, 23}},
		}, false},
		{"[4][4]int8", args{buildBuffer1(2 * 4 * 4), make([][4][4]int8, 2)}, [][4][4]int8{
			{{1, 5, 9, 13}, {2, 6, 10, 14}, {3, 7, 11, 15}, {4, 8, 12, 16}},
			{{17, 21, 25, 29}, {18, 22, 26, 30}, {19, 23, 27, 31}, {20, 24, 28, 32}},
		}, false},
		{"uint8", args{buildBuffer1(4), make([]uint8, 4)}, []uint8{1, 2, 3, 4}, false},
		{"uint8-FE", args{[]byte{0xFE}, make([]uint8, 1)}, []uint8{254}, false},
		{"[2]uint8", args{[]byte{1, 2, 0, 0, 3, 4, 0, 0}, make([][2]uint8, 2)}, [][2]uint8{{1, 2}, {3, 4}}, false},
		{"[3]uint8", args{[]byte{1, 2, 3, 0, 4, 5, 6, 0}, make([][3]uint8, 2)}, [][3]uint8{{1, 2, 3}, {4, 5, 6}}, false},
		{"[4]uint8", args{buildBuffer1(2 * 4), make([][4]uint8, 2)}, [][4]uint8{{1, 2, 3, 4}, {5, 6, 7, 8}}, false},
		{"[2][2]uint8", args{buildBuffer1(16, 2, 3, 6, 7, 10, 11), make([][2][2]uint8, 2)}, [][2][2]uint8{
			{{1, 5}, {2, 6}},
			{{9, 13}, {10, 14}},
		}, false},
		{"[3][3]uint8", args{buildBuffer1(32, 3, 7, 11, 15, 19, 23), make([][3][3]uint8, 2)}, [][3][3]uint8{
			{{1, 5, 9}, {2, 6, 10}, {3, 7, 11}},
			{{13, 17, 21}, {14, 18, 22}, {15, 19, 23}},
		}, false},
		{"[4][4]uint8", args{buildBuffer1(32), make([][4][4]uint8, 2)}, [][4][4]uint8{
			{{1, 5, 9, 13}, {2, 6, 10, 14}, {3, 7, 11, 15}, {4, 8, 12, 16}},
			{{17, 21, 25, 29}, {18, 22, 26, 30}, {19, 23, 27, 31}, {20, 24, 28, 32}},
		}, false},
		{"int16", args{buildBuffer2(4), make([]int16, 4)}, []int16{1, 2, 3, 4}, false},
		{"int16-FE", args{[]byte{0xFE, 0xFF}, make([]int16, 1)}, []int16{-2}, false},
		{"[2]int16", args{buildBuffer2(2 * 2), make([][2]int16, 2)}, [][2]int16{{1, 2}, {3, 4}}, false},
		{"[3]int16", args{[]byte{1, 0, 2, 0, 3, 0, 0, 0, 4, 0, 5, 0, 6, 0, 0, 0}, make([][3]int16, 2)}, [][3]int16{{1, 2, 3}, {4, 5, 6}}, false},
		{"[4]int16", args{buildBuffer2(2 * 4), make([][4]int16, 2)}, [][4]int16{{1, 2, 3, 4}, {5, 6, 7, 8}}, false},
		{"[2][2]int16", args{buildBuffer2(2 * 2 * 2), make([][2][2]int16, 2)}, [][2][2]int16{
			{{1, 3}, {2, 4}},
			{{5, 7}, {6, 8}},
		}, false},
		{"[3][3]int16", args{buildBuffer2(32, 6, 7, 14, 15, 22, 23, 30, 31, 14, 15, 22, 23, 38, 39), make([][3][3]int16, 2)}, [][3][3]int16{
			{{1, 5, 9}, {2, 6, 10}, {3, 7, 11}},
			{{13, 17, 21}, {14, 18, 22}, {15, 19, 23}},
		}, false},
		{"[4][4]int16", args{buildBuffer2(2 * 4 * 4), make([][4][4]int16, 2)}, [][4][4]int16{
			{{1, 5, 9, 13}, {2, 6, 10, 14}, {3, 7, 11, 15}, {4, 8, 12, 16}},
			{{17, 21, 25, 29}, {18, 22, 26, 30}, {19, 23, 27, 31}, {20, 24, 28, 32}},
		}, false},
		{"uint16", args{buildBuffer2(4), make([]uint16, 4)}, []uint16{1, 2, 3, 4}, false},
		{"uint16-FE", args{[]byte{0xFE, 0xFF}, make([]uint16, 1)}, []uint16{65534}, false},
		{"[2]uint16", args{buildBuffer2(2 * 2), make([][2]uint16, 2)}, [][2]uint16{{1, 2}, {3, 4}}, false},
		{"[3]uint16", args{[]byte{1, 0, 2, 0, 3, 0, 0, 0, 4, 0, 5, 0, 6, 0, 0, 0}, make([][3]uint16, 2)}, [][3]uint16{{1, 2, 3}, {4, 5, 6}}, false},
		{"[4]uint16", args{buildBuffer2(2 * 4), make([][4]uint16, 2)}, [][4]uint16{{1, 2, 3, 4}, {5, 6, 7, 8}}, false},
		{"[2][2]uint16", args{buildBuffer2(2 * 2 * 2), make([][2][2]uint16, 2)}, [][2][2]uint16{
			{{1, 3}, {2, 4}},
			{{5, 7}, {6, 8}},
		}, false},
		{"[3][3]uint16", args{buildBuffer2(32, 6, 7, 14, 15, 22, 23, 30, 31, 14, 15, 22, 23, 38, 39), make([][3][3]uint16, 2)}, [][3][3]uint16{
			{{1, 5, 9}, {2, 6, 10}, {3, 7, 11}},
			{{13, 17, 21}, {14, 18, 22}, {15, 19, 23}},
		}, false},
		{"[4][4]uint16", args{buildBuffer2(2 * 4 * 4), make([][4][4]uint16, 2)}, [][4][4]uint16{
			{{1, 5, 9, 13}, {2, 6, 10, 14}, {3, 7, 11, 15}, {4, 8, 12, 16}},
			{{17, 21, 25, 29}, {18, 22, 26, 30}, {19, 23, 27, 31}, {20, 24, 28, 32}},
		}, false},
		{"uint32", args{buildBuffer3(4), make([]uint32, 4)}, []uint32{1, 2, 3, 4}, false},
		{"uint32-FE", args{[]byte{0xFE, 0xFF, 0xFF, 0xFF}, make([]uint32, 1)}, []uint32{4294967294}, false},
		{"[2]uint32", args{buildBuffer3(2 * 2 * 4), make([][2]uint32, 2)}, [][2]uint32{{1, 2}, {3, 4}}, false},
		{"[3]uint32", args{buildBuffer3(2 * 3 * 4), make([][3]uint32, 2)}, [][3]uint32{{1, 2, 3}, {4, 5, 6}}, false},
		{"[4]uint32", args{buildBuffer3(2 * 4 * 4), make([][4]uint32, 2)}, [][4]uint32{{1, 2, 3, 4}, {5, 6, 7, 8}}, false},
		{"[2][2]uint32", args{buildBuffer3(2 * 2 * 2 * 4), make([][2][2]uint32, 2)}, [][2][2]uint32{
			{{1, 3}, {2, 4}},
			{{5, 7}, {6, 8}},
		}, false},
		{"[3][3]uint32", args{buildBuffer3(64 * 2), make([][3][3]uint32, 2)}, [][3][3]uint32{
			{{1, 4, 7}, {2, 5, 8}, {3, 6, 9}},
			{{10, 13, 16}, {11, 14, 17}, {12, 15, 18}},
		}, false},
		{"[4][4]uint32", args{buildBuffer3(2 * 4 * 4 * 4), make([][4][4]uint32, 2)}, [][4][4]uint32{
			{{1, 5, 9, 13}, {2, 6, 10, 14}, {3, 7, 11, 15}, {4, 8, 12, 16}},
			{{17, 21, 25, 29}, {18, 22, 26, 30}, {19, 23, 27, 31}, {20, 24, 28, 32}},
		}, false},
		{"float32", args{buildBufferF(4), make([]float32, 4)}, []float32{1, 2, 3, 4}, false},
		{"float32-FE", args{[]byte{0x00, 0x00, 0x00, 0xC0}, make([]float32, 1)}, []float32{-2}, false},
		{"[2]float32", args{buildBufferF(2 * 2 * 4), make([][2]float32, 2)}, [][2]float32{{1, 2}, {3, 4}}, false},
		{"[3]float32", args{buildBufferF(2 * 3 * 4), make([][3]float32, 2)}, [][3]float32{{1, 2, 3}, {4, 5, 6}}, false},
		{"[4]float32", args{buildBufferF(2 * 4 * 4), make([][4]float32, 2)}, [][4]float32{{1, 2, 3, 4}, {5, 6, 7, 8}}, false},
		{"[2][2]float32", args{buildBufferF(2 * 2 * 2 * 4), make([][2][2]float32, 2)}, [][2][2]float32{
			{{1, 3}, {2, 4}},
			{{5, 7}, {6, 8}},
		}, false},
		{"[3][3]float32", args{buildBufferF(64 * 2), make([][3][3]float32, 2)}, [][3][3]float32{
			{{1, 4, 7}, {2, 5, 8}, {3, 6, 9}},
			{{10, 13, 16}, {11, 14, 17}, {12, 15, 18}},
		}, false},
		{"[4][4]float32", args{buildBufferF(2 * 4 * 4 * 4), make([][4][4]float32, 2)}, [][4][4]float32{
			{{1, 5, 9, 13}, {2, 6, 10, 14}, {3, 7, 11, 15}, {4, 8, 12, 16}},
			{{17, 21, 25, 29}, {18, 22, 26, 30}, {19, 23, 27, 31}, {20, 24, 28, 32}},
		}, false},
		{"color.RGBA", args{buildBuffer1(2 * 4), make([]color.RGBA, 2)}, []color.RGBA{{1, 2, 3, 4}, {5, 6, 7, 8}}, false},
		{"color.RGBA64", args{buildBuffer2(2 * 4), make([]color.RGBA64, 2)}, []color.RGBA64{{1, 2, 3, 4}, {5, 6, 7, 8}}, false},
		{"gltf.RGBA", args{buildBufferF(2 * 4), make([]gltf.RGBA, 2)}, []gltf.RGBA{{R: 1, G: 2, B: 3, A: 4}, {R: 5, G: 6, B: 7, A: 8}}, false},
		{"gltf.RGB", args{buildBufferF(2 * 3), make([]gltf.RGB, 2)}, []gltf.RGB{{R: 1, G: 2, B: 3}, {R: 4, G: 5, B: 6}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Read(tt.args.b, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.args.data, tt.want) {
				t.Errorf("Read() error = %v, want %v", tt.args.data, tt.want)
			}
		})
	}
}

func TestWrite(t *testing.T) {
	type args struct {
		n    int
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"not suported", args{2, 1}, nil, true},
		{"small", args{2, []int8{1, 2, 3}}, []byte{1, 2, 3}, true},
		{"empty", args{0, []int8{}}, []byte{}, false},
		{"int8", args{4, []int8{1, 2, 3, 4}}, buildBuffer1(4), false},
		{"int8-FE", args{1, []int8{-2}}, []byte{0xFE}, false},
		{"[2]int8", args{8, [][2]int8{{1, 2}, {3, 4}}}, []byte{1, 2, 0, 0, 3, 4, 0, 0}, false},
		{"[3]int8", args{8, [][3]int8{{1, 2, 3}, {4, 5, 6}}}, []byte{1, 2, 3, 0, 4, 5, 6, 0}, false},
		{"[4]int8", args{8, [][4]int8{{1, 2, 3, 4}, {5, 6, 7, 8}}}, buildBuffer1(2 * 4), false},
		{"[2][2]int8", args{16, [][2][2]int8{
			{{1, 5}, {2, 6}},
			{{9, 13}, {10, 14}},
		}}, buildBuffer1(16, 2, 3, 6, 7, 10, 11, 14, 15), false},
		{"[3][3]int8", args{24, [][3][3]int8{
			{{1, 5, 9}, {2, 6, 10}, {3, 7, 11}},
			{{13, 17, 21}, {14, 18, 22}, {15, 19, 23}},
		}}, buildBuffer1(24, 3, 7, 11, 15, 19, 23), false},
		{"[4][4]int8", args{32, [][4][4]int8{
			{{1, 5, 9, 13}, {2, 6, 10, 14}, {3, 7, 11, 15}, {4, 8, 12, 16}},
			{{17, 21, 25, 29}, {18, 22, 26, 30}, {19, 23, 27, 31}, {20, 24, 28, 32}},
		}}, buildBuffer1(2 * 4 * 4), false},
		{"uint8", args{4, []uint8{1, 2, 3, 4}}, buildBuffer1(4), false},
		{"uint8-FE", args{1, []uint8{254}}, []byte{0xFE}, false},
		{"[2]uint8", args{8, [][2]uint8{{1, 2}, {3, 4}}}, []byte{1, 2, 0, 0, 3, 4, 0, 0}, false},
		{"[3]uint8", args{8, [][3]uint8{{1, 2, 3}, {4, 5, 6}}}, []byte{1, 2, 3, 0, 4, 5, 6, 0}, false},
		{"[4]uint8", args{8, [][4]uint8{{1, 2, 3, 4}, {5, 6, 7, 8}}}, buildBuffer1(2 * 4), false},
		{"[2][2]uint8", args{16, [][2][2]uint8{
			{{1, 5}, {2, 6}},
			{{9, 13}, {10, 14}},
		}}, buildBuffer1(16, 2, 3, 6, 7, 10, 11, 14, 15), false},
		{"[3][3]uint8", args{24, [][3][3]uint8{
			{{1, 5, 9}, {2, 6, 10}, {3, 7, 11}},
			{{13, 17, 21}, {14, 18, 22}, {15, 19, 23}},
		}}, buildBuffer1(24, 3, 7, 11, 15, 19, 23), false},
		{"[4][4]uint8", args{32, [][4][4]uint8{
			{{1, 5, 9, 13}, {2, 6, 10, 14}, {3, 7, 11, 15}, {4, 8, 12, 16}},
			{{17, 21, 25, 29}, {18, 22, 26, 30}, {19, 23, 27, 31}, {20, 24, 28, 32}},
		}}, buildBuffer1(2 * 4 * 4), false},
		{"int16", args{8, []int16{1, 2, 3, 4}}, buildBuffer2(4), false},
		{"int16-FE", args{2, []int16{-2}}, []byte{0xFE, 0xFF}, false},
		{"[2]int16", args{8, [][2]int16{{1, 2}, {3, 4}}}, buildBuffer2(2 * 2), false},
		{"[3]int16", args{16, [][3]int16{{1, 2, 3}, {4, 5, 6}}}, []byte{1, 0, 2, 0, 3, 0, 0, 0, 4, 0, 5, 0, 6, 0, 0, 0}, false},
		{"[4]int16", args{16, [][4]int16{{1, 2, 3, 4}, {5, 6, 7, 8}}}, buildBuffer2(2 * 4), false},
		{"[2][2]int16", args{16, [][2][2]int16{
			{{1, 3}, {2, 4}},
			{{5, 7}, {6, 8}},
		}}, buildBuffer2(2 * 2 * 2), false},
		{"[3][3]int16", args{48, [][3][3]int16{
			{{1, 5, 9}, {2, 6, 10}, {3, 7, 11}},
			{{13, 17, 21}, {14, 18, 22}, {15, 19, 23}},
		}}, buildBuffer2(24, 6, 7, 14, 15, 22, 23, 30, 31, 14, 15, 22, 23, 38, 39, 46, 47), false},
		{"[4][4]int16", args{64, [][4][4]int16{
			{{1, 5, 9, 13}, {2, 6, 10, 14}, {3, 7, 11, 15}, {4, 8, 12, 16}},
			{{17, 21, 25, 29}, {18, 22, 26, 30}, {19, 23, 27, 31}, {20, 24, 28, 32}},
		}}, buildBuffer2(2 * 4 * 4), false},
		{"uint16", args{8, []uint16{1, 2, 3, 4}}, buildBuffer2(4), false},
		{"uint16-FE", args{2, []uint16{65534}}, []byte{0xFE, 0xFF}, false},
		{"[2]uint16", args{8, [][2]uint16{{1, 2}, {3, 4}}}, buildBuffer2(2 * 2), false},
		{"[3]uint16", args{16, [][3]uint16{{1, 2, 3}, {4, 5, 6}}}, []byte{1, 0, 2, 0, 3, 0, 0, 0, 4, 0, 5, 0, 6, 0, 0, 0}, false},
		{"[4]uint16", args{16, [][4]uint16{{1, 2, 3, 4}, {5, 6, 7, 8}}}, buildBuffer2(2 * 4), false},
		{"[2][2]uint16", args{16, [][2][2]uint16{
			{{1, 3}, {2, 4}},
			{{5, 7}, {6, 8}},
		}}, buildBuffer2(2 * 2 * 2), false},
		{"[3][3]uint16", args{48, [][3][3]uint16{
			{{1, 5, 9}, {2, 6, 10}, {3, 7, 11}},
			{{13, 17, 21}, {14, 18, 22}, {15, 19, 23}},
		}}, buildBuffer2(24, 6, 7, 14, 15, 22, 23, 30, 31, 14, 15, 22, 23, 38, 39, 46, 47), false},
		{"[4][4]uint16", args{64, [][4][4]uint16{
			{{1, 5, 9, 13}, {2, 6, 10, 14}, {3, 7, 11, 15}, {4, 8, 12, 16}},
			{{17, 21, 25, 29}, {18, 22, 26, 30}, {19, 23, 27, 31}, {20, 24, 28, 32}},
		}}, buildBuffer2(2 * 4 * 4), false},
		{"uint32", args{16, []uint32{1, 2, 3, 4}}, buildBuffer3(4), false},
		{"uint32-FE", args{4, []uint32{4294967294}}, []byte{0xFE, 0xFF, 0xFF, 0xFF}, false},
		{"[2]uint32", args{16, [][2]uint32{{1, 2}, {3, 4}}}, buildBuffer3(4), false},
		{"[3]uint32", args{24, [][3]uint32{{1, 2, 3}, {4, 5, 6}}}, buildBuffer3(6), false},
		{"[4]uint32", args{32, [][4]uint32{{1, 2, 3, 4}, {5, 6, 7, 8}}}, buildBuffer3(8), false},
		{"[2][2]uint32", args{32, [][2][2]uint32{
			{{1, 3}, {2, 4}},
			{{5, 7}, {6, 8}},
		}}, buildBuffer3(8), false},
		{"[3][3]uint32", args{72, [][3][3]uint32{
			{{1, 4, 7}, {2, 5, 8}, {3, 6, 9}},
			{{10, 13, 16}, {11, 14, 17}, {12, 15, 18}},
		}}, buildBuffer3(18), false},
		{"[4][4]uint32", args{128, [][4][4]uint32{
			{{1, 5, 9, 13}, {2, 6, 10, 14}, {3, 7, 11, 15}, {4, 8, 12, 16}},
			{{17, 21, 25, 29}, {18, 22, 26, 30}, {19, 23, 27, 31}, {20, 24, 28, 32}},
		}}, buildBuffer3(32), false},
		{"float32", args{16, []float32{1, 2, 3, 4}}, buildBufferF(4), false},
		{"float32-FE", args{4, []float32{-2}}, []byte{0x00, 0x00, 0x00, 0xC0}, false},
		{"[2]float32", args{16, [][2]float32{{1, 2}, {3, 4}}}, buildBufferF(4), false},
		{"[3]float32", args{24, [][3]float32{{1, 2, 3}, {4, 5, 6}}}, buildBufferF(6), false},
		{"[4]float32", args{32, [][4]float32{{1, 2, 3, 4}, {5, 6, 7, 8}}}, buildBufferF(8), false},
		{"[2][2]float32", args{32, [][2][2]float32{
			{{1, 3}, {2, 4}},
			{{5, 7}, {6, 8}},
		}}, buildBufferF(8), false},
		{"[3][3]float32", args{72, [][3][3]float32{
			{{1, 4, 7}, {2, 5, 8}, {3, 6, 9}},
			{{10, 13, 16}, {11, 14, 17}, {12, 15, 18}},
		}}, buildBufferF(18), false},
		{"[4][4]float32", args{128, [][4][4]float32{
			{{1, 5, 9, 13}, {2, 6, 10, 14}, {3, 7, 11, 15}, {4, 8, 12, 16}},
			{{17, 21, 25, 29}, {18, 22, 26, 30}, {19, 23, 27, 31}, {20, 24, 28, 32}},
		}}, buildBufferF(32), false},
		{"color.RGBA", args{8, []color.RGBA{{1, 2, 3, 4}, {5, 6, 7, 8}}}, buildBuffer1(2 * 4), false},
		{"color.RGBA64", args{16, []color.RGBA64{{1, 2, 3, 4}, {5, 6, 7, 8}}}, buildBuffer2(2 * 4), false},
		{"gltf.RGBA", args{32, []gltf.RGBA{{R: 1, G: 2, B: 3, A: 4}, {R: 5, G: 6, B: 7, A: 8}}}, buildBufferF(2 * 4), false},
		{"gltf.RGB", args{24, []gltf.RGB{{R: 1, G: 2, B: 3}, {R: 4, G: 5, B: 6}}}, buildBufferF(2 * 3), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := make([]byte, tt.args.n)
			if err := Write(b, 0, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(b, tt.want) {
				t.Errorf("Write() error = %v, want %v", b, tt.want)
			}
		})
	}
}

func Test_ubyteComponent_Scalar(t *testing.T) {
	tests := []struct {
		name string
		want uint8
	}{
		{"base", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := make([]byte, 1)
			Ubyte.PutScalar(b, tt.want)
			if got := Ubyte.Scalar(b); got != tt.want {
				t.Errorf("ubyteComponent.Scalar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_byteComponent_Scalar(t *testing.T) {
	tests := []struct {
		name string
		want int8
	}{
		{"base", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := make([]byte, 1)
			Byte.PutScalar(b, tt.want)
			if got := Byte.Scalar(b); got != tt.want {
				t.Errorf("byteComponent.Scalar() = %v, want %v", got, tt.want)
			}
		})
	}
}
