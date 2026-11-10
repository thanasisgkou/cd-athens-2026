package game

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// ParticleKind defines the visual and physical behavior of a particle.
type ParticleKind int

const (
	PartBloodDrop ParticleKind = iota
	PartBloodDecal
	PartBone
	PartSandal
	PartFreddoCup
	PartCoffeeDrop
	PartDebris
	PartDust
	PartFireSpark
	PartLightning
	PartSmoke
)

// Particle represents a single visual effect item.
type Particle struct {
	Active   bool
	Kind     ParticleKind
	X, Y     float64
	Vx, Vy   float64
	Rot      float64
	RotSpeed float64
	Scale    float64
	Color    color.RGBA
	Life     float64
	MaxLife  float64
	Gravity  float64
	GroundY  float64
	Bounce   float64
	Width    float64
	Height   float64
}

// FloatingText represents an arcade popup text (e.g. "+250 HP", "CRUNCH!", "ΩΧ Ο ΚΑΦΕΣ!").
type FloatingText struct {
	Active  bool
	Text    string
	X, Y    float64
	Vy      float64
	Life    float64
	MaxLife float64
	Color   color.RGBA
	Scale   float64
}

const (
	MaxParticles = 1024
	MaxFloats    = 64
)

// ParticleSystem handles all gore, debris, and floating text effects.
type ParticleSystem struct {
	particles  [MaxParticles]Particle
	floatTexts [MaxFloats]FloatingText
	pixelImg   *ebiten.Image
	drawOp     ebiten.DrawImageOptions
}

// NewParticleSystem creates and pre-allocates the particle pool.
func NewParticleSystem() *ParticleSystem {
	ps := &ParticleSystem{
		pixelImg: ebiten.NewImage(1, 1),
	}
	ps.pixelImg.Fill(color.White)
	return ps
}

// SpawnBloodFountain explodes a burst of crimson and bright red blood droplets.
func (ps *ParticleSystem) SpawnBloodFountain(x, y float64, count int, dirX float64) {
	reds := []color.RGBA{
		{210, 10, 20, 255},
		{160, 5, 10, 255},
		{240, 30, 40, 255},
		{120, 0, 5, 255},
	}

	for i := 0; i < count; i++ {
		slot := ps.allocParticle()
		if slot == nil {
			break
		}

		angle := (rand.Float64()*math.Pi - math.Pi/2) * 0.8
		speed := 180.0 + rand.Float64()*260.0
		vx := math.Cos(angle)*speed*dirX + (rand.Float64()*80.0 - 40.0)
		vy := -math.Abs(math.Sin(angle)*speed) - (rand.Float64() * 120.0)

		slot.Active = true
		slot.Kind = PartBloodDrop
		slot.X = x + (rand.Float64()*12.0 - 6.0)
		slot.Y = y + (rand.Float64()*12.0 - 6.0)
		slot.Vx = vx
		slot.Vy = vy
		slot.Color = reds[rand.Intn(len(reds))]
		slot.Life = 0.8 + rand.Float64()*0.6
		slot.MaxLife = slot.Life
		slot.Gravity = 580.0
		slot.GroundY = 660.0
		slot.Bounce = 0.25 + rand.Float64()*0.2
		slot.Width = 3.0 + rand.Float64()*3.0
		slot.Height = slot.Width
	}
}

