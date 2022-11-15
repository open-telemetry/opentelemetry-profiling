# Open Telemetry profiling benchmarks

This repository consists of :
* data sets (`/profiles`)
* tool to convert (`/cmd/convert`)
* implementations of Open Telemetry profiling format converter (`/implementations/*`)
* tool to generate benchmarking reports (`/cmd/report`)

![diagram describing relationships between main components](./diagram.drawio.svg)

### Usage

#### To convert from source formats to intermediary format

```bash
# this will take files from /profiles/src and put converted files to /profiles/intermediary
make convert
```

#### To run the benchmarks and generate reports

```bash
make benchmark
```

### How to contribute

You can contribute by creating issues or pull requests. Best ways to contribute are:

* You can add profiles to the data set. In order to do that first add your profiles to `/profiles/src` directory. You can then run `make convert` to perform anonymization of symbol names and conversion into intermediary format.

* You can modify the encoder implementations (`/implementations`) and make improvements.

* You can make any other improvements to the repository.


### TODOs

This repository is work-in-progress. Most functionality is not yet implemented.

TODO:
* [ ] populate source data set (`/profiles/src`)
* [ ] implement converters / anonymizers for conversions from source formats to the intermediate format
* [ ] implement a tool for generating benchmark reports
* [ ] create a github actions workflow to run benchmarks on every PR and print results comparing the change to the main branch
