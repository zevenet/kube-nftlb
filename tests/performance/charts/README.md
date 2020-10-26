# tests/performance/charts

Directory where `gnuplot` charts are stored after running every `chart.gp` script.

## Y axes

### time

How much time (in milliseconds, _ms_) does it take to create or delete a Service for N endpoints?

- Represented as a [boxplot](https://en.wikipedia.org/wiki/Box_plot).

### rules

Hoy many rules are there after creating or deleting a Service for N endpoints?

- Represented as a [bar chart](https://en.wikipedia.org/wiki/Bar_chart).

## X axes

### by-count-type

Named internally as `create-service` or `delete-service`, they represent the action they do (creating a Service or deleting a Service).

### by-endpoints-number

Number of replicas defined in a deployment.

## `chart.gp` files

They are `gnuplot` scripts meant to generate a chart related to where they are located.

## Expected output

After running `generate_charts.sh`, this directory should look like this:

```
charts
├── README.md
├── rules
│   ├── by-count-type
│   │   ├── chart.gp
│   │   ├── replicas-test-010.png
│   │   ├── replicas-test-050.png
│   │   └── replicas-test-100.png
│   └── by-endpoints-number
│       ├── chart.gp
│       ├── create-service.png
│       └── delete-service.png
└── time
    ├── by-count-type
    │   ├── chart.gp
    │   ├── replicas-test-010.png
    │   ├── replicas-test-050.png
    │   └── replicas-test-100.png
    └── by-endpoints-number
        ├── chart.gp
        ├── create-service.png
        └── delete-service.png

6 directories, 15 files
```
