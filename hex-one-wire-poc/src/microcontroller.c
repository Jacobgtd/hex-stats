#include <stdint.h>
#include <stdio.h>
#include <windows.h>
#include "signals.h"
#include "microcontroller.h"

Microcontroller *create_microcontroller(uint32_t id)
{
    Microcontroller *mc = malloc(sizeof(Microcontroller));
    if (!mc)
    {
        return NULL;
    }
    mc->id = id;
    mc->r = 0;
    mc->g = 0;
    mc->b = 0;

    mc->current_action = IDLE;

    for (int i = 0; i < 6; i++)
    {
        mc->wires[i] = wire_create(&mc->pins[i], NULL);
    }

    mc->running = 1;
    int result = pthread_create(&mc->operator_thread, NULL, microcontroller_operate_thread, mc);
    if (result != 0)
    {
        mc->running = 0;
        mc->error = MICROCONTROLLER_COULD_NOT_START_OPERATOR_THREAD;
    }

    return mc;
}

void microcontroller_stop(Microcontroller *mc)
{
    mc->running = 0;
    pthread_join(mc->operator_thread, NULL);
    free(mc);
}

void *microcontroller_operate_thread(void *arg)
{
    Microcontroller *mc = (Microcontroller *)arg;
    microcontroller_operate(mc);
    return NULL;
}

void microcontroller_operate(Microcontroller *mc)
{

    while (mc->running)
    {
        Sleep(1);

        if (mc->doTiming && mc->timer_1 > 20 * COMM_PERIOD && microcontroller_in_listening_state(mc))
        {
            microcontroller_start_timer1(mc);
            mc->enable_interrupts = 1;
            mc->current_action = IDLE;
        }

        for (int i = 0; i < 6; i++)
        {

            State pin_now = microcontroller_read_pin(mc, i);
            if (pin_now == STATE_L && mc->last_edges[i] == STATE_H && mc->enable_interrupts)
            {

                if (microcontroller_in_listening_state(mc))
                {
                    microcontroller_start_timer1(mc);
                }

                switch (mc->current_action)
                {
                case IDLE:
                    mc->enable_interrupts = 0;
                    mc->current_action = LISTEN_PREAMBLE;
                // Writer states
                case WRITE_PREAMBLE:
                    break;
                case LISTEN_HS:
                    break;
                case WRITE_MSG:
                    break;
                case LISTEN_ACK:
                    break;
                // Reader states
                case LISTEN_PREAMBLE:
                    break;
                case WRITE_HS:
                    break;
                case LISTEN_MSG:
                    break;
                case WRITE_ACK:
                    break;
                }
            }

            mc->last_edges[i] = pin_now;
        }

        if (mc->sleep_counter)
        {
            mc->sleep_counter--;
            continue;
        }

        switch (mc->current_action)
        {
        // Writer states
        case IDLE:
            break;
        case WRITE_PREAMBLE:
            microcontroller_send_shout(mc, mc->pin_action, LISTEN_HS);
            break;
        case LISTEN_HS:
            microcontroller_read_shout(mc, mc->pin_action, WRITE_HS);
            break;
        case WRITE_MSG:
            // TODO
            break;
        case LISTEN_ACK:
            microcontroller_read_shout(mc, mc->pin_action, IDLE);
            break;
        // Reader states
        case LISTEN_PREAMBLE:
            microcontroller_read_shout(mc, mc->pin_action, WRITE_HS);
            break;
        case WRITE_HS:
            microcontroller_send_shout(mc, mc->pin_action, LISTEN_MSG);
            break;
        case LISTEN_MSG:
            // TODO
            break;
        case WRITE_ACK:
            microcontroller_send_shout(mc, mc->pin_action, IDLE);
            break;
        }
    }
}

char microcontroller_in_listening_state(Microcontroller *mc)
{
    return mc->current_action == LISTEN_PREAMBLE ||
           mc->current_action == LISTEN_HS ||
           mc->current_action == LISTEN_MSG ||
           mc->current_action == LISTEN_ACK;
}

void microcontroller_read_shout(Microcontroller *mc, uint8_t pin, MicrocontrollerAction next_action)
{
    microcontroller_sleep(mc, (SHOUT_LENGTH / 2) * COMM_PERIOD);
    State state = microcontroller_read_pin(mc, mc->pin_action);
    if (state == STATE_L)
    {
        mc->current_action = next_action;
    }
    else
    {
        mc->current_action = IDLE;
        mc->enable_interrupts = 1;
    }
}

void microcontroller_send_shout(Microcontroller *mc, uint8_t pin, MicrocontrollerAction next_action)
{
    microcontroller_set_pin_mode(mc, pin, PINMODE_OUTPUT);
    microcontroller_set_pin_output(mc, pin, 0);
    microcontroller_sleep(mc, SHOUT_LENGTH * COMM_PERIOD);
    microcontroller_set_pin_mode(mc, pin, PINMODE_INPUT);
    mc->current_action = next_action;
}

void microcontroller_sleep(Microcontroller *mc, int duration)
{
    mc->sleep_counter = duration;
}

void microcontroller_queue_send_led_color(Microcontroller *mc, uint8_t pin, uint8_t destination, uint8_t r, uint8_t g, uint8_t b)
{
    mc->current_action = WRITE_LED;
    mc->pin_action = pin;
    mc->destination = destination;
    mc->send_r = r;
    mc->send_g = g;
    mc->send_b = b;
}

void microcontroller_start_timer1(Microcontroller *mc)
{
    mc->timer_1 = 0;
    mc->doTiming = 1;
}

void microcontroller_set_led_color(Microcontroller *mc, uint8_t r, uint8_t g, uint8_t b)
{
    mc->r = r;
    mc->g = g;
    mc->b = b;
}

void microcontroller_wire(
    Microcontroller *mc0,
    uint8_t mc0_pin,
    Microcontroller *mc1,
    uint8_t mc1_pin)
{
    Wire wire = wire_create(&mc0->pins[mc0_pin], &mc1->pins[mc1_pin]);
    mc0->wires[mc0_pin] = wire;
    mc1->wires[mc1_pin] = wire;
}

void microcontroller_set_pin_mode(Microcontroller *mc, uint8_t pin, PinMode mode)
{
    if (pin < 6)
    {
        pin_set_mode(&mc->pins[pin], mode);
    }
}

void microcontroller_set_pin_output(Microcontroller *mc, uint8_t pin, uint8_t value)
{
    microcontroller_set_pin_mode(mc, pin, PINMODE_OUTPUT);
    if (pin < 6)
    {
        pin_set_output(&mc->pins[pin], value);
    }
}

State microcontroller_read_pin(Microcontroller *mc, uint8_t pin)
{

    if (mc->pins[pin].mode == PINMODE_INPUT)
    {
        return wire_read(mc->wires[pin]);
    }

    return mc->pins[pin].out ? STATE_H : STATE_L;
}