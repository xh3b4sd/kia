# kia

Opinionated Kubernetes infrastructure automation.



### Project Structure

General non sensitive configuration is stored in the `env` directory.

- `env/def` contains all templates applied to all Kubernetes environments. The
  defaults configured here should reliable work regardless the underlying
  infrastructure provider they are applied to.
- `env/eks` contains all templates applied to the cloud provider AWS. The
  patches configured here should reliable work for EKS on AWS.
- `env/knd` contains all templates applied to local machines running on darwin
  architectures. The patches configured here should reliable work for Kind
  containers.



### Project Installation

```
go get github.com/xh3b4sd/kia
```



### Project Prerequisites

- [aws] CLI for working with AWS and EKS clusters. Make sure to have your
  credentials set up according to the AWS account you want to use if you intend
  to work with EKS clusters.
- [aws-iam-authenticator] for AWS authentication. `eksctl` is using this tool
  under the hood in order to generate kube configs for your EKS cluster.
- [docker] for local setup using `kind`. Make sure to upper the RAM that Docker
  is allowed to use. You will need at least 4GB in order to get the local
  environment using Kind up and running.
- [eksctl] for working with AWS and EKS clusters. `kia` uses that to create EKS
  clusters. A config file is provided when creating the EKS cluster.
- [golang] for installation of `kia` and `red` on your machine. You will also
  need go if you intend to develop the project further.
- [grpcurl] for playing with gRPC apis. Always useful to have when working with
  gRPC so that requests can be made and responses can be inspected.
- [helm] for deploying applications into Kubernetes clusters. Anything deployed
  to a Kubernetes cluster should be packaged as a Helm chart.
- [istioctl] for setting up the service mesh within Kubernetes clusters. Does
  also provide Envoy as ingress controller which proxies grpc-web for frontend
  applications.
- [kind] for working with Kubernetes clusters locally. Necessary to get a
  "production like" environment on your machine. See the `docker` section above
  for eventual RAM issues.
- [kubectl] for working with the Kubernetes api. Used under the hood by `kia`.
  Also always useful to get around the Kubernetes cluster. You want to check
  runtime objects like pods and scale deployments.
- [red] for secret data management. Also see the project configuration section
  below. You may need to understand the file structure of the `sec` repository
  if you want to manage your own secret data.



### Project Configuration

- https://github.com/xh3b4sd/kia as of now needs to be cloned on your machine in
  order to find cluster specific configuration, because `kia` needs these files
  when you execute the command line tool. That is why the location of the
  repository needs to be configured in the kia config file as shown below.
- https://github.com/xh3b4sd/sec is an example repository containing secret
  data. You need to create your own. As of now this repository that you create
  for your own environment needs to be cloned on your machine in order to find
  cluster specific secrets, because `kia` needs these files when you execute the
  command line tool. That is why the location of the repository needs to be
  configured in the kia config file as shown below.

```
$ kia update org -h
Update the current organization name. Managing clusters for different
organizations requires kia to know where its own assets are and where to find
secret data. The latter is managed via the red command line tool. See
https://github.com/xh3b4sd/red for more information. Below is shown the
expected config file location on your file system, including the required
structure and its associated values.

    $ cat ~/.config/kia/config.yaml
    kia: "~/project/xh3b4sd/kia/"
    org:
      list:
        - org: "xh3b4sd"
          sec: "~/project/xh3b4sd/sec/"
        - org: "yourorg"
          sec: "~/project/yourorg/sec/"
      selected: "xh3b4sd"

Given the example config file above the organization used by kia can be
changed as shown below.

    kia update org --selected yourorg

Usage:
  kia update org [flags]

Flags:
  -h, --help              help for org
      --selected string   Select the given organization for current use.
```



### Cluster Creation

