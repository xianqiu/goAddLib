package addlib

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"unicode"
)

// 地址库采用{键: 值}存储, 其中键为地址名称的编码.
// 值采用如下的数据结构(libItem).
type libItem struct {
	code     string
	name     string
	parent   string
	children []string
}

// 地址库
var libItems = make(map[string]*libItem)

// 地址库的中文索引: 中文地址名称 -> 地址编码
// 使用场景: 用非标准的名称查询编码, 然后得到标准的地址名称
var libIndex = make(map[string]string)

// 根节点, 作为"省"的父节点.
const ROOT = "ROOT"

// 统计区划级别
const (
	levelProvince = "province"
	levelCity     = "city"
	levelDistrict = "district"
)

// 数据文件的名称
const (
	dataProvince = "provinces.data"
	dataCity     = "cities.data"
	dataDistrict = "districts.data"
)

func init() {
	githubPath := "github.com\\nullgo\\addlib\\lib.add"
	defaultDataPath := path.Join(os.Getenv("GOPATH"), "src", githubPath)
	err := Init(defaultDataPath)
	if err == nil {
		return
	}
	defaultDataPath = "lib.add"
	Init(defaultDataPath)
}

// 初始化
func Init(dataPath string) error {
	dataFiles := []string{dataProvince, dataCity, dataDistrict}
	if !checkData(dataPath, dataFiles) {
		return errors.New("miss data files")
	}
	libIndexCache := make(map[string]string) // 用于记录key对应的标准地址名称
	// 初始化全局变量(防止Init函数被多次调用)
	libItems = make(map[string]*libItem)
	libIndex = make(map[string]string)
	for _, file := range dataFiles {
		err := loadSingleDataFile(path.Join(dataPath, file), &libItems, &libIndex, &libIndexCache)
		if err != nil {
			return err
		}
	}

	// 删除autoIndex产生的空索引.
	cleanIndex(&libIndex)
	return nil
}

// 按行异步读取单个数据文件
// 初始化LibItems和LibIndex
func loadSingleDataFile(filePath string, ptLibItems *map[string]*libItem,
	ptLibIndex *map[string]string, ptLibIndexCache *map[string]string) error {

	lines := make(chan string, 1)
	go readLines(filePath, lines)
	for line := range lines {
		row := strings.Split(line, "\t")
		level, _ := parseLevelFromFilePath(filePath)
		if !checkRow(row, level) {
			msg := fmt.Sprintf("wrong data format, level = %s, row = %s", level, row)
			return errors.New(msg)
		}
		initLibItems(ptLibItems, row, level)
		err := initLibIndex(ptLibIndex, ptLibIndexCache, row, level)
		if err != nil {
			return err
		}
	}
	return nil
}

// 检查数据文件是否存在
// 输入: filePath - 数据的文件夹名
// 输入: dataFiles - 数据的文件名列表
func checkData(filePath string, dataFiles []string) bool {
	for _, file := range dataFiles {
		_, err := os.Stat(path.Join(filePath, file))
		if err != nil && os.IsNotExist(err) {
			fmt.Println(path.Join(filePath, file))
			return false
		}
	}
	return true
}

// 输入行, 检查数据的格式是否正确
// 规则如下:
// 1. dataProvince: 第1列为自身的编码(字母数字), 第2列为地址名称(仅汉字).
// 2. dataCity和dataDistrict: 第1列为parent编码, 第2列为自身的编码, 第3列为地址名称.
func checkRow(row []string, level string) bool {
	if level == levelProvince {
		if len(row) != 2 {
			return false
		}
		if !(isAllDigitAbc(row[0]) && isAllHanChar(row[1])) {
			return false
		}
	} else if level == levelCity || level == levelDistrict {
		if len(row) != 3 {
			return false
		}
		if !(isAllDigitAbc(row[0]) && isAllDigitAbc(row[1]) && isAllHanChar(row[2])) {
			return false
		}
	}
	return true
}

// 判断字符串是否全部汉字.
func isAllHanChar(str string) bool {
	for _, s := range str {
		if !unicode.Is(unicode.Scripts["Han"], s) {
			return false
		}
	}
	return true
}

// 判断字符串是否全部字母(a-z)或数字(0-9)
func isAllDigitAbc(str string) bool {
	r := []rune(str)
	for _, s := range r {
		if !unicode.IsLetter(s) && !unicode.IsDigit(s) {
			return false
		}
	}
	return true
}

// 初始化libItems.
// 输入每一行数据(row), 需要更新父节点和当前节点.
// 注意:
// 1. 数据文件(provinces.data)由于没有父节点, 其第0列为code, 第1列为name.
// 2. 数据文件(cities.data和districts.data), 其第0列为parent, 第1列为code, 第2列为name.
func initLibItems(ptLibItems *map[string]*libItem, row []string, level string) {
	if level == levelProvince {
		updateParent(ptLibItems, ROOT, row[0])
		updateSelf(ptLibItems, ROOT, row[0], row[1])
	} else {
		updateParent(ptLibItems, row[0], row[1])
		updateSelf(ptLibItems, row[0], row[1], row[2])
	}
}

// 通过文件名确定地址数据的统计区划级别.
// 一共三种级别: PROVINCE, CITY, DISTRICT.
func parseLevelFromFilePath(filePath string) (string, error) {
	switch base := path.Base(filePath); base {
	case dataProvince:
		return levelProvince, nil
	case dataCity:
		return levelCity, nil
	case dataDistrict:
		return levelDistrict, nil
	default:
		msg := fmt.Sprintf("incorrect path, filePath = %s", filePath)
		return "", errors.New(msg)
	}
}

