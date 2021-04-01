package pipeline

import (
	"fmt"
	"github.com/linuxsuren/jcli-ks-plugin/cmd/common"
	"github.com/spf13/cobra"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func newRestoreCommand() (cmd *cobra.Command) {
	opt := &restoreOption{}

	cmd = &cobra.Command{
		Use:   "restore",
		Short: "Restore KubeSphere Pipeline to Jenkins",
		Long: `Restore KubeSphere Pipeline to Jenkins
It only restore the config of jobs
This command rely on kubectl`,
		PreRunE: opt.preRunE,
		RunE:    opt.runE,
	}

	flags := cmd.Flags()
	flags.StringVarP(&opt.input, "input", "i", "", "The input directory which store the backup files")
	return
}

type restoreOption struct {
	input string
}

func (o *restoreOption) preRunE(cmd *cobra.Command, args []string) (err error) {
	if o.input == "" {
		o.input = path.Dir(".")
	}
	return
}

func (o *restoreOption) runE(cmd *cobra.Command, args []string) (err error) {
	var podName string
	if podName, err = getJenkinsPodName(); err != nil {
		err = fmt.Errorf("cannot get ks-jenkins pod, %v", err)
		return
	}

	err = filepath.Walk(fmt.Sprintf("%s/jobs", o.input), func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, "config.xml") {
			cmd := fmt.Sprintf("kubectl cp -n kubesphere-devops-system %s %s:var/jenkins_home/%s", path, podName, path)
			fmt.Println("start to restore", path)
			err = common.ExecCommand("kubectl", strings.Split(cmd, " ")[1:]...)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return
}
