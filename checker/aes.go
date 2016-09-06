package checker

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"
)

// BASE64
const (
	base64Table = "0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmno"
)

var coder = base64.NewEncoding(base64Table).WithPadding('|')

func Base64Encode(src []byte) []byte {
	return []byte(coder.EncodeToString(src))
}

func Base64Decode(src []byte) ([]byte, error) {
	return coder.DecodeString(string(src))
}

var sBox []byte = []byte{0x63, 0x7c, 0x77, 0x7b, 0xf2, 0x6b, 0x6f, 0xc5, 0x30, 0x01, 0x67, 0x2b, 0xfe, 0xd7, 0xab, 0x76, 0xca, 0x82, 0xc9, 0x7d, 0xfa, 0x59, 0x47, 0xf0, 0xad, 0xd4, 0xa2, 0xaf, 0x9c, 0xa4, 0x72, 0xc0, 0xb7, 0xfd, 0x93, 0x26, 0x36, 0x3f, 0xf7, 0xcc, 0x34, 0xa5, 0xe5, 0xf1, 0x71, 0xd8, 0x31, 0x15, 0x04, 0xc7, 0x23, 0xc3, 0x18, 0x96, 0x05, 0x9a, 0x07, 0x12, 0x80, 0xe2, 0xeb, 0x27, 0xb2, 0x75, 0x09, 0x83, 0x2c, 0x1a, 0x1b, 0x6e, 0x5a, 0xa0, 0x52, 0x3b, 0xd6, 0xb3, 0x29, 0xe3, 0x2f, 0x84, 0x53, 0xd1, 0x00, 0xed, 0x20, 0xfc, 0xb1, 0x5b, 0x6a, 0xcb, 0xbe, 0x39, 0x4a, 0x4c, 0x58, 0xcf, 0xd0, 0xef, 0xaa, 0xfb, 0x43, 0x4d, 0x33, 0x85, 0x45, 0xf9, 0x02, 0x7f, 0x50, 0x3c, 0x9f, 0xa8, 0x51, 0xa3, 0x40, 0x8f, 0x92, 0x9d, 0x38, 0xf5, 0xbc, 0xb6, 0xda, 0x21, 0x10, 0xff, 0xf3, 0xd2, 0xcd, 0x0c, 0x13,
	0xec, 0x5f, 0x97, 0x44, 0x17, 0xc4, 0xa7, 0x7e, 0x3d, 0x64, 0x5d, 0x19, 0x73, 0x60, 0x81, 0x4f, 0xdc, 0x22, 0x2a, 0x90, 0x88, 0x46, 0xee, 0xb8, 0x14, 0xde, 0x5e, 0x0b, 0xdb, 0xe0, 0x32, 0x3a, 0x0a, 0x49, 0x06, 0x24, 0x5c, 0xc2, 0xd3, 0xac, 0x62, 0x91, 0x95, 0xe4, 0x79, 0xe7, 0xc8, 0x37, 0x6d, 0x8d, 0xd5, 0x4e, 0xa9, 0x6c, 0x56, 0xf4, 0xea, 0x65, 0x7a, 0xae, 0x08, 0xba, 0x78, 0x25, 0x2e, 0x1c, 0xa6, 0xb4, 0xc6, 0xe8, 0xdd, 0x74, 0x1f, 0x4b, 0xbd, 0x8b, 0x8a, 0x70, 0x3e, 0xb5, 0x66, 0x48, 0x03, 0xf6, 0x0e, 0x61, 0x35, 0x57, 0xb9, 0x86, 0xc1, 0x1d, 0x9e, 0xe1, 0xf8, 0x98, 0x11, 0x69, 0xd9, 0x8e, 0x94, 0x9b, 0x1e, 0x87, 0xe9, 0xce, 0x55, 0x28, 0xdf, 0x8c, 0xa1, 0x89, 0x0d, 0xbf, 0xe6, 0x42, 0x68, 0x41, 0x99, 0x2d, 0x0f, 0xb0, 0x54, 0xbb, 0x16}

var rCon [][]byte = [][]byte{{0x00, 0x00, 0x00, 0x00}, {0x01, 0x00, 0x00, 0x00}, {0x02, 0x00, 0x00, 0x00}, {0x04, 0x00, 0x00, 0x00}, {0x08, 0x00, 0x00, 0x00}, {0x10, 0x00, 0x00, 0x00}, {0x20, 0x00, 0x00, 0x00}, {0x40, 0x00, 0x00, 0x00}, {0x80, 0x00, 0x00, 0x00}, {0x1b, 0x00, 0x00, 0x00}, {0x36, 0x00, 0x00, 0x00}}

