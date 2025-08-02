<div align="center" style="text-align: center; padding: 20px">
    <h2>🚀 Rockets Message Processor 🚀</h2>
</div>

![Go](https://img.shields.io/badge/Go-1.24.1-blue.svg?style=for-the-badge)

## ⚙️ Requirements

- [Go](https://golang.org/doc/install) `1.24 or earlier`
- [Docker](https://docs.docker.com/get-docker/) or [Orbstack](https://orbstack.dev/download)
- [Just](https://github.com/casey/just#installation)
- [golang-ci lint](https://github.com/golangci/golangci-lint)
- [Moq](https://github.com/matryer/moq)

## 🤔 Challenge decisions

In order to keep this README clean and focused on the project setup, I've documented my decisions and design choices in
a separate file called [DECISIONS.md](./DECISIONS.md)

## 📥 Installation

Clone the repository using Git:

```bash
git clone https://github.com/soulcodex/rockets-message-processor.git
```

### ⚙️ Install `just`

The repository uses [`just`](https://github.com/casey/just) for task automation. Install it with:

```bash
curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh | sudo bash -s -- --to /usr/local/bin
```

### 📦 Install Tools & Dependencies

To install all project external dependencies and tools, run:

```bash
just install
```

## ▶️ Running the rockets message processor

To start one the rockets message processor service in local, execute the following command:

```bash
just run
```

## 🐳 Running using a `docker-compose` stack

To start the rockets message processor **Docker Compose** stack, run:

```bash
just up
```

To shut down the stack, use:

```bash
just down
```

> This executes the `docker-compose.yml` file located in `deployments/docker-compose/`.

## 📁 Project Structure

For a deeper dive into project structure best practices, check out the
**[Go Standard Layout](https://github.com/golang-standards/project-layout)** repository.
