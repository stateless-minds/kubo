package plugin

import (
	"github.com/stateless-minds/kubo/repo/fsrepo"
)

// PluginDatastore is an interface that can be implemented to add handlers for
// for different datastores.
type PluginDatastore interface {
	Plugin

	DatastoreTypeName() string
	DatastoreConfigParser() fsrepo.ConfigFromMap
}
