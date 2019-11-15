package binary

import (
	"encoding/binary"
	"math"
	"reflect"
	"testing"
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
		{"[2]int8", args{buildBuffer1(2 * 2), make([][2]int8, 2)}, [][2]int8{{1, 2}, {3, 4}}, false},
		{"[3]int8", args{buildBuffer1(2 * 3), make([][3]int8, 2)}, [][3]int8{{1, 2, 3}, {4, 5, 6}}, false},
		{"[4]int8", args{buildBuffer1(2 * 4), make([][4]int8, 2)}, [][4]int8{{1, 2, 3, 4}, {5, 6, 7, 8}}, false},
		{"[2][2]int8", args{buildBuffer1(16, 2, 3, 6, 7, 10, 11), make([][2][2]int8, 2)}, [][2][2]int8{
			{{1, 5}, {2, 6}},
			{{9, 13}, {10, 14}},
		}, false},
		{"[3][3]int8", args{buildBuffer1(32, 3, 7, 11, 15, 19, 23), make([][3][3]int8, 2)}, [][3][3]int8{
			{{1, 5, 9}, {2, 6, 10}, {3, 7, 11}},
			{{13, 17, 21}, {14, 18, 22}, {15, 19, 23}},
		}, false},
		{"[4][4]int8", args{buildBuffer1(2 * 4 * 4), make([][4][4]int8, 2)}, [][4][4]int8{
			{{1, 5, 9, 13}, {2, 6, 10, 14}, {3, 7, 11, 15}, {4, 8, 12, 16}},
			{{17, 21, 25, 29}, {18, 22, 26, 30}, {19, 23, 27, 31}, {20, 24, 28, 32}},
		}, false},
		{"uint8", args{buildBuffer1(4), make([]uint8, 4)}, []uint8{1, 2, 3, 4}, false},
		{"uint8-FE", args{[]byte{0xFE}, make([]uint8, 1)}, []uint8{254}, false},
		{"[2]uint8", args{buildBuffer1(2 * 2), make([][2]uint8, 2)}, [][2]uint8{{1, 2}, {3, 4}}, false},
		{"[3]uint8", args{buildBuffer1(2 * 3), make([][3]uint8, 2)}, [][3]uint8{{1, 2, 3}, {4, 5, 6}}, false},
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
		{"int16", args{buildBuffer2(2 * 4), make([]int16, 4)}, []int16{1, 2, 3, 4}, false},
		{"int16-FE", args{[]byte{0xFE, 0xFF}, make([]int16, 1)}, []int16{-2}, false},
		{"[2]int16", args{buildBuffer2(2 * 2 * 2), make([][2]int16, 2)}, [][2]int16{{1, 2}, {3, 4}}, false},
		{"[3]int16", args{buildBuffer2(2 * 3 * 2), make([][3]int16, 2)}, [][3]int16{{1, 2, 3}, {4, 5, 6}}, false},
		{"[4]int16", args{buildBuffer2(2 * 4 * 2), make([][4]int16, 2)}, [][4]int16{{1, 2, 3, 4}, {5, 6, 7, 8}}, false},
		{"[2][2]int16", args{buildBuffer2(2 * 2 * 2 * 2), make([][2][2]int16, 2)}, [][2][2]int16{
			{{1, 3}, {2, 4}},
			{{5, 7}, {6, 8}},
		}, false},
		{"[3][3]int16", args{buildBuffer2(32*2, 6, 7, 14, 15, 22, 23, 30, 31, 14, 15, 22, 23, 38, 39), make([][3][3]int16, 2)}, [][3][3]int16{
			{{1, 5, 9}, {2, 6, 10}, {3, 7, 11}},
			{{13, 17, 21}, {14, 18, 22}, {15, 19, 23}},
		}, false},
		{"[4][4]int16", args{buildBuffer2(2 * 4 * 4 * 2), make([][4][4]int16, 2)}, [][4][4]int16{
			{{1, 5, 9, 13}, {2, 6, 10, 14}, {3, 7, 11, 15}, {4, 8, 12, 16}},
			{{17, 21, 25, 29}, {18, 22, 26, 30}, {19, 23, 27, 31}, {20, 24, 28, 32}},
		}, false},
		{"uint16", args{buildBuffer2(2 * 4), make([]uint16, 4)}, []uint16{1, 2, 3, 4}, false},
		{"uint16-FE", args{[]byte{0xFE, 0xFF}, make([]uint16, 1)}, []uint16{65534}, false},
		{"[2]uint16", args{buildBuffer2(2 * 2 * 2), make([][2]uint16, 2)}, [][2]uint16{{1, 2}, {3, 4}}, false},
		{"[3]uint16", args{buildBuffer2(2 * 3 * 2), make([][3]uint16, 2)}, [][3]uint16{{1, 2, 3}, {4, 5, 6}}, false},
		{"[4]uint16", args{buildBuffer2(2 * 4 * 2), make([][4]uint16, 2)}, [][4]uint16{{1, 2, 3, 4}, {5, 6, 7, 8}}, false},
		{"[2][2]uint16", args{buildBuffer2(2 * 2 * 2 * 2), make([][2][2]uint16, 2)}, [][2][2]uint16{
			{{1, 3}, {2, 4}},
			{{5, 7}, {6, 8}},
		}, false},
		{"[3][3]uint16", args{buildBuffer2(32*2, 6, 7, 14, 15, 22, 23, 30, 31, 14, 15, 22, 23, 38, 39), make([][3][3]uint16, 2)}, [][3][3]uint16{
			{{1, 5, 9}, {2, 6, 10}, {3, 7, 11}},
			{{13, 17, 21}, {14, 18, 22}, {15, 19, 23}},
		}, false},
		{"[4][4]uint16", args{buildBuffer2(2 * 4 * 4 * 2), make([][4][4]uint16, 2)}, [][4][4]uint16{
			{{1, 5, 9, 13}, {2, 6, 10, 14}, {3, 7, 11, 15}, {4, 8, 12, 16}},
			{{17, 21, 25, 29}, {18, 22, 26, 30}, {19, 23, 27, 31}, {20, 24, 28, 32}},
		}, false},
		{"uint32", args{buildBuffer3(2 * 8), make([]uint32, 4)}, []uint32{1, 2, 3, 4}, false},
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
		{"uint16", args{buildBuffer2(2 * 4), make([]uint16, 4)}, []uint16{1, 2, 3, 4}, false},
		{"uint16-FE", args{[]byte{0xFE, 0xFF}, make([]uint16, 1)}, []uint16{65534}, false},
		{"[2]uint16", args{buildBuffer2(2 * 2 * 2), make([][2]uint16, 2)}, [][2]uint16{{1, 2}, {3, 4}}, false},
		{"[3]uint16", args{buildBuffer2(2 * 3 * 2), make([][3]uint16, 2)}, [][3]uint16{{1, 2, 3}, {4, 5, 6}}, false},
		{"[4]uint16", args{buildBuffer2(2 * 4 * 2), make([][4]uint16, 2)}, [][4]uint16{{1, 2, 3, 4}, {5, 6, 7, 8}}, false},
		{"[2][2]uint16", args{buildBuffer2(2 * 2 * 2 * 2), make([][2][2]uint16, 2)}, [][2][2]uint16{
			{{1, 3}, {2, 4}},
			{{5, 7}, {6, 8}},
		}, false},
		{"[3][3]uint16", args{buildBuffer2(32*2, 6, 7, 14, 15, 22, 23, 30, 31, 14, 15, 22, 23, 38, 39), make([][3][3]uint16, 2)}, [][3][3]uint16{
			{{1, 5, 9}, {2, 6, 10}, {3, 7, 11}},
			{{13, 17, 21}, {14, 18, 22}, {15, 19, 23}},
		}, false},
		{"[4][4]uint16", args{buildBuffer2(2 * 4 * 4 * 2), make([][4][4]uint16, 2)}, [][4][4]uint16{
			{{1, 5, 9, 13}, {2, 6, 10, 14}, {3, 7, 11, 15}, {4, 8, 12, 16}},
			{{17, 21, 25, 29}, {18, 22, 26, 30}, {19, 23, 27, 31}, {20, 24, 28, 32}},
		}, false},
		{"float32", args{buildBufferF(2 * 8), make([]float32, 4)}, []float32{1, 2, 3, 4}, false},
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
