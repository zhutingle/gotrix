package ecdh

import (
	"fmt"
	"log"
	"math"
	"reflect"
	"strconv"
)

func keep() {
	log.Println(reflect.TypeOf("Hello World!!!"))
	fmt.Println("Hello World!!!")
}

const (
	DB = 26
	DM = (1 << uint(DB)) - 1
	DV = 1 << uint(DB)

	BI_FP = 52
	FV    = 1 << uint(BI_FP)
	F1    = BI_FP - DB
	F2    = 2*DB - BI_FP

	BI_RM = "0123456789abcdefghijklmnopqrstuvwxyz"
)

// 大数 0
var ZERO = NewBigInteger()

// 大数 1
var ONE = NewBigInteger()

// map[uint8]int 以字符为 key ，这个字符所代表的数据为 value
var BI_RC = initBI_RC()

// 初始化 BI_RC
func initBI_RC() map[uint8]int64 {
	ZERO.FromInt(0)
	ONE.FromInt(1)

	tempMap := make(map[uint8]int64)
	for i := uint8(0); i < 10; i++ {
		tempMap["0"[0]+i] = int64(i)
	}

	for i := uint8(10); i < 36; i++ {
		tempMap["A"[0]+i-uint8(10)] = int64(i)
		tempMap["a"[0]+i-uint8(10)] = int64(i)
	}

	return tempMap
}

// 大数据计算的结构体
// S 用来代表正负
// T 代表当前这个大数需要几个 int 型的对象来表示
// V 大数分隔成 int 类型后保存的地方
type BigInteger struct {
	T int64
	S int64
	V map[int64]int64
}

// 包中包含的函数：
func NewBigInteger() *BigInteger {
	return &BigInteger{T: 0, S: 0, V: make(map[int64]int64)}
}

func NewBigIntegerFromString(s string, b int64) *BigInteger {
	r := NewBigInteger()
	r.FromString(s, b)
	return r
}

func NewBigIntegerFromInt(x int64) *BigInteger {
	r := NewBigInteger()
	r.FromInt(x)
	return r
}

func NewBigIntegerFromRandom(length int64, secureRandom *SecureRandom) *BigInteger {
	r := NewBigInteger()
	var t int64 = length & 7
	var arrayLength int = int(length>>3) + 1
	var x []byte = make([]byte, arrayLength)
	secureRandom.Next(x)
	if t > 0 {
		x[0] &= (1 << uint(t)) - 1
	} else {
		x[0] = 0
	}
	r.FromString(string(x), 256)
	return r
}

func Int2char(n int) string {
	return string(BI_RM[n])
}

// 结构体 BigInteger 的成员函数
func IntAt(str string, index int) int64 {
	if BI_RC[str[index]] > 0 || str[index] == "0"[0] {
		return BI_RC[str[index]]
	} else {
		return -1
	}
}

func (this *BigInteger) ChunkSize(r int64) int64 {
	return int64(math.Floor(math.Ln2 * float64(DB) / math.Log(float64(r))))
}

func (this *BigInteger) FromInt(x int64) {
	if x > 0 {
		i := int64(0)
		for {
			div := x / DV
			rem := x % DV
			if div == 0 && rem == 0 {
				break
			}
			this.V[i] = rem
			i++
			x = div
		}
		this.T = i
		this.S = 0
	} else if x < -1 {
		i := int64(0)
		for {
			div := x / DV
			rem := x % DV
			if div == 0 && rem == 0 {
				break
			}
			if i == 0 {
				this.V[i] = rem + DV
			} else {
				this.V[i] = rem + DM
			}
			i++
			x = div
		}
		this.T = i
		this.S = -1
	} else {
		this.T = 0
		this.S = x
	}
}

func (this *BigInteger) FromRadix(s string, b int64) {
	this.FromInt(0)
	if b <= 0 {
		b = 10
	}
	var cs int64 = this.ChunkSize(b)
	var d = int64(math.Pow(float64(b), float64(cs)))
	var mi bool = false
	var j int64 = 0
	var w int64 = 0
	for i := 0; i < len(s); i++ {
		x := IntAt(s, i)
		if x < 0 {
			if s[i] == '-' && this.Signum() == 0 {
				mi = true
			}
			continue
		}
		w = b*w + x
		j++
		if j >= cs {
			this.DMultiply(d)
			this.DAddOffset(w, 0)
			j = 0
			w = 0
		}
	}
	if j > 0 {
		this.DMultiply(int64(math.Pow(float64(b), float64(j))))
		this.DAddOffset(w, 0)
	}
	if mi {
		ZERO.SubTo(this, this)
	}
}

