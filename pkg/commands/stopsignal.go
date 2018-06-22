/*
Copyright 2018 Google LLC

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

package commands

import (
	"strings"

	"github.com/GoogleContainerTools/kaniko/pkg/dockerfile"
	"github.com/GoogleContainerTools/kaniko/pkg/util"
	"github.com/docker/docker/builder/dockerfile/instructions"
	"github.com/docker/docker/pkg/signal"
	"github.com/google/go-containerregistry/pkg/v1"
	"github.com/sirupsen/logrus"
)

type StopSignalCommand struct {
	cmd *instructions.StopSignalCommand
}

// ExecuteCommand handles command processing similar to CMD and RUN,
func (s *StopSignalCommand) ExecuteCommand(config *v1.Config, buildArgs *dockerfile.BuildArgs) error {
	logrus.Info("cmd: STOPSIGNAL")

	// resolve possible environment variables
	replacementEnvs := buildArgs.ReplacementEnvs(config.Env)
	resolvedEnvs, err := util.ResolveEnvironmentReplacementList([]string{s.cmd.Signal}, replacementEnvs, false)
	if err != nil {
		return err
	}
	stopsignal := resolvedEnvs[0]

	// validate stopsignal
	_, err = signal.ParseSignal(stopsignal)
	if err != nil {
		return err
	}

	logrus.Infof("Replacing StopSignal in config with %v", stopsignal)
	config.StopSignal = stopsignal
	return nil
}

// FilesToSnapshot returns an empty array since this is a metadata command
func (s *StopSignalCommand) FilesToSnapshot() []string {
	return []string{}
}

// CreatedBy returns some information about the command for the image config history
func (s *StopSignalCommand) CreatedBy() string {
	entrypoint := []string{"STOPSIGNAL"}

	return strings.Join(append(entrypoint, s.cmd.Signal), " ")
}