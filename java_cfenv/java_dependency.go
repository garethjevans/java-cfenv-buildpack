package java_cfenv

import (
	"fmt"
	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/sherpa"
	"os"
	"path/filepath"
)

/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

type JavaDependency struct {
	LayerContributor libpak.DependencyLayerContributor
	Logger           bard.Logger
}

func NewJavaDependency(dependency libpak.BuildpackDependency, cache libpak.DependencyCache) (JavaDependency, libcnb.BOMEntry) {
	contributor, entry := libpak.NewDependencyLayer(dependency, cache, libcnb.LayerTypes{
		Launch: true,
	})
	return JavaDependency{
		LayerContributor: contributor,
	}, entry
}

func (j JavaDependency) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	j.LayerContributor.Logger = j.Logger

	return j.LayerContributor.Contribute(layer, func(artifact *os.File) (libcnb.Layer, error) {
		j.Logger.Bodyf("Copying to %s", layer.Path)

		file := filepath.Join(layer.Path, filepath.Base(artifact.Name()))
		if err := sherpa.CopyFile(artifact, file); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to copy %s to %s\n%w", artifact.Name(), file, err)
		}

		layer.LaunchEnvironment.Appendf("CLASSPATH", ":", "%s", file)

		return layer, nil
	})

}

func (j JavaDependency) Name() string {
	return j.LayerContributor.LayerName()
}
