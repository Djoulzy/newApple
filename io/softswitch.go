package io

const (
	// https://www.kreativekorp.com/miscpages/a2info/iomemory.shtml
	//
	// Comp:  O = Apple II+  E = Apple IIe  C = Apple IIc  G = Apple IIgs
	// Act:   R = Read       W = Write      7 = Bit 7      V = Value

	KBD          = 0x0000 // OECG  R   Last Key Pressed + 128
	STOREOFF     = 0x0000 //  ECG W    Use $C002-$C005 for Aux Memory
	STOREON      = 0x0001 //  ECG W    Use PAGE2 for Aux Memory
	RDMAINRAM    = 0x0002 //  ECG W    If 80STORE Off: Read Main Mem $0200-$BFFF
	RDCARDRAM    = 0x0003 //  ECG W    If 80STORE Off: Read Aux Mem $0200-$BFFF
	WRMAINRAM    = 0x0004 //  ECG W    If 80STORE Off: Write Main Mem $0200-$BFFF
	WRCARDRAM    = 0x0005 //  ECG W    If 80STORE Off: Write Aux Mem $0200-$BFFF
	SETSLOTCXROM = 0x0006 //  E G W    Peripheral ROM ($C100-$CFFF)
	SETINTCXROM  = 0x0007 //  E G W    Internal ROM ($C100-$CFFF)
	SETSTDZP     = 0x0008 //  ECG W    Main Stack and Zero Page
	SETALTZP     = 0x0009 //  ECG W    Aux Stack and Zero Page
	SETINTC3ROM  = 0x000A //  E G W    ROM in Slot 3
	SETSLOTC3ROM = 0x000B //  E G W    ROM in Aux Slot
	CLR80VID     = 0x000C //  ECG W    40 Columns
	SET80VID     = 0x000D //  ECG W    80 Columns
	CLRALTCHAR   = 0x000E //  ECG W    Primary Character Set
	SETALTCHAR   = 0x000F //  ECG W    Alternate Character Set
	KBDSTRB      = 0x0010 // OECG WR   Keyboard Strobe
	RDLCBNK2     = 0x0011 //  ECG  R7  Status of Selected $Dx Bank
	RDLCRAM      = 0x0012 //  ECG  R7  Status of $Dx ROM / $Dx RAM
	RDRAMRD      = 0x0013 //  ECG  R7  Status of Main/Aux RAM Reading
	RDRAMWRT     = 0x0014 //  ECG  R7  Status of Main/Aux RAM Writing
	RDCXROM      = 0x0015 //  E G  R7  Status of Periph/ROM Access
	RSTXINT      = 0x0015 //   C   R   Reset Mouse X0 Interrupt
	RDALTZP      = 0x0016 //  ECG  R7  Status of Main/Aux Stack and Zero Page
	RDC3ROM      = 0x0017 //  E G  R7  Status of Slot 3/Aux Slot ROM
	RSTYINT      = 0x0017 //   C   R   Reset Mouse Y0 Interrupt
	RD80STORE    = 0x0018 //  ECG  R7  Status of $C002-$C005/PAGE2 for Aux Mem
	RDVBL        = 0x0019 //  E G  R7  Vertical Blanking (E:1=drawing G:0=drawing)
	RSTVBL       = 0x0019 //   C   R   Reset Vertical Blanking Interrupt
	RDTEXT       = 0x001A //  ECG  R7  Status of Text/Graphics
	RDMIXED      = 0x001B //  ECG  R7  Status of Full Screen/Mixed Graphics
	RDPAGE2      = 0x001C //  ECG  R7  Status of Page 1/Page 2
	RDHIRES      = 0x001D //  ECG  R7  Status of LoRes/HiRes
	RDALTCHAR    = 0x001E //  ECG  R7  Status of Primary/Alternate Character Set
	RD80VID      = 0x001F //  ECG  R7  Status of 40/80 Columns
	TAPEOUT      = 0x0020 // OE    R   Toggle Cassette Tape Output
	MONOCOLOR    = 0x0021 //    G W 7  Color/Mono
	TBCOLOR      = 0x0022 //    G    V Screen Color: Low Nibble is BG, High Nibble is Text
	VGCINT       = 0x0023 //    G    V Video Graphics Controller Interrupts:
	//										b0-2=ext,scan,1sec enable b4-7=ext,scan,1sec,VGC
	MOUSEDATA = 0x0024 //       G    V Mouse Data: High Bit is Button, Other Bits are Movement
	KEYMODREG = 0x0025 //       G    V Modifier Keys: Bit 7: Command, Bit 6: Option,
	//										Bit 5: NotUsed, Bit 4: Keypad,
	//										Bit 3: Repeat,  Bit 2: Caps,
	//										Bit 1: Control, Bit 0: Shift
	DATAREG = 0x0026 //         G    V ADB Command/Data b0-2=# b3=valid b4=clr buf
	//										b5=reboot b6=abort b7=status
	KMSTATUS = 0x0027 //        G    V ADB Status: b0=cmdFull b1=mouseX b2=keyIntr b3=key
	//										b4=cmdIntr b5=data 6=mouseInt 7=mouse
	ROMBANK   = 0x0028 //    ????      ROM bank select toggle
	NEWVIDEO  = 0x0029 //       G    V New Video: 129=SHR, 1=None, Bit 6=Linearize, Bit 5=BW
	LANGSEL   = 0x002B //       G      Bit 3=Secondary Bit 4=50Hz Bits 5-7=Display Language
	CHARROM   = 0x002C //    ????      Addr for test mode read of character ROM
	SLTROMSEL = 0x002D //       G      Slot Register; Bits 1-7=use slot card
	VERTCNT   = 0x002E //    ????      Addr for read of video cntr bits V5-VB
	HORIZCNT  = 0x002F //    ????      Addr for read of video cntr bits VA-H0
	SPKR      = 0x0030 //    OECG  R   Toggle Speaker
	DISKREG   = 0x0031 //       G      Disk Interface: Bit 6=3.5 Bit 7=RWHead 1
	SCANINT   = 0x0032 //       G    V VGC Interrupt-Clear
	CLOCKDATA = 0x0033 //       G      Interface to Battery RAM (undocumented)
	CLOCKCTL  = 0x0034 //       G      b0-3=borderColor b5=stopBit b6=read b7=start
	SHADOW    = 0x0035 //       G      Inhibit Shadowing: Bit 6: I/O Memory, Bit 5: Alternate
	//										Display Mode, Bit 4: Auxilary HGR,
	//										Bit 3: Super HiRes, Bit 2: HiRes
	//										Page 2, Bit 1: HiRes Page 1,
	//										Bit 0: Text/LoRes
	CYAREG   = 0x0036 //        G      Bits 0-3=Disk Detect Bit 4=Shadow All Banks Bit 7=Fast
	BMAREG   = 0x0037 //        G      Bit 5=BW
	SCCBREG  = 0x0038 //        G      SCC Command Channel B
	SCCAREG  = 0x0039 //        G      SCC Command Channel A
	SCCBDATA = 0x003A //        G      SCC Data Channel B
	SCCADATA = 0x003B //        G      SCC Data Channel A
	SOUNDCTL = 0x003C //        G    V Sound Settings: Bits 0-3=Volume Bit 5=AutoIncr
	//										Bit 6=RAM Bit 7=Busy
	SOUNDDATA = 0x003D //       G      Sound Data
	SOUNDADRL = 0x003E //       G      Address Pointer L
	SOUNDADRH = 0x003F //       G      Address Pointer H
	STROBE    = 0x0040 //    OE    R   Game I/O Strobe Output
	RDXYMSK   = 0x0040 //      C   R7  Read X0/Y0 Interrupt
	RDVBLMSK  = 0x0041 //      C   R7  Read VBL Interrupt
	RDX0EDGE  = 0x0042 //      C   R7  Read X0 Edge Selector
	RDY0EDGE  = 0x0043 //      C   R7  Read Y0 Edge Selector
	MMDELTAX  = 0x0044 //       G    V Mega II Mouse Delta Movement X
	MMDELTAY  = 0x0045 //       G    V Mega II Mouse Delta Movement Y
	DIAGTYPE  = 0x0046 //    ????      Self or Burn-In diagdistics: Bit 7=burn-in diag
	INTFLAG   = 0x0046 //    ????      b0=IRQ b1=MMmov b2=MMbut b3=VBL b4=qsec
	//										b5=AN3 b6=mouse was down b7=mouse is down
	CLRVBLINT = 0x0047 //    ????      Clear VBL Interrupt
	CLRXYINT  = 0x0048 //    ????      Clear MM Interrupt
	RSTXY     = 0x0048 //      C  WR   Reset X and Y Interrupts
	EMUBYTE   = 0x004F //         WR   Emulation ID byte: write once, then read once for program
	//										being used, read again for version number.
	//										$FE=Bernie, $16=Sweet16, $4B=KEGS, $AB=Appleblossom
	TXTCLR   = 0x0050 //     OECG WR   Display Graphics
	TXTSET   = 0x0051 //     OECG WR   Display Text
	MIXCLR   = 0x0052 //     OECG WR   Display Full Screen
	MIXSET   = 0x0053 //     OECG WR   Display Split Screen
	TXTPAGE1 = 0x0054 //     OECG WR   Display Page 1
	TXTPAGE2 = 0x0055 //     OECG WR   If 80STORE Off: Display Page 2
	//                        ECG WR   If 80STORE On: Read/Write Aux Display Mem
	LORES     = 0x0056 //    OECG WR   Display LoRes Graphics
	HIRES     = 0x0057 //    OECG WR   Display HiRes Graphics
	CLRAN0    = 0x0058 //    OE G WR   If IOUDIS off: Annunciator 0 Off
	DISXY     = 0x0058 //      C  WR   If IOUDIS on: Mask X0/Y0 Move Interrupts
	SETAN0    = 0x0059 //    OE G WR   If IOUDIS off: Annunciator 0 On
	ENBXY     = 0x0059 //      C  WR   If IOUDIS on: Allow X0/Y0 Move Interrupts
	CLRAN1    = 0x005A //    OE G WR   If IOUDIS off: Annunciator 1 Off
	DISVBL    = 0x005A //      C  WR   If IOUDIS on: Disable VBL Interrupts
	SETAN1    = 0x005B //    OE G WR   If IOUDIS off: Annunciator 1 On
	ENVBL     = 0x005B //      C  WR   If IOUDIS on: Enable VBL Interrupts
	CLRAN2    = 0x005C //    OE G WR   If IOUDIS off: Annunciator 2 Off
	X0EDGEUP  = 0x005C //      C  WR   If IOUDIS on: Interrupt on X0 Rising
	SETAN2    = 0x005D //    OE G WR   If IOUDIS off: Annunciator 2 On
	X0EDGEDWN = 0x005D //      C  WR   If IOUDIS on: Interrupt on X0 Falling
	CLRAN3    = 0x005E //    OE G WR   If IOUDIS off: Annunciator 3 Off
	Y0EDGEUP  = 0x005E //      C  WR   If IOUDIS on: Interrupt on Y0 Rising
	DHIRESON  = 0x005E //     ECG WR   In 80-Column Mode: Double Width Graphics
	SETAN3    = 0x005F //    OE G WR   If IOUDIS off: Annunciator 3 On
	Y0EDGEDWN = 0x005F //      C  WR   If IOUDIS on: Interrupt on Y0 Falling
	DHIRESOFF = 0x005F //     ECG WR   In 80-Column Mode: Single Width Graphics
	TAPEIN    = 0x0060 //    OE    R7  Read Cassette Input
	//                         C   R7  Status of 80/40 Column Switch
	BUTN3  = 0x0060 //          G  R7  Switch Input 3
	RDBTN0 = 0x0061 //        ECG  R7  Switch Input 0 / Open Apple
	BUTN1  = 0x0062 //        E G  R7  Switch Input 1 / Solid Apple
	RD63   = 0x0063 //        E G  R7  Switch Input 2 / Shift Key
	//                         C   R7  Bit 7 = Mouse Button Not Pressed
	PADDL0   = 0x0064 //     OECG  R7  Analog Input 0
	PADDL1   = 0x0065 //     OECG  R7  Analog Input 1
	PADDL2   = 0x0066 //     OE G  R7  Analog Input 2
	RDMOUX1  = 0x0066 //       C   R7  Mouse Horiz Position
	PADDL3   = 0x0067 //     OE G  R7  Analog Input 3
	RDMOUY1  = 0x0067 //       C   R7  Mouse Vert Position
	STATEREG = 0x0068 //        G    V b0=INTCXROM b1=ROMBANK b2=LCBNK2 b3=RDROM
	//										b4=RAMWRT b5=RAMRD b6=PAGE2 b7=ALTZP
	TESTREG = 0x006D //      ????      Test Mode Bit Register
	CLRTM   = 0x006E //      ????      Clear Test Mode
	ENTM    = 0x006F //      ????      Enable Test Mode
	PTRIG   = 0x0070 //       E    R   Analog Input Reset
	//                         C  WR   Analog Input Reset + Reset VBLINT Flag
	BANKSEL = 0x0073 //       ECG W    Memory Bank Select for > 128K
	BLOSSOM = 0x0077 //           W    Appleblossom Special I/O Address
	//										$C1=Install clock driver, $CC=Get time in input
	//										buffer, $CF=get time in ProDOS global page.
	DSBLIOU   = 0x0078 //      C  W    Disable IOU Access
	ENBLIOU   = 0x0079 //      C  W    Enable IOU Access
	IOUDISON  = 0x007E //     EC  W    Disable IOU
	RDIOUDIS  = 0x007E //     EC   R7  Status of IOU Disabling
	IOUDISOFF = 0x007F //     EC  W    Enable IOU
	RDDHIRES  = 0x007F //     EC   R7  Status of Double HiRes
	RAMROB2   = 0x0080 //    OECG  R   Read RAM; no write; use $D000 bank 2
	ROMWB2    = 0x0081 //    OECG  RR  Read ROM; write RAM; use $D000 bank 2
	ROMROB2   = 0x0082 //    OECG  R   Read ROM; no write; use $D000 bank 2
	RAMRWB2   = 0x0083 //    OECG  RR  Read and write RAM; use $D000 bank 2
	// SAT_RDRAMB2 = 0x0084 //  OECG  R   Read RAM bank 2; no write
	// SAT_ROMIN   = 0x0085 //  OECG  RR  Read ROM; write RAM bank 2
	// RDROMNOWR   = 0x0086 //  OECG  R   Read ROM; no write
	// LCBANK2     = 0x0087 //  OECG  RR  Read/write RAM bank 2
	RAMROB1 = 0x0088 //      OECG  R   Read RAM; no write; use $D000 bank 1
	ROMWB1  = 0x0089 //      OECG  RR  Read ROM; write RAM; use $D000 bank 1
	ROMROB1 = 0x008A //      OECG  R   Read ROM; no write; use $D000 bank 1
	RAMRWB1 = 0x008B //      OECG  RR  Read and write RAM; use $D000 bank 1
	// RBANK1NOW   = 0x008C //  OECG  R   Read RAM bank 1; no write
	// RDROMWB1    = 0x008D //  OECG  RR  Read ROM; write RAM bank 1
	// RDROMNOWR   = 0x008E //  OECG  R   Read ROM; no write
	// RWB1        = 0x008F //  OECG  RR  Read/write RAM bank 1

	SLOT0_OFFSET  = 0x0090 // SLOTS
	SLOT1_OFFSET  = 0x0090
	SLOT2_OFFSET  = 0x00A0
	SLOT3_OFFSET  = 0x00B0
	SLOT4_OFFSET  = 0x00C0
	SLOT5_OFFSET  = 0x00D0
	SLOT6_OFFSET  = 0x00E0
	SLOT7_OFFSET  = 0x00F0

	// DRIVE OPERATIONS
	DRVSM0   = 0x00 // Q0
	DRVSM1   = 0x02 // Q1
	DRVSM2   = 0x04 // Q2
	DRVSM3   = 0x06 // Q3
	DRIVE    = 0x08 // Q4
	DRVSEL   = 0x0A // Q5
	DRVDATA  = 0x0C // Q6
	DRVWRITE = 0x0E // Q7
)
