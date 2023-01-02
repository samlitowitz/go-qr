package symbol

type Mode int

const (
	M_Numeric      Mode = iota
	M_Alphanumeric Mode = iota
	M_Byte         Mode = iota
)

//  ec level, mode, characters

func ChooseVersion(eclevel ErrorCorrectionLevel, m Mode, charLen int) {

}
