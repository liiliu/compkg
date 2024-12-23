package ai

import (
	"testing"
	"weihu_server/library/util"
)

func TestAddAnalysisTask(t *testing.T) {
	InitTest()
	type args struct {
		params AnalysisParams
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
				params: AnalysisParams{
					TaskName:       "analysis",
					CallbackUrl:    "http://116.211.150.24:13006/openApi/v1/aiAnalysisCallback",
					CallbackParams: "{\"conversationId\": \"10444\"}",
					Config:         "{\"input\": {\"llm\": \"qwen:long\", \"file_url\": \"https://smartcard-1253080096.cos.ap-guangzhou.myqcloud.com/audio/2024/09/09/E4B0638557BC_1723195210_E4B0638557BC-1723195210_asr_result.json\"}}",
				},
				token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSIsImV4cGlyZV90aW1lIjoxNzMzODg3MzkxLCJpc3MiOiJhdWRpb19vcGVuX3NlcnZlciIsImV4cCI6MTczMzg4NzM5MSwibmJmIjoxNzMxMjk1MzkxLCJpYXQiOjE3MzEyOTUzOTEsImp0aSI6IjEifQ.aBNz-Yu-h_FY8z2wDxXMdtIE9lM15qBOZP67nnXOPb8",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTaskId, err := AddAnalysisTask(tt.args.params, tt.args.token)
			t.Logf("tt.args.params: %s", util.JsonToString(tt.args.params))
			if (err != nil) != tt.wantErr {
				t.Errorf("AddAnalysisTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else {
				t.Logf("gotTaskId: %s", gotTaskId)
			}
		})
	}
}
