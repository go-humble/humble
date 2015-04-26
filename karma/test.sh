# This script compiles all the files in karma/go to karma/js
# and then runs karma with the mac configuration. You must be
# in the project root to run it.
echo "--> compiling karma tests to js..."
for gofile in $(find karma/go/*.go); do
	jsfile=${gofile//go/js}
	echo "    $gofile -> $jsfile"
	gopherjs build $gofile -o $jsfile | sed 's/^/    /'
done
echo "--> running karma tests"
# Change this line if you are on a platform other than mac
karma run karma/test-mac.conf.js | sed 's/^/    /'
echo "DONE"
