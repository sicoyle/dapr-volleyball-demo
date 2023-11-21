namespace RoundWorkflow.Models
{
    // Processing
    public record RoundInput(string[] Teams, int Round);

    public record Pair(string TeamA, string TeamB);

    public record GameResult(Dictionary<string, int> TeamA, Dictionary<string, int> TeamB);

    // Output
    public record RoundResult(bool Succeeded, List<string>? WinningTeams);

}