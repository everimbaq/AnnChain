// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package vm

import (
	"fmt"
	"math/big"

	"github.com/dappledger/AnnChain/genesis/eth/common"
	"github.com/dappledger/AnnChain/genesis/eth/common/math"
	"github.com/dappledger/AnnChain/genesis/eth/params"
)

// Type is the VM type accepted by **NewVm**
type Type byte

const (
	StdVmTy Type = iota // Default standard VM
	JitVmTy             // LLVM JIT VM
	MaxVmTy
)

var (
	ADDRESS20 = []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}
	ADDRESS32 = []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}
)

var (
	Pow256 = common.BigPow(2, 256) // Pow256 is 2**256

	U256 = common.U256 // Shortcut to common.U256
	S256 = common.S256 // Shortcut to common.S256

	Zero = common.Big0 // Shortcut to common.Big0
	One  = common.Big1 // Shortcut to common.Big1

	max = big.NewInt(math.MaxInt64) // Maximum 64 bit integer
)

// calculates the memory size required for a step
func calcMemSize(off, l *big.Int) *big.Int {
	if l.Cmp(common.Big0) == 0 {
		return common.Big0
	}

	return new(big.Int).Add(off, l)
}

// getDataBig returns a slice from the data based on the start and size and pads
// up to size with zero's. This function is overflow safe.
func getDataBig(data []byte, start *big.Int, size *big.Int) []byte {
	dlen := big.NewInt(int64(len(data)))

	s := math.BigMin(start, dlen)
	e := math.BigMin(new(big.Int).Add(s, size), dlen)
	return common.RightPadBytes(data[s.Uint64():e.Uint64()], int(size.Uint64()))
}

// bigUint64 returns the integer casted to a uint64 and returns whether it
// overflowed in the process.
func bigUint64(v *big.Int) (uint64, bool) {
	return v.Uint64(), v.BitLen() > 64
}

// calculates the quadratic gas
func quadMemGas(mem *Memory, newMemSize, gas *big.Int) {
	if newMemSize.Cmp(common.Big0) > 0 {
		newMemSizeWords := toWordSize(newMemSize)
		newMemSize.Mul(newMemSizeWords, u256(32))

		if newMemSize.Cmp(u256(int64(mem.Len()))) > 0 {
			// be careful reusing variables here when changing.
			// The order has been optimised to reduce allocation
			oldSize := toWordSize(big.NewInt(int64(mem.Len())))
			pow := new(big.Int).Exp(oldSize, common.Big2, Zero)
			linCoef := oldSize.Mul(oldSize, params.MemoryGas)
			quadCoef := new(big.Int).Div(pow, params.QuadCoeffDiv)
			oldTotalFee := new(big.Int).Add(linCoef, quadCoef)

			pow.Exp(newMemSizeWords, common.Big2, Zero)
			linCoef = linCoef.Mul(newMemSizeWords, params.MemoryGas)
			quadCoef = quadCoef.Div(pow, params.QuadCoeffDiv)
			newTotalFee := linCoef.Add(linCoef, quadCoef)

			fee := newTotalFee.Sub(newTotalFee, oldTotalFee)
			gas.Add(gas, fee)
		}
	}
}

// Simple helper
func u256(n int64) *big.Int {
	return big.NewInt(n)
}

// Mainly used for print variables and passing to Print*
func toValue(val *big.Int) interface{} {
	// Let's assume a string on right padded zero's
	b := val.Bytes()
	if b[0] != 0 && b[len(b)-1] == 0x0 && b[len(b)-2] == 0x0 {
		return string(b)
	}

	return val
}

// getData returns a slice from the data based on the start and size and pads
// up to size with zero's. This function is overflow safe.
func getData(data []byte, start, size *big.Int) []byte {
	length := uint64(len(data))
	if start.Uint64() > length {
		start = new(big.Int).SetUint64(length)
	}
	start.Add(start, size)
	end := start
	if end.Uint64() > length {
		end = new(big.Int).SetUint64(length)
	}
	return common.RightPadBytes(data[start.Int64():end.Int64()], int(size.Int64()))
}

// useGas attempts to subtract the amount of gas and returns whether it was
// successful
func useGas(gas, amount *big.Int) bool {
	if gas.Cmp(amount) < 0 {
		return false
	}
	// Sub the amount of gas from the remaining
	gas.Sub(gas, amount)
	return true
}
