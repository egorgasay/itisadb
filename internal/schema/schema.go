package schema

type KeyValue struct {
	Key   string `bson:"Key"`
	Value string `bson:"Value"`
}
