package apientry

import api "github.com/dfklegend/cell2/apimapper"

type EntryWithOptions struct {
	Entry api.IAPIEntry
	Opts  []Option
}

type APIEntries struct {
	entries []EntryWithOptions
}

// Register registers a component to hub with options
func (a *APIEntries) Register(e api.IAPIEntry, options ...Option) {
	a.entries = append(a.entries, EntryWithOptions{e, options})
}

// List returns all components with it's options
func (a *APIEntries) List() []EntryWithOptions {
	return a.entries
}
