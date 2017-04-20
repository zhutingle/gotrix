package test

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/zhutingle/gotrix/ecdh"
)

func TestIntAt(t *testing.T) {

	str1 := "0123456789abcdefghijklmnopqrstuvwxyz"
	for i := 0; i < len(str1); i++ {
		if int64(i) == ecdh.IntAt(str1, i) {
			t.Logf("字符[%c] IntAt 测试通过", str1[i])
		} else {
			t.Errorf("字符[%c] IntAt 测试不通过", str1[i])
		}
	}

	str2 := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for i := 0; i < len(str2); i++ {
		if int64(i) == ecdh.IntAt(str2, i) {
			t.Logf("字符[%c] IntAt 测试通过", str2[i])
		} else {
			t.Errorf("字符[%c] IntAt 测试不通过", str2[i])
		}
	}

	str3 := "~!@#$%^&*()_+-={}|\\][:\"';<>?/.,	 \n"
	for i := 0; i < len(str3); i++ {
		if -1 == ecdh.IntAt(str3, i) {
			t.Logf("字符[%c] IntAt 测试通过", str3[i])
		} else {
			t.Errorf("字符[%c] IntAt 测试不通过", str3[i])
		}
	}

}

func TestChunkSize(t *testing.T) {

	bi := ecdh.NewBigInteger()

	answers := map[int64]int64{2: 26, 3: 16, 4: 13, 5: 11, 6: 10, 7: 9, 8: 8, 9: 8, 10: 7, 11: 7, 12: 7, 13: 7, 14: 6, 15: 6, 16: 6, 17: 6, 18: 6, 19: 6, 20: 6, 21: 5, 22: 5, 23: 5, 24: 5, 25: 5, 26: 5, 27: 5, 28: 5, 29: 5, 30: 5, 31: 5, 32: 5, 33: 5, 34: 5, 35: 5, 36: 5, 37: 4, 38: 4, 39: 4, 40: 4, 41: 4, 42: 4, 43: 4, 44: 4, 45: 4, 46: 4, 47: 4, 48: 4, 49: 4, 50: 4, 51: 4, 52: 4, 53: 4, 54: 4, 55: 4, 56: 4, 57: 4, 58: 4, 59: 4, 60: 4, 61: 4, 62: 4, 63: 4}

	for k, v := range answers {
		if v == bi.ChunkSize(k) {
			t.Logf("值为[%d]时 ChunkSize 测试通过", k)
		} else {
			t.Errorf("值为[%d]时 ChunkSize 测试不通过", k)
		}
	}

}

func calculateVaue(bi *ecdh.BigInteger) int64 {
	if bi.T == 0 {
		return bi.S
	}
	value := int64(0)
	if bi.S == 0 {
		for j := int64(0); j < bi.T; j++ {
			value += bi.V[j] * int64(math.Pow(float64(ecdh.DV), float64(j)))
		}
	} else if bi.S < 0 {
		for j := int64(0); j < bi.T; j++ {
			if j == 0 {
				value += (ecdh.DV - bi.V[j]) * int64(math.Pow(float64(ecdh.DV), float64(j)))
			} else {
				value += (ecdh.DM - bi.V[j]) * int64(math.Pow(float64(ecdh.DV), float64(j)))
			}
		}
		value *= bi.S
	}
	return value
}

func TestFromInt(t *testing.T) {

	rand.Seed(time.Now().Unix())

	bi := ecdh.NewBigInteger()

	for i := int64(0); i < 100; i++ {

		randInt := rand.Int63()
		if rand.Intn(2) == 0 {
			randInt = -randInt
		}

		bi.FromInt(randInt)

		if calculateVaue(bi) == randInt {
			t.Logf("值为[%d]时 FromInt 测试通过", randInt)
		} else {
			t.Errorf("值为[%d]时 FromInt 测试不通过[%d]", randInt, calculateVaue(bi))
		}
	}
}

