# kia

Opinionated kubernetes infrastructure automation.



### Project Structure

General non sensitive configuration is stored in the `env` directory.

* `env/def` contains all templates applied to all kubernetes environments. The
  defaults configured here should reliable work regardless the underlying
  infrastructure provider they are applied to.
* `env/eks` contains all templates applied to the cloud provider AWS. The
  patches configured here should reliable work for EKS on AWS.
* `env/osx` contains all templates applied to local machines running on darwin
  architectures. The patches configured here should reliable work for Kind
  containers.



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