func (this *BigInteger) FromString(s string, b int64) {
	var k int64
	if b == 16 {
		k = 4
	} else if b == 8 {
		k = 3
	} else if b == 256 {
		k = 8
	} else if b == 2 {
		k = 1
	} else if b == 32 {
		k = 5
	} else if b == 4 {
		k = 2
	} else {
		this.FromRadix(s, b)
		return
	}
	this.T = 0
	this.S = 0
	var i int = len(s)
	var mi bool = false
	var sh int64 = 0
	for i--; i >= 0; i-- {
		var x int64
		if k == 8 {
			x = int64(s[i] & 0xFF)
		} else {
			x = IntAt(s, i)
		}
		if x < 0 {
			if s[i:i+1] == "-" {
				mi = true
			}
			continue
		}
		mi = false
		if sh == 0 {
			this.V[this.T] = x
			this.T++
		} else if sh+k > DB {
			this.V[this.T-1] |= (x & ((1 << uint(DB-sh)) - 1)) << uint(sh)
			this.V[this.T] = x >> uint(DB-sh)
			this.T++
		} else {
			this.V[this.T-1] |= x << uint(sh)
		}
		sh += k
		if sh >= DB {
			sh -= DB
		}
	}
	if k == 8 && (s[0]&0x80) != 0 {
		this.S = -1
		if sh > 0 {
			this.V[this.T-1] |= ((1 << uint(DB-sh)) - 1) << uint(sh)
		}
	}
	this.clamp()
	if mi {
		ZERO.SubTo(this, this)
	}
}

func (this *BigInteger) Signum() int {
	if this.S < 0 {
		return -1
	} else if this.T <= 0 || (this.T == 1 && this.V[0] <= 0) {
		return 0
	} else {
		return 1
	}
}

func (this *BigInteger) AM(i int64, x int64, w *BigInteger, j int64, c int64, n int64) int64 {
	n--
	for n >= 0 {
		n--
		v := x*this.V[i] + w.V[j] + c
		i++
		c = int64(math.Floor(float64(v / int64(0x4000000))))
		w.V[j] = v & 0x3FFFFFF
		j++
	}
	return c
}

func (this *BigInteger) clamp() {
	c := this.S & DM
	for this.T > 0 && this.V[this.T-1] == int64(c) {
		this.T--
	}
}

// (public) return string representation in given radix
func (this *BigInteger) ToString(b int64) string {
	if this.S < 0 {
		return "-" + this.Negate().ToString(b)
	}
	var k int64
	if b == 16 {
		k = 4
	} else if b == 8 {
		k = 3
	} else if b == 2 {
		k = 1
	} else if b == 32 {
		k = 5
	} else if b == 4 {
		k = 2
	} else {
		return this.ToRadix(b)
	}
	var km int64 = (1 << uint(k)) - 1
	var d int64
	var m bool = false
	var r string = ""
	var i int64 = this.T
	var p int64 = DB - (i*DB)%k
	if i > 0 {
		i--
		if p < DB {
			d = this.V[i] >> uint(p)
			if d > 0 {
				m = true
				r = Int2char(int(d))
			}
		}
		for i >= 0 {
			if p < k {
				d = this.V[i] & ((1 << uint(p)) - 1) << uint(k-p)
				i--
				p += DB - k
				d |= this.V[i] >> uint(p)
			} else {
				p -= k
				d = (this.V[i] >> uint(p)) & km
				if p <= 0 {
					p += DB
					i--
				}
			}
			if d > 0 {
				m = true
			}
			if m {
				r += Int2char(int(d))
			}
		}
	}
	if m {
		return r
	} else {
		return "0"
	}
}