func TestFromRadix(t *testing.T) {

	rand.Seed(time.Now().Unix())

	bi := ecdh.NewBigInteger()

	for i := int64(0); i < 100; i++ {

		randInt := rand.Int63()
		if rand.Intn(2) == 0 {
			randInt = -randInt
		}

		bi.FromRadix(fmt.Sprintf("%d", randInt), 10)

		if calculateVaue(bi) == randInt {
			t.Logf("值为[%d]时 FromRadix 测试通过", randInt)
		} else {
			t.Errorf("值为[%d]时 FromRadix 测试不通过[%d]", randInt, calculateVaue(bi))
		}
	}
}

func TestSubTo(t *testing.T) {

	rand.Seed(time.Now().Unix())

	bi1 := ecdh.NewBigInteger()
	bi2 := ecdh.NewBigInteger()
	bi3 := ecdh.NewBigInteger()
	randInt1 := int64(0)
	randInt2 := int64(0)
	randInt3 := int64(0)
	for i := 0; i < 100; i++ {
		randInt1 = rand.Int63()
		if rand.Intn(2) == 0 {
			randInt1 = -randInt1
		}
		bi1.FromInt(randInt1)

		randInt2 = rand.Int63()
		if rand.Intn(2) == 0 {
			randInt2 = -randInt2
		}
		bi2.FromInt(randInt2)

		bi1.SubTo(bi2, bi3)
		randInt3 = randInt1 - randInt2

		if calculateVaue(bi3) == randInt3 {
			t.Logf(" SubTo 测试通过 %d - %d = %d ", randInt1, randInt2, randInt3)
		} else {
			t.Errorf(" SubTo 测试通过 %d - %d != %d ", randInt1, randInt2, randInt3)
		}

	}
}

func TestNegate(t *testing.T) {

	rand.Seed(time.Now().Unix())

	bi := ecdh.NewBigInteger()
	randInt := int64(0)

	for i := 0; i < 100; i++ {
		randInt = rand.Int63()
		if rand.Intn(2) == 0 {
			randInt = -randInt
		}
		bi.FromInt(randInt)

		bi = bi.Negate()
		randInt = -randInt

		if calculateVaue(bi) == randInt {
			t.Logf(" Negate 测试通过，正确值和实际值都为：%d", randInt)
		} else {
			t.Errorf(" Negate 测试不通过，正确值为：%d 实际值为：%d", randInt, calculateVaue(bi))
		}
	}
}

func TestDMultiply(t *testing.T) {

	rand.Seed(time.Now().Unix())

	bi := ecdh.NewBigInteger()
	randInt := int64(0)
	randInt2 := int64(0)

	for i := 0; i < 100; i++ {
		randInt = int64(rand.Int31())
		bi.FromInt(randInt)

		randInt2 = int64(rand.Int31())

		bi.DMultiply(randInt2)

		if calculateVaue(bi) == randInt * randInt2 {
			t.Log(" DMultiply 测试通过\n")
		} else {
			t.Errorf(" DMultiply 测试不通过：randInt=%d,randInt2=%d,randInt*randInt2=%d,bi=%d", randInt, randInt2, randInt * randInt2, calculateVaue(bi))
		}
	}
}

func TestLShiftTo(t *testing.T) {

	rand.Seed(time.Now().Unix())

	bi := ecdh.NewBigInteger()
	bi2 := ecdh.NewBigInteger()
	randInt := int64(0)
	randInt2 := int64(0)

	for i := 0; i < 100; i++ {
		randInt = int64(rand.Int31())
		if rand.Intn(2) == 0 {
			randInt = -randInt
		}
		bi.FromInt(randInt)

		randInt2 = rand.Int63n(32)

		bi.LShiftTo(randInt2, bi2)

		if calculateVaue(bi2) == (randInt << uint(randInt2)) {
			t.Log(" LShiftTo 测试通过\n")
		} else {
			t.Errorf(" LShiftTo 测试不通过：randInt=%d randInt2=%d 正确结果为：%d 实际结果为：%d", randInt, randInt2, randInt << uint(randInt2), calculateVaue(bi2))
		}
	}
}

