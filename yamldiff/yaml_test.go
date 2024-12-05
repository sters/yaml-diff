package yamldiff

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func confirmYaml(t *testing.T, yamlA, yamlB RawYamlList, want string, opts []DoOptionFunc) {
	t.Helper()

	result := "\n"
	for _, diff := range Do(yamlA, yamlB, opts...) {
		result += diff.Dump()
		result += "\n"
	}

	assert.Equal(t, strings.TrimSpace(want), strings.TrimSpace(result))
}

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

	confirmYaml(t, yamlsA, yamlsB, resultOfNoOptions(), []DoOptionFunc{})
	confirmYaml(t, yamlsA, yamlsB, resultOfWithEmpty(), []DoOptionFunc{EmptyAsNull()})
	confirmYaml(t, yamlsA, yamlsB, resultOfWithZero(), []DoOptionFunc{ZeroAsNull()})
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

func TestForSpecificIssues(t *testing.T) {
	tests := map[string]struct {
		yamlA  string
		yamlB  string
		option []DoOptionFunc
		want   string
	}{
		// TODO: It should be extracted as multiline format to understand easily. It needs to refactor (*diff).dump()
		"#52": {
			yamlA: `
data:
  config: |
    logging.a: false
    logging.b: false`,
			yamlB: `
data:
  config: |
    logging.a: false
    logging.c: false`,
			want: `
  data:
-   config: "logging.a: false\nlogging.b: false"
+   config: "logging.a: false\nlogging.c: false"`,
		},
		"#29": {
			yamlA: `
value: |-
  foo
  bar
  baz
  special
    multiline
`,
			yamlB: `
value: "foo\nbar\nbaz\n\
special\n\
\  multiline"
`,
			want: `value: "foo\nbar\nbaz\nspecial\n  multiline"`,
		},
	}

	for key, test := range tests {
		t.Run(key, func(t *testing.T) {
			yamlA, err := Load(test.yamlA)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v", err)
			}

			yamlB, err := Load(test.yamlB)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v", err)
			}

			confirmYaml(t, yamlA, yamlB, test.want, test.option)
		})
	}
}
