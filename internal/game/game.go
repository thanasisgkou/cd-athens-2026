package game

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"math"
	"math/rand"
	"time"

	"athens-rampage/assets"
	"athens-rampage/internal/sprite"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// GameScene enumerates the high-level game states.
type GameScene int

const (
	SceneTitle GameScene = iota
	SceneTitanSelect
	SceneStageIntro
	ScenePlay
	SceneStageClear
	SceneVictory
	SceneGameOver
)

// AthensGame is the top-level Ebitengine Game implementation.
type AthensGame struct {
	Scene         GameScene
	CurrentStage  int // 0 to 4 (corresponding to AthensStages[0..4])
	SelectedTitan int // 0, 1, 2 (Kronos, Typhon, Polyphemus)

	// Systems
	Audio     *AudioSystem
	Particles *ParticleSystem

	// Fonts & Text Faces
	fontSource *text.GoTextFaceSource
	faceTitle  text.Face
	faceSub    text.Face
	faceHUD    text.Face
	faceSign   text.Face
	faceBubble text.Face
	faceStats  text.Face

	// Pre-loaded Artwork & Sprites
	titleBanner   *ebiten.Image
	titanPortrait *ebiten.Image
	pixelImg      *ebiten.Image
	kronosSprite  *sprite.CompiledSprite
	victimSprite  *sprite.CompiledSprite
	chopperSprite *sprite.CompiledSprite

	// Active Stage Entities
	Titan        *Titan
	Buildings    []*Building
	Citizens     []*Citizen
	Choppers     []*PoliceChopper
	Traffic      []*StreetVehicle
	Boss         *MechaMinotaur
	Bullets      []Bullet
	ChopperTimer float64
	TrafficTimer float64

	// Stage Statistics & Tracking
	StageTimer       float64
	StageVictimsEaten int
	StageBuildingsDown int
	IntroTimer       float64
	ClearTimer       float64

	// Screen Shake / Camera Trauma
	Trauma float64

	// Last timestamp for dt calculation
	lastTime time.Time
}

// NewAthensGame creates and fully initializes Athens Rampage.
func NewAthensGame() (*AthensGame, error) {
	audioSys := GetAudioSystem()
	particleSys := NewParticleSystem()

	g := &AthensGame{
		Scene:         SceneTitle,
		Audio:         audioSys,
		Particles:     particleSys,
		lastTime:      time.Now(),
		pixelImg:      ebiten.NewImage(1, 1),
		SelectedTitan: 0,
		CurrentStage:  0,
	}
	g.pixelImg.Fill(color.White)

	// 1. Load Fonts
	fontData, err := assets.FS.ReadFile("fonts/font.ttf")
	if err != nil {
		return nil, fmt.Errorf("failed to read font.ttf: %w", err)
	}
	src, err := text.NewGoTextFaceSource(bytes.NewReader(fontData))
	if err != nil {
		return nil, fmt.Errorf("failed to parse font: %w", err)
	}
	g.fontSource = src
	g.faceTitle = &text.GoTextFace{Source: src, Size: 32}
	g.faceSub = &text.GoTextFace{Source: src, Size: 20}
	g.faceHUD = &text.GoTextFace{Source: src, Size: 16}
	g.faceSign = &text.GoTextFace{Source: src, Size: 10}
	g.faceBubble = &text.GoTextFace{Source: src, Size: 11}
	g.faceStats = &text.GoTextFace{Source: src, Size: 24}

	// 2. Load Banner & Portrait Images
	bannerData, err := assets.FS.ReadFile("images/title_banner.png")
	if err == nil {
		img, _, err := image.Decode(bytes.NewReader(bannerData))
		if err == nil {
			g.titleBanner = ebiten.NewImageFromImage(img)
		}
	}

	portraitData, err := assets.FS.ReadFile("images/kronos_portrait.png")
	if err == nil {
		img, _, err := image.Decode(bytes.NewReader(portraitData))
		if err == nil {
			g.titanPortrait = ebiten.NewImageFromImage(img)
		}
	}

	// 3. Compile Procedural Sprites
	titanData, err := assets.FS.ReadFile("sprites/kronos_titan.sprite")
	if err == nil {
		g.kronosSprite, err = sprite.CompileSpriteDSL(string(titanData))
		if err != nil {
			fmt.Printf("[Game] Titan sprite compile warning: %v\n", err)
		}
	}

	vicData, err := assets.FS.ReadFile("sprites/athens_victim.sprite")
	if err == nil {
		g.victimSprite, err = sprite.CompileSpriteDSL(string(vicData))
		if err != nil {
			fmt.Printf("[Game] Victim sprite compile warning: %v\n", err)
		}
	}

	chopData, err := assets.FS.ReadFile("sprites/police_chopper.sprite")
	if err == nil {
		g.chopperSprite, err = sprite.CompileSpriteDSL(string(chopData))
		if err != nil {
			fmt.Printf("[Game] Chopper sprite compile warning: %v\n", err)
		}
	}

	// Start looping procedural BGM
	g.Audio.PlayBGM()

	return g, nil
}

