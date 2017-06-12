// Copyright (c) 2005  Tom Wu
// All Rights Reserved.
// See "LICENSE" for details.

// Basic JavaScript BN library - subset useful for RSA encryption.

// Bits per digit
var dbits;

// JavaScript engine analysis
var canary = 0xdeadbeefcafe;
var j_lm = ((canary & 0xffffff) == 0xefcafe);

// (public) Constructor
function BigInteger(a, b, c) {
    if (a != null)
        if ("number" == typeof a)
            this.fromNumber(a, b, c);
        else if (b == null && "string" != typeof a)
            this.fromString(a, 256);
        else
            this.fromString(a, b);
}

// return new, unset BigInteger
function nbi() {
    return new BigInteger(null);
}

// am: Compute w_j += (x*this_i), propagate carries,
// c is initial carry, returns final carry.
// c < 3*dvalue, x < 2*dvalue, this_i < dvalue
// We need to select the fastest one that works in this environment.

// am1: use a single mult and divide to get the high bits,
// max digit bits should be 26 because
// max internal value = 2*dvalue^2-2*dvalue (< 2^53)
function am1(i, x, w, j, c, n) {
    while (--n >= 0) {
        var v = x * this[i++] + w[j] + c;
        c = Math.floor(v / 0x4000000);
        w[j++] = v & 0x3ffffff;
    }
    return c;
}
// am2 avoids a big mult-and-extract completely.
// Max digit bits should be <= 30 because we do bitwise ops
// on values up to 2*hdvalue^2-hdvalue-1 (< 2^31)
function am2(i, x, w, j, c, n) {
    var xl = x & 0x7fff, xh = x >> 15;
    while (--n >= 0) {
        var l = this[i] & 0x7fff;
        var h = this[i++] >> 15;
        var m = xh * l + h * xl;
        l = xl * l + ((m & 0x7fff) << 15) + w[j] + (c & 0x3fffffff);
        c = (l >>> 30) + (m >>> 15) + xh * h + (c >>> 30);
        w[j++] = l & 0x3fffffff;
    }
    return c;
}
// Alternately, set max digit bits to 28 since some
// browsers slow down when dealing with 32-bit numbers.
function am3(i, x, w, j, c, n) {
    var xl = x & 0x3fff, xh = x >> 14;
    while (--n >= 0) {
        var l = this[i] & 0x3fff;
        var h = this[i++] >> 14;
        var m = xh * l + h * xl;
        l = xl * l + ((m & 0x3fff) << 14) + w[j] + c;
        c = (l >> 28) + (m >> 14) + xh * h;
        w[j++] = l & 0xfffffff;
    }
    return c;
}
if (j_lm && (navigator.appName == "Microsoft Internet Explorer")) {
    BigInteger.prototype.am = am2;
    dbits = 30;
} else if (j_lm && (navigator.appName != "Netscape")) {
    BigInteger.prototype.am = am1;
    dbits = 26;
} else { // Mozilla/Netscape seems to prefer am3
    BigInteger.prototype.am = am3;
    dbits = 28;

    BigInteger.prototype.am = am1;
    dbits = 26;
}

BigInteger.prototype.DB = dbits;
BigInteger.prototype.DM = ((1 << dbits) - 1);
BigInteger.prototype.DV = (1 << dbits);

var BI_FP = 52;
BigInteger.prototype.FV = Math.pow(2, BI_FP);
BigInteger.prototype.F1 = BI_FP - dbits;
BigInteger.prototype.F2 = 2 * dbits - BI_FP;

// Digit conversions
var BI_RM = "0123456789abcdefghijklmnopqrstuvwxyz";
var BI_RC = new Array();
var rr, vv;
rr = "0".charCodeAt(0);
for (vv = 0; vv <= 9; ++vv)
    BI_RC[rr++] = vv;
rr = "a".charCodeAt(0);
for (vv = 10; vv < 36; ++vv)
    BI_RC[rr++] = vv;
rr = "A".charCodeAt(0);
for (vv = 10; vv < 36; ++vv)
    BI_RC[rr++] = vv;

function int2char(n) {
    return BI_RM.charAt(n);
}
function intAt(s, i) {
    var c = BI_RC[s.charCodeAt(i)];
    return (c == null) ? -1 : c;
}

// (protected) copy this to r
function bnpCopyTo(r) {
    for (var i = this.t - 1; i >= 0; --i)
        r[i] = this[i];
    r.t = this.t;
    r.s = this.s;
}

// (protected) set from integer value x, -DV <= x < DV
function bnpFromInt(x) {
    this.t = 1;
    this.s = (x < 0) ? -1 : 0;
    if (x > 0)
        this[0] = x;
    else if (x < -1)
        this[0] = x + this.DV;
    else
        this.t = 0;
}

// return bigint initialized to value
function nbv(i) {
    var r = nbi();
    r.fromInt(i);
    return r;
}

// (protected) set from string and radix
function bnpFromString(s, b) {
    var k;
    if (b == 16)
        k = 4;
    else if (b == 8)
        k = 3;
    else if (b == 256)
        k = 8; // byte array
    else if (b == 2)
        k = 1;
    else if (b == 32)
        k = 5;
    else if (b == 4)
        k = 2;
    else {
        this.fromRadix(s, b);
        return;
    }
    this.t = 0;
    this.s = 0;
    var i = s.length, mi = false, sh = 0;
    while (--i >= 0) {
        var x = (k == 8) ? s[i] & 0xff : intAt(s, i);
        if (x < 0) {
            if (s.charAt(i) == "-")
                mi = true;
            continue;
        }
        mi = false;
        if (sh == 0)
            this[this.t++] = x;
        else if (sh + k > this.DB) {
            this[this.t - 1] |= (x & ((1 << (this.DB - sh)) - 1)) << sh;
            this[this.t++] = (x >> (this.DB - sh));
        } else
            this[this.t - 1] |= x << sh;
        sh += k;
        if (sh >= this.DB)
            sh -= this.DB;
    }
    if (k == 8 && (s[0] & 0x80) != 0) {
        this.s = -1;
        if (sh > 0)
            this[this.t - 1] |= ((1 << (this.DB - sh)) - 1) << sh;
    }
    this.clamp();
    if (mi)
        BigInteger.ZERO.subTo(this, this);
}

// (protected) clamp off excess high words
function bnpClamp() {
    var c = this.s & this.DM;
    while (this.t > 0 && this[this.t - 1] == c)
        --this.t;
}

// (public) return string representation in given radix
function bnToString(b) {
    if (this.s < 0)
        return "-" + this.negate().toString(b);
    var k;
    if (b == 16)
        k = 4;
    else if (b == 8)
        k = 3;
    else if (b == 2)
        k = 1;
    else if (b == 32)
        k = 5;
    else if (b == 4)
        k = 2;
    else
        return this.toRadix(b);
    var km = (1 << k) - 1, d, m = false, r = "", i = this.t;
    var p = this.DB - (i * this.DB) % k;
    if (i-- > 0) {
        if (p < this.DB && (d = this[i] >> p) > 0) {
            m = true;
            r = int2char(d);
        }
        while (i >= 0) {
            if (p < k) {
                d = (this[i] & ((1 << p) - 1)) << (k - p);
                d |= this[--i] >> (p += this.DB - k);
            } else {
                d = (this[i] >> (p -= k)) & km;
                if (p <= 0) {
                    p += this.DB;
                    --i;
                }
            }
            if (d > 0)
                m = true;
            if (m)
                r += int2char(d);
        }
    }
    return m ? r : "0";
}

// (public) -this
function bnNegate() {
    var r = nbi();
    BigInteger.ZERO.subTo(this, r);
    return r;
}

// (public) |this|
function bnAbs() {
    return (this.s < 0) ? this.negate() : this;
}

// (public) return + if this > a, - if this < a, 0 if equal
function bnCompareTo(a) {
    var r = this.s - a.s;
    if (r != 0)
        return r;
    var i = this.t;
    r = i - a.t;
    if (r != 0)
        return (this.s < 0) ? -r : r;
    while (--i >= 0)
        if ((r = this[i] - a[i]) != 0)
            return r;
    return 0;
}

// returns bit length of the integer x
function nbits(x) {
    var r = 1, t;
    if ((t = x >>> 16) != 0) {
        x = t;
        r += 16;
    }
    if ((t = x >> 8) != 0) {
        x = t;
        r += 8;
    }
    if ((t = x >> 4) != 0) {
        x = t;
        r += 4;
    }
    if ((t = x >> 2) != 0) {
        x = t;
        r += 2;
    }
    if ((t = x >> 1) != 0) {
        x = t;
        r += 1;
    }
    return r;
}

// (public) return the number of bits in "this"
function bnBitLength() {
    if (this.t <= 0)
        return 0;
    return this.DB * (this.t - 1) + nbits(this[this.t - 1] ^ (this.s & this.DM));
}

// (protected) r = this << n*DB
function bnpDLShiftTo(n, r) {
    var i;
    for (i = this.t - 1; i >= 0; --i)
        r[i + n] = this[i];
    for (i = n - 1; i >= 0; --i)
        r[i] = 0;
    r.t = this.t + n;
    r.s = this.s;
}

// (protected) r = this >> n*DB
function bnpDRShiftTo(n, r) {
    for (var i = n; i < this.t; ++i)
        r[i - n] = this[i];
    r.t = Math.max(this.t - n, 0);
    r.s = this.s;
}

// (protected) r = this << n
function bnpLShiftTo(n, r) {
    var bs = n % this.DB;
    var cbs = this.DB - bs;
    var bm = (1 << cbs) - 1;
    var ds = Math.floor(n / this.DB), c = (this.s << bs) & this.DM, i;
    for (i = this.t - 1; i >= 0; --i) {
        r[i + ds + 1] = (this[i] >> cbs) | c;
        c = (this[i] & bm) << bs;
    }
    for (i = ds - 1; i >= 0; --i)
        r[i] = 0;
    r[ds] = c;
    r.t = this.t + ds + 1;
    r.s = this.s;
    r.clamp();
}

// (protected) r = this >> n
function bnpRShiftTo(n, r) {
    r.s = this.s;
    var ds = Math.floor(n / this.DB);
    if (ds >= this.t) {
        r.t = 0;
        return;
    }
    var bs = n % this.DB;
    var cbs = this.DB - bs;
    var bm = (1 << bs) - 1;
    r[0] = this[ds] >> bs;
    for (var i = ds + 1; i < this.t; ++i) {
        r[i - ds - 1] |= (this[i] & bm) << cbs;
        r[i - ds] = this[i] >> bs;
    }
    if (bs > 0)
        r[this.t - ds - 1] |= (this.s & bm) << cbs;
    r.t = this.t - ds;
    r.clamp();
}

// (protected) r = this - a
function bnpSubTo(a, r) {
    var i = 0, c = 0, m = Math.min(a.t, this.t);
    while (i < m) {
        c += this[i] - a[i];
        r[i++] = c & this.DM;
        c >>= this.DB;
    }
    if (a.t < this.t) {
        c -= a.s;
        while (i < this.t) {
            c += this[i];
            r[i++] = c & this.DM;
            c >>= this.DB;
        }
        c += this.s;
    } else {
        c += this.s;
        while (i < a.t) {
            c -= a[i];
            r[i++] = c & this.DM;
            c >>= this.DB;
        }
        c -= a.s;
    }
    r.s = (c < 0) ? -1 : 0;
    if (c < -1)
        r[i++] = this.DV + c;
    else if (c > 0)
        r[i++] = c;
    r.t = i;
    r.clamp();
}

// (protected) r = this * a, r != this,a (HAC 14.12)
// "this" should be the larger one if appropriate.
function bnpMultiplyTo(a, r) {
    var x = this.abs(), y = a.abs();
    var i = x.t;
    r.t = i + y.t;
    while (--i >= 0)
        r[i] = 0;
    for (i = 0; i < y.t; ++i)
        r[i + x.t] = x.am(0, y[i], r, i, 0, x.t);
    r.s = 0;
    r.clamp();
    if (this.s != a.s)
        BigInteger.ZERO.subTo(r, r);
}

// (protected) r = this^2, r != this (HAC 14.16)
function bnpSquareTo(r) {
    var x = this.abs();
    var i = r.t = 2 * x.t;
    while (--i >= 0)
        r[i] = 0;
    for (i = 0; i < x.t - 1; ++i) {
        var c = x.am(i, x[i], r, 2 * i, 0, 1);
        if ((r[i + x.t] += x.am(i + 1, 2 * x[i], r, 2 * i + 1, c, x.t - i - 1)) >= x.DV) {
            r[i + x.t] -= x.DV;
            r[i + x.t + 1] = 1;
        }
    }
    if (r.t > 0)
        r[r.t - 1] += x.am(i, x[i], r, 2 * i, 0, 1);
    r.s = 0;
    r.clamp();
}

// (protected) divide this by m, quotient and remainder to q, r (HAC 14.20)
// r != q, this != m.  q or r may be null.
function bnpDivRemTo(m, q, r) {
    var pm = m.abs();
    if (pm.t <= 0)
        return;
    var pt = this.abs();
    if (pt.t < pm.t) {
        if (q != null)
            q.fromInt(0);
        if (r != null)
            this.copyTo(r);
        return;
    }
    if (r == null)
        r = nbi();
    var y = nbi(), ts = this.s, ms = m.s;
    var nsh = this.DB - nbits(pm[pm.t - 1]); // normalize modulus
    if (nsh > 0) {
        pm.lShiftTo(nsh, y);
        pt.lShiftTo(nsh, r);
    } else {
        pm.copyTo(y);
        pt.copyTo(r);
    }
    var ys = y.t;
    var y0 = y[ys - 1];
    if (y0 == 0)
        return;
    var yt = y0 * (1 << this.F1) + ((ys > 1) ? y[ys - 2] >> this.F2 : 0);
    var d1 = this.FV / yt, d2 = (1 << this.F1) / yt, e = 1 << this.F2;
    var i = r.t, j = i - ys, t = (q == null) ? nbi() : q;
    y.dlShiftTo(j, t);
    if (r.compareTo(t) >= 0) {
        r[r.t++] = 1;
        r.subTo(t, r);
    }
    BigInteger.ONE.dlShiftTo(ys, t);
    t.subTo(y, y); // "negative" y so we can replace sub with am later
    while (y.t < ys)
        y[y.t++] = 0;
    while (--j >= 0) {
        // Estimate quotient digit
        var qd = (r[--i] == y0) ? this.DM : Math.floor(r[i] * d1 + (r[i - 1] + e) * d2);
        if ((r[i] += y.am(0, qd, r, j, 0, ys)) < qd) { // Try it out
            y.dlShiftTo(j, t);
            r.subTo(t, r);
            while (r[i] < --qd)
                r.subTo(t, r);
        }
    }
    if (q != null) {
        r.drShiftTo(ys, q);
        if (ts != ms)
            BigInteger.ZERO.subTo(q, q);
    }
    r.t = ys;
    r.clamp();
    if (nsh > 0)
        r.rShiftTo(nsh, r); // Denormalize remainder
    if (ts < 0)
        BigInteger.ZERO.subTo(r, r);
}

// (public) this mod a
function bnMod(a) {
    var r = nbi();
    this.abs().divRemTo(a, null, r);
    if (this.s < 0 && r.compareTo(BigInteger.ZERO) > 0)
        a.subTo(r, r);
    return r;
}

// Modular reduction using "classic" algorithm
function Classic(m) {
    this.m = m;
}
function cConvert(x) {
    if (x.s < 0 || x.compareTo(this.m) >= 0)
        return x.mod(this.m);
    else
        return x;
}
function cRevert(x) {
    return x;
}
function cReduce(x) {
    x.divRemTo(this.m, null, x);
}
function cMulTo(x, y, r) {
    x.multiplyTo(y, r);
    this.reduce(r);
}
function cSqrTo(x, r) {
    x.squareTo(r);
    this.reduce(r);
}

Classic.prototype.convert = cConvert;
Classic.prototype.revert = cRevert;
Classic.prototype.reduce = cReduce;
Classic.prototype.mulTo = cMulTo;
Classic.prototype.sqrTo = cSqrTo;

// (protected) return "-1/this % 2^DB"; useful for Mont. reduction
// justification:
// xy == 1 (mod m)
// xy = 1+km
// xy(2-xy) = (1+km)(1-km)
// x[y(2-xy)] = 1-k^2m^2
// x[y(2-xy)] == 1 (mod m^2)
// if y is 1/x mod m, then y(2-xy) is 1/x mod m^2
// should reduce x and y(2-xy) by m^2 at each step to keep size bounded.
// JS multiply "overflows" differently from C/C++, so care is needed here.
function bnpInvDigit() {
    if (this.t < 1)
        return 0;
    var x = this[0];
    if ((x & 1) == 0)
        return 0;
    var y = x & 3; // y == 1/x mod 2^2
    y = (y * (2 - (x & 0xf) * y)) & 0xf; // y == 1/x mod 2^4
    y = (y * (2 - (x & 0xff) * y)) & 0xff; // y == 1/x mod 2^8
    y = (y * (2 - (((x & 0xffff) * y) & 0xffff))) & 0xffff; // y == 1/x mod 2^16
    // last step - calculate inverse mod DV directly;
    // assumes 16 < DB <= 32 and assumes ability to handle 48-bit ints
    y = (y * (2 - x * y % this.DV)) % this.DV; // y == 1/x mod 2^dbits
    // we really want the negative inverse, and -DV < y < DV
    return (y > 0) ? this.DV - y : -y;
}

// Montgomery reduction
function Montgomery(m) {
    this.m = m;
    this.mp = m.invDigit();
    this.mpl = this.mp & 0x7fff;
    this.mph = this.mp >> 15;
    this.um = (1 << (m.DB - 15)) - 1;
    this.mt2 = 2 * m.t;
}

// xR mod m
function montConvert(x) {
    var r = nbi();
    x.abs().dlShiftTo(this.m.t, r);
    r.divRemTo(this.m, null, r);
    if (x.s < 0 && r.compareTo(BigInteger.ZERO) > 0)
        this.m.subTo(r, r);
    return r;
}

// x/R mod m
function montRevert(x) {
    var r = nbi();
    x.copyTo(r);
    this.reduce(r);
    return r;
}

// x = x/R mod m (HAC 14.32)
function montReduce(x) {
    while (x.t <= this.mt2)
        // pad x so am has enough room later
        x[x.t++] = 0;
    for (var i = 0; i < this.m.t; ++i) {
        // faster way of calculating u0 = x[i]*mp mod DV
        var j = x[i] & 0x7fff;
        var u0 = (j * this.mpl + (((j * this.mph + (x[i] >> 15) * this.mpl) & this.um) << 15)) & x.DM;
        // use am to combine the multiply-shift-add into one call
        j = i + this.m.t;
        x[j] += this.m.am(0, u0, x, i, 0, this.m.t);
        // propagate carry
        while (x[j] >= x.DV) {
            x[j] -= x.DV;
            x[++j]++;
        }
    }
    x.clamp();
    x.drShiftTo(this.m.t, x);
    if (x.compareTo(this.m) >= 0)
        x.subTo(this.m, x);
}

// r = "x^2/R mod m"; x != r
function montSqrTo(x, r) {
    x.squareTo(r);
    this.reduce(r);
}

// r = "xy/R mod m"; x,y != r
function montMulTo(x, y, r) {
    x.multiplyTo(y, r);
    this.reduce(r);
}

Montgomery.prototype.convert = montConvert;
Montgomery.prototype.revert = montRevert;
Montgomery.prototype.reduce = montReduce;
Montgomery.prototype.mulTo = montMulTo;
Montgomery.prototype.sqrTo = montSqrTo;

// (protected) true iff this is even
function bnpIsEven() {
    return ((this.t > 0) ? (this[0] & 1) : this.s) == 0;
}

// (protected) this^e, e < 2^32, doing sqr and mul with "r" (HAC 14.79)
function bnpExp(e, z) {
    if (e > 0xffffffff || e < 1)
        return BigInteger.ONE;
    var r = nbi(), r2 = nbi(), g = z.convert(this), i = nbits(e) - 1;
    g.copyTo(r);
    while (--i >= 0) {
        z.sqrTo(r, r2);
        if ((e & (1 << i)) > 0)
            z.mulTo(r2, g, r);
        else {
            var t = r;
            r = r2;
            r2 = t;
        }
    }
    return z.revert(r);
}

// (public) this^e % m, 0 <= e < 2^32
function bnModPowInt(e, m) {
    var z;
    if (e < 256 || m.isEven())
        z = new Classic(m);
    else
        z = new Montgomery(m);
    return this.exp(e, z);
}

// protected
BigInteger.prototype.copyTo = bnpCopyTo;
BigInteger.prototype.fromInt = bnpFromInt;
BigInteger.prototype.fromString = bnpFromString;
BigInteger.prototype.clamp = bnpClamp;
BigInteger.prototype.dlShiftTo = bnpDLShiftTo;
BigInteger.prototype.drShiftTo = bnpDRShiftTo;
BigInteger.prototype.lShiftTo = bnpLShiftTo;
BigInteger.prototype.rShiftTo = bnpRShiftTo;
BigInteger.prototype.subTo = bnpSubTo;
BigInteger.prototype.multiplyTo = bnpMultiplyTo;
BigInteger.prototype.squareTo = bnpSquareTo;
BigInteger.prototype.divRemTo = bnpDivRemTo;
BigInteger.prototype.invDigit = bnpInvDigit;
BigInteger.prototype.isEven = bnpIsEven;
BigInteger.prototype.exp = bnpExp;

// public
BigInteger.prototype.toString = bnToString;
BigInteger.prototype.negate = bnNegate;
BigInteger.prototype.abs = bnAbs;
BigInteger.prototype.compareTo = bnCompareTo;
BigInteger.prototype.bitLength = bnBitLength;
BigInteger.prototype.mod = bnMod;
BigInteger.prototype.modPowInt = bnModPowInt;

// "constants"
BigInteger.ZERO = nbv(0);
BigInteger.ONE = nbv(1);

// Copyright (c) 2005-2009  Tom Wu
// All Rights Reserved.
// See "LICENSE" for details.

// Extended JavaScript BN functions, required for RSA private ops.

// Version 1.1: new BigInteger("0", 10) returns "proper" zero
// Version 1.2: square() API, isProbablePrime fix

// (public)
function bnClone() { var r = nbi(); this.copyTo(r); return r; }

// (public) return value as integer
function bnIntValue() {
    if(this.s < 0) {
        if(this.t == 1) return this[0]-this.DV;
        else if(this.t == 0) return -1;
    }
    else if(this.t == 1) return this[0];
    else if(this.t == 0) return 0;
    // assumes 16 < DB < 32
    return ((this[1]&((1<<(32-this.DB))-1))<<this.DB)|this[0];
}

// (public) return value as byte
function bnByteValue() { return (this.t==0)?this.s:(this[0]<<24)>>24; }

// (public) return value as short (assumes DB>=16)
function bnShortValue() { return (this.t==0)?this.s:(this[0]<<16)>>16; }

