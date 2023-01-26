package yamldiff

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	yamlA := `
apiVersion: v1
kind: Service
metadata:
  name: my-service
spec:
  selector:
    app: MyApp
  ports:
  - protocol: TCP
    port: 80
    targetPort: 9376
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app-deployment
  labels:
    app: MyApp
spec:
  replicas: 3
  selector:
    matchLabels:
      app: MyApp
  template:
    metadata:
      labels:
        app: MyApp
    spec:
      containers:
      - name: app
        image: my-app:1.0.0
        ports:
        - containerPort: 9376
---
foo: missing-in-b
---
this:
  is:
    the: same
    empty:
---
someStr: foo
zeroStr: ""
someInt: 5
zeroInt: 0
differs: fromA
`
	yamlB := `
metadata:
  name: my-service
spec:
  ports:
  - protocol: TCP
    targetPort: 9376
    port: 8080
  selector:
    app: MyApp
apiVersion: v1
kind: Service
---
apiVersion: apps/v1
kind: Deployment
spec:
  selector:
    matchLabels:
      app: MyApp
  replicas: 10
  template:
    metadata:
      labels:
        app: MyApp
    spec:
      containers:
      - name: app
        ports:
        - containerPort: 9376
        image: my-app:1.1.0
metadata:
  name: app-deployment
  labels:
    app: MyApp
---
bar:
  - missing in a.yaml
---
baz:
  - missing in a.yaml
---
this:
  is:
    the: same
---
someStr: foo
someInt: 5
differs: fromB
`

	yamlsA, err := Load(yamlA)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
	}

	yamlsB, err := Load(yamlB)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
	}

	confirm := func(expect string, opts []DoOptionFunc) {
		result := "\n"
		for _, diff := range Do(yamlsA, yamlsB, opts...) {
			result += diff.Dump()
			result += "\n"
		}

		assert.Equal(t, expect, result)
	}

	confirm(resultOfNoOptions(), []DoOptionFunc{})
	confirm(resultOfWithEmpty(), []DoOptionFunc{EmptyAsNull()})
	confirm(resultOfWithZero(), []DoOptionFunc{ZeroAsNull()})
}

func resultOfNoOptions() string {
	return `
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
-     empty:

  someStr: "foo"
- zeroStr: ""
  someInt: 5
- zeroInt: 0
- differs: "fromA"
+ differs: "fromB"

+ bar:
+   - "missing in a.yaml"

+ baz:
+   - "missing in a.yaml"

`
}

func resultOfWithEmpty() string {
	return `
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
      empty:

  someStr: "foo"
- zeroStr: ""
  someInt: 5
- zeroInt: 0
- differs: "fromA"
+ differs: "fromB"

+ bar:
+   - "missing in a.yaml"

+ baz:
+   - "missing in a.yaml"

`
}

func resultOfWithZero() string {
	return `
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
-     empty:

  someStr: "foo"
  zeroStr: ""
  someInt: 5
  zeroInt: 0
- differs: "fromA"
+ differs: "fromB"

+ bar:
+   - "missing in a.yaml"

+ baz:
+   - "missing in a.yaml"

`
}