// 初始化中文索引.
// 输入一行数据(row), 用标准名称的前k个汉字作为key.
// 按如下规则:
// 1. k的初始值为2.
// 2. 若key唯一, 则把key当作键值.
// 3. 若key已经存在. 我们把key对应的项叫"旧项", 当前的row叫做"新项".
//    若key的长度小于旧项标准名称的长度, 则把key删除. 然后对新项和旧项重新建索引, 此时k值加1.
// 4. 注意: 旧项的名字不能是新项名字的子字符串(否则算法需要修改).
func initLibIndex(ptLibIndex *map[string]string, ptLibIndexCache *map[string]string,
	row []string, level string) error {

	var err error
	if level == levelProvince {
		// 注意: 省的key名为: <level>-<省名>. 例如: PROVINCE-北京.
		err = autoIndex(ptLibIndex, ptLibIndexCache, level, row[0], row[1], 2)
	} else if level == levelCity {
		// 市的key名为: <level>-<市名>. 例如: CITY-杭州.
		err = autoIndex(ptLibIndex, ptLibIndexCache, level, row[1], row[2], 2)
	} else if level == levelDistrict {
		// 区的key名为: <市编码>-<区名>. 例如: CN033001000-西湖区.
		err = autoIndex(ptLibIndex, ptLibIndexCache, row[0], row[1], row[2], 2)
	}
	if err != nil {
		return err
	}

	return nil
}

// 用递归的方式建立索引.
// 结果写入指针ptLibIndexCache对应的map.
func autoIndex(ptLibIndex *map[string]string, ptLibIndexCache *map[string]string,
	prefix string, code string, name string, initKeySize int) error {

	key, err := formatKey(prefix, name, initKeySize)
	if err != nil {
		return err
	}
	// 若存在重名, 我们把已经存在的code叫"旧项", 另一个待索引的code叫"新项".
	// 需要对旧项或新项设置新的key(通过增加key的长度实现).
	if oldCode, ok := (*ptLibIndex)[key]; ok {
		keySize := getKeySize(key) + 1
		// 数据有重复
		if keySize > len([]rune(name)) {
			fmt.Println(key, initKeySize)
			return errors.New("data has identical items")
		}
		// 如果旧项key的长度无法再增加, 说明旧项是新项的子字符串, 则返回错误.
		if n := len([]rune((*ptLibIndex)[key])); n != 0 && keySize > n {
			return errors.New("sub-name exists")
		}
		// 把旧项key对应的value标记为空 (仅标记删除).
		// 暂时不能删旧项已存在的key, 防止有其它项也有前k个相同的汉字.
		// 如果旧项对应的value已经为空, 则不需要重复设置旧项的新key.
		if (*ptLibIndex)[key] != "" {
			(*ptLibIndex)[key] = ""
			// 旧项设置新key
			err := autoIndex(ptLibIndex, ptLibIndexCache, prefix, oldCode, (*ptLibIndexCache)[key], keySize)
			if err != nil {
				return err
			}
		}
		// 新项设置新key
		err = autoIndex(ptLibIndex, ptLibIndexCache, prefix, code, name, keySize)
		if err != nil {
			return err
		}
	} else {
		(*ptLibIndex)[key] = code
		(*ptLibIndexCache)[key] = name
		return nil
	}
	return nil
}

// 删除空索引
// 即: 删除autoIndex方法中标记为空的索引
func cleanIndex(ptLibIndex *map[string]string) {
	for key := range *ptLibIndex {
		if (*ptLibIndex)[key] == "" {
			delete(*ptLibIndex, key)
		}
	}
}

// 格式化索引的名称.
// 格式: <前缀>-<标准地址名称前k个汉字>
// 输入: prefix - 前缀, name - 标准的地址名称, keySize - 截取的汉字个数
func formatKey(prefix string, name string, keySize int) (string, error) {
	if keySize > len([]rune(name)) {
		msg := fmt.Sprintf("keySize is larger than name size. name = %s, keySize = %d", name, keySize)
		return "", errors.New(msg)
	}
	if prefix == "" {
		return string([]rune(name)[0:keySize]), nil
	}
	return prefix + "-" + string([]rune(name)[0:keySize]), nil
}

// 去除索引的prefix和连接符"-", 然后计算剩余汉字的个数.
func getKeySize(key string) int {
	s := strings.Split(key, "-")
	return len([]rune(s[len(s)-1]))
}

// 按行读文件
// 如果出错, 则想lines写入空字符串.
func readLines(filePath string, lines chan string) {
	file, err := os.Open(filePath)
	if err != nil {
		lines <- ""
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		lines <- scanner.Text()
	}
	defer close(lines)
}

// 给定当前节点(例如杭州市), 创建或更新父节点(例如浙江省)
func updateParent(ptLibItems *map[string]*libItem, parent string, self string) {
	if _, ok := (*ptLibItems)[parent]; !ok {
		// 节点不存在则新增一个空的父节点
		(*ptLibItems)[parent] = &libItem{parent, "", "", make([]string, 0)}
	}
	// 把self作为父节点的孩子(添加到children列表中).
	(*ptLibItems)[parent].children = append((*ptLibItems)[parent].children, self)
}

// 创建或更新当前节点
func updateSelf(ptLibItems *map[string]*libItem, parent string, self string, selfName string) {
	if p, ok := (*ptLibItems)[self]; ok {
		p.name = selfName
		p.parent = parent
	} else {
		(*ptLibItems)[self] = &libItem{self, selfName, parent, make([]string, 0)}
	}
}
