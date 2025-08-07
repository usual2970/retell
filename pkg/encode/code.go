package encode

import "fmt"

func Code(length int) string {
	max := 10
	for i := 0; i < length-1; i++ {
		max *= 10
	}

	rs := rng.Intn(max)
	return fmt.Sprintf("%0*d", length, rs)
}
