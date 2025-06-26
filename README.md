<div align="center">
  <img src="./.media/logo.png" alt="GoDyno Logo" />
  <br>
  <h3 align="center">Type-safe DynamoDB code generation</h3>
</div>

---

[![Go Version](https://img.shields.io/github/go-mod/go-version/Mad-Pixels/go-dyno?style=flat-square&logo=go&logoColor=white)](https://golang.org/) [![License](https://img.shields.io/github/license/Mad-Pixels/go-dyno?style=flat-square)](LICENSE) [![Latest Release](https://img.shields.io/github/v/release/Mad-Pixels/go-dyno?style=flat-square&logo=github)](https://github.com/Mad-Pixels/go-dyno/releases/latest)

**`GoDyno`** is a cli-tool for generating type-safe Go code from JSON schemas for DynamoDB.

### Idea
**The Pain:** DynamoDB integration in Go often leads to runtime errors, endless string literals, and fragile code that breaks when schemas change. Developers waste time writing boilerplate query builders, managing attribute mappings, and debugging type mismatches that could be caught at compile-time.
**The Solution:** GoDyno eliminates this friction by generating strongly-typed Go code directly from your DynamoDB schema definitions. Write your table schema once in JSON, and get production-ready Go code with full IDE support, automatic index selection, and compile-time safety. No more guessing attribute names or debugging marshaling errors at runtime.

### Why GoDyno?
- **Code Generation:** Produces clean, dependency-free Go code directly into your project.
- **Type Safety:** Ensures compile-time checks and full IDE autocompletion. 
- **Unified Schema:** Maintains a single source of truth—use one JSON schema for both Terraform and Go.
- **Production Ready:** Generates optimized queries with intelligent automatic index selection.
- **AWS SDK v2 Compatibility:** Full support for AWS SDK v2, including handling composite keys.

### Documentation
All our docs placed [here](https://go-dyno.madpixels.io/).

### Contributing
We're open to any new ideas and contributions.
Found a bug? Have an idea? We welcome pull requests and issues.
