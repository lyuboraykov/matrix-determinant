import java.util.*;
import java.util.concurrent.ThreadLocalRandom;
import java.util.stream.*;


public class Determinant {

    private static int matrixSize = 10;
    private static int threadCount = 8;
    private static Boolean quiet = false;

    public static void main(String[] args) {
        int UPPER_BOUND = 10;

        parseArguments(args);

        int[][] matrix = new int[matrixSize][matrixSize];
        for (int i = 0; i < matrixSize; i++) {
            for (int j = 0; j < matrixSize; j++) {
                matrix[i][j] = ThreadLocalRandom.current().nextInt(UPPER_BOUND);
            }
        }

        int[][] permutations = permute(matrixSize);

        Date start = new Date();
        Runnable[] tasks = new Runnable[threadCount];
        Thread[] threads = new Thread[threadCount];
        int[] partialResults = new int[threadCount];
        int permutationsSliceSize = permutations.length / threadCount;
        for (int i = 0; i < threadCount; i++) {
            final int resultIdx = i;
            int startIndex = resultIdx * permutationsSliceSize;
            int endIndex = (resultIdx + 1) * permutationsSliceSize;
            int[][] permutationsPiece = Arrays.copyOfRange(permutations, startIndex, endIndex);
            tasks[i] = () ->  partialResults[resultIdx] = calculateDeterminant(matrix, permutationsPiece);
            threads[i] = new Thread(tasks[i]);
            threads[i].start();
        }
        for (Thread thread: threads) {
            try {
                thread.join();
            } catch (InterruptedException e) {
                e.printStackTrace();
            }
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

    private static int determinantSign(int[] permutation) {
        int inversions = 0;
        for (int i = 0; i < permutation.length; i++) {
            for (int j = i; j < permutation.length; j++) {
                if (permutation[j] < permutation[i]) {
                    inversions++;
                }
            }
        }
        return (int)Math.pow(-1, inversions % 2);
    }

    private static int calculateDeterminant(int[][] matrix, int[][] permutations) {
        int determinant = 0;
        for (int[] permutation: permutations) {
            int multiplier = determinantSign(permutation);
            int multiple = 1;
            for (int i = 0; i < permutation.length; i++) {
                multiple *= matrix[i][permutation[i]];
            }
            determinant += multiplier * multiple;
        }
        return determinant;
    }
}


