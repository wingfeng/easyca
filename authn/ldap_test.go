package authn

import (
	"easyca/conf"
	"testing"
)

func Test_validateUser(t *testing.T) {
	conf.InitConfig("../conf/config.yaml")
	type args struct {
		username string
		password string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "admin",
			args: args{
				username: "caadmin",
				password: "pass@word1",
			},
			want: true,
		},
		{
			name: "wrong password",
			args: args{
				username: "caadmin",
				password: "1pass@word1",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateUser(tt.args.username, tt.args.password); (got == nil) == tt.want {
				t.Errorf("validateUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
