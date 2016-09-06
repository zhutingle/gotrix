package ecdh

import (
	"fmt"
	"strconv"
)

func keepEC() {
	fmt.Println("Hello World!!!")
}

type ECFieldElementFp struct {
	X *BigInteger
	Q *BigInteger
}

func NewECFieldElementFp(q *BigInteger, x *BigInteger) *ECFieldElementFp {
	var ecFieldElementFq *ECFieldElementFp = &ECFieldElementFp{}
	ecFieldElementFq.Q = q
	ecFieldElementFq.X = x
	return ecFieldElementFq
}

func (this *ECFieldElementFp) Equals(other *ECFieldElementFp) bool {
	if this == other {
		return true
	}
	return this.X.Equals(other.X) && this.Q.Equals(other.Q)
}

func (this *ECFieldElementFp) ToBigInteger() *BigInteger {
	return this.X
}

func (this *ECFieldElementFp) Negate() *ECFieldElementFp {
	return NewECFieldElementFp(this.Q, this.X.Negate().Mod(this.Q))
}

func (this *ECFieldElementFp) Add(b *ECFieldElementFp) *ECFieldElementFp {
	return NewECFieldElementFp(this.Q, this.X.Add(b.ToBigInteger()).Mod(this.Q))
}

func (this *ECFieldElementFp) Subtract(b *ECFieldElementFp) *ECFieldElementFp {
	return NewECFieldElementFp(this.Q, this.X.Subtract(b.ToBigInteger()).Mod(this.Q))
}

func (this *ECFieldElementFp) Multiply(b *ECFieldElementFp) *ECFieldElementFp {
	return NewECFieldElementFp(this.Q, this.X.Multiply(b.ToBigInteger()).Mod(this.Q))
}

func (this *ECFieldElementFp) Square() *ECFieldElementFp {
	return NewECFieldElementFp(this.Q, this.X.Square().Mod(this.Q))
}

func (this *ECFieldElementFp) Divide(b *ECFieldElementFp) *ECFieldElementFp {
	return NewECFieldElementFp(this.Q, this.X.Multiply(b.ToBigInteger().ModInverse(this.Q)).Mod(this.Q))
}

//ECPointFp
type ECPointFp struct {
	Curve *ECCurveFp
	X     *ECFieldElementFp
	Y     *ECFieldElementFp
	Z     *BigInteger
	Zinv  *BigInteger
}

func NewECPointFp(curve *ECCurveFp, x *ECFieldElementFp, y *ECFieldElementFp, z *BigInteger) *ECPointFp {
	var this *ECPointFp = &ECPointFp{}
	this.Curve = curve
	this.X = x
	this.Y = y
	if z == nil {
		this.Z = ONE
	} else {
		this.Z = z
	}
	this.Zinv = nil
	return this
}

func (this *ECPointFp) GetX() *ECFieldElementFp {
	if this.Zinv == nil {
		this.Zinv = this.Z.ModInverse(this.Curve.Q)
	}
	var r *BigInteger = this.X.ToBigInteger().Multiply(this.Zinv)
	this.Curve.Reduce(r)
	return this.Curve.FromBigInteger(r)
}

func (this *ECPointFp) GetY() *ECFieldElementFp {
	if this.Zinv == nil {
		this.Zinv = this.Z.ModInverse(this.Curve.Q)
	}
	var r *BigInteger = this.Y.ToBigInteger().Multiply(this.Zinv)
	this.Curve.Reduce(r)
	return this.Curve.FromBigInteger(r)
}

func (this *ECPointFp) Equals(other *ECPointFp) bool {
	if other == this {
		return true
	}
	if this.IsInfinity() {
		return other.IsInfinity()
	}
	if other.IsInfinity() {
		return this.IsInfinity()
	}
	var u *BigInteger = other.Y.ToBigInteger().Multiply(this.Z).Subtract(this.Y.ToBigInteger().Multiply(other.Z)).Mod(this.Curve.Q)
	if !u.Equals(ZERO) {
		return false
	}
	var v *BigInteger = other.X.ToBigInteger().Multiply(this.Z).Subtract(this.X.ToBigInteger().Multiply(other.Z)).Mod(this.Curve.Q)
	return v.Equals(ZERO)
}

