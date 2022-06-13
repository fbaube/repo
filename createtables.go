package repo

// AllTableConfigs configures the three key tables.
var AllTableConfigs = []TableConfig{
	TableConfig_Inbatch,
	TableConfig_Contentity,
	TableConfig_Topicref,
}

// Topicref describes a reference from a Map (i.e. TOC) to a Topic.
// Note that "Topic" does NOT necessarily refer to a DITA `topictref`
// element!
//
// The relationship is N-to-N btwn Maps and Topics, so `Topicref` might
// not be unique because a topic might be explicitly referenced more than
// once by a map. So for simplicity, let's create only one `Topicref` per
// topic per map file, and see if it creates problems elsewhere later on.
//
type Topicref struct {
	Idx_Topicref       int
	Idx_Map_Contentity int
	Idx_Tpc_Contentity int
}

// TableSpec_Topicref describes the table.
var TableSpec_Topicref = DbTblSpec{D_TBL,
	"TRF", "topicref", "Reference from map to topic"}

var ColumnSpecs_Topicref = []DbColSpec{
	// NONE!
}

var TableConfig_Topicref = TableConfig{
	"topicref",
	// ONLY foreign keys
	[]string{"map_contentity", "tpc_contentity"},
	ColumnSpecs_Topicref,
}
