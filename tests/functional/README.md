# tests/functional

## Requirements

- `curl`

## Setup scripts

```console
# Give them execute permissions
root@debian:kube-nftlb/tests/functional# chmod +x *.sh
```

## Test resources

### Individually

```console
root@debian:kube-nftlb/tests/functional# ./test.sh 001_test_dir
```

### Multiple

```console
root@debian:kube-nftlb/tests/functional# ./test.sh 001_test_dir 002_test_dir 003_test_dir
```

### Everything

```console
root@debian:kube-nftlb/tests/functional# ./test.sh
```

## Clean resources

If anything goes wrong, you can run this script to clean every deployment or service pending.

```console
root@debian:kube-nftlb/tests/functional# ./clean.sh
```
