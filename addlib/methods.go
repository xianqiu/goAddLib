package addlib

import (
	"errors"
	"fmt"
)

type Address struct {
	Province string
	City     string
	District string
}

type AddressCodes struct {
	ProvinceCode string
	CityCode     string
	DistrictCode string
}


// ---------------
// 所有的外部方法 |
// --------------

// 输出所有省编码
func ProvinceCodes(mainland bool) []string {
	islands := [3]string{"CN002000000", "CN027000000", "CN029000000"}
	fullResult := libItems[ROOT].children
	//剔除港澳台
	noIslands := make([]string, len(fullResult) - 3)
	isIsland := false
	if mainland == true {
		k := 0
		for _, v := range fullResult {
			for i := 0; i < len(islands); i++ {
				if v == islands[i] {
					isIsland = true
					break
				}
			}
			if isIsland {
				isIsland = false
				continue
			}
			noIslands[k] = v
			k++
		}
		return noIslands
	} else {
		return fullResult
	}
}

// 输入省编码, 输出它所管辖的市编码
// 若输入错误, 则返回空[]
func CityCodes(ofProvinceCode string) []string {
	if p, ok := libItems[ofProvinceCode]; ok {
		return p.children
	}
	return make([]string, 0)
}

// 输入市编码, 输出它所管辖的区编码
// 若输入错误, 则返回空[]
func DistrictCodes(ofCityCode string) []string {
	if p, ok := libItems[ofCityCode]; ok {
		return p.children
	}
	return make([]string, 0)
}

// 输入编码, 输出其标准地址名称
// 若输入错误, 则返回""
func GetName(code string) string {
	if item, ok := libItems[code]; ok {
		return item.name
	}
	return ""
}

// 输入地址名称, 输出对应的编码
// 说明:
// 1. 查询区名时, 必须指定市名
// 2. 查询市名时, 可以不指定省名
// 3. 如果省市区的名字全部指定, 则按照区->市->省的顺序查找, 并返回第一个有效的编码
// 4. 要求输入的汉字至少是2个, 而且前两个汉字和标准名称的前2个汉字相同(相同的汉字个数越多, 查询到正确编码的机会越大)
func GetCode(provinceName string, cityName string, districtName string) string {

	provinceCode, cityCode, districtCode := "", "", ""
	if provinceName != "" {
		provinceCode = GetProvinceCode(provinceName)
	}
	if cityName != "" {
		cityCode = GetCityCode(cityName)
		if districtName != "" {
			districtCode = GetDistrictCode(cityName, districtName)
		}
	}
	if districtCode != "" {
		return districtCode
	}
	if cityCode != "" {
		return cityCode
	}
	if provinceCode != "" {
		return provinceCode
	}
	return ""
}

// 输入省名, 输出省编码
func GetProvinceCode(provinceName string) string {
	key, _ := formatKey(levelProvince, provinceName, 2)
	if code, ok := libIndex[key]; ok {
		return code
	}
	return ""
}

// 输入市名, 输出市编码
func GetCityCode(cityName string) string {
	minKeySize, maxKeySize := 2, len([]rune(cityName))
	for keySize := minKeySize; keySize <= maxKeySize; keySize++ {
		key, _ := formatKey(levelCity, cityName, keySize)
		if code, ok := libIndex[key]; ok {
			return code
		}
	}
	return ""
}

// 输入市名和区名, 输出区编码
func GetDistrictCode(cityName string, districtName string) string {

	maxKeySize := len([]rune(districtName))
	cityCode := GetCityCode(cityName)
	if cityCode == "" {
		return ""
	}
	for keySize := 2; keySize <= maxKeySize; keySize++ {
		key, _ := formatKey(cityCode, districtName, keySize)
		if code, ok := libIndex[key]; ok {
			return code
		}
	}
	return ""
}

// 输入地址编码, 输出其所属省市区编码.
// 例如: 浙江省杭州市西湖区 = CN033001012
// 		 ParseCode(CN033001012) -> {CN033000000 CN033001000 CN033001012}
func ParseCode(code string) (AddressCodes, error) {

	addc := AddressCodes{"", "", ""}
	parsedCodes := make([]string, 0)
	autoParseCodes(code, &parsedCodes)

	k := len(parsedCodes)
	switch k {
	case 1:
		addc.ProvinceCode = parsedCodes[0]
	case 2:
		addc.ProvinceCode = parsedCodes[1]
		addc.CityCode = parsedCodes[0]
	case 3:
		addc.ProvinceCode = parsedCodes[2]
		addc.CityCode = parsedCodes[1]
		addc.DistrictCode = parsedCodes[0]
	default:
		msg := fmt.Sprintf("invalid code, code = %s", code)
		return addc, errors.New(msg)
	}
	return addc, nil
}

// 递归查找code的父节点, 直到ROOT.
// 结果保存到ptParsedCodes(ROOT节点不保存).
func autoParseCodes(code string, ptParsedCodes *[]string) {

	ptSelf, ok := libItems[code]
	if ok {
		*ptParsedCodes = append(*ptParsedCodes, code)
	} else {
		return
	}

	if parent := ptSelf.parent; parent == ROOT {
		return
	} else {
		autoParseCodes(parent, ptParsedCodes)
	}
}

// 输入省市区名称, 解析其标准三级地址名称
// 例: ParseAddress("", "杭州", "西湖") -> {浙江省 杭州市 西湖区}
func ParseAddress(provinceName string, cityName string, districtName string) (Address, error) {

	add := Address{"", "", ""}
	if districtName != "" && cityName != "" {
		if code := GetDistrictCode(cityName, districtName); code != "" {
			add.District = GetName(code)
			cityCode := libItems[code].parent
			add.City = GetName(cityCode)
			add.Province = GetName(libItems[cityCode].parent)
			return add, nil
		}
	}

	if cityName != "" {
		if code := GetCityCode(cityName); code != "" {
			add.City = GetName(code)
			add.Province = GetName(libItems[code].parent)
			return add, nil
		}
	}

	if provinceName != "" {
		if code := GetProvinceCode(provinceName); code != "" {
			add.Province = GetName(code)
			return add, nil
		}
	}

	msg := fmt.Sprintf("address not found, provinceName = %s, cityName = %s, districtName = %s",
		provinceName, cityName, districtName)
	return add, errors.New(msg)
}

// 输出省名列表
func Provinces(mainland bool) []string {
	provinces := make([]string, 0)
	for _, code := range ProvinceCodes(mainland) {
		provinces = append(provinces, GetName(code))
	}
	return provinces
}

// 输入省名, 输出它管辖的所有市名
func Cities(ofProvince string) []string {
	cities := make([]string, 0)
	provinceCode := GetProvinceCode(ofProvince)
	if provinceCode == "" {
		return cities
	}
	for _, code := range CityCodes(provinceCode) {
		cities = append(cities, GetName(code))
	}
	return cities
}

// 输入市名, 输出它管辖的所有区名
func Districts(ofCity string) []string {
	districts := make([]string, 0)
	cityCode := GetCityCode(ofCity)
	if cityCode == "" {
		return districts
	}
	for _, code := range DistrictCodes(cityCode) {
		districts = append(districts, GetName(code))
	}
	return districts
}
