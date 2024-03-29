package artist

import (
	"concert-manager/data"
	"concert-manager/ui/textui/input"
	"concert-manager/ui/textui/output"
	"concert-manager/ui/textui/screens"
)

type Search interface {
	FindFuzzyArtistMatchesByName(string) []data.Artist
}

type Editor struct {
	Search       Search
	SelectScreen screens.Screen
	actions      []string
	artist       *data.Artist
	tempArtist   data.Artist
	returnScreen screens.Screen
}

const (
	search = iota + 1
	setName
	setGenre
	save
	cancel
)

func NewEditScreen() *Editor {
	e := Editor{}
	e.actions = []string{"Search Artists", "Set Name", "Set Genre", "Save Artist", "Cancel"}
	return &e
}

func (e *Editor) AddContext(context screens.ScreenContext) {
	if context.ContextType == screens.Selector {
		e.tempArtist = context.Props[0].(data.Artist)
		return
	}

	e.returnScreen = context.ReturnScreen
	props := context.Props
	e.artist = props[0].(*data.Artist)
	e.tempArtist.Name = e.artist.Name
	e.tempArtist.Genre = e.artist.Genre
}

func (e Editor) Title() string {
	return "Edit Artist"
}

func (e Editor) DisplayData() {
	output.Displayf("%+v\n", e.tempArtist)
}

func (e Editor) Actions() []string {
	return e.actions
}

func (e *Editor) NextScreen(i int) (screens.Screen, *screens.ScreenContext) {
	switch i {
	case search:
		name := input.PromptAndGetInput("artist name", input.NoValidation)
		matches := e.Search.FindFuzzyArtistMatchesByName(name)
		return e.SelectScreen, screens.NewScreenContext(e, matches)
	case setName:
		e.tempArtist.Name = input.PromptAndGetInput("artist name", input.NoValidation)
	case setGenre:
		e.tempArtist.Genre = input.PromptAndGetInput("artist genre", input.NoValidation)
	case save:
		if e.tempArtist.Populated() {
			*e.artist = e.tempArtist
			return e.returnScreen, nil
		} else {
			output.Displayln("Failed to save artist: all fields are required")
		}
	case cancel:
		return e.returnScreen, nil
	}
	return e, nil
}
