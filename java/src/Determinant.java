import java.util.*;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.ThreadLocalRandom;
import java.util.stream.*;
import java.util.concurrent.ForkJoinPool;
import java.util.stream.Collectors;


public class Determinant {

    private static int matrixSize = 10;
    private static int threadCount = 2;
    private static Boolean quiet = false;

    public static void main(String[] args) {
        Date start = new Date();
        int UPPER_BOUND = 10;

        parseArguments(args);

        int[][] matrix = new int[matrixSize][matrixSize];
        for (int i = 0; i < matrixSize; i++) {
            for (int j = 0; j < matrixSize; j++) {
                matrix[i][j] = ThreadLocalRandom.current().nextInt(UPPER_BOUND);
            }
        }


        ForkJoinPool workerPool = new ForkJoinPool(threadCount);

        try {
            List<Integer> nums =  workerPool.submit(() -> IntStream.range(0, matrixSize)
                                                                   .boxed()
                                                                   .parallel()
                                                                   .collect(Collectors.toList())).get();
            List<Stream<Integer>> streams = workerPool.submit(() -> Permutations.of(nums)
                                                                                .parallel()
                                                                                .collect(Collectors.toList())).get();

            List<List<Integer>> permutations;

            permutations = workerPool.submit(() -> streams.stream().parallel()
                                                                   .map((Stream<Integer> stream) -> stream.parallel()
                                                                                                          .collect(Collectors.toList()))
                                                                   .collect(Collectors.toList())).get();

            Runnable[] tasks = new Runnable[threadCount];
            int[] partialResults = new int[threadCount];
            int permutationsSliceSize = permutations.size() / threadCount;
            for (int i = 0; i < threadCount; i++) {
                final int resultIdx = i;
                int startIndex = resultIdx * permutationsSliceSize;
                int endIndex = (resultIdx + 1) * permutationsSliceSize;
                List<List<Integer>> permutationsPiece = permutations.subList(startIndex, endIndex);
                tasks[i] = () ->  partialResults[resultIdx] = calculateDeterminant(matrix, permutationsPiece);
                new Thread(tasks[i]).start();
            }
            int determinant = IntStream.of(partialResults).sum();
            for (int[] row: matrix) {
                for(int el: row) {
                    System.out.print(el + " ");
                }
                System.out.println();
            }
            Date end = new Date();
            System.out.println((end.getTime() - start.getTime()) / 1000.0);
            System.out.println(determinant);
        } catch (InterruptedException | ExecutionException exc) {
            System.out.println(exc);
        }
    }

    private static final Map<String, Argument> ARGUMENTS;
    static {
        String QUIET = "-q";
        String THREAD_COUNT = "-t";
        ARGUMENTS = new HashMap<>();
        ARGUMENTS.put(QUIET, new Argument(QUIET,
                                          "Disable verbose printing",
                                          (int noop) -> quiet = true));
        ARGUMENTS.put(THREAD_COUNT, new Argument(THREAD_COUNT,
                                                 "-t <int> set the number of threads to run",
                                                 (int count) -> threadCount = count));
    }


    private static void parseArguments(String[] args) {
        // TODO: Implement
        return;
    }

    private static int[][] permute(int n) {
        int[][] permutations = new int[factorial(n)][n];
        int[] numbers = IntStream.range(0, n).toArray();
        doPermute(permutations, numbers, 0);
        return permutations;
    }

    private static int rowNumber = 0;

    private static void doPermute(int[][] permutations, int[] numbers, int col) {
        int n = numbers.length;
        System.out.println(n);
        if (col == n) {
            permutations[rowNumber] = Arrays.copyOf(numbers, numbers.length);
            rowNumber++;
        }

        for (int i = col; i < n; i++) {
            swap(numbers, col, i);
            doPermute(permutations, numbers, col + 1);
            swap(numbers, col, i);
        }
    }

    private static void swap(int[] arr, int i, int j) {
        int temp = arr[i];
        arr[i] = arr[j];
        arr[j] = temp;
    }

    private static int factorial(int n) {
        int fact = 1;
        for (int i = n; i > 0; i--) {
            fact *= i;
        }
        return fact;
    }

    private static int determinantSign(List<Integer> permutation) {
        int inversions = 0;
        for (int i = 0; i < permutation.size(); i++) {
            for (int j = i; j < permutation.size(); j++) {
                if (permutation.get(j) < permutation.get(i)) {
                    inversions++;
                }
            }
        }
        return (int)Math.pow(-1, inversions % 2);
    }

    private static int calculateDeterminant(int[][] matrix, List<List<Integer>> permutations) {
        int determinant = 0;
        for (List<Integer> permutation: permutations) {
            int multiplier = determinantSign(permutation);
            int multiple = 1;
            for (int i = 0; i < permutation.size(); i++) {
                multiple *= matrix[i][permutation.get(i)];
            }
            determinant += multiplier * multiple;
        }
        return determinant;
    }
}


