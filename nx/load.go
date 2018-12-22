package nx

import (
	"fmt"

	"github.com/Hucaru/gonx"
)

var items map[int32]Item
var maps map[int32]Map
var mobs map[int32]Mob
var playerSkills map[int32][]PlayerSkill
var mobSkills map[int32][]MobSkill

// LoadFile into useable types
func LoadFile(fname string) {
	nodes, textLookup, _, _, err := gonx.Parse(fname)

	if err != nil {
		panic(err)
	}

	items = extractItems(nodes, textLookup)
	maps = extractMaps(nodes, textLookup)
	mobs = extractMobs(nodes, textLookup)
	playerSkills, mobSkills = extractSkills(nodes, textLookup)
}

// GetItem from loaded nx
func GetItem(id int32) (Item, error) {
	if _, ok := items[id]; !ok {
		return Item{}, fmt.Errorf("invalid item id: %v", id)
	}

	return items[id], nil
}

// GetMap from loaded nx
func GetMap(id int32) (Map, error) {
	if _, ok := maps[id]; !ok {
		return Map{}, fmt.Errorf("invalid map id: %v", id)
	}

	return maps[id], nil
}

// GetMob from loaded nx
func GetMob(id int32) (Mob, error) {
	if _, ok := mobs[id]; !ok {
		return Mob{}, fmt.Errorf("invalid mob id: %v", id)
	}

	return mobs[id], nil
}

// GetPlayerSkill from loaded nx
func GetPlayerSkill(id int32) ([]PlayerSkill, error) {
	if _, ok := playerSkills[id]; !ok {
		return []PlayerSkill{}, fmt.Errorf("Invalid player skill id: %v", id)
	}

	return playerSkills[id], nil
}

// GetMobSkill from loaded nx
func GetMobSkill(id int32) ([]MobSkill, error) {
	if _, ok := mobSkills[id]; !ok {
		return []MobSkill{}, fmt.Errorf("Invalid mob skill id: %v", id)
	}

	return mobSkills[id], nil
}
