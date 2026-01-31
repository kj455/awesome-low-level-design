# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Overview

This is an educational repository for Low Level Design (LLD) / Object Oriented Design (OOD) interview preparation. It contains:
- 34+ LLD interview problems with solutions in 6 languages (Java, Python, C++, C#, Go, TypeScript)
- Design pattern implementations with explanations
- UML class diagrams for each problem

## CLI Commands

The repository includes an interactive CLI tool for practicing LLD problems in Go.

**Prerequisites:** `fzf` must be installed (`brew install fzf`)

```bash
# Start a problem - opens interactive selector, creates workspace
make lld start

# Reset a problem's main.go to template
make lld reset
```

The CLI creates workspaces in `workspace/go/<problem-id>/` with:
- `main.go` - starter template
- `problem.md` - problem requirements
- `main_test.go` - test file (user-created)

## Running Tests

```bash
# Run tests for a specific problem
cd workspace/go/<problem-id>
go test -v

# Run a single test
go test -v -run TestParkingLot
```

Tests use the `testify` library for assertions (`assert.NoError`, `assert.Equal`).

## Project Structure

```
problems/           # Problem specifications (markdown)
solutions/          # Reference implementations
  ├── java/src/     # Java solutions
  ├── python/       # Python solutions
  ├── cpp/          # C++ solutions
  ├── csharp/       # C# solutions
  ├── golang/       # Go solutions
  └── typescript/   # TypeScript solutions
design-patterns/    # Pattern implementations with README explanations
class-diagrams/     # UML diagrams (PNG)
workspace/go/       # Active Go workspace for practicing
scripts/lld/        # CLI tool source (Cobra-based)
```

## Code Conventions

**Go solutions:**
- Constructor functions: `NewParkingLot()`, `NewVehicle()`
- Thread safety: use `sync.Mutex` for concurrent access
- ID generation: use `atomic.AddInt64` for thread-safe counters
- Package per problem: `parkinglot`, `vendingmachine`

**All languages:**
- Each solution is self-contained and runnable independently
- Solutions use standard library only (no external frameworks)
- Problems map to corresponding UML diagrams in `class-diagrams/`

## Design Patterns Coverage

Creational: Singleton, Factory Method, Abstract Factory, Builder, Prototype
Structural: Adapter, Bridge, Composite, Decorator, Facade, Flyweight, Proxy
Behavioral: Iterator, Observer, Strategy, Command, State, Template Method, Visitor, Mediator, Memento, Chain of Responsibility
