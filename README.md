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

## 3. openapiにエンドポイントを追加してみる

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

## 4. entスキーマをpostgresに反映する

`docker-compose.yml`を作成し、postgresを起動します。

```yaml
version: '3'

services:
  postgresql:
    image: postgres:14.0-alpine
    ports:
      - 5432:5432
    environment:
      POSTGRES_USER: admin
      POSTGRES_PASSWORD: admin
      POSTGRES_DB: db
```

```sh
$ docker-compose up -d --renew-anon-volumes
```

「[versioned-migration](https://entgo.io/ja/docs/versioned-migrations)」を利用することでentスキーマの差分をpostgresのマイグレーションスクリプトとして吐き出すことが出来ます。

entは以下のマイグレーションツールに対応しているようです。

- Atlas
- golang-migrate/migrate
- pressly/goose
- amacneil/dbmate
- Flyway
- Liquibase

今回はgooseを利用してみます。

versioned-migrationを有効化するために、`ent/generate.go`に以下を追記します。

```go
//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/versioned-migration ./schema
```

`ent/generate.go`を変更したので再度`go generate`を実施します。

```sh
$ go generate ./ent 
```

`ent/migrate/main.go`を作成します。

```go
//go:build ignore

package main

import (
	"context"
	"log"
	"os"

	"2022-dev-career-boost-handon/ent/migrate"

	"ariga.io/atlas/sql/sqltool"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
	_ "github.com/lib/pq"
)

func main() {
	ctx := context.Background()
	dir, err := sqltool.NewGooseDir("ent/migrate/migrations")
	if err != nil {
		log.Fatalf("failed creating atlas migration directory: %v", err)
	}
	opts := []schema.MigrateOption{
		schema.WithDir(dir),                         
		schema.WithMigrationMode(schema.ModeReplay), 
		schema.WithDialect(dialect.Postgres),        
	}
	if len(os.Args) != 2 {
		log.Fatalln("migration name is required. Use: 'go run -mod=mod ent/migrate/main.go <name>'")
	}
	err = migrate.NamedDiff(ctx, "postgresql://admin:admin@localhost:5432/db?sslmode=disable", os.Args[1], opts...)
	if err != nil {
		log.Fatalf("failed generating migration file: %v", err)
	}
}
```

`ent/migrate/migrations`ディレクトリを用意し、`ent/migrate/main.go`を実行します。

```sh
$ mkdir ent/migrate/migrations
$ go run -mod=mod ent/migrate/main.go create_organizations_and_users
$ ls ent/migrate/migrations   
20221201131619_create_organizations_and_users.sql       atlas.sum
$ cat ent/migrate/migrations/20221201131619_create_organizations_and_users.sql
-- +goose Up
-- create "organizations" table
CREATE TABLE "organizations" ("id" bigint NOT NULL GENERATED BY DEFAULT AS IDENTITY, "name" character varying NOT NULL DEFAULT 'unknown', PRIMARY KEY ("id"));
-- create "users" table
CREATE TABLE "users" ("id" bigint NOT NULL GENERATED BY DEFAULT AS IDENTITY, "name" character varying NOT NULL DEFAULT 'unknown', PRIMARY KEY ("id"));
-- create "organization_users" table
CREATE TABLE "organization_users" ("organization_id" bigint NOT NULL, "user_id" bigint NOT NULL, PRIMARY KEY ("organization_id", "user_id"), CONSTRAINT "organization_users_organization_id" FOREIGN KEY ("organization_id") REFERENCES "organizations" ("id") ON DELETE CASCADE, CONSTRAINT "organization_users_user_id" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE);

-- +goose Down
-- reverse: create "organization_users" table
DROP TABLE "organization_users";
-- reverse: create "users" table
DROP TABLE "users";
-- reverse: create "organizations" table
DROP TABLE "organizations";
```

マイグレーションスクリプトを適用します。

```sh
$ goose -dir ./ent/migrate/migrations postgres "host=localhost port=5432 user=admin password=admin dbname=db sslmode=disable" up
```

DBにテーブルが作成されていることを確認してみましょう。

```
$ PGPASSWORD=admin docker-compose exec postgresql psql -d db -U admin -c "\dt"
              List of relations
 Schema |        Name        | Type  | Owner 
--------+--------------------+-------+-------
 public | goose_db_version   | table | admin
 public | organization_users | table | admin
 public | organizations      | table | admin
 public | users              | table | admin
(4 rows)

$ PGPASSWORD=admin docker-compose exec postgresql psql -d db -U admin -c "\d users"
                                 Table "public.users"
 Column |       Type        | Collation | Nullable |             Default              
--------+-------------------+-----------+----------+----------------------------------
 id     | bigint            |           | not null | generated by default as identity
 name   | character varying |           | not null | 'unknown'::character varying
Indexes:
    "users_pkey" PRIMARY KEY, btree (id)
Referenced by:
    TABLE "organization_users" CONSTRAINT "organization_users_user_id" FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE

$ PGPASSWORD=admin docker-compose exec postgresql psql -d db -U admin -c "\d organizations"
                             Table "public.organizations"
 Column |       Type        | Collation | Nullable |             Default              
--------+-------------------+-----------+----------+----------------------------------
 id     | bigint            |           | not null | generated by default as identity
 name   | character varying |           | not null | 'unknown'::character varying
Indexes:
    "organizations_pkey" PRIMARY KEY, btree (id)
Referenced by:
    TABLE "organization_users" CONSTRAINT "organization_users_organization_id" FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE
```

## 5. スキーマを変更する

Userに年齢を表す`age`とアクティブかどうかを表す`active`フィールドを追加してみます。

`ent/schema/user.go`の`Fields()`を以下のように変更します。

```go
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Default("unknown"),
		field.Int("age").
			Positive().
			Optional(),
		field.Bool("active").
			Default(true),
	}
}
```

`age`は0以上の整数かつ`Optional()`(任意)にしていることに注意してください。

また、entではフィールドを明示的に`Optional()`としない場合はデフォルトで`Required()`になります([参考](https://entgo.io/ja/docs/schema-fields#%E3%82%AA%E3%83%97%E3%82%B7%E3%83%A7%E3%83%8A%E3%83%AB))。

`go generate`でアセットを生成します。この段階でopenapi.jsonも生成されます。

```sh
$ go generate ./ent 
$ docker run -p 8081:8080 -e SWAGGER_JSON=/openapi.json -v $(pwd)/ent/openapi.json:/openapi.json swaggerapi/swagger-ui
```

swagger-uiで確認するとちゃんとUserに`age`と`active`が追加されています。

![swagger6](./swagger6.png)

次にpostgresにマイグレーションをかけてみます。

```sh
$ go run -mod=mod ent/migrate/main.go add_age_and_active_fields_to_users
$ cat ent/migrate/migrations/20221201155546_add_age_and_active_fields_to_users.sql
-- +goose Up
-- modify "users" table
ALTER TABLE "users" ADD COLUMN "age" bigint NULL, ADD COLUMN "active" boolean NOT NULL DEFAULT true;

-- +goose Down
-- reverse: modify "users" table
ALTER TABLE "users" DROP COLUMN "active", DROP COLUMN "age";
$ goose -dir ./ent/migrate/migrations postgres "host=localhost port=5432 user=admin password=admin dbname=db sslmode=disable" up
```

postgresのusersテーブルを確認してみます。

```sh
PGPASSWORD=admin docker-compose exec postgresql psql -d db -U admin -c "\d users"
                                 Table "public.users"
 Column |       Type        | Collation | Nullable |             Default              
--------+-------------------+-----------+----------+----------------------------------
 id     | bigint            |           | not null | generated by default as identity
 name   | character varying |           | not null | 'unknown'::character varying
 age    | bigint            |           |          | 
 active | boolean           |           | not null | true
Indexes:
    "users_pkey" PRIMARY KEY, btree (id)
Referenced by:
    TABLE "organization_users" CONSTRAINT "organization_users_user_id" FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
```

狙い通りテーブルにColumnが追加されています。

