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
yaml-diff path/to/foo.yaml path/to/bar.yaml
```

If the given yaml has a [`---` separated structure](https://yaml.org/spec/1.2/spec.html#id2760395), then the two yaml's will get all the differences in their respective structures. The one with the smallest difference is considered to be the same structure and the difference is displayed.

Also, it's displayed in the form of a golang object, and you won't know the rows and other information about the specified yaml itself for now.

## Example

<details><summary>`make run-example`</summary>

```text
go run cmd/yaml-diff/main.go example/a.yaml example/b.yaml
--- example/a.yaml
+++ example/b.yaml

  Map{
  	"this": Map{"is": Map{"the": String("same")}},
  }

  Map{
  	"apiVersion": String("v1"),
  	"kind":       String("Service"),
  	"metadata":   Map{"name": String("my-service")},
  	"spec": Map{
  		"ports": List{
  			Map{
- 				"port":       Number(80),
+ 				"port":       Number(8080),
  				"protocol":   String("TCP"),
  				"targetPort": Number(9376),
  			},
  		},
  		"selector": Map{"app": String("MyApp")},
  	},
  }

  Map{
  	"apiVersion": String("apps/v1"),
  	"kind":       String("Deployment"),
  	"metadata":   Map{"labels": Map{"app": String("MyApp")}, "name": String("app-deployment")},
  	"spec": Map{
- 		"replicas": uNumber(3),
+ 		"replicas": uNumber(10),
  		"selector": Map{"matchLabels": Map{"app": String("MyApp")}},
  		"template": Map{
  			"metadata": Map{"labels": Map{"app": String("MyApp")}},
  			"spec": Map{
  				"containers": List{
  					Map{
- 						"image": String("my-app:1.0.0"),
+ 						"image": String("my-app:1.1.0"),
  						"name":  String("app"),
  						"ports": List{Map{"containerPort": uNumber(9376)}},
  					},
  				},
  			},
  		},
  	},
  }

  map[String]interface{}{
+ 	"bar": List{String("missing in a.yaml")},
- 	"foo": String("missing-in-b"),
  }

+ Map{"baz": List{String("missing in a.yaml")}}

--------------------
go run cmd/yaml-diff/main.go example/b.yaml example/a.yaml
--- example/b.yaml
+++ example/a.yaml

  Map{
  	"this": Map{"is": Map{"the": String("same")}},
  }

  Map{
  	"apiVersion": String("v1"),
  	"kind":       String("Service"),
  	"metadata":   Map{"name": String("my-service")},
  	"spec": Map{
  		"ports": List{
  			Map{
- 				"port":       Number(8080),
+ 				"port":       Number(80),
  				"protocol":   String("TCP"),
  				"targetPort": Number(9376),
  			},
  		},
  		"selector": Map{"app": String("MyApp")},
  	},
  }

  Map{
  	"apiVersion": String("apps/v1"),
  	"kind":       String("Deployment"),
  	"metadata":   Map{"labels": Map{"app": String("MyApp")}, "name": String("app-deployment")},
  	"spec": Map{
- 		"replicas": Number(10),
+ 		"replicas": Number(3),
  		"selector": Map{"matchLabels": Map{"app": String("MyApp")}},
  		"template": Map{
  			"metadata": Map{"labels": Map{"app": String("MyApp")}},
  			"spec": Map{
  				"containers": List{
  					Map{
- 						"image": String("my-app:1.1.0"),
+ 						"image": String("my-app:1.0.0"),
  						"name":  String("app"),
  						"ports": List{Map{"containerPort": Number(9376)}},
  					},
  				},
  			},
  		},
  	},
  }

  map[String]interface{}{
- 	"bar": List{String("missing in a.yaml")},
+ 	"foo": String("missing-in-b"),
  }

- Map{"baz": List{String("missing in a.yaml")}}

```

</details>
