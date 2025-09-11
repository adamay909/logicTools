package gentzen

import "fmt"

type inferenceRule struct {
	name       string
	fullName   string
	ruleType   uint8
	premises   []sequent //the most significant comes first
	conclusion []sequent // alternative conclusions
	spec       string
}

type InfRuleTemplate struct {
	Name       string
	FullName   string
	ruleType   uint8
	Premises   []string
	Conclusion []string
	Spec       string
}

const (
	predicatelogicrule uint8 = 1 << iota
	introductionrule
	frontmanipulationrule
)

var infrule map[string]inferenceRule

func (i inferenceRule) isPLrule() bool {
	return i.ruleType&predicatelogicrule != 0
}

func (i inferenceRule) isFrontManipulationRule() bool {
	return i.ruleType&frontmanipulationrule != 0
}

func (i inferenceRule) isIntroductionRule() bool {
	return i.ruleType&introductionrule != 0
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
			Name:     "^I",
			FullName: "Conjunction Introduction",
			ruleType: introductionrule,
			Premises: []string{
				"/L_1:s_1",
				"/L_2:s_2",
			},
			Conclusion: []string{
				"/L_1,/L_2:Ks_1s_2",
			},
		},

		InfRuleTemplate{
			Name:     "^E",
			FullName: "Conjunction Elimination",
			Premises: []string{
				"/L:Ks_1s_2",
			},
			Conclusion: []string{
				"/L:s_1",
				"/L:s_2",
			},
		},

		InfRuleTemplate{
			Name:     "vI",
			FullName: "Disjunction Introduction",
			ruleType: introductionrule,
			Premises: []string{
				"/L_1:s_1",
			},
			Conclusion: []string{
				"/L_1:As_1s_2",
				"/L_1:As_2s_1",
			},
		},

		InfRuleTemplate{
			Name:     "vE",
			FullName: "Disjunction Elimination",
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
			Name:     ">I",
			FullName: "Conditional Introduction",
			ruleType: introductionrule,
			Premises: []string{
				"s_1,/L:s_2",
			},
			Conclusion: []string{
				"/L:Cs_1s_2",
			},
		},

		InfRuleTemplate{
			Name:     ">E",
			FullName: "Conditional Elimination",
			Premises: []string{
				"/L_1:Cs_1s_2",
				"/L_2:s_1",
			},
			Conclusion: []string{
				"/L_1,/L_2:s_2",
			},
		},

		InfRuleTemplate{
			Name:     "-I",
			FullName: "Negation Introduction",
			ruleType: introductionrule,
			Premises: []string{
				"s_1,/L_1:s_2",
				"s_1,/L_2:Ns_2",
			},
			Conclusion: []string{
				"/L_1,/L_2:Ns_1",
			},
		},

		InfRuleTemplate{
			Name:     "-E",
			FullName: "Negation Elimination",
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
			ruleType: frontmanipulationrule,
			Premises: []string{
				"/L_1:s_1",
			},
			Conclusion: []string{
				"/L_1,/L_2:s_1",
			},
		},

		InfRuleTemplate{
			Name:     "rewrite",
			FullName: "Sequent Rewrite",
			ruleType: frontmanipulationrule,
		},

		InfRuleTemplate{
			Name:     "UE",
			FullName: "Universal Quantifier Elimination",
			ruleType: predicatelogicrule,
			Premises: []string{
				"/L_1:Ux_1Fx_1",
			},
			Conclusion: []string{
				"/L_1:Fa",
			},
		},

		InfRuleTemplate{
			Name:     "UI",
			FullName: "Universal Quantifier Introduction",
			ruleType: predicatelogicrule | introductionrule,
			Premises: []string{
				"/L_1:Fa",
			},
			Conclusion: []string{
				"/L_1:UxFx",
			},
			Spec: "constants unique",
		},

		InfRuleTemplate{
			Name:     "XI",
			FullName: "Existential Quantifier Introduction",
			ruleType: predicatelogicrule | introductionrule,
			Premises: []string{
				"/L_1:Fa",
			},
			Conclusion: []string{
				"/L_1:XxFx",
			},
		},

		InfRuleTemplate{
			Name:     "XE",
			FullName: "Existential Quantifier Elimination",
			ruleType: predicatelogicrule,
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

func DeclareInferenceRule(i InfRuleTemplate) {

	rpl := oPL
	defer func() {
		oPL = rpl
	}()

	ir := inferenceRule{
		name:     i.Name,
		fullName: i.FullName,
		ruleType: i.ruleType,
		spec:     i.Spec,
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
