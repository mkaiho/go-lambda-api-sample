package infrastructure

import (
	"fmt"
	"testing"

	"github.com/mkaiho/go-lambda-api-sample/entity"
	"github.com/stretchr/testify/assert"
)

var testUsers = createDummyUsers(3)

func assertEqualUser(t assert.TestingT, want entity.User, got entity.User) bool {
	if want == nil || got == nil {
		return assert.Equal(t, want, got, "User = %v, want %v", got, want)
	}
	if !assert.Equal(t, want.UserID().Value(), got.UserID().Value(), "User.UserID().Value() = %v, want %v", got.UserID().Value(), want.UserID().Value()) {
		return false
	}
	if !assert.Equal(t, want.Name(), got.Name(), "User.Name() = %v, want %v", got.Name(), want.Name()) {
		return false
	}
	return true
}

func Test_dummyUsersReader_FindAll(t *testing.T) {
	tests := []struct {
		name    string
		want    assert.ValueAssertionFunc
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Return dummy users",
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				for i, user := range i1.([]entity.User) {
					wantUserID := fmt.Sprintf("%03d", i+1)
					if got := user.UserID().Value(); got != wantUserID {
						t.Errorf("dummyUsersReader.FindAll() got.UserID().Value() = %v, want.UserID().Value() %v", got, wantUserID)
						return false
					}
					wantUserName := fmt.Sprintf("testuser_%03d", i+1)
					if got := user.Name(); got != wantUserName {
						t.Errorf("dummyUsersReader.FindAll() got.Name() = %v, want.Name() %v", got, wantUserName)
						return false
					}
				}
				return true
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &dummyUsersReader{}
			got, err := r.FindAll()
			if !tt.wantErr(t, err) {
				return
			}
			tt.want(t, got)
		})
	}
}

func Test_dummyUsersReader_FindByID(t *testing.T) {
	type args struct {
		id entity.UserID
	}
	tests := []struct {
		name    string
		args    args
		want    entity.User
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Return user",
			args: args{
				id: testUsers[0].UserID(),
			},
			want:    testUsers[0],
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &dummyUsersReader{}
			got, err := r.FindByID(tt.args.id)
			if !tt.wantErr(t, err) {
				return
			}
			assertEqualUser(t, tt.want, got)
		})
	}
}
