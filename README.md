# Aimpanel v2

The most easy to use game panel as a service

## Swagger

* https://github.com/go-swagger/go-swagger

Swagger allow us to automatically build API docs and client SDK in the future, 
all this without much work from devs

### Generate

This will generate .json file with swagger spec

```bash
./Taskfile.sh swagger-gen
```

### Serve

This will serve UI for generated swagger spec

```bash
./Taskfile.sh swagger-serve
```
