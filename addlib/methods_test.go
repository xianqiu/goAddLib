package addlib

import (
	"fmt"
	"testing"
)

func TestProvinceCodes(t *testing.T) {
	expected := 34
	expectedMainland := 31
	if got := len(ProvinceCodes(false)); got != expected {
		t.Errorf("expected: %d, got: %d", expected, got)
	}
	if got := len(ProvinceCodes(true)); got != expectedMainland {
		t.Errorf("expected: %d, got: %d", expectedMainland, got)
	}
}

func TestCityCodes(t *testing.T) {

	tests := []struct {
		in       string
		expected int
	}{
		{GetCode("浙江", "", ""), 11},
		{GetCode("广西", "", ""), 14},
		{GetCode("新疆", "", ""), 15},
		{GetCode("上海", "", ""), 1},
		{"foo", 0},
	}

	for _, tt := range tests {
		if got := len(CityCodes(tt.in)); got != tt.expected {
			t.Errorf("expected: %d, got: %d", tt.expected, got)
		}
	}
}

func TestDistrictCodes(t *testing.T) {
	tests := []struct {
		in       string
		expected int
	}{
		{GetCode("", "杭州", ""), 13},
		{GetCode("", "北京", ""), 16},
		{GetCode("", "重庆", ""), 38},
		{GetCode("", "呼和浩特", ""), 9},
		{"foo", 0},
	}

	for _, tt := range tests {
		if got := len(DistrictCodes(tt.in)); got != tt.expected {
			t.Errorf("expected: %d, got: %d", tt.expected, got)
		}
	}
}

func TestGetName(t *testing.T) {

	s0 := "浙江省"
	s1 := "杭州市"
	s2 := "西湖区"
	s3 := "北京市"

	tests := []struct {
		in       string
		expected string
	}{
		{GetCode("", s1, ""), s1},
		{GetCode(s0, "", ""), s0},
		{GetCode("", s1, s2), s2},
		{GetCode(s3, s3, ""), s3},
		{"foo", ""},
	}

	for _, tt := range tests {
		if got := GetName(tt.in); got != tt.expected {
			t.Errorf("expected: %s, got: %s", tt.expected, got)
		}
	}
}

func TestGetCode(t *testing.T) {
	tests := []struct {
		inProvince string
		inCity     string
		inDistrict string
		expected   string
	}{
		{"浙江", "杭州", "西湖", "CN033001012"},
		{"", "杭州", "西湖", "CN033001012"},
		{"浙江", "杭州", "", "CN033001000"},
		{"浙江", "", "", "CN033000000"},
		{"", "巴彦淖尔市", "乌拉特后旗", "CN019003004"},
		{"", "巴彦淖尔市", "乌拉特中旗", "CN019003006"},
		{"", "巴彦淖尔市", "乌拉特前旗", "CN019003005"},
		{"foo", "bar", "", ""},
		{"啊", "", "", ""},
		{"", "哈", "", ""},
		{"", "哈", "哦", ""},
	}

	for _, tt := range tests {
		if got := GetCode(tt.inProvince, tt.inCity, tt.inDistrict); got != tt.expected {
			t.Errorf("expected: %s, got: %s", tt.expected, got)
		}
	}
}

func TestGetProvinceCode(t *testing.T) {

	tests := []struct {
		in       string
		expected string
	}{
		{"浙江", "CN033000000"},
		{"澳门", "CN002000000"},
		{"天津", "CN028000000"},
		{"宁夏", "CN020000000"},
		{"foo", ""},
	}

	for _, tt := range tests {
		if got := GetProvinceCode(tt.in); got != tt.expected {
			t.Errorf("expected: %s, got: %s", tt.expected, got)
		}
	}
}

func TestGetCityCode(t *testing.T) {
	tests := []struct {
		in       string
		expected string
	}{
		{"张家界", "CN014012000"},
		{"张家口", "CN010012000"},
		{"常德", "CN014001000"},
		{"重庆", "CN034002000"},
		{"foo", ""},
	}

	for _, tt := range tests {
		if got := GetCityCode(tt.in); got != tt.expected {
			t.Errorf("expected: %s, got: %s", tt.expected, got)
		}
	}
}

