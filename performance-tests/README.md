# performance-tests

## Requirements

- **gnuplot** (to make benchmark charts)

```console
root@debian:~# apt-get install gnuplot
```

## Setup scripts

```console
# Copy and rename expected-rule-count.sh
root@debian:kube-nftlb/performance-tests# cp expected-rule-count.sh.example expected-rule-count.sh

# Give them execute permissions
root@debian:kube-nftlb/performance-tests# chmod +x *.sh
```

## Get expected rule count

⚠ Edit `expected-rule-count.sh` before continuing. This must be done for every file in `kubes/`, because counting rules can differ between each kube. Read `test.sh` to know already defined functions for this.

```console
# Pass your kube-test file as the first parameter (this is an example)
root@debian:kube-nftlb/performance-tests# ./expected-rule-count.sh ./kubes/kube-test.yaml
```

```console
# To specify a single (or several) deployment(s), pass them after specifying the first parameter (this is an example)
root@debian:kube-nftlb/performance-tests# ./expected-rule-count.sh ./kubes/kube-test.yaml ./testdata/deployments/resource-test.yaml ./testdata/deployments/resource-test-2.yaml
```

## Test resources

### Only once

```console
# Save your results by piping the output to a file
root@debian:kube-nftlb/performance-tests# ./test.sh > results/your-test.txt
```

### Several times (recommended)

```console
# Set how many times the tests will be run by changing REPEATS
root@debian:kube-nftlb/performance-tests# REPEATS=5 ; for i in $(seq 1 $REPEATS) ; do ./test.sh > "results/your-test-$i.txt" ; done
```

## Filter results

### Individually

```console
# Beware, this will edit your file! If you only want to see the filtered output, remove the "-i" flag
root@debian:kube-nftlb/performance-tests# sed -i -f filters/result.sed results/your-result.txt
```

### Everything (recommended)

```console
# Save your filtered results by piping the output to a file
root@debian:kube-nftlb/performance-tests# sed -f filters/result.sed results/* > filtered-results.txt
```

## Make benchmark charts

⚠ `filtered-results.txt` is required.

```console
# A filepath with all your filtered results must be specified as first parameter
root@debian:kube-nftlb/performance-tests# ./generate_pool_data.sh filtered-results.txt

# Two daemonset *names* must be specified as parameters
root@debian:kube-nftlb/performance-tests# ./generate_charts.sh kube-test-1 kube-test-2
```

## Clean resources

If anything goes wrong, you can run this script to clean every deployment, service or daemonset pending.

```console
root@debian:kube-nftlb/performance-tests# ./clean.sh
```
