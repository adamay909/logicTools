package gentzen

import "github.com/adamay909/logicTools/ju"

func init() {

	ju.SetStringFunc[syntaxNode](stringFunc)

	infrule = make(map[string]inferenceRule)

	SetStandardInferenceRules()

}
