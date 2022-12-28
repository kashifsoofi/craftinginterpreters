using cslox.Parser;
using cslox.Scanning;

class Lox
{
    static bool hadError = false;

    public static int Main(string[] args)
    {
        Expression expression = new Binary(
            new Unary(
                new Token(TokenType.MINUS, "-", null, 1),
                new Literal("123")),
            new Token(TokenType.STAR, "*", null, 1),
            new Grouping(
                new Literal(45.67)));

        Console.WriteLine(new AstPrinter().Print(expression));

        if (args.Length > 1)
        {
            Console.WriteLine("Usage: cslox [script]");
            return 64;
        }
        else if (args.Length == 1)
        {
            RunFile(args[0]);
        }
        else
        {
            RunPrompt();
        }

        if (hadError)
        {
            return 65;
        }
        return 0;
    }

    static void RunFile(string path)
    {
        string source = File.ReadAllText(path);
        Run(source);
    }

    static void RunPrompt()
    {
        while (true)
        {
            Console.Write("> ");
            var line = Console.ReadLine();
            if (line == null)
            {
                break;
            }
            Run(line);
        }
    }

    static void Run(string source)
    {
        var scanner = new Scanner(source);
        var tokens = scanner.ScanTokens();

        foreach (var token in tokens)
        {
            Console.WriteLine(token);
        }
    }

    public static void Error(int line, string message)
    {
        Report(line, "", message);
    }

    static void Report(int line, string where, string message)
    {
        Console.WriteLine($"[line {line}] Error{where}: {message}");
        hadError = true;
    }
}
