package admincli

import (
	m "raft/admin/src/internal/admin_cli/model"
  clusterform "raft/admin/src/internal/admin_cli/model/cluster_form"
)

type Cli interface {
  Start() (map[string]string, error)
}

type cli struct {}

func NewCli() *cli {
  return new(cli)
}

func (this *cli) Start() (map[string]string, error) {
  var text_input m.Model = clusterform.NewClusterForm(clusterform.ADD)
  return text_input.Show()
}



