package ai

import "testing"

func TestGetToken(t *testing.T) {
	InitTest()
	tests := []struct {
		name            string
		wantAccessToken string
		wantExpiredTime int64
		wantErr         bool
	}{
		// TODO: Add test cases.
		{
			name:            "test case 1",
			wantAccessToken: "your_access_token",
			wantExpiredTime: 1626640000, // 设置一个固定的过期时间
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAccessToken, gotExpiredTime, err := GetToken()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else {
				t.Logf("gotAccessToken: %s", gotAccessToken)
				t.Logf("gotExpiredTime: %d", gotExpiredTime)
			}
		})
	}
}
