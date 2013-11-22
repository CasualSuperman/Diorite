package multiverse

// The colors of mana that exist in the Multiverse.
var ManaColors = struct {
	Colorless, White, Blue, Black, Red, Green ManaColor
}{0, 1, 2, 4, 8, 16}

// The borders that cards have.
var BorderColors = struct {
	White, Black, Silver BorderColor
}{1, 2, 3}

// Rarities of cards.
var Rarities = struct {
	Common, Uncommon, Rare, Mythic, Basic, Special Rarity
}{1, 2, 3, 4, 5, 6}

// Set types.
var SetTypes = struct {
	Core, Expansion, Reprint, Box, Un, FromTheVault, PremiumDeck, DuelDeck, Starter, Commander, Planechase, Archenemy SetType
}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

// SuperType values.
var SuperTypes = struct {
	Basic, Elite, Legendary, Ongoing, Snow, World SuperType
}{32, 16, 8, 4, 2, 1}

// Type values.
var Types = struct {
	Artifact, Creature, Enchantment, Instant, Land, Planeswalker, Sorcery, Tribal Type
}{128, 64, 32, 16, 8, 4, 2, 1}
