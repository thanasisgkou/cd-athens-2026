package game

import (
	"image/color"
	"math"
	"math/rand"
	"time"

	"athens-rampage/internal/sprite"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// CitizenState tracks the life and terror state of Athenian victims.
type CitizenState int

const (
	CitizenInWindow CitizenState = iota
	CitizenPanicking
	CitizenFleeing
	CitizenGrabbed
	CitizenEaten
	CitizenFallen
)

// Citizen represents an Athenian resident or tourist.
type Citizen struct {
	Dialog       VictimDialog
	State        CitizenState
	X, Y         float64
	Vx, Vy       float64
	FacingLeft   bool
	AnimCtrl     *sprite.AnimationController
	Building     *Building
	Floor        int
	Col          int
	SpeechText   string
	SpeechTimer  float64
	GrabTimer    float64
	EatenTimer   float64
	PanicCooldown float64
}

// NewCitizen creates a resident from stage victim dialog specs.
func NewCitizen(dialog VictimDialog, cs *sprite.CompiledSprite, b *Building, f, c int) *Citizen {
	var ctrl *sprite.AnimationController
	if cs != nil {
		c, err := cs.ToEbitenController()
		if err == nil {
			ctrl = c
			_ = ctrl.Play("idle")
		}
	}

	slot := b.Slots[f][c]
	x := slot.X + slot.W/2 - 12.0
	y := slot.Y + slot.H - 24.0

	return &Citizen{
		Dialog:     dialog,
		State:      CitizenInWindow,
		X:          x,
		Y:          y,
		AnimCtrl:   ctrl,
		Building:   b,
		Floor:      f,
		Col:        c,
		FacingLeft: rand.Float64() < 0.5,
	}
}

// TriggerPanic makes the citizen scream their authentic quote and play audio.
func (c *Citizen) TriggerPanic(audio *AudioSystem) {
	if c.SpeechTimer > 0 {
		return
	}
	c.SpeechText = c.Dialog.PanicQuote
	c.SpeechTimer = 2.4
	c.State = CitizenPanicking

	if c.AnimCtrl != nil {
		_ = c.AnimCtrl.Play("panic")
	}

	// Browser speech synthesis for spoken Greek/English quip
	SpeakLine(c.Dialog.PanicQuote, c.Dialog.VoiceLang, 1.25, 1.1)

	// Play scream or yiayia shock
	if c.Dialog.Archetype == "Γιαγιά στο Μπαλκόνι" {
		audio.PlaySFX("yiayia_shock", 800*time.Millisecond)
	} else {
		audio.PlaySFX("victim_scream", 600*time.Millisecond)
	}
}

// Update advances citizen physics, panic runs, and speech bubble timers.
func (c *Citizen) Update(dt float64, titanX, titanY float64, audio *AudioSystem) {
	if c.SpeechTimer > 0 {
		c.SpeechTimer -= dt
	}

	if c.AnimCtrl != nil {
		c.AnimCtrl.Update(time.Duration(dt * float64(time.Second)))
	}

	switch c.State {
	case CitizenInWindow:
		// Check proximity to titan
		distX := math.Abs(c.X - titanX)
		distY := math.Abs(c.Y - titanY)
		if distX < 90.0 && distY < 90.0 {
			c.TriggerPanic(audio)
		}

	case CitizenPanicking:
		// Jitter around window / balcony
		c.X += (rand.Float64()*4.0 - 2.0)
		distX := math.Abs(c.X - titanX)
		distY := math.Abs(c.Y - titanY)
		if distX < 120.0 && distY < 120.0 && rand.Float64() < 0.02 {
			c.TriggerPanic(audio)
		}

	case CitizenFleeing:
		// Runs along street sidewalk (groundY = 636)
		dir := 1.0
		if c.FacingLeft {
			dir = -1.0
		}
		c.X += dir * 75.0 * dt
		c.Y = 636.0

		if c.X < -50 || c.X > 1330 {
			c.State = CitizenEaten // despawned offscreen
		}

	case CitizenFallen:
		// Falling from a collapsing building
		c.Vy += 600.0 * dt
		c.Y += c.Vy * dt
		if c.Y >= 636.0 {
			c.Y = 636.0
			c.Vy = 0
			c.State = CitizenFleeing
		}
	}
}

// Draw renders citizen sprite and speech bubble if active.
func (c *Citizen) Draw(screen *ebiten.Image, pixelImg *ebiten.Image, bubbleFace text.Face) {
	if c.State == CitizenEaten {
		return
	}

	if c.AnimCtrl != nil {
		c.AnimCtrl.Draw(screen, c.X, c.Y, c.FacingLeft)
	} else {
		// Fallback simple pixel citizen
		drawFilledRect(screen, pixelImg, c.X+6, c.Y+4, 12, 16, c.Dialog.ShirtColor)
		drawFilledRect(screen, pixelImg, c.X+8, c.Y-2, 8, 8, color.RGBA{240, 200, 160, 255})
	}

	// Speech Bubble
	if c.SpeechTimer > 0 && bubbleFace != nil && len(c.SpeechText) > 0 {
		drawSpeechBubble(screen, pixelImg, bubbleFace, c.SpeechText, c.X+12.0, c.Y-18.0)
	}
}

// PoliceChopper ("ΕΛ.ΑΣ.") flies through the air, searchlights, and shoots at the Titan.
type PoliceChopper struct {
	Active     bool
	X, Y       float64
	Vx, Vy     float64
	Health     float64
	MaxHealth  float64
	AnimCtrl   *sprite.AnimationController
	FacingLeft bool
	ShootTimer float64
	W, H       float64
}

// NewPoliceChopper spawns an Athenian police helicopter.
func NewPoliceChopper(cs *sprite.CompiledSprite, startX, startY float64, facingLeft bool) *PoliceChopper {
	var ctrl *sprite.AnimationController
	if cs != nil {
		c, err := cs.ToEbitenController()
		if err == nil {
			ctrl = c
			_ = ctrl.Play("fly")
		}
	}

	return &PoliceChopper{
		Active:     true,
		X:          startX,
		Y:          startY,
		Health:     120.0,
		MaxHealth:  120.0,
		AnimCtrl:   ctrl,
		FacingLeft: facingLeft,
		W:          48.0,
		H:          32.0,
		ShootTimer: 1.5,
	}
}

// Update advances chopper patrol, aiming, and firing at the Titan.
func (pc *PoliceChopper) Update(dt float64, titanX, titanY float64, bullets *[]Bullet, ps *ParticleSystem) {
	if !pc.Active {
		return
	}

	if pc.AnimCtrl != nil {
		pc.AnimCtrl.Update(time.Duration(dt * float64(time.Second)))
	}

	// Hover & follow titan smoothly
	targetX := titanX + 180.0
	if pc.FacingLeft {
		targetX = titanX + 180.0
	} else {
		targetX = titanX - 180.0
	}
	targetY := math.Max(80.0, titanY-60.0)

	pc.X += (targetX - pc.X) * 1.5 * dt
	pc.Y += (targetY - pc.Y) * 1.2 * dt

	pc.FacingLeft = (pc.X > titanX)

	// Fire machinegun bursts
	pc.ShootTimer -= dt
	if pc.ShootTimer <= 0 {
		pc.ShootTimer = 2.2 + rand.Float64()*1.2

		// Spawn 3 tracer bullets towards Titan
		dx := (titanX + 32.0) - pc.X
		dy := (titanY + 32.0) - pc.Y
		dist := math.Hypot(dx, dy)
		if dist > 0 {
			vx := (dx / dist) * 380.0
			vy := (dy / dist) * 380.0

			*bullets = append(*bullets, Bullet{
				Active: true,
				X:      pc.X + 24.0,
				Y:      pc.Y + 16.0,
				Vx:     vx,
				Vy:     vy,
				Damage: 18.0,
			})
		}
	}
}

// TakeDamage damages the chopper, exploding it if health reaches zero.
func (pc *PoliceChopper) TakeDamage(dmg float64, ps *ParticleSystem, audio *AudioSystem) bool {
	pc.Health -= dmg
	ps.SpawnBuildingDebris(pc.X+pc.W/2, pc.Y+pc.H/2, 8, color.RGBA{45, 85, 185, 255}) // Police blue shards

	if pc.Health <= 0 {
		pc.Active = false
		ps.SpawnBuildingDebris(pc.X+pc.W/2, pc.Y+pc.H/2, 28, color.RGBA{240, 160, 40, 255})
		ps.SpawnWrathSparks(pc.X+pc.W/2, pc.Y+pc.H/2, color.RGBA{255, 100, 20, 255})
		audio.PlaySFX("building_collapse", 0)
		return true
	}
	return false
}

// Draw renders the police chopper and rotor blades.
func (pc *PoliceChopper) Draw(screen *ebiten.Image, pixelImg *ebiten.Image) {
	if !pc.Active {
		return
	}

	if pc.AnimCtrl != nil {
		var op ebiten.DrawImageOptions
		if pc.FacingLeft {
			op.GeoM.Scale(-1.2, 1.2)
			op.GeoM.Translate(pc.X+pc.W, pc.Y)
		} else {
			op.GeoM.Scale(1.2, 1.2)
			op.GeoM.Translate(pc.X, pc.Y)
		}
		frame := pc.AnimCtrl.Sheet.Frame(pc.AnimCtrl.CurrentFrame)
		if frame != nil {
			screen.DrawImage(frame, &op)
		}
	} else {
		// Fallback blue police chopper box
		drawFilledRect(screen, pixelImg, pc.X, pc.Y, pc.W, pc.H, color.RGBA{30, 70, 160, 255})
	}
}

// Bullet fired by police choppers at the Titan.
type Bullet struct {
	Active bool
	X, Y   float64
	Vx, Vy float64
	Damage float64
}

// Update advances bullet position.
func (b *Bullet) Update(dt float64) {
	b.X += b.Vx * dt
	b.Y += b.Vy * dt
	if b.X < -50 || b.X > 1330 || b.Y < -50 || b.Y > 750 {
		b.Active = false
	}
}

// StreetVehicle represents a Yellow Athens Taxi or Delivery Scooter driving along the road.
type StreetVehicle struct {
	Active     bool
	IsTaxi     bool // true = Yellow Mercedes Taxi, false = Delivery Scooter (παπάκι)
	X, Y       float64
	Vx         float64
	W, H       float64
	FacingLeft bool
	Crushed    bool
}

// NewStreetVehicle spawns street traffic at sidewalk level.
func NewStreetVehicle(isTaxi bool, startX float64, facingLeft bool) *StreetVehicle {
	vx := 90.0
	if isTaxi {
		vx = 135.0
	}
	if facingLeft {
		vx = -vx
	}

	w, h := 22.0, 18.0
	if isTaxi {
		w, h = 48.0, 22.0
	}

	return &StreetVehicle{
		Active:     true,
		IsTaxi:     isTaxi,
		X:          startX,
		Y:          640.0 - h,
		Vx:         vx,
		W:          w,
		H:          h,
		FacingLeft: facingLeft,
	}
}

// Update advances street vehicles.
func (v *StreetVehicle) Update(dt float64) {
	if !v.Active || v.Crushed {
		return
	}
	v.X += v.Vx * dt
	if v.X < -100 || v.X > 1380 {
		v.Active = false
	}
}

// Crush smashes the vehicle into scrap metal when stomped or punched by Titan.
func (v *StreetVehicle) Crush(ps *ParticleSystem, audio *AudioSystem) {
	if v.Crushed {
		return
	}
	v.Crushed = true
	v.H = 6.0
	v.Y = 640.0 - 6.0

	metalCol := color.RGBA{245, 215, 30, 255} // Taxi yellow
	if !v.IsTaxi {
		metalCol = color.RGBA{240, 90, 30, 255} // Delivery orange
	}
	ps.SpawnBuildingDebris(v.X+v.W/2, v.Y, 18, metalCol)
	audio.PlaySFX("titan_punch", 0)
}

// Draw renders street vehicles.
func (v *StreetVehicle) Draw(screen, pixelImg *ebiten.Image) {
	if !v.Active {
		return
	}

	if v.IsTaxi {
		// Yellow Athens Taxi
		bodyCol := color.RGBA{255, 215, 20, 255}
		if v.Crushed {
			bodyCol = color.RGBA{180, 150, 15, 255}
		}
		drawFilledRect(screen, pixelImg, v.X, v.Y, v.W, v.H, bodyCol)
		drawRectOutline(screen, pixelImg, v.X, v.Y, v.W, v.H, color.RGBA{35, 30, 20, 255}, 1.0)
		if !v.Crushed {
			// Taxi roof sign
			drawFilledRect(screen, pixelImg, v.X+v.W/2-7.0, v.Y-5.0, 14.0, 5.0, color.RGBA{255, 255, 255, 255})
		}
	} else {
		// Delivery Scooter (Παπάκι με κουτί delivery)
		drawFilledRect(screen, pixelImg, v.X, v.Y, v.W, v.H, color.RGBA{240, 95, 30, 255})
		if !v.Crushed {
			// Thermal box on back
			drawFilledRect(screen, pixelImg, v.X+2.0, v.Y-6.0, 8.0, 8.0, color.RGBA{220, 30, 40, 255})
		}
	}
}

// MechaMinotaur is the final Boss of Stage 5: The Tax Authority (ΑΑΔΕ) Robot Minotaur!
type MechaMinotaur struct {
	Active     bool
	X, Y       float64
	Vx, Vy     float64
	Health     float64
	MaxHealth  float64
	W, H       float64
	FacingLeft bool
	State      string // "idle", "charge", "laser", "stomp"
	Timer      float64
	LaserTimer float64
	LaserActive bool
	Enraged    bool
}

// NewMechaMinotaur spawns the final boss on the Acropolis.
func NewMechaMinotaur(x, y float64) *MechaMinotaur {
	return &MechaMinotaur{
		Active:     true,
		X:          x,
		Y:          y,
		Health:     2500.0,
		MaxHealth:  2500.0,
		W:          96.0,
		H:          96.0,
		FacingLeft: true,
		State:      "idle",
		Timer:      2.0,
	}
}

// Update advances the boss AI, VAT laser beams, and seismic stomps.
func (boss *MechaMinotaur) Update(dt float64, titanX, titanY float64, ps *ParticleSystem, audio *AudioSystem) {
	if !boss.Active {
		return
	}

	boss.Timer -= dt
	boss.FacingLeft = (boss.X > titanX)

	// Enrage at 40% health
	if boss.Health <= boss.MaxHealth*0.40 && !boss.Enraged {
		boss.Enraged = true
		ps.SpawnFloatingText("ΕΝΕΡΓΟΠΟΙΗΣΗ ΕΛΕΓΧΟΥ ΑΑΔΕ!", boss.X+boss.W/2, boss.Y-20.0, color.RGBA{255, 30, 40, 255}, 1.2)
		audio.PlaySFX("titan_wrath", 0)
	}

	if boss.Timer <= 0 {
		actions := []string{"charge", "laser", "stomp"}
		boss.State = actions[rand.Intn(len(actions))]
		boss.Timer = 3.0
		if boss.Enraged {
			boss.Timer = 1.8
		}

		if boss.State == "laser" {
			boss.LaserActive = true
			boss.LaserTimer = 1.2
			ps.SpawnFloatingText("ΦΠΑ 24% LASER!", boss.X+boss.W/2, boss.Y-20.0, color.RGBA{255, 80, 80, 255}, 1.1)
		} else if boss.State == "stomp" {
			audio.PlaySFX("building_collapse", 0)
			ps.SpawnBuildingDebris(boss.X+boss.W/2, boss.Y+boss.H, 20, color.RGBA{240, 230, 210, 255})
		}
	}

	if boss.State == "charge" {
		dir := -1.0
		if !boss.FacingLeft {
			dir = 1.0
		}
		speed := 220.0
		if boss.Enraged {
			speed = 320.0
		}
		boss.X += dir * speed * dt
	}

	if boss.LaserActive {
		boss.LaserTimer -= dt
		if boss.LaserTimer <= 0 {
			boss.LaserActive = false
		}
	}
}

// Draw renders the Mecha-Minotaur and its laser beam.
func (boss *MechaMinotaur) Draw(screen, pixelImg *ebiten.Image) {
	if !boss.Active {
		return
	}

	bodyCol := color.RGBA{140, 145, 155, 255}
	if boss.Enraged {
		bodyCol = color.RGBA{200, 60, 60, 255}
	}

	// Mechanical Minotaur Torso & Head
	drawFilledRect(screen, pixelImg, boss.X, boss.Y, boss.W, boss.H, bodyCol)
	drawRectOutline(screen, pixelImg, boss.X, boss.Y, boss.W, boss.H, color.RGBA{25, 25, 30, 255}, 2.5)

	// Glowing Red Eyes / Sensor Bar
	drawFilledRect(screen, pixelImg, boss.X+20.0, boss.Y+18.0, boss.W-40.0, 10.0, color.RGBA{255, 20, 30, 255})

	// Bronze Horns (Mechanical)
	drawFilledRect(screen, pixelImg, boss.X-12.0, boss.Y-18.0, 16.0, 24.0, color.RGBA{220, 160, 40, 255})
	drawFilledRect(screen, pixelImg, boss.X+boss.W-4.0, boss.Y-18.0, 16.0, 24.0, color.RGBA{220, 160, 40, 255})

	// VAT Laser Beam
	if boss.LaserActive {
		beamW := 480.0
		beamX := boss.X - beamW
		if !boss.FacingLeft {
			beamX = boss.X + boss.W
		}
		drawFilledRect(screen, pixelImg, beamX, boss.Y+22.0, beamW, 14.0, color.RGBA{255, 30, 40, 220})
	}
}

// Helper: draws a speech bubble with dark background and crisp outline
func drawSpeechBubble(screen, pixelImg *ebiten.Image, face text.Face, txt string, x, y float64) {
	bubbleW := float64(len([]rune(txt)))*11.0 + 18.0
	bubbleH := 24.0
	bX := x - bubbleW/2
	bY := y - bubbleH

	drawFilledRect(screen, pixelImg, bX, bY, bubbleW, bubbleH, color.RGBA{15, 18, 25, 235})
	drawRectOutline(screen, pixelImg, bX, bY, bubbleW, bubbleH, color.RGBA{240, 200, 60, 255}, 1.5)

	var op text.DrawOptions
	op.GeoM.Translate(x, bY+bubbleH/2)
	op.ColorScale.Scale(1.0, 1.0, 1.0, 1.0)
	op.PrimaryAlign = text.AlignCenter
	op.SecondaryAlign = text.AlignCenter
	text.Draw(screen, txt, face, &op)
}
