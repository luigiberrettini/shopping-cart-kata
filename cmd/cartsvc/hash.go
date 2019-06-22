package main

import (
	"fmt"
	"strconv"
)

func (a *App) encode(id int64) (string, error) {
	return a.HashGen.EncodeHex(fmt.Sprintf("%x", id))
}

func (a *App) decode(hash string) (int64, error) {
	s, err := a.HashGen.DecodeHex(hash)
	if err != nil {
		return 0, err
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return int64(0), err
	}
	return int64(i), nil
}
