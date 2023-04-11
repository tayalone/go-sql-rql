package param

// QueryParams is Obj with provide SQL Query Params
type QueryParams interface {
	SetFields(fields []string)
	GetFields() []string
	SetSkip(skip int)
	GetSkip() int
}

// QueryParams struct
type queryParams struct {
	fields []string
	skip   int
}

// New QueryParams
func New() QueryParams {
	return &queryParams{
		fields: []string{},
		skip:   0,
	}
}

func (qp *queryParams) SetFields(fields []string) {
	qp.fields = fields
}

func (qp *queryParams) GetFields() []string {
	return qp.fields
}

func (qp *queryParams) SetSkip(skip int) {
	qp.skip = skip
}

func (qp *queryParams) GetSkip() int {
	return qp.skip
}
