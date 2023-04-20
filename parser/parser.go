package parser

// Operation of query
var Operation = map[string]interface{}{
	"$eq":         struct{}{},
	"$is":         struct{}{},
	"$gt":         struct{}{},
	"$gte":        struct{}{},
	"$lt":         struct{}{},
	"$lte":        struct{}{},
	"$between":    struct{}{},
	"$notBetween": struct{}{},
	"$in":         struct{}{},
}

// LinkOperation of query
var LinkOperation = map[string]interface{}{
	"$and": struct{}{},
	"$or":  struct{}{},
}

// NotOperation of query
var NotOperation = map[string]interface{}{
	"$not": struct{}{},
}