// Layout maintains the fixed 1280x720 arcade virtual resolution.
func (g *AthensGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 1280, 720
}

// Update advances the scene state machine and physics.
func (g *AthensGame) Update() error {
	now := time.Now()
	dt := now.Sub(g.lastTime).Seconds()
	if dt <= 0 || dt > 0.1 {
		dt = 1.0 / 60.0
	}
	g.lastTime = now

	// Toggle mute with 'M'
	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		g.Audio.ToggleMute()
	}

	// Toggle fullscreen with F11 or Alt+Enter
	if inpututil.IsKeyJustPressed(ebiten.KeyF11) || (ebiten.IsKeyPressed(ebiten.KeyAlt) && inpututil.IsKeyJustPressed(ebiten.KeyEnter)) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	// Camera trauma decay
	if g.Trauma > 0 {
		g.Trauma = math.Max(0, g.Trauma-dt*1.2)
	}

	switch g.Scene {
	case SceneTitle:
		g.updateTitle()
	case SceneTitanSelect:
		g.updateTitanSelect()
	case SceneStageIntro:
		g.updateStageIntro(dt)
	case ScenePlay:
		g.updatePlay(dt)
	case SceneStageClear:
		g.updateStageClear(dt)
	case SceneVictory:
		g.updateVictory()
	case SceneGameOver:
		g.updateGameOver()
	}

	return nil
}

func (g *AthensGame) updateTitle() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.Scene = SceneTitanSelect
		g.Audio.PlaySFX("pickup_gyros", 0)
	}
}

func (g *AthensGame) updateTitanSelect() {
	// Left/Right or A/D to cycle Titans
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.SelectedTitan = (g.SelectedTitan - 1 + len(TitanRoster)) % len(TitanRoster)
		g.Audio.PlaySFX("pickup_gyros", 0)
		t := TitanRoster[g.SelectedTitan]
		SpeakLine(t.VoiceIntro, "el-GR", 1.0, 1.0)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.SelectedTitan = (g.SelectedTitan + 1) % len(TitanRoster)
		g.Audio.PlaySFX("pickup_gyros", 0)
		t := TitanRoster[g.SelectedTitan]
		SpeakLine(t.VoiceIntro, "el-GR", 1.0, 1.0)
	}

	// Lock in selection
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.initStage(g.CurrentStage)
		g.Scene = SceneStageIntro
		g.IntroTimer = 2.4
		g.Audio.PlaySFX("titan_wrath", 0)

		stage := AthensStages[g.CurrentStage]
		SpeakLine(stage.IntroQuip, "el-GR", 1.0, 1.0)
	}
}

func (g *AthensGame) updateStageIntro(dt float64) {
	g.IntroTimer -= dt
	if g.IntroTimer <= 0 || inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.Scene = ScenePlay
	}
}

func (g *AthensGame) initStage(stageIdx int) {
	spec := AthensStages[stageIdx]
	g.Buildings = nil
	g.Citizens = nil
	g.Choppers = nil
	g.Traffic = nil
	g.Bullets = nil
	g.Boss = nil
	g.StageTimer = 0
	g.StageVictimsEaten = 0
	g.StageBuildingsDown = 0
	g.ChopperTimer = 3.0
	g.TrafficTimer = 1.0

	// 1. Build Buildings
	startX := 40.0
	bCount := len(spec.Buildings)
	spacing := (1280.0 - 80.0) / float64(bCount)

	for i, bSpec := range spec.Buildings {
		posX := startX + float64(i)*spacing
		bldg := NewBuilding(bSpec, posX, 640.0)
		g.Buildings = append(g.Buildings, bldg)

		// Populate windows with citizens
		for f := 0; f < bldg.Floors; f++ {
			for c := 0; c < bldg.Cols; c++ {
				if rand.Float64() < 0.40 && len(spec.Victims) > 0 {
					dialog := spec.Victims[rand.Intn(len(spec.Victims))]
					citizen := NewCitizen(dialog, g.victimSprite, bldg, f, c)
					bldg.Slots[f][c].Victim = citizen
					g.Citizens = append(g.Citizens, citizen)
				}
			}
		}
	}

	// 2. Spawn Player Titan
	selectedClass := TitanRoster[g.SelectedTitan]
	g.Titan = NewTitan(selectedClass, g.kronosSprite, 120.0)

	// 3. Stage 5 Boss
	if spec.HasBoss {
		g.Boss = NewMechaMinotaur(1050.0, 640.0-96.0)
	}
}

