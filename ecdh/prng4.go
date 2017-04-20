package ecdh

import ()

type ArcFour struct {
	I int
	J int
	S map[int]byte
}

func NewArcFour() *ArcFour {
	return &ArcFour{I: 0, J: 0, S: make(map[int]byte)}
}

func (this *ArcFour) Init(key map[int]byte) {
	var i int = 0
	var j int = 0
	var t byte = 0
	for i = 0; i < 256; i++ {
		this.S[i] = byte(i)
	}
	for i = 0; i < 256; i++ {
		j = (j + int(this.S[i]) + int(key[i%len(key)])) & 255
		t = this.S[i]
		this.S[i] = this.S[j]
		this.S[j] = t
	}
	this.I = 0
	this.J = 0
}

func (this *ArcFour) Next() byte {
	var t byte
	this.I = (this.I + 1) & 255
	this.J = (this.J + int(this.S[this.I])) & 255
	t = this.S[this.I]
	this.S[this.I] = this.S[this.J]
	this.S[this.J] = t
	return this.S[int((t+this.S[this.I])&255)]
}
