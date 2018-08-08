## 接口描述 (version 1.0)

### 初始化

* **Init(dataPath string) error**  
    
    * 说明: 读取地址库输入到内存. dataPath为标准地址库的路径. 例如`/home/usr/data/lib.add`.
    * 注意: addlib包中已经包含了地址库文件. 默认情况下不需要初始化.  

### 获取名称

* **Provinces() []string**

    * 说明: 输出所有省名.

* **Cities(ofProvince string) []string**
    
    * 说明: 输入省名, 获取它管辖的所有市名.

* **Districts(ofCity string) []string**
    
    * 说明: 输入市名, 获取它管辖的所有区名.

* **GetName(code string) string**
    
    * 说明: 通过编码查询准地址名称. 若输入错误, 则返回"".
    
* **ParseAddress(provinceName string, cityName string, districtName string) (Address, error)**
    
    * 说明: 输入省市区名称, 解析其标准的三级地址名称.  
    例: `ParseAddress("", "杭州", "西湖") -> {浙江省 杭州市 西湖区}`

### 获取编码 

* **ProvinceCodes() []string**  
    
    * 输入: 无  
    * 输出: 所有省编码

* **CityCodes(ofProvinceCode string) []string**
    
    * 输入: 省编码
    * 输出: 它所管辖的市编码. 若输入错误, 则返回空[]

* **DistrictCodes(ofCityCode string) []string**
    
    * 输入: 市编码
    * 输出: 它所管辖的区编码. 若输入错误, 则返回空[]

* **GetCode(provinceName string, cityName string, districtName string) string**
    
    * 说明: 查询地址名称对应的编码.
    
    注意:
    1. 查询区名时, 必须指定市名
    1. 查询市名时, 可以不指定省名
    1. 如果省市区的名字全部指定, 则按照区->市->省的顺序查找, 并返回第一个有效的编码
    1. 要求输入的汉字至少是2个, 而且前两个汉字和标准名称的前2个汉字相同(相同的汉字个数越多, 查询到正确编码的机会越大)

* **GetProvinceCode(provinceName string) string**
    
    * 说明: 查询省编码

* **GetCityCode(cityName string) string**
    
    * 说明: 查询市编码

* **GetDistrictCode(cityName string, districtName string) string**
    
    * 说明: 查询区编码
    
* **ParseCode(code string) (AddressCodes, error)**
    
    * 说明: 输入地址编码, 输出其所属省市区编码.  
    例如: 浙江省杭州市西湖区 = CN033001012.  
    `ParseCode(CN033001012) -> {CN033000000 CN033001000 CN033001012}`