func TestRShiftTo(t *testing.T) {

	rand.Seed(time.Now().Unix())

	bi := ecdh.NewBigInteger()
	bi2 := ecdh.NewBigInteger()
	randInt := int64(0)
	randInt2 := int64(0)

	for i := 0; i < 100; i++ {

		randInt = int64(rand.Int31())
		if rand.Intn(2) == 0 {
			randInt = -randInt
		}
		bi.FromInt(randInt)

		randInt2 = rand.Int63n(32)

		bi.RShiftTo(randInt2, bi2)

		if calculateVaue(bi2) == (randInt >> uint(randInt2)) {
			t.Log(" RShiftTo 测试通过\n")
		} else {
			t.Errorf(" RShiftTo 测试不通过：randInt=%d randInt2=%d 正确结果为：%d 实际结果为：%d", randInt, randInt2, randInt >> uint(randInt2), calculateVaue(bi2))
		}
	}
}

func TestCompareTo(t *testing.T) {

	rand.Seed(time.Now().Unix())

	bi1 := ecdh.NewBigInteger()
	bi2 := ecdh.NewBigInteger()
	randInt1 := int64(0)
	randInt2 := int64(0)

	for i := 0; i < 100; i++ {
		randInt1 = int64(rand.Int31())
		if rand.Intn(2) == 0 {
			randInt1 = -randInt1
		}
		bi1.FromInt(randInt1)

		randInt2 = int64(rand.Int31())
		if rand.Intn(2) == 0 {
			randInt2 = -randInt2
		}
		bi2.FromInt(randInt2)

		realResult := randInt1 - randInt2
		compareResult := bi1.CompareTo(bi2)

		if (realResult == 0 && compareResult == 0) || (realResult > 0 && compareResult > 0) || (realResult < 0 && compareResult < 0) {
			t.Log(" CompareTo 测试通过\n")
		} else {
			t.Errorf(" CompareTo 测试不通过：randInt1=%d randInt2=%d 正确结果为：%d 实际结果为：%d", randInt1, randInt2, realResult, compareResult)
		}

	}
}

func TestBitLength(t *testing.T) {

	rand.Seed(time.Now().Unix())
	bi1 := ecdh.NewBigInteger()
	randInt1 := int64(0)

	for i := 0; i < 100; i++ {
		randInt1 = rand.Int63()
		bi1.FromInt(randInt1)

		realResult := int64(0)
		for x := randInt1; x > 0; x = x >> 1 {
			realResult++
		}
		funcResult := bi1.BitLength()

		if realResult == funcResult {
			t.Logf("BitLength 测试通过：randInt1=%d,realResult=%d,funcResult=%d\n", randInt1, realResult, funcResult)
		} else {
			t.Errorf("BitLength 测试不通过：randInt1=%d,realResult=%d,funcResult=%d\n", randInt1, realResult, funcResult)
		}

	}
}

func TestMultiplyTo(t *testing.T) {

	rand.Seed(time.Now().Unix())

	bi1 := ecdh.NewBigInteger()
	bi2 := ecdh.NewBigInteger()
	randInt1 := int64(0)
	randInt2 := int64(0)
	biTemp := ecdh.NewBigInteger()

	for i := 0; i < 100; i++ {
		randInt1 = int64(rand.Int31())
		if rand.Intn(2) == 0 {
			randInt1 = -randInt1
		}
		bi1.FromInt(randInt1)

		randInt2 = int64(rand.Int31())
		if rand.Intn(2) == 0 {
			randInt2 = -randInt2
		}
		bi2.FromInt(randInt2)

		realResult := randInt1 * randInt2

		bi1.MultiplyTo(bi2, biTemp)
		funcResult := calculateVaue(biTemp)

		if realResult == funcResult {
			t.Log(" MultiplyTo 测试通过\n")
		} else {
			t.Errorf(" MultiplyTo 测试不通过：randInt1=%d randInt2=%d 正确结果为：%d 实际结果为：%d", randInt1, randInt2, realResult, funcResult)
		}

	}
}

