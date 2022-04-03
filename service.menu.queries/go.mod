module github.com/Resta-Inc/resta/menu.queries

go 1.18

replace github.com/Resta-Inc/resta v0.0.0 => ../

require (
	github.com/Resta-Inc/resta v0.0.0
	github.com/benbjohnson/clock v1.3.0
	github.com/gofrs/uuid v4.2.0+incompatible
	github.com/gorilla/mux v1.8.0
	github.com/lib/pq v1.10.4
	github.com/stretchr/testify v1.7.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.1.1 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
)
