/*
Shell sort is a way to deal with worst case scenarios of the insertion sort

Let's assume a sorting problem below:
10, 9, 8, 7, 6, 5, 4, 3, 2, 1

- We pick a number that "skip": 3

We would then split the above problem to the following smaller arrays that s to be sorted:
(10, 7, 4, 1), (9, 6, 3), (8, 5, 2)

After doing a insertion sort on the above... we would arrange it as the following:
(1, 4, 7, 10), (3, 6, 9), (2, 5, 8)

Which then we would piece back tgt:
1, 3, 2, 4, 6, 5, 7, 9, 8, 10

Doing a another round of "skip": 2
(1, 2, 6, 7, 8), (3, 4, 5, 9, 10)

Which would piece to:
1, 3, 2, 4, 6, 5, 7, 9, 8, 10

After this would be the normal insertion sort across the whole array
1, 2, 3, 4, 5, 6, 7, 8, 9, 10


*/