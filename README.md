# Panawa batch
Backend jobs of Panawa service

## Build & Run
```sh
# build executable
$ make

# run locally
$ ./batch {command} --config=sample-config.yaml
```

## Commands
### `backfill`
- Fetches and saves price data of date range
- Set `backfill.startdate` and `backfill.enddate` on cli flag or config file

### `fetch`
- Fetches and saves price data of current date