// (protected) return x s.t. r^x < DV
function bnpChunkSize(r) {	return Math.floor(Math.LN2*this.DB/Math.log(r)); }

// (public) 0 if this == 0, 1 if this > 0
function bnSigNum() {
    if(this.s < 0) return -1;
    else if(this.t <= 0 || (this.t == 1 && this[0] <= 0)) return 0;
    else return 1;
}

// (protected) convert to radix string
function bnpToRadix(b) {
    if(b == null) b = 10;
    if(this.signum() == 0 || b < 2 || b > 36) return "0";
    var cs = this.chunkSize(b);
    var a = Math.pow(b,cs);
    var d = nbv(a), y = nbi(), z = nbi(), r = "";
    this.divRemTo(d,y,z);
    while(y.signum() > 0) {
        r = (a+z.intValue()).toString(b).substr(1) + r;
        y.divRemTo(d,y,z);
    }
    return z.intValue().toString(b) + r;
}

// (protected) convert from radix string
function bnpFromRadix(s,b) {
    this.fromInt(0);
    if(b == null) b = 10;
    var cs = this.chunkSize(b);
    var d = Math.pow(b,cs), mi = false, j = 0, w = 0;
    for(var i = 0; i < s.length; ++i) {
        var x = intAt(s,i);
        if(x < 0) {
            if(s.charAt(i) == "-" && this.signum() == 0) mi = true;
            continue;
        }
        w = b*w+x;
        if(++j >= cs) {
            this.dMultiply(d);
            this.dAddOffset(w,0);
            j = 0;
            w = 0;
        }
    }
    if(j > 0) {
        this.dMultiply(Math.pow(b,j));
        this.dAddOffset(w,0);
    }
    if(mi) BigInteger.ZERO.subTo(this,this);
}

// (protected) alternate constructor
function bnpFromNumber(a,b,c) {
    if("number" == typeof b) {
        // new BigInteger(int,int,RNG)
        if(a < 2) this.fromInt(1);
        else {
            this.fromNumber(a,c);
            if(!this.testBit(a-1))	// force MSB set
                this.bitwiseTo(BigInteger.ONE.shiftLeft(a-1),op_or,this);
            if(this.isEven()) this.dAddOffset(1,0); // force odd
            while(!this.isProbablePrime(b)) {
                this.dAddOffset(2,0);
                if(this.bitLength() > a) this.subTo(BigInteger.ONE.shiftLeft(a-1),this);
            }
        }
    }
    else {
        // new BigInteger(int,RNG)
        var x = new Array(), t = a&7;
        x.length = (a>>3)+1;
        b.nextBytes(x);
        if(t > 0) x[0] &= ((1<<t)-1); else x[0] = 0;
        this.fromString(x,256);
    }
}

// (public) convert to bigendian byte array
function bnToByteArray() {
    var i = this.t, r = new Array();
    r[0] = this.s;
    var p = this.DB-(i*this.DB)%8, d, k = 0;
    if(i-- > 0) {
        if(p < this.DB && (d = this[i]>>p) != (this.s&this.DM)>>p)
            r[k++] = d|(this.s<<(this.DB-p));
        while(i >= 0) {
            if(p < 8) {
                d = (this[i]&((1<<p)-1))<<(8-p);
                d |= this[--i]>>(p+=this.DB-8);
            }
            else {
                d = (this[i]>>(p-=8))&0xff;
                if(p <= 0) { p += this.DB; --i; }
            }
            if((d&0x80) != 0) d |= -256;
            if(k == 0 && (this.s&0x80) != (d&0x80)) ++k;
            if(k > 0 || d != this.s) r[k++] = d;
        }
    }
    return r;
}

function bnEquals(a) { return(this.compareTo(a)==0); }
function bnMin(a) { return(this.compareTo(a)<0)?this:a; }
function bnMax(a) { return(this.compareTo(a)>0)?this:a; }

// (protected) r = this op a (bitwise)
function bnpBitwiseTo(a,op,r) {
    var i, f, m = Math.min(a.t,this.t);
    for(i = 0; i < m; ++i) r[i] = op(this[i],a[i]);
    if(a.t < this.t) {
        f = a.s&this.DM;
        for(i = m; i < this.t; ++i) r[i] = op(this[i],f);
        r.t = this.t;
    }
    else {
        f = this.s&this.DM;
        for(i = m; i < a.t; ++i) r[i] = op(f,a[i]);
        r.t = a.t;
    }
    r.s = op(this.s,a.s);
    r.clamp();
}

// (public) this & a
function op_and(x,y) { return x&y; }
function bnAnd(a) { var r = nbi(); this.bitwiseTo(a,op_and,r); return r; }

// (public) this | a
function op_or(x,y) { return x|y; }
function bnOr(a) { var r = nbi(); this.bitwiseTo(a,op_or,r); return r; }

// (public) this ^ a
function op_xor(x,y) { return x^y; }
function bnXor(a) { var r = nbi(); this.bitwiseTo(a,op_xor,r); return r; }

// (public) this & ~a
function op_andnot(x,y) { return x&~y; }
function bnAndNot(a) { var r = nbi(); this.bitwiseTo(a,op_andnot,r); return r; }

// (public) ~this
function bnNot() {
    var r = nbi();
    for(var i = 0; i < this.t; ++i) r[i] = this.DM&~this[i];
    r.t = this.t;
    r.s = ~this.s;
    return r;
}

// (public) this << n
function bnShiftLeft(n) {
    var r = nbi();
    if(n < 0) this.rShiftTo(-n,r); else this.lShiftTo(n,r);
    return r;
}

// (public) this >> n
function bnShiftRight(n) {
    var r = nbi();
    if(n < 0) this.lShiftTo(-n,r); else this.rShiftTo(n,r);
    return r;
}

// return index of lowest 1-bit in x, x < 2^31
function lbit(x) {
    if(x == 0) return -1;
    var r = 0;
    if((x&0xffff) == 0) { x >>= 16; r += 16; }
    if((x&0xff) == 0) { x >>= 8; r += 8; }
    if((x&0xf) == 0) { x >>= 4; r += 4; }
    if((x&3) == 0) { x >>= 2; r += 2; }
    if((x&1) == 0) ++r;
    return r;
}

// (public) returns index of lowest 1-bit (or -1 if none)
function bnGetLowestSetBit() {
    for(var i = 0; i < this.t; ++i)
        if(this[i] != 0) return i*this.DB+lbit(this[i]);
    if(this.s < 0) return this.t*this.DB;
    return -1;
}

// return number of 1 bits in x
function cbit(x) {
    var r = 0;
    while(x != 0) { x &= x-1; ++r; }
    return r;
}

// (public) return number of set bits
function bnBitCount() {
    var r = 0, x = this.s&this.DM;
    for(var i = 0; i < this.t; ++i) r += cbit(this[i]^x);
    return r;
}

// (public) true iff nth bit is set
function bnTestBit(n) {
    var j = Math.floor(n/this.DB);
    if(j >= this.t) return(this.s!=0);
    return((this[j]&(1<<(n%this.DB)))!=0);
}

// (protected) this op (1<<n)
function bnpChangeBit(n,op) {
    var r = BigInteger.ONE.shiftLeft(n);
    this.bitwiseTo(r,op,r);
    return r;
}

// (public) this | (1<<n)
function bnSetBit(n) { return this.changeBit(n,op_or); }

// (public) this & ~(1<<n)
function bnClearBit(n) { return this.changeBit(n,op_andnot); }

// (public) this ^ (1<<n)
function bnFlipBit(n) { return this.changeBit(n,op_xor); }

// (protected) r = this + a
function bnpAddTo(a,r) {
    var i = 0, c = 0, m = Math.min(a.t,this.t);
    while(i < m) {
        c += this[i]+a[i];
        r[i++] = c&this.DM;
        c >>= this.DB;
    }
    if(a.t < this.t) {
        c += a.s;
        while(i < this.t) {
            c += this[i];
            r[i++] = c&this.DM;
            c >>= this.DB;
        }
        c += this.s;
    }
    else {
        c += this.s;
        while(i < a.t) {
            c += a[i];
            r[i++] = c&this.DM;
            c >>= this.DB;
        }
        c += a.s;
    }
    r.s = (c<0)?-1:0;
    if(c > 0) r[i++] = c;
    else if(c < -1) r[i++] = this.DV+c;
    r.t = i;
    r.clamp();
}

// (public) this + a
function bnAdd(a) { var r = nbi(); this.addTo(a,r); return r; }

// (public) this - a
function bnSubtract(a) { var r = nbi(); this.subTo(a,r); return r; }

// (public) this * a
function bnMultiply(a) { var r = nbi(); this.multiplyTo(a,r); return r; }

// (public) this^2
function bnSquare() { var r = nbi(); this.squareTo(r); return r; }

// (public) this / a
function bnDivide(a) { var r = nbi(); this.divRemTo(a,r,null); return r; }

// (public) this % a
function bnRemainder(a) { var r = nbi(); this.divRemTo(a,null,r); return r; }

// (public) [this/a,this%a]
function bnDivideAndRemainder(a) {
    var q = nbi(), r = nbi();
    this.divRemTo(a,q,r);
    return new Array(q,r);
}

// (protected) this *= n, this >= 0, 1 < n < DV
function bnpDMultiply(n) {
    this[this.t] = this.am(0,n-1,this,0,0,this.t);
    ++this.t;
    this.clamp();
}

// (protected) this += n << w words, this >= 0
function bnpDAddOffset(n,w) {
    if(n == 0) return;
    while(this.t <= w) this[this.t++] = 0;
    this[w] += n;
    while(this[w] >= this.DV) {
        this[w] -= this.DV;
        if(++w >= this.t) this[this.t++] = 0;
        ++this[w];
    }
}

// A "null" reducer
function NullExp() {}
function nNop(x) { return x; }
function nMulTo(x,y,r) { x.multiplyTo(y,r); }
function nSqrTo(x,r) { x.squareTo(r); }

NullExp.prototype.convert = nNop;
NullExp.prototype.revert = nNop;
NullExp.prototype.mulTo = nMulTo;
NullExp.prototype.sqrTo = nSqrTo;

// (public) this^e
function bnPow(e) { return this.exp(e,new NullExp()); }

// (protected) r = lower n words of "this * a", a.t <= n
// "this" should be the larger one if appropriate.
function bnpMultiplyLowerTo(a,n,r) {
    var i = Math.min(this.t+a.t,n);
    r.s = 0; // assumes a,this >= 0
    r.t = i;
    while(i > 0) r[--i] = 0;
    var j;
    for(j = r.t-this.t; i < j; ++i) r[i+this.t] = this.am(0,a[i],r,i,0,this.t);
    for(j = Math.min(a.t,n); i < j; ++i) this.am(0,a[i],r,i,0,n-i);
    r.clamp();
}

// (protected) r = "this * a" without lower n words, n > 0
// "this" should be the larger one if appropriate.
function bnpMultiplyUpperTo(a,n,r) {
    --n;
    var i = r.t = this.t+a.t-n;
    r.s = 0; // assumes a,this >= 0
    while(--i >= 0) r[i] = 0;
    for(i = Math.max(n-this.t,0); i < a.t; ++i)
        r[this.t+i-n] = this.am(n-i,a[i],r,0,0,this.t+i-n);
    r.clamp();
    r.drShiftTo(1,r);
}

// Barrett modular reduction
function Barrett(m) {
    // setup Barrett
    this.r2 = nbi();
    this.q3 = nbi();
    BigInteger.ONE.dlShiftTo(2*m.t,this.r2);
    this.mu = this.r2.divide(m);
    this.m = m;
}

function barrettConvert(x) {
    if(x.s < 0 || x.t > 2*this.m.t) return x.mod(this.m);
    else if(x.compareTo(this.m) < 0) return x;
    else { var r = nbi(); x.copyTo(r); this.reduce(r); return r; }
}

function barrettRevert(x) { return x; }

// x = x mod m (HAC 14.42)
function barrettReduce(x) {
    x.drShiftTo(this.m.t-1,this.r2);
    if(x.t > this.m.t+1) { x.t = this.m.t+1; x.clamp(); }
    this.mu.multiplyUpperTo(this.r2,this.m.t+1,this.q3);
    this.m.multiplyLowerTo(this.q3,this.m.t+1,this.r2);
    while(x.compareTo(this.r2) < 0) x.dAddOffset(1,this.m.t+1);
    x.subTo(this.r2,x);
    while(x.compareTo(this.m) >= 0) x.subTo(this.m,x);
}

// r = x^2 mod m; x != r
function barrettSqrTo(x,r) { x.squareTo(r); this.reduce(r); }

// r = x*y mod m; x,y != r
function barrettMulTo(x,y,r) { x.multiplyTo(y,r); this.reduce(r); }

Barrett.prototype.convert = barrettConvert;
Barrett.prototype.revert = barrettRevert;
Barrett.prototype.reduce = barrettReduce;
Barrett.prototype.mulTo = barrettMulTo;
Barrett.prototype.sqrTo = barrettSqrTo;

// (public) this^e % m (HAC 14.85)
function bnModPow(e,m) {
    var i = e.bitLength(), k, r = nbv(1), z;
    if(i <= 0) return r;
    else if(i < 18) k = 1;
    else if(i < 48) k = 3;
    else if(i < 144) k = 4;
    else if(i < 768) k = 5;
    else k = 6;
    if(i < 8)
        z = new Classic(m);
    else if(m.isEven())
        z = new Barrett(m);
    else
        z = new Montgomery(m);

    // precomputation
    var g = new Array(), n = 3, k1 = k-1, km = (1<<k)-1;
    g[1] = z.convert(this);
    if(k > 1) {
        var g2 = nbi();
        z.sqrTo(g[1],g2);
        while(n <= km) {
            g[n] = nbi();
            z.mulTo(g2,g[n-2],g[n]);
            n += 2;
        }
    }

    var j = e.t-1, w, is1 = true, r2 = nbi(), t;
    i = nbits(e[j])-1;
    while(j >= 0) {
        if(i >= k1) w = (e[j]>>(i-k1))&km;
        else {
            w = (e[j]&((1<<(i+1))-1))<<(k1-i);
            if(j > 0) w |= e[j-1]>>(this.DB+i-k1);
        }

        n = k;
        while((w&1) == 0) { w >>= 1; --n; }
        if((i -= n) < 0) { i += this.DB; --j; }
        if(is1) {	// ret == 1, don't bother squaring or multiplying it
            g[w].copyTo(r);
            is1 = false;
        }
        else {
            while(n > 1) { z.sqrTo(r,r2); z.sqrTo(r2,r); n -= 2; }
            if(n > 0) z.sqrTo(r,r2); else { t = r; r = r2; r2 = t; }
            z.mulTo(r2,g[w],r);
        }

        while(j >= 0 && (e[j]&(1<<i)) == 0) {
            z.sqrTo(r,r2); t = r; r = r2; r2 = t;
            if(--i < 0) { i = this.DB-1; --j; }
        }
    }
    return z.revert(r);
}

// (public) gcd(this,a) (HAC 14.54)
function bnGCD(a) {
    var x = (this.s<0)?this.negate():this.clone();
    var y = (a.s<0)?a.negate():a.clone();
    if(x.compareTo(y) < 0) { var t = x; x = y; y = t; }
    var i = x.getLowestSetBit(), g = y.getLowestSetBit();
    if(g < 0) return x;
    if(i < g) g = i;
    if(g > 0) {
        x.rShiftTo(g,x);
        y.rShiftTo(g,y);
    }
    while(x.signum() > 0) {
        if((i = x.getLowestSetBit()) > 0) x.rShiftTo(i,x);
        if((i = y.getLowestSetBit()) > 0) y.rShiftTo(i,y);
        if(x.compareTo(y) >= 0) {
            x.subTo(y,x);
            x.rShiftTo(1,x);
        }
        else {
            y.subTo(x,y);
            y.rShiftTo(1,y);
        }
    }
    if(g > 0) y.lShiftTo(g,y);
    return y;
}

// (protected) this % n, n < 2^26
function bnpModInt(n) {
    if(n <= 0) return 0;
    var d = this.DV%n, r = (this.s<0)?n-1:0;
    if(this.t > 0)
        if(d == 0) r = this[0]%n;
        else for(var i = this.t-1; i >= 0; --i) r = (d*r+this[i])%n;
    return r;
}

// (public) 1/this % m (HAC 14.61)
function bnModInverse(m) {
    var ac = m.isEven();
    if((this.isEven() && ac) || m.signum() == 0) return BigInteger.ZERO;
    var u = m.clone(), v = this.clone();
    var a = nbv(1), b = nbv(0), c = nbv(0), d = nbv(1);
    while(u.signum() != 0) {
        while(u.isEven()) {
            u.rShiftTo(1,u);
            if(ac) {
                if(!a.isEven() || !b.isEven()) { a.addTo(this,a); b.subTo(m,b); }
                a.rShiftTo(1,a);
            }
            else if(!b.isEven()) b.subTo(m,b);
            b.rShiftTo(1,b);
        }
        while(v.isEven()) {
            v.rShiftTo(1,v);
            if(ac) {
                if(!c.isEven() || !d.isEven()) { c.addTo(this,c); d.subTo(m,d); }
                c.rShiftTo(1,c);
            }
            else if(!d.isEven()) d.subTo(m,d);
            d.rShiftTo(1,d);
        }
        if(u.compareTo(v) >= 0) {
            u.subTo(v,u);
            if(ac) a.subTo(c,a);
            b.subTo(d,b);
        }
        else {
            v.subTo(u,v);
            if(ac) c.subTo(a,c);
            d.subTo(b,d);
        }
    }
    if(v.compareTo(BigInteger.ONE) != 0) return BigInteger.ZERO;
    if(d.compareTo(m) >= 0) return d.subtract(m);
    if(d.signum() < 0) d.addTo(m,d); else return d;
    if(d.signum() < 0) return d.add(m); else return d;
}

var lowprimes = [2,3,5,7,11,13,17,19,23,29,31,37,41,43,47,53,59,61,67,71,73,79,83,89,97,101,103,107,109,113,127,131,137,139,149,151,157,163,167,173,179,181,191,193,197,199,211,223,227,229,233,239,241,251,257,263,269,271,277,281,283,293,307,311,313,317,331,337,347,349,353,359,367,373,379,383,389,397,401,409,419,421,431,433,439,443,449,457,461,463,467,479,487,491,499,503,509,521,523,541,547,557,563,569,571,577,587,593,599,601,607,613,617,619,631,641,643,647,653,659,661,673,677,683,691,701,709,719,727,733,739,743,751,757,761,769,773,787,797,809,811,821,823,827,829,839,853,857,859,863,877,881,883,887,907,911,919,929,937,941,947,953,967,971,977,983,991,997];
var lplim = (1<<26)/lowprimes[lowprimes.length-1];

// (public) test primality with certainty >= 1-.5^t
function bnIsProbablePrime(t) {
    var i, x = this.abs();
    if(x.t == 1 && x[0] <= lowprimes[lowprimes.length-1]) {
        for(i = 0; i < lowprimes.length; ++i)
            if(x[0] == lowprimes[i]) return true;
        return false;
    }
    if(x.isEven()) return false;
    i = 1;
    while(i < lowprimes.length) {
        var m = lowprimes[i], j = i+1;
        while(j < lowprimes.length && m < lplim) m *= lowprimes[j++];
        m = x.modInt(m);
        while(i < j) if(m%lowprimes[i++] == 0) return false;
    }
    return x.millerRabin(t);
}

// (protected) true if probably prime (HAC 4.24, Miller-Rabin)
function bnpMillerRabin(t) {
    var n1 = this.subtract(BigInteger.ONE);
    var k = n1.getLowestSetBit();
    if(k <= 0) return false;
    var r = n1.shiftRight(k);
    t = (t+1)>>1;
    if(t > lowprimes.length) t = lowprimes.length;
    var a = nbi();
    for(var i = 0; i < t; ++i) {
        //Pick bases at random, instead of starting at 2
        a.fromInt(lowprimes[Math.floor(Math.random()*lowprimes.length)]);
        var y = a.modPow(r,this);
        if(y.compareTo(BigInteger.ONE) != 0 && y.compareTo(n1) != 0) {
            var j = 1;
            while(j++ < k && y.compareTo(n1) != 0) {
                y = y.modPowInt(2,this);
                if(y.compareTo(BigInteger.ONE) == 0) return false;
            }
            if(y.compareTo(n1) != 0) return false;
        }
    }
    return true;
}

// protected
BigInteger.prototype.chunkSize = bnpChunkSize;
BigInteger.prototype.toRadix = bnpToRadix;
BigInteger.prototype.fromRadix = bnpFromRadix;
BigInteger.prototype.fromNumber = bnpFromNumber;
BigInteger.prototype.bitwiseTo = bnpBitwiseTo;
BigInteger.prototype.changeBit = bnpChangeBit;
BigInteger.prototype.addTo = bnpAddTo;
BigInteger.prototype.dMultiply = bnpDMultiply;
BigInteger.prototype.dAddOffset = bnpDAddOffset;
BigInteger.prototype.multiplyLowerTo = bnpMultiplyLowerTo;
BigInteger.prototype.multiplyUpperTo = bnpMultiplyUpperTo;
BigInteger.prototype.modInt = bnpModInt;
BigInteger.prototype.millerRabin = bnpMillerRabin;

// public
BigInteger.prototype.clone = bnClone;
BigInteger.prototype.intValue = bnIntValue;
BigInteger.prototype.byteValue = bnByteValue;
BigInteger.prototype.shortValue = bnShortValue;
BigInteger.prototype.signum = bnSigNum;
BigInteger.prototype.toByteArray = bnToByteArray;
BigInteger.prototype.equals = bnEquals;
BigInteger.prototype.min = bnMin;
BigInteger.prototype.max = bnMax;
BigInteger.prototype.and = bnAnd;
BigInteger.prototype.or = bnOr;
BigInteger.prototype.xor = bnXor;
BigInteger.prototype.andNot = bnAndNot;
BigInteger.prototype.not = bnNot;
BigInteger.prototype.shiftLeft = bnShiftLeft;
BigInteger.prototype.shiftRight = bnShiftRight;
BigInteger.prototype.getLowestSetBit = bnGetLowestSetBit;
BigInteger.prototype.bitCount = bnBitCount;
BigInteger.prototype.testBit = bnTestBit;
BigInteger.prototype.setBit = bnSetBit;
BigInteger.prototype.clearBit = bnClearBit;
BigInteger.prototype.flipBit = bnFlipBit;
BigInteger.prototype.add = bnAdd;
BigInteger.prototype.subtract = bnSubtract;
BigInteger.prototype.multiply = bnMultiply;
BigInteger.prototype.divide = bnDivide;
BigInteger.prototype.remainder = bnRemainder;
BigInteger.prototype.divideAndRemainder = bnDivideAndRemainder;
BigInteger.prototype.modPow = bnModPow;
BigInteger.prototype.modInverse = bnModInverse;
BigInteger.prototype.pow = bnPow;
BigInteger.prototype.gcd = bnGCD;
BigInteger.prototype.isProbablePrime = bnIsProbablePrime;

