// +build gameboyadvance

package machine

import (
	"image/color"
	"runtime"
	"runtime/volatile"
	"unsafe"
)

var (
	// IO maps the I/O Peripherals.
	IO = (*IORegs)(unsafe.Pointer(uintptr(0x04000000)))

	// Display maps the display memory in various modes.
	Display Displays

	// Tile maps the VRAM as tiles.
	Tile = TileRAM{
		Blocks4: (*[6][32 * 16]Tile4)(unsafe.Pointer(uintptr(0x06000000))),
		Blocks8: (*[6][16 * 16]Tile8)(unsafe.Pointer(uintptr(0x06000000))),
	}

	// BackgroundPalette maps the Background PaletteRAM as palettes.
	BackgroundPalette = PaletteRAM{
		Full: (*Palette)(unsafe.Pointer(uintptr(0x05000000))),
		Bank: (*[16]PaletteBank)(unsafe.Pointer(uintptr(0x05000000))),
	}

	// SpritePalette maps the Sprite PaletteRAM.
	SpritePalette = PaletteRAM{
		Full: (*Palette)(unsafe.Pointer(uintptr(0x05000200))),
		Bank: (*[16]PaletteBank)(unsafe.Pointer(uintptr(0x05000200))),
	}

	// Sprites maps the Object Attribute Memory for managing sprites.
	Sprites = SpriteRAM{
		Sprite: (*[128]SpriteAttrs)(unsafe.Pointer(uintptr(0x07000000))),
	}
)

// IOMap is the memory mapping of the IO Registers.
//
// Full IO Map:
//   https://www.akkit.org/info/gbatek.htm#gbaiomap
type IORegs struct {
	LCD    LCDRegs
	Sound  SoundRegs
	_      [0x100 - 0x0B0]byte // DMA
	Timer  [4]Timer            // Timers
	_      [0x10]byte          // Unused
	_      [0x130 - 0x120]byte // Serial 1
	Keypad Keypad
	_      [0x200 - 0x134]byte // Serial 2
	Int    InterruptRegs
}

type LCDRegs struct {
	DISPCNT  volatile.Register16 // R/W - LCD Control
	_        volatile.Register16 // R/W - Undocumented - Green Swap
	DISPSTAT volatile.Register16 // R/W - General LCD Status (STAT,LYC)
	VCOUNT   volatile.Register16 // R   - Vertical Counter (LY)
	BG0CNT   volatile.Register16 // R/W - BG0 Control
	BG1CNT   volatile.Register16 // R/W - BG1 Control
	BG2CNT   volatile.Register16 // R/W - BG2 Control
	BG3CNT   volatile.Register16 // R/W - BG3 Control
	BG0HOFS  volatile.Register16 // W   - BG0 X-Offset
	BG0VOFS  volatile.Register16 // W   - BG0 Y-Offset
	BG1HOFS  volatile.Register16 // W   - BG1 X-Offset
	BG1VOFS  volatile.Register16 // W   - BG1 Y-Offset
	BG2HOFS  volatile.Register16 // W   - BG2 X-Offset
	BG2VOFS  volatile.Register16 // W   - BG2 Y-Offset
	BG3HOFS  volatile.Register16 // W   - BG3 X-Offset
	BG3VOFS  volatile.Register16 // W   - BG3 Y-Offset
	BG2PA    volatile.Register16 // W   - BG2 Rotation/Scaling Parameter A (dx)
	BG2PB    volatile.Register16 // W   - BG2 Rotation/Scaling Parameter B (dmx)
	BG2PC    volatile.Register16 // W   - BG2 Rotation/Scaling Parameter C (dy)
	BG2PD    volatile.Register16 // W   - BG2 Rotation/Scaling Parameter D (dmy)
	BG2X     volatile.Register32 // W   - BG2 Reference Point X-Coordinate
	BG2Y     volatile.Register32 // W   - BG2 Reference Point Y-Coordinate
	BG3PA    volatile.Register16 // W   - BG3 Rotation/Scaling Parameter A (dx)
	BG3PB    volatile.Register16 // W   - BG3 Rotation/Scaling Parameter B (dmx)
	BG3PC    volatile.Register16 // W   - BG3 Rotation/Scaling Parameter C (dy)
	BG3PD    volatile.Register16 // W   - BG3 Rotation/Scaling Parameter D (dmy)
	BG3X     volatile.Register32 // W   - BG3 Reference Point X-Coordinate
	BG3Y     volatile.Register32 // W   - BG3 Reference Point Y-Coordinate
	WIN0H    volatile.Register16 // W   - Window 0 Horizontal Dimensions
	WIN1H    volatile.Register16 // W   - Window 1 Horizontal Dimensions
	WIN0V    volatile.Register16 // W   - Window 0 Vertical Dimensions
	WIN1V    volatile.Register16 // W   - Window 1 Vertical Dimensions
	WININ    volatile.Register16 // R/W - Inside of Window 0 and 1
	WINOUT   volatile.Register16 // R/W - Inside of OBJ Window & Outside of Windows
	MOSAIC   volatile.Register16 // W   - Mosaic Size
	_        volatile.Register16 // -   - Not used
	BLDCNT   volatile.Register16 // R/W - Color Special Effects Selection
	BLDALPHA volatile.Register16 // W   - Alpha Blending Coefficients
	BLDY     volatile.Register16 // W   - Brightness (Fade-In/Out) Coefficient
	_        volatile.Register32 // -   - Not used
	_        volatile.Register32 // -   - Not used
}

