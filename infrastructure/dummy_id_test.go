package infrastructure

import (
	"testing"

	"github.com/mkaiho/go-lambda-api-sample/entity"
	"github.com/stretchr/testify/assert"
)

func Test_dummyID_Value(t *testing.T) {
	tests := []struct {
		name string
		id   dummyID
		want string
	}{
		{
			name: "Return string value",
			id:   dummyID("test_001"),
			want: "test_001",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.id.Value()
			assert.Equal(t, tt.want, got, "dummyID.Value() = %v, want %v", got, tt.want)
		})
	}
}

func Test_dummyID_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		id   dummyID
		want bool
	}{
		{
			name: "return true when value size is 0",
			id:   "",
			want: true,
		},
		{
			name: "return true when value size is 1",
			id:   "x",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.id.IsEmpty()
			assert.Equal(t, tt.want, got, "dummyID.IsEmpty() = %v, want %v", got, tt.want)
		})
	}
}

func Test_dummyIDValidator_Validate(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Return no error",
			args: args{
				value: "test",
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &dummyIDValidator{}
			err := v.Validate(tt.args.value)
			tt.wantErr(t, err)
		})
	}
}

func Test_dummyIDGenerator_From(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    entity.ID
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Return ID",
			args: args{
				value: "test_001",
			},
			want:    dummyID("test_001"),
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &dummyIDGenerator{
				validator: dummyIDValidator{},
			}
			got, err := g.From(tt.args.value)
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got, "dummyIDGenerator.From() = %v, want %v", got, tt.want)
		})
	}
}

func Test_dummyIDGenerator_Generate(t *testing.T) {
	tests := []struct {
		name string
		want entity.ID
	}{
		{
			name: "Return new ID",
			want: dummyID("001"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &dummyIDGenerator{
				validator: dummyIDValidator{},
			}
			got := g.Generate()
			assert.Equal(t, tt.want, got, "dummyIDGenerator.Generate() = %v, want %v", got, tt.want)
		})
	}
}
