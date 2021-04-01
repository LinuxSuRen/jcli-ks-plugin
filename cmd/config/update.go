package config

import (
	"context"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/linuxsuren/jcli-ks-plugin/cmd/common"
	kstypes "github.com/linuxsuren/ks/kubectl-plugin/types"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

func newUpdateCmd() (cmd *cobra.Command) {
	opt := &updateOption{}

	cmd = &cobra.Command{
		Use:     "update",
		Aliases: []string{"up"},
		Short:   "Update the config item",
		PreRunE: opt.preRunE,
		RunE:    opt.runE,
	}
	return
}

type updateOption struct {
	name  string
	token string
}

// kubeSphereConfig is the config object of KubeSphere
// currently it's partial
type kubeSphereConfig struct {
	DevOps struct {
		Password string `yaml:"password"`
	} `yaml:"devops"`
}

func (o *updateOption) preRunE(_ *cobra.Command, args []string) (err error) {
	if len(args) > 0 {
		o.name = args[0]
	}

	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	var config *rest.Config
	var client dynamic.Interface

	if config, err = clientcmd.BuildConfigFromFlags("", kubeconfig); err != nil {
		fmt.Println(err)
	} else {
		if client, err = dynamic.NewForConfig(config); err != nil {
			return
		}
	}

	ctx := context.TODO()
	var rawConfigMap *unstructured.Unstructured
	if rawConfigMap, err = client.Resource(kstypes.GetConfigMapSchema()).Namespace("kubesphere-system").
		Get(ctx, "kubesphere-config", metav1.GetOptions{}); err == nil {
		data := rawConfigMap.Object["data"]
		dataMap := data.(map[string]interface{})
		configStr := dataMap["kubesphere.yaml"].(string)

		ksCfg := kubeSphereConfig{}
		if err = yaml.Unmarshal([]byte(configStr), &ksCfg); err == nil {
			o.token = ksCfg.DevOps.Password
		} else {
			err = fmt.Errorf("unable to parse KubeSphere config from configmap, error: %v", err)
		}
	}

	if o.token == "" {
		err = fmt.Errorf("unable to get Jenkins token from KubeSphere, please check if it was enabled")
	}
	return
}

func (o *updateOption) runE(_ *cobra.Command, args []string) (err error) {
	err = common.ExecCommand("jcli", "config", "update", "--token", o.token, o.name)
	return
}
