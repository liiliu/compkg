package ai

import (
	"testing"
	"weihu_server/library/util"
)

func TestAddAsrTask(t *testing.T) {
	InitTest()
	type args struct {
		params AsrParams
		token  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{
				params: AsrParams{
					TaskName:       "asr",
					AudioFileUrl:   "https://smartcard-1253080096.cos.ap-guangzhou.myqcloud.com/audio/2024/11/14/6055F9F9F015_1731554768.mp3",
					CallbackUrl:    "http://116.211.150.24:13006/openApi/v1/aiAsrCallback",
					CallbackParams: "{\"conversationId\":1945,\"aiAnalysisId\":1935}",
					Config:         "{\"input\": {\"need_normalize\": false, \"need_denoise\": false, \"is_stereo\": true, \"asr_component\": \"tencent\"},\"out\": {\"audio_type\": \"mp3\", \"asr_type\": \"file\"}}",
				},
				token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSIsImV4cGlyZV90aW1lIjoxNzMzODg3MzkxLCJpc3MiOiJhdWRpb19vcGVuX3NlcnZlciIsImV4cCI6MTczMzg4NzM5MSwibmJmIjoxNzMxMjk1MzkxLCJpYXQiOjE3MzEyOTUzOTEsImp0aSI6IjEifQ.aBNz-Yu-h_FY8z2wDxXMdtIE9lM15qBOZP67nnXOPb8",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTaskId, err := AddAsrTask(tt.args.params, tt.args.token)
			t.Logf("tt.args.params: %s", util.JsonToString(tt.args.params))

			if (err != nil) != tt.wantErr {
				t.Errorf("AddAsrTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else {
				t.Logf("gotTaskId: %s", gotTaskId)
			}
		})
	}
}