func (this *BigInteger) ToRadix(b int64) string {
	if this.Signum() == 0 || b < 2 || b > 36 {
		return "0"
	}
	var cs int64 = this.ChunkSize(b)
	var a int64 = int64(math.Pow(float64(b), float64(cs)))
	var d *BigInteger = NewBigIntegerFromInt(a)
	var y *BigInteger = NewBigInteger()
	var z *BigInteger = NewBigInteger()
	var r string = ""
	this.DivRemTo(d, y, z)
	for y.Signum() > 0 {
		r = strconv.FormatInt(a+z.IntValue(), int(b))[1:] + r
		y.DivRemTo(d, y, z)
	}
	return strconv.FormatInt(z.IntValue(), int(b)) + r
	//  if(b == null) b = 10;
	//  if(this.signum() == 0 || b < 2 || b > 36) return "0";
	//  var cs = this.chunkSize(b);
	//  var a = Math.pow(b,cs);
	//  var d = nbv(a), y = nbi(), z = nbi(), r = "";
	//  this.divRemTo(d,y,z);
	//  while(y.signum() > 0) {
	//    r = (a+z.intValue()).toString(b).substr(1) + r;
	//    y.divRemTo(d,y,z);
	//  }
	//  return z.intValue().toString(b) + r;
}

func (this *BigInteger) DMultiply(n int64) {
	this.V[this.T] = this.AM(0, n-1, this, 0, 0, this.T)
	this.T++
	this.clamp()
}

func (this *BigInteger) DAddOffset(n int64, w int64) {
	if n == 0 {
		return
	}
	for this.T <= w {
		this.V[this.T] = 0
		this.T++
	}
	this.V[w] += n
	for this.V[w] >= DV {
		this.V[w] -= DV
		w++
		if w >= this.T {
			this.V[this.T] = 0
			this.T++
		}
		this.V[w]++
	}
}

func (this *BigInteger) SubTo(a *BigInteger, r *BigInteger) {
	var i int64 = 0
	var c int64 = 0
	var m int64 = int64(math.Min(float64(a.T), float64(this.T)))
	for i < m {
		c += this.V[i] - a.V[i]
		r.V[i] = c & DM
		i++
		c = c >> uint(DB)
	}
	if a.T < this.T {
		c -= a.S
		for i < this.T {
			c += this.V[i]
			r.V[i] = c & DM
			i++
			c = c >> uint(DB)
		}
		c += this.S
	} else {
		c += this.S
		for i < a.T {
			c -= a.V[i]
			r.V[i] = c & DM
			i++
			c = c >> DB
		}
		c -= a.S
	}
	if c < 0 {
		r.S = -1
	} else {
		r.S = 0
	}
	if c < -1 {
		r.V[i] = DV + c
		i++
	} else if c > 0 {
		r.V[i] = c
		i++
	}
	r.T = i
	r.clamp()
}

func (this *BigInteger) IntValue() int64 {
	if this.S < 0 {
		if this.T == 1 {
			return this.V[0] - DV
		} else if this.T == 0 {
			return -1
		}
	} else if this.T == 1 {
		return this.V[0]
	} else if this.T == 0 {
		return 0
	}
	return (this.V[1] & ((1 << uint(32-DB)) - 1) << uint(DB)) | this.V[0]
}

func (this *BigInteger) ByteValue() int64 {
	if this.T == 0 {
		return this.S
	} else {
		return (this.V[0] << 24) >> 24
	}
}

func (this *BigInteger) ShortValue() int64 {
	if this.T == 0 {
		return this.S
	} else {
		return (this.V[0] << 16) >> 16
	}
}

func (this *BigInteger) Negate() *BigInteger {
	r := NewBigInteger()
	ZERO.SubTo(this, r)
	return r
}

func (this *BigInteger) CopyTo(r *BigInteger) {
	for i := this.T - 1; i >= 0; i-- {
		r.V[i] = this.V[i]
	}
	r.T = this.T
	r.S = this.S
}

// (protected) r = this << n*DB
func (this *BigInteger) DLShiftTo(n int64, r *BigInteger) {
	for i := this.T - 1; i >= 0; i-- {
		r.V[i+n] = this.V[i]
	}
	for i := n - 1; i >= 0; i-- {
		r.V[i] = 0
	}
	r.T = this.T + n
	r.S = this.S
}

