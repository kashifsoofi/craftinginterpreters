// See https://aka.ms/new-console-template for more information
var outputDir = "../../../../cslox/Parser";
var types = new Dictionary<string, string[]>
{
    ["Assign"] = new[] { "Token Name", "Expr Value" },
    ["Binary"] = new[] { "Expr Left", "Token Operator", "Expr Right" },
    ["Call"] = new[] { "Expr Callee", "Token Paren", "List<Expr> Arguments" },
    ["Get"] = new[] { "Expr Object", "Token Name" },
    ["Grouping"] = new[] { "Expr Expression" },
    ["Literal"] = new[] { "object? Value" },
    ["Logical"] = new[] { "Expr Left", "Token Operator", "Expr Right" },
    ["Set"] = new[] { "Expr Object", "Token Name", "Expr Value" },
    ["This"] = new[] { "Token Keyword" },
    ["Unary"] = new[] { "Token Operator", "Expr Right" },
    ["Variable"] = new[] { "Token Name" },
};
DefineAst(outputDir, "Expr", types);

var stmtTypes = new Dictionary<string, string[]>
{
    ["Block"] = new[] { "List<Stmt> Statements" },
    ["Class"] = new[] { "Token Name", "List<Function> Methods" },
    ["ExpressionStmt"] = new[] { "Expr Expression" },
    ["Function"] = new[] { "Token Name", "List<Token> Parameters", "List<Stmt> Body" },
    ["If"] = new[] { "Expr Condition", "Stmt ThenBranch", "Stmt? ElseBranch" },
    ["Print"] = new[] { "Expr Expression" },
    ["Return"] = new[] { "Token Keyword", "Expr? Value" },
    ["Var"] = new[] { "Token Name", "Expr? Initializer" },
    ["While"] = new[] { "Expr Condition", "Stmt Body" }, 
};
DefineAst(outputDir, "Stmt", stmtTypes);

static void DefineAst(string outputDir, string baseName, Dictionary<string, string[]> types)
{
    var path = $"{outputDir}/{baseName}.cs";
    FileStream stream = new FileStream(path, FileMode.Create);
    using var writer = new StreamWriter(stream);
    
    writer.WriteLine("using cslox.Scanning;");
    writer.WriteLine("");
    writer.WriteLine("namespace cslox.Parser;");

    DefineVisitor(writer, baseName, types);

    writer.WriteLine("");
    writer.WriteLine($"abstract class {baseName}");
    writer.WriteLine("{");
    writer.WriteLine($"\tpublic abstract T Accept<T>(I{baseName}Visitor<T> visitor);");
    writer.WriteLine("}");

    foreach (var (className, fields) in types)
    {
        DefineType(writer, baseName, className, fields);
    }

}

static void DefineType(StreamWriter writer, string baseName, string className, string[] fields)
{
    var fieldList = "";
    var fieldNames = new List<string>();
    foreach (var field in fields)
    {
        var typeAndName = field.Split(" ");
        fieldNames.Add(typeAndName[1]);

        var parameterName = typeAndName[1].ToLower();
        if (parameterName == "operator" || parameterName == "object")
        {
            parameterName = $"@{parameterName}";
        }
        if (fieldList != "")
        {
            fieldList += ", ";
        }
        fieldList += $"{typeAndName[0]} {parameterName}";
    }

    writer.WriteLine("");
    writer.WriteLine($"class {className} : {baseName}");
    writer.WriteLine("{");
    writer.WriteLine($"\tpublic {className}({fieldList})");
    writer.WriteLine("\t{");
    foreach (var fieldName in fieldNames)
    {
        var parameterName = fieldName.ToLower();
        if (parameterName == "operator" || parameterName == "object")
        {
            parameterName = $"@{parameterName}";
        }

        writer.WriteLine($"\t\t{fieldName} = {parameterName};");
    }
    writer.WriteLine("\t}");
    writer.WriteLine("");
    foreach (var field in fields)
    {
        writer.WriteLine($"\tpublic {field} {{ get; }}");
    }
    writer.WriteLine("");
    writer.WriteLine($"\tpublic override T Accept<T>(I{baseName}Visitor<T> visitor)");
    writer.WriteLine("\t{");
    var methodName = $"Visit{className}";
    if (!methodName.EndsWith(baseName))
    {
        methodName += baseName;
    }
    writer.WriteLine($"\t\treturn visitor.{methodName}(this);");
    writer.WriteLine("\t}");
    writer.WriteLine("}");
}

static void DefineVisitor(StreamWriter writer, string baseName, Dictionary<string, string[]> types)
{
    writer.WriteLine("");
    writer.WriteLine($"interface I{baseName}Visitor<T>");
    writer.WriteLine("{");
    foreach (var (typeName, fields) in types)
    {
        var methodName = $"Visit{typeName}";
        if (!methodName.EndsWith(baseName))
        {
            methodName += baseName;
        }
        writer.WriteLine($"\tT {methodName}({typeName} {baseName.ToLower()});");
    }
    writer.WriteLine("}");
}