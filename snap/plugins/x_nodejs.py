# -*- Mode:Python; indent-tabs-mode:nil; tab-width:4 -*-
#
# Modified by Anthony Fok on 2018-10-01 to add support for ppc64el and s390x
#
# Copyright (C) 2015-2017 Canonical Ltd
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License version 3 as
# published by the Free Software Foundation.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.

"""The nodejs plugin is useful for node/npm based parts.

The plugin uses node to install dependencies from `package.json`. It
also sets up binaries defined in `package.json` into the `PATH`.

This plugin uses the common plugin keywords as well as those for "sources".
For more information check the 'plugins' topic for the former and the
'sources' topic for the latter.

Additionally, this plugin uses the following plugin-specific keywords:

    - node-packages:
      (list)
      A list of dependencies to fetch using npm.
    - node-engine:
      (string)
      The version of nodejs you want the snap to run on.
    - npm-run:
      (list)
      A list of targets to `npm run`.
      These targets will be run in order, after `npm install`
    - npm-flags:
      (list)
      A list of flags for npm.
    - node-package-manager
      (string; default: npm)
      The language package manager to use to drive installation
      of node packages. Can be either `npm` (default) or `yarn`.
"""

import collections
import contextlib
import json
import logging
import os
import shutil
import subprocess
import sys

import snapcraft
from snapcraft import sources
from snapcraft.file_utils import link_or_copy_tree
from snapcraft.internal import errors

logger = logging.getLogger(__name__)

_NODEJS_BASE = "node-v{version}-linux-{arch}"
_NODEJS_VERSION = "8.12.0"
_NODEJS_TMPL = "https://nodejs.org/dist/v{version}/{base}.tar.gz"
_NODEJS_ARCHES = {"i386": "x86", "amd64": "x64", "armhf": "armv7l", "arm64": "arm64", "ppc64el": "ppc64le", "s390x": "s390x"}
_YARN_URL = "https://yarnpkg.com/latest.tar.gz"


