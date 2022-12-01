# 2022-dev-career-boost-handon

## 初期セットアップ

```sh
$ go mod init 2022-dev-career-boost-handon
$ go get -d entgo.io/ent/cmd/ent
```

## 1. Organization, Userスキーマを作成する

OrganizationとUserのスキーマを作成してみます。

```sh
$ go run -mod=mod entgo.io/ent/cmd/ent init Organization User
$ tree
.
├── README.md
├── ent
│   ├── generate.go
│   └── schema
│       ├── organization.go
│       └── user.go
├── go.mod
└── go.sum
```

作成されたテンプレートの中身はこんな感じ。

```sh
$ cat ent/schema/organization.go 
package schema

import "entgo.io/ent"

// Organization holds the schema definition for the Organization entity.
type Organization struct {
        ent.Schema
}

// Fields of the Organization.
func (Organization) Fields() []ent.Field {
        return nil
}

// Edges of the Organization.
func (Organization) Edges() []ent.Edge {
        return nil
}
```

Organizationに`name`フィールドを追加してみましょう。

```go
func (Organization) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Default("unknown"),
	}
}
```

次にOrganizationは複数のUserを持つこととしてみます。
entではEdgesとして表現されます。

```go
func (Organization) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("users", User.Type),
	}
}
```

同様にUserにも`name`フィールドとEdgeを追加してみましょう。
Userは必ずOrganizationに所属することとして、EdgeはRequired()されていることに注意してください。

```go
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Default("unknown"),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("organizations", Organization.Type).Ref("users").Required(),
	}
}
```

定義したスキーマに基づいてアセットを生成してみます。

```sh
$ go generate ./ent
$ tree
.
├── README.md
├── ent
│   ├── client.go
│   ├── config.go
│   ├── context.go
│   ├── ent.go
│   ├── enttest
│   │   └── enttest.go
│   ├── generate.go
│   ├── hook
│   │   └── hook.go
│   ├── migrate
│   │   ├── migrate.go
│   │   └── schema.go
│   ├── mutation.go
│   ├── organization
│   │   ├── organization.go
│   │   └── where.go
│   ├── organization.go
│   ├── organization_create.go
│   ├── organization_delete.go
│   ├── organization_query.go
│   ├── organization_update.go
│   ├── predicate
│   │   └── predicate.go
│   ├── runtime
│   │   └── runtime.go
│   ├── runtime.go
│   ├── schema
│   │   ├── organization.go
│   │   └── user.go
│   ├── tx.go
│   ├── user
│   │   ├── user.go
│   │   └── where.go
│   ├── user.go
|   ...
├── go.mod
└── go.sum
```

[公式ドキュメント](https://entgo.io/ja/docs/code-gen/#%E3%82%A2%E3%82%BB%E3%83%83%E3%83%88%E3%81%AE%E7%94%9F%E6%88%90)によると、このとき以下のアセットが生成されています。

- Client と Tx オブジェクトはグラフとのやり取りに使用されます。
- 各スキーマ型の CRUD ビルダー。 
- 各スキーマタイプのエンティティオブジェクト。
- ビルダーとの相互作用に使用される定数と述語を含むパッケージ。
- SQLの migrate パッケージ。 
- mutationミドルウェアを追加するためのhook パッケージ。

entvizを用いると、作成したスキーマを可視化することができます。

```sh
$ go get github.com/hedwigz/entviz/cmd/entviz
$ go run github.com/hedwigz/entviz/cmd/entviz ./ent/schema
$ open ./schema-viz.html
```

![entviz1](./entviz1.png)

## 2. openapiを自動生成する

公式ドキュメントの(「Ent と ogen で REST CRUD APIを自動生成」)[https://entgo.io/ja/blog/2022/02/15/generate-rest-crud-with-ent-and-ogen/]の手順に従います。

`ent/entc.go`ファイルを作成し、以下を記述します。

```go
//go:build ignore

package main

import (
    "log"

    "ariga.io/ogent"
    "entgo.io/contrib/entoas"
    "entgo.io/ent/entc"
    "entgo.io/ent/entc/gen"
    "github.com/ogen-go/ogen"
)

func main() {
    spec := new(ogen.Spec)
    oas, err := entoas.NewExtension(entoas.Spec(spec))
    if err != nil {
        log.Fatalf("creating entoas extension: %v", err)
    }
    ogent, err := ogent.NewExtension(spec)
    if err != nil {
        log.Fatalf("creating ogent extension: %v", err)
    }
    err = entc.Generate("./schema", &gen.Config{}, entc.Extensions(ogent, oas))
    if err != nil {
        log.Fatalf("running ent codegen: %v", err)
    }
}
```

`ent/generate.go`を編集し以下を追記します。

```go
//go:generate go run -mod=mod entc.go
```

`go generate`を実行するとopenapiが生成されます。

```sh
$ go generate ./ent
$ cat ent/openapi.json
```

swagger-uiを利用してopenapiをビジュアライズしてみます。

```sh
$ docker run -p 8081:8080 -e SWAGGER_JSON=/openapi.json -v $(pwd)/ent/openapi.json:/openapi.json swaggerapi/swagger-ui
```

http://localhost:8081 をブラウザで開くとCRUD用APIの定義やスキーマが生成されていることがわかります。

![swagger1](./swagger1.png)

![swagger2](./swagger2.png)

![swagger3](./swagger3.png)

![swagger4](./swagger4.png)

## 3. エンドポイントを追加してみる

`ent/entc.go`を書き換えることで生成されるopenapiに手を加えることができます。

`PATCH /organizations/{id}/addUser`を追加してみましょう。

`ent/entc.go`を以下のように変更してみます。

```go
...
func main() {
	oas, err := entoas.NewExtension(
		entoas.Spec(spec),
		entoas.Mutations(func(_ *gen.Graph, spec *ogen.Spec) error {
			spec.AddPathItem("/organizations/{id}/addUser", ogen.NewPathItem().
				SetPatch(ogen.NewOperation().
					SetOperationID("addUser").
					AddTags("Organization").
					AddResponse("200", ogen.NewResponse()),
				).
				AddParameters(
					ogen.NewParameter().
						InPath().
						SetName("id").
						SetRequired(true).
						SetSchema(ogen.Int()),
					ogen.NewParameter().
						InQuery().
						SetName("user_id").
						SetRequired(true).
						SetSchema(ogen.Int()),
				),
			)
			return nil
		}),
	)
...
```

openapiを生成して再度結果を見てみます。

```sh
$ go generate ./ent
$ docker run -p 8081:8080 -e SWAGGER_JSON=/openapi.json -v $(pwd)/ent/openapi.json:/openapi.json swaggerapi/swagger-ui
```

![swagger5](./swagger5.png)