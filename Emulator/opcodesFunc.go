package emulator

import (
	"math/rand"
	"time"
)

func (c *Cpu) StackPush(address uint16) {
	// Vérifiez que le pointeur de pile (SP) est dans la plage valide (0-15).
	c.Stack[c.Sp] = address
	c.Sp++
}

func (c *Cpu) StackPop() uint16 {
	c.Sp--
	return c.Stack[c.Sp]
}

// Opcode 00E0 - Effacer l'écran =
// Clear the display.
func (c *Cpu) op00E0() {
	for x := 0; x < 64; x++ {
		for y := 0; y < 32; y++ {

			if x > 64 || y > 31 {
				x = x % 64
				y = y % 32

			}
			c.Gfx[x][y] = 0
		}
	}
}

// stock les données
func (c *Cpu) op6XNN(opcodeX, opcodeNNN byte) {
	c.Registre[opcodeX] = opcodeNNN
}

// Opcode 00EE - Retour de sous-routine =
// Return from a subroutine.The interpreter sets the program counter to the address at the top of the stack,
// then subtracts 1 from the stack pointer
func (c *Cpu) op00EE() {
	c.Pc = c.StackPop()
}

// Opcode 1NNN - Saut
// Jump to location nnn. The interpreter sets the program counter to nnn.
func (c *Cpu) op1nnn(address uint16) {
	c.Pc = address - 2
}

// Opcode 2NNN - Appel de sous-routine =
// Call subroutine at nnn. The interpreter increments the stack pointer, then puts the current PC on the top
// of the stack. The PC is then set to nnn.
func (c *Cpu) op2nnn(address uint16) {
	// Vérifie que le pointeur de pile (SP) est dans la plage valide (0-15).
	c.StackPush(c.Pc)
	c.Pc = address - 2
}

// Opcode 3XNN - Saut conditionnel (égal) =
// Skip next instruction if Vx = kk. The interpreter compares register Vx to kk, and if they are equal,
// increments the program counter by 2.
func (c *Cpu) op3nnn(opcodeX, opcodeNNN byte) {
	if c.Registre[opcodeX] == opcodeNNN {
		c.Pc += 2
	}
}

// Opcode 4XNN - Saut conditionnel (différent)
// Skip next instruction if Vx != kk. The interpreter compares register Vx to kk, and if they are not equal,
// increments the program counter by 2.
func (c *Cpu) op4nnn(opcodeX, opcodeNN byte) {
	if c.Registre[opcodeX] != opcodeNN {
		c.Pc += 2
	}
}

// Opcode 5XY0 - Saut conditionnel (égalité de registres)
// Skip next instruction if Vx = Vy. The interpreter compares register Vx to register Vy, and if they are equal,
// increments the program counter by 2.
func (c *Cpu) op5nnn(opcodeX, opcodeY byte) {
	if c.Registre[opcodeX] == c.Registre[opcodeY] {
		c.Pc += 2
	}
}

// Opcode 6XNN - Chargement de valeur constante
// Set Vx = kk. The interpreter puts the value kk into register Vx
func (c *Cpu) op6nnn(opcodeX, opcodeNN byte) {
	c.Registre[opcodeX] = opcodeNN
}

// Opcode 7XNN - Ajout de valeur constante =
// Set Vx = Vx + kk. Adds the value kk to the value of register Vx, then stores the result in Vx.
func (c *Cpu) op7nnn(opcodeX, opcodeNN byte) {
	c.Registre[opcodeX] = c.Registre[opcodeX] + opcodeNN
}

// Opcode 8XY0 - Copie de Registre =
// Set Vx = Vy. Stores the value of register Vy in register Vx.
func (c *Cpu) op8nn0(opcodeX, opcodeY byte) {
	c.Registre[opcodeX] = c.Registre[opcodeY]
}

// Opcode 8XY1 - Opération OU (bitwise OR) =
// Set Vx = Vx OR Vy. Performs a bitwise OR on the values of Vx and Vy, then stores the result in Vx. A
// bit wise OR compares the corresponding bits from two values, and if either bit is 1, then the same bit in the
// result is also 1. Otherwise, it is 0.
func (c *Cpu) op8xy1(x, y byte) {
	c.Registre[x] |= c.Registre[y]
	c.Registre[0xF] = 0
}

// Opcode 8XY2 - Opération ET (bitwise AND) =
// Set Vx = Vx AND Vy. Performs a bitwise AND on the values of Vx and Vy, then stores the result in Vx.
// A bitwise AND compares the corresponding bits from two values, and if both bits are 1, then the same bit
// in the result is also 1. Otherwise, it is 0.
func (c *Cpu) op8xy2(opcodeX, opcodeY byte) {
	c.Registre[opcodeX] &= c.Registre[opcodeY]
	c.Registre[0xF] = 0
}

