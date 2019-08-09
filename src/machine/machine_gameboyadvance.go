// +build gameboyadvance

package machine

import (
	"image/color"
	"runtime/volatile"
	"unsafe"
)

// Memory-mapped I/O
//
// Full IO Map:
//   https://www.akkit.org/info/gbatek.htm#gbaiomap
var (
	// DisplayControl maps the display control memory.
	DisplayControl = (*DisplayRegs)(unsafe.Pointer(uintptr(0x04000000)))

	// Interrupts maps the interrupt control registers.
	Interrupts = (*InterruptRegs)(unsafe.Pointer(uintptr(0x04000200)))

	// InterruptBIOS maps the BIOS interrupt control registers.
	InterruptBIOS = (*InterruptBIOSRegs)(unsafe.Pointer(uintptr(0x03007FF8)))

	// Keypad maps the keypad control and status registers.
	Keypad = (*KeypadRegs)(unsafe.Pointer(uintptr(0x04000130)))

	// The Timer array maps the four timer control registers.
	Timer = (*[4]TimerRegs)(unsafe.Pointer(uintptr(0x04000100)))
)

// Displays:

// Display is a convenience alias for Mode3Framebuffer.
var Display Mode3Framebuffer

var mode3Framebuffer = (*[160][240]volatile.Register16)(unsafe.Pointer(uintptr(0x06000000)))

// Mode3Framebuffer is a convenience wrapper for setting pixels in the 240x160 16bpp bitmap mode.
type Mode3Framebuffer struct{}

func (Mode3Framebuffer) Configure() {
	DisplayControl.DISPCNT.Set(DISPCNT_MODE3 | DISPCNT_DISPLAY_BG2)
}

func (Mode3Framebuffer) Size() (w, h int16) {
	return 240, 160
}

func (Mode3Framebuffer) SetPixel(x, y int16, c color.RGBA) {
	mode3Framebuffer[y][x].Set(uint16(c.R)&0x1f | uint16(c.G)&0x1f<<5 | uint16(c.B)&0x1f<<10)
}

func (Mode3Framebuffer) Display() error {
	return nil
}

type DisplayRegs struct {
	DISPCNT  volatile.Register16 // R/W - LCD Control
	_        volatile.Register16 // R/W - Undocumented - Green Swap
	DISPSTAT volatile.Register16 // R/W - General LCD Status (STAT,LYC)
	VCOUNT   volatile.Register16 // R   - Vertical Counter (LY)
}

// DISPCNT Register Constants
const (
	DISPCNT_MODE1 = 1
	DISPCNT_MODE2 = 2
	DISPCNT_MODE3 = 3
	DISPCNT_MODE4 = 4
	DISPCNT_MODE5 = 5

	DISPCNT_FRAME_SELECT         = 1 << 4
	DISPCNT_HBLANK_INTERVAL_FREE = 1 << 5
	DISPCNT_OBJ_VRAM_MAP_1D      = 1 << 6
	DISPCNT_FORCED_BLANK         = 1 << 7
	DISPCNT_DISPLAY_BG0          = 1 << 8
	DISPCNT_DISPLAY_BG1          = 1 << 9
	DISPCNT_DISPLAY_BG2          = 1 << 10
	DISPCNT_DISPLAY_BG3          = 1 << 11
	DISPCNT_DISPLAY_OBJ          = 1 << 12

	DISPCNT_MODE_MASK = 0x7
	DISPCNT_BG_MASK   = 0xf << 8
)

// DISPSTAT Register Constants
const (
	DISPSTAT_VBLANK       = 1 << iota // V-Blank flag   (Read only) (1=VBlank) (set in line 160..226; not 227)
	DISPSTAT_HBLANK                   // H-Blank flag   (Read only) (1=HBlank) (toggled in all lines, 0..227)
	DISPSTAT_VCOUNTER                 // V-Counter flag (Read only) (1=Match)  (set in selected line)
	DISPSTAT_VBLANK_IRQ               // V-Blank IRQ Enable         (1=Enable)
	DISPSTAT_HBLANK_IRQ               // H-Blank IRQ Enable         (1=Enable)
	DISPSTAT_VCOUNTER_IRQ             // V-Counter IRQ Enable       (1=Enable)
)

// Interrupts:

type InterruptRegs struct {
	Request volatile.Register16 // R/W - Interrupt Request
	Ack     volatile.Register16 // R/W - Interrupt Active (R) / Acknowledge (W)
	_       volatile.Register32 // -   - Not used
	Enable  volatile.Register16 // R/W - Interrupt Master Enable
}

type InterruptBIOSRegs struct {
	Ack volatile.Register16 // R/W - BIOS Interrupt acknoweldgements
}

// Interrupt constants
const (
	INT_VBLANK   = iota // LCD V-Blank
	INT_HBLANK          // LCD H-Blank
	INT_VCOUNTER        // LCD V-Counter Match
	INT_TIMER0          // Timer 0 Overflow
	INT_TIMER1          // Timer 1 Overflow
	INT_TIMER2          // Timer 2 Overflow
	INT_TIMER3          // Timer 3 Overflow
	INT_SERIAL          // Serial Communication
	INT_DMA0            // DMA 0
	INT_DMA1            // DMA 1
	INT_DMA2            // DMA 2
	INT_DMA3            // DMA 3
	INT_KEYPAD          // Keypad
	INT_GAMEPAK         // Game Pak (external IRQ source)
	InterruptCount
)

// Keypad:

type KeypadRegs struct {
	Status  volatile.Register16 // R   - Key Status
	Control volatile.Register16 // R/W - Key Interrupt Control
}

// Keypad constants
const (
	KEY_A = 1 << iota
	KEY_B
	KEY_SELECT
	KEY_START
	KEY_RIGHT
	KEY_LEFT
	KEY_UP
	KEY_DOWN
	KEY_RB
	KEY_LB

	// KEY_ANY has the bits for every key set.
	KEY_ANY = 0x3FF

	KEY_IRQ_ENABLE = 1 << 14
	KEY_IRQ_ALL    = 1 << 15
)

// Timers:

type TimerRegs struct {
	Counter volatile.Register16 // R/W - Counter (R) / Reload (W)
	Control volatile.Register16 // R/W - Flags (see TIMER_*)
}

// Timer constants
const (
	TIMER_PRESCALE_1      = 0
	TIMER_PRESCALE_64     = 1
	TIMER_PRESCALE_256    = 2
	TIMER_PRESCALE_1024   = 3
	TIMER_COUNT_OVERFLOWS = 1 << 2
	TIMER_IRQ_ENABLE      = 1 << 6
	TIMER_START           = 1 << 7
)
