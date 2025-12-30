### The accounter application is used to keep track of working hours.

#### Project structure:

- **cmd**:
    - **main.go**: backend app entry point;
    - **frontend**: **UI** app entry point;

- **config**:
    - *config.go*: application config entity;
    - *config.yaml*: **YAML** config file;

- **internal**: internal packages;
    - **app**: application engine;
    - **infrastucture**: application infrastructure layer;
        - **adapter_sql**: repository adapters (**SQL**);
        - **auth**: auth service (**JWT**);
        - **renderers**: output renderers (**HTML/PDF/CSV/DocX/XLSX/JSON**);
        - **protocols**:
            - *websocket.go*: Websocket service (**gorila**);
        - **rest_api**: RestAPI HTTP server:
            - *server.go*: HTTP server wrapper (**Gin**);
            - **endpoints**: RestAPI endpoints;
                - **v1**: HTTP engine version 1;
                    - *auth.go*: authorization;
                    - *engine.go*: V1 HTTP engine;
                    - *middlewares.go*: request middlewares;
                    - *task.go*: task endpoints;
                    - *user.go*: user endpoints;
        - **ui**: fronted UI:
            - **wasm-goapp**: WASM UI based on [**go-app**]( https://go-app.dev/);
    - **domain**: application domain layer:
        - **task**:
            - *params.go*: task query params;
            - *ports.go*: task ports;
            - *service.go*: task use cases;
            - *task.go*: task entity;
        - **user**:
            - *ports.go*: user ports;
            - *service.go*: user use cases;
            - *user.go*: user entity;

- **pkg**: common usage packages
    - **logger**: tool for logging;
    - **tools**:
        - *containers.go*: data containers;
        - *dateutils.go*: utils for working with date/time;
        - *documents.go*:  utils for working with documents;
        - *formatters.go*: utils for formatting values;
        - *helpers.go*: utils for working with types;
        - *requests.go*: utils for HTTP requests;
        - *validators.go*: input data validators;

- **tools**: application tools;
    - **migrator**:
        - *migrator.go*: migrations for databases;

- **web**:
    - **fonts**: UI fonts;
    - **icons**: UI icons;
    - **js**: UI javascript libraries;
    - **styles**: CSS styles;
    - *app.wasm* (on building adapters/ui/wasm-goapp): UI builded **WASM** file;

- **templates**: static templates;

- **scripts**:
    - *migrate.sh*: launch migrator;
    - *start.sh*: application build script; 