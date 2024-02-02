package structs

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
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

	var err error
	tm, err := time.Parse(`"`+DateFormat+`"`, string(data))
	*t = Date(tm)
	return err
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

type ZeroCoreStructs struct {
	ZeroMeta

	ID         string                 `json:"id,omitempty"`
	CreateTime Date                   `json:"createTime,omitempty"`
	UpdateTime Date                   `json:"updateTime,omitempty"`
	Features   map[string]interface{} `json:"features,omitempty"`
}

func (e *ZeroCoreStructs) UInt8ToString(bs []uint8) string {
	ba := []byte{}
	for _, b := range bs {
		ba = append(ba, byte(b))
	}
	return string(ba)
}

func (e *ZeroCoreStructs) InitDefault() error {
	uid, err := uuid.NewV4()
	if err != nil {
		return err
	}
	e.ID = uid.String()
	e.CreateTime = Date(time.Now())
	e.UpdateTime = Date(time.Now())
	if e.Features == nil {
		e.Features = make(map[string]interface{})
	}
	return nil
}

func (e *ZeroCoreStructs) GetJSONFeature() string {
	if e.Features == nil {
		mjson, err := json.Marshal(make(map[string]string))
		if err != nil {
			panic(err)
		}
		return string(mjson)
	} else {
		mjson, err := json.Marshal(e.Features)
		if err != nil {
			panic(err)
		}
		return string(mjson)
	}
}

func (e *ZeroCoreStructs) SetJSONFeature(jsonString string) {
	var jsonMap map[string]interface{}
	err := json.Unmarshal([]byte(jsonString), &jsonMap)
	if err != nil {
		panic(err)
	}
	e.Features = jsonMap
}

func (e *ZeroCoreStructs) LoadRowData(rowmap map[string]interface{}) {
	_, ok := rowmap["id"]
	if ok {
		e.ID = e.UInt8ToString(rowmap["id"].([]uint8))
	}

	_, ok = rowmap["create_time"]
	if ok {
		e.CreateTime = Date(rowmap["create_time"].(time.Time))
	}

	_, ok = rowmap["update_time"]
	if ok {
		e.UpdateTime = Date(rowmap["update_time"].(time.Time))
	}

	_, ok = rowmap["features"]
	if ok {
		e.SetJSONFeature(e.UInt8ToString(rowmap["features"].([]uint8)))
	}
}

func (e *ZeroCoreStructs) String() string {
	mjson, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return string(mjson)
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
