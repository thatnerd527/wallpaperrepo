public record MouseEvent
{
    public double ClientX { get; init; } = 0;
    public double ClientY { get; init; } = 0;

    public double ScreenX { get; init; }= 0;
    public double ScreenY { get; init; }= 0;
    public double PageX { get; init; }= 0;
    public double PageY { get; init; }= 0;
    public double OffsetX { get; init; }= 0;
    public double OffsetY { get; init; }= 0;
    public double MovementX { get; init; }= 0;
    public double MovementY { get; init; }= 0;

    public MouseEvent(double clientX, double clientY, double screenX, double screenY, double pageX, double pageY, double offsetX, double offsetY, double movementX, double movementY)
    {
        this.ClientX = clientX;
        this.ClientY = clientY;
        this.ScreenX = screenX;
        this.ScreenY = screenY;
        this.PageX = pageX;
        this.PageY = pageY;
        this.OffsetX = offsetX;
        this.OffsetY = offsetY;
        this.MovementX = movementX;
        this.MovementY = movementY;

    }
}

public class MouseHandler
{
    private Action<MouseEvent> action;

    private Action destroyAction;

    public MouseHandler(Action<MouseEvent> action, Action destroyAction)
    {
        this.action = action;
        this.destroyAction = destroyAction;
    }

    public void Invoke(MouseEvent e)
    {
        action(e);
    }

    public void Destroy()
    {
        destroyAction();
    }
}

public class MouseUpHandler
{
    private Action action;

    private Action destroyAction;

    public MouseUpHandler(Action action, Action destroyAction)
    {
        this.action = action;
        this.destroyAction = destroyAction;
    }

    public void Invoke()
    {
        action();
    }

    public void Destroy()
    {
        destroyAction();
    }
}