// SpawnGoreGibs spawns flying bones, Greek sandals, and Freddo cups when a victim is devoured.
func (ps *ParticleSystem) SpawnGoreGibs(x, y float64) {
	// Flying Bone/Skull
	for i := 0; i < 3; i++ {
		slot := ps.allocParticle()
		if slot == nil {
			break
		}
		slot.Active = true
		slot.Kind = PartBone
		slot.X = x
		slot.Y = y
		slot.Vx = (rand.Float64()*240.0 - 120.0)
		slot.Vy = -180.0 - rand.Float64()*180.0
		slot.Rot = rand.Float64() * math.Pi * 2
		slot.RotSpeed = (rand.Float64()*12.0 - 6.0)
		slot.Color = color.RGBA{245, 245, 235, 255}
		slot.Life = 1.4
		slot.MaxLife = 1.4
		slot.Gravity = 520.0
		slot.GroundY = 660.0
		slot.Bounce = 0.35
		slot.Width = 6.0
		slot.Height = 4.0
	}

	// Flying Greek Sandal / Παντόφλα
	if slot := ps.allocParticle(); slot != nil {
		slot.Active = true
		slot.Kind = PartSandal
		slot.X = x
		slot.Y = y
		slot.Vx = (rand.Float64()*200.0 - 100.0)
		slot.Vy = -220.0 - rand.Float64()*120.0
		slot.Rot = rand.Float64() * math.Pi * 2
		slot.RotSpeed = (rand.Float64()*10.0 - 5.0)
		slot.Color = color.RGBA{180, 110, 50, 255} // Greek leather sandal
		slot.Life = 1.5
		slot.MaxLife = 1.5
		slot.Gravity = 480.0
		slot.GroundY = 660.0
		slot.Bounce = 0.4
		slot.Width = 9.0
		slot.Height = 5.0
	}

	// Flying Freddo Espresso Cup (50% chance)
	if rand.Float64() < 0.55 {
		if slot := ps.allocParticle(); slot != nil {
			slot.Active = true
			slot.Kind = PartFreddoCup
			slot.X = x
			slot.Y = y
			slot.Vx = (rand.Float64()*180.0 - 90.0)
			slot.Vy = -200.0 - rand.Float64()*150.0
			slot.Rot = rand.Float64() * math.Pi * 2
			slot.RotSpeed = (rand.Float64()*8.0 - 4.0)
			slot.Color = color.RGBA{110, 60, 25, 255} // Coffee brown
			slot.Life = 1.6
			slot.MaxLife = 1.6
			slot.Gravity = 450.0
			slot.GroundY = 660.0
			slot.Bounce = 0.3
			slot.Width = 8.0
			slot.Height = 11.0
		}
	}
}

// SpawnBuildingDebris bursts concrete chunks and dust puffs when punching structures.
func (ps *ParticleSystem) SpawnBuildingDebris(x, y float64, count int, wallCol color.RGBA) {
	for i := 0; i < count; i++ {
		slot := ps.allocParticle()
		if slot == nil {
			break
		}
		slot.Active = true
		slot.Kind = PartDebris
		slot.X = x + (rand.Float64()*16.0 - 8.0)
		slot.Y = y + (rand.Float64()*16.0 - 8.0)
		slot.Vx = (rand.Float64()*260.0 - 130.0)
		slot.Vy = -100.0 - rand.Float64()*180.0
		slot.Rot = rand.Float64() * math.Pi * 2
		slot.RotSpeed = (rand.Float64()*8.0 - 4.0)
		slot.Color = wallCol
		slot.Life = 0.7 + rand.Float64()*0.5
		slot.MaxLife = slot.Life
		slot.Gravity = 620.0
		slot.GroundY = 660.0
		slot.Bounce = 0.2
		slot.Width = 4.0 + rand.Float64()*4.0
		slot.Height = slot.Width
	}

	// Dust cloud
	for i := 0; i < 4; i++ {
		slot := ps.allocParticle()
		if slot == nil {
			break
		}
		slot.Active = true
		slot.Kind = PartDust
		slot.X = x + (rand.Float64()*20.0 - 10.0)
		slot.Y = y + (rand.Float64()*20.0 - 10.0)
		slot.Vx = (rand.Float64()*60.0 - 30.0)
		slot.Vy = -30.0 - rand.Float64()*40.0
		slot.Color = color.RGBA{180, 175, 170, 160}
		slot.Life = 0.6 + rand.Float64()*0.4
		slot.MaxLife = slot.Life
		slot.Gravity = -10.0 // Floats gently
		slot.Width = 8.0 + rand.Float64()*8.0
		slot.Height = slot.Width
	}
}

// SpawnWrathSparks spawns lightning and fire particles during Titan Wrath.
func (ps *ParticleSystem) SpawnWrathSparks(x, y float64, auraCol color.RGBA) {
	for i := 0; i < 18; i++ {
		slot := ps.allocParticle()
		if slot == nil {
			break
		}
		angle := rand.Float64() * math.Pi * 2
		speed := 120.0 + rand.Float64()*240.0
		slot.Active = true
		slot.Kind = PartLightning
		slot.X = x
		slot.Y = y
		slot.Vx = math.Cos(angle) * speed
		slot.Vy = math.Sin(angle) * speed
		slot.Color = auraCol
		slot.Life = 0.4 + rand.Float64()*0.3
		slot.MaxLife = slot.Life
		slot.Gravity = 0
		slot.Width = 3.0 + rand.Float64()*4.0
		slot.Height = slot.Width
	}
}

