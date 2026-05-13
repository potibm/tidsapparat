# Tidsapparat

## Configuration

The application can be configured in multiple ways. This includes CLI parameters, environment variables and config files.

They are processed in the following order:

- CLI parameters
- Environment variables
- Config file
- Defaults.

### CLI Parameters

You can call the `tidsapparat serve` command with the following flags

```
-p, --port          int             # the port for the server
--otel-endpoint     string          # hostname and port for otel
--db-file           string          # the database file to use (without suffix)
--log-format        string          # log format, any of json, txt
--log-level         string          # log level, any of debug, info, warn, error
```

For details call `tidsapparat help serve`.

### Config

The application will use the config in `config/config.yaml` and `config/config.local.yaml` (which is ignored by git and useful for local overrides).

To generate a basic config you can use `tidsapparat config create`.

The yaml consists of those parts:

#### app

The app section defines core server settings, database paths, and basic network policies.

````
app:
    cors_allow_origins:
        - https://localhost:3300            
    db_filename: tidsapparat
    env: production 
    frontend_url: https://localhost:3300
    log_format: json
    log_level: info
    redis_url: "redis://pass@localhost:3305"
````

#### party config

Set up some information on your event.

`````
party:
    default_address: Prins Jørgens Gård 5, 1218 København K, Denmark
    end_date: "2026-05-17"
    start_date: "2026-05-15"
    timezone: Europe/Copenhagen
`````

#### event duration

Within the frontend you will see these quicklinks for the event duration:

`````
event_durations:
    - 0
    - 15
    - 30
    - 60
    - 90
    - 120
`````

#### exporter

You can set up multiple exporters that generate given filetypes to a destination.

`````
exporter:
    - name: ical_file_exporter              # any distinctive name
      type: ical                            # currently ical only
      destination: file                     # any of file or s3
      filename: events                      # the filename to write to (without suffix)
      options:                              # options for the destination
        dir: ./exports
      enabled: false                        # any of true or false
`````


#### s3 client

When you are using an exporter that writes to s3, you need to configure it accordingly.

`````
s3_client:
    access_key_id: accesskey
    secret_access_key: secretkey
    endpoint: http://localhost:9000
    region: us-east-1
    use_path_style: true
`````

#### format

Used to localize how dates and times are rendered in the generated feeds or frontend.

````
format:
    date:
        locale: da-DK           # the locale to use to display dates in the frontend
        options:                # javascript DateTimeFormat-Options, see below
            hour: 2-digit
            minute: 2-digit
            weekday: short
````

The date options are according to the documentation at https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Intl/DateTimeFormat/DateTimeFormat#options.

#### sentry

When using Sentry, provide the project specific DSN and further information here:

`````
sentry:
    dsn: ""
    environment: production
    replay_error_sample_rate: 0.1
    replay_session_sample_rate: 0.1
    trace_sample_rate: 0.1
    version: dev
`````

## Environment

You can overwrite all the config properties using environment variables. 

`app.env` translates to `APP_ENV`,  `app.redis_url` to `APP_REDIS_URL` and so on.

`APP_CORS_ALLOW_ORIGINS` can be provided as a comma separated list of values.

For more complex types (exporters e.g.) better use the config file mentioned above.