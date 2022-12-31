package reedsolomon

// GenerateGeneratorPolynomial returns a generator polynomial with n+1 terms in the Galois Field GF(256) in alpha notation.
// Terms are from least significant, x^0, at index 0 to most significant, x^n, at index n.
func GenerateGeneratorPolynomial(n int) []byte {
	g := make([]byte, n+1)
	for terms := 1; terms < n; terms++ {
		for j := terms; j > 0; j-- {
			a := gfMult(g[j], byte(terms))
			b := gfMult(g[j-1], 0)
			g[j] = logTable[expTable[a]^expTable[b]]
		}
		g[0] = gfMult(g[0], byte(terms))
	}
	return g
}
