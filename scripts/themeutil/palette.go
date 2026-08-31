package themeutil

type Palette struct {
	Meta     Meta              `json:"meta"`
	Roles    map[string]string `json:"roles"`
	Semantic map[string]string `json:"semantic"`
	Derived  map[string]string `json:"derived,omitempty"`
	Accents  []string          `json:"accents"`
	Terminal map[string]string `json:"terminal"`
	Style    map[string]any    `json:"style"`
	Alpha    AlphaConfig       `json:"alpha"`
}

type Meta struct {
	Name                 string `json:"name"`
	Author               string `json:"author"`
	Appearance           string `json:"appearance"`
	ThemeName            string `json:"theme_name"`
	BackgroundAppearance string `json:"background_appearance"`
	BlurMode             string `json:"blur_mode,omitempty"`
}
