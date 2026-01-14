### The accounter application is used to keep track of working hours.

#### Project structure:

- **cmd**:
    - **backend**: backend entry point;
    - **frontend**: **UI** app entry point;
    - **migrator**: databases migrations;

- **config**: application configuration;

- **internal**: internal packages;
    - **app**: application engine;
    - **infrastructure**: infrastructure layer;
        - **adapters**: repository adapters;
            - **postgres**: **PostgreSQL** repositories;
        - **auth**: authorization service (**JWT**);
        - **renderers**: output renderers (**HTML/PDF/CSV/DocX/XLSX/JSON**);
        - **comparators**: tools for determining the difference between objects;
        - **messaging**: MQ brokers;
            - **kafka**: Kafka MQ;
        - **transport**: transport layers;
            - **websocket**: websocket service (**gorilla**)
            - **rest**: REST API HTTP server:
                - *server.go*: HTTP server wrapper (**Gin**);
                - **endpoints**: RestAPI endpoints;
                    - **v1**: HTTP engine version 1;
        - **ui**: fronted UI:
            - **wasm-goapp**: WASM UI based on [**go-app**]( https://go-app.dev/);

    - **domain**: application domain layer:
        - **common**: common use domain;
        - **task**: task domain;
        - **user**: user domain;
        - **event**: event domain;

- **pkg**: common usage packages:
    - **logger**: tool for logging;
    - **utils**: publicly available utilities;

- **web**: static web components:
    - **fonts**: UI fonts;
    - **icons**: UI icons;
    - **js**: UI javascript libraries;
    - **styles**: CSS styles;
    - *app.wasm* (on building adapters/ui/wasm-goapp): UI builded **WASM** file;

- **templates**: static templates;
- **migrations**: database migration scripts;
- **scripts**: scripts to perform various build and migration analysis operations:
    - *migrate.sh*: launch migrator;
    - *start.sh*: application build script;
- **deployments**: deployments tools (Dockerfile, Docker compose);