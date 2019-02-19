package main

import "C"
import (
	"goAddLib/addlib"
	"strings"
)

//export provinces
func provinces(mainland bool) *C.char {
	str := strings.Join(addlib.Provinces(mainland), "\t")
	return C.CString(str)
}

//export cities
func cities(ofProvince *C.char) *C.char {
	str := strings.Join(addlib.Cities(C.GoString(ofProvince)), "\t")
	return C.CString(str)
}

//export districts
func districts(ofCity *C.char) *C.char {
	str := strings.Join(addlib.Districts(C.GoString(ofCity)), "\t")
	return C.CString(str)
}

//export getName
func getName(code *C.char) *C.char {
	return C.CString(addlib.GetName(C.GoString(code)))
}

//export parseAddress
func parseAddress(provinceName *C.char, cityName *C.char, districtName *C.char) *C.char {
	add, err := addlib.ParseAddress(C.GoString(provinceName), C.GoString(cityName), C.GoString(districtName))
	if err != nil {
		return C.CString("")
	}
	return C.CString(add.Province + "\t" + add.City + "\t" + add.District)
}

//export provinceCodes
func provinceCodes(mainland bool) *C.char {
	str := strings.Join(addlib.ProvinceCodes(mainland), "\t")
	return C.CString(str)
}

//export cityCodes
func cityCodes(ofProvinceCode *C.char) *C.char {
	str := strings.Join(addlib.CityCodes(C.GoString(ofProvinceCode)), "\t")
	return C.CString(str)
}

//export districtCodes
func districtCodes(ofCityCode *C.char) *C.char {
	str := strings.Join(addlib.DistrictCodes(C.GoString(ofCityCode)), "\t")
	return C.CString(str)
}

//export getCode
func getCode(provinceName *C.char, cityName *C.char, districtName *C.char) *C.char {
	return C.CString(addlib.GetCode(C.GoString(provinceName), C.GoString(cityName), C.GoString(districtName)))
}

//export getProvinceCode
func getProvinceCode(provinceName *C.char) *C.char {
	return C.CString(addlib.GetProvinceCode(C.GoString(provinceName)))
}

//export getCityCode
func getCityCode(cityName *C.char) *C.char {
	return C.CString(addlib.GetCityCode(C.GoString(cityName)))
}

//export getDistrictCode
func getDistrictCode(cityName *C.char, districtName *C.char) *C.char {
	return C.CString(addlib.GetDistrictCode(C.GoString(cityName), C.GoString(districtName)))
}

//export parseCode
func parseCode(code *C.char) *C.char {
	addCodes, err := addlib.ParseCode(C.GoString(code))
	if err != nil {
		return C.CString("")
	}
	return C.CString(addCodes.ProvinceCode + "\t" + addCodes.CityCode + "\t" + addCodes.DistrictCode)
}

//export initialize
func initialize(dataPath *C.char) *C.char {
	err := addlib.Init(C.GoString(dataPath))
	if err != nil {
		return C.CString(err.Error())
	}
	return C.CString("")
}


/**
 进入goAddLib目录, 然后用如下命令生成动态链接库:
 go build -buildmode=c-shared -o addlib.so ./export.go
 */
func main() {

}