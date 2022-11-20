namespace cslox.Scanning;

using static TokenType;

class Scanner
{
	private readonly string source;
	private List<Token> tokens { get; } = new List<Token>();
	private int start = 0;
	private int current = 0;
	private int line = 1;

	private static Dictionary<string, TokenType> Keywords = new Dictionary<string, TokenType>
	{
		["and"] = AND,
		["class"] = CLASS,
		["else"] = ELSE,
		["false"] = FALSE,
		["for"] = FOR,
		["fun"] = FUN,
		["if"] = IF,
		["nil"] = NIL,
		["or"] = OR,
		["print"] = PRINT,
		["return"] = RETURN,
		["super"] = SUPER,
		["this"] = THIS,
		["true"] = TRUE,
		["var"] = VAR,
		["while"] = WHILE,
	};

	public Scanner(string source)
	{
		this.source = source;
	}

	private bool IsAtEnd() => current >= source.Length;
	private char Advance() => source[current++];

	public List<Token> ScanTokens()
	{
		while (!IsAtEnd())
		{
			start = current;
			ScanToken();
		}

		tokens.Add(new Token(EOF, "", null, line));
		return tokens;
	}

	private void AddToken(TokenType type) => AddToken(type, null);

	private void AddToken(TokenType type, object? literal)
	{
		var text = source.Substring(start, current - start);
		tokens.Add(new Token(type, text, literal, line));
	}

	private bool Match(char expected)
	{
		if (IsAtEnd())
		{
			return false;
		}

		if (source[current] != expected)
		{
			return false;
		}

		current++;
		return true;
	}

	private char Peek()
	{
		if (IsAtEnd())
		{
			return '\0';
		}

		return source[current];
	}

	private char PeekNext()
	{
		if (current + 1 > source.Length)
		{
			return '\0';
		}
		return source[current + 1];
	}

	private void ScanString()
	{
		while (Peek() != '"' && !IsAtEnd())
		{
			if (Peek() == '\n')
			{
				line++;
			}
			Advance();
		}

		if (IsAtEnd())
		{
			Lox.Error(line, "Unterminated string.");
			return;
		}

		// The closing ".
		Advance();

		// Trim the surrounding quotes.
		string value = source.Substring(start + 1, current - start - 1);
		AddToken(STRING, value);
	}

	private bool IsDigit(char c) => (c >= '0' && c <= '9');

	private void ScanNumber()
	{
		while (IsDigit(Peek()))
		{
			Advance();
		}

		if (Peek() == '.' && IsDigit(PeekNext()))
		{
			// Consume the "."
			Advance();

			while (IsDigit(Peek()))
			{
				Advance();
			}
		}

		AddToken(NUMBER, Double.Parse(source.Substring(start, current - start)));
	}

	private bool IsAlpha(char c) =>
		(c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_';

	private bool IsAlphaNumeric(char c) => IsAlpha(c) || IsDigit(c);

	private void ScanIdentifier()
	{
		while (IsAlphaNumeric(Peek()))
		{
			Advance();
		}

		var text = source.Substring(start, current - start);
		var type = IDENTIFIER;
		if (Keywords.ContainsKey(text))
		{
			type = Keywords[text];
		}

		AddToken(type);
	}

	private void ScanBlockComment()
	{
		var depth = 0;

        while (!IsAtEnd())
		{
            var c = Advance();
            switch (c)
            {
                case '/':
                    if (Match('*'))
                    {
                        depth++;
                    }
                    break;
                case '*':
                    if (Match('/'))
                    {
                        if (depth == 0)
                        {
                            return;
                        }
                        depth--;
                    }
                    break;
                case '\n':
                    line++;
                    break;
            }
        }
		Lox.Error(line, "Unclosed block comment");
    }

    private void ScanToken()
	{
		char c = Advance();
		switch (c)
		{
			case '(': AddToken(LEFT_PAREN); break;
			case ')': AddToken(RIGHT_PAREN); break;
			case '{': AddToken(LEFT_BRACE); break;
			case '}': AddToken(RIGHT_BRACE); break;
			case ',': AddToken(COMMA); break;
			case '.': AddToken(DOT); break;
			case '-': AddToken(MINUS); break;
			case '+': AddToken(PLUS); break;
			case ';': AddToken(SEMICOLON); break;
			case '*': AddToken(STAR); break;
			case '!':
				AddToken(Match('=') ? BANG_EQUAL : BANG);
				break;
			case '=':
				AddToken(Match('=') ? EQUAL_EQUAL : EQUAL);
				break;
			case '<':
				AddToken(Match('=') ? LESS_EQUAL : LESS);
				break;
			case '>':
				AddToken(Match('=') ? GREATER_EQUAL : GREATER);
				break;
			case '/':
				if (Match('/'))
				{
					// A comment goes until the end of the line.
					while (Peek() != '\n' && !IsAtEnd())
					{
						Advance();
					}
                }
				else if (Match('*'))
				{
					ScanBlockComment();
				}
				else
				{
					AddToken(SLASH);
				}
				break;
			case ' ':
			case '\t':
			case '\r':
				// Ignore whitespace
				break;
			case '\n':
				line++;
				break;
			case '"': ScanString(); break;

            default:
				if (IsDigit(c))
				{
					ScanNumber();
				}
				else if (IsAlpha(c))
				{
					ScanIdentifier();
				}
				else
				{
					Lox.Error(line, "Unexpected character.");
				}
				break;
		}
	}
}