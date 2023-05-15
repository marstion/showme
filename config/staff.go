package config

import (
	"bufio"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/dexyk/stringosim"
)

type Staff struct {
	Staffinfo []StaffInfo
}

type StaffInfo struct {
	Id     string // 序号
	Name   string // 姓名
	Idcard string // 身份证
	Phone  string // 手机号

	score float64 // 匹配成绩
}

// 载入员工台账
func LoadStaffConfig(name string) (s Staff, err error) {

	file, err := os.Open(name)
	if err != nil {
		return
	}
	defer file.Close()
	br := bufio.NewReader(file)

	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		infoList := strings.Split(string(a), "\t")
		if len(infoList) >= 3 {
			var info StaffInfo

			info.Id = strings.Replace(infoList[0], " ", "", -1)
			info.Name = strings.Replace(infoList[1], " ", "", -1)
			info.Idcard = strings.Replace(infoList[2], " ", "", -1)
			// 排除台账不带电话号码的情况
			if len(infoList) == 4 {
				info.Phone = strings.Replace(infoList[3], " ", "", -1)
			}

			s.Staffinfo = append(s.Staffinfo, info)
		}
	}

	return
}

func (s Staff) ComparisonIdCard(idList []string) (StaffInfo, bool) {
	var sl []StaffInfo
	// 累计成员得分， 按分数排序，并返回最优结果
	for _, si := range s.Staffinfo {
		for _, str := range idList {
			si.score = stringosim.Jaro([]rune(si.Idcard), []rune(str))
			if si.score >= 0.74 {
				sl = append(sl, si)
			}
		}
	}
	if len(sl) > 0 {
		sort.Slice(sl, func(i, j int) bool { return sl[i].score > sl[j].score })
		return sl[0], true
	}
	return StaffInfo{}, false
}

func (s Staff) ComparisonName(nameList []string) (StaffInfo, bool) {
	for _, si := range s.Staffinfo {
		for _, str := range nameList {
			if stringosim.Jaro([]rune(si.Name), []rune(str)) >= 0.94 {
				return si, true
			}
		}
	}
	return StaffInfo{}, false
}

func (s Staff) ComparisonPhone(phoneList []string) (StaffInfo, bool) {

	for _, si := range s.Staffinfo {
		for _, str := range phoneList {
			if stringosim.Jaro([]rune(si.Phone), []rune(str)) >= 0.96 {
				return si, true
			}
		}
	}
	return StaffInfo{}, false

}
