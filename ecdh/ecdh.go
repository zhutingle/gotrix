package ecdh

var Q string = "1461501637330902918203684832716283019653785059327"
var A string = "1461501637330902918203684832716283019653785059324"
var B string = "163235791306168110546604919403271579530548345413"
var GX string = "425826231723888350446541592701409065913635568770"
var GY string = "203520114162904107873991457957346892027982641970"
var N string = "1461501637330902918203687197606826779884643492439"

var rng *SecureRandom

func Get_rng() *SecureRandom {
	if rng == nil {
		rng = &SecureRandom{}
	}
	return rng
}

func Get_curve() *ECCurveFp {
	return NewECCurveFp(NewBigIntegerFromString(Q, 10), NewBigIntegerFromString(A, 10), NewBigIntegerFromString(B, 10))
}

func Get_G(curve *ECCurveFp) *ECPointFp {
	return NewECPointFp(curve, curve.FromBigInteger(NewBigIntegerFromString(GX, 10)), curve.FromBigInteger(NewBigIntegerFromString(GY, 10)), nil)
}

func Rand() *BigInteger {
	n := NewBigInteger()
	n.FromRadix(N, 10)
	n1 := n.Subtract(ONE)
	r := NewBigIntegerFromRandom(n.BitLength(), Get_rng())
	return r.Mod(n1).Add(ONE)
}

func SecretKey(a *BigInteger) *ECPointFp {
	var curve *ECCurveFp = Get_curve()
	var G *ECPointFp = Get_G(curve)
	var P *ECPointFp = G.Multiply(a)
	return P
}

func PublicKey(a *BigInteger, x string, y string) *ECPointFp {
	var curve *ECCurveFp = Get_curve()
	var G *ECPointFp = NewECPointFp(curve, curve.FromBigInteger(NewBigIntegerFromString(x, 16)), curve.FromBigInteger(NewBigIntegerFromString(y, 16)), nil)
	var P *ECPointFp = G.Multiply(a)
	return P
}
