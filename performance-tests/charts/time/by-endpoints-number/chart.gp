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
set title font "Helvetica,24"
set key outside

# Chart style
set style line 1 lc rgb "green"
set style line 2 lc rgb "light-blue"
set style fill solid border rgb "black"
set style boxplot outliers pointtype 7
set style data boxplot
set boxwidth 0.3
set pointsize 0.5
set border 3

# x-axis
set xtics
set xtics border in scale 0,0 nomirror norotate
set xtics norangelimit
set xtics ( "10 endpoints" 1, "50 endpoints" 2, "100 endpoints" 3)

# y-axis
set yrange [0:]
set ytics border in scale 1,0.5 nomirror norotate
set ylabel "Time (in milliseconds)"

# Draw plot and save it
plot \
    for [i=1:3] DATA_PATH_A \
                using (i-0.15):i \
                title (i == 1 ? KUBE_NAME_A : "") \
                ls 1, \
    for [i=1:3] DATA_PATH_B \
                using (i+0.15):i \
                title (i == 1 ? KUBE_NAME_B : "") \
                ls 2
