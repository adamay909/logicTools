package gentzen

import "fmt"

type inferenceRule struct {
	name        string
	displayName string
	latexName   string
	fullName    string
	ruleType    uint8
	patterns    [][]sequent //the most significant comes first
	spec        string
}

/*InfRuleTemplate is used to declare inference rules.*/
type InfRuleTemplate struct {
	Name        string     //name of inference rule
	DisplayName string     // how the inference rule is to be displayed as text. Defaults to Name.
	LatexName   string     // latex commands for displaying the name. Defaults to DisplayName.
	FullName    string     //full, unabreviated name of rule. Defaults to Name.
	RuleType    uint8      //for available types, see below
	Patterns    [][]string //the premises in sequent form.
	Spec        string     //further specifications. Currently, only "unique constants" is supported that is used to restrict ∀I and ∃E.
}

// The following constants define available inference rule types.
const (
	RTpredicateLogic    uint8 = 1 << iota //rule is predicate logic specific.
	RTintroduction                        //rule is an introduction rule
	RTfrontmanipulation                   //rule is for manipulating front formulas. NOTE: currently unsupported.
)

var infrule map[string]inferenceRule

func (i inferenceRule) isPLrule() bool {
	return i.ruleType&RTpredicateLogic != 0
}

func (i inferenceRule) isFrontManipulationRule() bool {
	return i.ruleType&RTfrontmanipulation != 0
}

func (i inferenceRule) isIntroductionRule() bool {
	return i.ruleType&RTintroduction != 0
}

