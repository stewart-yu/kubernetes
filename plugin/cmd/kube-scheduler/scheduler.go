/*
Copyright 2014 The Kubernetes Authors.

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

package main

import (
	goflag "flag"
	"os"

	"github.com/spf13/pflag"

	utilflag "k8s.io/apiserver/pkg/util/flag"
	"k8s.io/apiserver/pkg/util/logs"
	_ "k8s.io/kubernetes/pkg/client/metrics/prometheus" // for client metric registration
	_ "k8s.io/kubernetes/pkg/version/prometheus"        // for version metric registration
	"k8s.io/kubernetes/plugin/cmd/kube-scheduler/app"
	"k8s.io/kubernetes/staging/src/k8s.io/apiserver/pkg/util/flag"
	"k8s.io/kubernetes/pkg/version/verflag"
	"github.com/golang/glog"
)

func main() {

	// scheduler 启动所需要配置信息的对象
	s := options.NewSchedulerServer()
	// 解析命令行的参数，对结构体中的内容进行赋值
	s.AddFlags(pflag.CommandLine)

	flag.InitFlags()
	logs.InitLogs()
	defer logs.FlushLogs()

	verflag.PrintAndExitIfRequested()

	// app.Runs(s) 根据配置信息构建出来各种实例，然后运行 scheduler 的核心逻辑，
	// 这个函数会一直运行，不会退出。
	if err := app.Run(s); err != nil {
		glog.Fatalf("scheduler app failed to run: %v", err)

	}
}