// Opcode 8XY3 - Opération XOR (bitwise XOR) =
// Set Vx = Vx XOR Vy. Performs a bitwise exclusive OR on the values of Vx and Vy, then stores the result
// in Vx. An exclusive OR compares the corresponding bits from two values, and if the bits are not both the
// same, then the corresponding bit in the result is set to 1. Otherwise, it is 0.
func (c *Cpu) op8xy3(opcodeX, opcodeY byte) {
	c.Registre[opcodeX] ^= c.Registre[opcodeY]
	c.Registre[0xF] = 0
}

// Opcode 8XY4 - Ajout avec retenue =
//Set Vx = Vx + Vy, set VF = carry. The values of Vx and Vy are added together. If the result is greater
//than 8 bits (i.e., ¿ 255,) VF is set to 1, otherwise 0. Only the lowest 8 bits of the result are kept, and stored
//in Vx.

// Vx += Vy
func (c *Cpu) op8xy4(opcodeX, opcodeY byte) {

	var a uint16 = uint16(c.Registre[opcodeX])
	var b uint16 = uint16(c.Registre[opcodeY])

	final := a + b
	if final > 255 {
		c.Registre[0xF] = 0x1
	} else {
		c.Registre[0xF] = 0x0
	}
	if opcodeX != 0xF {
		c.Registre[opcodeX] = byte(final)
	}
}

// Opcode 8XY5 - Soustraction avec retenue
// Vx -= Vy
// Set Vx = Vx - Vy, set VF = NOT borrow. If Vx ¿ Vy, then VF is set to 1, otherwise 0. Then Vy is
// subtracted from Vx, and the results stored in Vx.
func (c *Cpu) op8xy5(opcodeX, opcodeY byte) {
	var X = int16(c.Registre[opcodeX])
	var Y = int16(c.Registre[opcodeY])
	res := X - Y
	if res < 0 {
		c.Registre[0xF] = 0
	} else {
		c.Registre[0xF] = 1
	}
	if opcodeX != 0xF {
		c.Registre[opcodeX] = byte(res)
	}
}

// Opcode 8XY6 - Décalage à droite
// Set Vx = Vx SHR 1. If the least-significant bit of Vx is 1, then VF is set to 1, otherwise 0. Then Vx is
// divided 	by 2
func (c *Cpu) op8xy6(opcodeX, opcodeY byte) {
	if (c.Registre[opcodeX] % 2) == 1 {
		c.Registre[0xF] = 1
	} else {
		c.Registre[0xF] = 0
	}
	if opcodeX != 0xF {
		// c.Registre[opcodeX] = c.Registre[opcodeX] / 2
		c.Registre[opcodeX] = c.Registre[opcodeY] / 2
	}
	// c.Registre[opcodeX]++
}

// Opcode 8XY7 - Soustraction inversée avec retenue =
// Set Vx = Vy - Vx, set VF = NOT borrow. If Vy ¿ Vx, then VF is set to 1, otherwise 0. Then Vx is
// subtracted from Vy, and the results stored in Vx.
func (c *Cpu) op8xy7(opcodeX, opcodeY byte) {
	c.Registre[opcodeX] = c.Registre[opcodeY] - c.Registre[opcodeX]
	if c.Registre[opcodeY] >= c.Registre[opcodeX] {
		c.Registre[0xF] = 1
	} else {
		c.Registre[0xF] = 0
	}
}

// Opcode 8XYE - Décalage à gauche =
// Set Vx = Vx SHL 1. If the most-significant bit of Vx is 1, then VF is set to 1, otherwise to 0. Then Vx is
// multiplied by 2
func (c *Cpu) op8xyE(opcodeX, opcodeY byte) {
	if (c.Registre[opcodeY] >> 7) == 1 {
		c.Registre[0xF] = 1
	} else {
		c.Registre[0xF] = 0
	}
	if opcodeX != 0xF {
		// c.Registre[opcodeX] = c.Registre[opcodeX] * 2
		c.Registre[opcodeX] = c.Registre[opcodeY] * 2
	}
	// c.Registre[opcodeX]++
}

// Opcode 9XY0 - Saut conditionnel (différents registres)=
// Skip next instruction if Vx != Vy. The values of Vx and Vy are compared, and if they are not equal, the
// program counter is increased by 2
func (c *Cpu) op9nn0(opcodeX, opcodeY byte) {
	if c.Registre[opcodeX] != c.Registre[opcodeY] {
		c.Pc += 2
	}
}

// Opcode ANNN - Chargement de l'index (I) =
// Set I = nnn. The value of register I is set to nnn.
func (c *Cpu) opAnnn(address uint16) { // verifier si nnn = opcodennn ou 0
	c.I = address
}

// Opcode BNNN - Saut avec offset =
// Jump to location nnn + V0. The program counter is set to nnn plus the value of V0
func (c *Cpu) opBnnn(address uint16) {
	c.Pc = address + uint16(c.Registre[4])
}

