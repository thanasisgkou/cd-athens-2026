//go:build !js || !wasm

package game

// SpeakLine is a no-op on native desktop builds (audio cues play via Ebitengine PCM).
func SpeakLine(text string, lang string, pitch float64, rate float64) {
}
