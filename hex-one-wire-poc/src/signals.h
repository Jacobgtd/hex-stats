#ifndef SIGNALS_H
#define SIGNALS_H

#include <stdint.h>

typedef enum {
    PINMODE_INPUT,
    PINMODE_OUTPUT
} PinMode;

typedef enum {
    STATE_H,
    STATE_L
} State;

typedef struct {
    uint8_t out;
    uint8_t mode;
} Pin;

void pin_set_mode(Pin* pin, PinMode mode);

void pin_set_output(Pin* pin, uint8_t value);

typedef struct {
    Pin* pin0;
    Pin* pin1;
} Wire;

Wire wire_create(Pin* pin0, Pin* pin1);

State wire_read(Wire wire);

#endif // SIGNALS_H