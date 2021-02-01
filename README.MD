# Virtual Assistant Voice Gateway - GO POC

---
[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-blue)](https://golang.org/)
[![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](../../graphs/commit-activity)

VA Voice GW POC written in GO. Main purpose is to compare with Rust based version in terms of:

- performance (CPU, throughput)
- memory consumption
- Last but not least: SW design (ease of coding, maintainability, etc.)

## How to format the code
```
go fmt github.com/va-voice-gateway
```
## Go design patterns
- http://blog.ralch.com/tutorial/design-patterns/golang-singleton/
- https://www.sohamkamani.com/golang/2018-06-20-golang-factory-patterns/

## To check
- [ ] Best way how to write constructor method (return type X vs *X)
- [ ] Custom features + conditional compilation based on them. Possible?
- [ ] Best way how to represent nullable values(Rust Option<T>). https://stackoverflow.com/questions/9993178/is-there-a-nice-way-to-simulate-a-maybe-or-option-type-in-go
- [ ] How to run unit tests in all folder (instead of go test -v ./nlp)