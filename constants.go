package staticfile

const (
	StaticfileDependency = "staticfile"
	NginxDependency      = "nginx"

	LayerNameStaticfile = "staticfile"

	StartLoggingContents = `
cat < $APP_ROOT/logs/nginx/access.log &
(>&2 cat) < $APP_ROOT/logs/nginx/error.log &
`

	InitScriptContents = `
# ------------------------------------------------------------------------------------------------
# Copyright 2013 Jordon Bedwell.
# Apache License.
#
# Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
# except  in compliance with the License. You may obtain a copy of the License at:
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software distributed under the
# License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
# either express or implied. See the License for the specific language governing permissions
# and  limitations under the License.
# ------------------------------------------------------------------------------------------------

export LD_LIBRARY_PATH=$APP_ROOT/nginx/lib:$LD_LIBRARY_PATH

mkdir -p $APP_ROOT/logs/nginx
if [[ ! -f $APP_ROOT/logs/nginx/access.log ]]; then
		mkfifo $APP_ROOT/logs/nginx/access.log
fi

if [[ ! -f $APP_ROOT/logs/nginx/error.log ]]; then
    mkfifo $APP_ROOT/logs/nginx/error.log
fi
`
)
