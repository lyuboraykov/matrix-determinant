import java.util.function.IntConsumer;

public class Argument {
    public String lineOption;
    public String description;
    public IntConsumer trigger;
    public Argument(String lineOption, String description, IntConsumer trigger) {
        lineOption = lineOption;
        description = description;
        trigger = trigger;
    }
}