func (g *AthensGame) updatePlay(dt float64) {
	g.StageTimer += dt

	// 1. Update Particles
	g.Particles.Update(dt)

	// 2. Update Buildings & Check Collapse
	allBuildingsCollapsed := true
	for _, b := range g.Buildings {
		b.Update(dt, g.Particles)
		if !b.Collapsed {
			allBuildingsCollapsed = false
		}
	}

	// 3. Update Player Titan
	g.Titan.Update(dt, g.Buildings, g.Choppers, g.Traffic, g.Boss, g.Particles, g.Audio)

	// Check Game Over
	if g.Titan.Health <= 0 {
		g.Scene = SceneGameOver
		g.Audio.PlaySFX("building_collapse", 0)
		return
	}

	// 4. Update Citizens
	for _, c := range g.Citizens {
		c.Update(dt, g.Titan.X, g.Titan.Y, g.Audio)
	}

	// 5. Spawn & Update Police Choppers ("ΕΛ.ΑΣ.")
	g.ChopperTimer -= dt
	if g.ChopperTimer <= 0 && len(g.Choppers) < 3 {
		g.ChopperTimer = 4.5 + rand.Float64()*3.0
		spawnLeft := rand.Float64() < 0.5
		spawnX := -60.0
		if !spawnLeft {
			spawnX = 1340.0
		}
		chopper := NewPoliceChopper(g.chopperSprite, spawnX, 100.0+rand.Float64()*120.0, !spawnLeft)
		g.Choppers = append(g.Choppers, chopper)
	}

	for _, pc := range g.Choppers {
		pc.Update(dt, g.Titan.X, g.Titan.Y, &g.Bullets, g.Particles)
	}

	// 6. Spawn & Update Street Traffic
	g.TrafficTimer -= dt
	if g.TrafficTimer <= 0 && len(g.Traffic) < 4 {
		g.TrafficTimer = 2.8 + rand.Float64()*2.0
		isTaxi := rand.Float64() < 0.65
		facingLeft := rand.Float64() < 0.5
		startX := -80.0
		if facingLeft {
			startX = 1360.0
		}
		veh := NewStreetVehicle(isTaxi, startX, facingLeft)
		g.Traffic = append(g.Traffic, veh)
	}

	for _, v := range g.Traffic {
		v.Update(dt)
	}

	// 7. Update Bullets & Collision with Titan
	for i := range g.Bullets {
		b := &g.Bullets[i]
		if !b.Active {
			continue
		}
		b.Update(dt)

		// Hit Titan
		if rectsOverlap(b.X, b.Y, 4, 4, g.Titan.X+12, g.Titan.Y+12, g.Titan.W-24, g.Titan.H-24) {
			b.Active = false
			g.Titan.Health -= b.Damage
			g.Trauma = math.Min(1.0, g.Trauma+0.08)
			g.Particles.SpawnBloodFountain(b.X, b.Y, 2, 0.5)
		}
	}

	// 8. Update Boss
	if g.Boss != nil {
		g.Boss.Update(dt, g.Titan.X, g.Titan.Y, g.Particles, g.Audio)

		// Boss laser hits Titan
		if g.Boss.LaserActive {
			laserW := 480.0
			laserX := g.Boss.X - laserW
			if !g.Boss.FacingLeft {
				laserX = g.Boss.X + g.Boss.W
			}
			if rectsOverlap(laserX, g.Boss.Y+22.0, laserW, 14.0, g.Titan.X, g.Titan.Y, g.Titan.W, g.Titan.H) {
				g.Titan.Health -= 90.0 * dt
				g.Trauma = math.Min(1.0, g.Trauma+0.25)
				g.Particles.SpawnBloodFountain(g.Titan.X+g.Titan.W/2, g.Titan.Y+20.0, 3, 1.0)
			}
		}
	}

	// 9. Stage Clear Check
	stageWon := false
	if g.Boss != nil {
		if g.Boss.Health <= 0 {
			stageWon = true
		}
	} else if allBuildingsCollapsed {
		stageWon = true
	}

	if stageWon {
		g.Scene = SceneStageClear
		g.ClearTimer = 3.0
		g.Audio.PlaySFX("pickup_gyros", 0)
	}
}

