SHELL = bash
.ONESHELL:
VERSION = 0.0.1
FORMAT = file
OS = linux

clean:
	rm -fr build

.PHONY: build
build: create-package build/java-cfenv-buildpack build/package.toml
	pack buildpack package "garethjevans-java-cfenv-buildpack-$(VERSION).cnb" --config ./build/package.toml --format "$(FORMAT)"

create-package:
	GO111MODULE=on go install github.com/paketo-buildpacks/libpak/cmd/create-package

.PHONY: build/java-cfenv-buildpack
build/java-cfenv-buildpack: create-package
	create-package --destination ./build/java-cfenv-buildpack --version "0.0.1"

.PHONY: build/package.toml
build/package.toml:
	./scripts/create-package-config.sh ./build/package.toml ./java-cfenv-buildpack "$(OS)"
	cat ./build/package.toml
