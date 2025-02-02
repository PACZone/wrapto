package database

type Log struct {
	Actor       string `bson:"actor"`
	Description string `bson:"description"`
	Trace       string `bson:"trace"`
	OrderID     string `bson:"order_id"`
	CreatedAt   int64  `bson:"created_at"`
}

type State struct {
	Pactus  uint32 `bson:"pactus"`
	Polygon uint32 `bson:"polygon"`
}
