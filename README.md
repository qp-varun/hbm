# HBM (Harbormaster)

Harbormaster is a basic extendable Docker Engine [access authorization plugin](https://docs.docker.com/engine/extend/plugins_authorization/) that runs on directly on the host.

By default, Harbormaster plugin prevents from executing commands with certain parameters.

 1. Docker commands
 2. Pull images
 3. Start containers with specific parameters

* `--privileged`
* `--ipc=host`
* `--net=host`
* `--pid=host`
* `--userns=host`
* `--uts=host`
* any Linux capabilities with parameter `--cap-add=[]`
* any devices added with parameter `--device=[]`
* any dns servers added with parameter `--dns`
* any ports added with parameter `--port`
* any volumes mounted with parameter `-v`
* any logging with parameters `--log-driver` and `--log-opt`
* `--sysctl`
* `--security-opt`

## Versions

Supported Docker versions with HBM.

| HBM Version | Docker Version | Docker API |
|-------------|----------------|------------|
| 0.19.x      | 27.x           | 1.47    |

## Open Source Licenses

We depend on the many great open source licenses, listed below:

* [Docker types and client](github.com/moby/moby) by Docker. View [License](https://github.com/moby/moby/blob/master/LICENSE)
* [Docker plugins helpers](github.com/docker/go-plugins-helpers) by Docker. View [License](https://github.com/docker/go-plugins-helpers/blob/master/LICENSE)
* [Gorm](github.com/jinzhu/gorm) by Jinzhu. View [License](https://github.com/jinzhu/gorm/blob/master/License)
* [Cobra](github.com/spf13/cobra) by Steve Francia. View [License](https://github.com/spf13/cobra/blob/master/LICENSE.txt)
