package game

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// SlotState tracks destruction level of an individual window/balcony section.
type SlotState int

const (
	SlotIntact SlotState = iota
	SlotDamaged
	SlotBroken
	SlotEmptyHole
)

// WindowSlot represents a single smashable room/balcony on a floor.
type WindowSlot struct {
	FloorIdx    int
	ColIdx      int
	X, Y        float64
	W, H        float64
	State       SlotState
	Health      float64
	MaxHealth   float64
	HasBalcony  bool
	HasAwning   bool
	AwningColor color.RGBA
	Victim      *Citizen
	ItemGyro    bool
	ItemCoffee  bool
}

// Building represents an Athenian apartment block (polykatoikia) or Acropolis monument.
type Building struct {
	Spec          BuildingSpec
	X, Y          float64
	Width, Height float64
	Floors        int
	Cols          int
	Slots         [][]*WindowSlot // [floor][col]
	Health        float64
	MaxHealth     float64
	Collapsed     bool
	CollapseProg  float64 // 0.0 intact -> 1.0 fully sunk into rubble
	SinkY         float64
	ShakeOffsetX  float64
	ShakeOffsetY  float64
	WallColor     color.RGBA
	AwningColor   color.RGBA
	IsMarble      bool
	SignTitle     string
	Subtitle      string
	RubbleHeight  float64
	SolarHeaterX  float64
	SolarHeaterY  float64
}

// NewBuilding constructs and lays out a building according to the stage spec.
func NewBuilding(spec BuildingSpec, posX, groundY float64) *Building {
	const (
		floorHeight = 64.0
		colWidth    = 58.0
		paddingX    = 18.0
		roofHeight  = 30.0
	)

	bWidth := float64(spec.Cols)*colWidth + paddingX*2.0
	bHeight := float64(spec.Floors)*floorHeight + roofHeight
	bY := groundY - bHeight

	b := &Building{
		Spec:         spec,
		X:            posX,
		Y:            bY,
		Width:        bWidth,
		Height:       bHeight,
		Floors:       spec.Floors,
		Cols:         spec.Cols,
		WallColor:    spec.WallColor,
		AwningColor:  spec.AwningColor,
		IsMarble:     spec.IsMarble,
		SignTitle:    spec.SignTitle,
		Subtitle:     spec.Subtitle,
		SolarHeaterX: posX + bWidth*0.5,
		SolarHeaterY: bY - 14.0,
	}

	b.Slots = make([][]*WindowSlot, spec.Floors)
	totalHealth := 0.0

	for f := 0; f < spec.Floors; f++ {
		b.Slots[f] = make([]*WindowSlot, spec.Cols)
		floorY := bY + roofHeight + float64(f)*floorHeight

		for c := 0; c < spec.Cols; c++ {
			slotX := posX + paddingX + float64(c)*colWidth
			slotW := colWidth - 8.0
			slotH := floorHeight - 12.0

			slotHealth := 120.0
			totalHealth += slotHealth

			slot := &WindowSlot{
				FloorIdx:    f,
				ColIdx:      c,
				X:           slotX,
				Y:           floorY + 6.0,
				W:           slotW,
				H:           slotH,
				State:       SlotIntact,
				Health:      slotHealth,
				MaxHealth:   slotHealth,
				HasBalcony:  !spec.IsMarble && f > 0,
				HasAwning:   !spec.IsMarble && f > 0 && rand.Float64() < 0.85,
				AwningColor: spec.AwningColor,
			}

			// Random gyros or coffee pickup inside windows
			if rand.Float64() < 0.22 {
				slot.ItemGyro = true
			} else if rand.Float64() < 0.22 {
				slot.ItemCoffee = true
			}

			b.Slots[f][c] = slot
		}
	}

	b.MaxHealth = totalHealth
	b.Health = totalHealth
	return b
}

// DamageSlot inflicts punch damage to a window slot, returning damage dealt and if destroyed.
func (b *Building) DamageSlot(f, c int, dmg float64, ps *ParticleSystem) (float64, bool) {
	if b.Collapsed || f < 0 || f >= b.Floors || c < 0 || c >= b.Cols {
		return 0, false
	}

	slot := b.Slots[f][c]
	if slot.State == SlotEmptyHole {
		return 0, false
	}

	actualDmg := math.Min(slot.Health, dmg)
	slot.Health -= actualDmg
	b.Health -= actualDmg

	// Debris burst
	ps.SpawnBuildingDebris(slot.X+slot.W/2, slot.Y+slot.H/2, 6, b.WallColor)

	destroyed := false
	if slot.Health <= 0 {
		slot.State = SlotEmptyHole
		destroyed = true
		ps.SpawnBuildingDebris(slot.X+slot.W/2, slot.Y+slot.H/2, 14, b.WallColor)
	} else if slot.Health <= slot.MaxHealth*0.35 {
		slot.State = SlotBroken
	} else if slot.Health <= slot.MaxHealth*0.70 {
		slot.State = SlotDamaged
	}

	// Check total collapse threshold (e.g. 25% health remaining triggers collapse)
	if b.Health <= b.MaxHealth*0.20 && !b.Collapsed {
		b.TriggerCollapse(ps)
	}

	return actualDmg, destroyed
}