func (g *AthensGame) updateStageClear(dt float64) {
	g.ClearTimer -= dt
	if g.ClearTimer <= 0 || inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.CurrentStage++
		if g.CurrentStage >= len(AthensStages) {
			g.Scene = SceneVictory
		} else {
			g.initStage(g.CurrentStage)
			g.Scene = SceneStageIntro
			g.IntroTimer = 2.4

			stage := AthensStages[g.CurrentStage]
			SpeakLine(stage.IntroQuip, "el-GR", 1.0, 1.0)
		}
	}
}

func (g *AthensGame) updateVictory() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.CurrentStage = 0
		g.Scene = SceneTitle
	}
}

func (g *AthensGame) updateGameOver() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.initStage(g.CurrentStage)
		g.Scene = ScenePlay
	}
}

// Draw renders all scenes, parallax Athens silhouettes, entities, gore, and HUD.
func (g *AthensGame) Draw(screen *ebiten.Image) {
	switch g.Scene {
	case SceneTitle:
		g.drawTitle(screen)
	case SceneTitanSelect:
		g.drawTitanSelect(screen)
	case SceneStageIntro:
		g.drawStageIntro(screen)
	case ScenePlay:
		g.drawPlay(screen)
	case SceneStageClear:
		g.drawStageClear(screen)
	case SceneVictory:
		g.drawVictory(screen)
	case SceneGameOver:
		g.drawGameOver(screen)
	}
}

func (g *AthensGame) drawTitle(screen *ebiten.Image) {
	// Dark Athens Night Background
	screen.Fill(color.RGBA{16, 12, 28, 255})

	// Title Banner Image
	if g.titleBanner != nil {
		var op ebiten.DrawImageOptions
		bW := float64(g.titleBanner.Bounds().Dx())
		scale := 1280.0 / bW
		op.GeoM.Scale(scale, scale)
		screen.DrawImage(g.titleBanner, &op)
	}

	// Dark Gradient Overlay for Crisp Arcade Text
	drawFilledRect(screen, g.pixelImg, 0, 480, 1280, 240, color.RGBA{10, 8, 18, 220})

	// Arcade Title Text
	var textOp text.DrawOptions
	textOp.GeoM.Translate(640, 520)
	textOp.ColorScale.Scale(1.0, 0.85, 0.2, 1.0) // Gold
	textOp.PrimaryAlign = text.AlignCenter
	textOp.SecondaryAlign = text.AlignCenter
	text.Draw(screen, "ATHENS RAMPAGE: ΤΙΤΑΝΕΣ ΣΤΗΝ ΑΚΡΟΠΟΛΗ", g.faceTitle, &textOp)

	// Blinking Prompt
	pulse := 0.5 + 0.5*math.Sin(float64(time.Now().UnixNano())/200000000.0)
	textOp.GeoM.Reset()
	textOp.GeoM.Translate(640, 580)
	textOp.ColorScale.Reset()
	textOp.ColorScale.Scale(1.0, 0.2, 0.3, float32(pulse))
	text.Draw(screen, "PRESS ENTER OR CLICK TO RAMPAGE!", g.faceSub, &textOp)

	// Subtitle inside jokes
	textOp.GeoM.Reset()
	textOp.GeoM.Translate(640, 640)
	textOp.ColorScale.Reset()
	textOp.ColorScale.Scale(0.8, 0.8, 0.85, 1.0)
	text.Draw(screen, "Γκρεμίστε τις Πολυκατοικίες, Φάτε τους Τουρίστες, Νικήστε την ΑΑΔΕ!", g.faceHUD, &textOp)
}

