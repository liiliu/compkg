package ai

import (
	"fmt"
	"testing"
	"weihu_server/library/util"
)

func TestAddCustomAnalysisTask(t *testing.T) {
	InitTest()
	type args struct {
		params CustomAnalysisParams
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
				params: CustomAnalysisParams{
					TaskName:       "run",
					CallbackUrl:    "http://116.211.150.24:13006/openApi/v1/aiAnalysisCallback",
					CallbackParams: "{\"conversationId\": \"1999\"}",
					Config:         "{\"input\":{\"llm\":\"qwen:long\",\"prompts\":[\"以下是一段对话记录：\\n {conversation}\\n \\n根据这段对话记录，分析出这段对话的主要目的或者意图\\n输出结果的指定语言为：{language}\"],\"input_params\":\"{\\\"conversation\\\":\\\"__url_https://smartcard-1253080096.cos.ap-guangzhou.myqcloud.com/audio/2024/11/13/D83BDA891D80_1731465345_asr_result.json\\\",\\\"language\\\":\\\"中文\\\"}\",\"out_formatter\":\"{{\\\"result\\\":{{\\\"purpose\\\": \\\"对话目的\\\"}}}}\"}}",
					//Config:         "{\"input\":{\"llm\":\"qwen:long\",\"prompts\":[\"以下是一段对话记录：\\n{conversation}\\n\\n将上述对话记录提炼为包含标题和综合摘要\\n请直接给出明确的结果。\\n特别强调：\\n- title部分的返回内容限制在100个字数以内\\n- 如果有多个标题，仅生成一条最合适的结果\\n- 摘要内容以条目化列表方式返回呈现\\n\\n输出结果的指定语言为：\\n{language}\\n\"],\"input_params\":\"{\\\"conversation\\\":\\\"__url_https://smartcard-1253080096.cos.ap-guangzhou.myqcloud.com/audio/2024/11/13/D83BDA891D80_1731465345_asr_result.json\\\",\\\"language\\\":\\\"中文\\\"}\",\"out_formatter\":\"{\\\"title\\\": \\\"标题\\\",\\\"summary_list\\\": [{\\\"content\\\": \\\"摘要内容\\\"}]}\"}}",
				},
				token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSIsImV4cGlyZV90aW1lIjoxNzMzODg3MzkxLCJpc3MiOiJhdWRpb19vcGVuX3NlcnZlciIsImV4cCI6MTczMzg4NzM5MSwibmJmIjoxNzMxMjk1MzkxLCJpYXQiOjE3MzEyOTUzOTEsImp0aSI6IjEifQ.aBNz-Yu-h_FY8z2wDxXMdtIE9lM15qBOZP67nnXOPb8",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(tt.args.params.Config)
			gotTaskId, err := AddCustomAnalysisTask(tt.args.params, tt.args.token)
			t.Logf("tt.args.params: %s", util.JsonToString(tt.args.params))
			if (err != nil) != tt.wantErr {
				t.Errorf("AddCustomAnalysisTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else {
				t.Logf("gotTaskId: %s", gotTaskId)
			}
		})
	}
}
