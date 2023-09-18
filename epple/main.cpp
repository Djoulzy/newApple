#include "diskcontroller.h"

uint8_t getNextBit(Disk2Drive drive) {
    uint8_t bit = drive.readPulse();
    drive.rotateDiskOneBit();
    return bit;
}

uint8_t getByte(Disk2Drive drive) {
    uint8_t bit, result;

	result = 0;
	for (bit = 0; bit == 0; bit = getNextBit(drive)) {
	}
	result = 0x80; // the bit we just retrieved is the high bit
	for (int i = 6; i >= 0; i--) {
		result |= getNextBit(drive) << i;
	}

	// fmt.Printf("Trk: %05.02f = %02X\n", W.physicalTrack, result)
	return result;
}

int main(int argc, char *argv[]) {
    bool debugfirst = true;
    uint8_t data, debugoutpos = 0;

    DiskController test(6, false, 12);

    test.loadMedia(0, argv[1]);
    test.reset();
    test.io(0xC0E0 + 0x0009, '\0', false);
    test.io(0xC0E0 + 0x000E, '\0', false);

    printf("Motor: %d\n", test.isMotorOn());
    printf("Track: 0x%02X\n", test.getTrack());

    for (int i=0; i<100; i++) {
        data = test.io(0xC0E0 + 0x000C, '\0', false);
        if ((data & 0x80u) != 0u) {
            if (debugfirst) {
                debugfirst = false;
                for (int i = 0; i < 128; ++i) {
                    printf("%02X", i);
                }
                printf("\n");
            }
            printf("%02X", data);
            ++debugoutpos;
            if (128 <= debugoutpos) {
                debugoutpos = 0;
                printf("\n");
                test.dumpLss();
            }
        }
    }

    return 0;
}