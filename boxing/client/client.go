/*
Copyright 2022 QuanxiangCloud Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
     http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"
	"time"

	"github.com/quanxiang-cloud/cabin/boxing/header"
	"github.com/quanxiang-cloud/cabin/boxing/resp"
	e "github.com/quanxiang-cloud/cabin/error"
)

// Config client config
type Config struct {
	Timeout      time.Duration
	MaxIdleConns int
}

// New new a http client
func New(conf Config) http.Client {
	return http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(conf.Timeout * time.Second)
				c, err := net.DialTimeout(netw, addr, time.Second*conf.Timeout)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
			MaxIdleConns: conf.MaxIdleConns,
		},
	}
}

// POST http post
func POST(ctx context.Context, client *http.Client, uri string, params interface{}, entity interface{}) error {
	if reflect.ValueOf(entity).Kind() != reflect.Ptr {
		return errors.New("the entity type must be a pointer")
	}

	paramByte, err := json.Marshal(params)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(paramByte)
	req, err := http.NewRequest("POST", uri, reader)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add(header.GetRequestIDKV(ctx).Wreck())
	req.Header.Add(header.GetTimezone(ctx).Wreck())

	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("expected state value is 200, actually %d", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return decomposeBody(body, entity)
}

func decomposeBody(body []byte, entity interface{}) error {
	r := new(resp.Resp)
	r.Data = entity

	err := json.Unmarshal(body, r)
	if err != nil {
		return err
	}

	if r.Code != e.Success {
		return r.Error
	}

	return nil
}