// JSBN-specific extension
BigInteger.prototype.square = bnSquare;

// BigInteger interfaces not implemented in jsbn:

// BigInteger(int signum, byte[] magnitude)
// double doubleValue()
// float floatValue()
// int hashCode()
// long longValue()
// static BigInteger valueOf(long val)

// prng4.js - uses Arcfour as a PRNG

function Arcfour() {
    this.i = 0;
    this.j = 0;
    this.S = new Array();
}

// Initialize arcfour context from key, an array of ints, each from [0..255]
function ARC4init(key) {
    var i, j, t;
    for(i = 0; i < 256; ++i)
        this.S[i] = i;
    j = 0;
    for(i = 0; i < 256; ++i) {
        j = (j + this.S[i] + key[i % key.length]) & 255;
        t = this.S[i];
        this.S[i] = this.S[j];
        this.S[j] = t;
    }
    this.i = 0;
    this.j = 0;
}

function ARC4next() {
    var t;
    this.i = (this.i + 1) & 255;
    this.j = (this.j + this.S[this.i]) & 255;
    t = this.S[this.i];
    this.S[this.i] = this.S[this.j];
    this.S[this.j] = t;
    return this.S[(t + this.S[this.i]) & 255];
}

Arcfour.prototype.init = ARC4init;
Arcfour.prototype.next = ARC4next;

// Plug in your RNG constructor here
function prng_newstate() {
    return new Arcfour();
}

// Pool size must be a multiple of 4 and greater than 32.
// An array of bytes the size of the pool will be passed to init()
var rng_psize = 256;


// Random number generator - requires a PRNG backend, e.g. prng4.js

// For best results, put code like
// <body onClick='rng_seed_time();' onKeyPress='rng_seed_time();'>
// in your main HTML document.

var rng_state;
var rng_pool;
var rng_pptr;

// Mix in a 32-bit integer into the pool
function rng_seed_int(x) {
    rng_pool[rng_pptr++] ^= x & 255;
    rng_pool[rng_pptr++] ^= (x >> 8) & 255;
    rng_pool[rng_pptr++] ^= (x >> 16) & 255;
    rng_pool[rng_pptr++] ^= (x >> 24) & 255;
    if(rng_pptr >= rng_psize) rng_pptr -= rng_psize;
}

// Mix in the current time (w/milliseconds) into the pool
function rng_seed_time() {
    rng_seed_int(new Date().getTime());
}

// Initialize the pool with junk if needed.
if(rng_pool == null) {
    rng_pool = new Array();
    rng_pptr = 0;
    var t;
    if(window.crypto && window.crypto.getRandomValues) {
        // Use webcrypto if available
        var ua = new Uint8Array(32);
        window.crypto.getRandomValues(ua);
        for(t = 0; t < 32; ++t)
            rng_pool[rng_pptr++] = ua[t];
    }
    if(navigator.appName == "Netscape" && navigator.appVersion < "5" && window.crypto) {
        // Extract entropy (256 bits) from NS4 RNG if available
        var z = window.crypto.random(32);
        for(t = 0; t < z.length; ++t)
            rng_pool[rng_pptr++] = z.charCodeAt(t) & 255;
    }
    while(rng_pptr < rng_psize) {  // extract some randomness from Math.random()
        t = Math.floor(65536 * Math.random());
        rng_pool[rng_pptr++] = t >>> 8;
        rng_pool[rng_pptr++] = t & 255;
    }
    rng_pptr = 0;
    rng_seed_time();
    //rng_seed_int(window.screenX);
    //rng_seed_int(window.screenY);
}

function rng_get_byte() {
    if(rng_state == null) {
        rng_seed_time();
        rng_state = prng_newstate();
        rng_state.init(rng_pool);
        for(rng_pptr = 0; rng_pptr < rng_pool.length; ++rng_pptr)
            rng_pool[rng_pptr] = 0;
        rng_pptr = 0;
        //rng_pool = null;
    }
    // TODO: allow reseeding after first request
    return rng_state.next();
}

function rng_get_bytes(ba) {
    var i;
    for(i = 0; i < ba.length; ++i) ba[i] = rng_get_byte();
}

function SecureRandom() {}

SecureRandom.prototype.nextBytes = rng_get_bytes;


// Basic Javascript Elliptic Curve implementation
// Ported loosely from BouncyCastle's Java EC code
// Only Fp curves implemented for now

// Requires jsbn.js and jsbn2.js

// ----------------
// ECFieldElementFp

// constructor
function ECFieldElementFp(q,x) {
    this.x = x;
    // TODO if(x.compareTo(q) >= 0) error
    this.q = q;
}

function feFpEquals(other) {
    if(other == this) return true;
    return (this.q.equals(other.q) && this.x.equals(other.x));
}

function feFpToBigInteger() {
    return this.x;
}

function feFpNegate() {
    return new ECFieldElementFp(this.q, this.x.negate().mod(this.q));
}

function feFpAdd(b) {
    return new ECFieldElementFp(this.q, this.x.add(b.toBigInteger()).mod(this.q));
}

function feFpSubtract(b) {
    return new ECFieldElementFp(this.q, this.x.subtract(b.toBigInteger()).mod(this.q));
}

function feFpMultiply(b) {
    return new ECFieldElementFp(this.q, this.x.multiply(b.toBigInteger()).mod(this.q));
}

function feFpSquare() {
    return new ECFieldElementFp(this.q, this.x.square().mod(this.q));
}

function feFpDivide(b) {
    return new ECFieldElementFp(this.q, this.x.multiply(b.toBigInteger().modInverse(this.q)).mod(this.q));
}

ECFieldElementFp.prototype.equals = feFpEquals;
ECFieldElementFp.prototype.toBigInteger = feFpToBigInteger;
ECFieldElementFp.prototype.negate = feFpNegate;
ECFieldElementFp.prototype.add = feFpAdd;
ECFieldElementFp.prototype.subtract = feFpSubtract;
ECFieldElementFp.prototype.multiply = feFpMultiply;
ECFieldElementFp.prototype.square = feFpSquare;
ECFieldElementFp.prototype.divide = feFpDivide;

// ----------------
// ECPointFp

// constructor
function ECPointFp(curve,x,y,z) {
    this.curve = curve;
    this.x = x;
    this.y = y;
    // Projective coordinates: either zinv == null or z * zinv == 1
    // z and zinv are just BigIntegers, not fieldElements
    if(z == null) {
        this.z = BigInteger.ONE;
    }
    else {
        this.z = z;
    }
    this.zinv = null;
    //TODO: compression flag
}

function pointFpGetX() {
    if(this.zinv == null) {
        this.zinv = this.z.modInverse(this.curve.q);
    }
    var r = this.x.toBigInteger().multiply(this.zinv);
    this.curve.reduce(r);
    return this.curve.fromBigInteger(r);
}

function pointFpGetY() {
    if(this.zinv == null) {
        this.zinv = this.z.modInverse(this.curve.q);
    }
    var r = this.y.toBigInteger().multiply(this.zinv);
    this.curve.reduce(r);
    return this.curve.fromBigInteger(r);
}

function pointFpEquals(other) {
    if(other == this) return true;
    if(this.isInfinity()) return other.isInfinity();
    if(other.isInfinity()) return this.isInfinity();
    var u, v;
    // u = Y2 * Z1 - Y1 * Z2
    u = other.y.toBigInteger().multiply(this.z).subtract(this.y.toBigInteger().multiply(other.z)).mod(this.curve.q);
    if(!u.equals(BigInteger.ZERO)) return false;
    // v = X2 * Z1 - X1 * Z2
    v = other.x.toBigInteger().multiply(this.z).subtract(this.x.toBigInteger().multiply(other.z)).mod(this.curve.q);
    return v.equals(BigInteger.ZERO);
}

function pointFpIsInfinity() {
    if((this.x == null) && (this.y == null)) return true;
    return this.z.equals(BigInteger.ZERO) && !this.y.toBigInteger().equals(BigInteger.ZERO);
}

function pointFpNegate() {
    return new ECPointFp(this.curve, this.x, this.y.negate(), this.z);
}

function pointFpAdd(b) {
    if(this.isInfinity()) return b;
    if(b.isInfinity()) return this;

    // u = Y2 * Z1 - Y1 * Z2
    var u = b.y.toBigInteger().multiply(this.z).subtract(this.y.toBigInteger().multiply(b.z)).mod(this.curve.q);
    // v = X2 * Z1 - X1 * Z2
    var v = b.x.toBigInteger().multiply(this.z).subtract(this.x.toBigInteger().multiply(b.z)).mod(this.curve.q);

    if(BigInteger.ZERO.equals(v)) {
        if(BigInteger.ZERO.equals(u)) {
            return this.twice(); // this == b, so double
        }
        return this.curve.getInfinity(); // this = -b, so infinity
    }

    var THREE = new BigInteger("3");
    var x1 = this.x.toBigInteger();
    var y1 = this.y.toBigInteger();
    var x2 = b.x.toBigInteger();
    var y2 = b.y.toBigInteger();

    var v2 = v.square();
    var v3 = v2.multiply(v);
    var x1v2 = x1.multiply(v2);
    var zu2 = u.square().multiply(this.z);

    // x3 = v * (z2 * (z1 * u^2 - 2 * x1 * v^2) - v^3)
    var x3 = zu2.subtract(x1v2.shiftLeft(1)).multiply(b.z).subtract(v3).multiply(v).mod(this.curve.q);
    // y3 = z2 * (3 * x1 * u * v^2 - y1 * v^3 - z1 * u^3) + u * v^3
    var y3 = x1v2.multiply(THREE).multiply(u).subtract(y1.multiply(v3)).subtract(zu2.multiply(u)).multiply(b.z).add(u.multiply(v3)).mod(this.curve.q);
    // z3 = v^3 * z1 * z2
    var z3 = v3.multiply(this.z).multiply(b.z).mod(this.curve.q);

    return new ECPointFp(this.curve, this.curve.fromBigInteger(x3), this.curve.fromBigInteger(y3), z3);
}

function pointFpTwice() {
    if(this.isInfinity()) return this;
    if(this.y.toBigInteger().signum() == 0) return this.curve.getInfinity();

    // TODO: optimized handling of constants
    var THREE = new BigInteger("3");
    var x1 = this.x.toBigInteger();
    var y1 = this.y.toBigInteger();

    var y1z1 = y1.multiply(this.z);
    var y1sqz1 = y1z1.multiply(y1).mod(this.curve.q);
    var a = this.curve.a.toBigInteger();

    // w = 3 * x1^2 + a * z1^2
    var w = x1.square().multiply(THREE);
    if(!BigInteger.ZERO.equals(a)) {
        w = w.add(this.z.square().multiply(a));
    }
    w = w.mod(this.curve.q);
    //this.curve.reduce(w);
    // x3 = 2 * y1 * z1 * (w^2 - 8 * x1 * y1^2 * z1)
    var x3 = w.square().subtract(x1.shiftLeft(3).multiply(y1sqz1)).shiftLeft(1).multiply(y1z1).mod(this.curve.q);
    // y3 = 4 * y1^2 * z1 * (3 * w * x1 - 2 * y1^2 * z1) - w^3
    var y3 = w.multiply(THREE).multiply(x1).subtract(y1sqz1.shiftLeft(1)).shiftLeft(2).multiply(y1sqz1).subtract(w.square().multiply(w)).mod(this.curve.q);
    // z3 = 8 * (y1 * z1)^3
    var z3 = y1z1.square().multiply(y1z1).shiftLeft(3).mod(this.curve.q);

    return new ECPointFp(this.curve, this.curve.fromBigInteger(x3), this.curve.fromBigInteger(y3), z3);
}

// Simple NAF (Non-Adjacent Form) multiplication algorithm
// TODO: modularize the multiplication algorithm
function pointFpMultiply(k) {
    if(this.isInfinity()) return this;
    if(k.signum() == 0) return this.curve.getInfinity();

    var e = k;
    var h = e.multiply(new BigInteger("3"));

    var neg = this.negate();
    var R = this;

    var i;
    for(i = h.bitLength() - 2; i > 0; --i) {
        R = R.twice();

        var hBit = h.testBit(i);
        var eBit = e.testBit(i);

        if (hBit != eBit) {
            R = R.add(hBit ? this : neg);
        }
    }

    return R;
}

// Compute this*j + x*k (simultaneous multiplication)
function pointFpMultiplyTwo(j,x,k) {
    var i;
    if(j.bitLength() > k.bitLength())
        i = j.bitLength() - 1;
    else
        i = k.bitLength() - 1;

    var R = this.curve.getInfinity();
    var both = this.add(x);
    while(i >= 0) {
        R = R.twice();
        if(j.testBit(i)) {
            if(k.testBit(i)) {
                R = R.add(both);
            }
            else {
                R = R.add(this);
            }
        }
        else {
            if(k.testBit(i)) {
                R = R.add(x);
            }
        }
        --i;
    }

    return R;
}

ECPointFp.prototype.getX = pointFpGetX;
ECPointFp.prototype.getY = pointFpGetY;
ECPointFp.prototype.equals = pointFpEquals;
ECPointFp.prototype.isInfinity = pointFpIsInfinity;
ECPointFp.prototype.negate = pointFpNegate;
ECPointFp.prototype.add = pointFpAdd;
ECPointFp.prototype.twice = pointFpTwice;
ECPointFp.prototype.multiply = pointFpMultiply;
ECPointFp.prototype.multiplyTwo = pointFpMultiplyTwo;

// ----------------
// ECCurveFp

// constructor
function ECCurveFp(q,a,b) {
    this.q = q;
    this.a = this.fromBigInteger(a);
    this.b = this.fromBigInteger(b);
    this.infinity = new ECPointFp(this, null, null);
    this.reducer = new Barrett(this.q);
}

function curveFpGetQ() {
    return this.q;
}

function curveFpGetA() {
    return this.a;
}

function curveFpGetB() {
    return this.b;
}

function curveFpEquals(other) {
    if(other == this) return true;
    return(this.q.equals(other.q) && this.a.equals(other.a) && this.b.equals(other.b));
}

function curveFpGetInfinity() {
    return this.infinity;
}

function curveFpFromBigInteger(x) {
    return new ECFieldElementFp(this.q, x);
}

function curveReduce(x) {
    this.reducer.reduce(x);
}

// for now, work with hex strings because they're easier in JS
function curveFpDecodePointHex(s) {
    switch(parseInt(s.substr(0,2), 16)) { // first byte
        case 0:
            return this.infinity;
        case 2:
        case 3:
            // point compression not supported yet
            return null;
        case 4:
        case 6:
        case 7:
            var len = (s.length - 2) / 2;
            var xHex = s.substr(2, len);
            var yHex = s.substr(len+2, len);

            return new ECPointFp(this,
                this.fromBigInteger(new BigInteger(xHex, 16)),
                this.fromBigInteger(new BigInteger(yHex, 16)));

        default: // unsupported
            return null;
    }
}

function curveFpEncodePointHex(p) {
    if (p.isInfinity()) return "00";
    var xHex = p.getX().toBigInteger().toString(16);
    var yHex = p.getY().toBigInteger().toString(16);
    var oLen = this.getQ().toString(16).length;
    if ((oLen % 2) != 0) oLen++;
    while (xHex.length < oLen) {
        xHex = "0" + xHex;
    }
    while (yHex.length < oLen) {
        yHex = "0" + yHex;
    }
    return "04" + xHex + yHex;
}

ECCurveFp.prototype.getQ = curveFpGetQ;
ECCurveFp.prototype.getA = curveFpGetA;
ECCurveFp.prototype.getB = curveFpGetB;
ECCurveFp.prototype.equals = curveFpEquals;
ECCurveFp.prototype.getInfinity = curveFpGetInfinity;
ECCurveFp.prototype.fromBigInteger = curveFpFromBigInteger;
ECCurveFp.prototype.reduce = curveReduce;
ECCurveFp.prototype.decodePointHex = curveFpDecodePointHex;
ECCurveFp.prototype.encodePointHex = curveFpEncodePointHex;


// Named EC curves

// Requires ec.js, jsbn.js, and jsbn2.js

// ----------------
// X9ECParameters

// constructor
function X9ECParameters(curve,g,n,h) {
    this.curve = curve;
    this.g = g;
    this.n = n;
    this.h = h;
}

function x9getCurve() {
    return this.curve;
}

function x9getG() {
    return this.g;
}

function x9getN() {
    return this.n;
}

function x9getH() {
    return this.h;
}

X9ECParameters.prototype.getCurve = x9getCurve;
X9ECParameters.prototype.getG = x9getG;
X9ECParameters.prototype.getN = x9getN;
X9ECParameters.prototype.getH = x9getH;

// ----------------
// SECNamedCurves

function fromHex(s) { return new BigInteger(s, 16); }

function secp128r1() {
    // p = 2^128 - 2^97 - 1
    var p = fromHex("FFFFFFFDFFFFFFFFFFFFFFFFFFFFFFFF");
    var a = fromHex("FFFFFFFDFFFFFFFFFFFFFFFFFFFFFFFC");
    var b = fromHex("E87579C11079F43DD824993C2CEE5ED3");
    //byte[] S = Hex.decode("000E0D4D696E6768756151750CC03A4473D03679");
    var n = fromHex("FFFFFFFE0000000075A30D1B9038A115");
    var h = BigInteger.ONE;
    var curve = new ECCurveFp(p, a, b);
    var G = curve.decodePointHex("04"
        + "161FF7528B899B2D0C28607CA52C5B86"
        + "CF5AC8395BAFEB13C02DA292DDED7A83");
    return new X9ECParameters(curve, G, n, h);
}

function secp160k1() {
    // p = 2^160 - 2^32 - 2^14 - 2^12 - 2^9 - 2^8 - 2^7 - 2^3 - 2^2 - 1
    var p = fromHex("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFAC73");
    var a = BigInteger.ZERO;
    var b = fromHex("7");
    //byte[] S = null;
    var n = fromHex("0100000000000000000001B8FA16DFAB9ACA16B6B3");
    var h = BigInteger.ONE;
    var curve = new ECCurveFp(p, a, b);
    var G = curve.decodePointHex("04"
        + "3B4C382CE37AA192A4019E763036F4F5DD4D7EBB"
        + "938CF935318FDCED6BC28286531733C3F03C4FEE");
    return new X9ECParameters(curve, G, n, h);
}

function secp160r1() {
    // p = 2^160 - 2^31 - 1
    var p = fromHex("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF7FFFFFFF");
    var a = fromHex("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF7FFFFFFC");
    var b = fromHex("1C97BEFC54BD7A8B65ACF89F81D4D4ADC565FA45");
    //byte[] S = Hex.decode("1053CDE42C14D696E67687561517533BF3F83345");
    var n = fromHex("0100000000000000000001F4C8F927AED3CA752257");
    var h = BigInteger.ONE;
    var curve = new ECCurveFp(p, a, b);
    var G = curve.decodePointHex("04"
        + "4A96B5688EF573284664698968C38BB913CBFC82"
        + "23A628553168947D59DCC912042351377AC5FB32");
    return new X9ECParameters(curve, G, n, h);
}

function secp192k1() {
    // p = 2^192 - 2^32 - 2^12 - 2^8 - 2^7 - 2^6 - 2^3 - 1
    var p = fromHex("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFEE37");
    var a = BigInteger.ZERO;
    var b = fromHex("3");
    //byte[] S = null;
    var n = fromHex("FFFFFFFFFFFFFFFFFFFFFFFE26F2FC170F69466A74DEFD8D");
    var h = BigInteger.ONE;
    var curve = new ECCurveFp(p, a, b);
    var G = curve.decodePointHex("04"
        + "DB4FF10EC057E9AE26B07D0280B7F4341DA5D1B1EAE06C7D"
        + "9B2F2F6D9C5628A7844163D015BE86344082AA88D95E2F9D");
    return new X9ECParameters(curve, G, n, h);
}

function secp192r1() {
    // p = 2^192 - 2^64 - 1
    var p = fromHex("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFFFFFFFFFFFF");
    var a = fromHex("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFFFFFFFFFFFC");
    var b = fromHex("64210519E59C80E70FA7E9AB72243049FEB8DEECC146B9B1");
    //byte[] S = Hex.decode("3045AE6FC8422F64ED579528D38120EAE12196D5");
    var n = fromHex("FFFFFFFFFFFFFFFFFFFFFFFF99DEF836146BC9B1B4D22831");
    var h = BigInteger.ONE;
    var curve = new ECCurveFp(p, a, b);
    var G = curve.decodePointHex("04"
        + "188DA80EB03090F67CBF20EB43A18800F4FF0AFD82FF1012"
        + "07192B95FFC8DA78631011ED6B24CDD573F977A11E794811");
    return new X9ECParameters(curve, G, n, h);
}

function secp224r1() {
    // p = 2^224 - 2^96 + 1
    var p = fromHex("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF000000000000000000000001");
    var a = fromHex("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFE");
    var b = fromHex("B4050A850C04B3ABF54132565044B0B7D7BFD8BA270B39432355FFB4");
    //byte[] S = Hex.decode("BD71344799D5C7FCDC45B59FA3B9AB8F6A948BC5");
    var n = fromHex("FFFFFFFFFFFFFFFFFFFFFFFFFFFF16A2E0B8F03E13DD29455C5C2A3D");
    var h = BigInteger.ONE;
    var curve = new ECCurveFp(p, a, b);
    var G = curve.decodePointHex("04"
        + "B70E0CBD6BB4BF7F321390B94A03C1D356C21122343280D6115C1D21"
        + "BD376388B5F723FB4C22DFE6CD4375A05A07476444D5819985007E34");
    return new X9ECParameters(curve, G, n, h);
}

function secp256r1() {
    // p = 2^224 (2^32 - 1) + 2^192 + 2^96 - 1
    var p = fromHex("FFFFFFFF00000001000000000000000000000000FFFFFFFFFFFFFFFFFFFFFFFF");
    var a = fromHex("FFFFFFFF00000001000000000000000000000000FFFFFFFFFFFFFFFFFFFFFFFC");
    var b = fromHex("5AC635D8AA3A93E7B3EBBD55769886BC651D06B0CC53B0F63BCE3C3E27D2604B");
    //byte[] S = Hex.decode("C49D360886E704936A6678E1139D26B7819F7E90");
    var n = fromHex("FFFFFFFF00000000FFFFFFFFFFFFFFFFBCE6FAADA7179E84F3B9CAC2FC632551");
    var h = BigInteger.ONE;
    var curve = new ECCurveFp(p, a, b);
    var G = curve.decodePointHex("04"
        + "6B17D1F2E12C4247F8BCE6E563A440F277037D812DEB33A0F4A13945D898C296"
        + "4FE342E2FE1A7F9B8EE7EB4A7C0F9E162BCE33576B315ECECBB6406837BF51F5");
    return new X9ECParameters(curve, G, n, h);
}

