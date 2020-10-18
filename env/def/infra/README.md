# infra

This helm chart is used to install general infrastructure related resources into
kubernetes clusters. One example is to ensure a pull secret so that containers
can authenticate against some container registry when pulling their desired
container images. Usually a pull secret contains docker credentials.
