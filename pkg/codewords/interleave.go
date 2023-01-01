package codewords

import (
	"fmt"
	"sort"

	"github.com/samlitowitz/go-qr/pkg/codewords/errorcorrection/reedsolomon"
)

type block struct {
	dataStart int
	dcc       int
	ecStart   int
	eccc      int
}

func blocksByDataCodewordCountInAscendingOrder(ecCfg *reedsolomon.Config) ([]*block, error) {
	dataPos := 0
	ecPos := 0

	blockCount := 0
	blockDCCs := make([]int, 0, 2)
	blocksByDCC := make(map[int][]*block)
	for _, ec := range ecCfg.ErrorCorrections {
		dcc := ec.K
		eccc := ec.N - dcc
		if _, ok := blocksByDCC[dcc]; !ok {
			blockDCCs = append(blockDCCs, dcc)
			blocksByDCC[dcc] = make([]*block, 0, 1)
		}

		for i := 0; i < ec.Blocks; i++ {
			blocksByDCC[dcc] = append(
				blocksByDCC[dcc],
				&block{
					dataStart: dataPos,
					dcc:       dcc,
					ecStart:   ecPos,
					eccc:      eccc,
				},
			)
			blockCount++
			dataPos += dcc
			ecPos += eccc
		}
	}
	sort.Ints(blockDCCs)

	blocks := make([]*block, 0, blockCount)
	for _, blockDataCodewordCount := range blockDCCs {
		if _, ok := blocksByDCC[blockDataCodewordCount]; !ok {
			return nil, fmt.Errorf("Invalid error correction configuration")
		}

		for _, block := range blocksByDCC[blockDataCodewordCount] {
			blocks = append(blocks, block)
		}
	}
	return blocks, nil
}

func Interleave(ecCfg *reedsolomon.Config, dataCodewords, ecCodewords []byte) ([]byte, error) {
	if len(ecCfg.ErrorCorrections) == 1 {
		return append(dataCodewords, ecCodewords...), nil
	}

	blocks, err := blocksByDataCodewordCountInAscendingOrder(ecCfg)
	if err != nil {
		return nil, err
	}

	prevDCC := 0
	prevDCCBlock := 0
	blockCount := len(blocks)
	totalDCC := len(dataCodewords)
	msg := make([]byte, totalDCC+len(ecCodewords))

	for curBlock, block := range blocks {
		for j := 0; j < prevDCC; j++ {
			msg[j*blockCount+curBlock] = dataCodewords[block.dataStart+j]
		}
		for j := prevDCC; j < block.dcc; j++ {
			msg[j*blockCount+curBlock-prevDCCBlock] = dataCodewords[block.dataStart+j]
		}
		for j := 0; j < block.eccc; j++ {
			msg[totalDCC+j*blockCount+curBlock] = ecCodewords[block.ecStart+j]
		}
	}

	return msg, nil
}