// SpawnFloatingText triggers an arcade text popup.
func (ps *ParticleSystem) SpawnFloatingText(txt string, x, y float64, col color.RGBA, scale float64) {
	for i := range ps.floatTexts {
		if !ps.floatTexts[i].Active {
			ps.floatTexts[i] = FloatingText{
				Active:  true,
				Text:    txt,
				X:       x,
				Y:       y,
				Vy:      -65.0,
				Life:    1.2,
				MaxLife: 1.2,
				Color:   col,
				Scale:   scale,
			}
			return
		}
	}
}

func (ps *ParticleSystem) allocParticle() *Particle {
	for i := range ps.particles {
		if !ps.particles[i].Active {
			return &ps.particles[i]
		}
	}
	return nil
}

// Update advances all particles and text floats by dt seconds.
func (ps *ParticleSystem) Update(dt float64) {
	for i := range ps.particles {
		p := &ps.particles[i]
		if !p.Active {
			continue
		}

		p.Life -= dt
		if p.Life <= 0 {
			p.Active = false
			continue
		}

		p.Vy += p.Gravity * dt
		p.X += p.Vx * dt
		p.Y += p.Vy * dt
		p.Rot += p.RotSpeed * dt

		// Ground bounce
		if p.GroundY > 0 && p.Y >= p.GroundY {
			p.Y = p.GroundY
			if math.Abs(p.Vy) > 30.0 {
				p.Vy = -p.Vy * p.Bounce
				p.Vx *= 0.65
			} else {
				p.Vy = 0
				p.Vx *= 0.4
			}
		}
	}

	for i := range ps.floatTexts {
		ft := &ps.floatTexts[i]
		if !ft.Active {
			continue
		}
		ft.Life -= dt
		if ft.Life <= 0 {
			ft.Active = false
			continue
		}
		ft.Y += ft.Vy * dt
		ft.Vy *= 0.94
	}
}

// Draw renders all active particles with zero dynamic heap allocations.
func (ps *ParticleSystem) Draw(screen *ebiten.Image) {
	for i := range ps.particles {
		p := &ps.particles[i]
		if !p.Active {
			continue
		}

		alpha := math.Min(1.0, p.Life/math.Max(0.01, p.MaxLife*0.4))
		colR := float32(p.Color.R) / 255.0
		colG := float32(p.Color.G) / 255.0
		colB := float32(p.Color.B) / 255.0
		colA := float32(p.Color.A) / 255.0 * float32(alpha)

		ps.drawOp.GeoM.Reset()
		ps.drawOp.ColorScale.Reset()

		if p.Rot != 0 {
			ps.drawOp.GeoM.Translate(-p.Width/2, -p.Height/2)
			ps.drawOp.GeoM.Rotate(p.Rot)
			ps.drawOp.GeoM.Translate(p.Width/2, p.Height/2)
		}

		ps.drawOp.GeoM.Scale(p.Width, p.Height)
		ps.drawOp.GeoM.Translate(p.X, p.Y)
		ps.drawOp.ColorScale.Scale(colR, colG, colB, colA)

		screen.DrawImage(ps.pixelImg, &ps.drawOp)
	}
}

// DrawFloatingTexts renders text popups using the provided font face.
func (ps *ParticleSystem) DrawFloatingTexts(screen *ebiten.Image, face text.Face) {
	if face == nil {
		return
	}

	for i := range ps.floatTexts {
		ft := &ps.floatTexts[i]
		if !ft.Active {
			continue
		}

		alpha := math.Min(1.0, ft.Life/math.Max(0.01, ft.MaxLife*0.3))
		colR := float32(ft.Color.R) / 255.0
		colG := float32(ft.Color.G) / 255.0
		colB := float32(ft.Color.B) / 255.0
		colA := float32(ft.Color.A) / 255.0 * float32(alpha)

		var op text.DrawOptions
		op.GeoM.Translate(ft.X, ft.Y)
		op.ColorScale.Scale(colR, colG, colB, colA)
		op.PrimaryAlign = text.AlignCenter
		op.SecondaryAlign = text.AlignCenter

		text.Draw(screen, ft.Text, face, &op)
	}
}
