package assets

import "embed"

//go:embed sprites/*.sprite audio/*.sfx audio/*.music fonts/font.ttf images/*.png
var FS embed.FS
