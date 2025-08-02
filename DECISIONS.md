# Backend Engineer Challenge Decisions 🤔 ✍

## Architecture & Design 🤠

* In terms of architecture, I've decided to create a simple service that's
  able to process the message received through the HTTP endpoint `/messages`
  like SQS does but without processing it in background and more oriented to
  a push model rather than a pull model like I normally do consuming messages
  from SQS.
* Observability is a key aspect but is out of scope for this challenge, but as you can see there are some components
  of mine that are already taking into account things like semantic conventions for logging, error handling and
  metrics, but I haven't implemented any observability solution in this service to keep it simple.
* The service is designed following DDD (Domain Driven Design) principles, offering
  a clear separation of concerns, modular structure and a focus on business logic.
* I usually split the repositories in writer and reader interfaces but in this case I've decided to keep it simple
  and have a single interface for the rocket repository, but it can be easily split later if needed.
* In terms of persistence, I've chosen to use PostgresSQL as the database at the beginning, which is a robust and
  reliable choice for handling structured data but to make it simple and focused on the message processing logic I've
  decided to use instead an in memory implementation of the rocket repository to keep my implementation simple, but in a
  real-world scenario, it would be essential to use a persistent storage solution like PostgresSQL or another database
  for obvious reasons.
* Due to the absence of an identifier in the message received, I've decided to use the elements in the message metadata
  to create a unique identifier and implement a deduplication mechanism in memory to avoid processing the same message
  multiple times, which is a common requirement in message processing systems.
* The deduplication mechanism is implemented in memory using a map as structure to store the processed messages,
  which is a simple and clear way but in a realistic scenario, where you might have a large amount of messages and
  service run in various instances (replicas) it would be essential to think about an state of deduplication
  distributed, otherwise you wont get the expected results.
* I do believe that the messages structure created on my own to approach the problem can be absolutely better in terms
  of its basis. To foster a common base structure for all the messages I've to be more consistent, offer a better DX and
  reduce the boilerplate code that's going to make the code even more open to changes and evolution in the future.
* I've decided to add an small DI (Dependency Injection) container to the service to manage dependencies and
  facilitate testing, which is a common practice in Go applications to improve modularity and testability, but it can
  be easily improved using a more sophisticated DI framework if needed in the future such as `wire` or similar.
* The receiver HTTP handler is simple but isn't that clean and could be improved a lot in terms of
  error handling, validation, components used on it and response formatting, but it can be easily improved later if
  needed.

## Tooling 🔧

* The service is built using Go, a language known for its performance and simplicity.
* The project uses `just` for task automation, making it easier to manage development tasks.
* The service uses docker for containerization, allowing a seamless development experience.
* I've added `golangci-lint` for code quality checks, ensuring that the code adheres to best
  practices and standards accompanied by a git hook automatically installed on setup.
* I've decided to user `moq` for mocking dependencies in tests, which simplifies the process of creating mock
  implementations for interfaces.

## Shortcuts ✍

* I've decided to omit the healthcheck endpoint as it is not a requirement for the challenge, but
  it can be easily added later even more if we're aiming for run our service in a container runtime such as k8s.
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
* I do believe that events serialization could be better handled. I'd rather in fact create an standard message
  structure for all the messaging to avoid repeat code, foster consistency across messaging and make event creation
  easier to maintain and evolve in the future, but I didn't it in this case to keep the service simple and focused on
  the message processing logic, but in a real-world scenario, it would be essential to handle message serialization
  and deserialization properly.
* I've taken an small shortcut in terms of HTTP error handling, I've decided to use a simple error handling strategy
  that returns a generic error messages instead of a more detailed one, but it can be easily improved later if needed.
