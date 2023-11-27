package database

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"

	"github.com/0meet1/zero-framework/global"
)

var (
	serverAddr string
	user       string
	auth       string
)

func InitElasticDatabase() {
	serverAddr = global.StringValue("zero.elastic.serverAddr")
	user = global.StringValue("zero.elastic.user")
	auth = global.StringValue("zero.elastic.auth")
}

type EQueryRequest struct {
	indexName   string
	documentID  string
	deleteQuery bool

	Query interface{} `json:"query,omitempty"`
}

func (query *EQueryRequest) Init(indexName string, documentID string, deleteQuery bool) {
	query.indexName = indexName
	query.documentID = documentID
	query.deleteQuery = deleteQuery
}

func (query *EQueryRequest) InitIndex(indexName string) {
	query.indexName = indexName
	query.documentID = ""
	query.deleteQuery = false
}

func (query *EQueryRequest) Append() error {
	jsonbytes, err := json.Marshal(query.Query)
	if err != nil {
		return err
	}
	request, err := http.NewRequest("POST", "http://"+path.Join(serverAddr, query.indexName, "_doc", query.documentID), bytes.NewBuffer(jsonbytes))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Connection", "close")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("request failed status code : %d", resp.StatusCode))
	}
	return nil
}

func (query *EQueryRequest) Update() error {
	jsonbytes, err := json.Marshal(query.Query)
	if err != nil {
		return err
	}
	request, err := http.NewRequest("PUT", "http://"+path.Join(serverAddr, query.indexName, "_doc", query.documentID), bytes.NewBuffer(jsonbytes))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Connection", "close")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("request failed status code : %d", resp.StatusCode))
	}
	return nil
}

func (query *EQueryRequest) Delete() error {
	jsonbytes, err := json.Marshal(query.Query)
	if err != nil {
		return err
	}
	request, err := http.NewRequest("DELETE", "http://"+path.Join(serverAddr, query.indexName, "_doc", query.documentID), bytes.NewBuffer(jsonbytes))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Connection", "close")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("request failed status code : %d", resp.StatusCode))
	}
	return nil
}

func (query *EQueryRequest) Get() (*EQueryResponse, error) {
	jsonbytes, err := json.Marshal(query.Query)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("GET", "http://"+path.Join(serverAddr, query.indexName, "_doc", query.documentID), bytes.NewBuffer(jsonbytes))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Connection", "close")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.New(fmt.Sprintf("request failed status code : %d", resp.StatusCode))
	}
	qreqp := &EQueryResponse{}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, qreqp.ParserError(resp)
	}
	qreqp.ParserData(resp)
	return qreqp, nil
}

func (query *EQueryRequest) Search() (*EQueryResponse, error) {
	jsonbytes, err := json.Marshal(query.Query)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", "http://"+path.Join(serverAddr, query.indexName, "_search"), bytes.NewBuffer(jsonbytes))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Connection", "close")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	qreqp := &EQueryResponse{}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, qreqp.ParserError(resp)
	}
	qreqp.ParserData(resp)
	return qreqp, nil
}

func (query *EQueryRequest) DeleteByQuery() error {
	jsonbytes, err := json.Marshal(query)
	if err != nil {
		return err
	}
	request, err := http.NewRequest("POST", "http://"+path.Join(serverAddr, query.indexName, "_delete_by_query"), bytes.NewBuffer(jsonbytes))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Connection", "close")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("request failed status code : %d", resp.StatusCode))
	}
	return nil
}

type EQueryResponse struct {
	Datas []interface{}
	Total int
	Error string
}

func (qresp *EQueryResponse) ParserData(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var datas map[string]interface{}
	err = json.Unmarshal(body, &datas)
	if err != nil {
		return err
	}

	hits, ok := datas["hits"]
	if ok {
		qresp.Total = int(hits.(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))
		_hits, ok := hits.(map[string]interface{})["hits"]
		if ok {
			hitDatas := _hits.([]interface{})
			qresp.Datas = make([]interface{}, len(hitDatas))
			for i, hitData := range hitDatas {
				_source := hitData.(map[string]interface{})["_source"].(map[string]interface{})
				delete(_source, "tags")
				qresp.Datas[i] = _source
			}
		}
	}

	_, ok = datas["_source"]
	if ok {
		qresp.Total = 1
		qresp.Datas = make([]interface{}, 1)
		qresp.Datas[0] = datas
	}
	return nil
}

func (qresp *EQueryResponse) ParserError(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	qresp.Error = string(body)
	return errors.New(qresp.Error)
}

type EQuerySearch struct {
	Source         []string      `json:"_source,omitempty"`
	Size           int           `json:"size,omitempty"`
	From           int           `json:"from,omitempty"`
	Sort           []interface{} `json:"sort,omitempty"`
	Query          interface{}   `json:"query,omitempty"`
	TrackTotalHits int           `json:"track_total_hits,omitempty"`
}
