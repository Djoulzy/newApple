package config

// Globals : Partie globale du fichier de conf
type Globals struct {
	StartLogging bool
	FileLog      string
	Disassamble  bool
	LogLevel     int
	Debug        bool
	Display      bool
	Model        string
	ColorDisplay bool
	CPUModel     string
	Mhz          int
}

type Slots struct {
	Slot1   string
	Slot2   string
	Slot3   string
	Slot4   string
	Slot5   string
	Slot6   string
	Slot7   string
	Catalog [8]string
}

type Disks struct {
	Disk1 string
	Disk2 string
}

type DebugMode struct {
	Breakpoint uint16
	Dump       uint16
	Zone       int
}

// ConfigData : Data structure du fichier de conf
type ConfigData struct {
	Globals
	Slots
	Disks
	DebugMode
}
