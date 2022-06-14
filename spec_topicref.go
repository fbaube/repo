package repo

// Topicref describes a reference from a Map (i.e. TOC) to a Topic.
// Note that "Topic" does NOT necessarily refer to a DITA `topictref`
// element!
//
// The relationship is N-to-N btwn Maps and Topics, so `Topicref` might
// not be unique because a topic might be explicitly referenced more than
// once by a map. So for simplicity, let's create only one `Topicref` per
// topic per map file, and see if it creates problems elsewhere later on.
//
// Note also that if we decide to use multi-trees, then perhaps these links
// can count not just as kids for maps, but also as parents for topics.
//
type Topicref struct {
	Idx_Topicref       int
	Idx_Map_Contentity int
	Idx_Tpc_Contentity int
}

// TableSpec_Topicref describes table TOPICREF.
//
var TableSpec_Topicref = DbTblSpec{D_TBL,
	"TRF", "topicref", "Reference from map to topic"}

// ColumnSpecs_Topicref is empty, cos the table contains only foreign keys.
//
var ColumnSpecs_Topicref = []DbColSpec{
	// NONE!
}

// TableConfig_Topicref specifies only two foreign keys.
//
var TableConfig_Topicref = TableConfig{
	"topicref",
	// ONLY foreign keys
	[]string{"map_contentity", "tpc_contentity"},
	ColumnSpecs_Topicref,
}