// (protected) r = this >> n*DB
func (this *BigInteger) DRShiftTo(n int64, r *BigInteger) {
	for i := n; i < this.T; i++ {
		r.V[i-n] = this.V[i]
	}
	r.T = int64(math.Max(float64(this.T-n), 0))
	r.S = this.S
}

// (protected) r = this << n
func (this *BigInteger) LShiftTo(n int64, r *BigInteger) {
	var bs int64 = n % DB
	var cbs int64 = DB - bs
	var bm int64 = (1 << uint(cbs)) - 1
	var ds int64 = n / DB
	var c = (this.S << uint(bs)) & DM
	for i := this.T - 1; i >= 0; i-- {
		r.V[i+ds+1] = (this.V[i] >> uint(cbs)) | c
		c = (this.V[i] & bm) << uint(bs)
	}
	for i := ds - 1; i >= 0; i-- {
		r.V[i] = 0
	}
	r.V[ds] = c
	r.T = this.T + ds + 1
	r.S = this.S
	r.clamp()
}

// (protected) r = this >> n
func (this *BigInteger) RShiftTo(n int64, r *BigInteger) {
	r.S = this.S
	var ds int64 = n / DB
	if ds >= this.T {
		r.T = 0
		return
	}
	var bs int64 = n % DB
	var cbs int64 = DB - bs
	var bm int64 = (1 << uint(bs)) - 1
	r.V[0] = this.V[ds] >> uint(bs)
	for i := ds + 1; i < this.T; i++ {
		r.V[i-ds-1] |= (this.V[i] & bm) << uint(cbs)
		r.V[i-ds] = this.V[i] >> uint(bs)
	}
	if bs > 0 {
		r.V[this.T-ds-1] |= (this.S & bm) << uint(cbs)
	}
	r.T = this.T - ds
	r.clamp()
}

// (public) this << n
func (this *BigInteger) ShiftLeft(n int64) *BigInteger {
	var r *BigInteger = NewBigInteger()
	if n < 0 {
		this.RShiftTo(-n, r)
	} else {
		this.LShiftTo(n, r)
	}
	return r
}

// (public) this >> n
func (this *BigInteger) ShiftRight(n int64) *BigInteger {
	var r *BigInteger = NewBigInteger()
	if n < 0 {
		this.LShiftTo(-n, r)
	} else {
		this.RShiftTo(n, r)
	}
	return r
}

func (this *BigInteger) Abs() *BigInteger {
	if this.S < 0 {
		return this.Negate()
	} else {
		return this
	}
}

// (public) return + if this > a, - if this < a, 0 if equal
func (this *BigInteger) CompareTo(a *BigInteger) int64 {
	var r int64 = this.S - a.S
	if r != 0 {
		return r
	}
	var i int64 = this.T
	r = i - a.T
	if r != 0 {
		if this.S < 0 {
			return -r
		} else {
			return r
		}
	}
	i--
	for i >= 0 {
		r = this.V[i] - a.V[i]
		if r != 0 {
			return r
		}
		i--
	}
	return 0
}

// returns bit length of the integer x
func Nbits(x int64) int64 {
	var r int64 = 1
	var t int64 = 0
	t = int64(uint64(x) >> 16)
	if t != 0 {
		x = t
		r += 16
	}
	t = int64(uint64(x) >> 8)
	if t != 0 {
		x = t
		r += 8
	}
	t = int64(uint64(x) >> 4)
	if t != 0 {
		x = t
		r += 4
	}
	t = int64(uint64(x) >> 2)
	if t != 0 {
		x = t
		r += 2
	}
	t = int64(uint64(x) >> 1)
	if t != 0 {
		x = t
		r += 1
	}
	return r
}

// (public) return the number of bits in "this"
func (this *BigInteger) BitLength() int64 {
	if this.T <= 0 {
		return 0
	}
	return DB*(this.T-1) + Nbits(this.V[this.T-1]^(this.S&DM))
}

// (public) true iff nth bit is set
func (this *BigInteger) TestBit(n int64) bool {
	var j int64 = int64(math.Floor(float64(n) / float64(DB)))
	if j >= this.T {
		return this.S != 0
	}
	return this.V[j]&(1<<uint(n%DB)) != 0
}