func cipher(input []byte, w [][]byte) []byte {
	var Nb int = 4
	var Nr int = len(w)/Nb - 1

	var state [][]byte = make([][]byte, 4)
	for i := 0; i < len(state); i++ {
		state[i] = make([]byte, Nb)
	}
	for i := 0; i < 4*Nb; i++ {
		state[i%4][i/4] = input[i]
	}

	state = addRoundKey(state, w, 0, Nb)

	for round := 1; round < Nr; round++ {
		state = subBytes(state, Nb)
		state = shiftRows(state, Nb)
		state = mixColumns(state, Nb)
		state = addRoundKey(state, w, round, Nb)
	}

	state = subBytes(state, Nb)
	state = shiftRows(state, Nb)
	state = addRoundKey(state, w, Nr, Nb)

	var output []byte = make([]byte, 4*Nb)
	for i := 0; i < 4*Nb; i++ {
		output[i] = state[i%4][i/4]
	}
	return output
}

func Test(key []byte) [][]byte {
	return keyExpansion(key)
}

func keyExpansion(key []byte) [][]byte {
	var Nb int = 4
	var Nk int = len(key) / 4
	var Nr int = Nk + 6

	var w [][]byte = make([][]byte, Nb*(Nr+1))
	var temp []byte = make([]byte, 4)

	for i := 0; i < Nk; i++ {
		w[i] = []byte{key[4*i], key[4*i+1], key[4*i+2], key[4*i+3]}
	}

	for i := Nk; i < (Nb * (Nr + 1)); i++ {
		w[i] = make([]byte, 4)
		for t := 0; t < 4; t++ {
			temp[t] = w[i-1][t]
		}
		if i%Nk == 0 {
			temp = subWord(rotWord(temp))
			for t := 0; t < 4; t++ {
				temp[t] ^= rCon[i/Nk][t]
			}
		} else if Nk > 6 && i%Nk == 4 {
			temp = subWord(temp)
		}
		for t := 0; t < 4; t++ {
			w[i][t] = w[i-Nk][t] ^ temp[t]
		}
	}

	return w
}

func subBytes(s [][]byte, Nb int) [][]byte {
	for r := 0; r < 4; r++ {
		for c := 0; c < Nb; c++ {
			s[r][c] = sBox[s[r][c]]
		}
	}
	return s
}

func shiftRows(s [][]byte, Nb int) [][]byte {
	var t [4]byte
	for r := 1; r < 4; r++ {
		for c := 0; c < 4; c++ {
			t[c] = s[r][(c+r)%Nb]
		}
		for c := 0; c < 4; c++ {
			s[r][c] = t[c]
		}
	}
	return s
}

func mixColumns(s [][]byte, Nb int) [][]byte {
	for c := 0; c < 4; c++ {
		var a []byte = make([]byte, 4)
		var b []byte = make([]byte, 4)
		for i := 0; i < 4; i++ {
			a[i] = s[i][c]
			if s[i][c]&0x80 != 0 {
				b[i] = byte(int(s[i][c])<<1 ^ 0x011b)
			} else {
				b[i] = s[i][c] << 1
			}
		}
		s[0][c] = b[0] ^ a[1] ^ b[1] ^ a[2] ^ a[3]
		s[1][c] = a[0] ^ b[1] ^ a[2] ^ b[2] ^ a[3]
		s[2][c] = a[0] ^ a[1] ^ b[2] ^ a[3] ^ b[3]
		s[3][c] = a[0] ^ b[0] ^ a[1] ^ a[2] ^ b[3]
	}
	return s
}

func addRoundKey(state [][]byte, w [][]byte, rnd int, Nb int) [][]byte {
	for r := 0; r < 4; r++ {
		for c := 0; c < Nb; c++ {
			state[r][c] ^= w[rnd*4+c][r]
		}
	}
	return state
}

func subWord(w []byte) []byte {
	for i := 0; i < 4; i++ {
		w[i] = sBox[w[i]]
	}
	return w
}

func rotWord(w []byte) []byte {
	tmp := w[0]
	for i := 0; i < 3; i++ {
		w[i] = w[i+1]
	}
	w[3] = tmp
	return w
}

