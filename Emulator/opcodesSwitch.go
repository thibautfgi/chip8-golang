package emulator

// Fonction pour décoder les opcodes
func (c *Cpu) Decode(opcode uint16) {

	// Diviser l'opcode en parties individuelles pour faciliter le décodage
	opcodeX := byte(opcode>>8) & 0x000F // Bits 8 à 11
	opcodeY := byte(opcode>>4) & 0x000F // Bits 4 à 7
	opcodeNNN := opcode & 0x0FFF        // Bits 0 à 11
	opcodeNN := byte(opcode & 0x00FF)   // Bits 0 à 7
	opcodeN4 := byte(opcode & 0x000F)   // 4 derniers bits

	switch opcode & 0xF000 {
	case 0x0000:
		switch opcode {
		case 0x00E0:
			c.op00E0()
		case 0x00EE:
			c.op00EE()
		default:
		}
	case 0x1000:
		c.op1nnn(opcodeNNN)
	case 0x2000:
		c.op2nnn(opcodeNNN)
	case 0x3000:
		c.op3nnn(opcodeX, opcodeNN)
	case 0x4000:
		c.op4nnn(opcodeX, opcodeNN)
	case 0x5000:
		c.op5nnn(opcodeX, opcodeY)
	case 0x6000:
		c.op6nnn(opcodeX, opcodeNN)
	case 0x7000:
		c.op7nnn(opcodeX, opcodeNN)
	case 0x8000:
		switch opcode & 0x000F {
		case 0x0000:
			c.op8nn0(opcodeX, opcodeY)
		case 0x0001:
			c.op8xy1(opcodeX, opcodeY)
		case 0x0002:
			c.op8xy2(opcodeX, opcodeY)
		case 0x0003:
			c.op8xy3(opcodeX, opcodeY)
		case 0x0004:
			c.op8xy4(opcodeX, opcodeY)
		case 0x0005:
			c.op8xy5(opcodeX, opcodeY)
		case 0x0006:
			c.op8xy6(opcodeX, opcodeY)
		case 0x0007:
			c.op8xy7(opcodeX, opcodeY)
		case 0xE:
			c.op8xyE(opcodeX, opcodeY)
		default:
		}
	case 0x9000:
		c.op9nn0(opcodeX, opcodeY)
	case 0xA000:
		c.opAnnn(opcodeNNN)
	case 0xB000:
		c.opBnnn(opcodeNNN)
	case 0xC000:
		c.opCxkk(opcodeX, opcodeNN)
	case 0xD000:
		c.opDxyn(opcodeX, opcodeY, opcodeN4)
	case 0xE000:
		switch opcode & 0x000F {
		case 0x000E:
			c.opEx9E(opcodeX)
		case 0x0001:
			c.opExA1(opcodeX)
		default:
		}
	case 0xF000:
		switch opcodeNNN & 0xFF {
		case 0x07:
			c.opFx07(opcodeX)
		case 0x0A:
			c.opFx0A(opcodeX)
		case 0x15:
			c.opFx15(opcodeX)
		case 0x18:
			c.opFx18(opcodeX)
		case 0x1E:
			c.opFx1E(opcodeX)
		case 0x29:
			c.opFx29(opcodeX)
		case 0x33:
			c.opFx33(opcodeX)
		case 0x55:
			c.opFx55(opcodeX)
		case 0x65:
			c.opFx65(opcodeX)
		default:
		}
	default:
	}
}
