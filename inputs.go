package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	pngembed "github.com/sabhiram/png-embed"
)

func inputs(g *Game) {
	if contains(inpututil.PressedKeys(), ebiten.KeyS) && contains(inpututil.PressedKeys(), ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyS) {
		savePNGFile(g)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyO) {
		extractPNG()
	}
}

func contains(s []ebiten.Key, e ebiten.Key) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func savePNGFile(g *Game) {
	println("saved")
	g.Canvas.GetLevel().ExportPNG()
}

func extractPNG() {
	println("opened")
	bs, _ := ioutil.ReadFile("sample.png")
	fileBytes, err := pngembed.Extract(bs)
	if err != nil {
		println(err.Error())
		return
	}

	data, err := base64.StdEncoding.DecodeString(string(fileBytes["FOO"]))
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("%q\n", data)
}
