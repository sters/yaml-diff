# yaml-diff

[![go](https://github.com/sters/yaml-diff/workflows/Go/badge.svg)](https://github.com/sters/yaml-diff/actions?query=workflow%3AGo)
[![codecov](https://codecov.io/gh/sters/yaml-diff/branch/main/graph/badge.svg)](https://codecov.io/gh/sters/yaml-diff)
[![go-report](https://goreportcard.com/badge/github.com/sters/yaml-diff)](https://goreportcard.com/report/github.com/sters/yaml-diff)

## Usage

```
go install github.com/sters/yaml-diff/cmd/yaml-diff@latest
```
or download from [Releases](https://github.com/sters/yaml-diff/releases).

```
yaml-diff path/to/foo.yaml path/to/bar.yaml
```

If the given yaml has a [`---` separated structure](https://yaml.org/spec/1.2/spec.html#id2760395), then the two yaml's will get all the differences in their respective structures. The one with the smallest difference is considered to be the same structure and the difference is displayed.

The result structure is the same as based or target yaml but format (includes map fields order) is different.

## Example

<details><summary>You can try example directory.</summary>

```text
$ go run cmd/yaml-diff/main.go example/a.yaml example/b.yaml
--- example/a.yaml
+++ example/b.yaml

  apiVersion: "v1"
  kind: "Service"
  metadata:
    name: "my-service"
  spec:
    selector:
      app: "MyApp"
    ports:
      -
        protocol: "TCP"
-       port: 80
+       port: 8080
        targetPort: 9376

  apiVersion: "apps/v1"
  kind: "Deployment"
  metadata:
    name: "app-deployment"
    labels:
      app: "MyApp"
  spec:
-   replicas: 3
+   replicas: 10
    selector:
      matchLabels:
        app: "MyApp"
    template:
      metadata:
        labels:
          app: "MyApp"
      spec:
        containers:
          -
            name: "app"
-           image: "my-app:1.0.0"
+           image: "my-app:1.1.0"
            ports:
              -
                containerPort: 9376

- foo: "missing-in-b"

  this:
    is:
      the: "same"

+ bar:
+   - "missing in a.yaml"

+ baz:
+   - "missing in a.yaml"
```

Even if it reverse order, it also worked properly.

```text
$ go run cmd/yaml-diff/main.go example/b.yaml example/a.yaml
--- example/b.yaml
+++ example/a.yaml

  metadata:
    name: "my-service"
  spec:
    ports:
      -
        protocol: "TCP"
        targetPort: 9376
-       port: 8080
+       port: 80
    selector:
      app: "MyApp"
  apiVersion: "v1"
  kind: "Service"

  apiVersion: "apps/v1"
  kind: "Deployment"
  spec:
    selector:
      matchLabels:
        app: "MyApp"
-   replicas: 10
+   replicas: 3
    template:
      metadata:
        labels:
          app: "MyApp"
      spec:
        containers:
          -
            name: "app"
            ports:
              -
                containerPort: 9376
-           image: "my-app:1.1.0"
+           image: "my-app:1.0.0"
  metadata:
    name: "app-deployment"
    labels:
      app: "MyApp"

- bar:
-   - "missing in a.yaml"

- baz:
-   - "missing in a.yaml"

  this:
    is:
      the: "same"

+ foo: "missing-in-b"
```

</details>
