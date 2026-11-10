package game

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"athens-rampage/internal/sprite"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// TitanState defines the active animation and movement mode of the Titan.
type TitanState int

const (
	StateIdle TitanState = iota
	StateWalking
	StateClimbing
	StatePunching
	StateDevouring
	StateWrath
	StateFalling
)

// Titan represents the giant mythological player character.
type Titan struct {
	Class        TitanClass
	State        TitanState
	X, Y         float64
	Vx, Vy       float64
	W, H         float64
	FacingLeft   bool
	Health       float64
	MaxHealth    float64
	Rage         float64 // 0 to 100%
	Score        int
	GoreCombo    int
	ComboTimer   float64
	AnimCtrl     *sprite.AnimationController
	ClimbingBldg *Building
	ClimbSide    int // -1 = left wall, 0 = facade, +1 = right wall
	PunchTimer   float64
	DevourTimer  float64
	WrathTimer   float64
	HeldVictim   *Citizen
	InvulnTimer  float64
}

// NewTitan initializes the player with their selected Titan class.
func NewTitan(class TitanClass, cs *sprite.CompiledSprite, startX float64) *Titan {
	var ctrl *sprite.AnimationController
	if cs != nil {
		c, err := cs.ToEbitenController()
		if err == nil {
			ctrl = c
			_ = ctrl.Play("idle")
		}
	}

	const (
		titanW  = 72.0
		titanH  = 72.0
		groundY = 640.0
	)

	return &Titan{
		Class:      class,
		State:      StateIdle,
		X:          startX,
		Y:          groundY - titanH,
		W:          titanW,
		H:          titanH,
		Health:     1000.0,
		MaxHealth:  1000.0,
		Rage:       25.0,
		AnimCtrl:   ctrl,
		FacingLeft: false,
	}
}

