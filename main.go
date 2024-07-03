/*
Copyright The Ratify Authors.
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
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	kubletConfigPath = "/etc/kubernetes/kubelet.conf"
	configMapName    = "ratify-config"
	filePathPrefix   = ".ratify"
	defaultHomeDir   = "/root"
	certsFolderName  = "certs"
)

var homedir string

func main() {
	logrus.Debugf("using kubeconfig: %s", kubletConfigPath)
	var err error
	homedir, err = os.UserHomeDir()
	if homedir == "" || err != nil {
		logrus.Errorf("Unable to get home directory. Using default path %s", defaultHomeDir)
		homedir = defaultHomeDir
	}
	logrus.Debugf("using home directory: %s", homedir)

	// uses the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubletConfigPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logrus.Errorf("failed to create clientset: %v", err)
		return
	}

	// watch for changes in the ConfigMap
	watchForChanges(clientset, "default")
}

// based on: https://github.com/ScarletTanager/configmap-watcher-example/blob/main/watch/watch.go
// watchForChanges registers a watcher on the ConfigMap and listens for events
func watchForChanges(clientset *kubernetes.Clientset, namespace string) {
	for {
		watcher, err := clientset.CoreV1().ConfigMaps(namespace).Watch(context.TODO(),
			v1.SingleObject(v1.ObjectMeta{Name: configMapName, Namespace: namespace}))
		if err != nil {
			panic("unable to create watcher")
		}
		logrus.Info("watcher created")
		updateConfigMap(watcher.ResultChan())
	}
}

// updateConfigMap updates the files in the directory with the data from the ConfigMap
// Listens for events on the eventChannel and updates the files accordingly
func updateConfigMap(eventChannel <-chan watch.Event) {
	for {
		event, open := <-eventChannel
		if open {
			switch event.Type {
			case watch.Added:
				logrus.Info("configmap added")
				fallthrough
			case watch.Modified:
				logrus.Info("configmap modified. writing to file(s)")
				if modifiedConfigMap, ok := event.Object.(*corev1.ConfigMap); ok {
					if len(modifiedConfigMap.Data) == 0 {
						logrus.Warning("configmap has no data")
						return
					} else {
						for filename, value := range modifiedConfigMap.Data {
							// write value to file with name key
							filePath := fmt.Sprintf("%s/%s/%s", homedir, filePathPrefix, filename)
							if strings.Contains(filename, ".crt") {
								filePath = fmt.Sprintf("%s/%s/%s/%s", homedir, filePathPrefix, certsFolderName, filename)
							}
							logrus.Debugf("writing to file %s", filePath)
							err := writeFile(filePath, value)
							if err != nil {
								logrus.Errorf("failed to write to file: %v", err)
							}
						}
					}
				} else {
					logrus.Error("unable to cast object to ConfigMap")
				}
			case watch.Deleted:
				// delete the files associated with the ConfigMap
				logrus.Info("configmap deleted. deleting file(s)")
				if deletedConfigMap, ok := event.Object.(*corev1.ConfigMap); ok {
					for filename := range deletedConfigMap.Data {
						filePath := fmt.Sprintf("%s/%s/%s", homedir, filePathPrefix, filename)
						logrus.Infof("deleting file %s", filePath)
						err := os.Remove(filePath)
						if err != nil {
							logrus.Errorf("error deleting file %s: %v", filePath, err)
						}
					}
				} else {
					logrus.Error("unable to cast object to ConfigMap")
				}
			default:
				// do nothing
			}
		} else {
			// if eventChannel is closed, it means the server has closed the connection
			return
		}
	}
}

// writeFile writes the contents to the file at filePath
// It creates the directory if it does not exist
func writeFile(filePath string, contents string) error {
	// Ensure the directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	err := os.WriteFile(filePath, []byte(contents), 0644)
	if err != nil {
		logrus.Errorf("failed writing to file %s: %v", filePath, err)
		return err
	}
	logrus.Infof("successfully wrote to file %s", filePath)
	return nil
}