// TODO: make this into a proper hashtable
function getSECCurveByName(name) {
    if(name == "secp128r1") return secp128r1();
    if(name == "secp160k1") return secp160k1();
    if(name == "secp160r1") return secp160r1();
    if(name == "secp192k1") return secp192k1();
    if(name == "secp192r1") return secp192r1();
    if(name == "secp224r1") return secp224r1();
    if(name == "secp256r1") return secp256r1();
    return null;
}


/**
 *
 * MD5 (Message-Digest Algorithm) http://www.webtoolkit.info/
 *
 */

window.MD5 = function(string) {

    function RotateLeft(lValue, iShiftBits) {
        return (lValue << iShiftBits) | (lValue >>> (32 - iShiftBits));
    }

    function AddUnsigned(lX, lY) {
        var lX4, lY4, lX8, lY8, lResult;
        lX8 = (lX & 0x80000000);
        lY8 = (lY & 0x80000000);
        lX4 = (lX & 0x40000000);
        lY4 = (lY & 0x40000000);
        lResult = (lX & 0x3FFFFFFF) + (lY & 0x3FFFFFFF);
        if (lX4 & lY4) {
            return (lResult ^ 0x80000000 ^ lX8 ^ lY8);
        }
        if (lX4 | lY4) {
            if (lResult & 0x40000000) {
                return (lResult ^ 0xC0000000 ^ lX8 ^ lY8);
            } else {
                return (lResult ^ 0x40000000 ^ lX8 ^ lY8);
            }
        } else {
            return (lResult ^ lX8 ^ lY8);
        }
    }

    function F(x, y, z) {
        return (x & y) | ((~x) & z);
    }
    function G(x, y, z) {
        return (x & z) | (y & (~z));
    }
    function H(x, y, z) {
        return (x ^ y ^ z);
    }
    function I(x, y, z) {
        return (y ^ (x | (~z)));
    }

    function FF(a, b, c, d, x, s, ac) {
        a = AddUnsigned(a, AddUnsigned(AddUnsigned(F(b, c, d), x), ac));
        return AddUnsigned(RotateLeft(a, s), b);
    }
    ;

    function GG(a, b, c, d, x, s, ac) {
        a = AddUnsigned(a, AddUnsigned(AddUnsigned(G(b, c, d), x), ac));
        return AddUnsigned(RotateLeft(a, s), b);
    }
    ;

    function HH(a, b, c, d, x, s, ac) {
        a = AddUnsigned(a, AddUnsigned(AddUnsigned(H(b, c, d), x), ac));
        return AddUnsigned(RotateLeft(a, s), b);
    }
    ;

    function II(a, b, c, d, x, s, ac) {
        a = AddUnsigned(a, AddUnsigned(AddUnsigned(I(b, c, d), x), ac));
        return AddUnsigned(RotateLeft(a, s), b);
    }
    ;

    function ConvertToWordArray(string) {
        var lWordCount;
        var lMessageLength = string.length;
        var lNumberOfWords_temp1 = lMessageLength + 8;
        var lNumberOfWords_temp2 = (lNumberOfWords_temp1 - (lNumberOfWords_temp1 % 64)) / 64;
        var lNumberOfWords = (lNumberOfWords_temp2 + 1) * 16;
        var lWordArray = Array(lNumberOfWords - 1);
        var lBytePosition = 0;
        var lByteCount = 0;
        while (lByteCount < lMessageLength) {
            lWordCount = (lByteCount - (lByteCount % 4)) / 4;
            lBytePosition = (lByteCount % 4) * 8;
            lWordArray[lWordCount] = (lWordArray[lWordCount] | (string.charCodeAt(lByteCount) << lBytePosition));
            lByteCount++;
        }
        lWordCount = (lByteCount - (lByteCount % 4)) / 4;
        lBytePosition = (lByteCount % 4) * 8;
        lWordArray[lWordCount] = lWordArray[lWordCount] | (0x80 << lBytePosition);
        lWordArray[lNumberOfWords - 2] = lMessageLength << 3;
        lWordArray[lNumberOfWords - 1] = lMessageLength >>> 29;
        return lWordArray;
    }
    ;

    function WordToHex(lValue) {
        var WordToHexValue = "", WordToHexValue_temp = "", lByte, lCount;
        for (lCount = 0; lCount <= 3; lCount++) {
            lByte = (lValue >>> (lCount * 8)) & 255;
            WordToHexValue_temp = "0" + lByte.toString(16);
            WordToHexValue = WordToHexValue + WordToHexValue_temp.substr(WordToHexValue_temp.length - 2, 2);
        }
        return WordToHexValue;
    }
    ;

    function Utf8Encode(string) {
        string = string.replace(/\r\n/g, "\n");
        var utftext = "";

        for (var n = 0; n < string.length; n++) {

            var c = string.charCodeAt(n);

            if (c < 128) {
                utftext += String.fromCharCode(c);
            } else if ((c > 127) && (c < 2048)) {
                utftext += String.fromCharCode((c >> 6) | 192);
                utftext += String.fromCharCode((c & 63) | 128);
            } else {
                utftext += String.fromCharCode((c >> 12) | 224);
                utftext += String.fromCharCode(((c >> 6) & 63) | 128);
                utftext += String.fromCharCode((c & 63) | 128);
            }

        }

        return utftext;
    }
    ;

    var x = Array();
    var k, AA, BB, CC, DD, a, b, c, d;
    var S11 = 7, S12 = 12, S13 = 17, S14 = 22;
    var S21 = 5, S22 = 9, S23 = 14, S24 = 20;
    var S31 = 4, S32 = 11, S33 = 16, S34 = 23;
    var S41 = 6, S42 = 10, S43 = 15, S44 = 21;

    string = Utf8Encode(string);

    x = ConvertToWordArray(string);

    a = 0x67452301;
    b = 0xEFCDAB89;
    c = 0x98BADCFE;
    d = 0x10325476;

    for (k = 0; k < x.length; k += 16) {
        AA = a;
        BB = b;
        CC = c;
        DD = d;
        a = FF(a, b, c, d, x[k + 0], S11, 0xD76AA478);
        d = FF(d, a, b, c, x[k + 1], S12, 0xE8C7B756);
        c = FF(c, d, a, b, x[k + 2], S13, 0x242070DB);
        b = FF(b, c, d, a, x[k + 3], S14, 0xC1BDCEEE);
        a = FF(a, b, c, d, x[k + 4], S11, 0xF57C0FAF);
        d = FF(d, a, b, c, x[k + 5], S12, 0x4787C62A);
        c = FF(c, d, a, b, x[k + 6], S13, 0xA8304613);
        b = FF(b, c, d, a, x[k + 7], S14, 0xFD469501);
        a = FF(a, b, c, d, x[k + 8], S11, 0x698098D8);
        d = FF(d, a, b, c, x[k + 9], S12, 0x8B44F7AF);
        c = FF(c, d, a, b, x[k + 10], S13, 0xFFFF5BB1);
        b = FF(b, c, d, a, x[k + 11], S14, 0x895CD7BE);
        a = FF(a, b, c, d, x[k + 12], S11, 0x6B901122);
        d = FF(d, a, b, c, x[k + 13], S12, 0xFD987193);
        c = FF(c, d, a, b, x[k + 14], S13, 0xA679438E);
        b = FF(b, c, d, a, x[k + 15], S14, 0x49B40821);
        a = GG(a, b, c, d, x[k + 1], S21, 0xF61E2562);
        d = GG(d, a, b, c, x[k + 6], S22, 0xC040B340);
        c = GG(c, d, a, b, x[k + 11], S23, 0x265E5A51);
        b = GG(b, c, d, a, x[k + 0], S24, 0xE9B6C7AA);
        a = GG(a, b, c, d, x[k + 5], S21, 0xD62F105D);
        d = GG(d, a, b, c, x[k + 10], S22, 0x2441453);
        c = GG(c, d, a, b, x[k + 15], S23, 0xD8A1E681);
        b = GG(b, c, d, a, x[k + 4], S24, 0xE7D3FBC8);
        a = GG(a, b, c, d, x[k + 9], S21, 0x21E1CDE6);
        d = GG(d, a, b, c, x[k + 14], S22, 0xC33707D6);
        c = GG(c, d, a, b, x[k + 3], S23, 0xF4D50D87);
        b = GG(b, c, d, a, x[k + 8], S24, 0x455A14ED);
        a = GG(a, b, c, d, x[k + 13], S21, 0xA9E3E905);
        d = GG(d, a, b, c, x[k + 2], S22, 0xFCEFA3F8);
        c = GG(c, d, a, b, x[k + 7], S23, 0x676F02D9);
        b = GG(b, c, d, a, x[k + 12], S24, 0x8D2A4C8A);
        a = HH(a, b, c, d, x[k + 5], S31, 0xFFFA3942);
        d = HH(d, a, b, c, x[k + 8], S32, 0x8771F681);
        c = HH(c, d, a, b, x[k + 11], S33, 0x6D9D6122);
        b = HH(b, c, d, a, x[k + 14], S34, 0xFDE5380C);
        a = HH(a, b, c, d, x[k + 1], S31, 0xA4BEEA44);
        d = HH(d, a, b, c, x[k + 4], S32, 0x4BDECFA9);
        c = HH(c, d, a, b, x[k + 7], S33, 0xF6BB4B60);
        b = HH(b, c, d, a, x[k + 10], S34, 0xBEBFBC70);
        a = HH(a, b, c, d, x[k + 13], S31, 0x289B7EC6);
        d = HH(d, a, b, c, x[k + 0], S32, 0xEAA127FA);
        c = HH(c, d, a, b, x[k + 3], S33, 0xD4EF3085);
        b = HH(b, c, d, a, x[k + 6], S34, 0x4881D05);
        a = HH(a, b, c, d, x[k + 9], S31, 0xD9D4D039);
        d = HH(d, a, b, c, x[k + 12], S32, 0xE6DB99E5);
        c = HH(c, d, a, b, x[k + 15], S33, 0x1FA27CF8);
        b = HH(b, c, d, a, x[k + 2], S34, 0xC4AC5665);
        a = II(a, b, c, d, x[k + 0], S41, 0xF4292244);
        d = II(d, a, b, c, x[k + 7], S42, 0x432AFF97);
        c = II(c, d, a, b, x[k + 14], S43, 0xAB9423A7);
        b = II(b, c, d, a, x[k + 5], S44, 0xFC93A039);
        a = II(a, b, c, d, x[k + 12], S41, 0x655B59C3);
        d = II(d, a, b, c, x[k + 3], S42, 0x8F0CCC92);
        c = II(c, d, a, b, x[k + 10], S43, 0xFFEFF47D);
        b = II(b, c, d, a, x[k + 1], S44, 0x85845DD1);
        a = II(a, b, c, d, x[k + 8], S41, 0x6FA87E4F);
        d = II(d, a, b, c, x[k + 15], S42, 0xFE2CE6E0);
        c = II(c, d, a, b, x[k + 6], S43, 0xA3014314);
        b = II(b, c, d, a, x[k + 13], S44, 0x4E0811A1);
        a = II(a, b, c, d, x[k + 4], S41, 0xF7537E82);
        d = II(d, a, b, c, x[k + 11], S42, 0xBD3AF235);
        c = II(c, d, a, b, x[k + 2], S43, 0x2AD7D2BB);
        b = II(b, c, d, a, x[k + 9], S44, 0xEB86D391);
        a = AddUnsigned(a, AA);
        b = AddUnsigned(b, BB);
        c = AddUnsigned(c, CC);
        d = AddUnsigned(d, DD);
    }

    var temp = WordToHex(a) + WordToHex(b) + WordToHex(c) + WordToHex(d);

    return temp.toLowerCase();
};

/*
 * A JavaScript implementation of the Secure Hash Algorithm, SHA-512, as defined
 * in FIPS 180-2
 * Version 2.2 Copyright Anonymous Contributor, Paul Johnston 2000 - 2009.
 * Other contributors: Greg Holt, Andrew Kepert, Ydnar, Lostinet
 * Distributed under the BSD License
 * See http://pajhome.org.uk/crypt/md5 for details.
 * http://www.sharejs.com/codes
 */

/*
 * Configurable variables. You may need to tweak these to be compatible with
 * the server-side, but the defaults work in most cases.
 */
var hexcase = 0; /* hex output format. 0 - lowercase; 1 - uppercase */
var b64pad = ""; /* base-64 pad character. "=" for strict RFC compliance */

/*
 * These are the functions you'll usually want to call They take string arguments and return either hex or base-64 encoded strings
 */
function hex_sha512(s) {
    return rstr2hex(rstr_sha512(str2rstr_utf8(s)));
}
function b64_sha512(s) {
    return rstr2b64(rstr_sha512(str2rstr_utf8(s)));
}
function any_sha512(s, e) {
    return rstr2any(rstr_sha512(str2rstr_utf8(s)), e);
}
function hex_hmac_sha512(k, d) {
    return rstr2hex(rstr_hmac_sha512(str2rstr_utf8(k), str2rstr_utf8(d)));
}
function b64_hmac_sha512(k, d) {
    return rstr2b64(rstr_hmac_sha512(str2rstr_utf8(k), str2rstr_utf8(d)));
}
function any_hmac_sha512(k, d, e) {
    return rstr2any(rstr_hmac_sha512(str2rstr_utf8(k), str2rstr_utf8(d)), e);
}

/*
 * Perform a simple self-test to see if the VM is working
 */
function sha512_vm_test() {
    return hex_sha512("abc").toLowerCase() == "ddaf35a193617abacc417349ae20413112e6fa4e89a97ea20a9eeee64b55d39a" + "2192992a274fc1a836ba3c23a3feebbd454d4423643ce80e2a9ac94fa54ca49f";
}

/*
 * Calculate the SHA-512 of a raw string
 */
function rstr_sha512(s) {
    return binb2rstr(binb_sha512(rstr2binb(s), s.length * 8));
}

/*
 * Calculate the HMAC-SHA-512 of a key and some data (raw strings)
 */
function rstr_hmac_sha512(key, data) {
    var bkey = rstr2binb(key);
    if (bkey.length > 32)
        bkey = binb_sha512(bkey, key.length * 8);

    var ipad = Array(32), opad = Array(32);
    for (var i = 0; i < 32; i++) {
        ipad[i] = bkey[i] ^ 0x36363636;
        opad[i] = bkey[i] ^ 0x5C5C5C5C;
    }

    var hash = binb_sha512(ipad.concat(rstr2binb(data)), 1024 + data.length * 8);
    return binb2rstr(binb_sha512(opad.concat(hash), 1024 + 512));
}

/*
 * Convert a raw string to a hex string
 */
function rstr2hex(input) {
    try {
        hexcase
    } catch (e) {
        hexcase = 0;
    }
    var hex_tab = hexcase ? "0123456789ABCDEF" : "0123456789abcdef";
    var output = "";
    var x;
    for (var i = 0; i < input.length; i++) {
        x = input.charCodeAt(i);
        output += hex_tab.charAt((x >>> 4) & 0x0F) + hex_tab.charAt(x & 0x0F);
    }
    return output;
}

/*
 * Convert a raw string to a base-64 string
 */
function rstr2b64(input) {
    try {
        b64pad
    } catch (e) {
        b64pad = '';
    }
    var tab = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";
    var output = "";
    var len = input.length;
    for (var i = 0; i < len; i += 3) {
        var triplet = (input.charCodeAt(i) << 16) | (i + 1 < len ? input.charCodeAt(i + 1) << 8 : 0) | (i + 2 < len ? input.charCodeAt(i + 2) : 0);
        for (var j = 0; j < 4; j++) {
            if (i * 8 + j * 6 > input.length * 8)
                output += b64pad;
            else
                output += tab.charAt((triplet >>> 6 * (3 - j)) & 0x3F);
        }
    }
    return output;
}

/*
 * Convert a raw string to an arbitrary string encoding
 */
function rstr2any(input, encoding) {
    var divisor = encoding.length;
    var i, j, q, x, quotient;

    /* Convert to an array of 16-bit big-endian values, forming the dividend */
    var dividend = Array(Math.ceil(input.length / 2));
    for (i = 0; i < dividend.length; i++) {
        dividend[i] = (input.charCodeAt(i * 2) << 8) | input.charCodeAt(i * 2 + 1);
    }

    /*
     * Repeatedly perform a long division. The binary array forms the dividend, the length of the encoding is the divisor. Once computed, the quotient forms the dividend for the next step. All remainders are stored for later use.
     */
    var full_length = Math.ceil(input.length * 8 / (Math.log(encoding.length) / Math.log(2)));
    var remainders = Array(full_length);
    for (j = 0; j < full_length; j++) {
        quotient = Array();
        x = 0;
        for (i = 0; i < dividend.length; i++) {
            x = (x << 16) + dividend[i];
            q = Math.floor(x / divisor);
            x -= q * divisor;
            if (quotient.length > 0 || q > 0)
                quotient[quotient.length] = q;
        }
        remainders[j] = x;
        dividend = quotient;
    }

    /* Convert the remainders to the output string */
    var output = "";
    for (i = remainders.length - 1; i >= 0; i--)
        output += encoding.charAt(remainders[i]);

    return output;
}

/*
 * Encode a string as utf-8. For efficiency, this assumes the input is valid utf-16.
 */
function str2rstr_utf8(input) {
    var output = "";
    var i = -1;
    var x, y;

    while (++i < input.length) {
        /* Decode utf-16 surrogate pairs */
        x = input.charCodeAt(i);
        y = i + 1 < input.length ? input.charCodeAt(i + 1) : 0;
        if (0xD800 <= x && x <= 0xDBFF && 0xDC00 <= y && y <= 0xDFFF) {
            x = 0x10000 + ((x & 0x03FF) << 10) + (y & 0x03FF);
            i++;
        }

        /* Encode output as utf-8 */
        if (x <= 0x7F)
            output += String.fromCharCode(x);
        else if (x <= 0x7FF)
            output += String.fromCharCode(0xC0 | ((x >>> 6) & 0x1F), 0x80 | (x & 0x3F));
        else if (x <= 0xFFFF)
            output += String.fromCharCode(0xE0 | ((x >>> 12) & 0x0F), 0x80 | ((x >>> 6) & 0x3F), 0x80 | (x & 0x3F));
        else if (x <= 0x1FFFFF)
            output += String.fromCharCode(0xF0 | ((x >>> 18) & 0x07), 0x80 | ((x >>> 12) & 0x3F), 0x80 | ((x >>> 6) & 0x3F), 0x80 | (x & 0x3F));
    }
    return output;
}

/*
 * Encode a string as utf-16
 */
function str2rstr_utf16le(input) {
    var output = "";
    for (var i = 0; i < input.length; i++)
        output += String.fromCharCode(input.charCodeAt(i) & 0xFF, (input.charCodeAt(i) >>> 8) & 0xFF);
    return output;
}

function str2rstr_utf16be(input) {
    var output = "";
    for (var i = 0; i < input.length; i++)
        output += String.fromCharCode((input.charCodeAt(i) >>> 8) & 0xFF, input.charCodeAt(i) & 0xFF);
    return output;
}

/*
 * Convert a raw string to an array of big-endian words Characters >255 have their high-byte silently ignored.
 */
function rstr2binb(input) {
    var output = Array(input.length >> 2);
    for (var i = 0; i < output.length; i++)
        output[i] = 0;
    for (var i = 0; i < input.length * 8; i += 8)
        output[i >> 5] |= (input.charCodeAt(i / 8) & 0xFF) << (24 - i % 32);
    return output;
}

/*
 * Convert an array of big-endian words to a string
 */
function binb2rstr(input) {
    var output = "";
    for (var i = 0; i < input.length * 32; i += 8)
        output += String.fromCharCode((input[i >> 5] >>> (24 - i % 32)) & 0xFF);
    return output;
}

/*
 * Calculate the SHA-512 of an array of big-endian dwords, and a bit length
 */
