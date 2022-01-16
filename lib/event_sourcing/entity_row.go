package eventsourcing

type EntityRow struct {
	EntityType    string
	EntityUuid    string
	EntityVersion int
}