class NodePlugin(snapcraft.BasePlugin):
    @classmethod
    def schema(cls):
        schema = super().schema()

        schema["properties"]["node-packages"] = {
            "type": "array",
            "minitems": 1,
            "uniqueItems": True,
            "items": {"type": "string"},
            "default": [],
        }
        schema["properties"]["node-engine"] = {
            "type": "string",
            "default": _NODEJS_VERSION,
        }
        schema["properties"]["node-package-manager"] = {
            "type": "string",
            "default": "npm",
            "enum": ["npm", "yarn"],
        }
        schema["properties"]["npm-run"] = {
            "type": "array",
            "minitems": 1,
            "uniqueItems": False,
            "items": {"type": "string"},
            "default": [],
        }
        schema["properties"]["npm-flags"] = {
            "type": "array",
            "minitems": 1,
            "uniqueItems": False,
            "items": {"type": "string"},
            "default": [],
        }

        if "required" in schema:
            del schema["required"]

        return schema

    @classmethod
    def get_build_properties(cls):
        # Inform Snapcraft of the properties associated with building. If these
        # change in the YAML Snapcraft will consider the build step dirty.
        return ["node-packages", "npm-run", "npm-flags"]

    @classmethod
    def get_pull_properties(cls):
        # Inform Snapcraft of the properties associated with pulling. If these
        # change in the YAML Snapcraft will consider the build step dirty.
        return ["node-engine", "node-package-manager"]

    @property
    def _nodejs_tar(self):
        if self._nodejs_tar_handle is None:
            self._nodejs_tar_handle = sources.Tar(
                self._nodejs_release_uri, self._npm_dir
            )
        return self._nodejs_tar_handle

    @property
    def _yarn_tar(self):
        if self._yarn_tar_handle is None:
            self._yarn_tar_handle = sources.Tar(_YARN_URL, self._npm_dir)
        return self._yarn_tar_handle

    def __init__(self, name, options, project):
        super().__init__(name, options, project)
        self._source_package_json = os.path.join(
            os.path.abspath(self.options.source), "package.json"
        )
        self._npm_dir = os.path.join(self.partdir, "npm")
        self._manifest = collections.OrderedDict()
        self._nodejs_release_uri = get_nodejs_release(
            self.options.node_engine, self.project.deb_arch
        )
        self._nodejs_tar_handle = None
        self._yarn_tar_handle = None

    def pull(self):
        super().pull()
        os.makedirs(self._npm_dir, exist_ok=True)
        self._nodejs_tar.download()
        if self.options.node_package_manager == "yarn":
            self._yarn_tar.download()
        # do the install in the pull phase to download all dependencies.
        if self.options.node_package_manager == "npm":
            self._npm_install(rootdir=self.sourcedir)
        else:
            self._yarn_install(rootdir=self.sourcedir)

    def clean_pull(self):
        super().clean_pull()

        # Remove the npm directory (if any)
        if os.path.exists(self._npm_dir):
            shutil.rmtree(self._npm_dir)

    def build(self):
        super().build()
        if self.options.node_package_manager == "npm":
            installed_node_packages = self._npm_install(rootdir=self.builddir)
            # Copy the content of the symlink to the build directory
            # LP: #1702661
            modules_dir = os.path.join(self.installdir, "lib", "node_modules")
            _copy_symlinked_content(modules_dir)
        else:
            installed_node_packages = self._yarn_install(rootdir=self.builddir)
            lock_file_path = os.path.join(self.sourcedir, "yarn.lock")
            if os.path.isfile(lock_file_path):
                with open(lock_file_path) as lock_file:
                    self._manifest["yarn-lock-contents"] = lock_file.read()

        self._manifest["node-packages"] = [
            "{}={}".format(name, installed_node_packages[name])
            for name in installed_node_packages
        ]

    def _npm_install(self, rootdir):
        self._nodejs_tar.provision(
            self.installdir, clean_target=False, keep_tarball=True
        )
        npm_cmd = ["npm"] + self.options.npm_flags
        npm_install = npm_cmd + ["--cache-min=Infinity", "install"]
        for pkg in self.options.node_packages:
            self.run(npm_install + ["--global"] + [pkg], cwd=rootdir)
        if os.path.exists(os.path.join(rootdir, "package.json")):
            self.run(npm_install, cwd=rootdir)
            self.run(npm_install + ["--global"], cwd=rootdir)
        for target in self.options.npm_run:
            self.run(npm_cmd + ["run", target], cwd=rootdir)
        return self._get_installed_node_packages("npm", self.installdir)

    def _yarn_install(self, rootdir):
        self._nodejs_tar.provision(
            self.installdir, clean_target=False, keep_tarball=True
        )
        self._yarn_tar.provision(self._npm_dir, clean_target=False, keep_tarball=True)
        yarn_cmd = [os.path.join(self._npm_dir, "bin", "yarn")]
        yarn_cmd.extend(self.options.npm_flags)
        if "http_proxy" in os.environ:
            yarn_cmd.extend(["--proxy", os.environ["http_proxy"]])
        if "https_proxy" in os.environ:
            yarn_cmd.extend(["--https-proxy", os.environ["https_proxy"]])
        flags = []
        if rootdir == self.builddir:
            yarn_add = yarn_cmd + ["global", "add"]
            flags.extend(
                [
                    "--offline",
                    "--prod",
                    "--global-folder",
                    self.installdir,
                    "--prefix",
                    self.installdir,
                ]
            )
        else:
            yarn_add = yarn_cmd + ["add"]
        for pkg in self.options.node_packages:
            self.run(yarn_add + [pkg] + flags, cwd=rootdir)

        # local packages need to be added as if they were remote, we
        # remove the local package.json so `yarn add` doesn't pollute it.
        if os.path.exists(self._source_package_json):
            with contextlib.suppress(FileNotFoundError):
                os.unlink(os.path.join(rootdir, "package.json"))
            shutil.copy(
                self._source_package_json, os.path.join(rootdir, "package.json")
            )
            self.run(yarn_add + ["file:{}".format(rootdir)] + flags, cwd=rootdir)

        # npm run would require to bring back package.json
        if self.options.npm_run and os.path.exists(self._source_package_json):
            # The current package.json is the yarn prefilled one.
            with contextlib.suppress(FileNotFoundError):
                os.unlink(os.path.join(rootdir, "package.json"))
            os.link(self._source_package_json, os.path.join(rootdir, "package.json"))
        for target in self.options.npm_run:
            self.run(
                yarn_cmd + ["run", target],
                cwd=rootdir,
                env=self._build_environment(rootdir),
            )
        return self._get_installed_node_packages("npm", self.installdir)

    def _get_installed_node_packages(self, package_manager, cwd):
        try:
            output = self.run_output(
                [package_manager, "ls", "--global", "--json"], cwd=cwd
            )
        except subprocess.CalledProcessError as error:
            # XXX When dependencies have missing dependencies, an error like
            # this is printed to stderr:
            # npm ERR! peer dep missing: glob@*, required by glob-promise@3.1.0
            # retcode is not 0, which raises an exception.
            output = error.output.decode(sys.getfilesystemencoding()).strip()
        packages = collections.OrderedDict()
        dependencies = json.loads(output, object_pairs_hook=collections.OrderedDict)[
            "dependencies"
        ]
        while dependencies:
            key, value = dependencies.popitem(last=False)
            # XXX Just as above, dependencies without version are the ones
            # missing.
            if "version" in value:
                packages[key] = value["version"]
            if "dependencies" in value:
                dependencies.update(value["dependencies"])
        return packages

    def get_manifest(self):
        return self._manifest

    def _build_environment(self, rootdir):
        env = os.environ.copy()
        if rootdir.endswith("src"):
            hidden_path = os.path.join(rootdir, "node_modules", ".bin")
            if env.get("PATH"):
                new_path = "{}:{}".format(hidden_path, env.get("PATH"))
            else:
                new_path = hidden_path
            env["PATH"] = new_path
        return env


def _get_nodejs_base(node_engine, machine):
    if machine not in _NODEJS_ARCHES:
        raise errors.SnapcraftEnvironmentError(
            "architecture not supported ({})".format(machine)
        )
    return _NODEJS_BASE.format(version=node_engine, arch=_NODEJS_ARCHES[machine])


def get_nodejs_release(node_engine, arch):
    return _NODEJS_TMPL.format(
        version=node_engine, base=_get_nodejs_base(node_engine, arch)
    )


def _copy_symlinked_content(modules_dir):
    """Copy symlinked content.

    When running newer versions of npm, symlinks to the local tree are
    created from the part's installdir to the root of the builddir of the
    part (this only affects some build configurations in some projects)
    which is valid when running from the context of the part but invalid
    as soon as the artifacts migrate across the steps,
    i.e.; stage and prime.

    If modules_dir does not exist we simply return.
    """
    if not os.path.exists(modules_dir):
        return
    modules = [os.path.join(modules_dir, d) for d in os.listdir(modules_dir)]
    symlinks = [l for l in modules if os.path.islink(l)]
    for link_path in symlinks:
        link_target = os.path.realpath(link_path)
        os.unlink(link_path)
        link_or_copy_tree(link_target, link_path)
