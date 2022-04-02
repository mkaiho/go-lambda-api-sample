package infrastructure

import (
	"testing"

	"github.com/mkaiho/go-lambda-api-sample/entity"
	"github.com/stretchr/testify/assert"
)

const dummyUUIDValue = "b8ce9565-6cd1-4e07-924e-efb2105ccca2"

func Test_uuid_Value(t *testing.T) {
	type fields struct {
		value string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "return string value",
			fields: fields{
				value: dummyUUIDValue,
			},
			want: dummyUUIDValue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &uuid{
				value: tt.fields.value,
			}
			got := id.Value()
			assert.Equal(t, tt.want, got, "uuid.Value() = %v, want %v", got, tt.want)
		})
	}
}

func Test_uuid_IsEmpty(t *testing.T) {
	type fields struct {
		value string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "return true when value size is 0",
			fields: fields{
				value: "",
			},
			want: true,
		},
		{
			name: "return false when value is set",
			fields: fields{
				value: dummyUUIDValue,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &uuid{
				value: tt.fields.value,
			}
			got := id.IsEmpty()
			assert.Equal(t, tt.want, got, "uuid.IsEmpty() = %v, want %v", got, tt.want)
		})
	}
}

func Test_uuidValidator_Validate(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		v       *uuidValidator
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "return no error when value is valid",
			v:    &uuidValidator{},
			args: args{
				value: dummyUUIDValue,
			},
			wantErr: assert.NoError,
		},
		{
			name: "return error when value lengs is shorter than 36",
			v:    &uuidValidator{},
			args: args{
				value: dummyUUIDValue[:35],
			},
			wantErr: assert.Error,
		},
		{
			name: "return error when value lengs is longer than 36",
			v:    &uuidValidator{},
			args: args{
				value: dummyUUIDValue + "0",
			},
			wantErr: assert.Error,
		},
		{
			name: "return error when value contains invalid character",
			v:    &uuidValidator{},
			args: args{
				value: dummyUUIDValue[:35] + "@",
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &uuidValidator{}
			err := v.Validate(tt.args.value)
			tt.wantErr(t, err, "uuidValidator.Validate() error = %v", err)
		})
	}
}

func Test_uuidGenerator_From(t *testing.T) {
	type fields struct {
		validator uuidValidator
	}
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    entity.ID
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "return ID",
			args: args{
				value: dummyUUIDValue,
			},
			want: &uuid{
				value: dummyUUIDValue,
			},
			wantErr: assert.NoError,
		},
		{
			name: "return error when value is invalid",
			args: args{
				value: "@",
			},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &uuidGenerator{
				validator: tt.fields.validator,
			}
			got, err := g.From(tt.args.value)
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got, "uuidGenerator.From() = %v, want %v", got, tt.want)
		})
	}
}
