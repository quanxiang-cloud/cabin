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
package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Hosts      []string
	Direct     bool
	Credential struct {
		AuthMechanism           string
		AuthMechanismProperties map[string]string
		AuthSource              string
		Username                string
		Password                string
		PasswordSet             bool
	}
}

func New(conf *Config) (*mongo.Client, error) {
	opts := options.Client().
		SetDirect(conf.Direct).
		SetHosts(conf.Hosts).SetAuth(options.Credential{
		AuthMechanism:           conf.Credential.AuthMechanism,
		AuthMechanismProperties: conf.Credential.AuthMechanismProperties,
		AuthSource:              conf.Credential.AuthSource,
		Username:                conf.Credential.Username,
		Password:                conf.Credential.Password,
		PasswordSet:             conf.Credential.PasswordSet,
	})

	return mongo.Connect(context.TODO(), opts)
}
