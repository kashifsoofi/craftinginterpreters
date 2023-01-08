using cslox.Interpreter;
using cslox.Parser;
using cslox.Scanning;

class Lox
{
    static Interpreter interpreter = new Interpreter();
    static bool hadError = false;
    static bool hadRuntimeError = false;

    public static int Main(string[] args)
    {
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
        if (hadRuntimeError)
        {
            return 70;
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
        // Reset errors
        hadError = false;
        hadRuntimeError = false;

        var scanner = new Scanner(source);
        var tokens = scanner.ScanTokens();
        var parser = new Parser(tokens);
        var statements = parser.Parse();

        // Stop if there was a syntax error.
        if (hadError) return;

        var resolver = new Resolver(interpreter);
        resolver.Resolve(statements);

        // Stop if there was a resolution error.
        if (hadError) return;

        interpreter.Interpret(statements);
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

    public static void Error(Token token, string message)
    {
        if (token.Type == TokenType.EOF)
        {
            Report(token.Line, " at end", message);
        }
        else
        {
            Report(token.Line, $" at `{token.Lexeme}`", message);
        }
    }

    public static void RuntimeError(RuntimeError error)
    {
        Console.WriteLine($"{error.Message}\n[line {error.Token.Line}]");
        hadRuntimeError = true;
    }
}
