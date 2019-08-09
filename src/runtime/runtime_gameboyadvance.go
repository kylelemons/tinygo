// +build gameboyadvance

package runtime

import (
	"machine"
	"runtime/volatile"
)

// An Interrupt represents an interrupt for the GameBoyAdvance.
type Interrupt struct {
	n uint
}

var interrupts = [machine.InterruptCount]struct {
	handler func()
	periph  *volatile.Register16
	bit     uint16
}{
	machine.INT_VBLANK:   {periph: &machine.DisplayControl.DISPSTAT, bit: machine.DISPSTAT_VBLANK_IRQ},
	machine.INT_HBLANK:   {periph: &machine.DisplayControl.DISPSTAT, bit: machine.DISPSTAT_HBLANK_IRQ},
	machine.INT_VCOUNTER: {periph: &machine.DisplayControl.DISPSTAT, bit: machine.DISPSTAT_VCOUNTER_IRQ},
	machine.INT_TIMER0:   {periph: &machine.Timer[0].Control, bit: machine.TIMER_IRQ_ENABLE},
	machine.INT_TIMER1:   {periph: &machine.Timer[1].Control, bit: machine.TIMER_IRQ_ENABLE},
	machine.INT_TIMER2:   {periph: &machine.Timer[2].Control, bit: machine.TIMER_IRQ_ENABLE},
	machine.INT_TIMER3:   {periph: &machine.Timer[3].Control, bit: machine.TIMER_IRQ_ENABLE},
	machine.INT_SERIAL:   {}, // TODO
	machine.INT_DMA0:     {}, // TODO
	machine.INT_DMA1:     {}, // TODO
	machine.INT_DMA2:     {}, // TODO
	machine.INT_DMA3:     {}, // TODO
	machine.INT_KEYPAD:   {periph: &machine.Keypad.Control, bit: machine.KEY_IRQ_ENABLE},
	machine.INT_GAMEPAK:  {}, // TODO
}

// NewInterrupt returns an Interrupt for controlling GBA interrupts.
//
// This must be called in the global scope and assigned to a global variable.
//
// The behavior is undefined if NewInterrupt is called twice for the same
// interrupt in the same executable.
func NewInterrupt(interrupt uintptr, handler func()) Interrupt {
	// For now, our "undefined behavior" is simply "first call wins."
	if interrupts[interrupt].handler == nil {
		interrupts[interrupt].handler = handler
	}
	return Interrupt{
		n: uint(interrupt),
	}
}

// Enable enables the interrupt, allowing it to fire.
func (i Interrupt) Enable() {
	// defer ReenableInterrupts(DisableInterrupts())

	machine.Interrupts.Request.SetBits(1 << i.n)
	if reg := interrupts[i.n].periph; reg != nil {
		reg.SetBits(interrupts[i.n].bit)
	}
}

// Disable disables the interrupt, eventually preventing it from firing.
//
// Note that it is possible for an interrupt to happen after calling Disable if
// the peripheral fires the interrupt before the disable is registered.
func (i Interrupt) Disable() {
	// defer ReenableInterrupts(DisableInterrupts())

	machine.Interrupts.Request.ClearBits(1 << i.n)
	if reg := interrupts[i.n].periph; reg != nil {
		reg.ClearBits(interrupts[i.n].bit)
	}
}

// DisableInterrupts disables interrupts and returns an opaque state value for
// use in restoring the current interrupt state.
//
// In the case of the GameBoyAdvance, only the Global Interrupt Enable is saved
// and restored by these calls.
//
// Typically, this will be immediately followed by a deferred call to
// ReenableInterrupts.
//
// The returned state value is opaque and should not be interpreted.
func DisableInterrupts() (state uintptr) {
	return uintptr(machine.Interrupts.Enable.Get())
}

// ReenableInterrupts restores the interrupts given by mask.
//
// Interrupts that were enabled between the call to Disable and the call to
// Reenable will be maintained.
func ReenableInterrupts(state uintptr) {
	machine.Interrupts.Enable.Set(uint16(state))
}

//go:export runtime_isr_trampoline
func isr() {
	caught := machine.Interrupts.Ack.Get()

	// Ack the hardware interrupts
	machine.Interrupts.Ack.SetBits(caught)

	// Ack the interrupts to the BIOS too
	machine.InterruptBIOS.Ack.SetBits(caught)

	requested := machine.Interrupts.Request.Get()
	for i, intr := range interrupts {
		mask := uint16(1) << uint(i)
		if caught&requested&mask == mask {
			intr.handler()
		}
	}
}