func (g *AthensGame) drawTitanSelect(screen *ebiten.Image) {
	screen.Fill(color.RGBA{22, 18, 32, 255})

	// Header
	var textOp text.DrawOptions
	textOp.GeoM.Translate(640, 50)
	textOp.ColorScale.Scale(1.0, 0.85, 0.25, 1.0)
	textOp.PrimaryAlign = text.AlignCenter
	textOp.SecondaryAlign = text.AlignCenter
	text.Draw(screen, "ΕΠΙΛΕΞΤΕ ΤΟΝ ΤΙΤΑΝΑ ΣΑΣ", g.faceTitle, &textOp)

	textOp.GeoM.Reset()
	textOp.GeoM.Translate(640, 90)
	textOp.ColorScale.Reset()
	textOp.ColorScale.Scale(0.7, 0.7, 0.8, 1.0)
	textOp.PrimaryAlign = text.AlignCenter
	textOp.SecondaryAlign = text.AlignCenter
	text.Draw(screen, "[Χρησιμοποιήστε ΑΡΙΣΤΕΡΑ / ΔΕΞΙΑ ή A / D και πατήστε ENTER]", g.faceHUD, &textOp)

	// 3 Titan Cards
	cardW := 360.0
	cardH := 460.0
	cardY := 140.0

	for i, titan := range TitanRoster {
		cardX := 60.0 + float64(i)*400.0
		isSelected := (i == g.SelectedTitan)

		cardBg := color.RGBA{35, 30, 48, 240}
		borderCol := color.RGBA{80, 75, 95, 255}
		if isSelected {
			cardBg = color.RGBA{55, 45, 75, 255}
			borderCol = color.RGBA{255, 215, 40, 255}
		}

		drawFilledRect(screen, g.pixelImg, cardX, cardY, cardW, cardH, cardBg)
		drawRectOutline(screen, g.pixelImg, cardX, cardY, cardW, cardH, borderCol, 2.5)

		// Titan Name
		textOp.GeoM.Reset()
		textOp.GeoM.Translate(cardX+cardW/2, cardY+35.0)
		textOp.ColorScale.Reset()
		textOp.ColorScale.Scale(
			float32(titan.TintColor.R)/255.0,
			float32(titan.TintColor.G)/255.0,
			float32(titan.TintColor.B)/255.0,
			1.0,
		)
		textOp.PrimaryAlign = text.AlignCenter
		text.Draw(screen, titan.Name, g.faceTitle, &textOp)

		// Title
		textOp.GeoM.Reset()
		textOp.GeoM.Translate(cardX+cardW/2, cardY+70.0)
		textOp.ColorScale.Reset()
		textOp.ColorScale.Scale(1.0, 0.9, 0.4, 1.0)
		text.Draw(screen, titan.Title, g.faceHUD, &textOp)

		// Portrait for Kronos or Sprite Box
		if i == 0 && g.titanPortrait != nil {
			var pOp ebiten.DrawImageOptions
			pScale := 180.0 / float64(g.titanPortrait.Bounds().Dx())
			pOp.GeoM.Scale(pScale, pScale)
			pOp.GeoM.Translate(cardX+cardW/2-90.0, cardY+100.0)
			screen.DrawImage(g.titanPortrait, &pOp)
		} else {
			// Colored Titan Silhouette Box
			drawFilledRect(screen, g.pixelImg, cardX+cardW/2-70.0, cardY+110.0, 140.0, 160.0, titan.TintColor)
			drawRectOutline(screen, g.pixelImg, cardX+cardW/2-70.0, cardY+110.0, 140.0, 160.0, color.RGBA{20, 20, 25, 255}, 2.0)
		}

		// Stats & Perks
		textOp.GeoM.Reset()
		textOp.GeoM.Translate(cardX+cardW/2, cardY+310.0)
		textOp.ColorScale.Reset()
		textOp.ColorScale.Scale(0.9, 0.9, 0.95, 1.0)
		text.Draw(screen, titan.Description, g.faceHUD, &textOp)

		// Special Move
		textOp.GeoM.Reset()
		textOp.GeoM.Translate(cardX+cardW/2, cardY+380.0)
		textOp.ColorScale.Reset()
		textOp.ColorScale.Scale(
			float32(titan.AuraColor.R)/255.0,
			float32(titan.AuraColor.G)/255.0,
			float32(titan.AuraColor.B)/255.0,
			1.0,
		)
		text.Draw(screen, "SPECIAL: "+titan.SpecialName, g.faceSign, &textOp)
	}

	// Bottom prompt
	textOp.GeoM.Reset()
	textOp.GeoM.Translate(640, 660)
	textOp.ColorScale.Reset()
	textOp.ColorScale.Scale(1.0, 0.85, 0.3, 1.0)
	text.Draw(screen, "ΠΑΤΗΣΤΕ ENTER ΓΙΑ ΝΑ ΞΕΚΙΝΗΣΕΤΕ ΤΗΝ ΚΑΤΑΣΤΡΟΦΗ!", g.faceSub, &textOp)
}