// TriggerCollapse initiates the building sinking collapse animation.
func (b *Building) TriggerCollapse(ps *ParticleSystem) {
	if b.Collapsed {
		return
	}
	b.Collapsed = true
	b.CollapseProg = 0.01

	// Massive dust cloud and rubble bursts across the building
	for f := 0; f < b.Floors; f++ {
		for c := 0; c < b.Cols; c++ {
			s := b.Slots[f][c]
			ps.SpawnBuildingDebris(s.X+s.W/2, s.Y+s.H/2, 8, b.WallColor)
		}
	}

	audio := GetAudioSystem()
	audio.PlaySFX("building_collapse", 0)
}

// Update advances destruction animations and collapse physics.
func (b *Building) Update(dt float64, ps *ParticleSystem) {
	if !b.Collapsed {
		return
	}

	if b.CollapseProg < 1.0 {
		b.CollapseProg += dt * 0.45
		if b.CollapseProg > 1.0 {
			b.CollapseProg = 1.0
		}

		// Sinking distance towards ground (660)
		groundY := 660.0
		targetSink := (groundY - b.Y) - 48.0 // Leaves 48px rubble
		b.SinkY = targetSink * math.Pow(b.CollapseProg, 1.8)

		// Violent earthquake shaking while collapsing
		shakeMag := (1.0 - b.CollapseProg) * 8.0
		b.ShakeOffsetX = (rand.Float64()*2.0 - 1.0) * shakeMag
		b.ShakeOffsetY = (rand.Float64()*2.0 - 1.0) * shakeMag

		// Continual dust puffs during collapse
		if rand.Float64() < 0.4 {
			dustX := b.X + rand.Float64()*b.Width
			dustY := b.Y + b.SinkY + rand.Float64()*b.Height*0.5
			ps.SpawnBuildingDebris(dustX, dustY, 3, b.WallColor)
		}
	} else {
		b.ShakeOffsetX = 0
		b.ShakeOffsetY = 0
	}
}