// (protected) r = this * a, r != this,a (HAC 14.12)
// "this" should be the larger one if appropriate.
func (this *BigInteger) MultiplyTo(a *BigInteger, r *BigInteger) {
	x := this.Abs()
	y := a.Abs()
	i := x.T
	r.T = i + y.T
	for i--; i >= 0; i-- {
		r.V[i] = 0
	}
	for i = 0; i < y.T; i++ {
		r.V[i+x.T] = x.AM(0, y.V[i], r, i, 0, x.T)
	}
	r.S = 0
	r.clamp()
	if this.S != a.S {
		ZERO.SubTo(r, r)
	}
}

// (protected) r = this^2, r != this (HAC 14.16)
func (this *BigInteger) SquareTo(r *BigInteger) {
	var x *BigInteger = this.Abs()
	r.T = 2 * x.T
	var i int64 = r.T
	for i--; i >= 0; i-- {
		r.V[i] = 0
	}
	for i = 0; i < x.T-1; i++ {
		var c int64 = x.AM(i, x.V[i], r, 2*i, 0, 1)
		r.V[i+x.T] += x.AM(i+1, 2*x.V[i], r, 2*i+1, c, x.T-i-1)
		if r.V[i+x.T] >= DV {
			r.V[i+x.T] -= DV
			r.V[i+x.T+1] = 1
		}
	}
	if r.T > 0 {
		r.V[r.T-1] += x.AM(i, x.V[i], r, 2*i, 0, 1)
	}
	r.S = 0
	r.clamp()
}

// (protected) divide this by m, quotient and remainder to q, r (HAC 14.20)
// r != q, this != m.  q or r may be null.
func (this *BigInteger) DivRemTo(m *BigInteger, q *BigInteger, r *BigInteger) {
	var pm *BigInteger = m.Abs()
	if pm.T <= 0 {
		return
	}
	var pt *BigInteger = this.Abs()
	if pt.T < pm.T {
		if q != nil {
			q.FromInt(0)
		}
		if r != nil {
			this.CopyTo(r)
		}
		return
	}
	if r == nil {
		r = NewBigInteger()
	}
	var y *BigInteger = NewBigInteger()
	var ts int64 = this.S
	var ms int64 = m.S
	var nsh int64 = DB - Nbits(pm.V[pm.T-1])
	if nsh > 0 {
		pm.LShiftTo(nsh, y)
		pt.LShiftTo(nsh, r)
	} else {
		pm.CopyTo(y)
		pt.CopyTo(r)
	}
	var ys int64 = y.T
	var y0 int64 = y.V[ys-1]
	if y0 == 0 {
		return
	}
	var yt int64 = y0 * (1 << uint(F1))
	if ys > 1 {
		yt += y.V[ys-2] >> uint(F2)
	}
	var d1 float64 = float64(FV) / float64(yt)
	var d2 float64 = float64(1<<uint(F1)) / float64(yt)
	var e int64 = 1 << uint(F2)
	var i int64 = r.T
	var j int64 = i - ys
	var t *BigInteger = q
	if q == nil {
		t = NewBigInteger()
	}
	y.DLShiftTo(j, t)
	if r.CompareTo(t) >= 0 {
		r.V[r.T] = 1
		r.T++
		r.SubTo(t, r)
	}
	ONE.DLShiftTo(ys, t)
	t.SubTo(y, y)
	for j--; j >= 0; j-- {
		var qd int64 = DM
		i--
		if r.V[i] != y0 {
			qd = int64(math.Floor(float64(r.V[i])*d1 + float64(r.V[i-1]+e)*d2))
		}
		r.V[i] += y.AM(0, qd, r, j, 0, ys)
		if (r.V[i]) > qd {
			y.DLShiftTo(j, t)
			r.SubTo(t, r)
			for qd--; r.V[i] < qd; qd-- {
				r.SubTo(t, r)
			}
		}
	}
	if q != nil {
		r.DRShiftTo(ys, q)
		if ts != ms {
			ZERO.SubTo(q, q)
		}
	}
	r.T = ys
	r.clamp()
	if nsh > 0 {
		r.RShiftTo(nsh, r)
	}
	if ts < 0 {
		ZERO.SubTo(r, r)
	}
}

