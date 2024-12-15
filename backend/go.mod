// defines my Go project's module path and dependencies, like a package.json
// go.sum contains cryptographic checksums of specific versions of dependencies, ensuring reproducible builds. like a package-lock.json
module go-fitsync/backend

go 1.22.2

require github.com/lib/pq v1.10.9