// Draw renders the building facade, windows, awnings, solar heaters, and signs.
func (b *Building) Draw(screen *ebiten.Image, pixelImg *ebiten.Image, signFace, subFace text.Face) {
	drawY := b.Y + b.SinkY + b.ShakeOffsetY
	drawX := b.X + b.ShakeOffsetX

	var op ebiten.DrawImageOptions

	// 1. Building Main Wall / Facade
	op.GeoM.Reset()
	op.GeoM.Scale(b.Width, b.Height)
	op.GeoM.Translate(drawX, drawY)
	op.ColorScale.Reset()
	op.ColorScale.Scale(
		float32(b.WallColor.R)/255.0,
		float32(b.WallColor.G)/255.0,
		float32(b.WallColor.B)/255.0,
		1.0,
	)
	screen.DrawImage(pixelImg, &op)

	// Wall border / pillars outline
	drawRectOutline(screen, pixelImg, drawX, drawY, b.Width, b.Height, color.RGBA{45, 40, 35, 255}, 2.0)

	// 2. Acropolis Classical Columns or Polykatoikia Ground Shop Sign
	if b.IsMarble {
		// Draw Classical Columns (Parthenon style fluted marble pillars)
		columnWidth := 14.0
		colCount := b.Cols + 2
		spacing := (b.Width - float64(colCount)*columnWidth) / float64(colCount-1)
		for i := 0; i < colCount; i++ {
			colX := drawX + float64(i)*(columnWidth+spacing)
			drawFilledRect(screen, pixelImg, colX, drawY, columnWidth, b.Height, color.RGBA{255, 252, 245, 230})
			drawRectOutline(screen, pixelImg, colX, drawY, columnWidth, b.Height, color.RGBA{195, 180, 150, 255}, 1.5)
		}
	} else {
		// Ground Floor Shop Signboard
		signH := 24.0
		signY := drawY + b.Height - signH - 8.0
		drawFilledRect(screen, pixelImg, drawX+6.0, signY, b.Width-12.0, signH, color.RGBA{30, 28, 25, 240})
		drawRectOutline(screen, pixelImg, drawX+6.0, signY, b.Width-12.0, signH, color.RGBA{220, 160, 40, 255}, 1.5)

		if signFace != nil {
			var textOp text.DrawOptions
			textOp.GeoM.Translate(drawX+b.Width/2, signY+signH/2)
			textOp.ColorScale.Scale(1.0, 0.95, 0.3, 1.0)
			textOp.PrimaryAlign = text.AlignCenter
			textOp.SecondaryAlign = text.AlignCenter
			text.Draw(screen, b.SignTitle, signFace, &textOp)
		}

		// Roof Solar Water Heater (Ηλιακός Θερμοσίφωνας)
		heaterX := drawX + b.Width*0.45
		heaterY := drawY - 18.0
		// Silver Tank Cylinder
		drawFilledRect(screen, pixelImg, heaterX, heaterY, 26.0, 12.0, color.RGBA{200, 205, 215, 255})
		drawRectOutline(screen, pixelImg, heaterX, heaterY, 26.0, 12.0, color.RGBA{90, 95, 105, 255}, 1.0)
		// Blue Glass Solar Panel tilted
		drawFilledRect(screen, pixelImg, heaterX-14.0, heaterY+4.0, 14.0, 14.0, color.RGBA{35, 75, 140, 255})
	}

	// 3. Window Slots & Balconies
	for f := 0; f < b.Floors; f++ {
		for c := 0; c < b.Cols; c++ {
			slot := b.Slots[f][c]
			slotX := slot.X + b.ShakeOffsetX
			slotY := slot.Y + b.SinkY + b.ShakeOffsetY

			// Window Interior
			var winCol color.RGBA
			switch slot.State {
			case SlotIntact:
				winCol = color.RGBA{18, 24, 38, 255} // Dark room glass
			case SlotDamaged:
				winCol = color.RGBA{35, 30, 25, 255}
			case SlotBroken:
				winCol = color.RGBA{65, 20, 15, 255}
			case SlotEmptyHole:
				winCol = color.RGBA{12, 10, 8, 255} // Complete dark gaping hole
			}

			drawFilledRect(screen, pixelImg, slotX, slotY, slot.W, slot.H, winCol)
			drawRectOutline(screen, pixelImg, slotX, slotY, slot.W, slot.H, color.RGBA{40, 35, 30, 255}, 1.0)

			// Balcony Railing (Iconic Athenian Polykatoikia metal bars)
			if slot.HasBalcony && slot.State != SlotEmptyHole {
				railH := 14.0
				railY := slotY + slot.H - railH
				drawFilledRect(screen, pixelImg, slotX-2.0, railY, slot.W+4.0, railH, color.RGBA{65, 70, 75, 200})
				drawRectOutline(screen, pixelImg, slotX-2.0, railY, slot.W+4.0, railH, color.RGBA{30, 32, 35, 255}, 1.0)
			}

			// Orange / Blue Awning (Τέντα)
			if slot.HasAwning && slot.State != SlotEmptyHole {
				awnH := 10.0
				awnY := slotY - awnH + 2.0
				drawFilledRect(screen, pixelImg, slotX-4.0, awnY, slot.W+8.0, awnH, slot.AwningColor)
				drawRectOutline(screen, pixelImg, slotX-4.0, awnY, slot.W+8.0, awnH, color.RGBA{20, 20, 20, 255}, 1.0)
			}

			// Draw Pickups (Gyros or Freddo) if room is exposed
			if slot.State == SlotBroken || slot.State == SlotDamaged {
				if slot.ItemGyro {
					// Pita Gyros pickup icon
					drawFilledRect(screen, pixelImg, slotX+slot.W/2-6.0, slotY+slot.H/2-4.0, 12.0, 8.0, color.RGBA{225, 140, 30, 255})
				} else if slot.ItemCoffee {
					// Freddo Espresso cup pickup icon
					drawFilledRect(screen, pixelImg, slotX+slot.W/2-4.0, slotY+slot.H/2-6.0, 8.0, 12.0, color.RGBA{110, 55, 20, 255})
				}
			}
		}
	}
}

// Helper drawing utilities for crisp pixel geometry with zero dynamic allocations
func drawFilledRect(screen, pixelImg *ebiten.Image, x, y, w, h float64, col color.RGBA) {
	var op ebiten.DrawImageOptions
	op.GeoM.Scale(w, h)
	op.GeoM.Translate(x, y)
	op.ColorScale.Scale(
		float32(col.R)/255.0,
		float32(col.G)/255.0,
		float32(col.B)/255.0,
		float32(col.A)/255.0,
	)
	screen.DrawImage(pixelImg, &op)
}

func drawRectOutline(screen, pixelImg *ebiten.Image, x, y, w, h float64, col color.RGBA, thickness float64) {
	// Top, Bottom, Left, Right
	drawFilledRect(screen, pixelImg, x, y, w, thickness, col)
	drawFilledRect(screen, pixelImg, x, y+h-thickness, w, thickness, col)
	drawFilledRect(screen, pixelImg, x, y, thickness, h, col)
	drawFilledRect(screen, pixelImg, x+w-thickness, y, thickness, h, col)
}