var sha512_k;
function binb_sha512(x, len) {
    if (sha512_k == undefined) {
        // SHA512 constants
        sha512_k = new Array(new int64(0x428a2f98, -685199838), new int64(0x71374491, 0x23ef65cd), new int64(-1245643825, -330482897), new int64(-373957723, -2121671748), new int64(0x3956c25b, -213338824), new int64(0x59f111f1, -1241133031), new int64(-1841331548, -1357295717), new int64(-1424204075, -630357736), new int64(-670586216, -1560083902), new int64(0x12835b01, 0x45706fbe), new int64(0x243185be, 0x4ee4b28c), new int64(0x550c7dc3, -704662302), new int64(0x72be5d74, -226784913), new int64(-2132889090, 0x3b1696b1), new int64(-1680079193, 0x25c71235), new int64(-1046744716, -815192428), new int64(-459576895, -1628353838), new int64(-272742522, 0x384f25e3), new int64(0xfc19dc6, -1953704523), new int64(0x240ca1cc, 0x77ac9c65), new int64(0x2de92c6f, 0x592b0275), new int64(0x4a7484aa,
            0x6ea6e483), new int64(0x5cb0a9dc, -1119749164), new int64(0x76f988da, -2096016459), new int64(-1740746414, -295247957), new int64(-1473132947, 0x2db43210), new int64(-1341970488, -1728372417), new int64(-1084653625, -1091629340), new int64(-958395405, 0x3da88fc2), new int64(-710438585, -1828018395), new int64(0x6ca6351, -536640913), new int64(0x14292967, 0xa0e6e70), new int64(0x27b70a85, 0x46d22ffc), new int64(0x2e1b2138, 0x5c26c926), new int64(0x4d2c6dfc, 0x5ac42aed), new int64(0x53380d13, -1651133473), new int64(0x650a7354, -1951439906), new int64(0x766a0abb, 0x3c77b2a8), new int64(-2117940946, 0x47edaee6), new int64(-1838011259, 0x1482353b), new int64(-1564481375, 0x4cf10364), new int64(-1474664885, -1136513023), new int64(-1035236496, -789014639), new int64(
            -949202525, 0x654be30), new int64(-778901479, -688958952), new int64(-694614492, 0x5565a910), new int64(-200395387, 0x5771202a), new int64(0x106aa070, 0x32bbd1b8), new int64(0x19a4c116, -1194143544), new int64(0x1e376c08, 0x5141ab53), new int64(0x2748774c, -544281703), new int64(0x34b0bcb5, -509917016), new int64(0x391c0cb3, -976659869), new int64(0x4ed8aa4a, -482243893), new int64(0x5b9cca4f, 0x7763e373), new int64(0x682e6ff3, -692930397), new int64(0x748f82ee, 0x5defb2fc), new int64(0x78a5636f, 0x43172f60), new int64(-2067236844, -1578062990), new int64(-1933114872, 0x1a6439ec), new int64(-1866530822, 0x23631e28), new int64(-1538233109, -561857047), new int64(-1090935817, -1295615723), new int64(-965641998, -479046869), new int64(-903397682, -366583396), new int64(
            -779700025, 0x21c0c207), new int64(-354779690, -840897762), new int64(-176337025, -294727304), new int64(0x6f067aa, 0x72176fba), new int64(0xa637dc5, -1563912026), new int64(0x113f9804, -1090974290), new int64(0x1b710b35, 0x131c471b), new int64(0x28db77f5, 0x23047d84), new int64(0x32caab7b, 0x40c72493), new int64(0x3c9ebe0a, 0x15c9bebc), new int64(0x431d67c4, -1676669620), new int64(0x4cc5d4be, -885112138), new int64(0x597f299c, -60457430), new int64(0x5fcb6fab, 0x3ad6faec), new int64(0x6c44198c, 0x4a475817));
    }

    // Initial hash values
    var H = new Array(new int64(0x6a09e667, -205731576), new int64(-1150833019, -2067093701), new int64(0x3c6ef372, -23791573), new int64(-1521486534, 0x5f1d36f1), new int64(0x510e527f, -1377402159), new int64(-1694144372, 0x2b3e6c1f), new int64(0x1f83d9ab, -79577749), new int64(0x5be0cd19, 0x137e2179));

    var T1 = new int64(0, 0), T2 = new int64(0, 0), a = new int64(0, 0), b = new int64(0, 0), c = new int64(0, 0), d = new int64(0, 0), e = new int64(0, 0), f = new int64(0, 0), g = new int64(0, 0), h = new int64(0, 0),
        // Temporary variables not specified by the document
        s0 = new int64(0, 0), s1 = new int64(0, 0), Ch = new int64(0, 0), Maj = new int64(0, 0), r1 = new int64(0, 0), r2 = new int64(0, 0), r3 = new int64(0, 0);
    var j, i;
    var W = new Array(80);
    for (i = 0; i < 80; i++)
        W[i] = new int64(0, 0);

    // append padding to the source string. The format is described in the FIPS.
    x[len >> 5] |= 0x80 << (24 - (len & 0x1f));
    x[((len + 128 >> 10) << 5) + 31] = len;

    for (i = 0; i < x.length; i += 32) // 32 dwords is the block size
    {
        int64copy(a, H[0]);
        int64copy(b, H[1]);
        int64copy(c, H[2]);
        int64copy(d, H[3]);
        int64copy(e, H[4]);
        int64copy(f, H[5]);
        int64copy(g, H[6]);
        int64copy(h, H[7]);

        for (j = 0; j < 16; j++) {
            W[j].h = x[i + 2 * j];
            W[j].l = x[i + 2 * j + 1];
        }

        for (j = 16; j < 80; j++) {
            // sigma1
            int64rrot(r1, W[j - 2], 19);
            int64revrrot(r2, W[j - 2], 29);
            int64shr(r3, W[j - 2], 6);
            s1.l = r1.l ^ r2.l ^ r3.l;
            s1.h = r1.h ^ r2.h ^ r3.h;
            // sigma0
            int64rrot(r1, W[j - 15], 1);
            int64rrot(r2, W[j - 15], 8);
            int64shr(r3, W[j - 15], 7);
            s0.l = r1.l ^ r2.l ^ r3.l;
            s0.h = r1.h ^ r2.h ^ r3.h;

            int64add4(W[j], s1, W[j - 7], s0, W[j - 16]);
        }

        for (j = 0; j < 80; j++) {
            // Ch
            Ch.l = (e.l & f.l) ^ (~e.l & g.l);
            Ch.h = (e.h & f.h) ^ (~e.h & g.h);

            // Sigma1
            int64rrot(r1, e, 14);
            int64rrot(r2, e, 18);
            int64revrrot(r3, e, 9);
            s1.l = r1.l ^ r2.l ^ r3.l;
            s1.h = r1.h ^ r2.h ^ r3.h;

            // Sigma0
            int64rrot(r1, a, 28);
            int64revrrot(r2, a, 2);
            int64revrrot(r3, a, 7);
            s0.l = r1.l ^ r2.l ^ r3.l;
            s0.h = r1.h ^ r2.h ^ r3.h;

            // Maj
            Maj.l = (a.l & b.l) ^ (a.l & c.l) ^ (b.l & c.l);
            Maj.h = (a.h & b.h) ^ (a.h & c.h) ^ (b.h & c.h);

            int64add5(T1, h, s1, Ch, sha512_k[j], W[j]);
            int64add(T2, s0, Maj);

            int64copy(h, g);
            int64copy(g, f);
            int64copy(f, e);
            int64add(e, d, T1);
            int64copy(d, c);
            int64copy(c, b);
            int64copy(b, a);
            int64add(a, T1, T2);
        }
        int64add(H[0], H[0], a);
        int64add(H[1], H[1], b);
        int64add(H[2], H[2], c);
        int64add(H[3], H[3], d);
        int64add(H[4], H[4], e);
        int64add(H[5], H[5], f);
        int64add(H[6], H[6], g);
        int64add(H[7], H[7], h);
    }

    // represent the hash as an array of 32-bit dwords
    var hash = new Array(16);
    for (i = 0; i < 8; i++) {
        hash[2 * i] = H[i].h;
        hash[2 * i + 1] = H[i].l;
    }
    return hash;
}

// A constructor for 64-bit numbers
function int64(h, l) {
    this.h = h;
    this.l = l;
    // this.toString = int64toString;
}

// Copies src into dst, assuming both are 64-bit numbers
function int64copy(dst, src) {
    dst.h = src.h;
    dst.l = src.l;
}

// Right-rotates a 64-bit number by shift
// Won't handle cases of shift>=32
// The function revrrot() is for that
function int64rrot(dst, x, shift) {
    dst.l = (x.l >>> shift) | (x.h << (32 - shift));
    dst.h = (x.h >>> shift) | (x.l << (32 - shift));
}

// Reverses the dwords of the source and then rotates right by shift.
// This is equivalent to rotation by 32+shift
function int64revrrot(dst, x, shift) {
    dst.l = (x.h >>> shift) | (x.l << (32 - shift));
    dst.h = (x.l >>> shift) | (x.h << (32 - shift));
}

// Bitwise-shifts right a 64-bit number by shift
// Won't handle shift>=32, but it's never needed in SHA512
function int64shr(dst, x, shift) {
    dst.l = (x.l >>> shift) | (x.h << (32 - shift));
    dst.h = (x.h >>> shift);
}

// Adds two 64-bit numbers
// Like the original implementation, does not rely on 32-bit operations
function int64add(dst, x, y) {
    var w0 = (x.l & 0xffff) + (y.l & 0xffff);
    var w1 = (x.l >>> 16) + (y.l >>> 16) + (w0 >>> 16);
    var w2 = (x.h & 0xffff) + (y.h & 0xffff) + (w1 >>> 16);
    var w3 = (x.h >>> 16) + (y.h >>> 16) + (w2 >>> 16);
    dst.l = (w0 & 0xffff) | (w1 << 16);
    dst.h = (w2 & 0xffff) | (w3 << 16);
}

// Same, except with 4 addends. Works faster than adding them one by one.
function int64add4(dst, a, b, c, d) {
    var w0 = (a.l & 0xffff) + (b.l & 0xffff) + (c.l & 0xffff) + (d.l & 0xffff);
    var w1 = (a.l >>> 16) + (b.l >>> 16) + (c.l >>> 16) + (d.l >>> 16) + (w0 >>> 16);
    var w2 = (a.h & 0xffff) + (b.h & 0xffff) + (c.h & 0xffff) + (d.h & 0xffff) + (w1 >>> 16);
    var w3 = (a.h >>> 16) + (b.h >>> 16) + (c.h >>> 16) + (d.h >>> 16) + (w2 >>> 16);
    dst.l = (w0 & 0xffff) | (w1 << 16);
    dst.h = (w2 & 0xffff) | (w3 << 16);
}

// Same, except with 5 addends
function int64add5(dst, a, b, c, d, e) {
    var w0 = (a.l & 0xffff) + (b.l & 0xffff) + (c.l & 0xffff) + (d.l & 0xffff) + (e.l & 0xffff);
    var w1 = (a.l >>> 16) + (b.l >>> 16) + (c.l >>> 16) + (d.l >>> 16) + (e.l >>> 16) + (w0 >>> 16);
    var w2 = (a.h & 0xffff) + (b.h & 0xffff) + (c.h & 0xffff) + (d.h & 0xffff) + (e.h & 0xffff) + (w1 >>> 16);
    var w3 = (a.h >>> 16) + (b.h >>> 16) + (c.h >>> 16) + (d.h >>> 16) + (e.h >>> 16) + (w2 >>> 16);
    dst.l = (w0 & 0xffff) | (w1 << 16);
    dst.h = (w2 & 0xffff) | (w3 << 16);
}
// 该代码片段来自于: http://www.sharejs.com/codes/javascript/8904


//var password = 'L0ck it up saf3';
//var plaintext = 'pssst ... đon’t tell anyøne!';
//var ciphertext = Aes.Ctr.encrypt(plaintext, password, 256);
//var origtext = Aes.Ctr.decrypt(ciphertext, password, 256);

/* - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - */
/* AES implementation in JavaScript (c) Chris Veness 2005-2011 */
/* - see http://csrc.nist.gov/publications/PubsFIPS.html#197 */
/* - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - */

window.Aes = {}; // Aes namespace

/**
 * AES Cipher function: encrypt 'input' state with Rijndael algorithm applies Nr rounds (10/12/14) using key schedule w for 'add round key' stage
 *
 * @param {Number[]}
 *            input 16-byte (128-bit) input state array
 * @param {Number[][]}
 *            w Key schedule as 2D byte-array (Nr+1 x Nb bytes)
 * @returns {Number[]} Encrypted output state array
 */
Aes.cipher = function(input, w) { // main Cipher function [§5.1]
    var Nb = 4; // block size (in words): no of columns in state (fixed at 4 for
    // AES)
    var Nr = w.length / Nb - 1; // no of rounds: 10/12/14 for 128/192/256-bit
    // keys

    var state = [ [], [], [], [] ]; // initialise 4xNb byte-array 'state' with
    // input [§3.4]
    for (var i = 0; i < 4 * Nb; i++)
        state[i % 4][Math.floor(i / 4)] = input[i];

    state = Aes.addRoundKey(state, w, 0, Nb);

    for (var round = 1; round < Nr; round++) {
        state = Aes.subBytes(state, Nb);
        state = Aes.shiftRows(state, Nb);
        state = Aes.mixColumns(state, Nb);
        state = Aes.addRoundKey(state, w, round, Nb);
    }

    state = Aes.subBytes(state, Nb);
    state = Aes.shiftRows(state, Nb);
    state = Aes.addRoundKey(state, w, Nr, Nb);

    var output = new Array(4 * Nb); // convert state to 1-d array before
    // returning [§3.4]
    for (var i = 0; i < 4 * Nb; i++)
        output[i] = state[i % 4][Math.floor(i / 4)];
    return output;
};

/**
 * Perform Key Expansion to generate a Key Schedule
 *
 * @param {Number[]}
 *            key Key as 16/24/32-byte array
 * @returns {Number[][]} Expanded key schedule as 2D byte-array (Nr+1 x Nb bytes)
 */
Aes.keyExpansion = function(key) { // generate Key Schedule (byte-array Nr+1 x
    // Nb) from Key [§5.2]
    var Nb = 4; // block size (in words): no of columns in state (fixed at 4 for
    // AES)
    var Nk = key.length / 4; // key length (in words): 4/6/8 for
    // 128/192/256-bit
    // keys
    var Nr = Nk + 6; // no of rounds: 10/12/14 for 128/192/256-bit keys

    var w = new Array(Nb * (Nr + 1));
    var temp = new Array(4);

    for (var i = 0; i < Nk; i++) {
        var r = [ key[4 * i], key[4 * i + 1], key[4 * i + 2], key[4 * i + 3] ];
        w[i] = r;
    }

    for (var i = Nk; i < (Nb * (Nr + 1)); i++) {
        w[i] = new Array(4);
        for (var t = 0; t < 4; t++)
            temp[t] = w[i - 1][t];
        if (i % Nk == 0) {
            temp = Aes.subWord(Aes.rotWord(temp));
            for (var t = 0; t < 4; t++)
                temp[t] ^= Aes.rCon[i / Nk][t];
        } else if (Nk > 6 && i % Nk == 4) {
            temp = Aes.subWord(temp);
        }
        for (var t = 0; t < 4; t++)
            w[i][t] = w[i - Nk][t] ^ temp[t];
    }

    return w;
};

/*
 * ---- remaining routines are private, not called externally ----
 */

Aes.subBytes = function(s, Nb) { // apply SBox to state S [§5.1.1]
    for (var r = 0; r < 4; r++) {
        for (var c = 0; c < Nb; c++)
            s[r][c] = Aes.sBox[s[r][c]];
    }
    return s;
};

Aes.shiftRows = function(s, Nb) { // shift row r of state S left by r bytes
    // [§5.1.2]
    var t = new Array(4);
    for (var r = 1; r < 4; r++) {
        for (var c = 0; c < 4; c++)
            t[c] = s[r][(c + r) % Nb]; // shift into temp copy
        for (var c = 0; c < 4; c++)
            s[r][c] = t[c]; // and copy back
    } // note that this will work for Nb=4,5,6, but not 7,8 (always 4 for
    // AES):
    return s; // see
    // asmaes.sourceforge.net/rijndael/rijndaelImplementation.pdf
};

Aes.mixColumns = function(s, Nb) { // combine bytes of each col of state S
    // [§5.1.3]
    for (var c = 0; c < 4; c++) {
        var a = new Array(4); // 'a' is a copy of the current column from 's'
        var b = new Array(4); // 'b' is a•{02} in GF(2^8)
        for (var i = 0; i < 4; i++) {
            a[i] = s[i][c];
            b[i] = s[i][c] & 0x80 ? s[i][c] << 1 ^ 0x011b : s[i][c] << 1;

        }
        // a[n] ^ b[n] is a•{03} in GF(2^8)
        s[0][c] = b[0] ^ a[1] ^ b[1] ^ a[2] ^ a[3]; // 2*a0 + 3*a1 + a2 + a3
        s[1][c] = a[0] ^ b[1] ^ a[2] ^ b[2] ^ a[3]; // a0 * 2*a1 + 3*a2 + a3
        s[2][c] = a[0] ^ a[1] ^ b[2] ^ a[3] ^ b[3]; // a0 + a1 + 2*a2 + 3*a3
        s[3][c] = a[0] ^ b[0] ^ a[1] ^ a[2] ^ b[3]; // 3*a0 + a1 + a2 + 2*a3
    }
    return s;
};

Aes.addRoundKey = function(state, w, rnd, Nb) { // xor Round Key into state S
    // [§5.1.4]
    for (var r = 0; r < 4; r++) {
        for (var c = 0; c < Nb; c++)
            state[r][c] ^= w[rnd * 4 + c][r];
    }
    return state;
};

Aes.subWord = function(w) { // apply SBox to 4-byte word w
    for (var i = 0; i < 4; i++)
        w[i] = Aes.sBox[w[i]];
    return w;
};

Aes.rotWord = function(w) { // rotate 4-byte word w left by one byte
    var tmp = w[0];
    for (var i = 0; i < 3; i++)
        w[i] = w[i + 1];
    w[3] = tmp;
    return w;
};

// sBox is pre-computed multiplicative inverse in GF(2^8) used in subBytes and
// keyExpansion [§5.1.1]
Aes.sBox = [ 0x63, 0x7c, 0x77, 0x7b, 0xf2, 0x6b, 0x6f, 0xc5, 0x30, 0x01, 0x67, 0x2b, 0xfe, 0xd7, 0xab, 0x76, 0xca, 0x82, 0xc9, 0x7d, 0xfa, 0x59, 0x47, 0xf0, 0xad, 0xd4, 0xa2, 0xaf, 0x9c, 0xa4, 0x72, 0xc0, 0xb7, 0xfd, 0x93, 0x26, 0x36, 0x3f, 0xf7, 0xcc, 0x34, 0xa5, 0xe5, 0xf1, 0x71, 0xd8, 0x31, 0x15, 0x04, 0xc7, 0x23, 0xc3, 0x18, 0x96, 0x05, 0x9a, 0x07, 0x12, 0x80, 0xe2, 0xeb, 0x27, 0xb2, 0x75, 0x09, 0x83, 0x2c, 0x1a, 0x1b, 0x6e, 0x5a, 0xa0, 0x52, 0x3b, 0xd6, 0xb3, 0x29, 0xe3, 0x2f, 0x84, 0x53, 0xd1, 0x00, 0xed, 0x20, 0xfc, 0xb1, 0x5b, 0x6a, 0xcb, 0xbe, 0x39, 0x4a, 0x4c, 0x58, 0xcf, 0xd0, 0xef, 0xaa, 0xfb, 0x43, 0x4d, 0x33, 0x85, 0x45, 0xf9, 0x02, 0x7f, 0x50, 0x3c, 0x9f, 0xa8, 0x51, 0xa3, 0x40, 0x8f, 0x92, 0x9d, 0x38, 0xf5, 0xbc, 0xb6, 0xda, 0x21, 0x10, 0xff, 0xf3, 0xd2, 0xcd, 0x0c, 0x13, 0xec, 0x5f, 0x97, 0x44, 0x17, 0xc4, 0xa7, 0x7e, 0x3d, 0x64, 0x5d, 0x19, 0x73, 0x60, 0x81, 0x4f, 0xdc, 0x22, 0x2a, 0x90, 0x88, 0x46, 0xee, 0xb8, 0x14, 0xde, 0x5e, 0x0b, 0xdb, 0xe0, 0x32, 0x3a, 0x0a,
    0x49, 0x06, 0x24, 0x5c, 0xc2, 0xd3, 0xac, 0x62, 0x91, 0x95, 0xe4, 0x79, 0xe7, 0xc8, 0x37, 0x6d, 0x8d, 0xd5, 0x4e, 0xa9, 0x6c, 0x56, 0xf4, 0xea, 0x65, 0x7a, 0xae, 0x08, 0xba, 0x78, 0x25, 0x2e, 0x1c, 0xa6, 0xb4, 0xc6, 0xe8, 0xdd, 0x74, 0x1f, 0x4b, 0xbd, 0x8b, 0x8a, 0x70, 0x3e, 0xb5, 0x66, 0x48, 0x03, 0xf6, 0x0e, 0x61, 0x35, 0x57, 0xb9, 0x86, 0xc1, 0x1d, 0x9e, 0xe1, 0xf8, 0x98, 0x11, 0x69, 0xd9, 0x8e, 0x94, 0x9b, 0x1e, 0x87, 0xe9, 0xce, 0x55, 0x28, 0xdf, 0x8c, 0xa1, 0x89, 0x0d, 0xbf, 0xe6, 0x42, 0x68, 0x41, 0x99, 0x2d, 0x0f, 0xb0, 0x54, 0xbb, 0x16 ];

// rCon is Round Constant used for the Key Expansion [1st col is 2^(r-1) in
// GF(2^8)] [§5.2]
Aes.rCon = [ [ 0x00, 0x00, 0x00, 0x00 ], [ 0x01, 0x00, 0x00, 0x00 ], [ 0x02, 0x00, 0x00, 0x00 ], [ 0x04, 0x00, 0x00, 0x00 ], [ 0x08, 0x00, 0x00, 0x00 ], [ 0x10, 0x00, 0x00, 0x00 ], [ 0x20, 0x00, 0x00, 0x00 ], [ 0x40, 0x00, 0x00, 0x00 ], [ 0x80, 0x00, 0x00, 0x00 ], [ 0x1b, 0x00, 0x00, 0x00 ], [ 0x36, 0x00, 0x00, 0x00 ] ];

/* - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - */
/* AES Counter-mode implementation in JavaScript (c) Chris Veness 2005-2011 */
/* - see http://csrc.nist.gov/publications/nistpubs/800-38a/sp800-38a.pdf */
/* - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - */

Aes.Ctr = {}; // Aes.Ctr namespace: a subclass or extension of Aes

/**
 * Encrypt a text using AES encryption in Counter mode of operation
 *
 * Unicode multi-byte character safe
 *
 * @param {String}
 *            plaintext Source text to be encrypted
 * @param {String}
 *            password The password to use to generate a key
 * @param {Number}
 *            nBits Number of bits to be used in the key (128, 192, or 256)
 * @returns {string} Encrypted text
 */
Aes.Ctr.encrypt = function(plaintext, password, nBits) {
    var blockSize = 16; // block size fixed at 16 bytes / 128 bits (Nb=4) for
    // AES
    if (!(nBits == 128 || nBits == 192 || nBits == 256))
        return ''; // standard allows 128/192/256 bit keys
    plaintext = Utf8.encode(plaintext);
    password = Utf8.encode(password);
    // var t = new Date(); // timer

    // use AES itself to encrypt password to get cipher key (using plain
    // password as source for key
    // expansion) - gives us well encrypted key
    var nBytes = nBits / 8; // no bytes in key
    var pwBytes = new Array(nBytes);
    for (var i = 0; i < nBytes; i++) {
        pwBytes[i] = isNaN(password.charCodeAt(i)) ? 0 : password.charCodeAt(i);
    }
    var key = Aes.cipher(pwBytes, Aes.keyExpansion(pwBytes)); // gives us
    // 16-byte key
    key = key.concat(key.slice(0, nBytes - 16)); // expand key to 16/24/32
    // bytes long

    // initialise 1st 8 bytes of counter block with nonce (NIST SP800-38A §B.2):
    // [0-1] = millisec,
    // [2-3] = random, [4-7] = seconds, together giving full sub-millisec
    // uniqueness up to Feb 2106
    var counterBlock = new Array(blockSize);

    // var nonce = (new Date()).getTime(); // timestamp: milliseconds since
    var nonce = 1467035047547;
    // 1-Jan-1970
    var nonceMs = nonce % 1000;
    var nonceSec = Math.floor(nonce / 1000);
    // var nonceRnd = Math.floor(Math.random() * 0xffff);
    var nonceRnd = 1;

    for (var i = 0; i < 2; i++)
        counterBlock[i] = (nonceMs >>> i * 8) & 0xff;
    for (var i = 0; i < 2; i++)
        counterBlock[i + 2] = (nonceRnd >>> i * 8) & 0xff;
    for (var i = 0; i < 4; i++)
        counterBlock[i + 4] = (nonceSec >>> i * 8) & 0xff;

    // and convert it to a string to go on the front of the ciphertext
    var ctrTxt = '';
    for (var i = 0; i < 8; i++)
        ctrTxt += String.fromCharCode(counterBlock[i]);

    // generate key schedule - an expansion of the key into distinct Key Rounds
    // for each round
    var keySchedule = Aes.keyExpansion(key);

    var blockCount = Math.ceil(plaintext.length / blockSize);
    var ciphertxt = new Array(blockCount); // ciphertext as array of strings

    for (var b = 0; b < blockCount; b++) {
        // set counter (block #) in last 8 bytes of counter block (leaving nonce
        // in 1st 8 bytes)
        // done in two stages for 32-bit ops: using two words allows us to go
        // past 2^32 blocks (68GB)
        for (var c = 0; c < 4; c++)
            counterBlock[15 - c] = (b >>> c * 8) & 0xff;
        for (var c = 0; c < 4; c++)
            counterBlock[15 - c - 4] = (b / 0x100000000 >>> c * 8);

        var cipherCntr = Aes.cipher(counterBlock, keySchedule); // -- encrypt
        // counter block
        // --

        // block size is reduced on final block
        var blockLength = b < blockCount - 1 ? blockSize : (plaintext.length - 1) % blockSize + 1;
        var cipherChar = new Array(blockLength);

        for (var i = 0; i < blockLength; i++) { // -- xor plaintext with
            // ciphered counter
            // char-by-char --
            cipherChar[i] = cipherCntr[i] ^ plaintext.charCodeAt(b * blockSize + i);
            cipherChar[i] = String.fromCharCode(cipherChar[i]);
        }
        ciphertxt[b] = cipherChar.join('');
    }

    // Array.join is more efficient than repeated string concatenation in IE
    var ciphertext = ctrTxt + ciphertxt.join('');
    for (var i = 0, str = []; i < ciphertext.length; i++) {
        str.push(ciphertext.charCodeAt(i));
    }
    ciphertext = Base64.encode(ciphertext); // encode in base64

    // alert((new Date()) - t);
    return ciphertext;
};