// (public) this mod a
func (this *BigInteger) Mod(a *BigInteger) *BigInteger {
	var r *BigInteger = NewBigInteger()
	this.Abs().DivRemTo(a, nil, r)
	if this.S < 0 && r.CompareTo(ZERO) > 0 {
		a.SubTo(r, r)
	}
	return r
}

func (this *BigInteger) Equals(a *BigInteger) bool {
	return this.CompareTo(a) == 0
}

func (this *BigInteger) Min(a *BigInteger) *BigInteger {
	if this.CompareTo(a) < 0 {
		return this
	} else {
		return a
	}
}

func (this *BigInteger) Max(a *BigInteger) *BigInteger {
	if this.CompareTo(a) > 0 {
		return this
	} else {
		return a
	}
}

func (this *BigInteger) AddTo(a *BigInteger, r *BigInteger) {
	var i int64 = 0
	var c int64 = 0
	var m int64 = int64(math.Min(float64(a.T), float64(this.T)))
	for i < m {
		c += this.V[i] + a.V[i]
		r.V[i] = c & DM
		i++
		c >>= uint(DB)
	}
	if a.T < this.T {
		c += a.S
		for i < this.T {
			c += this.V[i]
			r.V[i] = c & DM
			i++
			c >>= uint(DB)
		}
		c += this.S
	} else {
		c += this.S
		for i < a.T {
			c += a.V[i]
			r.V[i] = c & DM
			i++
			c >>= uint(DB)
		}
		c += a.S
	}
	if c < 0 {
		r.S = -1
	} else {
		r.S = 0
	}
	if c > 0 {
		r.V[i] = c
		i++
	} else if c < 0 {
		r.V[i] = DV + c
		i++
	}
	r.T = i
	r.clamp()
}

func (this *BigInteger) Add(a *BigInteger) *BigInteger {
	r := NewBigInteger()
	this.AddTo(a, r)
	return r
}

func (this *BigInteger) Subtract(a *BigInteger) *BigInteger {
	r := NewBigInteger()
	this.SubTo(a, r)
	return r
}

func (this *BigInteger) Multiply(a *BigInteger) *BigInteger {
	r := NewBigInteger()
	this.MultiplyTo(a, r)
	return r
}

func (this *BigInteger) Square() *BigInteger {
	r := NewBigInteger()
	this.SquareTo(r)
	return r
}

func (this *BigInteger) Divide(a *BigInteger) *BigInteger {
	r := NewBigInteger()
	this.DivRemTo(a, r, nil)
	return r
}

func (this *BigInteger) Remainder(a *BigInteger) *BigInteger {
	r := NewBigInteger()
	this.DivRemTo(a, nil, r)
	return r
}

// (protected) true iff this is even
func (this *BigInteger) IsEven() bool {
	if this.T > 0 {
		return (this.V[0] & 1) == 0
	} else {
		return this.S == 0
	}
}

func (this *BigInteger) Clone() *BigInteger {
	var r *BigInteger = NewBigInteger()
	this.CopyTo(r)
	return r
}

// (public) 1/this % m (HAC 14.61)
func (this *BigInteger) ModInverse(m *BigInteger) *BigInteger {
	var ac bool = m.IsEven()
	if (this.IsEven() && ac) || m.Signum() == 0 {
		return ZERO
	}
	var u *BigInteger = m.Clone()
	var v *BigInteger = this.Clone()
	var a *BigInteger = NewBigInteger()
	a.FromInt(1)
	var b *BigInteger = NewBigInteger()
	b.FromInt(0)
	var c *BigInteger = NewBigInteger()
	c.FromInt(0)
	var d *BigInteger = NewBigInteger()
	d.FromInt(1)
	for u.Signum() != 0 {
		for u.IsEven() {
			u.RShiftTo(1, u)
			if ac {
				if !a.IsEven() || !b.IsEven() {
					a.AddTo(this, a)
					b.SubTo(m, b)
				}
				a.RShiftTo(1, a)
			} else if !b.IsEven() {
				b.SubTo(m, b)
			}
			b.RShiftTo(1, b)
		}
		for v.IsEven() {
			v.RShiftTo(1, v)
			if ac {
				if !c.IsEven() || !d.IsEven() {
					c.AddTo(this, c)
					d.SubTo(m, d)
				}
				c.RShiftTo(1, c)
			} else if !d.IsEven() {
				d.SubTo(m, d)
			}
			d.RShiftTo(1, d)
		}
		if u.CompareTo(v) >= 0 {
			u.SubTo(v, u)
			if ac {
				a.SubTo(c, a)
			}
			b.SubTo(d, b)
		} else {
			v.SubTo(u, v)
			if ac {
				c.SubTo(a, c)
			}
			d.SubTo(b, d)
		}
	}
	if v.CompareTo(ONE) != 0 {
		return ZERO
	}
	if d.CompareTo(m) >= 0 {
		return d.Subtract(m)
	}
	if d.Signum() < 0 {
		d.AddTo(m, d)
	} else {
		return d
	}
	if d.Signum() < 0 {
		return d.Add(m)
	} else {
		return d
	}
}

