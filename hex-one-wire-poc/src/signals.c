#include <stdint.h>
#include "signals.h"

void pin_set_mode(Pin* pin, PinMode mode) {
    if (pin) {
        pin->mode = mode;
    }
}

void pin_set_output(Pin* pin, uint8_t value) {
    if (pin) {
        pin->out = value ? 1 : 0;
    }
}

Wire wire_create(Pin* pin0, Pin* pin1) {
    Wire wire;
    wire.pin0 = pin0;
    wire.pin1 = pin1;
    return wire;
}

State wire_read(Wire wire) {

    State state = STATE_H;
    if (wire.pin0 && wire.pin0->mode == PINMODE_OUTPUT && !wire.pin0->out) {
        state = STATE_L;
    }

    if (wire.pin1 && wire.pin1->mode == PINMODE_OUTPUT && !wire.pin1->out) {
        state = STATE_L;
    }

    return state;
}