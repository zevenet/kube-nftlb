#!/bin/bash
# to execute use bash script

# It reads the name of each directory, applies the creation / deletion configuration of pods and makes the curl requests to generate logs
# The tests are as generic as possible to parameterize the entire process. Therefore, there are some tests that need more in-depth tests. 
# in that case each directory can contain three different script types: 

# before.sh = runs before main script
# after.sh = run after the main script
# instead.sh = Some tests need specific tests, this script replaces the main script and does its own tests.

after="after.sh"
before="before.sh"
instead="instead.sh"

for directory in *; do
    if [ -d "$directory" ]; then
	sleep 10
        echo "launch test $directory"
        if [[ -n $(find $directory/ -name $before) ]] 
            then
                cd $directory/
                echo "launch script $before"
                cd ..
        elif [[ -n $(find $directory/ -name $instead) ]] 
            then
                echo "|==> launch only script $instead"
                cd $directory/
                bash $instead
                cd ..
        else
            kubectl apply -f $directory/ &>/dev/null
            # Sleeps are defined to leave a space for pods to rise and not cause problems when obtaining logs for non-existent pods.
            # Each creation or deletion process has an action time
            sleep 10
            curl --silent -H "Key: 12345" http://localhost:5555/farms/ > $directory/configCreation.nft
            kubectl delete -f $directory/ &>/dev/null
            sleep 10
            curl --silent -H "Key: 12345" http://localhost:5555/farms/ > $directory/configDelete.nft
            # We apply a format to eliminate all those parameters that may vary from one test to another. That is, all the parameters that kubernetes generates randomly.
            # This then allows us to better appreciate the difference between a correct and an incorrect test once uploaded to the repository.
            sed -i 's/\("virtual-addr": \)\(.*\)/\1"IP"/' $directory/configCreation.nft
            sed -i 's/\("ip-addr": \)\(.*\)/\1"IP"/' $directory/configCreation.nft
            sed -i -r 'N;s/.*(\n(\s*)"ip-addr":)/\2"name": "NAME",\1/;P;D' $directory/configCreation.nft
        fi
    fi
done
