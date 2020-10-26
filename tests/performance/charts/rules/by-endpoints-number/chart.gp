#!/usr/bin/env gnuplot

# Args
TITLE=ARG1
OUTPUT_FILE=ARG2
DATA_PATH_A=ARG3
DATA_PATH_B=ARG4
KUBE_NAME_A=ARG5
KUBE_NAME_B=ARG6

# Output image
set terminal png size 1500,1000
set output OUTPUT_FILE
set title TITLE
set title font ",24"
set key outside

# Chart style
set style line 1 lc rgb "green"
set style line 2 lc rgb "light-blue"
set style fill solid border rgb "black"
set boxwidth 0.4
set border 3

# x-axis
set xtics nomirror
set xtics ( "10 endpoints" 0.2, "50 endpoints" 1.2, "100 endpoints" 2.2)

# y-axis
set ytics 0, 10
set yrange [0:]
set ytics nomirror
set ylabel "Number of rules"

# Draw plot and save it
plot \
    DATA_PATH_A using ($0):1 \
                title KUBE_NAME_A \
                with boxes \
                ls 1, \
    DATA_PATH_B using ($0+0.4):1 \
                title KUBE_NAME_B \
                with boxes \
                ls 2
