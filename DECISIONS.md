# Backend Engineer Challenge Decisions 🤔 ✍

## Architecture 🤠

* In terms of architecture, I've decided to create a simple service that's
  able to process the message received through the HTTP endpoint `/messages`
  like SQS does but without processing it in background and more oriented to
  a push model rather than a pull model like I normally do consuming messages
  from SQS.
* The service is designed following DDD (Domain Driven Design) principles, offering
  a clear separation of concerns, modular structure and a focus on business logic.
* I usually split the repositories in writer and reader interfaces but in this case I've decided to keep it simple
  and have a single interface for the rocket repository, but it can be easily split later if needed.

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
* I've used a set of tools already written by me for my own usage in different side projects, such as
  Zerolog for logging, a retry mechanism without backoff, a distributed mutex, message bus, utils and semantic error
  handling
  if there's any doubt about these components, I'll be happy to explain them in detail.
* I've decided to not validate the rocket type at domain level due to not knowing the full list of rocket types
  that will be used now or added in the future, but it can be easily added later if the types are finite or static
  defined.
