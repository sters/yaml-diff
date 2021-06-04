# yaml-diff

[![go](https://github.com/sters/yaml-diff/workflows/Go/badge.svg)](https://github.com/sters/yaml-diff/actions?query=workflow%3AGo)
[![codecov](https://codecov.io/gh/sters/yaml-diff/branch/main/graph/badge.svg)](https://codecov.io/gh/sters/yaml-diff)
[![go-report](https://goreportcard.com/badge/github.com/sters/yaml-diff)](https://goreportcard.com/report/github.com/sters/yaml-diff)

## Usage

```
go get -u github.com/sters/yaml-diff/cmd/yaml-diff
```
or download from [Releases](https://github.com/sters/yaml-diff/releases).

```
yaml-diff -file1 path/to/foo.yaml -file2 path/to/bar.yaml
```

If the given yaml has a [`---` separated structure](https://yaml.org/spec/1.2/spec.html#id2760395), then the two yaml's will get all the differences in their respective structures. The one with the smallest difference is considered to be the same structure and the difference is displayed.

Also, it's displayed in the form of a golang object, and you won't know the rows and other information about the specified yaml itself for now.

## Example

```text
$ make run-example
go run cmd/yaml-diff/main.go -file1 example/a.yaml -file2 example/b.yaml
  map[string]interface{}{
  	"apiVersion": string("v1"),
  	"kind":       string("Service"),
  	"metadata":   map[string]interface{}{"name": string("my-service")},
  	"spec": map[string]interface{}{
  		"ports": []interface{}{
  			map[string]interface{}{
- 				"port":       uint64(80),
+ 				"port":       uint64(8080),
  				"protocol":   string("TCP"),
  				"targetPort": uint64(9376),
  			},
  		},
  		"selector": map[string]interface{}{"app": string("MyApp")},
  	},
  }

  map[string]interface{}{
  	"apiVersion": string("apps/v1"),
  	"kind":       string("Deployment"),
  	"metadata":   map[string]interface{}{"labels": map[string]interface{}{"app": string("MyApp")}, "name": string("app-deployment")},
  	"spec": map[string]interface{}{
- 		"replicas": uint64(3),
+ 		"replicas": uint64(10),
  		"selector": map[string]interface{}{"matchLabels": map[string]interface{}{"app": string("MyApp")}},
  		"template": map[string]interface{}{
  			"metadata": map[string]interface{}{"labels": map[string]interface{}{"app": string("MyApp")}},
  			"spec": map[string]interface{}{
  				"containers": []interface{}{
  					map[string]interface{}{
- 						"image": string("my-app:1.0.0"),
+ 						"image": string("my-app:1.1.0"),
  						"name":  string("app"),
  						"ports": []interface{}{map[string]interface{}{"containerPort": uint64(9376)}},
  					},
  				},
  			},
  		},
  	},
  }

  map[string]interface{}{
+ 	"bar": []interface{}{string("missing in a.yaml")},
- 	"foo": string("missing-in-b"),
  }

Not found on another one: map[baz:[missing in a.yaml]]...
--------------------
go run cmd/yaml-diff/main.go -file1 example/b.yaml -file2 example/a.yaml
  map[string]interface{}{
  	"apiVersion": string("v1"),
  	"kind":       string("Service"),
  	"metadata":   map[string]interface{}{"name": string("my-service")},
  	"spec": map[string]interface{}{
  		"ports": []interface{}{
  			map[string]interface{}{
- 				"port":       uint64(8080),
+ 				"port":       uint64(80),
  				"protocol":   string("TCP"),
  				"targetPort": uint64(9376),
  			},
  		},
  		"selector": map[string]interface{}{"app": string("MyApp")},
  	},
  }

  map[string]interface{}{
  	"apiVersion": string("apps/v1"),
  	"kind":       string("Deployment"),
  	"metadata":   map[string]interface{}{"labels": map[string]interface{}{"app": string("MyApp")}, "name": string("app-deployment")},
  	"spec": map[string]interface{}{
- 		"replicas": uint64(10),
+ 		"replicas": uint64(3),
  		"selector": map[string]interface{}{"matchLabels": map[string]interface{}{"app": string("MyApp")}},
  		"template": map[string]interface{}{
  			"metadata": map[string]interface{}{"labels": map[string]interface{}{"app": string("MyApp")}},
  			"spec": map[string]interface{}{
  				"containers": []interface{}{
  					map[string]interface{}{
- 						"image": string("my-app:1.1.0"),
+ 						"image": string("my-app:1.0.0"),
  						"name":  string("app"),
  						"ports": []interface{}{map[string]interface{}{"containerPort": uint64(9376)}},
  					},
  				},
  			},
  		},
  	},
  }

  map[string]interface{}{
- 	"bar": []interface{}{string("missing in a.yaml")},
+ 	"foo": string("missing-in-b"),
  }

Not found on another one: map[baz:[missing in a.yaml]]...
```