type SoundRegs struct {
	SOUND1CNT_L volatile.Register16    // R/W - Channel 1 Sweep register       (NR10)
	SOUND1CNT_H volatile.Register16    // R/W - Channel 1 Duty/Length/Envelope (NR11, NR12)
	SOUND1CNT_X volatile.Register16    // R/W - Channel 1 Frequency/Control    (NR13, NR14)
	_           volatile.Register16    // -   - Not used
	SOUND2CNT_L volatile.Register16    // R/W - Channel 2 Duty/Length/Envelope (NR21, NR22)
	_           volatile.Register16    // -   - Not used
	SOUND2CNT_H volatile.Register16    // R/W - Channel 2 Frequency/Control    (NR23, NR24)
	_           volatile.Register16    // -   - Not used
	SOUND3CNT_L volatile.Register16    // R/W - Channel 3 Stop/Wave RAM select (NR30)
	SOUND3CNT_H volatile.Register16    // R/W - Channel 3 Length/Volume        (NR31, NR32)
	SOUND3CNT_X volatile.Register16    // R/W - Channel 3 Frequency/Control    (NR33, NR34)
	_           volatile.Register16    // -   - Not used
	SOUND4CNT_L volatile.Register16    // R/W - Channel 4 Length/Envelope      (NR41, NR42)
	_           volatile.Register16    // -   - Not used
	SOUND4CNT_H volatile.Register16    // R/W - Channel 4 Frequency/Control    (NR43, NR44)
	_           volatile.Register16    // -   - Not used
	SOUNDCNT_L  volatile.Register16    // R/W - Control Stereo/Volume/Enable   (NR50, NR51)
	SOUNDCNT_H  volatile.Register16    // R/W - Control Mixing/DMA Control
	SOUNDCNT_X  volatile.Register16    // R/W - Control Sound on/off           (NR52)
	_           volatile.Register16    // -   - Not used
	SOUNDBIAS   volatile.Register16    // BIAS- Sound PWM Control
	_           [3]volatile.Register16 // -   - Not used
	WAVE_RAM    [2][8]byte             // R/W - Channel 3 Wave Pattern RAM (2 banks!!)
	FIFO_A      volatile.Register32    // W   - Channel A FIFO, Data 0-3
	FIFO_B      volatile.Register32    // W   - Channel B FIFO, Data 0-3
	_           [4]volatile.Register16 // Not used
}

type InterruptRegs struct {
	Request volatile.Register16 // R/W - Interrupt Request
	Ack     volatile.Register16 // R/W - Interrupt Active (R) / Acknowledge (W)
	_       volatile.Register32 // -   - Not used
	Enable  volatile.Register16 // R/W - Interrupt Master Enable
}

type PinMode uint8

// Set has not been implemented.
func (p Pin) Set(value bool) {
	// do nothing
}