func TestSquareTo(t *testing.T) {
	rand.Seed(time.Now().Unix())

	bi1 := ecdh.NewBigInteger()
	biResult := ecdh.NewBigInteger()
	randInt1 := int64(0)

	for i := 0; i < 100; i++ {
		randInt1 = int64(rand.Int31())
		if rand.Intn(2) == 0 {
			randInt1 = -randInt1
		}
		bi1.FromInt(randInt1)

		realResult := randInt1 * randInt1
		bi1.SquareTo(biResult)
		funcResult := calculateVaue(biResult)

		if realResult == funcResult {
			t.Log(" SquareTo 测试通过\n")
		} else {
			t.Errorf(" SquareTo 测试不通过：randInt1=%d,realResult=%d,funcResult=%d\n", randInt1, realResult, funcResult)
		}
	}
}

func TestDivRemTo(t *testing.T) {
	rand.Seed(time.Now().Unix())

	bi1 := ecdh.NewBigInteger()
	bi2 := ecdh.NewBigInteger()
	result1 := ecdh.NewBigInteger()
	result2 := ecdh.NewBigInteger()
	randInt1 := int64(0)
	randInt2 := int64(0)

	for i := 0; i < 100; i++ {
		randInt1 = rand.Int63()
		if rand.Intn(2) == 0 {
			randInt1 = -randInt1
		}
		bi1.FromInt(randInt1)

		randInt2 = int64(rand.Int31())
		if rand.Intn(2) == 0 {
			randInt2 = -randInt2
		}
		bi2.FromInt(randInt2)

		bi1.DivRemTo(bi2, result1, result2)

		if calculateVaue(result1) == randInt1 / randInt2 && calculateVaue(result2) == randInt1 % randInt2 {
			t.Log("DivRemTo 测试通过\n")
		} else {
			t.Errorf("DivRemTo 测试不通过：randInt1=%d,randInt2=%d,result1=%d,result2=%d\n", randInt1, randInt2, calculateVaue(result1), calculateVaue(result2))
		}

	}
}

func TestAddTo(t *testing.T) {

	rand.Seed(time.Now().Unix())

	bi1 := ecdh.NewBigInteger()
	bi2 := ecdh.NewBigInteger()
	bi3 := ecdh.NewBigInteger()
	randInt1 := int64(0)
	randInt2 := int64(0)
	randInt3 := int64(0)
	for i := 0; i < 100; i++ {
		randInt1 = rand.Int63()
		if rand.Intn(2) == 0 {
			randInt1 = -randInt1
		}
		bi1.FromInt(randInt1)

		randInt2 = rand.Int63()
		if rand.Intn(2) == 0 {
			randInt2 = -randInt2
		}
		bi2.FromInt(randInt2)

		bi1.AddTo(bi2, bi3)
		randInt3 = randInt1 + randInt2

		if calculateVaue(bi3) == randInt3 {
			t.Logf(" AddTo 测试通过 %d + %d = %d ", randInt1, randInt2, randInt3)
		} else {
			t.Errorf(" AddTo 测试通过 %d + %d != %d ", randInt1, randInt2, randInt3)
		}
	}
}

func TestInt2char(t *testing.T) {

	rand.Seed(time.Now().Unix())

	bi := ecdh.NewBigInteger()
	randInt := int64(0)

	for i := 0; i < 100; i++ {
		randInt = rand.Int63()
		if rand.Intn(2) == 0 {
			randInt = -randInt
		}
		bi.FromInt(randInt)

		if bi.ToString(16) == fmt.Sprintf("%x", randInt) {
			t.Log(" ToString(16) 测试成功\n")
		} else {
			t.Errorf(" ToString(16) 测试失败：randInt=%d", randInt)
		}
	}
}
