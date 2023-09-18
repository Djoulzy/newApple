#include "wozfile.h"

#include <istream>
#include <ostream>
#include <fstream>
#include <filesystem>
#include <cmath>
#include <cstring>

bool WozFile::trackIsZeroes(std::uint8_t qt) {
    for (std::uint16_t byt = 0; byt < this->trk_byts[qt]; ++byt) {
        if (this->trk[qt][byt]) {
            return false;
        }
    }
    return true;
}

bool WozFile::tracksAreIdentical(std::uint8_t qt1, std::uint8_t qt2) {
    if (this->trk_bits[qt1] != this->trk_bits[qt2]) {
        return false;
    }
    for (std::uint16_t byt = 0; byt < this->trk_byts[qt1]; ++byt) {
        if (this->trk[qt1][byt] != this->trk[qt2][byt]) {
            return false;
        }
    }
    return true;
}

// opposite of expandTracks()
void WozFile::reduceTracks() {
    for (std::uint8_t qt(0); qt < C_QTRACK; ++qt) {
        if (trackIsZeroes(this->tmap[qt])) {
            removeTrack(this->tmap[qt]);
            this->tmap[qt] = 0xFFu;
        }
    }
    for (std::uint8_t qt(0); qt < C_QTRACK-1; ++qt) {
        for (std::uint8_t qto(qt+1); qto < C_QTRACK; ++qto) {
            if (this->tmap[qt] != 0xFFu && this->tmap[qto] != 0xFFu && this->tmap[qt] != this->tmap[qto]) {
                if (tracksAreIdentical(this->tmap[qt], this->tmap[qto])) {
                    removeTrack(this->tmap[qto]);
                    this->tmap[qto] = this->tmap[qt];
                }
            }
        }
    }
    // kludge? remove track $22.25 if standard disk
    for (std::uint8_t qt(C_QTRACK-1); qt >= 0x89; --qt) {
        if (this->tmap[qt] != 0xFFu) {
            if (qt == 0x89) {
                if (this->tmap[qt] == this->tmap[qt-1]) {
                    this->tmap[qt] = 0xFFu;
                }
            }
            break;
        }
    }
//    dumpTmap();
//    dumpTracks();
}

static std::uint16_t bytesForBits(const std::uint32_t c_bits) {
    std::uint16_t c_bytes = (c_bits + 7) / 8;
    return static_cast<std::uint16_t>(c_bytes + 0x1FFu) / 0x200 * 0x200;
}

void WozFile::expandTracks() {
    for (std::uint8_t qt(0); qt < C_QTRACK; ++qt) {
        if (this->tmap[qt] != 0xFFu) {
            for (std::uint8_t qto(qt-1); qto != 0xFFu; --qto) {
                if (this->tmap[qt] == this->tmap[qto]) {
                    copyTrack(qt, qto);
                }
            }
        } else {
            createNewTrack(qt);
        }
    }
}



void WozFile::rawSet(std::uint8_t currentQuarterTrack, bool on) {
    if (on) {
        this->trk[this->tmap[currentQuarterTrack]][this->byt] |= this->bit;
    } else {
        this->trk[this->tmap[currentQuarterTrack]][this->byt] &= ~this->bit;
    }
}


void WozFile::setBit(std::uint8_t currentQuarterTrack, bool on) {
    if (!isLoaded()) {
        return; // there's no disk to write data to
    }

    if (!this->writable) {
        return; // write-protected
    }

//    printf("%c",(on?'1':'0')); fflush(stdout);
    rawSet(currentQuarterTrack, on);



    // also write preceding and following quarter tracks (at relative position, and if they exist)
    if (0 < currentQuarterTrack) {
        rawSet(currentQuarterTrack-1, on);
    }
    if (currentQuarterTrack < C_QTRACK-1) {
        rawSet(currentQuarterTrack+1, on);
    }

    this->modified = true;
}

void WozFile::removeTrack(const int trackIndex) {
    if (this->trk[trackIndex]) {
        delete [] this->trk[trackIndex];
        this->trk[trackIndex] = 0;
    }
    this->trk_bits[trackIndex] = 0;
    this->trk_byts[trackIndex] = 0;
}

void WozFile::copyTrack(std::uint8_t qt_dest, std::uint8_t qt_src) {
    for (std::uint8_t t(0); t < C_QTRACK; ++t) {
        if (!this->trk[t]) {
            this->tmap[qt_dest] = t;
            break;
        }
    }
    this->trk_bits[this->tmap[qt_dest]] = this->trk_bits[this->tmap[qt_src]];
    this->trk_byts[this->tmap[qt_dest]] = this->trk_byts[this->tmap[qt_src]];
    this->trk[this->tmap[qt_dest]] = new std::uint8_t[this->trk_byts[this->tmap[qt_dest]]];
    memcpy(this->trk[this->tmap[qt_dest]], this->trk[this->tmap[qt_src]], this->trk_byts[this->tmap[qt_dest]]);
}

void WozFile::createNewTrack(const std::uint8_t qt) {
    if (this->tmap[qt] != 0xFFu) { // track already exists
        return;
    }

    for (std::uint8_t t(0); t < C_QTRACK; ++t) {
        if (!this->trk[t]) {
            this->tmap[qt] = t;
            break;
        }
    }
    if (this->tmap[qt] == 0xFFu) {
        printf("Cannot create track %d\n", qt);
        return;
    }

    this->trk_bits[this->tmap[qt]] = calcNewTrackLengthBits(qt);
    this->trk_byts[this->tmap[qt]] = bytesForBits(this->trk_bits[this->tmap[qt]]);
    this->trk[this->tmap[qt]] = new std::uint8_t[this->trk_byts[this->tmap[qt]]];
    memset(this->trk[this->tmap[qt]], 0, this->trk_byts[this->tmap[qt]]);
}

/*
 * example:
 * tmap[]       track length in bits
 *  17 = 06 --> 50000
 *  18 = FF -X
 *  19 = 07 --> 50002
 *
 * calcNewTrackLengthBits(18) returns 50001
 */
std::uint32_t WozFile::calcNewTrackLengthBits(const std::uint8_t qt) {
    uint32_t t1 = 0;
    for (int t(qt-1); t >= 0; --t) {
        if (this->tmap[t] != 0xFFu) {
            t1 = this->trk_bits[this->tmap[t]];
            break;
        }
    }
    uint32_t t2 = 0;
    for (int t(qt+1); t < C_QTRACK; ++t) {
        if (this->tmap[t] != 0xFFu) {
            t2 = this->trk_bits[this->tmap[t]];
            break;
        }
    }
    // corner case: no existing tracks at all
    if (!t1 && !t2) {
        // 0xC780 yields COPY ][ PLUS disk speed of 200ms
        return 0xC780u;
    }
    // nominal case: average flanking tracks
    if (t1 && t2) {
        return (t1+t2)/2;
    }
    // odd cases: first or last track
    return t1 ? t1 : t2;
}