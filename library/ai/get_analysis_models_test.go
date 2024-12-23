package ai

import (
	"testing"
	"weihu_server/library/util"
)

func TestGetAnalysisModels(t *testing.T) {
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
				token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSIsImV4cGlyZV90aW1lIjoxNzMzODg3MzkxLCJpc3MiOiJhdWRpb19vcGVuX3NlcnZlciIsImV4cCI6MTczMzg4NzM5MSwibmJmIjoxNzMxMjk1MzkxLCJpYXQiOjE3MzEyOTUzOTEsImp0aSI6IjEifQ.aBNz-Yu-h_FY8z2wDxXMdtIE9lM15qBOZP67nnXOPb8",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			models, err := GetAnalysisModels(tt.args.token)
			t.Logf("tt.args.params: %s", util.JsonToString(tt.args.params))
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAnalysisModels() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else {
				t.Logf("models: %s", util.JsonToString(models))
			}
		})
	}
}