// Update handles player input, climbing physics, punching hitboxes, and devouring victims.
func (t *Titan) Update(dt float64, buildings []*Building, choppers []*PoliceChopper, traffic []*StreetVehicle, boss *MechaMinotaur, ps *ParticleSystem, audio *AudioSystem) {
	if t.InvulnTimer > 0 {
		t.InvulnTimer -= dt
	}

	if t.ComboTimer > 0 {
		t.ComboTimer -= dt
		if t.ComboTimer <= 0 {
			t.GoreCombo = 0
		}
	}

	if t.AnimCtrl != nil {
		t.AnimCtrl.Update(time.Duration(dt * float64(time.Second)))
	}

	// State-specific timers
	switch t.State {
	case StatePunching:
		t.PunchTimer -= dt
		if t.PunchTimer <= 0 {
			if t.ClimbingBldg != nil {
				t.State = StateClimbing
				t.playAnim("climb")
			} else {
				t.State = StateIdle
				t.playAnim("idle")
			}
		}
		return

	case StateDevouring:
		t.DevourTimer -= dt
		// Blood drips during chew
		if randFloat() < 0.35 {
			biteX := t.X + t.W*0.6
			if t.FacingLeft {
				biteX = t.X + t.W*0.4
			}
			ps.SpawnBloodFountain(biteX, t.Y+t.H*0.3, 3, 1.0)
		}

		if t.DevourTimer <= 0 {
			// Finish devouring victim
			t.finishDevour(ps, audio)
			if t.ClimbingBldg != nil {
				t.State = StateClimbing
				t.playAnim("climb")
			} else {
				t.State = StateIdle
				t.playAnim("idle")
			}
		}
		return

	case StateWrath:
		t.WrathTimer -= dt
		// Burst sparks & screen clear damage
		ps.SpawnWrathSparks(t.X+t.W/2, t.Y+t.H/2, t.Class.AuraColor)
		if t.WrathTimer <= 0 {
			if t.ClimbingBldg != nil {
				t.State = StateClimbing
				t.playAnim("climb")
			} else {
				t.State = StateIdle
				t.playAnim("idle")
			}
		}
		return
	}

	// 1. Movement Inputs
	moveX := 0.0
	moveY := 0.0

	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		moveX -= 1.0
		t.FacingLeft = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		moveX += 1.0
		t.FacingLeft = false
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		moveY -= 1.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		moveY += 1.0
	}

	// 2. Action Inputs
	isPunch := inpututil.IsKeyJustPressed(ebiten.KeyJ) || inpututil.IsKeyJustPressed(ebiten.KeyZ) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	isGrab := inpututil.IsKeyJustPressed(ebiten.KeyK) || inpututil.IsKeyJustPressed(ebiten.KeyX) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight)
	isWrath := inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyE)

	// Execute Wrath if fully charged
	if isWrath && t.Rage >= 100.0 {
		t.triggerWrath(buildings, choppers, boss, ps, audio)
		return
	}

	// Execute Grab / Eat
	if isGrab {
		if t.tryGrabCitizen(buildings, ps, audio) {
			return
		}
	}

	// Execute Punch
	if isPunch {
		t.triggerPunch(buildings, choppers, traffic, boss, ps, audio)
		return
	}

	// 3. Climbing vs Ground Physics
	groundY := 640.0 - t.H

	if t.State == StateClimbing {
		climbSpeed := 150.0 * t.Class.ClimbMul
		t.Y += moveY * climbSpeed * dt
		t.X += moveX * (climbSpeed * 0.75) * dt

		// Check if dropped below ground
		if t.Y >= groundY {
			t.Y = groundY
			t.ClimbingBldg = nil
			t.State = StateIdle
			t.playAnim("idle")
			return
		}

		// Check if climbed off top or side of building
		if t.ClimbingBldg != nil {
			if t.ClimbingBldg.Collapsed {
				// Building collapsed underneath Titan!
				t.ClimbingBldg = nil
				t.State = StateFalling
				t.Vy = 100.0
				return
			}

			// Stay bounded to building
			minX := t.ClimbingBldg.X - 16.0
			maxX := t.ClimbingBldg.X + t.ClimbingBldg.Width - t.W + 16.0
			if t.X < minX || t.X > maxX {
				// Reached edge
				t.X = math.Max(minX, math.Min(maxX, t.X))
			}

			topY := t.ClimbingBldg.Y + t.ClimbingBldg.SinkY - 10.0
			if t.Y < topY {
				t.Y = topY // Standing on roof
			}
		}

		if moveX != 0 || moveY != 0 {
			t.playAnim("climb")
		}
	} else {
		// Ground movement
		walkSpeed := 190.0 * t.Class.SpeedMul
		t.X += moveX * walkSpeed * dt

		// Clamp screen bounds
		t.X = math.Max(10.0, math.Min(1280.0-t.W-10.0, t.X))

		// Check for climbing attachment
		if moveY < 0 {
			for _, b := range buildings {
				if b.Collapsed {
					continue
				}
				if t.X+t.W > b.X && t.X < b.X+b.Width {
					// Latch onto building
					t.ClimbingBldg = b
					t.State = StateClimbing
					t.playAnim("climb")
					break
				}
			}
		}

		if t.State != StateClimbing {
			if moveX != 0 {
				t.State = StateWalking
				t.playAnim("walk")
			} else {
				t.State = StateIdle
				t.playAnim("idle")
			}
		}
	}
}

