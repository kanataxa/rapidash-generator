package generator

type Config struct {
	ShouldOverwrite bool
	FilePath        string
	Output          string
	Tag             string
	DependenceFiles []string
}
