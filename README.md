# Config Module for Tinh Tinh

<div align="center">
<img alt="GitHub Release" src="https://img.shields.io/github/v/release/tinh-tinh/config">
<img alt="GitHub License" src="https://img.shields.io/github/license/tinh-tinh/config">
<a href="https://codecov.io/gh/tinh-tinh/config" > 
 <img src="https://codecov.io/gh/tinh-tinh/config/graph/badge.svg?token=VK57E807N2"/> 
 </a>
<a href="https://pkg.go.dev/github.com/tinh-tinh/config"><img src="https://pkg.go.dev/badge/github.com/tinh-tinh/config.svg" alt="Go Reference"></a>
</div>

<div align="center">
    <img src="https://avatars.githubusercontent.com/u/178628733?s=400&u=2a8230486a43595a03a6f9f204e54a0046ce0cc4&v=4" width="200" alt="Tinh Tinh Logo">
</div>

## Install

```bash
go get -u github.com/tinh-tinh/config/v2
```

## Overview

The Config module for the Tinh Tinh framework provides flexible configuration management for Go applications. It supports loading environment variables from `.env` files and struct-based configs from YAML, with first-class integration for dependency injection and modular apps.

## Features

- Load configuration from `.env`, `.yaml`, or `.yml` files
- Strongly-typed config structs with tags for mapping, default values, and validation
- Namespace-based config injection for multiple configs (e.g., database, cache)
- Conditional module registration based on environment or custom logic
- Supports default values and validation via struct tags
- Seamless integration with Tinh Tinh modules and controllers

## Quick Start

### 1. Basic Usage with `.env`

```go
import "github.com/tinh-tinh/config/v2"

type AppConfig struct {
    NodeEnv   string        `mapstructure:"NODE_ENV"`
    Port      int           `mapstructure:"PORT"`
    ExpiresIn time.Duration `mapstructure:"EXPIRES_IN"`
    Log       bool          `mapstructure:"LOG"`
    Secret    string        `mapstructure:"SECRET"`
}

cfg, err := config.NewEnv[AppConfig](".env")
if err != nil {
    panic(err)
}
fmt.Println(cfg.Port)
```

### 2. Using with YAML

```go
type YamlConfig struct {
    Host string `yaml:"host"`
    Port int    `yaml:"port"`
}

cfg, err := config.NewYaml[YamlConfig]("config.yaml")
if err != nil {
    panic(err)
}
fmt.Println(cfg.Host)
```

---

## Tinh Tinh Module Integration

### Register as a Global Config

```go
import "github.com/tinh-tinh/tinhtinh/v2/core"

appModule := core.NewModule(core.NewModuleOptions{
    Imports: []core.Modules{
        config.ForRoot[AppConfig](".env"),
    },
})
```

### Inject Config Anywhere

```go
cfg := config.Inject[AppConfig](module)
fmt.Println(cfg.Port)
```

## Contributing

We welcome contributions! Please feel free to submit a Pull Request.

## Support

If you encounter any issues or need help, you can:
- Open an issue in the GitHub repository
- Check our documentation
- Join our community discussions