/**
 * Decrypt a text encrypted by AES in counter mode of operation
 *
 * @param {String}
 *            ciphertext Source text to be encrypted
 * @param {String}
 *            password The password to use to generate a key
 * @param {Number}
 *            nBits Number of bits to be used in the key (128, 192, or 256)
 * @returns {String} Decrypted text
 */
Aes.Ctr.decrypt = function(ciphertext, password, nBits) {
    var blockSize = 16; // block size fixed at 16 bytes / 128 bits (Nb=4) for
    // AES
    if (!(nBits == 128 || nBits == 192 || nBits == 256))
        return ''; // standard allows 128/192/256 bit keys
    ciphertext = Base64.decode(ciphertext);
    password = Utf8.encode(password);
    // var t = new Date(); // timer

    // use AES to encrypt password (mirroring encrypt routine)
    var nBytes = nBits / 8; // no bytes in key
    var pwBytes = new Array(nBytes);
    for (var i = 0; i < nBytes; i++) {
        pwBytes[i] = isNaN(password.charCodeAt(i)) ? 0 : password.charCodeAt(i);
    }
    var key = Aes.cipher(pwBytes, Aes.keyExpansion(pwBytes));
    key = key.concat(key.slice(0, nBytes - 16)); // expand key to 16/24/32
    // bytes long

    // recover nonce from 1st 8 bytes of ciphertext
    var counterBlock = new Array(8);
    ctrTxt = ciphertext.slice(0, 8);
    for (var i = 0; i < 8; i++)
        counterBlock[i] = ctrTxt.charCodeAt(i);

    // generate key schedule
    var keySchedule = Aes.keyExpansion(key);

    // separate ciphertext into blocks (skipping past initial 8 bytes)
    var nBlocks = Math.ceil((ciphertext.length - 8) / blockSize);
    var ct = new Array(nBlocks);
    for (var b = 0; b < nBlocks; b++)
        ct[b] = ciphertext.slice(8 + b * blockSize, 8 + b * blockSize + blockSize);
    ciphertext = ct; // ciphertext is now array of block-length strings

    // plaintext will get generated block-by-block into array of block-length
    // strings
    var plaintxt = new Array(ciphertext.length);

    for (var b = 0; b < nBlocks; b++) {
        // set counter (block #) in last 8 bytes of counter block (leaving nonce
        // in 1st 8 bytes)
        for (var c = 0; c < 4; c++)
            counterBlock[15 - c] = ((b) >>> c * 8) & 0xff;
        for (var c = 0; c < 4; c++)
            counterBlock[15 - c - 4] = (((b + 1) / 0x100000000 - 1) >>> c * 8) & 0xff;

        var cipherCntr = Aes.cipher(counterBlock, keySchedule); // encrypt
        // counter block

        var plaintxtByte = new Array(ciphertext[b].length);
        for (var i = 0; i < ciphertext[b].length; i++) {
            // -- xor plaintxt with ciphered counter byte-by-byte --
            plaintxtByte[i] = cipherCntr[i] ^ ciphertext[b].charCodeAt(i);
            plaintxtByte[i] = String.fromCharCode(plaintxtByte[i]);
        }
        plaintxt[b] = plaintxtByte.join('');
    }

    // join array of blocks into single plaintext string
    var plaintext = plaintxt.join('');
    plaintext = Utf8.decode(plaintext); // decode from UTF8 back to Unicode
    // multi-byte chars

    // alert((new Date()) - t);
    return plaintext;
};

/* - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - */
/* Base64 class: Base 64 encoding / decoding (c) Chris Veness 2002-2011 */
/* note: depends on Utf8 class */
/* - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - */

window.Base64 = {}; // Base64 namespace

Base64.code = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=";

/**
 * Encode string into Base64, as defined by RFC 4648 [http://tools.ietf.org/html/rfc4648] (instance method extending String object). As per RFC 4648, no newlines are added.
 *
 * @param {String}
 *            str The string to be encoded as base-64
 * @param {Boolean}
 *            [utf8encode=false] Flag to indicate whether str is Unicode string to be encoded to UTF8 before conversion to base64; otherwise string is assumed to be 8-bit characters
 * @returns {String} Base64-encoded string
 */
Base64.encode = function(str, utf8encode) { // http://tools.ietf.org/html/rfc4648
    utf8encode = (typeof utf8encode == 'undefined') ? false : utf8encode;
    var o1, o2, o3, bits, h1, h2, h3, h4, e = [], pad = '', c, plain, coded;
    var b64 = Base64.code;

    plain = utf8encode ? str.encodeUTF8() : str;

    c = plain.length % 3; // pad string to length of multiple of 3
    if (c > 0) {
        while (c++ < 3) {
            pad += '|';
            plain += '\0';
        }
    }
    // note: doing padding here saves us doing special-case packing for trailing
    // 1 or 2 chars

    for (c = 0; c < plain.length; c += 3) { // pack three octets into four
        // hexets
        o1 = plain.charCodeAt(c);
        o2 = plain.charCodeAt(c + 1);
        o3 = plain.charCodeAt(c + 2);

        bits = o1 << 16 | o2 << 8 | o3;

        h1 = bits >> 18 & 0x3f;
        h2 = bits >> 12 & 0x3f;
        h3 = bits >> 6 & 0x3f;
        h4 = bits & 0x3f;

        // use hextets to index into code string
        e[c / 3] = String.fromCharCode(h1 + 48, h2 + 48, h3 + 48, h4 + 48);
        // e[c / 3] = b64.charAt(h1) + b64.charAt(h2) + b64.charAt(h3) +
        // b64.charAt(h4);
    }
    coded = e.join(''); // join() is far faster than repeated string
    // concatenation in IE

    // replace 'A's from padded nulls with '='s
    coded = coded.slice(0, coded.length - pad.length) + pad;

    return coded;
};

/**
 * Decode string from Base64, as defined by RFC 4648 [http://tools.ietf.org/html/rfc4648] (instance method extending String object). As per RFC 4648, newlines are not catered for.
 *
 * @param {String}
 *            str The string to be decoded from base-64
 * @param {Boolean}
 *            [utf8decode=false] Flag to indicate whether str is Unicode string to be decoded from UTF8 after conversion from base64
 * @returns {String} decoded string
 */
Base64.decode = function(str, utf8decode) {
    utf8decode = (typeof utf8decode == 'undefined') ? false : utf8decode;
    var o1, o2, o3, h1, h2, h3, h4, bits, d = [], plain, coded;
    var b64 = Base64.code;

    coded = utf8decode ? str.decodeUTF8() : str;

    for (var c = 0; c < coded.length; c += 4) { // unpack four hexets into
        // three octets
        // h1 = b64.indexOf(coded.charAt(c));
        // h2 = b64.indexOf(coded.charAt(c + 1));
        // h3 = b64.indexOf(coded.charAt(c + 2));
        // h4 = b64.indexOf(coded.charAt(c + 3));
        h1 = coded.charCodeAt(c) - 48;
        h2 = coded.charCodeAt(c + 1) - 48;
        h3 = coded.charCodeAt(c + 2) - 48;
        h4 = coded.charCodeAt(c + 3) - 48;

        bits = h1 << 18 | h2 << 12 | h3 << 6 | h4;

        o1 = bits >>> 16 & 0xff;
        o2 = bits >>> 8 & 0xff;
        o3 = bits & 0xff;

        d[c / 4] = String.fromCharCode(o1, o2, o3);
        // check for padding
        if (h4 == 0x4C)
            d[c / 4] = String.fromCharCode(o1, o2);
        if (h3 == 0x4C)
            d[c / 4] = String.fromCharCode(o1);
    }
    plain = d.join(''); // join() is far faster than repeated string
    // concatenation in IE

    return utf8decode ? plain.decodeUTF8() : plain;
};

/* - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - */
/*
 * Utf8 class: encode / decode between multi-byte Unicode characters and UTF-8 multiple
 */
/* single-byte character encoding (c) Chris Veness 2002-2011 */
/* - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - */

window.Utf8 = {}; // Utf8 namespace

/**
 * Encode multi-byte Unicode string into utf-8 multiple single-byte characters (BMP / basic multilingual plane only)
 *
 * Chars in range U+0080 - U+07FF are encoded in 2 chars, U+0800 - U+FFFF in 3 chars
 *
 * @param {String}
 *            strUni Unicode string to be encoded as UTF-8
 * @returns {String} encoded string
 */
Utf8.encode = function(strUni) {
    // use regular expressions & String.replace callback function for better
    // efficiency
    // than procedural approaches
    var strUtf = strUni.replace(/[\u0080-\u07ff]/g, // U+0080 - U+07FF => 2
        // bytes 110yyyyy, 10zzzzzz
        function(c) {
            var cc = c.charCodeAt(0);
            return String.fromCharCode(0xc0 | cc >> 6, 0x80 | cc & 0x3f);
        });
    strUtf = strUtf.replace(/[\u0800-\uffff]/g, // U+0800 - U+FFFF => 3 bytes
        // 1110xxxx, 10yyyyyy, 10zzzzzz
        function(c) {
            var cc = c.charCodeAt(0);
            return String.fromCharCode(0xe0 | cc >> 12, 0x80 | cc >> 6 & 0x3F, 0x80 | cc & 0x3f);
        });
    return strUtf;
};

/**
 * Decode utf-8 encoded string back into multi-byte Unicode characters
 *
 * @param {String}
 *            strUtf UTF-8 string to be decoded back to Unicode
 * @returns {String} decoded string
 */
Utf8.decode = function(strUtf) {
    // note: decode 3-byte chars first as decoded 2-byte strings could appear to
    // be 3-byte char!
    var strUni = strUtf.replace(/[\u00e0-\u00ef][\u0080-\u00bf][\u0080-\u00bf]/g, // 3-byte
        // chars
        function(c) { // (note parentheses for precence)
            var cc = ((c.charCodeAt(0) & 0x0f) << 12) | ((c.charCodeAt(1) & 0x3f) << 6) | (c.charCodeAt(2) & 0x3f);
            return String.fromCharCode(cc);
        });
    strUni = strUni.replace(/[\u00c0-\u00df][\u0080-\u00bf]/g, // 2-byte chars
        function(c) { // (note parentheses for precence)
            var cc = (c.charCodeAt(0) & 0x1f) << 6 | c.charCodeAt(1) & 0x3f;
            return String.fromCharCode(cc);
        });
    return strUni;
};

/* - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - */

/*
 Copyright (c) 2012 Nevins Bartolomeo <nevins.bartolomeo@gmail.com>
 Copyright (c) 2012 Shane Girish <shaneGirish@gmail.com>
 Copyright (c) 2014 Daniel Wirtz <dcode@dcode.io>

 Redistribution and use in source and binary forms, with or without
 modification, are permitted provided that the following conditions
 are met:
 1. Redistributions of source code must retain the above copyright
 notice, this list of conditions and the following disclaimer.
 2. Redistributions in binary form must reproduce the above copyright
 notice, this list of conditions and the following disclaimer in the
 documentation and/or other materials provided with the distribution.
 3. The name of the author may not be used to endorse or promote products
 derived from this software without specific prior written permission.

 THIS SOFTWARE IS PROVIDED BY THE AUTHOR ``AS IS'' AND ANY EXPRESS OR
 IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES
 OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
 IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT, INDIRECT,
 INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT
 NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF
 THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

/**
 * @license bcrypt.js (c) 2013 Daniel Wirtz <dcode@dcode.io>
 * Released under the Apache License, Version 2.0
 * see: https://github.com/dcodeIO/bcrypt.js for details
 */
