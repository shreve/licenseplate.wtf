package model

import (
	"crypto/md5"
	_ "embed"
	"encoding/binary"
	"strings"
)

//go:embed data/adjectives.txt
var adjectiveStr string
var adjectives []string
var adjectiveLen uint64

//go:embed data/animals.txt
var animalStr string
var animals []string
var animalLen uint64

func init() {
	adjectives = strings.Split(adjectiveStr, "\n")
	adjectiveLen = uint64(len(adjectives))

	animals = strings.Split(animalStr, "\n")
	animalLen = uint64(len(animals))
}

// Take in some string and generate a name
// 1. Gen the md5 to diffuse some values
// 2. Split the bytes in half
// 3. Modulo the resulting numbers against list sizes
// 4. Format the strings
func NameHash(input string) string {
	bytes := md5.Sum([]byte(input))

	adji := binary.BigEndian.Uint64(bytes[0:8]) % adjectiveLen
	anmi := binary.BigEndian.Uint64(bytes[8:16]) % animalLen

	return adjectives[adji] + " " + animals[anmi]
}
