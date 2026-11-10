package game

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"athens-rampage/assets"
	"athens-rampage/internal/audio/procmusic"
	"athens-rampage/internal/audio/procsfx"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

// AudioSystem manages all sound effects and procedural background music for Athens Rampage.
type AudioSystem struct {
	mu           sync.Mutex
	ctx          *audio.Context
	sfxPCM       map[string][]byte
	lastPlayed   map[string]time.Time
	bgmPlayer    *audio.Player
	bgmPCM       []byte
	soundSys     *procsfx.SoundSystem
	musicSys     *procmusic.SoundSystem
	muted        bool
	masterVolume float64
}

var (
	globalAudio *AudioSystem
	audioOnce   sync.Once
)

// GetAudioSystem returns the singleton AudioSystem, initializing it on first call.
func GetAudioSystem() *AudioSystem {
	audioOnce.Do(func() {
		globalAudio = newAudioSystem()
	})
	return globalAudio
}

func newAudioSystem() *AudioSystem {
	ctx := audio.NewContext(procsfx.SampleRate)
	as := &AudioSystem{
		ctx:          ctx,
		sfxPCM:       make(map[string][]byte),
		lastPlayed:   make(map[string]time.Time),
		soundSys:     procsfx.NewSoundSystem(ctx),
		musicSys:     procmusic.NewSoundSystem(ctx),
		masterVolume: 1.0,
	}

	as.loadAllSFX()
	as.loadBGM()
	return as
}

func (as *AudioSystem) loadAllSFX() {
	sfxFiles := []string{
		"devour_gore",
		"victim_scream",
		"titan_punch",
		"building_collapse",
		"titan_wrath",
		"pickup_gyros",
		"yiayia_shock",
	}

	for _, name := range sfxFiles {
		path := fmt.Sprintf("audio/%s.sfx", name)
		data, err := assets.FS.ReadFile(path)
		if err != nil {
			fmt.Printf("[Audio] Warning: failed to read %s: %v\n", path, err)
			continue
		}

		pcm, err := procsfx.CompileDSL(string(data))
		if err != nil {
			fmt.Printf("[Audio] Warning: failed to compile %s: %v\n", path, err)
			continue
		}

		// Ensure 4-byte stereo alignment
		alignedLen := len(pcm) &^ 3
		if alignedLen > 0 {
			as.sfxPCM[name] = pcm[:alignedLen]
		}
	}
	fmt.Printf("[Audio] Successfully pre-compiled %d procedural SFX cues.\n", len(as.sfxPCM))
}

func (as *AudioSystem) loadBGM() {
	musicData, err := assets.FS.ReadFile("audio/athens_rampage.music")
	if err != nil {
		fmt.Printf("[Audio] Warning: failed to read BGM: %v\n", err)
		return
	}

	pcm, err := procmusic.CompileDSL(string(musicData))
	if err != nil {
		fmt.Printf("[Audio] Warning: failed to compile BGM: %v\n", err)
		return
	}

	alignedLen := len(pcm) &^ 3
	if alignedLen <= 0 {
		return
	}
	as.bgmPCM = pcm[:alignedLen]

	loopStream := audio.NewInfiniteLoop(bytes.NewReader(as.bgmPCM), int64(len(as.bgmPCM)))
	player, err := as.ctx.NewPlayer(loopStream)
	if err != nil {
		fmt.Printf("[Audio] Warning: failed to create BGM player: %v\n", err)
		return
	}

	player.SetVolume(0.55 * as.masterVolume)
	as.bgmPlayer = player
	fmt.Printf("[Audio] Pre-compiled and buffered BGM (%d bytes).\n", len(as.bgmPCM))
}

// PlayBGM starts looping the high-tempo SEGA Mega Drive FM background theme.
func (as *AudioSystem) PlayBGM() {
	as.mu.Lock()
	defer as.mu.Unlock()

	if as.muted || as.bgmPlayer == nil {
		return
	}
	if !as.bgmPlayer.IsPlaying() {
		_ = as.bgmPlayer.Rewind()
		as.bgmPlayer.Play()
	}
}

// StopBGM pauses the background music.
func (as *AudioSystem) StopBGM() {
	as.mu.Lock()
	defer as.mu.Unlock()

	if as.bgmPlayer != nil && as.bgmPlayer.IsPlaying() {
		as.bgmPlayer.Pause()
	}
}

// PlaySFX plays a pre-compiled procedural sound effect with cooldown protection.
func (as *AudioSystem) PlaySFX(cueName string, cooldown time.Duration) {
	as.mu.Lock()
	defer as.mu.Unlock()

	if as.muted {
		return
	}

	pcm, ok := as.sfxPCM[cueName]
	if !ok || len(pcm) == 0 {
		return
	}

	now := time.Now()
	if last, exists := as.lastPlayed[cueName]; exists {
		if now.Sub(last) < cooldown {
			return
		}
	}
	as.lastPlayed[cueName] = now

	// Play via audio player from buffer
	p, err := as.ctx.NewPlayer(bytes.NewReader(pcm))
	if err == nil {
		p.SetVolume(0.85 * as.masterVolume)
		p.Play()
	}
}

// ToggleMute toggles mute on/off.
func (as *AudioSystem) ToggleMute() bool {
	as.mu.Lock()
	defer as.mu.Unlock()

	as.muted = !as.muted
	if as.muted {
		if as.bgmPlayer != nil && as.bgmPlayer.IsPlaying() {
			as.bgmPlayer.Pause()
		}
	} else {
		if as.bgmPlayer != nil {
			as.bgmPlayer.Play()
		}
	}
	return as.muted
}

// IsMuted returns true if audio is muted.
func (as *AudioSystem) IsMuted() bool {
	as.mu.Lock()
	defer as.mu.Unlock()
	return as.muted
}
