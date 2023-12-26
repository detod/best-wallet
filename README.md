# best-wallet

- Use docker-compose to bring the environment up.
- Run tests with `go test -v ./...`

## Gameplan

There are 2 main things that guide the decisions I've made (or will make) while
implementing the solution:

1. This is an MVP.
    - Clearly define the "must-have" and separate it from the "nice-to-have".
        - Example "must-have":
            - Retain as much data as possible for auditing (legal) and debugging.
            - No negative balance allowed.
        - Example "nice-to-have":
            - User friendly error messages.
            - Automated instead of manual auditing process.
    - Stick to as little backend infra as possible, ideally start with 1 service.
        - Increase availability with more than one instance. 
        - Recover from panics to prevent the process from crashing (every
        spawned goroutine needs it's own "recover").
        - Backup state.
    - Reduce the amount of boilerplate and leverage off-the-shelf stuff as much
    as possible (observability, CI/CD, see https://encore.dev/).
2. It needs to scale during high peak activity period (no specific number).
    - Keep processing in user request threads (goroutines) short.
        - Fail or succeed early, defer as much as possible to background
        processing.
        - Things like KYC, KYT and notifications, can all be done on a
        background thread.
        - Use state machines and account for intermediary states in the UX.
        - Transition states with scheduled jobs (in batches), or message queues
        (in real-time).
            - To avoid introducing message queues from the beginning, consider
            kickstarting a background process by spawning a goroutine from the
            request thread (don't wait for result, return to the user).
        - Have a way to check for errors, retry if transient, otherwise escalate
        to humans (the state machine progresses towards an end state).
        - Users are responsible for retries on errors in request threads, the
        backend is responsible for retries on errors in background processing.
    - Use connection pools when talking to other services (postgres, redis).
    - Use shorter timeouts for IO.
    - Cache data that rarely or never changes (redis, in-process).
        - e.g. HMAC keys, user profile...
        - Don't forget cache invalidation.
    - Some information can be calculated ahead of time instead of on request
    threads (e.g. account balance).
    - Don't work with huge datasets, use pagination/compaction.
    - Index db tables properly.
    - Monitor latency, resource utilization, errors and alert accordingly.
    - Do load testing.
    - Be prepared to scale the DB (first vertically, then horizontally e.g.
    read-replicas).
    - Keep the service stateless (besides some local temporary in-memory cache)
    and scale when necessary (vertically and horizontally).
    - Consider splitting read and write state and scale them accordingly.
    - Favor eventual consistency where possible, less coordination and contention.
    - Do maintenance off peak hours (auditing, heavy batch jobs).

### Implementation details

- CreateAccount endpoint stores the intent to open an account in a DB row and
kickstarts a background process that will handle the KYC process. It returns
immediately to the client without waiting for a result. The process will result
with either an open account or a failed KYC.
- A transaction can have three possible states "pending", "cleared", or "failed".
    - Pending transactions are going through (or will go through) KYT.
    - If KYT fails, the transaction fails.
    - Failed transaction cannot be retried, new transaction is needed.
    - If KYT suceeds, the transaction will be subject to clearing.
- All transactions are recoreded in a ledger with double book-keeping.
- A transaction can be of type "debit" (deducts from the account) or "credit"
(adds to the account).
- Existing transactions in the ledger cannot be modified.
- The account balance is equal to the sum of all the cleared "credit" TXs minus
all the cleared "debit" TXs on that account. This is the "cleared" balance.
- We can keep track of another balance called "available balance", which
represents the "cleared" balance minus all the pending "debit" TXs on the
account. 
- If we want to guarantee that the "cleared" balance will not go below zero,
we need to validate every new "debit" TX against the "available" balance before
reserving the funds. Checking the "available" balance and reserving the funds
needs to be atomic + serialized per account.

### Open questions

- Unique customer identifier that correlates accounts of the same customer?
- Should KYC be done only once for a customer, regardless of the number of
accounts they open?
- Allow new account creation requests while there is one already pending for the
same user?
- Should deposits (money coming inside the wallet) have a KYT?
- Is exchaning funds with external systems in the scope of this solution?

### Code structure

A description of a code structure I really like, especially for MVPs.

Background:

I view the backend as a big state machine. State is often kept in
DBs (like Postgres and Redis here). State transitions on the state machine are
initiated by actions such as API requests from client apps, background jobs,
webhooks, raised events etc. The transitions are carried out by http handlers,
event consumers, scheduled jobs etc. I call these the "entrypoints" or
"use-cases". They carry names like "create account", "read account", "deposit",
"withdraw", "transfer" etc.

The code structure looks like this:

/cmd -> binaries (servers/consumers, background jobs, CLIs)
/internal -> anything that you don't want exposed to the outside world
    /handler -> http handlers
    /middleware -> http middleware (a special form of handler)
    /consumer -> event consumers
    /domain -> anything domain specific, used by folders above
    /util -> anything non-domain specific, used by folders above

It's organized in layers:

Layer 0: cmd
Layer 1: handler, middleware, consumer, job
Layer 2: domain
Layer 3: util
- - - - - - - -
Layer 4: 3rd party code on github, gitlab...
Layer 5: go stdlib

Semantics:

- A layer can depend on layers below, but not on layers above. Example, code in
Layer 1 can depend on code in Layer 2 or 3, but not the other way around.
- Elements in the same layer should not depend on eachother e.g handlers <->
consumers.
- As you go down the layers, the code gets more generic and reusable. It makes
less assumptions, it's narrowly focused, and can satisfy different use-cases.
- As you go up the layers, the code gets more specific and less reusable. It
makes business level decisions that are only sensible within the context of a
specific application.
- Any layer can talk to a database or do IO.
- Keep things flat without a lot of nesting (sub-layers). If the project grows,
you can introduce sub-layers for namespacing.

Layers are important in Go since circular package dependencies are not allowed.
Note: A lot of things here can be subjective, I'm not religious about them :)
