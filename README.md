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
