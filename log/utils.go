package log

import (
	"fmt"
	"strconv"
	"strings"
)

// RGB 得到一个RGB颜色的ANSI颜色代码
func RGB(rgb string) string {
	// 如果RGB字符串以#开头则去掉
	if strings.HasPrefix(rgb, "#") {
		rgb = rgb[1:]
	}

	// 将RGB字符串转换为整数
	r, err := strconv.ParseInt(rgb[0:2], 16, 64)
	if err != nil {
		r = 255
	}
	g, err := strconv.ParseInt(rgb[2:4], 16, 64)
	if err != nil {
		g = 255
	}
	b, err := strconv.ParseInt(rgb[4:6], 16, 64)
	if err != nil {
		b = 255
	}
	// 返回ANSI颜色代码
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
}
