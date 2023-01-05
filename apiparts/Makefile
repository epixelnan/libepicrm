default: oapi-api.gen.go oapi-server.gen.go oapi-types.gen.go

oapi-api.gen.go: api.yml
	~/go/bin/oapi-codegen -generate spec,skip-prune -o oapi-api.gen.go -package epicrm_apiparts api.yml

oapi-server.gen.go: api.yml
	~/go/bin/oapi-codegen -generate chi-server -o oapi-server.gen.go -package epicrm_apiparts api.yml

oapi-types.gen.go: api.yml
	# skip-prune is important because none of these types are utilized in this file
	~/go/bin/oapi-codegen -generate types,skip-prune -o oapi-types.gen.go -package epicrm_apiparts api.yml

clean:
	rm -f oapi-api.gen.go oapi-server.gen.go oapi-types.gen.go
