package usecase

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewErrNotFoundEntity(t *testing.T) {
	type args struct {
		entityName string
		params     map[string]string
	}
	tests := []struct {
		name       string
		args       args
		want       EntityNotFoundError
		matchTypes []error
	}{
		{
			name: "return not found entity error",
			args: args{
				entityName: "dummy_entity",
				params: map[string]string{
					"key1": "value1",
				},
			},
			want: EntityNotFoundError{
				useCaseError: useCaseError{
					errType: ErrorTypeNotFoundEntity,
					message: "dummy_entity is not found. params: key1=value1",
				},
				entityName: "dummy_entity",
				params: map[string]string{
					"key1": "value1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewErrNotFoundEntity(tt.args.entityName, tt.args.params)
			assert.Equal(t, tt.want, got, "NewErrNotFoundEntity() = %v, want %v", got, tt.want)

			var errNotFound EntityNotFoundError
			assert.True(t, errors.As(got, &errNotFound))
			var errUseCase UseCaseError
			assert.True(t, errors.As(got, &errUseCase))
			var errFatal FatalError
			assert.False(t, errors.As(got, &errFatal))
		})
	}
}

func TestNewErrFatal(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name string
		args args
		want FatalError
	}{
		{
			name: "return fatal error",
			args: args{
				message: "dummy error",
			},
			want: FatalError{
				useCaseError: useCaseError{
					errType: ErrorTypeFatal,
					message: "dummy error",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewErrFatal(tt.args.message)
			assert.Equal(t, tt.want, got, "NewErrFatal() = %v, want %v", got, tt.want)

			var errFatal FatalError
			assert.True(t, errors.As(got, &errFatal))
			var errUseCase UseCaseError
			assert.True(t, errors.As(got, &errUseCase))
			var errNotFound EntityNotFoundError
			assert.False(t, errors.As(got, &errNotFound))
		})
	}
}