func (g *AthensGame) drawStageIntro(screen *ebiten.Image) {
	screen.Fill(color.RGBA{14, 10, 20, 255})
	stage := AthensStages[g.CurrentStage]

	var textOp text.DrawOptions
	textOp.GeoM.Translate(640, 240)
	textOp.ColorScale.Scale(1.0, 0.85, 0.25, 1.0)
	textOp.PrimaryAlign = text.AlignCenter
	textOp.SecondaryAlign = text.AlignCenter
	text.Draw(screen, stage.Name, g.faceTitle, &textOp)

	textOp.GeoM.Reset()
	textOp.GeoM.Translate(640, 310)
	textOp.ColorScale.Reset()
	textOp.ColorScale.Scale(0.9, 0.4, 0.4, 1.0)
	textOp.PrimaryAlign = text.AlignCenter
	textOp.SecondaryAlign = text.AlignCenter
	text.Draw(screen, stage.Subtitle, g.faceSub, &textOp)

	textOp.GeoM.Reset()
	textOp.GeoM.Translate(640, 420)
	textOp.ColorScale.Reset()
	textOp.ColorScale.Scale(0.85, 0.85, 0.9, 1.0)
	textOp.PrimaryAlign = text.AlignCenter
	textOp.SecondaryAlign = text.AlignCenter
	text.Draw(screen, "\""+stage.IntroQuip+"\"", g.faceHUD, &textOp)
}

func (g *AthensGame) drawPlay(screen *ebiten.Image) {
	spec := AthensStages[g.CurrentStage]

	// Apply camera trauma screen shake
	var shakeX, shakeY float64
	if g.Trauma > 0 {
		shakeMag := math.Pow(g.Trauma, 1.8) * 16.0
		shakeX = (rand.Float64()*2.0 - 1.0) * shakeMag
		shakeY = (rand.Float64()*2.0 - 1.0) * shakeMag
	}

	// 1. Sky Gradient
	drawFilledRect(screen, g.pixelImg, 0, 0, 1280, 720, spec.SkyTop)
	drawFilledRect(screen, g.pixelImg, 0, 280, 1280, 440, spec.SkyBottom)

	// 2. Acropolis / Lycabettus Hill Silhouette in Distance
	drawFilledRect(screen, g.pixelImg, 180, 420, 320, 220, color.RGBA{22, 16, 32, 180})
	drawFilledRect(screen, g.pixelImg, 760, 400, 380, 240, color.RGBA{22, 16, 32, 180})

	// 3. Street Ground & Sidewalk (Y: 640)
	drawFilledRect(screen, g.pixelImg, 0, 640, 1280, 80, color.RGBA{42, 38, 45, 255})
	// Sidewalk curb
	drawFilledRect(screen, g.pixelImg, 0, 636, 1280, 4, color.RGBA{185, 180, 175, 255})

	// 4. Buildings (Polykatoikies / Acropolis)
	for _, b := range g.Buildings {
		b.Draw(screen, g.pixelImg, g.faceSign, g.faceSign)
	}

	// 5. Citizens
	for _, c := range g.Citizens {
		c.Draw(screen, g.pixelImg, g.faceBubble)
	}

	// 6. Street Traffic
	for _, v := range g.Traffic {
		v.Draw(screen, g.pixelImg)
	}

	// 7. Police Choppers
	for _, pc := range g.Choppers {
		pc.Draw(screen, g.pixelImg)
	}

	// 8. Bullets
	for _, bullet := range g.Bullets {
		if bullet.Active {
			drawFilledRect(screen, g.pixelImg, bullet.X, bullet.Y, 4, 4, color.RGBA{255, 240, 60, 255})
		}
	}

	// 9. Boss
	if g.Boss != nil {
		g.Boss.Draw(screen, g.pixelImg)
	}

	// 10. Titan Player
	g.Titan.Draw(screen, g.pixelImg)

	// 11. Gore & Debris Particles
	g.Particles.Draw(screen)
	g.Particles.DrawFloatingTexts(screen, g.faceHUD)

	// 12. HUD Overlay
	g.drawHUD(screen)

	// Offset screen by shake if active
	_ = shakeX
	_ = shakeY
}