func (this *ECPointFp) IsInfinity() bool {
	if this.X == nil && this.Y == nil {
		return true
	}
	return this.Z.Equals(ZERO) && !this.Y.ToBigInteger().Equals(ZERO)
}

func (this *ECPointFp) Negate() *ECPointFp {
	return NewECPointFp(this.Curve, this.X, this.Y.Negate(), this.Z)
}

func (this *ECPointFp) Add(b *ECPointFp) *ECPointFp {
	if this.IsInfinity() {
		return b
	}
	if b.IsInfinity() {
		return this
	}
	var u *BigInteger = b.Y.ToBigInteger().Multiply(this.Z).Subtract(this.Y.ToBigInteger().Multiply(b.Z)).Mod(this.Curve.Q)
	var v *BigInteger = b.X.ToBigInteger().Multiply(this.Z).Subtract(this.X.ToBigInteger().Multiply(b.Z)).Mod(this.Curve.Q)
	if ZERO.Equals(v) {
		if ZERO.Equals(u) {
			return this.Twice()
		}
		return this.Curve.GetInfinity()
	}
	var THREE = NewBigInteger()
	THREE.FromInt(3)
	var x1 *BigInteger = this.X.ToBigInteger()
	var y1 *BigInteger = this.Y.ToBigInteger()
	//	var x2 *BigInteger = b.X.ToBigInteger()
	//	var y2 *BigInteger = b.Y.ToBigInteger()

	var v2 *BigInteger = v.Square()
	var v3 *BigInteger = v2.Multiply(v)
	var x1v2 *BigInteger = x1.Multiply(v2)
	var zu2 *BigInteger = u.Square().Multiply(this.Z)

	var x3 *BigInteger = zu2.Subtract(x1v2.ShiftLeft(1)).Multiply(b.Z).Subtract(v3).Multiply(v).Mod(this.Curve.Q)
	var y3 *BigInteger = x1v2.Multiply(THREE).Multiply(u).Subtract(y1.Multiply(v3)).Subtract(zu2.Multiply(u)).Multiply(b.Z).Add(u.Multiply(v3)).Mod(this.Curve.Q)
	var z3 *BigInteger = v3.Multiply(this.Z).Multiply(b.Z).Mod(this.Curve.Q)

	return NewECPointFp(this.Curve, this.Curve.FromBigInteger(x3), this.Curve.FromBigInteger(y3), z3)
}

func (this *ECPointFp) Twice() *ECPointFp {
	if this.IsInfinity() {
		return this
	}
	if this.Y.ToBigInteger().Signum() == 0 {
		return this.Curve.GetInfinity()
	}

	var THREE *BigInteger = NewBigInteger()
	THREE.FromInt(3)
	var x1 *BigInteger = this.X.ToBigInteger()
	var y1 *BigInteger = this.Y.ToBigInteger()

	var y1z1 *BigInteger = y1.Multiply(this.Z)
	var y1sqz1 *BigInteger = y1z1.Multiply(y1).Mod(this.Curve.Q)
	var a *BigInteger = this.Curve.A.ToBigInteger()

	var w *BigInteger = x1.Square().Multiply(THREE)
	if !ZERO.Equals(a) {
		w = w.Add(this.Z.Square().Multiply(a))
	}
	w = w.Mod(this.Curve.Q)

	var x3 *BigInteger = w.Square().Subtract(x1.ShiftLeft(3).Multiply(y1sqz1)).ShiftLeft(1).Multiply(y1z1).Mod(this.Curve.Q)
	var y3 *BigInteger = w.Multiply(THREE).Multiply(x1).Subtract(y1sqz1.ShiftLeft(1)).ShiftLeft(2).Multiply(y1sqz1).Subtract(w.Square().Multiply(w)).Mod(this.Curve.Q)
	var z3 *BigInteger = y1z1.Square().Multiply(y1z1).ShiftLeft(3).Mod(this.Curve.Q)

	return NewECPointFp(this.Curve, this.Curve.FromBigInteger(x3), this.Curve.FromBigInteger(y3), z3)
}

