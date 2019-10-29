package main

func isRespawnPosition(r, g, b, a uint32) bool {
	if r == 0 && g == 0 && b == 255 {
		return true
	}

	return false
}

func isNormalGround(r, g, b, a uint32) bool {
	if r == 0 && g == 0 && b == 0 {
		return true
	}

	return false
}

func isWalkableGrass(r, g, b, a uint32) bool {
	if r == 127 && g == 127 && b == 127 {
		return true
	}

	return false
}

func isNormalGrass(r, g, b, a uint32) bool {
	if a == 0 || (r == 255 && g == 255 && b == 255) {
		return true
	}

	return false
}
