package main

import (
	"fmt"
	tl "github.com/JoelOtter/termloop"
	"math/rand"
	"time"
)

var rando = rand.New(rand.NewSource(time.Now().UnixNano()))

type Sprite struct {
	x         int
	y         int
	w         int
	h         int
	canvases  []*tl.Canvas // all preloaded images, loopSize-worth of each direction in order
	loopSize  int
	loopI     int
	direction int // up=0, right=1, down=2, left=3
	player    *Player
	enemy     *Enemy
	level     *tl.BaseLevel
}

func NewSprite(canvases []*tl.Canvas, loopSize, direction, loopI, x, y, w, h int, level *tl.BaseLevel) *Sprite {
	s := Sprite{canvases: canvases, loopSize: loopSize, direction: direction, loopI: loopI, x: x, y: y, w: w, h: h, level: level}
	return &s
}

func (sprite *Sprite) IsPlayer() bool {
	return sprite.player != nil
}

func (sprite *Sprite) IsEnemy() bool {
	return sprite.enemy != nil
}

func (sprite *Sprite) Draw(s *tl.Screen) {
	if (sprite.IsPlayer() && sprite.player.health <= 0) || (sprite.IsEnemy() && sprite.enemy.health <= 0) {
		sprite.level.RemoveEntity(sprite.CurrentEntity())
	} else {
		if sprite.IsPlayer() {
			screenWidthidth, screenh := s.Size()
			sprite.level.SetOffset(screenWidthidth/2-7-sprite.x, screenh/5-sprite.y)
		}
		sprite.CurrentEntity().Draw(s)
	}
}

func (sprite *Sprite) Position() (int, int)     { return sprite.x, sprite.y }
func (sprite *Sprite) Size() (int, int)         { return sprite.w, sprite.h }
func (sprite *Sprite) SetPosition(x int, y int) { sprite.x = x; sprite.y = y }
func (sprite *Sprite) UpdateAnimation(direction int) {
	baseI := sprite.loopSize * direction
	if direction != sprite.direction {
		sprite.direction = direction
		sprite.loopI = baseI
	} else {
		sprite.loopI = sprite.loopI + 1
		if sprite.loopI == baseI+sprite.loopSize {
			sprite.loopI = baseI
		}
	}
}

func (sprite *Sprite) CurrentEntity() *tl.Entity {
	c := sprite.canvases[sprite.loopI]
	e := tl.NewEntityFromCanvas(sprite.x, sprite.y, *c)
	return e
}

func (sprite *Sprite) Tick(ev tl.Event) {
	if ev.Type == tl.EventKey {
		x, y := sprite.Position()
		var direction = -1
		if sprite.IsPlayer() {
			switch ev.Key {
			case tl.KeyArrowRight:
				direction = 1
			case tl.KeyArrowLeft:
				direction = 3
			case tl.KeyArrowUp:
				direction = 0
			case tl.KeyArrowDown:
				direction = 2
			}
		} else {
			direction = rando.Intn(4)
		}
		if direction > -1 {
			switch direction {
			case 0:
				y -= 1
			case 1:
				x += 1
			case 2:
				y += 1
			case 3:
				x -= 1
			}
			sprite.UpdateAnimation(direction)
		}
		if x < 0 {
			x = 0
		}
		if y < 0 {
			y = 0
		}
		sprite.SetPosition(x, y)
	}
}

func (sprite *Sprite) Collide(other tl.Physical) {
	if otherSprite, ok := other.(*Sprite); ok {
		if sprite.IsPlayer() && otherSprite.IsEnemy() {
			sprite.player.attacker = otherSprite.enemy
		} else if otherSprite.IsPlayer() && sprite.IsEnemy() {
			otherSprite.player.attacker = sprite.enemy
		}
	}
}

type SpellEffect struct {
	x         int
	y         int
	w         int
	h         int
	direction int
	speed     int
	turns     int
	duration  int
	done      bool
	damage    int
	canvas    *tl.Canvas
	e         *tl.Entity
	level     *tl.BaseLevel
}

func NewSpellEffect(sprite *Sprite, canvas *tl.Canvas, level *tl.BaseLevel) *SpellEffect {
	var x, y, h, w int
	switch sprite.direction {
	case 0:
		h = 7
		w = 2
		x = sprite.x + sprite.w/2 - 1
		y = sprite.y - h - 1
	case 1:
		h = 2
		w = 7
		x = sprite.x + sprite.w + 1
		y = sprite.y + sprite.h/2
	case 2:
		h = 7
		w = 2
		x = sprite.x + sprite.w/2 - 1
		y = sprite.y + sprite.h + 1
	case 3:
		h = 2
		w = 7
		x = sprite.x - w
		y = sprite.y + sprite.h/2
	}
	se := SpellEffect{x: x, y: y, w: w, h: h, direction: sprite.direction, canvas: canvas, duration: 20, speed: 7, damage: 5, level: level}
	return &se
}

func (spellEffect *SpellEffect) Position() (int, int) { return spellEffect.x, spellEffect.y }
func (spellEffect *SpellEffect) Size() (int, int)     { return spellEffect.w, spellEffect.h }

func (spellEffect *SpellEffect) Draw(s *tl.Screen) {
	spellEffect.e = tl.NewEntityFromCanvas(spellEffect.x, spellEffect.y, *spellEffect.canvas)
	if spellEffect.turns == spellEffect.duration || spellEffect.done {
		spellEffect.level.RemoveEntity(spellEffect)
	} else {
		spellEffect.e.Draw(s)
	}
}

