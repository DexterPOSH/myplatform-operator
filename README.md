# MyPlatform Operator

Eventual goal is to build the MyPlatform operator in Go which abstracts an internal platform for hosting opnioniated applications.

## Get Started

Install the Operator SDK, it has two components:
- `operator-sdk` - CLI tool and SDK facilitates development of operators
- `operator lifecycle manager` - Facilitates installation, upgrade & RBAC of operators on a cluster

Installation:
- https://sdk.operatorframework.io/docs/installation/

### Creating a project

Run below to use the operator-skd cli to scaffold a project for developing the operator.

```bash
# Create a directory to store the operator
mkdir -p $HOME/projects/myplatform-operator

# switch to the directory created
cd $HOME/projects/myplatform-operator

# Force using Go modules
export GO111MODULE=on

# Run the operator-sdk CLI to scaffold the project structure
operator-sdk init --domain=dexterposh.github.io --repo=github.com/DexterPOSH/myplatform-operator --skip-go-version-check
```

> The operator-sdk init command generates a go.mod file to be used with Go modules. The --repo flag is required when creating a project outside of $GOPATH/src/, because generated files require a valid module path.

### PROJECT file

One important file of note is the PROJECT file. All the next commands we run use the information in this file.

### Manager

Quick look of the `main.go` file shows the code that initializes and runs the Manager. The manager is responsible for registering the scheme for all custom resource API definitions along with running controllers and webhooks.


### Create an API & Controller

```bash
operator-sdk create api --group=dev --version=v1alpha1 --kind=MyPlatform
```

Run below to update the generated code whenever the *_types.go files are modified.

Under the hood below command runs the controller-gen utility to implement `runtime.Object` interface automatically for our type.


```bash
make generate
```

Below would generate the CRDs automatically

```bash
make manifests
```