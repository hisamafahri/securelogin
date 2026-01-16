# Rules

- project/module name is `github.com/hisamafahri/securelogin`

## Rules for Structuring Code

- This is an HTTP service. There are multiple layers to construct an endpoint:
    - `handlers` layer
       * handle HTTP requests and validating its payload/params
       * handle calling the appropriate usecases functions
    - `usecases` layer
       * handle the business logic of the application
       * handle calling the appropriate services functions
    - `services` layer
       * handle the interaction with external services (e.g. database, third-party APIs, etc)
       * each service function is atomic (do one thing and do it well)
    - `repositories` layer
       * handle the interaction with the database (e.g. queries, transactions, etc)
       * each repository function is atomic (do one thing and do it well)
       * can only be called by the services layer
    - `views` layer
       * handle the rendering of the HTML response
       * can only be called by the usecase layer
- each layer has own package (e.g. `handlers`, `usecases`, `services`, etc)

## Rules for Commands

- Run available commands in the root of the project (where `Makefiled` is located)

## Rules for Comments & Documentations

- don't add comments in code at all, I'll ask you to add comments if I want it
- for docs, always answer in markdown format
- don't use any emojis
- don't use any markdown styling (bold, italic, etc)
- keep it brief, short, but cover the details

