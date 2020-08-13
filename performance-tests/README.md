# performance-tests

## Test

```console
# Save your results by piping the output to a file
root@debian:kubernetes-rules-test# ./test.sh > results/your-test.txt
```

## Parse results

```console
root@debian:kubernetes-rules-test# cat results/your-test.txt | sed -e '/^Starting/d' -e '/^Deleting/d' -e '/^serviceaccount\//d' -e '/^clusterrolebinding[.]/d' -e '/^daemonset[.]/d' -e '/^deployment[.]/d' -e '/^pod\//d' -e '/^service/d' -e '/^$/d'
```

## Clean

If anything goes wrong, you can run this script to clean every deployment,service or daemonset pending.

```console
root@debian:kubernetes-rules-test# ./clean.sh
```

