package game

import (
	"math/rand"
	"time"
)

var MasterFlowerPool = []string{
	// ── Originals ──────────────────────────────────────────────────────────
	"rose_red.svg", "rose_white.svg", "tulip_yellow.svg",
	"daisy_basic.svg", "lily_white.svg", "sunflower.svg",
	"orchid_purple.svg", "peony_pink.svg", "cherry_blossom.svg",
	"blue_hydrangea.svg",

	// ── Roses ──────────────────────────────────────────────────────────────
	"rose_pink.svg", "rose_yellow.svg", "rose_orange.svg",
	"rose_lavender.svg", "rose_peach.svg", "rose_crimson.svg",
	"rose_ivory.svg", "rose_coral.svg", "rose_magenta.svg",
	"rose_blush.svg",

	// ── Tulips ─────────────────────────────────────────────────────────────
	"tulip_red.svg", "tulip_pink.svg", "tulip_purple.svg",
	"tulip_white.svg", "tulip_orange.svg", "tulip_striped.svg",
	"tulip_parrot.svg", "tulip_peach.svg",

	// ── Lilies ─────────────────────────────────────────────────────────────
	"lily_pink.svg", "lily_orange.svg", "lily_yellow.svg",
	"lily_red.svg", "lily_purple.svg", "lily_calla_white.svg",
	"lily_calla_pink.svg", "lily_tiger.svg", "lily_stargazer.svg",

	// ── Daisies ────────────────────────────────────────────────────────────
	"daisy_white.svg", "daisy_pink.svg", "daisy_yellow.svg",
	"daisy_purple.svg", "daisy_gerbera_red.svg", "daisy_gerbera_orange.svg",
	"daisy_gerbera_yellow.svg", "daisy_gerbera_pink.svg", "daisy_oxeye.svg",

	// ── Orchids ────────────────────────────────────────────────────────────
	"orchid_white.svg", "orchid_pink.svg", "orchid_yellow.svg",
	"orchid_blue.svg", "orchid_red.svg", "orchid_spotted.svg",
	"orchid_dendrobium.svg", "orchid_cymbidium.svg",

	// ── Hydrangeas ─────────────────────────────────────────────────────────
	"hydrangea_pink.svg", "hydrangea_white.svg", "hydrangea_purple.svg",
	"hydrangea_green.svg", "hydrangea_red.svg",

	// ── Sunflowers & Dahlias ───────────────────────────────────────────────
	"sunflower_mini.svg", "sunflower_red.svg",
	"dahlia_red.svg", "dahlia_pink.svg", "dahlia_orange.svg",
	"dahlia_purple.svg", "dahlia_white.svg", "dahlia_yellow.svg",
	"dahlia_pompom.svg",

	// ── Peonies ────────────────────────────────────────────────────────────
	"peony_white.svg", "peony_red.svg", "peony_coral.svg",
	"peony_lavender.svg", "peony_blush.svg",

	// ── Wildflowers & Meadow ───────────────────────────────────────────────
	"lavender_sprig.svg", "wildflower_blue.svg", "wildflower_yellow.svg",
	"poppy_red.svg", "poppy_orange.svg", "poppy_white.svg",
	"cornflower_blue.svg", "forget_me_not.svg", "buttercup.svg",
	"clover_pink.svg", "dandelion.svg", "chamomile.svg",
	"cosmos_pink.svg", "cosmos_white.svg", "cosmos_purple.svg",

	// ── Tropical & Exotic ──────────────────────────────────────────────────
	"bird_of_paradise.svg", "anthurium_red.svg", "anthurium_pink.svg",
	"protea_pink.svg", "heliconia.svg", "ginger_flower.svg",
	"plumeria_white.svg", "plumeria_pink.svg", "hibiscus_red.svg",
	"hibiscus_yellow.svg",

	// ── Blossoms & Branches ────────────────────────────────────────────────
	"cherry_blossom_branch.svg", "plum_blossom.svg", "magnolia_white.svg",
	"magnolia_pink.svg", "wisteria_purple.svg", "wisteria_white.svg",
	"lily_of_the_valley.svg",

	// ── Greenery & Foliage ─────────────────────────────────────────────────
	"eucalyptus_sprig.svg", "fern_leaf.svg", "baby_breath.svg",
	"cotton_stem.svg", "thistle_purple.svg", "allium_purple.svg",
	"cattail.svg",

	// ── Romantic & Whimsical ───────────────────────────────────────────────
	"ranunculus_pink.svg", "ranunculus_white.svg", "ranunculus_orange.svg",
	"anemone_purple.svg", "anemone_red.svg", "anemone_white.svg",
	"sweet_pea_pink.svg", "sweet_pea_purple.svg",
	"snapdragon_pink.svg", "snapdragon_yellow.svg",
	"foxglove_purple.svg", "foxglove_white.svg",
	"lisianthus_purple.svg", "lisianthus_white.svg",
}

// DailyCount is how many flowers are shown each day.
// With 110 flowers in the pool, 8 gives good variety without overwhelming.
const DailyCount = 8

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

	dailySelection := poolCopy[:DailyCount]

	if hasSpecialUnlock {
		dailySelection[0] = "rare_glowing_lotus.svg"
	}

	return dailySelection
}
