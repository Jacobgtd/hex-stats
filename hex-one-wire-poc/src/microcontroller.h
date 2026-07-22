#ifndef MICROCONTROLLER_H
#define MICROCONTROLLER_H

#define COMM_PERIOD 30
#define MAX_RETRIES 3
#define SHOUT_LENGTH 8

#include <stdint.h>
#include "signals.h"
#include <pthread.h>

typedef enum {
    MICROCONTROLLER_NO_ERROR,
    MICROCONTROLLER_COULD_NOT_START_OPERATOR_THREAD
} MicrocontrollerError;

typedef enum {
    IDLE,
    WRITE_PREAMBLE,
    LISTEN_PREAMBLE,
    WRITE_HS,
    LISTEN_HS,
    LISTEN_MSG,
    WRITE_MSG,
    WRITE_ACK,
    LISTEN_ACK,
} MicrocontrollerAction;

typedef struct {
    uint32_t id;

    uint8_t r;
    uint8_t g;
    uint8_t b;

    Pin pins[6];
    Wire wires[6];
    char last_edges[6];
    char enable_interrupts;

    char running;
    pthread_t operator_thread;
    MicrocontrollerError error;

    MicrocontrollerAction current_action;
    uint8_t pin_action;
    uint8_t destination;

    uint8_t send_r;
    uint8_t send_g;
    uint8_t send_b;

    uint8_t timer_1;
    uint8_t doTiming;

    uint8_t retries;
    int sleep_counter;



} Microcontroller;

Microcontroller* create_microcontroller(uint32_t id);

void* microcontroller_operate_thread(void* arg);

void microcontroller_operate(Microcontroller* mc);

void microcontroller_stop(Microcontroller* mc);

void microcontroller_set_led_color(
    Microcontroller* mc,
    uint8_t r,
    uint8_t g,
    uint8_t b
);

void microcontroller_wire(
    Microcontroller* mc0,
    uint8_t mc0_pin,
    Microcontroller* mc1,
    uint8_t mc1_pin
);

void microcontroller_set_pin_mode(
    Microcontroller* mc,
    uint8_t pin,
    PinMode mode
);

void microcontroller_set_pin_output(
    Microcontroller* mc,
    uint8_t pin,
    uint8_t value
);

State microcontroller_read_pin(Microcontroller* mc, uint8_t pin);

#endif // MICROCONTROLLER_H