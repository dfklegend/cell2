package csvutils

import (
	"bufio"
	"os"
	"strings"

	"github.com/gocarina/gocsv"
)

// 自定义的属性只要有UnmarshalCSV接口即可被序列化读取
//type TypeUnmarshaller interface {
//	UnmarshalCSV(string) error
//}

// LoadFromFile 加载csv,默认表头行=2,起始行=3
func LoadFromFile(path string, out interface{}) error {
	return LoadFromFileWithIndex(path, out, 2, 3)
}

// LoadFromFileWithIndex 指定表头行和起始行
func LoadFromFileWithIndex(path string, out interface{}, headIdx int, beginIdx int) error {

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	sb := strings.Builder{}
	i := 0
	for scanner.Scan() {
		text := scanner.Text()
		i++
		if len(text) == 0 {
			continue
		}
		if i == headIdx {
			sb.WriteString(text)
			sb.WriteString("\r\n")
		}
		if i >= beginIdx {
			sb.WriteString(text)
			sb.WriteString("\r\n")
		}
	}
	return gocsv.UnmarshalString(sb.String(), out)
}
