package registry

import (
	"strings"

	api "github.com/dfklegend/cell2/apimapper"
	"github.com/dfklegend/cell2/apimapper/apientry"
)

func RegisterWithLowercaseName(entry api.IAPIEntry, collection string, entryName string) {
	Registry.AddCollection(collection).Register(entry, apientry.WithGroupName(entryName), apientry.WithNameFunc(strings.ToLower))
}
