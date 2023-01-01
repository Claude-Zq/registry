// Copyright 2022 CloudWeGo Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/hertz-contrib/registry/consul"
)

var wg sync.WaitGroup

func main() {
	config := consulapi.DefaultConfig()
	config.Address = "127.0.0.1:8500"
	consulClient, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatal(err)
		return
	}
	localIP, err := consul.GetLocalIPv4Address()
	if err != nil {
		log.Fatal(err)
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		addr := fmt.Sprintf("%s:8888", localIP)
		r := consul.NewConsulRegister(consulClient)
		h := server.Default(
			server.WithHostPorts(addr),
			server.WithRegistry(r, &registry.Info{
				ServiceName: "hertz.test.demo",
				Addr:        utils.NewNetAddr("tcp", addr),
				Weight:      10,
				Tags:        nil,
			}),
		)

		h.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
			ctx.JSON(consts.StatusOK, utils.H{"ping": "pong1"})
		})
		h.Spin()
	}()
	go func() {
		defer wg.Done()
		addr := fmt.Sprintf("%s:8889", localIP)
		r := consul.NewConsulRegister(consulClient)
		h := server.Default(
			server.WithHostPorts(addr),
			server.WithRegistry(r, &registry.Info{
				ServiceName: "hertz.test.demo",
				Addr:        utils.NewNetAddr("tcp", addr),
				Weight:      10,
				Tags:        nil,
			}),
		)

		h.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
			ctx.JSON(consts.StatusOK, utils.H{"ping": "pong2"})
		})
		h.Spin()
	}()

	wg.Wait()
}