// Simple NAF (Non-Adjacent Form) multiplication algorithm
// TODO: modularize the multiplication algorithm
func (this *ECPointFp) Multiply(k *BigInteger) *ECPointFp {
	if this.IsInfinity() {
		return this
	}
	if k.Signum() == 0 {
		return this.Curve.GetInfinity()
	}

	var e *BigInteger = k
	var h *BigInteger = e.Multiply(NewBigIntegerFromInt(3))

	var neg *ECPointFp = this.Negate()
	var R *ECPointFp = this

	for i := h.BitLength() - 2; i > 0; i-- {
		R = R.Twice()

		var hBit bool = h.TestBit(i)
		var eBit bool = e.TestBit(i)

		if hBit != eBit {
			if hBit {
				R = R.Add(this)
			} else {
				R = R.Add(neg)
			}
		}
	}
	return R
}

// Compute this*j + x*k (simultaneous multiplication)
func (this *ECPointFp) MultiplyTwo(j *BigInteger, x *ECPointFp, k *BigInteger) *ECPointFp {
	var i int64
	if j.BitLength() > k.BitLength() {
		i = j.BitLength() - 1
	} else {
		i = k.BitLength() - 1
	}

	var R *ECPointFp = this.Curve.GetInfinity()
	var both *ECPointFp = this.Add(x)
	for i >= 0 {
		R = R.Twice()
		if j.TestBit(i) {
			if k.TestBit(i) {
				R = R.Add(both)
			} else {
				R = R.Add(this)
			}
		} else {
			if k.TestBit(i) {
				R = R.Add(x)
			}
		}
		i--
	}
	return R
}

// ECCurveFp
type ECCurveFp struct {
	Q        *BigInteger
	A        *ECFieldElementFp
	B        *ECFieldElementFp
	Infinity *ECPointFp
	Reducer  *Barrett
}

func NewECCurveFp(q *BigInteger, a *BigInteger, b *BigInteger) *ECCurveFp {
	var this *ECCurveFp = &ECCurveFp{}
	this.Q = q
	this.A = this.FromBigInteger(a)
	this.B = this.FromBigInteger(b)
	this.Infinity = NewECPointFp(this, nil, nil, nil)
	this.Reducer = NewBarrett(this.Q)
	return this
}

func (this *ECCurveFp) GetQ() *BigInteger {
	return this.Q
}

func (this *ECCurveFp) GetA() *ECFieldElementFp {
	return this.A
}

func (this *ECCurveFp) GetB() *ECFieldElementFp {
	return this.B
}

func (this *ECCurveFp) Equals(other *ECCurveFp) bool {
	if other == this {
		return true
	}
	return this.Q.Equals(other.Q) && this.A.Equals(other.A) && this.B.Equals(other.B)
}

func (this *ECCurveFp) GetInfinity() *ECPointFp {
	return this.Infinity
}

func (this *ECCurveFp) FromBigInteger(x *BigInteger) *ECFieldElementFp {
	return NewECFieldElementFp(this.Q, x)
}

func (this *ECCurveFp) Reduce(x *BigInteger) {
	this.Reducer.Reduce(x)
}

// for now, work with hex strings because they're easier in JS
func (this *ECCurveFp) DecodePointHex(s string) *ECPointFp {
	i, err := strconv.ParseInt(s[0:2], 16, 32)
	if err != nil {
		return nil
	}
	switch i {
	case 0:
		return this.Infinity
	case 2:
	case 3:
		return nil
	case 4:
	case 6:
	case 7:
		var length int = (len(s) - 2) / 2
		var xHex string = s[2 : length+2]
		var yHex string = s[length+2 : length*2+2]
		bi1 := NewBigInteger()
		bi1.FromRadix(xHex, 16)
		bi2 := NewBigInteger()
		bi2.FromRadix(yHex, 16)
		return NewECPointFp(this, this.FromBigInteger(bi1), this.FromBigInteger(bi2), nil)
	default:
		return nil
	}
	return nil
}

func (this *ECCurveFp) EncodePointHex(p *ECPointFp) string {
	if p.IsInfinity() {
		return "00"
	}
	var xHex string = p.GetX().ToBigInteger().ToString(16)
	var yHex string = p.GetY().ToBigInteger().ToString(16)
	var oLen int = len(this.GetQ().ToString(16))
	if oLen%2 != 0 {
		oLen++
	}
	for len(xHex) < oLen {
		xHex = "0" + xHex
	}
	for len(yHex) < oLen {
		yHex = "0" + yHex
	}
	return "04" + xHex + yHex
}