func (spellEffect *SpellEffect) Tick(ev tl.Event) {
	spellEffect.turns = spellEffect.turns + 1
	switch spellEffect.direction {
	case 0:
		spellEffect.y = spellEffect.y - spellEffect.speed
	case 1:
		spellEffect.x = spellEffect.x + spellEffect.speed
	case 2:
		spellEffect.y = spellEffect.y + spellEffect.speed
	case 3:
		spellEffect.x = spellEffect.x - spellEffect.speed
	}
}

func (spellEffect *SpellEffect) Collide(other tl.Physical) {
	if sprite, ok := other.(*Sprite); ok {
		if sprite.IsEnemy() && !spellEffect.done {
			sprite.enemy.TakeDamage(spellEffect)
			spellEffect.done = true
		}
	}
}

type Player struct {
	health            int
	maxHealth         int
	mana              int
	maxMana           int
	gold              int
	experience        int
	experienceToLevel int
	level             int
	isCasting         bool
	attacker          *Enemy
	portrait          *tl.Canvas
	sprite            *Sprite
	spellCanvases     []*tl.Canvas
}

func NewPlayer(portrait *tl.Canvas, spellCanvases []*tl.Canvas, sprite *Sprite) *Player {
	p := Player{health: 100, maxHealth: 100, mana: 100, maxMana: 100, experienceToLevel: 100, portrait: portrait, level: 1, sprite: sprite, spellCanvases: spellCanvases}
	return &p
}

var StatsBG = tl.RgbTo256Color(170, 170, 170)

func (player *Player) Draw(s *tl.Screen) {
	screenWidthidth, screenh := s.Size()
	x := player.sprite.x + screenWidthidth/2 - 30
	y := player.sprite.y + screenh - 25
	bg := tl.NewRectangle(x, y, x+20, y+10, StatsBG)
	bg.Draw(s)

	health := tl.NewText(x+1, y+1, fmt.Sprintf("%3.f%% health", float32(player.health)/float32(player.maxHealth)*100), tl.ColorRed, StatsBG)
	health.Draw(s)

	mana := tl.NewText(x+27, y+1, fmt.Sprintf("%3.f%% mana", float32(player.mana)/float32(player.maxMana)*100), tl.ColorBlue, StatsBG)
	mana.Draw(s)

	gold := tl.NewText(x+1, y+12, fmt.Sprintf("%d gold", player.gold), tl.ColorYellow, StatsBG)
	gold.Draw(s)

	experience := tl.NewText(x+29, y+12, fmt.Sprintf("%3.f%% xp", float32(player.experience)/float32(player.experienceToLevel)*100), tl.ColorMagenta, StatsBG)
	experience.Draw(s)

	if player.isCasting {
		player.isCasting = false

		newSpellEffect := NewSpellEffect(player.sprite, player.spellCanvases[player.sprite.direction], player.sprite.level)
		player.sprite.level.AddEntity(newSpellEffect)
	}

	e := tl.NewEntityFromCanvas(x+12, y+1, *player.portrait)
	e.Draw(s)
}

func (player *Player) Tick(ev tl.Event) {
	if ev.Type == tl.EventKey {
		switch ev.Key {
		case tl.KeySpace:
			player.CastSpell()
		}
	}

	if player.attacker != nil {
		player.TakeDamage(player.attacker)
		player.attacker = nil
	}
}

func (player *Player) CastSpell() {
	if player.mana >= 5 {
		player.mana = player.mana - 5
		player.isCasting = true
	}
}

func (player *Player) TakeDamage(enemy *Enemy) {
	player.health = player.health - enemy.damage
}

type Enemy struct {
	damage    int
	health    int
	maxHealth int
}

func NewEnemy(damage, maxHealth int) *Enemy {
	e := Enemy{damage: damage, health: maxHealth, maxHealth: maxHealth}
	return &e
}

func (enemy *Enemy) TakeDamage(spellEffect *SpellEffect) {
	enemy.health = enemy.health - spellEffect.damage
}

func main() {
	g := tl.NewGame()

	level := tl.NewBaseLevel(tl.Cell{
		Bg: tl.ColorBlack,
	})

	background := tl.NewEntityFromCanvas(-50, -50, *tl.BackgroundCanvasFromFile("artwork/background/forest.png"))
	level.AddEntity(background)

	charSelect := "m_mage"
	var playerCanvases []*tl.Canvas
	for i := 0; i < 12; i++ {
		playerCanvases = append(playerCanvases, tl.BackgroundCanvasFromFile(fmt.Sprintf("artwork/sprites/player/%s/%d.png", charSelect, i)))
	}
	playerSprite := NewSprite(playerCanvases, 3, 2, 7, 0, 0, 14, 28, level)
	level.AddEntity(playerSprite)

	var enemyCanvases []*tl.Canvas
	for i := 0; i < 12; i++ {
		enemyCanvases = append(enemyCanvases, tl.BackgroundCanvasFromFile(fmt.Sprintf("artwork/sprites/enemy/flower/%d.png", i)))
	}

	enemySprite := NewSprite(enemyCanvases, 3, 2, 7, 70, 15, 14, 28, level)
	enemySprite.enemy = NewEnemy(5, 10)
	level.AddEntity(enemySprite)

	portrait := tl.BackgroundCanvasFromFile(fmt.Sprintf("artwork/sprites/player/%s/portrait.png", charSelect))
	var spellCanvases []*tl.Canvas
	for i := 0; i < 4; i++ {
		spellCanvases = append(spellCanvases, tl.BackgroundCanvasFromFile(fmt.Sprintf("artwork/spells/fireball/%d.png", i)))
	}
	player := NewPlayer(portrait, spellCanvases, playerSprite)
	playerSprite.player = player
	level.AddEntity(player)

	g.Screen().SetLevel(level)
	g.Start()
}
