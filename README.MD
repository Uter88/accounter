### The accounter application is used to keep track of working hours.

#### Project structure:

- **cmd**:
    - **main.go**: backend app entry point;
    - **frontend**: **UI** app entry point;

- **config**: application configuration;

- **internal**: internal packages;
    - **app**: application engine;
    - **infrastucture**: infrastructure layer;
        - **adapter_sql**: repository adapters (**SQL**);
        - **auth**: authorization service (**JWT**);
        - **renderers**: output renderers (**HTML/PDF/CSV/DocX/XLSX/JSON**);
        - **comporators**: tools for determining the difference between objects;
        - **queues**: MQ brokers;
            - *kafka.go*: Kafka MQ;
        - **protocols**:
            - *websocket.go*: Websocket service (**gorila**);
        - **rest**: REST API HTTP server:
            - *server.go*: HTTP server wrapper (**Gin**);
            - **endpoints**: RestAPI endpoints;
                - **v1**: HTTP engine version 1;
        - **ui**: fronted UI:
            - **wasm-goapp**: WASM UI based on [**go-app**]( https://go-app.dev/);

    - **domain**: application domain layer:
        - **shared**: common use domain;
        - **task**: task domain;
        - **user**: user domain;
        - **event**: event domain;

- **pkg**: common usage packages;
    - **logger**: tool for logging;
    - **tools**: publicly available tools;

- **tools**: application maintenance tools;
    - **migrator**:
        - *migrator.go*: migrations for databases;

- **web**: statics web components; 
    - **fonts**: UI fonts;
    - **icons**: UI icons;
    - **js**: UI javascript libraries;
    - **styles**: CSS styles;
    - *app.wasm* (on building adapters/ui/wasm-goapp): UI builded **WASM** file;

- **templates**: static templates;
- **migrations**: database migration scripts;
- **scripts**: scripts to perform various build and migration analysis operations;
    - *migrate.sh*: launch migrator;
    - *start.sh*: application build script;
- **deployments**: deplotments tools (Dockerfile, Docker compose) 