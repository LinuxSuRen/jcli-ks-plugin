package pipeline

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	jCLI "github.com/jenkins-zh/jenkins-cli/client"
	"github.com/linuxsuren/jcli-ks-plugin/cmd/common"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path"
	"strings"
)

func newBackupCommand() (cmd *cobra.Command) {
	opt := &backupOption{}

	cmd = &cobra.Command{
		Use:     "backup",
		Aliases: []string{"b"},
		Short:   "Backup KubeSphere Pipeline from Jenkins",
		Long: `Backup Ku beSphere Pipeline from Jenkins
It only backups the config of jobs
This command rely on kubectl`,
		PreRunE: opt.preRunE,
		RunE:    opt.runE,
	}

	flags := cmd.Flags()
	flags.StringVarP(&opt.output, "output", "o", "", "The output directory which store the backup files")
	flags.StringArrayVarP(&opt.jobs, "jobs", "j", []string{},
		"Which jobs do you want to backup")
	return
}

type backupOption struct {
	jobs   []string
	output string
}

func (o *backupOption) preRunE(cmd *cobra.Command, args []string) (err error) {
	if len(o.jobs) == 0 {
		// search jobs from Jenkins, then let users to choose
		jclient := &jCLI.JobClient{
			JenkinsCore: jCLI.JenkinsCore{
				RoundTripper: nil,
			},
		}
		if _, err = getCurrentJenkinsAndClient(&(jclient.JenkinsCore)); err != nil {
			return
		}

		var keyword string
		if len(args) > 0 {
			keyword = args[0]
		}

		var items []jCLI.JenkinsItem
		var allJobs []string
		if items, err = jclient.SearchViaBlue(keyword, 0, 1000); err == nil {
			allJobs = make([]string, len(items))

			for i, item := range items {
				allJobs[i] = item.FullName
			}
		}

		prompt := &survey.MultiSelect{
			Message: "Please select the pipelines that you want to backup:",
			Options: allJobs,
		}
		if err = survey.AskOne(prompt, &o.jobs); err == nil {
			for i, item := range o.jobs {
				o.jobs[i] = strings.Join(strings.Split("/"+item, "/"), "/job/")
			}
		}
	}

	if o.output == "" {
		o.output = path.Dir(".")
	}
	return
}

func (o *backupOption) runE(_ *cobra.Command, _ []string) (err error) {
	var podName string
	if podName, err = getJenkinsPodName(); err != nil {
		err = fmt.Errorf("cannot get ks-jenkins pod, %v", err)
		return
	}

	for _, job := range o.jobs {
		job = strings.TrimPrefix(job, "/")
		job = strings.TrimSuffix(job, "/")
		job = strings.ReplaceAll(job, "job/", "jobs/")

		if err = os.MkdirAll(path.Dir(podName), 0666); err != nil {
			err = fmt.Errorf("cannot mkdir %s, %v", podName, err)
			return
		}

		cmd := fmt.Sprintf("kubectl cp -n kubesphere-devops-system %s:var/jenkins_home/%s/config.xml %s/config.xml", podName, job, job)
		fmt.Println("start to backup", fmt.Sprintf("/var/jenkins_home/%s/config.xml", job))
		err = common.ExecCommand("kubectl", strings.Split(cmd, " ")[1:]...)
		if err != nil {
			return
		}
	}
	return
}

func getJenkinsPodName() (name string, err error) {
	var data []byte
	cmd := exec.Command("kubectl", strings.Split("-n kubesphere-devops-system get pod -l app=ks-jenkins -o custom-columns=NAME:.metadata.name --no-headers=true", " ")...)
	if data, err = cmd.Output(); err == nil {
		name = strings.TrimSpace(string(data))
	}
	return
}
