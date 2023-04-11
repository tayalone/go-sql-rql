package parser

import (
	"reflect"
	"testing"
)

var emptyMap = map[string]interface{}{}

func Test_parser_getFields(t *testing.T) {
	type fields struct {
		defaultLimit   *int
		allFields      map[string]interface{}
		disabledFields map[string]interface{}
	}
	type args struct {
		str string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		// TODO: Add test cases.
		{
			name: "Case: Empty `allFields` && Input Empty",
			fields: fields{
				defaultLimit:   nil,
				allFields:      emptyMap,
				disabledFields: emptyMap,
			},
			args: args{
				str: "",
			},
			want: []string{},
		},
		{
			name: "Case: Set `allFields` && Input empty",
			fields: fields{
				defaultLimit:   nil,
				allFields:      mockAllFields,
				disabledFields: emptyMap,
			},
			args: args{
				str: "id,name",
			},
			want: []string{"id", "name"},
		},
		{
			name: "Case: Set `allFields` && Input Coverage All Fields",
			fields: fields{
				defaultLimit:   nil,
				allFields:      mockAllFields,
				disabledFields: emptyMap,
			},
			args: args{
				str: "id,name",
			},
			want: []string{"id", "name"},
		},
		{
			name: "Case: Set `allFields` && Input Have Some Fields",
			fields: fields{
				defaultLimit:   nil,
				allFields:      mockAllFields,
				disabledFields: emptyMap,
			},
			args: args{
				str: "id",
			},
			want: []string{"id"},
		},
		{
			name: "Case: Set `allFields` && Some Input Not Have Existing in `allFields`",
			fields: fields{
				defaultLimit:   nil,
				allFields:      mockAllFields,
				disabledFields: emptyMap,
			},
			args: args{
				str: "id,not_fields",
			},
			want: []string{"id"},
		},
		{
			name: "Case: Set `allFields` && Set `disabledFields` && Input Have Existing in disabledFields",
			fields: fields{
				defaultLimit:   nil,
				allFields:      mockAllFields,
				disabledFields: mockDisabledFields,
			},
			args: args{
				str: "id,preventField",
			},
			want: []string{"id"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser{
				defaultLimit:   tt.fields.defaultLimit,
				allFields:      tt.fields.allFields,
				disabledFields: tt.fields.disabledFields,
			}
			if got := p.getFields(tt.args.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parser.getFields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parser_getSkip(t *testing.T) {
	type fields struct {
		defaultLimit   *int
		allFields      map[string]interface{}
		disabledFields map[string]interface{}
	}

	f := fields{
		defaultLimit:   nil,
		allFields:      emptyMap,
		disabledFields: emptyMap,
	}

	type args struct {
		str string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name:   "Case: empty skip",
			fields: f,
			args: args{
				str: "",
			},
			want: 0,
		}, {
			name:   "Case: skip is not number",
			fields: f,
			args: args{
				str: "one",
			},
			want: 0,
		}, {
			name:   "Case: skip less than 0",
			fields: f,
			args: args{
				str: "-1",
			},
			want: 0,
		}, {
			name:   "Case: skip is number",
			fields: f,
			args: args{
				str: "112",
			},
			want: 112,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &parser{
				defaultLimit:   tt.fields.defaultLimit,
				allFields:      tt.fields.allFields,
				disabledFields: tt.fields.disabledFields,
			}
			if got := p.getSkip(tt.args.str); got != tt.want {
				t.Errorf("parser.getSkip() = %v, want %v", got, tt.want)
			}
		})
	}
}