func AesEncrypt(plaintext []byte, password []byte, nBits int) (bytes []byte, err error) {
	var blockSize int = 16
	if !(nBits == 128 || nBits == 192 || nBits == 256) {
		err = errors.New("参数nBits只能为128、192、256中的一个")
		return
	}

	var nBytes int = nBits / 8
	var pwBytes []byte = make([]byte, nBytes)
	for i := 0; i < nBytes; i++ {
		if i < len(password) {
			pwBytes[i] = password[i]
		} else {
			pwBytes[i] = 0
		}
	}
	var key []byte = cipher(pwBytes, keyExpansion(pwBytes))
	key = append(key, key[0:nBytes-16]...)

	var counterBlock []byte = make([]byte, blockSize)

	var nonce int = int(time.Now().Nanosecond())
	var nonceMs int = nonce % 1000
	var nonceSec int = nonce / 1000
	var nonceRnd int = rand.Int() * 0xffff

	for i := 0; i < 2; i++ {
		counterBlock[i] = byte((uint(nonceMs) >> uint(i*8)) & 0xff)
	}
	for i := 0; i < 2; i++ {
		counterBlock[i+2] = byte((uint(nonceRnd) >> uint(i*8)) & 0xff)
	}
	for i := 0; i < 4; i++ {
		counterBlock[i+4] = byte((uint(nonceSec) >> uint(i*8)) & 0xff)
	}

	var keySchedule [][]byte = keyExpansion(key)
	var blockCount int = int(math.Ceil(float64(len(plaintext)) / float64(blockSize)))
	var bs []byte = make([]byte, (blockCount-1)*blockSize+(len(plaintext)-1)%blockSize+1+8)
	for i := 0; i < 8; i++ {
		bs[i] = counterBlock[i]
	}
	for b := 0; b < blockCount; b++ {

		for c := 0; c < 4; c++ {
			counterBlock[15-c] = byte((uint(b) >> uint(c*8)) & 0xff)
		}
		for c := 0; c < 4; c++ {
			counterBlock[15-c-4] = byte(uint(b/0x100000000) >> uint(c*8))
		}

		var cipherCntr []byte = cipher(counterBlock, keySchedule) // -- encrypt
		var blockLength int = 0
		if b < blockCount-1 {
			blockLength = blockSize
		} else {
			blockLength = (len(plaintext)-1)%blockSize + 1
		}

		var base int = b*blockSize + 8
		for i := 0; i < blockLength; i++ {
			bs[base+i] = cipherCntr[i] ^ plaintext[b*blockSize+i]
		}
	}
	var ciphertext []byte = Base64Encode(bs)
	return ciphertext, nil
}

func AesDecrypt(ciphertext []byte, password []byte, nBits int) (bytes []byte, err error) {
	var blockSize int = 16
	if !(nBits == 128 || nBits == 192 || nBits == 256) {
		err = errors.New("参数nBits只能为128、192、256中的一个")
		return
	}

	var cipherBytes []byte
	cipherBytes, err = Base64Decode(ciphertext)
	if err != nil {
		return
	}

	var nBytes int = nBits / 8
	var pwBytes []byte = make([]byte, nBytes)
	for i := 0; i < nBytes; i++ {
		if i < len(password) {
			pwBytes[i] = password[i]
		} else {
			pwBytes[i] = 0
		}
	}
	var key []byte = cipher(pwBytes, keyExpansion(pwBytes))
	key = append(key, key[0:nBytes-16]...)

	var counterBlock []byte = make([]byte, blockSize)

	for i := 0; i < 8; i++ {
		counterBlock[i] = cipherBytes[i]
	}

	var keySchedule [][]byte = keyExpansion(key)
	var blockCount int = int(math.Ceil(float64(len(cipherBytes)) / float64(blockSize)))
	var bs []byte = make([]byte, len(cipherBytes)-8)
	for b := 0; b < blockCount; b++ {

		for c := 0; c < 4; c++ {
			counterBlock[15-c] = byte((uint(b) >> uint(c*8)) & 0xff)
		}
		for c := 0; c < 4; c++ {
			counterBlock[15-c-4] = byte(uint(b/0x100000000) >> uint(c*8))
		}

		var cipherCntr []byte = cipher(counterBlock, keySchedule) // -- encrypt

		var base int = b * blockSize
		for i, blockLength := 0, int(math.Min(float64(len(cipherBytes)-8-base), float64(blockSize))); i < blockLength; i++ {
			bs[base+i] = cipherCntr[i] ^ cipherBytes[8+base+i]
		}
	}
	return bs, nil
}

func init() {
	fmt.Sprintln("%d,%d", rand.Int(), time.Now().Unix())
}
