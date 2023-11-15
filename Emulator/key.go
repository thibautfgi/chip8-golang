package emulator

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func (c *Cpu) GetKey() {
	// test := []string{"  0  ", "  1  ", "  2  ", "  3  ", "  4  ", "  5  ", "  6  ", "  7  ", "  8  ", "  9  ", "  A  ", "  B  ", "  C  ", "  D  ", "  E  ", "  F  "}
	c.Key = [16]bool{
		ebiten.IsKeyPressed(ebiten.Key0), //0
		ebiten.IsKeyPressed(ebiten.Key1), //1
		ebiten.IsKeyPressed(ebiten.Key2), //2
		ebiten.IsKeyPressed(ebiten.Key3), //3
		ebiten.IsKeyPressed(ebiten.Key4), //4
		ebiten.IsKeyPressed(ebiten.Key5), //5
		ebiten.IsKeyPressed(ebiten.Key6), //6
		ebiten.IsKeyPressed(ebiten.Key7), //7
		ebiten.IsKeyPressed(ebiten.Key8), //8
		ebiten.IsKeyPressed(ebiten.Key9), //9
		ebiten.IsKeyPressed(ebiten.KeyQ), //A
		ebiten.IsKeyPressed(ebiten.KeyB), //B
		ebiten.IsKeyPressed(ebiten.KeyC), //C
		ebiten.IsKeyPressed(ebiten.KeyD), //D
		ebiten.IsKeyPressed(ebiten.KeyE), //E
		ebiten.IsKeyPressed(ebiten.KeyF), //F
	}
}
