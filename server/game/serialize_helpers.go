package game

/*
	NOTE: All serializations are in big-endian form
				(most significant byte, MSB, first)
*/

/*
Get the byte at a particular index (0 = least significant byte,
1 = second least, etc.)
*/
func getByte[T uint8 | uint16 | uint32](num T, byteIdx int) byte {

	/*
		Uses bitwise operation magic (not really, look up how the >> and &
		operators work if you're interested)
	*/
	return byte((num >> (8 * byteIdx)) & 0xff)
}

// Serialize an individual byte (this should be simple, just getByte call)
func serUint8(num uint8, outputBuf []byte, startIdx int) int {
	outputBuf[startIdx] = getByte(num, 0)
	return startIdx + 1
}

// Serialize a uint16 (two getByte calls)
func serUint16(num uint16, outputBuf []byte, startIdx int) int {

	// Loop over each of the 2 bytes within the row (MSB first)
	for byteIdx := 1; byteIdx >= 0; byteIdx-- {

		// Serialize the byte
		outputBuf[startIdx] = getByte(num, byteIdx)

		// Add 1 to the start index, to prepare for serializing the next byte
		startIdx++
	}
	return startIdx
}

// Serialize a uint32 (four getByte calls)
func serUint32(num uint32, outputBuf []byte, startIdx int) int {

	// Loop over each of the 4 bytes within the row (MSB first)
	for byteIdx := 3; byteIdx >= 0; byteIdx-- {

		// Serialize the byte
		outputBuf[startIdx] = getByte(num, byteIdx)

		// Add 1 to the start index, to prepare for serializing the next byte
		startIdx++
	}

	// Return the starting index of the next field
	return startIdx
}

// Serialize a location (no getByte calls, serialized manually)
func serLocation(loc *locationState, outputBuf []byte, startIdx int) int {

	// Lock the location state (so no writes can be done simultaneously)
	loc.RLock()
	defer loc.RUnlock()

	// Cover each coordinate of the location, one at a time
	outputBuf[startIdx+0] = byte((loc.dRow() << 6) | loc.row)
	outputBuf[startIdx+1] = byte((loc.dCol() << 6) | loc.col)

	// Return the starting index of the next field
	return startIdx + 2
}
