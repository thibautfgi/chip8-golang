package emulator

import (
	"time"
)

type Cpu struct {
	Memory    [4096]byte
	Registre  [16]byte
	I         uint16
	Pc        uint16
	Gfx       [64][32]byte
	Stack     [16]uint16
	Sp        byte
	Romlength uint16

	Delay_timer byte
	Sound_timer byte
	timeStart   time.Time

	Key        [16]bool
	WaitForKey bool

	Opcode uint16
}

func (cpu *Cpu) initialiseFont() {
	//0x000-0x1FF - Chip 8 interpreter (contains font set in emu)
	//0x050-0x0A0 - Used for the built in 4x5 pixel font set (0-F)
	//0 0x050
	cpu.Memory[0x050] = 0xF0
	cpu.Memory[0x051] = 0x90
	cpu.Memory[0x052] = 0x90
	cpu.Memory[0x053] = 0x90
	cpu.Memory[0x054] = 0xF0
	// 1
	cpu.Memory[0x055] = 0x20
	cpu.Memory[0x056] = 0x60
	cpu.Memory[0x057] = 0x20
	cpu.Memory[0x058] = 0x20
	cpu.Memory[0x059] = 0x70
	// 2
	cpu.Memory[0x05A] = 0xF0
	cpu.Memory[0x05B] = 0x10
	cpu.Memory[0x05C] = 0xF0
	cpu.Memory[0x05D] = 0x80
	cpu.Memory[0x05E] = 0xF0
	// 3
	cpu.Memory[0x05F] = 0xF0
	cpu.Memory[0x060] = 0x10
	cpu.Memory[0x061] = 0xF0
	cpu.Memory[0x062] = 0x10
	cpu.Memory[0x063] = 0xF0
	// 4
	cpu.Memory[0x064] = 0x90
	cpu.Memory[0x065] = 0x90
	cpu.Memory[0x066] = 0xF0
	cpu.Memory[0x067] = 0x10
	cpu.Memory[0x068] = 0x10
	// 5
	cpu.Memory[0x069] = 0xF0
	cpu.Memory[0x06A] = 0x80
	cpu.Memory[0x06B] = 0xF0
	cpu.Memory[0x06C] = 0x10
	cpu.Memory[0x06D] = 0xF0
	// 6
	cpu.Memory[0x06E] = 0xF0
	cpu.Memory[0x06F] = 0x80
	cpu.Memory[0x070] = 0xF0
	cpu.Memory[0x071] = 0x90
	cpu.Memory[0x072] = 0xF0
	// 7
	cpu.Memory[0x073] = 0xF0
	cpu.Memory[0x074] = 0x10
	cpu.Memory[0x075] = 0x20
	cpu.Memory[0x076] = 0x40
	cpu.Memory[0x077] = 0x40
	// 8
	cpu.Memory[0x078] = 0xF0
	cpu.Memory[0x079] = 0x90
	cpu.Memory[0x07A] = 0xF0
	cpu.Memory[0x07B] = 0x90
	cpu.Memory[0x07C] = 0xF0
	// 9
	cpu.Memory[0x07D] = 0xF0
	cpu.Memory[0x07E] = 0x90
	cpu.Memory[0x07F] = 0xF0
	cpu.Memory[0x080] = 0x10
	cpu.Memory[0x081] = 0xF0
	// A
	cpu.Memory[0x082] = 0xF0
	cpu.Memory[0x083] = 0x90
	cpu.Memory[0x084] = 0xF0
	cpu.Memory[0x085] = 0x90
	cpu.Memory[0x086] = 0x90
	// B
	cpu.Memory[0x087] = 0xE0
	cpu.Memory[0x088] = 0x90
	cpu.Memory[0x089] = 0xE0
	cpu.Memory[0x08A] = 0x90
	cpu.Memory[0x08B] = 0xE0
	// C
	cpu.Memory[0x08C] = 0xF0
	cpu.Memory[0x08D] = 0x80
	cpu.Memory[0x08E] = 0x80
	cpu.Memory[0x08F] = 0x80
	cpu.Memory[0x090] = 0xF0
	// D
	cpu.Memory[0x091] = 0xE0
	cpu.Memory[0x092] = 0x90
	cpu.Memory[0x093] = 0x90
	cpu.Memory[0x094] = 0x90
	cpu.Memory[0x095] = 0xE0
	// E
	cpu.Memory[0x096] = 0xF0
	cpu.Memory[0x097] = 0x80
	cpu.Memory[0x098] = 0xF0
	cpu.Memory[0x099] = 0x80
	cpu.Memory[0x09A] = 0xF0
	// F
	cpu.Memory[0x09B] = 0xF0
	cpu.Memory[0x09C] = 0x80
	cpu.Memory[0x09D] = 0xF0
	cpu.Memory[0x09E] = 0x80
	cpu.Memory[0x09F] = 0x80
}

// Initialisation du cpu
func InitCpu(cpu *Cpu, rombytes []byte) {
	cpu.initialiseFont()
	cpu.loadROM(rombytes)
	cpu.Pc = 0x200 - 2
}

// Update du cpu
func (cpu *Cpu) Update() {
	if time.Now().Sub(cpu.timeStart) > time.Second/16 { // when one second has past
		if cpu.Delay_timer > 0 {
			cpu.Delay_timer -= 1
		}
		if cpu.Sound_timer > 0 {
			if cpu.Sound_timer == 1 {
				Song("Beep")
			} else if cpu.Sound_timer == 30 {
				Song("Over")
			}
			cpu.Sound_timer -= 1
		}
		cpu.timeStart = time.Now()
	}
	cpu.GetKey() // recupere la touche pressée (si il y en a une)
	cpu.Pc += 2
	op1 := cpu.Memory[cpu.Pc]
	op2 := cpu.Memory[cpu.Pc+1]
	cpu.Opcode = cpu.uint8ToUint16(op1, op2)
	cpu.Decode(cpu.Opcode)
}

// chargement du rom
func (cpu *Cpu) loadROM(rombytes []byte) {
	cpu.Romlength = uint16(len(rombytes))
	for i, byt := range rombytes {
		cpu.Memory[0x200+i] = byt
	}
}

// Fonction uint16 to uint8
func (c *Cpu) Uint16ToUint8(n uint16) (uint8, uint8) {
	return uint8(n >> 8), uint8(n & 0x00FF)
}

// unit 8 to uint16
func (c *Cpu) uint8ToUint16(n1 uint8, n2 uint8) uint16 {
	return uint16(uint16(n1)<<8 | uint16(n2))
}

// Fonction uint8 to uint4
func (c *Cpu) Uint8ToUint4(n uint8) (uint8, uint8) {
	return uint8(n >> 4), uint8(n & 0x0F)
}

// DrawSprite dessine un sprite à l'écran et renvoie true si un pixel a été effacé
func (c *Cpu) DrawSprite(x byte, y byte, row byte) bool {
	erased := false
	yIndex := y % 32
	for i := x; i < x+8; i++ {
		xIndex := i % 64
		wasSet := c.Gfx[xIndex][yIndex] == 1
		value := row >> (x + 8 - i - 1) & 0x01

		c.Gfx[xIndex][yIndex] ^= value

		if wasSet && c.Gfx[xIndex][yIndex] == 0 {
			erased = true
		}
	}
	return erased
}
