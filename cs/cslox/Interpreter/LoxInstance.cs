﻿using System;
namespace cslox.Interpreter;

class LoxInstance
{
	private readonly LoxClass klass;

	public LoxInstance(LoxClass klass)
	{
		this.klass = klass;
	}

    public override string ToString()
    {
		return $"{klass.Name} instance";
    }
}