(function(global, factory) {

    /* AMD */ if (typeof define === 'function' && define["amd"])
        define([], factory);
    /* CommonJS */ else if (typeof require === 'function' && typeof module === "object" && module && module["exports"])
        module["exports"] = factory();
    /* Global */ else
        (global["dcodeIO"] = global["dcodeIO"] || {})["bcrypt"] = factory();

}(this, function() {
    "use strict";

    /**
     * bcrypt namespace.
     * @type {Object.<string,*>}
     */
    var bcrypt = {};

    /**
     * The random implementation to use as a fallback.
     * @type {?function(number):!Array.<number>}
     * @inner
     */
    var randomFallback = null;

    /**
     * Generates cryptographically secure random bytes.
     * @function
     * @param {number} len Bytes length
     * @returns {!Array.<number>} Random bytes
     * @throws {Error} If no random implementation is available
     * @inner
     */
    function random(len) {
        /* node */ if (typeof module !== 'undefined' && module && module['exports'])
            try {
                return require("crypto")['randomBytes'](len);
            } catch (e) {}
        /* WCA */ try {
            var a; (self['crypto']||self['msCrypto'])['getRandomValues'](a = new Uint32Array(len));
            return Array.prototype.slice.call(a);
        } catch (e) {}
        /* fallback */ if (!randomFallback)
            throw Error("Neither WebCryptoAPI nor a crypto module is available. Use bcrypt.setRandomFallback to set an alternative");
        return randomFallback(len);
    }

    // Test if any secure randomness source is available
    var randomAvailable = false;
    try {
        random(1);
        randomAvailable = true;
    } catch (e) {}

    // Default fallback, if any
    randomFallback = null;

    /**
     * Sets the pseudo random number generator to use as a fallback if neither node's `crypto` module nor the Web Crypto
     *  API is available. Please note: It is highly important that the PRNG used is cryptographically secure and that it
     *  is seeded properly!
     * @param {?function(number):!Array.<number>} random Function taking the number of bytes to generate as its
     *  sole argument, returning the corresponding array of cryptographically secure random byte values.
     * @see http://nodejs.org/api/crypto.html
     * @see http://www.w3.org/TR/WebCryptoAPI/
     */
    bcrypt.setRandomFallback = function(random) {
        randomFallback = random;
    };

    /**
     * Synchronously generates a salt.
     * @param {number=} rounds Number of rounds to use, defaults to 10 if omitted
     * @param {number=} seed_length Not supported.
     * @returns {string} Resulting salt
     * @throws {Error} If a random fallback is required but not set
     * @expose
     */
    bcrypt.genSaltSync = function(rounds, seed_length) {
        rounds = rounds || GENSALT_DEFAULT_LOG2_ROUNDS;
        if (typeof rounds !== 'number')
            throw Error("Illegal arguments: "+(typeof rounds)+", "+(typeof seed_length));
        if (rounds < 4)
            rounds = 4;
        else if (rounds > 31)
            rounds = 31;
        var salt = [];
        salt.push("$2a$");
        if (rounds < 10)
            salt.push("0");
        salt.push(rounds.toString());
        salt.push('$');
        salt.push(base64_encode(random(BCRYPT_SALT_LEN), BCRYPT_SALT_LEN)); // May throw
        return salt.join('');
    };

    /**
     * Asynchronously generates a salt.
     * @param {(number|function(Error, string=))=} rounds Number of rounds to use, defaults to 10 if omitted
     * @param {(number|function(Error, string=))=} seed_length Not supported.
     * @param {function(Error, string=)=} callback Callback receiving the error, if any, and the resulting salt
     * @expose
     */
    bcrypt.genSalt = function(rounds, seed_length, callback) {
        if (typeof seed_length === 'function')
            callback = seed_length,
                seed_length = undefined; // Not supported.
        if (typeof rounds === 'function')
            callback = rounds,
                rounds = GENSALT_DEFAULT_LOG2_ROUNDS;
        if (typeof callback !== 'function')
            throw Error("Illegal callback: "+typeof(callback));
        if (typeof rounds !== 'number') {
            nextTick(callback.bind(this, Error("Illegal arguments: "+(typeof rounds))));
            return;
        }
        nextTick(function() { // Pretty thin, but salting is fast enough
            try {
                callback(null, bcrypt.genSaltSync(rounds));
            } catch (err) {
                callback(err);
            }
        });
    };

    /**
     * Synchronously generates a hash for the given string.
     * @param {string} s String to hash
     * @param {(number|string)=} salt Salt length to generate or salt to use, default to 10
     * @returns {string} Resulting hash
     * @expose
     */
    bcrypt.hashSync = function(s, salt) {
        if (typeof salt === 'undefined')
            salt = GENSALT_DEFAULT_LOG2_ROUNDS;
        if (typeof salt === 'number')
            salt = bcrypt.genSaltSync(salt);
        if (typeof s !== 'string' || typeof salt !== 'string')
            throw Error("Illegal arguments: "+(typeof s)+', '+(typeof salt));
        return _hash(s, salt);
    };

    /**
     * Asynchronously generates a hash for the given string.
     * @param {string} s String to hash
     * @param {number|string} salt Salt length to generate or salt to use
     * @param {function(Error, string=)} callback Callback receiving the error, if any, and the resulting hash
     * @param {function(number)=} progressCallback Callback successively called with the percentage of rounds completed
     *  (0.0 - 1.0), maximally once per `MAX_EXECUTION_TIME = 100` ms.
     * @expose
     */
    bcrypt.hash = function(s, salt, callback, progressCallback) {
        if (typeof callback !== 'function')
            throw Error("Illegal callback: "+typeof(callback));
        if (typeof s === 'string' && typeof salt === 'number')
            bcrypt.genSalt(salt, function(err, salt) {
                _hash(s, salt, callback, progressCallback);
            });
        else if (typeof s === 'string' && typeof salt === 'string')
            _hash(s, salt, callback, progressCallback);
        else
            nextTick(callback.bind(this, Error("Illegal arguments: "+(typeof s)+', '+(typeof salt))));
    };

    /**
     * Compares two strings of the same length in constant time.
     * @param {string} known Must be of the correct length
     * @param {string} unknown Must be the same length as `known`
     * @returns {boolean}
     * @inner
     */
    function safeStringCompare(known, unknown) {
        var right = 0,
            wrong = 0;
        for (var i=0, k=known.length; i<k; ++i) {
            if (known.charCodeAt(i) === unknown.charCodeAt(i))
                ++right;
            else
                ++wrong;
        }
        // Prevent removal of unused variables (never true, actually)
        if (right < 0)
            return false;
        return wrong === 0;
    }

    /**
     * Synchronously tests a string against a hash.
     * @param {string} s String to compare
     * @param {string} hash Hash to test against
     * @returns {boolean} true if matching, otherwise false
     * @throws {Error} If an argument is illegal
     * @expose
     */
    bcrypt.compareSync = function(s, hash) {
        if (typeof s !== "string" || typeof hash !== "string")
            throw Error("Illegal arguments: "+(typeof s)+', '+(typeof hash));
        if (hash.length !== 60)
            return false;
        return safeStringCompare(bcrypt.hashSync(s, hash.substr(0, hash.length-31)), hash);
    };

    /**
     * Asynchronously compares the given data against the given hash.
     * @param {string} s Data to compare
     * @param {string} hash Data to be compared to
     * @param {function(Error, boolean)} callback Callback receiving the error, if any, otherwise the result
     * @param {function(number)=} progressCallback Callback successively called with the percentage of rounds completed
     *  (0.0 - 1.0), maximally once per `MAX_EXECUTION_TIME = 100` ms.
     * @throws {Error} If the callback argument is invalid
     * @expose
     */
    bcrypt.compare = function(s, hash, callback, progressCallback) {
        if (typeof callback !== 'function')
            throw Error("Illegal callback: "+typeof(callback));
        if (typeof s !== "string" || typeof hash !== "string") {
            nextTick(callback.bind(this, Error("Illegal arguments: "+(typeof s)+', '+(typeof hash))));
            return;
        }
        if (hash.length !== 60) {
            nextTick(callback.bind(this, null, false));
            return;
        }
        bcrypt.hash(s, hash.substr(0, 29), function(err, comp) {
            if (err)
                callback(err);
            else
                callback(null, safeStringCompare(comp, hash));
        }, progressCallback);
    };

    /**
     * Gets the number of rounds used to encrypt the specified hash.
     * @param {string} hash Hash to extract the used number of rounds from
     * @returns {number} Number of rounds used
     * @throws {Error} If hash is not a string
     * @expose
     */
    bcrypt.getRounds = function(hash) {
        if (typeof hash !== "string")
            throw Error("Illegal arguments: "+(typeof hash));
        return parseInt(hash.split("$")[2], 10);
    };

    /**
     * Gets the salt portion from a hash. Does not validate the hash.
     * @param {string} hash Hash to extract the salt from
     * @returns {string} Extracted salt part
     * @throws {Error} If `hash` is not a string or otherwise invalid
     * @expose
     */
    bcrypt.getSalt = function(hash) {
        if (typeof hash !== 'string')
            throw Error("Illegal arguments: "+(typeof hash));
        if (hash.length !== 60)
            throw Error("Illegal hash length: "+hash.length+" != 60");
        return hash.substring(0, 29);
    };

    /**
     * Continues with the callback on the next tick.
     * @function
     * @param {function(...[*])} callback Callback to execute
     * @inner
     */
    var nextTick = typeof process !== 'undefined' && process && typeof process.nextTick === 'function'
        ? (typeof setImmediate === 'function' ? setImmediate : process.nextTick)
        : setTimeout;

    /**
     * Converts a JavaScript string to UTF8 bytes.
     * @param {string} str String
     * @returns {!Array.<number>} UTF8 bytes
     * @inner
     */
    function stringToBytes(str) {
        var out = [],
            i = 0;
        utfx.encodeUTF16toUTF8(function() {
            if (i >= str.length) return null;
            return str.charCodeAt(i++);
        }, function(b) {
            out.push(b);
        });
        return out;
    }

    // A base64 implementation for the bcrypt algorithm. This is partly non-standard.

    /**
     * bcrypt's own non-standard base64 dictionary.
     * @type {!Array.<string>}
     * @const
     * @inner
     **/
    var BASE64_CODE = "./ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789".split('');

    /**
     * @type {!Array.<number>}
     * @const
     * @inner
     **/
    var BASE64_INDEX = [-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
        -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
        -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 0,
        1, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, -1, -1, -1, -1, -1, -1,
        -1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
        20, 21, 22, 23, 24, 25, 26, 27, -1, -1, -1, -1, -1, -1, 28, 29, 30,
        31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47,
        48, 49, 50, 51, 52, 53, -1, -1, -1, -1, -1];

    /**
     * @type {!function(...number):string}
     * @inner
     */
    var stringFromCharCode = String.fromCharCode;

    /**
     * Encodes a byte array to base64 with up to len bytes of input.
     * @param {!Array.<number>} b Byte array
     * @param {number} len Maximum input length
     * @returns {string}
     * @inner
     */
    function base64_encode(b, len) {
        var off = 0,
            rs = [],
            c1, c2;
        if (len <= 0 || len > b.length)
            throw Error("Illegal len: "+len);
        while (off < len) {
            c1 = b[off++] & 0xff;
            rs.push(BASE64_CODE[(c1 >> 2) & 0x3f]);
            c1 = (c1 & 0x03) << 4;
            if (off >= len) {
                rs.push(BASE64_CODE[c1 & 0x3f]);
                break;
            }
            c2 = b[off++] & 0xff;
            c1 |= (c2 >> 4) & 0x0f;
            rs.push(BASE64_CODE[c1 & 0x3f]);
            c1 = (c2 & 0x0f) << 2;
            if (off >= len) {
                rs.push(BASE64_CODE[c1 & 0x3f]);
                break;
            }
            c2 = b[off++] & 0xff;
            c1 |= (c2 >> 6) & 0x03;
            rs.push(BASE64_CODE[c1 & 0x3f]);
            rs.push(BASE64_CODE[c2 & 0x3f]);
        }
        return rs.join('');
    }

    /**
     * Decodes a base64 encoded string to up to len bytes of output.
     * @param {string} s String to decode
     * @param {number} len Maximum output length
     * @returns {!Array.<number>}
     * @inner
     */
    function base64_decode(s, len) {
        var off = 0,
            slen = s.length,
            olen = 0,
            rs = [],
            c1, c2, c3, c4, o, code;
        if (len <= 0)
            throw Error("Illegal len: "+len);
        while (off < slen - 1 && olen < len) {
            code = s.charCodeAt(off++);
            c1 = code < BASE64_INDEX.length ? BASE64_INDEX[code] : -1;
            code = s.charCodeAt(off++);
            c2 = code < BASE64_INDEX.length ? BASE64_INDEX[code] : -1;
            if (c1 == -1 || c2 == -1)
                break;
            o = (c1 << 2) >>> 0;
            o |= (c2 & 0x30) >> 4;
            rs.push(stringFromCharCode(o));
            if (++olen >= len || off >= slen)
                break;
            code = s.charCodeAt(off++);
            c3 = code < BASE64_INDEX.length ? BASE64_INDEX[code] : -1;
            if (c3 == -1)
                break;
            o = ((c2 & 0x0f) << 4) >>> 0;
            o |= (c3 & 0x3c) >> 2;
            rs.push(stringFromCharCode(o));
            if (++olen >= len || off >= slen)
                break;
            code = s.charCodeAt(off++);
            c4 = code < BASE64_INDEX.length ? BASE64_INDEX[code] : -1;
            o = ((c3 & 0x03) << 6) >>> 0;
            o |= c4;
            rs.push(stringFromCharCode(o));
            ++olen;
        }
        var res = [];
        for (off = 0; off<olen; off++)
            res.push(rs[off].charCodeAt(0));
        return res;
    }

    /**
     * utfx-embeddable (c) 2014 Daniel Wirtz <dcode@dcode.io>
     * Released under the Apache License, Version 2.0
     * see: https://github.com/dcodeIO/utfx for details
     */
    var utfx = function() {
        "use strict";

        /**
         * utfx namespace.
         * @inner
         * @type {!Object.<string,*>}
         */
        var utfx = {};

        /**
         * Maximum valid code point.
         * @type {number}
         * @const
         */
        utfx.MAX_CODEPOINT = 0x10FFFF;

        /**
         * Encodes UTF8 code points to UTF8 bytes.
         * @param {(!function():number|null) | number} src Code points source, either as a function returning the next code point
         *  respectively `null` if there are no more code points left or a single numeric code point.
         * @param {!function(number)} dst Bytes destination as a function successively called with the next byte
         */
        utfx.encodeUTF8 = function(src, dst) {
            var cp = null;
            if (typeof src === 'number')
                cp = src,
                    src = function() { return null; };
            while (cp !== null || (cp = src()) !== null) {
                if (cp < 0x80)
                    dst(cp&0x7F);
                else if (cp < 0x800)
                    dst(((cp>>6)&0x1F)|0xC0),
                        dst((cp&0x3F)|0x80);
                else if (cp < 0x10000)
                    dst(((cp>>12)&0x0F)|0xE0),
                        dst(((cp>>6)&0x3F)|0x80),
                        dst((cp&0x3F)|0x80);
                else
                    dst(((cp>>18)&0x07)|0xF0),
                        dst(((cp>>12)&0x3F)|0x80),
                        dst(((cp>>6)&0x3F)|0x80),
                        dst((cp&0x3F)|0x80);
                cp = null;
            }
        };

        /**
         * Decodes UTF8 bytes to UTF8 code points.
         * @param {!function():number|null} src Bytes source as a function returning the next byte respectively `null` if there
         *  are no more bytes left.
         * @param {!function(number)} dst Code points destination as a function successively called with each decoded code point.
         * @throws {RangeError} If a starting byte is invalid in UTF8
         * @throws {Error} If the last sequence is truncated. Has an array property `bytes` holding the
         *  remaining bytes.
         */
        utfx.decodeUTF8 = function(src, dst) {
            var a, b, c, d, fail = function(b) {
                b = b.slice(0, b.indexOf(null));
                var err = Error(b.toString());
                err.name = "TruncatedError";
                err['bytes'] = b;
                throw err;
            };
            while ((a = src()) !== null) {
                if ((a&0x80) === 0)
                    dst(a);
                else if ((a&0xE0) === 0xC0)
                    ((b = src()) === null) && fail([a, b]),
                        dst(((a&0x1F)<<6) | (b&0x3F));
                else if ((a&0xF0) === 0xE0)
                    ((b=src()) === null || (c=src()) === null) && fail([a, b, c]),
                        dst(((a&0x0F)<<12) | ((b&0x3F)<<6) | (c&0x3F));
                else if ((a&0xF8) === 0xF0)
                    ((b=src()) === null || (c=src()) === null || (d=src()) === null) && fail([a, b, c ,d]),
                        dst(((a&0x07)<<18) | ((b&0x3F)<<12) | ((c&0x3F)<<6) | (d&0x3F));
                else throw RangeError("Illegal starting byte: "+a);
            }
        };

        /**
         * Converts UTF16 characters to UTF8 code points.
         * @param {!function():number|null} src Characters source as a function returning the next char code respectively
         *  `null` if there are no more characters left.
         * @param {!function(number)} dst Code points destination as a function successively called with each converted code
         *  point.
         */
        utfx.UTF16toUTF8 = function(src, dst) {
            var c1, c2 = null;
            while (true) {
                if ((c1 = c2 !== null ? c2 : src()) === null)
                    break;
                if (c1 >= 0xD800 && c1 <= 0xDFFF) {
                    if ((c2 = src()) !== null) {
                        if (c2 >= 0xDC00 && c2 <= 0xDFFF) {
                            dst((c1-0xD800)*0x400+c2-0xDC00+0x10000);
                            c2 = null; continue;
                        }
                    }
                }
                dst(c1);
            }
            if (c2 !== null) dst(c2);
        };

        /**
         * Converts UTF8 code points to UTF16 characters.
         * @param {(!function():number|null) | number} src Code points source, either as a function returning the next code point
         *  respectively `null` if there are no more code points left or a single numeric code point.
         * @param {!function(number)} dst Characters destination as a function successively called with each converted char code.
         * @throws {RangeError} If a code point is out of range
         */
        utfx.UTF8toUTF16 = function(src, dst) {
            var cp = null;
            if (typeof src === 'number')
                cp = src, src = function() { return null; };
            while (cp !== null || (cp = src()) !== null) {
                if (cp <= 0xFFFF)
                    dst(cp);
                else
                    cp -= 0x10000,
                        dst((cp>>10)+0xD800),
                        dst((cp%0x400)+0xDC00);
                cp = null;
            }
        };

        /**
         * Converts and encodes UTF16 characters to UTF8 bytes.
         * @param {!function():number|null} src Characters source as a function returning the next char code respectively `null`
         *  if there are no more characters left.
         * @param {!function(number)} dst Bytes destination as a function successively called with the next byte.
         */
        utfx.encodeUTF16toUTF8 = function(src, dst) {
            utfx.UTF16toUTF8(src, function(cp) {
                utfx.encodeUTF8(cp, dst);
            });
        };

        /**
         * Decodes and converts UTF8 bytes to UTF16 characters.
         * @param {!function():number|null} src Bytes source as a function returning the next byte respectively `null` if there
         *  are no more bytes left.
         * @param {!function(number)} dst Characters destination as a function successively called with each converted char code.
         * @throws {RangeError} If a starting byte is invalid in UTF8
         * @throws {Error} If the last sequence is truncated. Has an array property `bytes` holding the remaining bytes.
         */
        utfx.decodeUTF8toUTF16 = function(src, dst) {
            utfx.decodeUTF8(src, function(cp) {
                utfx.UTF8toUTF16(cp, dst);
            });
        };

        /**
         * Calculates the byte length of an UTF8 code point.
         * @param {number} cp UTF8 code point
         * @returns {number} Byte length
         */
        utfx.calculateCodePoint = function(cp) {
            return (cp < 0x80) ? 1 : (cp < 0x800) ? 2 : (cp < 0x10000) ? 3 : 4;
        };

        /**
         * Calculates the number of UTF8 bytes required to store UTF8 code points.
         * @param {(!function():number|null)} src Code points source as a function returning the next code point respectively
         *  `null` if there are no more code points left.
         * @returns {number} The number of UTF8 bytes required
         */
        utfx.calculateUTF8 = function(src) {
            var cp, l=0;
            while ((cp = src()) !== null)
                l += utfx.calculateCodePoint(cp);
            return l;
        };

        /**
         * Calculates the number of UTF8 code points respectively UTF8 bytes required to store UTF16 char codes.
         * @param {(!function():number|null)} src Characters source as a function returning the next char code respectively
         *  `null` if there are no more characters left.
         * @returns {!Array.<number>} The number of UTF8 code points at index 0 and the number of UTF8 bytes required at index 1.
         */
        utfx.calculateUTF16asUTF8 = function(src) {
            var n=0, l=0;
            utfx.UTF16toUTF8(src, function(cp) {
                ++n; l += utfx.calculateCodePoint(cp);
            });
            return [n,l];
        };

        return utfx;
    }();

    Date.now = Date.now || function() { return +new Date; };

    /**
     * @type {number}
     * @const
     * @inner
     */
    var BCRYPT_SALT_LEN = 16;

    /**
     * @type {number}
     * @const
     * @inner
     */
    var GENSALT_DEFAULT_LOG2_ROUNDS = 10;

    /**
     * @type {number}
     * @const
     * @inner
     */
    var BLOWFISH_NUM_ROUNDS = 16;

    /**
     * @type {number}
     * @const
     * @inner
     */
    var MAX_EXECUTION_TIME = 100;

    /**
     * @type {Array.<number>}
     * @const
     * @inner
     */
    var P_ORIG = [
        0x243f6a88, 0x85a308d3, 0x13198a2e, 0x03707344, 0xa4093822,
        0x299f31d0, 0x082efa98, 0xec4e6c89, 0x452821e6, 0x38d01377,
        0xbe5466cf, 0x34e90c6c, 0xc0ac29b7, 0xc97c50dd, 0x3f84d5b5,
        0xb5470917, 0x9216d5d9, 0x8979fb1b
    ];

    /**
     * @type {Array.<number>}
     * @const
     * @inner
     */
    var S_ORIG = [
        0xd1310ba6, 0x98dfb5ac, 0x2ffd72db, 0xd01adfb7, 0xb8e1afed,
        0x6a267e96, 0xba7c9045, 0xf12c7f99, 0x24a19947, 0xb3916cf7,
        0x0801f2e2, 0x858efc16, 0x636920d8, 0x71574e69, 0xa458fea3,
        0xf4933d7e, 0x0d95748f, 0x728eb658, 0x718bcd58, 0x82154aee,
        0x7b54a41d, 0xc25a59b5, 0x9c30d539, 0x2af26013, 0xc5d1b023,
        0x286085f0, 0xca417918, 0xb8db38ef, 0x8e79dcb0, 0x603a180e,
        0x6c9e0e8b, 0xb01e8a3e, 0xd71577c1, 0xbd314b27, 0x78af2fda,
        0x55605c60, 0xe65525f3, 0xaa55ab94, 0x57489862, 0x63e81440,
        0x55ca396a, 0x2aab10b6, 0xb4cc5c34, 0x1141e8ce, 0xa15486af,
        0x7c72e993, 0xb3ee1411, 0x636fbc2a, 0x2ba9c55d, 0x741831f6,
        0xce5c3e16, 0x9b87931e, 0xafd6ba33, 0x6c24cf5c, 0x7a325381,
        0x28958677, 0x3b8f4898, 0x6b4bb9af, 0xc4bfe81b, 0x66282193,
        0x61d809cc, 0xfb21a991, 0x487cac60, 0x5dec8032, 0xef845d5d,
        0xe98575b1, 0xdc262302, 0xeb651b88, 0x23893e81, 0xd396acc5,
        0x0f6d6ff3, 0x83f44239, 0x2e0b4482, 0xa4842004, 0x69c8f04a,
        0x9e1f9b5e, 0x21c66842, 0xf6e96c9a, 0x670c9c61, 0xabd388f0,
        0x6a51a0d2, 0xd8542f68, 0x960fa728, 0xab5133a3, 0x6eef0b6c,
        0x137a3be4, 0xba3bf050, 0x7efb2a98, 0xa1f1651d, 0x39af0176,
        0x66ca593e, 0x82430e88, 0x8cee8619, 0x456f9fb4, 0x7d84a5c3,
        0x3b8b5ebe, 0xe06f75d8, 0x85c12073, 0x401a449f, 0x56c16aa6,
        0x4ed3aa62, 0x363f7706, 0x1bfedf72, 0x429b023d, 0x37d0d724,
        0xd00a1248, 0xdb0fead3, 0x49f1c09b, 0x075372c9, 0x80991b7b,
        0x25d479d8, 0xf6e8def7, 0xe3fe501a, 0xb6794c3b, 0x976ce0bd,
        0x04c006ba, 0xc1a94fb6, 0x409f60c4, 0x5e5c9ec2, 0x196a2463,
        0x68fb6faf, 0x3e6c53b5, 0x1339b2eb, 0x3b52ec6f, 0x6dfc511f,
        0x9b30952c, 0xcc814544, 0xaf5ebd09, 0xbee3d004, 0xde334afd,
        0x660f2807, 0x192e4bb3, 0xc0cba857, 0x45c8740f, 0xd20b5f39,
        0xb9d3fbdb, 0x5579c0bd, 0x1a60320a, 0xd6a100c6, 0x402c7279,
        0x679f25fe, 0xfb1fa3cc, 0x8ea5e9f8, 0xdb3222f8, 0x3c7516df,
        0xfd616b15, 0x2f501ec8, 0xad0552ab, 0x323db5fa, 0xfd238760,
        0x53317b48, 0x3e00df82, 0x9e5c57bb, 0xca6f8ca0, 0x1a87562e,
        0xdf1769db, 0xd542a8f6, 0x287effc3, 0xac6732c6, 0x8c4f5573,
        0x695b27b0, 0xbbca58c8, 0xe1ffa35d, 0xb8f011a0, 0x10fa3d98,
        0xfd2183b8, 0x4afcb56c, 0x2dd1d35b, 0x9a53e479, 0xb6f84565,
        0xd28e49bc, 0x4bfb9790, 0xe1ddf2da, 0xa4cb7e33, 0x62fb1341,
        0xcee4c6e8, 0xef20cada, 0x36774c01, 0xd07e9efe, 0x2bf11fb4,
        0x95dbda4d, 0xae909198, 0xeaad8e71, 0x6b93d5a0, 0xd08ed1d0,
        0xafc725e0, 0x8e3c5b2f, 0x8e7594b7, 0x8ff6e2fb, 0xf2122b64,
        0x8888b812, 0x900df01c, 0x4fad5ea0, 0x688fc31c, 0xd1cff191,
        0xb3a8c1ad, 0x2f2f2218, 0xbe0e1777, 0xea752dfe, 0x8b021fa1,
        0xe5a0cc0f, 0xb56f74e8, 0x18acf3d6, 0xce89e299, 0xb4a84fe0,
        0xfd13e0b7, 0x7cc43b81, 0xd2ada8d9, 0x165fa266, 0x80957705,
        0x93cc7314, 0x211a1477, 0xe6ad2065, 0x77b5fa86, 0xc75442f5,
        0xfb9d35cf, 0xebcdaf0c, 0x7b3e89a0, 0xd6411bd3, 0xae1e7e49,
        0x00250e2d, 0x2071b35e, 0x226800bb, 0x57b8e0af, 0x2464369b,
        0xf009b91e, 0x5563911d, 0x59dfa6aa, 0x78c14389, 0xd95a537f,
        0x207d5ba2, 0x02e5b9c5, 0x83260376, 0x6295cfa9, 0x11c81968,
        0x4e734a41, 0xb3472dca, 0x7b14a94a, 0x1b510052, 0x9a532915,
        0xd60f573f, 0xbc9bc6e4, 0x2b60a476, 0x81e67400, 0x08ba6fb5,
        0x571be91f, 0xf296ec6b, 0x2a0dd915, 0xb6636521, 0xe7b9f9b6,
        0xff34052e, 0xc5855664, 0x53b02d5d, 0xa99f8fa1, 0x08ba4799,
        0x6e85076a, 0x4b7a70e9, 0xb5b32944, 0xdb75092e, 0xc4192623,
        0xad6ea6b0, 0x49a7df7d, 0x9cee60b8, 0x8fedb266, 0xecaa8c71,
        0x699a17ff, 0x5664526c, 0xc2b19ee1, 0x193602a5, 0x75094c29,
        0xa0591340, 0xe4183a3e, 0x3f54989a, 0x5b429d65, 0x6b8fe4d6,
        0x99f73fd6, 0xa1d29c07, 0xefe830f5, 0x4d2d38e6, 0xf0255dc1,
        0x4cdd2086, 0x8470eb26, 0x6382e9c6, 0x021ecc5e, 0x09686b3f,
        0x3ebaefc9, 0x3c971814, 0x6b6a70a1, 0x687f3584, 0x52a0e286,
        0xb79c5305, 0xaa500737, 0x3e07841c, 0x7fdeae5c, 0x8e7d44ec,
        0x5716f2b8, 0xb03ada37, 0xf0500c0d, 0xf01c1f04, 0x0200b3ff,
        0xae0cf51a, 0x3cb574b2, 0x25837a58, 0xdc0921bd, 0xd19113f9,
        0x7ca92ff6, 0x94324773, 0x22f54701, 0x3ae5e581, 0x37c2dadc,
        0xc8b57634, 0x9af3dda7, 0xa9446146, 0x0fd0030e, 0xecc8c73e,
        0xa4751e41, 0xe238cd99, 0x3bea0e2f, 0x3280bba1, 0x183eb331,
        0x4e548b38, 0x4f6db908, 0x6f420d03, 0xf60a04bf, 0x2cb81290,
        0x24977c79, 0x5679b072, 0xbcaf89af, 0xde9a771f, 0xd9930810,
        0xb38bae12, 0xdccf3f2e, 0x5512721f, 0x2e6b7124, 0x501adde6,
        0x9f84cd87, 0x7a584718, 0x7408da17, 0xbc9f9abc, 0xe94b7d8c,
        0xec7aec3a, 0xdb851dfa, 0x63094366, 0xc464c3d2, 0xef1c1847,
        0x3215d908, 0xdd433b37, 0x24c2ba16, 0x12a14d43, 0x2a65c451,
        0x50940002, 0x133ae4dd, 0x71dff89e, 0x10314e55, 0x81ac77d6,
        0x5f11199b, 0x043556f1, 0xd7a3c76b, 0x3c11183b, 0x5924a509,
        0xf28fe6ed, 0x97f1fbfa, 0x9ebabf2c, 0x1e153c6e, 0x86e34570,
        0xeae96fb1, 0x860e5e0a, 0x5a3e2ab3, 0x771fe71c, 0x4e3d06fa,
        0x2965dcb9, 0x99e71d0f, 0x803e89d6, 0x5266c825, 0x2e4cc978,
        0x9c10b36a, 0xc6150eba, 0x94e2ea78, 0xa5fc3c53, 0x1e0a2df4,
        0xf2f74ea7, 0x361d2b3d, 0x1939260f, 0x19c27960, 0x5223a708,
        0xf71312b6, 0xebadfe6e, 0xeac31f66, 0xe3bc4595, 0xa67bc883,
        0xb17f37d1, 0x018cff28, 0xc332ddef, 0xbe6c5aa5, 0x65582185,
        0x68ab9802, 0xeecea50f, 0xdb2f953b, 0x2aef7dad, 0x5b6e2f84,
        0x1521b628, 0x29076170, 0xecdd4775, 0x619f1510, 0x13cca830,
        0xeb61bd96, 0x0334fe1e, 0xaa0363cf, 0xb5735c90, 0x4c70a239,
        0xd59e9e0b, 0xcbaade14, 0xeecc86bc, 0x60622ca7, 0x9cab5cab,
        0xb2f3846e, 0x648b1eaf, 0x19bdf0ca, 0xa02369b9, 0x655abb50,
        0x40685a32, 0x3c2ab4b3, 0x319ee9d5, 0xc021b8f7, 0x9b540b19,
        0x875fa099, 0x95f7997e, 0x623d7da8, 0xf837889a, 0x97e32d77,
        0x11ed935f, 0x16681281, 0x0e358829, 0xc7e61fd6, 0x96dedfa1,
        0x7858ba99, 0x57f584a5, 0x1b227263, 0x9b83c3ff, 0x1ac24696,
        0xcdb30aeb, 0x532e3054, 0x8fd948e4, 0x6dbc3128, 0x58ebf2ef,
        0x34c6ffea, 0xfe28ed61, 0xee7c3c73, 0x5d4a14d9, 0xe864b7e3,
        0x42105d14, 0x203e13e0, 0x45eee2b6, 0xa3aaabea, 0xdb6c4f15,
        0xfacb4fd0, 0xc742f442, 0xef6abbb5, 0x654f3b1d, 0x41cd2105,
        0xd81e799e, 0x86854dc7, 0xe44b476a, 0x3d816250, 0xcf62a1f2,
        0x5b8d2646, 0xfc8883a0, 0xc1c7b6a3, 0x7f1524c3, 0x69cb7492,
        0x47848a0b, 0x5692b285, 0x095bbf00, 0xad19489d, 0x1462b174,
        0x23820e00, 0x58428d2a, 0x0c55f5ea, 0x1dadf43e, 0x233f7061,
        0x3372f092, 0x8d937e41, 0xd65fecf1, 0x6c223bdb, 0x7cde3759,
        0xcbee7460, 0x4085f2a7, 0xce77326e, 0xa6078084, 0x19f8509e,
        0xe8efd855, 0x61d99735, 0xa969a7aa, 0xc50c06c2, 0x5a04abfc,
        0x800bcadc, 0x9e447a2e, 0xc3453484, 0xfdd56705, 0x0e1e9ec9,
        0xdb73dbd3, 0x105588cd, 0x675fda79, 0xe3674340, 0xc5c43465,
        0x713e38d8, 0x3d28f89e, 0xf16dff20, 0x153e21e7, 0x8fb03d4a,
        0xe6e39f2b, 0xdb83adf7, 0xe93d5a68, 0x948140f7, 0xf64c261c,
        0x94692934, 0x411520f7, 0x7602d4f7, 0xbcf46b2e, 0xd4a20068,
        0xd4082471, 0x3320f46a, 0x43b7d4b7, 0x500061af, 0x1e39f62e,
        0x97244546, 0x14214f74, 0xbf8b8840, 0x4d95fc1d, 0x96b591af,
        0x70f4ddd3, 0x66a02f45, 0xbfbc09ec, 0x03bd9785, 0x7fac6dd0,
        0x31cb8504, 0x96eb27b3, 0x55fd3941, 0xda2547e6, 0xabca0a9a,
        0x28507825, 0x530429f4, 0x0a2c86da, 0xe9b66dfb, 0x68dc1462,
        0xd7486900, 0x680ec0a4, 0x27a18dee, 0x4f3ffea2, 0xe887ad8c,
        0xb58ce006, 0x7af4d6b6, 0xaace1e7c, 0xd3375fec, 0xce78a399,
        0x406b2a42, 0x20fe9e35, 0xd9f385b9, 0xee39d7ab, 0x3b124e8b,
        0x1dc9faf7, 0x4b6d1856, 0x26a36631, 0xeae397b2, 0x3a6efa74,
        0xdd5b4332, 0x6841e7f7, 0xca7820fb, 0xfb0af54e, 0xd8feb397,
        0x454056ac, 0xba489527, 0x55533a3a, 0x20838d87, 0xfe6ba9b7,
        0xd096954b, 0x55a867bc, 0xa1159a58, 0xcca92963, 0x99e1db33,
        0xa62a4a56, 0x3f3125f9, 0x5ef47e1c, 0x9029317c, 0xfdf8e802,
        0x04272f70, 0x80bb155c, 0x05282ce3, 0x95c11548, 0xe4c66d22,
        0x48c1133f, 0xc70f86dc, 0x07f9c9ee, 0x41041f0f, 0x404779a4,
        0x5d886e17, 0x325f51eb, 0xd59bc0d1, 0xf2bcc18f, 0x41113564,
        0x257b7834, 0x602a9c60, 0xdff8e8a3, 0x1f636c1b, 0x0e12b4c2,
        0x02e1329e, 0xaf664fd1, 0xcad18115, 0x6b2395e0, 0x333e92e1,
        0x3b240b62, 0xeebeb922, 0x85b2a20e, 0xe6ba0d99, 0xde720c8c,
        0x2da2f728, 0xd0127845, 0x95b794fd, 0x647d0862, 0xe7ccf5f0,
        0x5449a36f, 0x877d48fa, 0xc39dfd27, 0xf33e8d1e, 0x0a476341,
        0x992eff74, 0x3a6f6eab, 0xf4f8fd37, 0xa812dc60, 0xa1ebddf8,
        0x991be14c, 0xdb6e6b0d, 0xc67b5510, 0x6d672c37, 0x2765d43b,
        0xdcd0e804, 0xf1290dc7, 0xcc00ffa3, 0xb5390f92, 0x690fed0b,
        0x667b9ffb, 0xcedb7d9c, 0xa091cf0b, 0xd9155ea3, 0xbb132f88,
        0x515bad24, 0x7b9479bf, 0x763bd6eb, 0x37392eb3, 0xcc115979,
        0x8026e297, 0xf42e312d, 0x6842ada7, 0xc66a2b3b, 0x12754ccc,
        0x782ef11c, 0x6a124237, 0xb79251e7, 0x06a1bbe6, 0x4bfb6350,
        0x1a6b1018, 0x11caedfa, 0x3d25bdd8, 0xe2e1c3c9, 0x44421659,
        0x0a121386, 0xd90cec6e, 0xd5abea2a, 0x64af674e, 0xda86a85f,
        0xbebfe988, 0x64e4c3fe, 0x9dbc8057, 0xf0f7c086, 0x60787bf8,
        0x6003604d, 0xd1fd8346, 0xf6381fb0, 0x7745ae04, 0xd736fccc,
        0x83426b33, 0xf01eab71, 0xb0804187, 0x3c005e5f, 0x77a057be,
        0xbde8ae24, 0x55464299, 0xbf582e61, 0x4e58f48f, 0xf2ddfda2,
        0xf474ef38, 0x8789bdc2, 0x5366f9c3, 0xc8b38e74, 0xb475f255,
        0x46fcd9b9, 0x7aeb2661, 0x8b1ddf84, 0x846a0e79, 0x915f95e2,
        0x466e598e, 0x20b45770, 0x8cd55591, 0xc902de4c, 0xb90bace1,
        0xbb8205d0, 0x11a86248, 0x7574a99e, 0xb77f19b6, 0xe0a9dc09,
        0x662d09a1, 0xc4324633, 0xe85a1f02, 0x09f0be8c, 0x4a99a025,
        0x1d6efe10, 0x1ab93d1d, 0x0ba5a4df, 0xa186f20f, 0x2868f169,
        0xdcb7da83, 0x573906fe, 0xa1e2ce9b, 0x4fcd7f52, 0x50115e01,
        0xa70683fa, 0xa002b5c4, 0x0de6d027, 0x9af88c27, 0x773f8641,
        0xc3604c06, 0x61a806b5, 0xf0177a28, 0xc0f586e0, 0x006058aa,
        0x30dc7d62, 0x11e69ed7, 0x2338ea63, 0x53c2dd94, 0xc2c21634,
        0xbbcbee56, 0x90bcb6de, 0xebfc7da1, 0xce591d76, 0x6f05e409,
        0x4b7c0188, 0x39720a3d, 0x7c927c24, 0x86e3725f, 0x724d9db9,
        0x1ac15bb4, 0xd39eb8fc, 0xed545578, 0x08fca5b5, 0xd83d7cd3,
        0x4dad0fc4, 0x1e50ef5e, 0xb161e6f8, 0xa28514d9, 0x6c51133c,
        0x6fd5c7e7, 0x56e14ec4, 0x362abfce, 0xddc6c837, 0xd79a3234,
        0x92638212, 0x670efa8e, 0x406000e0, 0x3a39ce37, 0xd3faf5cf,
        0xabc27737, 0x5ac52d1b, 0x5cb0679e, 0x4fa33742, 0xd3822740,
        0x99bc9bbe, 0xd5118e9d, 0xbf0f7315, 0xd62d1c7e, 0xc700c47b,
        0xb78c1b6b, 0x21a19045, 0xb26eb1be, 0x6a366eb4, 0x5748ab2f,
        0xbc946e79, 0xc6a376d2, 0x6549c2c8, 0x530ff8ee, 0x468dde7d,
        0xd5730a1d, 0x4cd04dc6, 0x2939bbdb, 0xa9ba4650, 0xac9526e8,
        0xbe5ee304, 0xa1fad5f0, 0x6a2d519a, 0x63ef8ce2, 0x9a86ee22,
        0xc089c2b8, 0x43242ef6, 0xa51e03aa, 0x9cf2d0a4, 0x83c061ba,
        0x9be96a4d, 0x8fe51550, 0xba645bd6, 0x2826a2f9, 0xa73a3ae1,
        0x4ba99586, 0xef5562e9, 0xc72fefd3, 0xf752f7da, 0x3f046f69,
        0x77fa0a59, 0x80e4a915, 0x87b08601, 0x9b09e6ad, 0x3b3ee593,
        0xe990fd5a, 0x9e34d797, 0x2cf0b7d9, 0x022b8b51, 0x96d5ac3a,
        0x017da67d, 0xd1cf3ed6, 0x7c7d2d28, 0x1f9f25cf, 0xadf2b89b,
        0x5ad6b472, 0x5a88f54c, 0xe029ac71, 0xe019a5e6, 0x47b0acfd,
        0xed93fa9b, 0xe8d3c48d, 0x283b57cc, 0xf8d56629, 0x79132e28,
        0x785f0191, 0xed756055, 0xf7960e44, 0xe3d35e8c, 0x15056dd4,
        0x88f46dba, 0x03a16125, 0x0564f0bd, 0xc3eb9e15, 0x3c9057a2,
        0x97271aec, 0xa93a072a, 0x1b3f6d9b, 0x1e6321f5, 0xf59c66fb,
        0x26dcf319, 0x7533d928, 0xb155fdf5, 0x03563482, 0x8aba3cbb,
        0x28517711, 0xc20ad9f8, 0xabcc5167, 0xccad925f, 0x4de81751,
        0x3830dc8e, 0x379d5862, 0x9320f991, 0xea7a90c2, 0xfb3e7bce,
        0x5121ce64, 0x774fbe32, 0xa8b6e37e, 0xc3293d46, 0x48de5369,
        0x6413e680, 0xa2ae0810, 0xdd6db224, 0x69852dfd, 0x09072166,
        0xb39a460a, 0x6445c0dd, 0x586cdecf, 0x1c20c8ae, 0x5bbef7dd,
        0x1b588d40, 0xccd2017f, 0x6bb4e3bb, 0xdda26a7e, 0x3a59ff45,
        0x3e350a44, 0xbcb4cdd5, 0x72eacea8, 0xfa6484bb, 0x8d6612ae,
        0xbf3c6f47, 0xd29be463, 0x542f5d9e, 0xaec2771b, 0xf64e6370,
        0x740e0d8d, 0xe75b1357, 0xf8721671, 0xaf537d5d, 0x4040cb08,
        0x4eb4e2cc, 0x34d2466a, 0x0115af84, 0xe1b00428, 0x95983a1d,
        0x06b89fb4, 0xce6ea048, 0x6f3f3b82, 0x3520ab82, 0x011a1d4b,
        0x277227f8, 0x611560b1, 0xe7933fdc, 0xbb3a792b, 0x344525bd,
        0xa08839e1, 0x51ce794b, 0x2f32c9b7, 0xa01fbac9, 0xe01cc87e,
        0xbcc7d1f6, 0xcf0111c3, 0xa1e8aac7, 0x1a908749, 0xd44fbd9a,
        0xd0dadecb, 0xd50ada38, 0x0339c32a, 0xc6913667, 0x8df9317c,
        0xe0b12b4f, 0xf79e59b7, 0x43f5bb3a, 0xf2d519ff, 0x27d9459c,
        0xbf97222c, 0x15e6fc2a, 0x0f91fc71, 0x9b941525, 0xfae59361,
        0xceb69ceb, 0xc2a86459, 0x12baa8d1, 0xb6c1075e, 0xe3056a0c,
        0x10d25065, 0xcb03a442, 0xe0ec6e0e, 0x1698db3b, 0x4c98a0be,
        0x3278e964, 0x9f1f9532, 0xe0d392df, 0xd3a0342b, 0x8971f21e,
        0x1b0a7441, 0x4ba3348c, 0xc5be7120, 0xc37632d8, 0xdf359f8d,
        0x9b992f2e, 0xe60b6f47, 0x0fe3f11d, 0xe54cda54, 0x1edad891,
        0xce6279cf, 0xcd3e7e6f, 0x1618b166, 0xfd2c1d05, 0x848fd2c5,
        0xf6fb2299, 0xf523f357, 0xa6327623, 0x93a83531, 0x56cccd02,
        0xacf08162, 0x5a75ebb5, 0x6e163697, 0x88d273cc, 0xde966292,
        0x81b949d0, 0x4c50901b, 0x71c65614, 0xe6c6c7bd, 0x327a140a,
        0x45e1d006, 0xc3f27b9a, 0xc9aa53fd, 0x62a80f00, 0xbb25bfe2,
        0x35bdd2f6, 0x71126905, 0xb2040222, 0xb6cbcf7c, 0xcd769c2b,
        0x53113ec0, 0x1640e3d3, 0x38abbd60, 0x2547adf0, 0xba38209c,
        0xf746ce76, 0x77afa1c5, 0x20756060, 0x85cbfe4e, 0x8ae88dd8,
        0x7aaaf9b0, 0x4cf9aa7e, 0x1948c25c, 0x02fb8a8c, 0x01c36ae4,
        0xd6ebe1f9, 0x90d4f869, 0xa65cdea0, 0x3f09252d, 0xc208e69f,
        0xb74e6132, 0xce77e25b, 0x578fdfe3, 0x3ac372e6
    ];

    /**
     * @type {Array.<number>}
     * @const
     * @inner
     */
    var C_ORIG = [
        0x4f727068, 0x65616e42, 0x65686f6c, 0x64657253, 0x63727944,
        0x6f756274
    ];

    /**
     * @param {Array.<number>} lr
     * @param {number} off
     * @param {Array.<number>} P
     * @param {Array.<number>} S
     * @returns {Array.<number>}
     * @inner
     */
    function _encipher(lr, off, P, S) { // This is our bottleneck: 1714/1905 ticks / 90% - see profile.txt
        var n,
            l = lr[off],
            r = lr[off + 1];

        l ^= P[0];
        for (var i=0, k=BLOWFISH_NUM_ROUNDS-2; i<=k;)
            // Feistel substitution on left word
            n  = S[(l >> 24) & 0xff],
                n += S[0x100 | ((l >> 16) & 0xff)],
                n ^= S[0x200 | ((l >> 8) & 0xff)],
                n += S[0x300 | (l & 0xff)],
                r ^= n ^ P[++i],
                // Feistel substitution on right word
                n  = S[(r >> 24) & 0xff],
                n += S[0x100 | ((r >> 16) & 0xff)],
                n ^= S[0x200 | ((r >> 8) & 0xff)],
                n += S[0x300 | (r & 0xff)],
                l ^= n ^ P[++i];
        lr[off] = r ^ P[BLOWFISH_NUM_ROUNDS + 1];
        lr[off + 1] = l;
        return lr;
    }

    /**
     * @param {Array.<number>} data
     * @param {number} offp
     * @returns {{key: number, offp: number}}
     * @inner
     */
    function _streamtoword(data, offp) {
        for (var i = 0, word = 0; i < 4; ++i)
            word = (word << 8) | (data[offp] & 0xff),
                offp = (offp + 1) % data.length;
        return { key: word, offp: offp };
    }

    /**
     * @param {Array.<number>} key
     * @param {Array.<number>} P
     * @param {Array.<number>} S
     * @inner
     */
    function _key(key, P, S) {
        var offset = 0,
            lr = [0, 0],
            plen = P.length,
            slen = S.length,
            sw;
        for (var i = 0; i < plen; i++)
            sw = _streamtoword(key, offset),
                offset = sw.offp,
                P[i] = P[i] ^ sw.key;
        for (i = 0; i < plen; i += 2)
            lr = _encipher(lr, 0, P, S),
                P[i] = lr[0],
                P[i + 1] = lr[1];
        for (i = 0; i < slen; i += 2)
            lr = _encipher(lr, 0, P, S),
                S[i] = lr[0],
                S[i + 1] = lr[1];
    }

    /**
     * Expensive key schedule Blowfish.
     * @param {Array.<number>} data
     * @param {Array.<number>} key
     * @param {Array.<number>} P
     * @param {Array.<number>} S
     * @inner
     */
    function _ekskey(data, key, P, S) {
        var offp = 0,
            lr = [0, 0],
            plen = P.length,
            slen = S.length,
            sw;
        for (var i = 0; i < plen; i++)
            sw = _streamtoword(key, offp),
                offp = sw.offp,
                P[i] = P[i] ^ sw.key;
        offp = 0;
        for (i = 0; i < plen; i += 2)
            sw = _streamtoword(data, offp),
                offp = sw.offp,
                lr[0] ^= sw.key,
                sw = _streamtoword(data, offp),
                offp = sw.offp,
                lr[1] ^= sw.key,
                lr = _encipher(lr, 0, P, S),
                P[i] = lr[0],
                P[i + 1] = lr[1];
        for (i = 0; i < slen; i += 2)
            sw = _streamtoword(data, offp),
                offp = sw.offp,
                lr[0] ^= sw.key,
                sw = _streamtoword(data, offp),
                offp = sw.offp,
                lr[1] ^= sw.key,
                lr = _encipher(lr, 0, P, S),
                S[i] = lr[0],
                S[i + 1] = lr[1];
    }

    /**
     * Internaly crypts a string.
     * @param {Array.<number>} b Bytes to crypt
     * @param {Array.<number>} salt Salt bytes to use
     * @param {number} rounds Number of rounds
     * @param {function(Error, Array.<number>=)=} callback Callback receiving the error, if any, and the resulting bytes. If
     *  omitted, the operation will be performed synchronously.
     *  @param {function(number)=} progressCallback Callback called with the current progress
     * @returns {!Array.<number>|undefined} Resulting bytes if callback has been omitted, otherwise `undefined`
     * @inner
     */
    function _crypt(b, salt, rounds, callback, progressCallback) {
        var cdata = C_ORIG.slice(),
            clen = cdata.length,
            err;

        // Validate
        if (rounds < 4 || rounds > 31) {
            err = Error("Illegal number of rounds (4-31): "+rounds);
            if (callback) {
                nextTick(callback.bind(this, err));
                return;
            } else
                throw err;
        }
        if (salt.length !== BCRYPT_SALT_LEN) {
            err =Error("Illegal salt length: "+salt.length+" != "+BCRYPT_SALT_LEN);
            if (callback) {
                nextTick(callback.bind(this, err));
                return;
            } else
                throw err;
        }
        rounds = (1 << rounds) >>> 0;
        var P = P_ORIG.slice(),
            S = S_ORIG.slice(),
            i = 0, j;
        _ekskey(salt, b, P, S);

        /**
         * Calcualtes the next round.
         * @returns {Array.<number>|undefined} Resulting array if callback has been omitted, otherwise `undefined`
         * @inner
         */
        function next() {
            if (progressCallback)
                progressCallback(i / rounds);
            if (i < rounds) {
                var start = Date.now();
                for (; i < rounds;) {
                    i = i + 1;
                    _key(b, P, S);
                    _key(salt, P, S);
                    if (Date.now() - start > MAX_EXECUTION_TIME)
                        break;
                }
            } else {
                for (i = 0; i < 64; i++)
                    for (j = 0; j < (clen >> 1); j++)
                        _encipher(cdata, j << 1, P, S);
                var ret = [];
                for (i = 0; i < clen; i++)
                    ret.push(((cdata[i] >> 24) & 0xff) >>> 0),
                        ret.push(((cdata[i] >> 16) & 0xff) >>> 0),
                        ret.push(((cdata[i] >> 8) & 0xff) >>> 0),
                        ret.push((cdata[i] & 0xff) >>> 0);
                if (callback) {
                    callback(null, ret);
                    return;
                } else
                    return ret;
            }
            if (callback)
                nextTick(next);
        }

        // Async
        if (typeof callback !== 'undefined') {
            next();

            // Sync
        } else {
            var res;
            while (true)
                if (typeof(res = next()) !== 'undefined')
                    return res || [];
        }
    }

    /**
     * Internally hashes a string.
     * @param {string} s String to hash
     * @param {?string} salt Salt to use, actually never null
     * @param {function(Error, string=)=} callback Callback receiving the error, if any, and the resulting hash. If omitted,
     *  hashing is perormed synchronously.
     *  @param {function(number)=} progressCallback Callback called with the current progress
     * @returns {string|undefined} Resulting hash if callback has been omitted, otherwise `undefined`
     * @inner
     */
    function _hash(s, salt, callback, progressCallback) {
        var err;
        if (typeof s !== 'string' || typeof salt !== 'string') {
            err = Error("Invalid string / salt: Not a string");
            if (callback) {
                nextTick(callback.bind(this, err));
                return;
            }
            else
                throw err;
        }

        // Validate the salt
        var minor, offset;
        if (salt.charAt(0) !== '$' || salt.charAt(1) !== '2') {
            err = Error("Invalid salt version: "+salt.substring(0,2));
            if (callback) {
                nextTick(callback.bind(this, err));
                return;
            }
            else
                throw err;
        }
        if (salt.charAt(2) === '$')
            minor = String.fromCharCode(0),
                offset = 3;
        else {
            minor = salt.charAt(2);
            if ((minor !== 'a' && minor !== 'b' && minor !== 'y') || salt.charAt(3) !== '$') {
                err = Error("Invalid salt revision: "+salt.substring(2,4));
                if (callback) {
                    nextTick(callback.bind(this, err));
                    return;
                } else
                    throw err;
            }
            offset = 4;
        }

        // Extract number of rounds
        if (salt.charAt(offset + 2) > '$') {
            err = Error("Missing salt rounds");
            if (callback) {
                nextTick(callback.bind(this, err));
                return;
            } else
                throw err;
        }
        var r1 = parseInt(salt.substring(offset, offset + 1), 10) * 10,
            r2 = parseInt(salt.substring(offset + 1, offset + 2), 10),
            rounds = r1 + r2,
            real_salt = salt.substring(offset + 3, offset + 25);
        s += minor >= 'a' ? "\x00" : "";

        var passwordb = stringToBytes(s),
            saltb = base64_decode(real_salt, BCRYPT_SALT_LEN);

        /**
         * Finishes hashing.
         * @param {Array.<number>} bytes Byte array
         * @returns {string}
         * @inner
         */
        function finish(bytes) {
            var res = [];
            res.push("$2");
            if (minor >= 'a')
                res.push(minor);
            res.push("$");
            if (rounds < 10)
                res.push("0");
            res.push(rounds.toString());
            res.push("$");
            res.push(base64_encode(saltb, saltb.length));
            res.push(base64_encode(bytes, C_ORIG.length * 4 - 1));
            return res.join('');
        }

        // Sync
        if (typeof callback == 'undefined')
            return finish(_crypt(passwordb, saltb, rounds));

        // Async
        else {
            _crypt(passwordb, saltb, rounds, function(err, bytes) {
                if (err)
                    callback(err, null);
                else
                    callback(null, finish(bytes));
            }, progressCallback);
        }
    }

    /**
     * Encodes a byte array to base64 with up to len bytes of input, using the custom bcrypt alphabet.
     * @function
     * @param {!Array.<number>} b Byte array
     * @param {number} len Maximum input length
     * @returns {string}
     * @expose
     */
    bcrypt.encodeBase64 = base64_encode;

    /**
     * Decodes a base64 encoded string to up to len bytes of output, using the custom bcrypt alphabet.
     * @function
     * @param {string} s String to decode
     * @param {number} len Maximum output length
     * @returns {!Array.<number>}
     * @expose
     */
    bcrypt.decodeBase64 = base64_decode;

    return bcrypt;
}));