// Opcode CXNN - Génération d'un nombre aléatoire (0 à 255) =
// Set Vx = random byte AND kk. The interpreter generates a random number from 0 to 255, which is then
// ANDed with the value kk. The results are stored in Vx. See instruction 8xy2 for more information on AND.
func (c *Cpu) opCxkk(opcodeX, opcodeNN byte) {
	c.Registre[opcodeX] = byte(rand.Int()*256) & opcodeNN
}

// Opcode DXYN - Dessin à l'écran
func (c *Cpu) opDxyn(opcodeX, opcodeY, opcodeN byte) {
	xval := c.Registre[opcodeX]
	yval := c.Registre[opcodeY]
	c.Registre[0xF] = 0
	var i byte = 0
	for ; i < opcodeN; i++ {
		row := c.Memory[c.I+uint16(i)]
		if erased := c.DrawSprite(xval, yval+i, row); erased {
			c.Registre[0xF] = 1
		}
	}
}

// Opcode FX07 - Chargement du retard =
// Set Vx = delay timer value. The value of DT is placed into Vx.
func (c *Cpu) opFx07(opcodeX byte) {
	c.Registre[opcodeX] = c.Delay_timer
}

// Opcode FX15 - Réglage du retard =
// Set delay timer = Vx. Delay Timer is set equal to the value of Vx.
func (c *Cpu) opFx15(opcodeX byte) {
	c.Delay_timer = c.Registre[opcodeX]
}

// Opcode FX55 - Sauvegarde des registres
// Stores V0 to VX in memory starting at address I. I is then set to I + x + 1.
func (c *Cpu) opFx55(opcodeX byte) {
	for i := byte(0); i <= opcodeX; i++ {
		c.Memory[c.I+uint16(i)] = c.Registre[i]
	}
	c.I++

}

// Opcode FX29 - Chargement de l'emplacement du caractère
// Set I = location of sprite for digit Vx. The value of I is set to the location for the hexadecimal sprite
// corresponding to the value of Vx. See section 2.4, Display, for more information on the Chip-8 hexadecimal
// font. To obtain this value, multiply VX by 5 (all font data stored in first 80 bytes of memory).
func (c *Cpu) opFx29(opcodeX byte) {
	c.I = uint16(c.Registre[opcodeX]) * uint16(5)
}

// Fills V0 to VX with values from memory starting at address I. I is then set to I + x + 1.
func (c *Cpu) opFx65(opcodeX byte) {
	for i := byte(0); i <= opcodeX; i++ {
		c.Registre[i] = c.Memory[c.I+uint16(i)]
	}
	c.I++
}

// Opcode FX18 - Réglage du son =
// Set sound timer = Vx. Sound Timer is set equal to the value of Vx.
func (c *Cpu) opFx18(opcodeX byte) {
	c.Sound_timer = c.Registre[opcodeX]
}

// Opcode FX33 - Chargement des chiffres décimaux
// Store BCD representation of Vx in memory locations I, I+1, and I+2. The interpreter takes the decimal
// value of Vx, and places the hundreds digit in memory at location in I, the tens digit at location I+1, and
// the ones digit at location I+2.
func (c *Cpu) opFx33(opcodeX byte) {
	value := c.Registre[opcodeX]
	// Pour obtenir les chiffres individuels
	hundreds := value / 100
	value %= 100
	tens := value / 10
	ones := value % 10
	// Stock les chiffres décimaux dans la mémoire à partir de l'adresse I
	c.Memory[c.I] = hundreds
	c.Memory[c.I+1] = tens
	c.Memory[c.I+2] = ones
}

// Set I = I + Vx. The values of I and Vx are added, and the results are stored in I
func (c *Cpu) opFx1E(opcodeX byte) {
	c.I += uint16(c.Registre[opcodeX])
}

// Opcode EX9E - Saut si touche pressée
// Skip next instruction if key with the value of Vx is pressed. Checks the keyboard, and if the key corresponding
// to the value of Vx is currently in the down position, PC is increased by 2
func (c *Cpu) opEx9E(opcodeX byte) {
	if c.Key[c.Registre[opcodeX]] {
		c.Pc += 2
	}
}

// Opcode EXA1 - Saut si touche non pressée
// Skip next instruction if key with the value of Vx is not pressed. Checks the keyboard, and if the key
// corresponding to the value of Vx is currently in the up position, PC is increased by 2.
func (c *Cpu) opExA1(opcodeX byte) {
	if !c.Key[c.Registre[opcodeX]] {
		c.Pc += 2
	}
}

// Opcode FX0A - Attente de touche
// Wait for a key press, store the value of the key in Vx. All execution stops until a key is pressed, then the
// value of that key is stored in Vx.
func (c *Cpu) opFx0A(opcodeX byte) {
	c.WaitForKey = true
	c.GetKey()
	for i := 0; i < len(c.Key); i++ {
		if c.Key[i] {
			time.Sleep(time.Second)
			c.Registre[i] = opcodeX
			c.WaitForKey = false
		}
	}
	if c.WaitForKey {
		c.Pc -= 2
	}
}
