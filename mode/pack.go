package mode

const BitsPerByte = 8

func PackInt(v, numberOfBitsToPack, unusedBitsInByte, byteInBuf int, buf *[]byte) (int, int, int, error) {
	var toCopy int
	toCopy = v
	numberOfBitsUnpacked := numberOfBitsToPack
	for numberOfBitsPacked := 0; numberOfBitsPacked < numberOfBitsToPack; {
		switch true {
		case numberOfBitsUnpacked == unusedBitsInByte:
			// copy
			(*buf)[byteInBuf] |= byte(toCopy)

			// bookkeeping
			byteInBuf++
			unusedBitsInByte = BitsPerByte
			numberOfBitsPacked += numberOfBitsUnpacked
			numberOfBitsUnpacked = 0

		case numberOfBitsUnpacked < unusedBitsInByte:
			// copy
			(*buf)[byteInBuf] |= byte(toCopy) << (unusedBitsInByte - numberOfBitsUnpacked)

			// bookkeeping
			unusedBitsInByte -= numberOfBitsUnpacked
			numberOfBitsPacked += numberOfBitsUnpacked
			numberOfBitsUnpacked = 0

		case numberOfBitsUnpacked > unusedBitsInByte:
			// copy
			(*buf)[byteInBuf] |= byte(toCopy >> (numberOfBitsUnpacked - unusedBitsInByte))

			// bookkeeping
			numberOfBitsPacked += unusedBitsInByte
			numberOfBitsUnpacked -= unusedBitsInByte
			byteInBuf++
			unusedBitsInByte = BitsPerByte
		}
	}
	return numberOfBitsUnpacked, unusedBitsInByte, byteInBuf, nil
}