// (protected) r = lower n words of "this * a", a.t <= n
// "this" should be the larger one if appropriate.
func (this *BigInteger) MultiplyLowerTo(a *BigInteger, n int64, r *BigInteger) {
	var i int64 = int64(math.Min(float64(this.T+a.T), float64(n)))
	r.S = 0
	r.T = i
	for i > 0 {
		i--
		r.V[i] = 0
	}
	var j int64
	for j = r.T - this.T; i < j; i++ {
		r.V[i+this.T] = this.AM(0, a.V[i], r, i, 0, this.T)
	}
	for j = int64(math.Min(float64(a.T), float64(n))); i < j; i++ {
		this.AM(0, a.V[i], r, i, 0, n-i)
	}
	r.clamp()
}

// (protected) r = "this * a" without lower n words, n > 0
// "this" should be the larger one if appropriate.
func (this *BigInteger) MultiplyUpperTo(a *BigInteger, n int64, r *BigInteger) {
	n--
	r.T = this.T + a.T - n
	var i int64 = r.T
	r.S = 0
	for i--; i >= 0; i-- {
		r.V[i] = 0
	}
	for i = int64(math.Max(float64(n-this.T), 0)); i < a.T; i++ {
		r.V[this.T+i-n] = this.AM(n-i, a.V[i], r, 0, 0, this.T+i-n)
	}
	r.clamp()
	r.DRShiftTo(1, r)
}

// Barrett modular reduction
type Barrett struct {
	R2 *BigInteger
	Q3 *BigInteger
	MU *BigInteger
	M  *BigInteger
}

func NewBarrett(m *BigInteger) *Barrett {
	var this *Barrett = &Barrett{}
	this.R2 = NewBigInteger()
	this.Q3 = NewBigInteger()
	ONE.DLShiftTo(2*m.T, this.R2)
	this.MU = this.R2.Divide(m)
	this.M = m
	return this
}

func (this *Barrett) Convert(x *BigInteger) *BigInteger {
	if x.S < 0 || x.T > 2*this.M.T {
		return x.Mod(this.M)
	} else if x.CompareTo(this.M) < 0 {
		return x
	} else {
		var r *BigInteger = NewBigInteger()
		x.CopyTo(r)
		this.Reduce(r)
		return r
	}
}

func (this *Barrett) Revert(x *BigInteger) *BigInteger {
	return x
}

func (this *Barrett) Reduce(x *BigInteger) {
	x.DRShiftTo(this.M.T-1, this.R2)
	if x.T > this.M.T+1 {
		x.T = this.M.T + 1
		x.clamp()
	}
	this.MU.MultiplyUpperTo(this.R2, this.M.T+1, this.Q3)
	this.M.MultiplyLowerTo(this.Q3, this.M.T+1, this.R2)
	for x.CompareTo(this.R2) < 0 {
		x.DAddOffset(1, this.M.T+1)
	}
	x.SubTo(this.R2, x)
	for x.CompareTo(this.M) >= 0 {
		x.SubTo(this.M, x)
	}
}

func (this *Barrett) SqrtTo(x *BigInteger, r *BigInteger) {
	x.SquareTo(r)
	this.Reduce(r)
}

func (this *Barrett) MulTo(x *BigInteger, y *BigInteger, r *BigInteger) {
	x.MultiplyTo(y, r)
	this.Reduce(r)
}
