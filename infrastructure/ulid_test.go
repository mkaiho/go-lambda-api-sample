package infrastructure

import (
	"sync"
	"testing"

	"github.com/mkaiho/go-lambda-api-sample/entity"
	"github.com/stretchr/testify/assert"
)

const dummyULIDValue = "01fwndwhzdgzhjcmc82m73sh64"

func Test_ulid_Value(t *testing.T) {
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
				value: dummyULIDValue,
			},
			want: dummyULIDValue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &ulid{
				value: tt.fields.value,
			}
			got := id.Value()
			assert.Equal(t, tt.want, got, "ulid.Value() = %v, want %v", got, tt.want)
		})
	}
}

func Test_ulid_IsEmpty(t *testing.T) {
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
				value: dummyULIDValue,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &ulid{
				value: tt.fields.value,
			}
			got := id.IsEmpty()
			assert.Equal(t, tt.want, got, "ulid.IsEmpty() = %v, want %v", got, tt.want)
		})
	}
}

func Test_ulidValidator_Validate(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		v       *ulidValidator
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "return no error when value is valid",
			v:    &ulidValidator{},
			args: args{
				value: dummyULIDValue,
			},
			wantErr: assert.NoError,
		},
		{
			name: "return error when value lengs is shorter than 26",
			v:    &ulidValidator{},
			args: args{
				value: dummyULIDValue[:25],
			},
			wantErr: assert.Error,
		},
		{
			name: "return error when value lengs is longer than 26",
			v:    &ulidValidator{},
			args: args{
				value: dummyULIDValue + "0",
			},
			wantErr: assert.Error,
		},
		{
			name: "return error when value contains invalid character",
			v:    &ulidValidator{},
			args: args{
				value: dummyULIDValue[:25] + "@",
			},
			wantErr: assert.Error,
		},
		{
			name: "return error when value contains invalid character",
			v:    &ulidValidator{},
			args: args{
				value: "1l000000000000000000000000",
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &ulidValidator{}
			err := v.Validate(tt.args.value)
			tt.wantErr(t, err, "ulidValidator.Validate() error = %v", err)
		})
	}
}

func Test_ulidGenerator_From(t *testing.T) {
	type fields struct {
		pool      sync.Pool
		validator ulidValidator
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
				value: dummyULIDValue,
			},
			want: &ulid{
				value: dummyULIDValue,
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
			g := &ulidGenerator{
				pool:      tt.fields.pool,
				validator: tt.fields.validator,
			}
			got, err := g.From(tt.args.value)
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got, "ulidGenerator.From() = %v, want %v", got, tt.want)
		})
	}
}