// Displays is a convenience container for the various mode-based displays.
type Displays struct {
	Mode3 DisplayMode3
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

	DISPCNT_MODE_MASK       = 0x7
	DISPCNT_DISPLAY_BG_MASK = 0x3 << 8
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

var (
	// Mode 3 uses a single 16bpp 32k color 240x160 framebuffer.
	Mode3Frame0 = (*[160][240]volatile.Register16)(unsafe.Pointer(uintptr(0x06000000)))
)

// DisplayMode3 maps the Mode 3 Bitmap framebuffer.
type DisplayMode3 struct{}

func (DisplayMode3) Configure() {
	// Write into the I/O registers, setting video display parameters.
	//
	// Use video mode 3 (in BG2, a 16bpp bitmap in VRAM)
	IO.LCD.DISPCNT.ClearBits(DISPCNT_MODE_MASK | DISPCNT_DISPLAY_BG_MASK | DISPCNT_FORCED_BLANK)
	IO.LCD.DISPCNT.SetBits(DISPCNT_MODE3 | DISPCNT_DISPLAY_BG2)
}

func (DisplayMode3) Size() (x, y int16) {
	return 240, 160
}

func (DisplayMode3) SetPixel(x, y int16, c color.RGBA) {
	Mode3Frame0[y][x].Set(uint16(c.R)&0x1f | uint16(c.G)&0x1f<<5 | uint16(c.B)&0x1f<<10)
}

func (DisplayMode3) Display() error {
	// Nothing to do here.
	return nil
}

// Tile4 is an 8x8 4bpp tile.
type Tile4 [8]volatile.Register32

// Tile8 is an 8x8 8bpp tile.
type Tile8 [16]volatile.Register32

// TileRAM provides access to the tilesets.
//
// Be VERY careful if you use both 4bpp and 8bpp tiles!  They overlap in memory,
// so if you are not careful you may alias your tiles.
type TileRAM struct {
	Blocks4 *[6][32 * 16]Tile4 // tile memory mapped to 4bpp tiles (alias of Blocks8!)
	Blocks8 *[6][16 * 16]Tile8 // tile memory mapped to 8bpp tiles (alias of Blocks4!)
}

// TileRAM Block index constants.
const (
	TILE_BLOCK_B0 = iota // tile index starts from CBB and counts bitdepth
	TILE_BLOCK_B1        // tile index starts from CBB and counts bitdepth
	TILE_BLOCK_B2        // tile index starts from CBB and counts bitdepth
	TILE_BLOCK_B3        // tile index starts from CBB and counts bitdepth
	TILE_BLOCK_S0        // tile index starts from S0 and counts Tile4s
	TILE_BLOCK_S1        // tile index starts from S0 and counts Tile4s
)

type Palette [256]volatile.Register16
type PaletteBank [16]volatile.Register16

// PaletteRAM provides access to the palettes.
//
// Note that sprites and background sprites have separate ones!
type PaletteRAM struct {
	Full *Palette         // Full[PAL_BG | PAL_SPRITE][index] = color
	Bank *[16]PaletteBank // Bank[PAL_BG | PAL_SPRITE][bank][index] = color
}

// PaletteRAM index constants.
const (
	PAL_BG     = 0 // index for the background palette
	PAL_SPRITE = 1 // index for the sprite palette
)

type SpriteRAM struct {
	Sprite *[128]SpriteAttrs
}

// Enable2D sets the display to assume 2D sprites.
func (SpriteRAM) Enable2D() {
	IO.LCD.DISPCNT.ClearBits(DISPCNT_OBJ_VRAM_MAP_1D)
	IO.LCD.DISPCNT.SetBits(DISPCNT_DISPLAY_OBJ)
}

// Enable1D sets the display to assume 1D sprites.
func (SpriteRAM) Enable1D() {
	IO.LCD.DISPCNT.SetBits(DISPCNT_DISPLAY_OBJ | DISPCNT_OBJ_VRAM_MAP_1D)
}

// Disable disables sprite rendering.
func (SpriteRAM) Disable() {
	IO.LCD.DISPCNT.ClearBits(DISPCNT_DISPLAY_OBJ)
}

// SpriteAttrs map the Object Attribute registers for a single sprite.
type SpriteAttrs struct {
	Y volatile.Register16 // Y offset, mode, flags, bitdepth, shape, etc
	X volatile.Register16 // X offset, flipping, size, etc
	T volatile.Register16 // Tile offset, priority, etc
	_ volatile.Register16 // affine transform interleaved
}

// Sprite*Attrs constants help the compiler ensure that these otherwise-unitless
// values are only passed for the correct configuration options.
type (
	SpriteYAttrs uint16
	SpriteXAttrs uint16
)

// Sprite Attribute Constants
const (
	// Bit Depth (selected by SetupN, no need to pass them in)
	SPRITE_4BPP = 0 << 13
	SPRITE_8BPP = 1 << 13

	// Offsets (add to tile index if using TILE_BLOCK_S1)
	SPRITE_4BPP_S1_OFFSET = 512
	SPRITE_8BPP_S1_OFFSET = 256

	// Shape (Y Atribute)
	SPRITE_SQUARE SpriteYAttrs = 0 << 14
	SPRITE_WIDE   SpriteYAttrs = 1 << 14
	SPRITE_TALL   SpriteYAttrs = 2 << 14

	// Size (X Attribute)
	SPRITE_SIZE_S  SpriteXAttrs = 0 << 14
	SPRITE_SIZE_M  SpriteXAttrs = 1 << 14
	SPRITE_SIZE_L  SpriteXAttrs = 2 << 14
	SPRITE_SIZE_XL SpriteXAttrs = 3 << 14

	// Flipping (X Attribute, non-affine only)
	SPRITE_FLIP_H SpriteXAttrs = 1 << 12
	SPRITE_FLIP_V SpriteXAttrs = 1 << 13
)

// Setup4 sets up a 4bpp sprite.
//
// Note that if your tiles are in the second sprite tile block, you will need to
// add SPRITE_4BPP_S1_OFFSET before passing it in.
func (s *SpriteAttrs) Setup4(priority, tile, bank int, y SpriteYAttrs, x SpriteXAttrs) {
	s.Y.Set(uint16(y))
	s.X.Set(uint16(x))
	s.T.Set(uint16(bank&0xF)<<12 | uint16(priority&0x3)<<10 | uint16(tile&0x3FF))
}

// Setup4 sets up a 8bpp sprite.
//
// The tile index will be automatically scaled for 8bpp tiles, so this index
// will match the one used for the tile block.
//
// Note that if your tiles are in the second sprite tile block, you will need to
// add SPRITE_8BPP_S1_OFFSET before passing it in.
func (s *SpriteAttrs) Setup8(priority, tile int, y SpriteYAttrs, x SpriteXAttrs) {
	tile *= 2 // tile offset counts Tile4s
	s.Y.Set(uint16(y | SPRITE_8BPP))
	s.X.Set(uint16(x))
	s.T.Set(uint16(priority&0x3)<<10 | uint16(tile&0x3FF))
}

// TODO(kevlar): helpers for setting setup fields after setup:
//   - SetFlip

// SetPos sets the position of the sprite.
//
// The x coordinate is in the range [0,512).
// The y coordinate is in the range [0,256).
func (s *SpriteAttrs) SetPos(x, y int) {
	s.Y.ClearBits(0xFF)
	s.Y.SetBits(uint16(y & 0xFF))
	s.X.ClearBits(0xFF)
	s.X.SetBits(uint16(x & 0x1FF))
}

type Keypad struct {
	Status  volatile.Register16 // R   - Key Status
	Control volatile.Register16 // R/W - Key Interrupt Control
}

// A Key represents one of the possible keys.
type Key uint16

// Keypad constants
const (
	KEY_A Key = 1 << iota
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
	KEY_ANY Key = 0x3FF

	KEY_IRQ_ENABLE = 1 << 14
	KEY_IRQ_ALL    = 1 << 15
)

func (k *Keypad) WakeOn(keys ...Key) {
	var mask uint16
	for _, k := range keys {
		mask |= uint16(k)
	}
	k.Control.ClearBits(uint16(KEY_ANY))
	k.Control.SetBits(mask)
}

type Timer struct {
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

func (t *Timer) Start() {
	t.Control.SetBits(TIMER_START)
}
func (t *Timer) Stop() {
	t.Control.ClearBits(TIMER_START)
}

type Interrupt int

// Interrupt constants
const (
	INT_VBLANK   Interrupt = iota // LCD V-Blank
	INT_HBLANK                    // LCD H-Blank
	INT_VCOUNTER                  // LCD V-Counter Match
	INT_TIMER0                    // Timer 0 Overflow
	INT_TIMER1                    // Timer 1 Overflow
	INT_TIMER2                    // Timer 2 Overflow
	INT_TIMER3                    // Timer 3 Overflow
	INT_SERIAL                    // Serial Communication
	INT_DMA0                      // DMA 0
	INT_DMA1                      // DMA 1
	INT_DMA2                      // DMA 2
	INT_DMA3                      // DMA 3
	INT_KEYPAD                    // Keypad
	INT_GAMEPAK                   // Game Pak (external IRQ source)
	InterruptCount
)

type InterruptController struct {
	handlers [InterruptCount]func(Interrupt)
}

// Interrupts is a static interrupt controller.
var Interrupts InterruptController

var interruptEnabler = [InterruptCount]struct {
	Register *volatile.Register16
	Bit      uint16
}{
	INT_VBLANK:   {&IO.LCD.DISPSTAT, DISPSTAT_VBLANK_IRQ},
	INT_HBLANK:   {&IO.LCD.DISPSTAT, DISPSTAT_HBLANK_IRQ},
	INT_VCOUNTER: {&IO.LCD.DISPSTAT, DISPSTAT_VCOUNTER_IRQ},
	INT_TIMER0:   {&IO.Timer[0].Control, TIMER_IRQ_ENABLE},
	INT_TIMER1:   {&IO.Timer[1].Control, TIMER_IRQ_ENABLE},
	INT_TIMER2:   {&IO.Timer[2].Control, TIMER_IRQ_ENABLE},
	INT_TIMER3:   {&IO.Timer[3].Control, TIMER_IRQ_ENABLE},
	INT_SERIAL:   {}, // TODO
	INT_DMA0:     {}, // TODO
	INT_DMA1:     {}, // TODO
	INT_DMA2:     {}, // TODO
	INT_DMA3:     {}, // TODO
	INT_KEYPAD:   {&IO.Keypad.Control, KEY_IRQ_ENABLE},
	INT_GAMEPAK:  {}, // TODO
}

func (ic *InterruptController) Enable(f func(Interrupt), ints ...Interrupt) {
	IO.Int.Enable.Set(0)

	for _, intr := range ints {
		enabler := interruptEnabler[int(intr)]
		if enabler.Register != nil {
			enabler.Register.SetBits(enabler.Bit)
		}

		IO.Int.Request.SetBits(1 << uint16(intr))
		ic.handlers[intr] = f
	}

	IO.Int.Enable.Set(1)
}

func (ic *InterruptController) Disable(ints ...Interrupt) {
	ime := IO.Int.Enable.Get()
	IO.Int.Enable.Set(0)
	defer IO.Int.Enable.Set(ime)

	for _, intr := range ints {
		enabler := interruptEnabler[int(intr)]
		if enabler.Register != nil {
			enabler.Register.ClearBits(enabler.Bit)
		}

		IO.Int.Request.ClearBits(1 << uint16(intr))
		ic.handlers[intr] = nil
	}
}

func (ic *InterruptController) DisableAll() {
	IO.Int.Enable.Set(0)
	ic.handlers = [InterruptCount]func(Interrupt){}
}

func (ic *InterruptController) handle(caught uint16) {
	for i := Interrupt(0); i < InterruptCount; i++ {
		if caught&(1<<uint16(i)) != 0 && ic.handlers[i] != nil {
			ic.handlers[i](i)
		}
	}
}

func init() {
	runtime.UserISR = isr
}

func isr() {
	caught := IO.Int.Ack.Get()
	IO.Int.Ack.SetBits(caught) // ack all interrupts
	Interrupts.handle(caught)
}