func (g *AthensGame) drawHUD(screen *ebiten.Image) {
	// Top Banner Bar
	drawFilledRect(screen, g.pixelImg, 0, 0, 1280, 54, color.RGBA{18, 15, 24, 230})
	drawRectOutline(screen, g.pixelImg, 0, 0, 1280, 54, color.RGBA{60, 55, 70, 255}, 1.5)

	// Titan Name
	var textOp text.DrawOptions
	textOp.GeoM.Translate(20, 27)
	textOp.ColorScale.Scale(
		float32(g.Titan.Class.TintColor.R)/255.0,
		float32(g.Titan.Class.TintColor.G)/255.0,
		float32(g.Titan.Class.TintColor.B)/255.0,
		1.0,
	)
	textOp.PrimaryAlign = text.AlignStart
	textOp.SecondaryAlign = text.AlignCenter
	text.Draw(screen, g.Titan.Class.Name, g.faceHUD, &textOp)

	// Health Bar
	hpRatio := math.Max(0, g.Titan.Health/g.Titan.MaxHealth)
	barW := 220.0
	barH := 16.0
	barX := 130.0
	barY := 19.0
	drawFilledRect(screen, g.pixelImg, barX, barY, barW, barH, color.RGBA{60, 20, 25, 255})
	drawFilledRect(screen, g.pixelImg, barX, barY, barW*hpRatio, barH, color.RGBA{220, 35, 45, 255})
	drawRectOutline(screen, g.pixelImg, barX, barY, barW, barH, color.RGBA{240, 210, 80, 255}, 1.5)

	// Rage / Wrath Meter
	rageRatio := math.Max(0, g.Titan.Rage/100.0)
	rageX := 370.0
	rageW := 180.0
	drawFilledRect(screen, g.pixelImg, rageX, barY, rageW, barH, color.RGBA{30, 20, 50, 255})
	drawFilledRect(screen, g.pixelImg, rageX, barY, rageW*rageRatio, barH, color.RGBA{220, 90, 20, 255})
	drawRectOutline(screen, g.pixelImg, rageX, barY, rageW, barH, color.RGBA{255, 140, 40, 255}, 1.5)

	if g.Titan.Rage >= 100.0 {
		textOp.GeoM.Reset()
		textOp.GeoM.Translate(rageX+rageW/2, barY+barH/2)
		textOp.ColorScale.Reset()
		textOp.ColorScale.Scale(1.0, 1.0, 0.2, 1.0)
		textOp.PrimaryAlign = text.AlignCenter
		text.Draw(screen, "SPACE: WRATH!", g.faceSign, &textOp)
	}

	// Stage Title Center
	stage := AthensStages[g.CurrentStage]
	textOp.GeoM.Reset()
	textOp.GeoM.Translate(740, 27)
	textOp.ColorScale.Reset()
	textOp.ColorScale.Scale(1.0, 0.85, 0.3, 1.0)
	textOp.PrimaryAlign = text.AlignCenter
	text.Draw(screen, stage.Name, g.faceHUD, &textOp)

	// Score Right
	textOp.GeoM.Reset()
	textOp.GeoM.Translate(1160, 27)
	textOp.ColorScale.Reset()
	textOp.ColorScale.Scale(1.0, 1.0, 1.0, 1.0)
	textOp.PrimaryAlign = text.AlignEnd
	scoreStr := fmt.Sprintf("SCORE: %06d", g.Titan.Score)
	text.Draw(screen, scoreStr, g.faceHUD, &textOp)

	// Mute Indicator
	muteTxt := "[M] AUDIO ON"
	if g.Audio.IsMuted() {
		muteTxt = "[M] MUTED"
	}
	textOp.GeoM.Reset()
	textOp.GeoM.Translate(1260, 27)
	textOp.ColorScale.Reset()
	textOp.ColorScale.Scale(0.7, 0.7, 0.7, 1.0)
	textOp.PrimaryAlign = text.AlignEnd
	text.Draw(screen, muteTxt, g.faceSign, &textOp)
}

func (g *AthensGame) drawStageClear(screen *ebiten.Image) {
	screen.Fill(color.RGBA{18, 22, 36, 255})

	var textOp text.DrawOptions
	textOp.GeoM.Translate(640, 180)
	textOp.ColorScale.Scale(1.0, 0.85, 0.25, 1.0)
	textOp.PrimaryAlign = text.AlignCenter
	textOp.SecondaryAlign = text.AlignCenter
	text.Draw(screen, "Η ΠΕΡΙΟΧΗ ΚΑΤΕΔΑΦΙΣΤΗΚΕ!", g.faceTitle, &textOp)

	scoreStr := fmt.Sprintf("SCORE: %d", g.Titan.Score)
	textOp.GeoM.Reset()
	textOp.GeoM.Translate(640, 260)
	textOp.ColorScale.Reset()
	textOp.ColorScale.Scale(1.0, 1.0, 1.0, 1.0)
	textOp.PrimaryAlign = text.AlignCenter
	textOp.SecondaryAlign = text.AlignCenter
	text.Draw(screen, scoreStr, g.faceStats, &textOp)

	textOp.GeoM.Reset()
	textOp.GeoM.Translate(640, 340)
	textOp.ColorScale.Reset()
	textOp.ColorScale.Scale(0.9, 0.3, 0.3, 1.0)
	textOp.PrimaryAlign = text.AlignCenter
	textOp.SecondaryAlign = text.AlignCenter
	text.Draw(screen, "GORE RATING: SSS (ΑΙΜΟΒΟΡΟΣ ΤΙΤΑΝΑΣ)", g.faceSub, &textOp)

	textOp.GeoM.Reset()
	textOp.GeoM.Translate(640, 460)
	textOp.ColorScale.Reset()
	textOp.ColorScale.Scale(1.0, 0.85, 0.3, 1.0)
	textOp.PrimaryAlign = text.AlignCenter
	textOp.SecondaryAlign = text.AlignCenter
	text.Draw(screen, "ΠΑΤΗΣΤΕ ENTER ΓΙΑ ΤΗΝ ΕΠΟΜΕΝΗ ΠΙΣΤΑ", g.faceSub, &textOp)
}