func (t *Titan) triggerPunch(buildings []*Building, choppers []*PoliceChopper, traffic []*StreetVehicle, boss *MechaMinotaur, ps *ParticleSystem, audio *AudioSystem) {
	t.State = StatePunching
	t.PunchTimer = 0.22
	t.playAnim("punch")
	audio.PlaySFX("titan_punch", 100*time.Millisecond)

	// Punch Hitbox in front of Titan
	punchW := 42.0
	punchH := 36.0
	punchX := t.X + t.W - 6.0
	if t.FacingLeft {
		punchX = t.X - punchW + 6.0
	}
	punchY := t.Y + 16.0

	punchDamage := 65.0 * t.Class.SmashMul

	// 1. Damage Buildings
	for _, b := range buildings {
		if b.Collapsed {
			continue
		}
		for f := 0; f < b.Floors; f++ {
			for c := 0; c < b.Cols; c++ {
				slot := b.Slots[f][c]
				slotX := slot.X + b.ShakeOffsetX
				slotY := slot.Y + b.SinkY + b.ShakeOffsetY

				if rectsOverlap(punchX, punchY, punchW, punchH, slotX, slotY, slot.W, slot.H) {
					dmg, destroyed := b.DamageSlot(f, c, punchDamage, ps)
					if dmg > 0 {
						t.Score += int(dmg)
						t.Rage = math.Min(100.0, t.Rage+2.5)

						// Collect gyro / coffee pickups inside room
						if slot.ItemGyro {
							slot.ItemGyro = false
							t.Health = math.Min(t.MaxHealth, t.Health+120.0)
							ps.SpawnFloatingText("ΠΙΤΟΓΥΡΟ! +120 HP", t.X+t.W/2, t.Y-10.0, color.RGBA{240, 180, 40, 255}, 1.1)
							audio.PlaySFX("pickup_gyros", 0)
						} else if slot.ItemCoffee {
							slot.ItemCoffee = false
							t.Rage = math.Min(100.0, t.Rage+20.0)
							ps.SpawnFloatingText("FREDDO ESPRESSO! +20 RAGE", t.X+t.W/2, t.Y-10.0, color.RGBA{120, 220, 255, 255}, 1.1)
							audio.PlaySFX("pickup_gyros", 0)
						}

						if destroyed {
							t.Score += 250
							t.Rage = math.Min(100.0, t.Rage+5.0)
							ps.SpawnFloatingText("SMASHED!", slotX+slot.W/2, slotY, color.RGBA{255, 220, 80, 255}, 1.0)
						}
					}
				}
			}
		}
	}

	// 2. Swat Police Choppers
	for _, pc := range choppers {
		if !pc.Active {
			continue
		}
		if rectsOverlap(punchX, punchY, punchW, punchH, pc.X, pc.Y, pc.W, pc.H) {
			if pc.TakeDamage(punchDamage*1.4, ps, audio) {
				t.Score += 800
				t.Rage = math.Min(100.0, t.Rage+15.0)
				ps.SpawnFloatingText("ΕΛ.ΑΣ. DOWN! +800", pc.X+pc.W/2, pc.Y, color.RGBA{255, 80, 80, 255}, 1.2)
			}
		}
	}

	// 3. Crush Street Vehicles
	for _, v := range traffic {
		if !v.Active || v.Crushed {
			continue
		}
		if rectsOverlap(punchX, punchY, punchW, punchH, v.X, v.Y, v.W, v.H) {
			v.Crush(ps, audio)
			t.Score += 350
			t.Rage = math.Min(100.0, t.Rage+6.0)
			if v.IsTaxi {
				ps.SpawnFloatingText("ΤΑΞΙ SMASH! +350", v.X+v.W/2, v.Y-15.0, color.RGBA{255, 220, 30, 255}, 1.0)
			} else {
				ps.SpawnFloatingText("ΠΑΠΑΚΙ DELIVERY! +350", v.X+v.W/2, v.Y-15.0, color.RGBA{255, 140, 30, 255}, 1.0)
			}
		}
	}

	// 4. Punch Boss
	if boss != nil && boss.Active {
		if rectsOverlap(punchX, punchY, punchW, punchH, boss.X, boss.Y, boss.W, boss.H) {
			boss.Health -= punchDamage
			ps.SpawnBuildingDebris(boss.X+boss.W/2, boss.Y+boss.H/2, 10, color.RGBA{220, 50, 50, 255})
			t.Score += int(punchDamage * 2.0)
			t.Rage = math.Min(100.0, t.Rage+3.0)
			ps.SpawnFloatingText("BOSS HIT!", boss.X+boss.W/2, boss.Y-10.0, color.RGBA{255, 100, 100, 255}, 1.2)
		}
	}
}

func (t *Titan) tryGrabCitizen(buildings []*Building, ps *ParticleSystem, audio *AudioSystem) bool {
	grabRadius := 56.0

	for _, b := range buildings {
		if b.Collapsed {
			continue
		}
		for f := 0; f < b.Floors; f++ {
			for c := 0; c < b.Cols; c++ {
				slot := b.Slots[f][c]
				if slot.Victim == nil || slot.Victim.State == CitizenEaten {
					continue
				}

				vic := slot.Victim
				dist := math.Hypot((t.X+t.W/2)-(vic.X+12.0), (t.Y+t.H/2)-(vic.Y+12.0))
				if dist <= grabRadius {
					// Grab victim!
					t.HeldVictim = vic
					vic.State = CitizenGrabbed
					t.State = StateDevouring
					t.DevourTimer = 0.65
					t.playAnim("eat")

					// Speech quip on grab
					vic.SpeechText = vic.Dialog.EatenQuote
					vic.SpeechTimer = 2.0
					SpeakLine(vic.Dialog.EatenQuote, vic.Dialog.VoiceLang, 1.4, 1.2)

					audio.PlaySFX("victim_scream", 0)
					return true
				}
			}
		}
	}
	return false
}

