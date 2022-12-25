package mode

type Config struct {
	ModeIndicatorLength  int // The length of the mode indicator in bits
	ModeIndicator        byte
	CharacterCountLength int // The length of the character count in bits
}
