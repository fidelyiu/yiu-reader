package PathUtil

import "testing"

func TestIsValidDir(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "基础测试1",
			args: args{
				str: "E:\\学习文件\\前端",
			},
			want: true,
		},
		{
			name: "基础测试2",
			args: args{
				str: "\\学习文件\\前端",
			},
			want: false,
		},
		{
			name: "不存在的路径",
			args: args{
				str: "G:\\学习文件\\前端",
			},
			want: false,
		},
		{
			name: "错误的路径",
			args: args{
				str: "E:\\\\学习文件\\\\前端",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidDir(tt.args.str); got != tt.want {
				t.Errorf("IsValidDir() = %v, want %v", got, tt.want)
			}
		})
	}
}
