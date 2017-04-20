package ecdh

import (
	"math/rand"
	"time"
)

var rng_psize int = 256
var rng_state *ArcFour
var rng_pptr int
var rng_pool map[int]byte

func rng_seed_int(x int) {
	rng_pool[rng_pptr] ^= byte(x & 255)
	rng_pptr++
	rng_pool[rng_pptr] ^= byte((x >> uint(8)) & 255)
	rng_pptr++
	rng_pool[rng_pptr] ^= byte((x >> uint(16)) & 255)
	rng_pptr++
	rng_pool[rng_pptr] ^= byte((x >> uint(24)) & 255)
	rng_pptr++
	if rng_pptr >= rng_psize {
		rng_pptr -= rng_psize
	}
}

func rng_seed_time() {
	rng_seed_int(time.Now().Nanosecond())
}

func init_pool() {
	rng_pool = make(map[int]byte)
	rng_pptr = 0

	rand.Seed(time.Now().Unix())
	for rng_pptr < rng_psize {
		var r int = rand.Int()
		rng_pool[rng_pptr] = byte(r >> uint(8))
		rng_pptr++
		rng_pool[rng_pptr] = byte(r & 255)
		rng_pptr++
	}

	rng_pptr = 0
	rng_seed_time()
}

func rng_get_byte() byte {
	if rng_state == nil {
		init_pool()
		rng_state = NewArcFour()
		rng_state.Init(rng_pool)
		for rng_pptr = 0; rng_pptr < len(rng_pool); rng_pptr++ {
			rng_pool[rng_pptr] = 0
		}
		rng_pptr = 0
	}
	return rng_state.Next()
}

type SecureRandom struct {
}

func (this *SecureRandom) Next(ba []byte) {
	for i := 0; i < len(ba); i++ {
		ba[i] = rng_get_byte()
	}
}
