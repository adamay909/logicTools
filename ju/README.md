## Overview

Package ju provides basic generic definitions and methods for dealing with simple trees: each node has at most one parent, and zero or more child nodes, and each node holds a value of type T.

Given a value type T, you can define a node type that holds values of type T as Node\[T]. A node of type Node\[T] has at most one parent and zero or more child nodes, as well at most one left-sibling and at most one right-sibling.  The easiest way to use this package is to import it and do something like

	type mynodevalues struct {}
	type mynode = ju.Node[mynodevalues]

After this, you can think in terms of mynode with the methods defined in this package. Because mynode is a type alias for Node\[mynodevalues], you cannot define your own methods for mynode. You could try

	type mynode ju.Node[mynodevalues]

and deal with complications using reflection and unsafe stuff. I find it easier to use type aliases. Not having methods is not a big deal, and there are some helpers defined in this package that can use functions as arguments to methods (see package documentation).

Given two types T1 and T2, Node\[T1] and Node\[T2] are distinct types with the usual implications given Go's type safety rules.  In particular, you cannot have a Node\[T1] and a Node\[T2] in a single tree.

But notice that the type T you use to instantiate Node\[T] can be an interface type. You can let T be an interface implemented by various types. That way, you can have a single type Node\[T] that stores different value types in its Val field. You can retrieve the concrete type stored in Val using type assertion or, if you must, reflection.

This package is not tuned for performance or anything like that. The main purpose is to provide a means to eliminate writing duplicate code for simple trees.

### About the name of this package

In written Japanese there are two Chinese characters (Kanji) that can be used for tree: 木 and 樹. Both can be pronounced 'ki,' but only the latter can be  pronounced 'ju'. Crucially, the latter character means a living tree: something that can grow by adding branches, for instance. The former character (木) like the Japanse word 'ki' can also mean wood, as in the material. This package deals with trees in the sense of 'ju' (樹): something that can grow, but also be truncated without stopping being a 'ju'.


