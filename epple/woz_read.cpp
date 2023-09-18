#include "wozfile.h"

#include <istream>
#include <ostream>
#include <fstream>
#include <filesystem>
#include <cmath>
#include <cstring>


static std::uint8_t bc(std::uint8_t bit) {
    switch (bit) {
    case 0x80u: return 0u;
    case 0x40u: return 1u;
    case 0x20u: return 2u;
    case 0x10u: return 3u;
    case 0x08u: return 4u;
    case 0x04u: return 5u;
    case 0x02u: return 6u;
    case 0x01u: return 7u;
    }
    return 255u; // should never happen
}

static std::uint8_t cb(std::uint8_t bit) {
    switch (bit) {
    case 0u: return 0x80u;
    case 1u: return 0x40u;
    case 2u: return 0x20u;
    case 3u: return 0x10u;
    case 4u: return 0x08u;
    case 5u: return 0x04u;
    case 6u: return 0x02u;
    case 7u: return 0x01u;
    }
    return 255u; // should never happen
}

/*
 * Rotate the floppy disk by one bit.
 * In real life we don't care what track we're on, but for the
 * emulator we need to know. This is because the tracks within the
 * WOZ file could be different lengths. So in order to know when
 * we need to loop back to the beginning of the track (circular
 * track on the floppy), we need to know the actual bit length
 * of that track in our WOZ file.
 */
void WozFile::rotateOneBit(std::uint8_t currentQuarterTrack) {
    if (!isLoaded()) {
        return; // there's no disk to rotate
    }

    if (C_QTRACK <= currentQuarterTrack) {
        printf("attempt to move to illegal track.\n");
        return;
    }

    // Move to next bit
    this->bit >>= 1;

    // If we hit end of this byte, move on to beginning of next byte
    if (this->bit == 0) {
        ++this->byt;
        this->bit = 0x80u;
    }

    // this is an empty track, so any of the following byte/bit
    // adjustments don't apply now (they will be handled the
    // next time we hit a non-empty track)
    if (this->tmap[currentQuarterTrack] == 0xFFu) {
        return;
    }

    // this is the case where we are here for the first time,
    // and lastQuarterTrack has not been set to any valid prior track
    if (this->lastQuarterTrack == C_QTRACK) {
        this->lastQuarterTrack = currentQuarterTrack;
    }

    // If we changed tracks since the last time we were called,
    // we may need to adjust the rotational position. The new
    // position will be at the same relative position as the
    // previous, based on each track's length (tracks can be of
    // different lengths in the WOZ image).
    if (currentQuarterTrack != this->lastQuarterTrack) {
//        printf("switching from tmap[%02x] --> [%02x]\n", this->lastQuarterTrack, currentQuarterTrack);
        const double oldLen = this->trk_bits[this->tmap[this->lastQuarterTrack]];
        const double newLen = this->trk_bits[this->tmap[currentQuarterTrack]];
        const double ratio = newLen/oldLen;
        if (!(fabs(1-ratio) < 0.0001)) {
            const std::uint16_t newBit = static_cast<std::uint16_t>(round((this->byt*8+bc(this->bit)) * ratio));
            const std::uint8_t orig_bit = this->bit;
            const std::uint16_t orig_byt = this->byt;
            this->byt = newBit / 8;
            this->bit = cb(newBit % 8);
            printf("woz detected non 1:1 track size ratio: %f; adjusting byte/bit: %04X/%02X --> %04X/%02X\n",
                ratio, orig_byt, orig_bit, this->byt, this->bit);
        }
        this->lastQuarterTrack = currentQuarterTrack;
    }

    // Check for hitting the end of our track,
    // and if so, move back to the beginning.
    // This is how we emulate a circular track on the floppy.
    if (this->trk_bits[this->tmap[currentQuarterTrack]] <= this->byt*8u+bc(this->bit)) {
        this->byt = 0;
        this->bit = 0x80u;
    }
}

bool WozFile::exists(std::uint8_t currentQuarterTrack) {
    return isLoaded() && (this->tmap[currentQuarterTrack] != 0xFFu);
}

bool WozFile::getBit(std::uint8_t currentQuarterTrack) {
    return this->trk[this->tmap[currentQuarterTrack]][this->byt] & this->bit;
}

void WozFile::checkForWriteProtection() {
    if (!this->writable) {
        return;
    }

    if (!std::filesystem::exists(this->filePath)) {
        this->writable = false;
    }

    std::filesystem::path canon = std::filesystem::canonical(this->filePath);

    std::ofstream outf(canon, std::ios::binary|std::ios::app);
    this->writable = outf.is_open();
    outf.close();
}