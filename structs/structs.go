package structs

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ZeroMetaPtr struct {
	_self interface{}
}

type ZeroMetaDef interface {
	This() interface{}
	ThisDef(interface{})
}

type ZeroMeta struct {
	metaptr *ZeroMetaPtr
}

func (meta *ZeroMeta) This() interface{} {
	if meta.metaptr != nil {
		return meta.metaptr._self
	}
	return nil
}

func (meta *ZeroMeta) ThisDef(_self interface{}) {
	meta.metaptr = &ZeroMetaPtr{_self: _self}
}

const DateFormat = "2006-01-02T15:04:05Z"

type Date time.Time

func (t Date) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format(DateFormat))
	return []byte(stamp), nil
}

func (t *Date) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	local, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return err
	}
	tm, err := time.Parse(`"`+DateFormat+`"`, string(data))
	if err != nil {
		return err
	}
	tm = tm.In(local)
	*t = Date(tm)
	return nil
}

func (t Date) Time() time.Time {
	return time.Time(t)
}

type ZeroRequest struct {
	Querys  []interface{}          `json:"querys,omitempty"`
	Expands map[string]interface{} `json:"expands,omitempty"`
}

type ZeroResponse struct {
	Code    int                    `json:"code,omitempty"`
	Message string                 `json:"message,omitempty"`
	Datas   []interface{}          `json:"datas,omitempty"`
	Expands map[string]interface{} `json:"expands,omitempty"`
}

func CheckISO70641983MOD112(idCard string) bool {
	if len(idCard) != 18 {
		return false
	}
	multipliers := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	items := []uint8(strings.ToUpper(idCard))
	sum := 0
	for i := 0; i < 17; i++ {
		sum += int(items[i]-0x30) * multipliers[i]
	}
	check := []uint8{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}
	return items[17] == check[sum%11]
}

func BirthdayWithIDCard(idCard string) (*time.Time, error) {
	if len(idCard) != 18 {
		return nil, errors.New(fmt.Sprintf("error idCard length (%d）", len(idCard)))
	}
	t, err := time.Parse("20060102", idCard[6:14])
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func SexWithIDCard(idCard string) (int, error) {
	c17, err := strconv.Atoi(idCard[16:17])
	if err != nil {
		return -1, err
	}
	return c17 % 2, nil
}

func BytesString(bytes ...byte) string {
	bytesString := ""
	for i, b := range bytes {
		if i != 0 && i%8 == 0 {
			bytesString = fmt.Sprintf("%s\n ", bytesString)
		}
		bytesString = fmt.Sprintf("%s 0x%02X", bytesString, b)
	}
	return fmt.Sprintf("{%s }\n", bytesString)
}

func YearDuration(t time.Time) (time.Time, time.Time, error) {
	startTime := time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
	endTime := startTime.AddDate(1, 0, -1)
	duration1d, err := time.ParseDuration(fmt.Sprintf("%ds", 59+60*59+60*60*23))
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	endTime = endTime.Add(duration1d)
	return startTime, endTime, nil
}

func YearDurationString(t time.Time, xformat string) (string, string, error) {
	startTime, endTime, err := YearDuration(t)
	if err != nil {
		return "", "", err
	}
	return startTime.Format(xformat), endTime.Format(xformat), nil
}

func MonthDuration(t time.Time) (time.Time, time.Time, error) {
	startTime := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	endTime := startTime.AddDate(0, 1, -1)
	duration1d, err := time.ParseDuration(fmt.Sprintf("%ds", 59+60*59+60*60*23))
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	endTime = endTime.Add(duration1d)
	return startTime, endTime, nil
}

func MonthDurationString(t time.Time, xformat string) (string, string, error) {
	startTime, endTime, err := MonthDuration(t)
	if err != nil {
		return "", "", err
	}
	return startTime.Format(xformat), endTime.Format(xformat), nil
}

func DayDuration(t time.Time) (time.Time, time.Time, error) {
	startTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	endTime := startTime
	duration1d, err := time.ParseDuration(fmt.Sprintf("%ds", 59+60*59+60*60*23))
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	endTime = endTime.Add(duration1d)
	return startTime, endTime, nil
}

func DayDurationString(t time.Time, xformat string) (string, string, error) {
	startTime, endTime, err := DayDuration(t)
	if err != nil {
		return "", "", err
	}
	return startTime.Format(xformat), endTime.Format(xformat), nil
}