func (g *AthensGame) drawVictory(screen *ebiten.Image) {
	screen.Fill(color.RGBA{28, 12, 18, 255})

	var textOp text.DrawOptions
	textOp.GeoM.Translate(640, 180)
	textOp.ColorScale.Scale(1.0, 0.85, 0.25, 1.0)
	textOp.PrimaryAlign = text.AlignCenter
	textOp.SecondaryAlign = text.AlignCenter
	text.Draw(screen, "ΝΙΚΗ ΣΤΗΝ ΑΚΡΟΠΟΛΗ!", g.faceTitle, &textOp)

	textOp.GeoM.Reset()
	textOp.GeoM.Translate(640, 260)
	textOp.ColorScale.Reset()
	textOp.ColorScale.Scale(1.0, 0.4, 0.4, 1.0)
	textOp.PrimaryAlign = text.AlignCenter
	textOp.SecondaryAlign = text.AlignCenter
	text.Draw(screen, "Ο ΜΗΧΑΝΙΚΟΣ ΜΙΝΩΤΑΥΡΟΣ ΤΗΣ ΑΑΔΕ ΚΑΤΑΣΤΡΑΦΗΚΕ!", g.faceSub, &textOp)

	scoreStr := fmt.Sprintf("ΤΕΛΙΚΟ SCORE: %d", g.Titan.Score)
	textOp.GeoM.Reset()
	textOp.GeoM.Translate(640, 360)
	textOp.ColorScale.Reset()
	textOp.ColorScale.Scale(1.0, 1.0, 1.0, 1.0)
	textOp.PrimaryAlign = text.AlignCenter
	textOp.SecondaryAlign = text.AlignCenter
	text.Draw(screen, scoreStr, g.faceStats, &textOp)

	textOp.GeoM.Reset()
	textOp.GeoM.Translate(640, 500)
	textOp.ColorScale.Reset()
	textOp.ColorScale.Scale(1.0, 0.85, 0.3, 1.0)
	textOp.PrimaryAlign = text.AlignCenter
	textOp.SecondaryAlign = text.AlignCenter
	text.Draw(screen, "ΠΑΤΗΣΤΕ ENTER ΓΙΑ ΝΕΟ ΠΑΙΧΝΙΔΙ", g.faceSub, &textOp)
}

func (g *AthensGame) drawGameOver(screen *ebiten.Image) {
	screen.Fill(color.RGBA{22, 8, 10, 255})

	var textOp text.DrawOptions
	textOp.GeoM.Translate(640, 220)
	textOp.ColorScale.Scale(1.0, 0.15, 0.2, 1.0)
	textOp.PrimaryAlign = text.AlignCenter
	textOp.SecondaryAlign = text.AlignCenter
	text.Draw(screen, "Ο ΤΙΤΑΝΑΣ ΕΠΕΣΕ!", g.faceTitle, &textOp)

	textOp.GeoM.Reset()
	textOp.GeoM.Translate(640, 300)
	textOp.ColorScale.Reset()
	textOp.ColorScale.Scale(0.8, 0.8, 0.85, 1.0)
	textOp.PrimaryAlign = text.AlignCenter
	textOp.SecondaryAlign = text.AlignCenter
	text.Draw(screen, "ΤΑ ΑΠΛΗΡΩΤΑ ΚΟΙΝΟΧΡΗΣΤΑ & ΤΟ ΕΝΦΙΑ ΣΕ ΝΙΚΗΣΑΝ!", g.faceSub, &textOp)

	textOp.GeoM.Reset()
	textOp.GeoM.Translate(640, 420)
	textOp.ColorScale.Reset()
	textOp.ColorScale.Scale(1.0, 0.85, 0.3, 1.0)
	textOp.PrimaryAlign = text.AlignCenter
	textOp.SecondaryAlign = text.AlignCenter
	text.Draw(screen, "ΠΑΤΗΣΤΕ ENTER ΓΙΑ ΝΑ ΞΑΝΑΠΡΟΣΠΑΘΗΣΕΤΕ!", g.faceSub, &textOp)
}
