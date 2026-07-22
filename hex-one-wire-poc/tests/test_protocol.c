#include <stdio.h>
#include "microcontroller.h"

void run_test(char* name, char (*func)(void)) {
    if (func()) {
        printf("Test %s passed\n", name);
    } else {
        printf("Test %s failed\n", name);
        exit(1);
    }
}

char test_1(void){

    Microcontroller* mc1 = create_microcontroller(1);
    Microcontroller* mc2 = create_microcontroller(2);

    microcontroller_wire(mc1, 0, mc2, 0);
    microcontroller_stop(mc1);
    microcontroller_stop(mc2);
    return 1;

}

int main() {
    run_test("test_1", test_1);
    return 0;

}