/* SetStandardInferenceRules sets the inference rules to be those in the text book.*/
func SetStandardInferenceRules() {
	var infrules = []InfRuleTemplate{
		InfRuleTemplate{
			Name:     "A",
			FullName: "Assumption",
			Patterns: [][]string{
				[]string{
					"s_1:s_1",
				},
			},
		},

		InfRuleTemplate{
			Name:     "premise",
			FullName: "Premise",
			Patterns: [][]string{
				[]string{
					"/L:s_1",
				},
			},
		},

		InfRuleTemplate{
			Name:        "^I",
			FullName:    "Conjunction Introduction",
			DisplayName: "∧I",
			LatexName:   `\conjI`,
			RuleType:    RTintroduction,
			Patterns: [][]string{
				[]string{
					"/L_1:s_1",
					"/L_2:s_2",
					"/L_1,/L_2:Ks_1s_2",
				},
			},
		},

		InfRuleTemplate{
			Name:        "^E",
			DisplayName: "∧E",
			LatexName:   `\conjE`,
			FullName:    "Conjunction Elimination",
			Patterns: [][]string{
				[]string{
					"/L:Ks_1s_2",
					"/L:s_1",
					"/L:s_2",
				},
			},
		},

		InfRuleTemplate{
			Name:        "vI",
			DisplayName: "∨I",
			LatexName:   `\disjI`,
			FullName:    "Disjunction Introduction",
			RuleType:    RTintroduction,
			Patterns: [][]string{
				[]string{
					"/L_1:s_1",
					"/L_1:As_1s_2",
				},
				[]string{
					"/L_1:s_1",
					"/L_1:As_2s_1",
				},
			},
		},

		InfRuleTemplate{
			Name:        "vE",
			DisplayName: "∨E",
			LatexName:   `\disjE`,
			FullName:    "Disjunction Elimination",
			Patterns: [][]string{
				[]string{
					"/L_1:As_1s_2",
					"s_1,/L_2:s_3",
					"s_2,/L_3:s_3",
					"/L_1,/L_2,/L_3:s_3",
				},
			},
		},

		InfRuleTemplate{
			Name:        ">I",
			DisplayName: "⊃I",
			LatexName:   `\condI`,
			FullName:    "Conditional Introduction",
			RuleType:    RTintroduction,
			Patterns: [][]string{
				[]string{
					"s_1,/L:s_2",
					"/L:Cs_1s_2",
				},
			},
		},

		InfRuleTemplate{
			Name:        ">E",
			DisplayName: "⊃E",
			LatexName:   `\condE`,
			FullName:    "Conditional Elimination",
			Patterns: [][]string{
				[]string{
					"/L_1:Cs_1s_2",
					"/L_2:s_1",
					"/L_1,/L_2:s_2",
				},
			},
		},

		InfRuleTemplate{
			Name:        "-I",
			DisplayName: "¬I",
			LatexName:   `\negI`,
			FullName:    "Negation Introduction",
			RuleType:    RTintroduction,
			Patterns: [][]string{
				[]string{
					"s_1,/L_1:s_2",
					"s_1,/L_2:Ns_2",
					"/L_1,/L_2:Ns_1",
				},
			},
		},

		InfRuleTemplate{
			Name:        "-E",
			DisplayName: "¬I",
			LatexName:   `\negE`,
			FullName:    "Negation Elimination",
			Patterns: [][]string{
				[]string{
					"/L:NNs_1",
					"/L:s_1",
				},
			},
		},

		InfRuleTemplate{
			Name:     "M",
			FullName: "Monotonicity",
			RuleType: RTfrontmanipulation,
			Patterns: [][]string{
				[]string{
					"/L_1:s_1",
					"/L_1,/L_2:s_1",
				},
			},
		},

		InfRuleTemplate{
			Name:        "rewrite",
			DisplayName: " ",
			FullName:    "Sequent Rewrite",
			RuleType:    RTfrontmanipulation,
		},

		InfRuleTemplate{
			Name:        "UE",
			DisplayName: "∀E",
			LatexName:   `\uniE`,
			FullName:    "Universal Quantifier Elimination",
			RuleType:    RTpredicateLogic,
			Patterns: [][]string{
				[]string{
					"/L_1:Ux_1Fx_1",
					"/L_1:Fa",
				},
			},
		},

		InfRuleTemplate{
			Name:        "UI",
			DisplayName: "∀I",
			LatexName:   `\uniI`,
			FullName:    "Universal Quantifier Introduction",
			RuleType:    RTpredicateLogic | RTintroduction,
			Patterns: [][]string{
				[]string{
					"/L_1:Fa",
					"/L_1:UxFx",
				},
			},
			Spec: "constants unique",
		},

		InfRuleTemplate{
			Name:        "XI",
			DisplayName: "∃I",
			LatexName:   `\exI`,
			FullName:    "Existential Quantifier Introduction",
			RuleType:    RTpredicateLogic | RTintroduction,
			Patterns: [][]string{
				[]string{
					"/L_1:Fa",
					"/L_1:XxFx",
				},
			},
		},

		InfRuleTemplate{
			Name:        "XE",
			DisplayName: "∃E",
			LatexName:   `\exE`,
			FullName:    "Existential Quantifier Elimination",
			RuleType:    RTpredicateLogic,
			Patterns: [][]string{
				[]string{
					"/L_1:XxFx",
					"/L_2,Fa:Gb",
					"/L_1,/L_2:Gb",
				},
			},
			Spec: "constants unique",
		},
	}

	for i := range infrules {
		DeclareInferenceRule(infrules[i])
	}
}

/*DeclareInferenceRule makes the rule specified by i available to the proof checking backend.*/
func DeclareInferenceRule(i InfRuleTemplate) {

	rpl := oPL
	defer func() {
		oPL = rpl
	}()

	ir := inferenceRule{
		name:     i.Name,
		fullName: i.FullName,
		ruleType: i.RuleType,
		spec:     i.Spec,
	}

	if i.DisplayName == "" {
		ir.displayName = i.Name
	} else {
		ir.displayName = i.DisplayName
	}

	if i.LatexName == "" {
		ir.latexName = ir.displayName
	} else {
		ir.latexName = i.LatexName
	}

	rPL := oPL

	defer func() {
		SetPL(rPL)
	}()

	if ir.isPLrule() {
		SetPL(true)
	} else {
		SetPL(false)
	}

	for _, pattern := range i.Patterns {

		var irpattern []sequent
		for _, p := range pattern {
			seq, err := mkSequent(p, O_Polish)
			if err != nil {
				fmt.Println(err)
				panic("defintitions of inference rule " + i.FullName + " is bad")
			}
			irpattern = append(irpattern, seq)
		}
		ir.patterns = append(ir.patterns, irpattern)
	}

	infrule[ir.name] = ir
}

/*ClearInferenceRules removes all inference rules.*/
func ClearInferenceRules() {
	clear(infrule)
}

/*DeleteInferenceRule removes the named inference rule.*/
func DeleteInferenceRule(name string) {
	delete(infrule, name)
}
