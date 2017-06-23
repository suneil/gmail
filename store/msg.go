package store

type Message struct {
	ID           string   `bson:"_id,omitempty" json:"id,omitempty"`
	HistoryID    uint64   `bson:"historyId,omitempty"`
	InternalDate int64    `bson:"internalDate,omitempty"`
	ThreadID     string   `bson:"threadId,omitempty"`
	Snippet      string   `bson:"snippet,omitempty"`
	LabelIDs     []string `bson:"labelIds,omitempty"`
	Raw          string   `bson:"raw,omitempty"`
	From         string   `bson:"from,omitempty"`
	To           string   `bson:"to,omitempty"`
	Subject      string   `bson:"subject,omitempty"`
	Date         string   `bson:"date,omitempty"`
}


