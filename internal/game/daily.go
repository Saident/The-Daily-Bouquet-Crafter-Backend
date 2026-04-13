package game

import (
	"math/rand"
	"time"
)

var MasterFlowerPool = []string{
	"rose_red.png", "rose_white.png", "tulip_yellow.png",
	"daisy_basic.png", "lily_white.png", "sunflower.png",
	"orchid_purple.png", "peony_pink.png", "cherry_blossom.png",
	"blue_hydrangea.png",
}

func GetDailyInventory(hasSpecialUnlock bool) []string {
	loc := time.FixedZone("UTC+7", 7*60*60)
	now := time.Now().In(loc)
	seedDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	r := rand.New(rand.NewSource(seedDate.Unix()))

	poolCopy := make([]string, len(MasterFlowerPool))
	copy(poolCopy, MasterFlowerPool)

	r.Shuffle(len(poolCopy), func(i, j int) {
		poolCopy[i], poolCopy[j] = poolCopy[j], poolCopy[i]
	})

	dailySelection := poolCopy[:5]

	if hasSpecialUnlock {
		dailySelection[0] = "rare_glowing_lotus.gif"
	}

	return dailySelection
}