func (t *Titan) finishDevour(ps *ParticleSystem, audio *AudioSystem) {
	if t.HeldVictim == nil {
		return
	}

	vic := t.HeldVictim
	vic.State = CitizenEaten

	// Devour Gore FX
	biteX := t.X + t.W*0.6
	if t.FacingLeft {
		biteX = t.X + t.W*0.4
	}
	biteY := t.Y + t.H*0.3

	ps.SpawnBloodFountain(biteX, biteY, 32, 1.2)
	ps.SpawnGoreGibs(biteX, biteY)
	audio.PlaySFX("devour_gore", 0)

	// Heal Titan based on class multiplier!
	healAmount := 180.0 * t.Class.GoreHealMul
	t.Health = math.Min(t.MaxHealth, t.Health+healAmount)

	// Combo Multiplier
	t.GoreCombo++
	t.ComboTimer = 4.2
	comboMultiplier := t.GoreCombo
	if comboMultiplier > 5 {
		comboMultiplier = 5
	}

	t.Score += 600 * comboMultiplier
	t.Rage = math.Min(100.0, t.Rage+12.0)

	popupTxt := "DELICIOUS! +180 HP"
	if t.GoreCombo > 1 {
		popupTxt = fmt.Sprintf("GORE COMBO x%d! +HP", comboMultiplier)
	}
	ps.SpawnFloatingText(popupTxt, t.X+t.W/2, t.Y-18.0, color.RGBA{255, 30, 50, 255}, 1.25)

	t.HeldVictim = nil
}

func (t *Titan) triggerWrath(buildings []*Building, choppers []*PoliceChopper, boss *MechaMinotaur, ps *ParticleSystem, audio *AudioSystem) {
	t.Rage = 0
	t.State = StateWrath
	t.WrathTimer = 0.9
	t.playAnim("roar")

	audio.PlaySFX("titan_wrath", 0)
	SpeakLine(t.Class.VoiceIntro, "el-GR", 0.9, 1.0)

	ps.SpawnFloatingText(t.Class.SpecialName, t.X+t.W/2, t.Y-30.0, t.Class.AuraColor, 1.4)

	// Massive shockwave damage to all choppers and nearby buildings
	for _, pc := range choppers {
		if pc.Active {
			_ = pc.TakeDamage(300.0, ps, audio)
		}
	}

	for _, b := range buildings {
		if b.Collapsed {
			continue
		}
		for f := 0; f < b.Floors; f++ {
			for c := 0; c < b.Cols; c++ {
				_, _ = b.DamageSlot(f, c, 180.0, ps)
			}
		}
	}

	if boss != nil && boss.Active {
		boss.Health -= 400.0
		ps.SpawnBuildingDebris(boss.X+boss.W/2, boss.Y+boss.H/2, 25, color.RGBA{255, 50, 50, 255})
	}
}

func (t *Titan) playAnim(tag string) {
	if t.AnimCtrl != nil {
		_ = t.AnimCtrl.Play(tag)
	}
}

// Draw renders the Titan with class color scale and aura glows.
func (t *Titan) Draw(screen *ebiten.Image, pixelImg *ebiten.Image) {
	if t.AnimCtrl != nil && t.AnimCtrl.Sheet != nil {
		frame := t.AnimCtrl.Sheet.Frame(t.AnimCtrl.CurrentFrame)
		if frame != nil {
			var op ebiten.DrawImageOptions

			scale := 1.5 // 48px -> 72px giant Titan
			if t.FacingLeft {
				op.GeoM.Scale(-scale, scale)
				op.GeoM.Translate(t.X+t.W, t.Y)
			} else {
				op.GeoM.Scale(scale, scale)
				op.GeoM.Translate(t.X, t.Y)
			}

			// Apply Titan class tint
			op.ColorScale.Scale(
				float32(t.Class.TintColor.R)/255.0,
				float32(t.Class.TintColor.G)/255.0,
				float32(t.Class.TintColor.B)/255.0,
				1.0,
			)

			// Wrath aura glow
			if t.State == StateWrath {
				op.ColorScale.Scale(1.4, 1.2, 1.2, 1.0)
			}

			screen.DrawImage(frame, &op)
		}
	} else {
		// Fallback Titan rect
		drawFilledRect(screen, pixelImg, t.X, t.Y, t.W, t.H, t.Class.TintColor)
	}
}

func rectsOverlap(x1, y1, w1, h1, x2, y2, w2, h2 float64) bool {
	return x1 < x2+w2 && x1+w1 > x2 && y1 < y2+h2 && y1+h1 > y2
}

func randFloat() float64 {
	return float64(time.Now().UnixNano()%1000) / 1000.0
}
