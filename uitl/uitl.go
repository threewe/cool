package uitl

import (
	"fmt"
	"math/rand"
)
var randSeed int64

func init() {
	randSeed = 0
}
func RandRoomId(n int) string {
	randSeed++
	rand.Seed(randSeed)
	roomStr := ""
	for i := 0; i < n; i++ {
		d := rand.Int31n(10)
		roomStr =  fmt.Sprintf("%s%d", roomStr, d)
	}
	return roomStr
}