```
$ kia create eks -h
Create kubernetes infrastructure environments for eks. The basis for this
type of environment is a remote eks cluster. EKS stands for elastic
kubernetes service. For more information check the eks website.

    https://aws.amazon.com/eks

In order to create and setup the eks cluster we need to properly configure
the kia command line tool. This is done via its config file, tracked on the
local file system. The kia base path must be set. This is the local path of
the kia repository from which general templates are read. Add the following
line to your config file according to your local setup.

    kia: "~/project/xh3b4sd/kia/"

In order to create and setup the eks cluster we need to properly configure
secret data. This is done via a separate private repository containing the
secret data and the red command line tool.

    https://github.com/xh3b4sd/red

The sec base path must be set. This is the local path of the sec repository
from which the red command line tool reads the secret data. Add the following
line to your config file according to your local setup.

    sec: "~/project/xh3b4sd/sec/"

An eks cluster can be created like shown below once the kia config file and
the red command line tool are in place. Cluster creation requires a unique
cluster name. A simple convention could to use the kia prefix and a two digit
number.

    $ kia create eks -c kia02
    { "caller":"github.com/xh3b4sd/kia/cmd/create/eks/runner.go:88", "level":"info", "message":"decrypting local secrets", "time":"2020-10-24 20:42:24" }
    { "caller":"github.com/xh3b4sd/kia/cmd/create/eks/runner.go:125", "level":"info", "message":"creating eks cluster", "time":"2020-10-24 20:42:24" }
    { "caller":"github.com/xh3b4sd/kia/cmd/create/eks/runner.go:145", "level":"info", "message":"installing service mesh", "time":"2020-10-24 21:07:52" }
    { "caller":"github.com/xh3b4sd/kia/cmd/create/eks/runner.go:154", "level":"info", "message":"creating infra namespace", "time":"2020-10-24 21:08:30" }
    { "caller":"github.com/xh3b4sd/kia/cmd/create/eks/runner.go:163", "level":"info", "message":"configure istio injection", "time":"2020-10-24 21:08:32" }
    { "caller":"github.com/xh3b4sd/kia/cmd/create/eks/runner.go:172", "level":"info", "message":"installing infra chart", "time":"2020-10-24 21:08:32" }
    { "caller":"github.com/xh3b4sd/kia/cmd/create/eks/runner.go:198", "level":"info", "message":"installing external-dns chart", "time":"2020-10-24 21:08:35" }
    { "caller":"github.com/xh3b4sd/kia/cmd/create/eks/runner.go:235", "level":"info", "message":"installing cert-manager chart", "time":"2020-10-24 21:08:43" }
    { "caller":"github.com/xh3b4sd/kia/cmd/create/eks/runner.go:267", "level":"info", "message":"installing cert-asset chart", "time":"2020-10-24 21:08:54" }
    { "caller":"github.com/xh3b4sd/kia/cmd/create/eks/runner.go:307", "level":"info", "message":"installing istio-asset chart", "time":"2020-10-24 21:09:09" }

After some time the cluster created as shown above would be available
depending on the Route53 hosted zone configured in the secret data
repository.

    kia02.aws.example.com

Usage:
  kia create eks [flags]

Flags:
  -c, --cluster string   Cluster ID used for AWS and EKS resource naming.
  -h, --help             help for eks
  -k, --kia string       Kia base path on the local file system. (default "~/go/src/github.com/xh3b4sd/kia")
  -r, --region string    Region in which the EKS cluster gets created. (default "eu-central-1")
  -s, --sec string       Sec base path on the local file system. (default "~/projects/xh3b4sd/sec")
```



### Cluster Deletion

```
$ kia delete eks -h
Delete kubernetes infrastructure environments for eks. The deletion process
is mostly straight forward since eks takes care of most of the cloud provider
resources managed in aws. For now there is only one caveat to be aware of. We
use istio gateways and external-dns to register DNS records in Route53. In
order to cleanup the cluster specific DNS records we need to delete the
istio-asset chart first and let external-dns take care of the cleanup
procedure. For now the mechanism is purely time based, which means we just
wait for 5 minutes. This implies the cleanup might fail and we proceed
deleting the cluster regardless, leaving behind Route53 DNS records.

    $ kia delete eks -c kia02
    { "caller":"github.com/xh3b4sd/kia/cmd/delete/eks/runner.go:39", "level":"info", "message":"deleting istio-asset chart", "time":"2020-10-25 13:27:37" }
    { "caller":"github.com/xh3b4sd/kia/cmd/delete/eks/runner.go:55", "level":"info", "message":"waiting for cleanup", "time":"2020-10-25 13:27:39" }
    { "caller":"github.com/xh3b4sd/kia/cmd/delete/eks/runner.go:61", "level":"info", "message":"deleting cert-asset chart", "time":"2020-10-25 13:32:39" }
    { "caller":"github.com/xh3b4sd/kia/cmd/delete/eks/runner.go:70", "level":"info", "message":"deleting eks cluster", "time":"2020-10-25 13:32:42" }

Usage:
  kia delete eks [flags]

Flags:
  -c, --cluster string   Cluster ID used for AWS and EKS resource naming.
  -h, --help             help for eks
```



[aws]: https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-welcome.html
[aws-iam-authenticator]: https://docs.aws.amazon.com/eks/latest/userguide/install-aws-iam-authenticator.html
[docker]: https://www.docker.com/get-started
[eksctl]: https://eksctl.io
[golang]: https://github.com/moovweb/gvm
[grpcurl]: https://github.com/fullstorydev/grpcurl
[helm]: https://helm.sh/docs/intro/install
[istioctl]: https://istio.io/latest/docs/reference/commands/istioctl
[kind]: https://kind.sigs.k8s.io
[kubectl]: https://kubernetes.io/docs/tasks/tools/install-kubectl
[red]: https://github.com/xh3b4sd/red
