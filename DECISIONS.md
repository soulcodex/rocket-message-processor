# Backend Engineer Challenge Decisions 🤔 ✍

## Architecture 🤠

* In terms of architecture, I've decided to create a simple service that's
  able to process the message received through the HTTP endpoint `/messages`
  like SQS does but without processing it in background and more oriented to
  a push model rather than a pull model like I normally do consuming messages
  from SQS.
* The service is designed following DDD (Domain Driven Design) principles, offering
  a clear separation of concerns, modular structure and a focus on business logic.

## Tooling 🔧

* The service is built using Go, a language known for its performance and simplicity.
* The project uses `just` for task automation, making it easier to manage development tasks.
* The service uses docker for containerization, allowing a seamless development experience.
* I've added `golangci-lint` for code quality checks, ensuring that the code adheres to best
  practices and standards accompanied by a git hook automatically installed on setup.

## Shortcuts ✍

* I've decided to omit the healthcheck endpoint as it is not a requirement for the challenge, but
  it can be easily added later even more if we're aiming for run our service in a container runtime such as k8s.
* I've omitted the usage of a database migrations tool like `sql-migrate` to keep the service simple and focused on the
  message processing logic, but in a real-world scenario, it would be essential to manage database schema changes.
* I've omitted the usage of a message broker like SQS or Kafka to keep the service simple and focused on the
  message processing logic, but in a real-world scenario, it would be essential to handle message delivery and
  processing such as retries, dead-letter queues, etc.
* I've omitted the usage of `air` to provide hot reloading during development, but it can be easily added later if
  needed.