func TestGetDistrictCode(t *testing.T) {
	tests := []struct {
		inCity     string
		inDistrict string
		expected   string
	}{
		{"杭州", "西湖", "CN033001012"},
		{"南昌", "西湖", "CN016006008"},
		{"苏州", "昆山", "CN015007004"},
		{"foo", "bar", ""},
	}

	for _, tt := range tests {
		if got := GetDistrictCode(tt.inCity, tt.inDistrict); got != tt.expected {
			t.Errorf("expected: %s, got: %s", tt.expected, got)
		}
	}
}

func TestParseCode(t *testing.T) {
	tests := []struct {
		in       string
		expected AddressCodes
	}{
		{"CN033001012", AddressCodes{"CN033000000", "CN033001000", "CN033001012"}},
		{"CN033001000", AddressCodes{"CN033000000", "CN033001000", ""}},
		{"CN033000000", AddressCodes{"CN033000000", "", ""}},
		{"foo", AddressCodes{"", "", ""}},
	}
	for _, tt := range tests {
		got, _ := ParseCode(tt.in)
		strGot := fmt.Sprintf("%s-%s-%s", got.ProvinceCode, got.CityCode, got.DistrictCode)
		strExpected := fmt.Sprintf("%s-%s-%s", tt.expected.ProvinceCode, tt.expected.CityCode, tt.expected.DistrictCode)
		if strGot != strExpected {
			t.Errorf("expected: %s, got: %s", strExpected, strGot)
		}
	}
}

func TestParseAddress(t *testing.T) {
	tests := []struct {
		inProvince string
		inCity     string
		inDistrict string
		expected   Address
	}{
		{"浙江", "杭州", "西湖", Address{"浙江省", "杭州市", "西湖区"}},
		{"", "杭州", "西湖", Address{"浙江省", "杭州市", "西湖区"}},
		{"", "杭州", "", Address{"浙江省", "杭州市", ""}},
		{"浙江", "", "", Address{"浙江省", "", ""}},
		{"foo", "bar", "", Address{"", "", ""}},
	}
	for _, tt := range tests {
		got, _ := ParseAddress(tt.inProvince, tt.inCity, tt.inDistrict)
		strGot := fmt.Sprintf("%s-%s-%s", got.Province, got.City, got.District)
		strExpected := fmt.Sprintf("%s-%s-%s", tt.expected.Province, tt.expected.City, tt.expected.District)
		if strGot != strExpected {
			t.Errorf("expected: %s, got: %s", strExpected, strGot)
		}
	}
}

func TestProvinces(t *testing.T) {
	expected := 34
	expectedMainland := 31
	if got := len(Provinces(false)); got != expected {
		t.Errorf("expected: %d, got: %d", expected, got)
	}
	if got := len(Provinces(true)); got != expectedMainland {
		t.Errorf("expected: %d, got: %d", expected, got)
	}
}

func TestCities(t *testing.T) {
	tests := []struct {
		in       string
		expected int
	}{
		{"浙江", 11},
		{"广西", 14},
		{"新疆", 15},
		{"上海", 1},
		{"foo", 0},
	}

	for _, tt := range tests {
		if got := len(Cities(tt.in)); got != tt.expected {
			t.Errorf("expected: %d, got: %d", tt.expected, got)
		}
	}
}

func TestDistricts(t *testing.T) {
	tests := []struct {
		in       string
		expected int
	}{
		{"杭州", 13},
		{"北京", 16},
		{"重庆", 38},
		{"呼和浩特", 9},
		{"foo", 0},
	}

	for _, tt := range tests {
		if got := len(Districts(tt.in)); got != tt.expected {
			t.Errorf("expected: %d, got: %d", tt.expected, got)
		}
	}
}

func TestGetCode2(t *testing.T) {
	for _, p := range Provinces(false) {
		pGot := GetName(GetCode(p, "", ""))
		if p != pGot {
			t.Errorf("expected: %s, got: %s", p, pGot)
		}
		for _, c := range Cities(p) {
			cGot := GetName(GetCode(p, c, ""))
			if c != cGot {
				t.Errorf("expected: %s, got: %s", c, cGot)
			}
			for _, d := range Districts(c) {
				dGot := GetName(GetCode(p, c, d))
				if d != dGot {
					t.Errorf("expected: %s, got: %s", d, dGot)
				}
			}
		}
	}
}

func BenchmarkGetCode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, p := range Provinces(false) {
			for _, c := range Cities(p) {
				for _, d := range Districts(c) {
					fmt.Println(p, c, d)
				}
			}
		}
	}
}

// 加载测试数据test.add
func init() {
	testDataPath := "lib.add"
	Init(testDataPath)
}
