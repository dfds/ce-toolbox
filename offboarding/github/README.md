# GitHub offboarding tool

A tool written in Rust that checks all repositories within DFDS for deploy keys and collaborators.

## How?

Either

**A**: Grab a pre-compiled binary release at https://github.com/dfds/ce-toolbox/actions?query=workflow%3A%22github+offboarding+tool+-+Binary+builds%22++

or

**B**: Build it yourself

### Building the tool

If you don't already have Rust installed, you'll need that. See https://rustup.rs/ for instructions for your OS.

With Rust installed, run `cargo run` within this directory to build and run the tool.

### Usage

The tool will be expecting a GITHUB_TOKEN environment variable, e.g.

macOS/Unix:
```
GITHUB_TOKEN=1234 cargo run
```


Windows (Powershell):
```
$env:GITHUB_TOKEN=1234
cargo run
```
