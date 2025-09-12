package gentzen

import "fmt"

type inferenceRule struct {
	name        string
	displayName string
	latexName   string
	fullName    string
	ruleType    uint8
	premises    []sequent //the most significant comes first
	conclusion  []sequent // alternative conclusions
	spec        string
}

/*InfRuleTemplate is used to declare inference rules.*/
type InfRuleTemplate struct {
	Name        string   //name of inference rule
	DisplayName string   // how the inference rule is to be displayed as text. Defaults to Name.
	LatexName   string   // latex commands for displaying the name. Defaults to DisplayName.
	FullName    string   //full, unabreviated name of rule. Defaults to Name.
	RuleType    uint8    //for available types, see below
	Premises    []string //the premises in sequent form.
	Conclusion  []string //the conclusion in sequent form. It's a slice to allow alternative conclusions.
	Spec        string   //further specifications. Currently, only "unique constants" is supported that is used to restrict ∀I and ∃E.
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

func addTheoremRules(derived bool) {

	var theorems = [][]string{
		{"Identity", "ID", "Cpp"},
		{"NonContradiction", "NC", "NKpNp"},
		{"Excluded Middle", "EM", "ApNp"},
		{"Contraposition", "CP", "CCpqCNqNp"},
		{"Implication", "IM", "CCpqANpq"},
		{"Elimination", "EL", "CApqCNpq"},
		{"DeMorgan", "DM", "CNApqKNpNq"},
		{"DeMorgan", "DM", "CNKpqANpNq"},
		{"DeMorgan", "DM", "CANpNqNKpq"},
		{"DeMorgan", "DM", "CKNpNqNApq"},
		{"Commutativity of Conjunction", "CC", "CKpqKqp"},
		{"Commutatitivity of Disjunction", "CD", "CApqAqp"},
		{"Associativity of Conjunction", "AC", "CKKpqrKpKqr"},
		{"Associativity of Conjunction", "AC", "CKpKqrKKpqr"},
		{"Associativity of Disjunction", "AD", "CAApqrApAqr"},
		{"Associativity of Disjunction", "AD", "CApAqrAApqr"},
		{"Double Negation Introduction", "DN", "CpNNp"},
	}

	for _, t := range theorems {

		nr := inferenceRule{
			name:     t[1],
			fullName: t[0],
			premises: nil,
		}

		c, err := mkSequent(":"+t[2], O_Polish)
		if err != nil {
			panic("definition of theorem " + t[0] + " is bad")
		}
		nr.conclusion = []sequent{c}

		infrule[nr.name] = nr
	}

	if !derived {
		return
	}

	for _, t := range theorems {
		n := Parse(t[2], !allowGreekUpper)
		if n.Val.connective != Cond {
			continue
		}

		nr := inferenceRule{
			name:     t[1] + "R",
			fullName: t[0] + " (derived rule)",
		}

		p, _ := mkSequent("/L_1:"+n.Child(0).String(), O_Polish)
		c, _ := mkSequent("/L_1:"+n.Child(1).String(), O_Polish)

		nr.premises = []sequent{p}
		nr.conclusion = []sequent{c}

		infrule[nr.name] = nr
	}

}

/* SetStandardInferenceRules sets the inference rules to be those in the text book.*/
func SetStandardInferenceRules() {
	var infrules = []InfRuleTemplate{
		InfRuleTemplate{
			Name:     "A",
			FullName: "Assumption",
			Conclusion: []string{
				"s_1:s_1",
			},
		},

		InfRuleTemplate{
			Name:     "premise",
			FullName: "Premise",
			Conclusion: []string{
				"/L:s_1",
			},
		},

		InfRuleTemplate{
			Name:        "^I",
			FullName:    "Conjunction Introduction",
			DisplayName: "∧I",
			LatexName:   `\conjI`,
			RuleType:    RTintroduction,
			Premises: []string{
				"/L_1:s_1",
				"/L_2:s_2",
			},
			Conclusion: []string{
				"/L_1,/L_2:Ks_1s_2",
			},
		},

		InfRuleTemplate{
			Name:        "^E",
			DisplayName: "∧E",
			LatexName:   `\conjE`,
			FullName:    "Conjunction Elimination",
			Premises: []string{
				"/L:Ks_1s_2",
			},
			Conclusion: []string{
				"/L:s_1",
				"/L:s_2",
			},
		},

		InfRuleTemplate{
			Name:        "vI",
			DisplayName: "∨I",
			LatexName:   `\disjI`,
			FullName:    "Disjunction Introduction",
			RuleType:    RTintroduction,
			Premises: []string{
				"/L_1:s_1",
			},
			Conclusion: []string{
				"/L_1:As_1s_2",
				"/L_1:As_2s_1",
			},
		},

		InfRuleTemplate{
			Name:        "vE",
			DisplayName: "∨E",
			LatexName:   `\disjE`,
			FullName:    "Disjunction Elimination",
			Premises: []string{
				"/L_1:As_1s_2",
				"s_1,/L_2:s_3",
				"s_2,/L_3:s_3",
			},
			Conclusion: []string{
				"/L_1,/L_2,/L_3:s_3",
			},
		},

		InfRuleTemplate{
			Name:        ">I",
			DisplayName: "⊃I",
			LatexName:   `\condI`,
			FullName:    "Conditional Introduction",
			RuleType:    RTintroduction,
			Premises: []string{
				"s_1,/L:s_2",
			},
			Conclusion: []string{
				"/L:Cs_1s_2",
			},
		},

		InfRuleTemplate{
			Name:        ">E",
			DisplayName: "⊃E",
			LatexName:   `\condE`,
			FullName:    "Conditional Elimination",
			Premises: []string{
				"/L_1:Cs_1s_2",
				"/L_2:s_1",
			},
			Conclusion: []string{
				"/L_1,/L_2:s_2",
			},
		},

		InfRuleTemplate{
			Name:        "-I",
			DisplayName: "¬I",
			LatexName:   `\negI`,
			FullName:    "Negation Introduction",
			RuleType:    RTintroduction,
			Premises: []string{
				"s_1,/L_1:s_2",
				"s_1,/L_2:Ns_2",
			},
			Conclusion: []string{
				"/L_1,/L_2:Ns_1",
			},
		},

		InfRuleTemplate{
			Name:        "-E",
			DisplayName: "¬I",
			LatexName:   `\negE`,
			FullName:    "Negation Elimination",
			Premises: []string{
				"/L:NNs_1",
			},
			Conclusion: []string{
				"/L:s_1",
			},
		},

		InfRuleTemplate{
			Name:     "M",
			FullName: "Monotonicity",
			RuleType: RTfrontmanipulation,
			Premises: []string{
				"/L_1:s_1",
			},
			Conclusion: []string{
				"/L_1,/L_2:s_1",
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
			Premises: []string{
				"/L_1:Ux_1Fx_1",
			},
			Conclusion: []string{
				"/L_1:Fa",
			},
		},

		InfRuleTemplate{
			Name:        "UI",
			DisplayName: "∀I",
			LatexName:   `\uniI`,
			FullName:    "Universal Quantifier Introduction",
			RuleType:    RTpredicateLogic | RTintroduction,
			Premises: []string{
				"/L_1:Fa",
			},
			Conclusion: []string{
				"/L_1:UxFx",
			},
			Spec: "constants unique",
		},

		InfRuleTemplate{
			Name:        "XI",
			DisplayName: "∃I",
			LatexName:   `\exI`,
			FullName:    "Existential Quantifier Introduction",
			RuleType:    RTpredicateLogic | RTintroduction,
			Premises: []string{
				"/L_1:Fa",
			},
			Conclusion: []string{
				"/L_1:XxFx",
			},
		},

		InfRuleTemplate{
			Name:        "XE",
			DisplayName: "∃E",
			LatexName:   `\exE`,
			FullName:    "Existential Quantifier Elimination",
			RuleType:    RTpredicateLogic,
			Premises: []string{
				"/L_1:XxFx",
				"/L_2,Fa:Gb",
			},
			Conclusion: []string{
				"/L_1,/L_2:Gb",
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
	if ir.isPLrule() {
		SetPL(true)
	} else {
		SetPL(false)
	}

	for _, p := range i.Premises {
		seq, err := mkSequent(p, O_Polish)
		if err != nil {
			fmt.Println(err)
			panic("defintitions of inference rule " + i.FullName + " is bad")
		}
		ir.premises = append(ir.premises, seq)
	}

	for _, c := range i.Conclusion {
		seq, err := mkSequent(c, O_Polish)
		if err != nil {
			fmt.Println(err)
			panic("defintitions of inference rule " + i.FullName + " is bad")
		}
		ir.conclusion = append(ir.conclusion, seq)
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
