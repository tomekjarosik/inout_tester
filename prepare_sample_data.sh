
DATA_DIR="problems/multiply_by_2"
mkdir -p "${DATA_DIR}"

echo "Preparing test 1"
echo "1" > ${DATA_DIR}/t1.in
echo "2" > ${DATA_DIR}/t1.out

echo "Preparing test 2"
echo "123" > ${DATA_DIR}/t2.in
echo "246" > ${DATA_DIR}/t2.out

echo "Preparing test 3"
echo "-22" > ${DATA_DIR}/t3.in
echo "-44" > ${DATA_DIR}/t3.out