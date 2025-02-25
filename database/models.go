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

type Announcement struct {
	Title       string `bson:"title"`
	Description string `bson:"desc"`
	Link        string `bson:"link"`
	Show        bool   `bson:"show"`
}
