import java.util.Scanner;

class Solution {
    public static void main(String[] args) {
        Scanner scanner = new Scanner(System.in);
        String line = scanner.nextLine();  // reads the single input line from the console
        String[] strings = line.split(" ");  // splits the string wherever a space character is encountered, returns the result as a String[]
        int first = Integer.parseInt(strings[0]);
        int second = Integer.parseInt(strings[1]);
        System.out.println("First number = " + first + ", second number = " + second + ".");
    }
}