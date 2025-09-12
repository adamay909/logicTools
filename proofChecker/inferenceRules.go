package main

import "github.com/adamay909/logicTools/gentzen"

func setInferenceRules() {
	var infrules = []gentzen.InfRuleTemplate{
		gentzen.InfRuleTemplate{
			Name:     "A",
			FullName: "Assumption",
			Conclusion: []string{
				"s_1:s_1",
			},
		},

		gentzen.InfRuleTemplate{
			Name:     "premise",
			FullName: "Premise",
			Conclusion: []string{
				"/L:s_1",
			},
		},

		gentzen.InfRuleTemplate{
			Name:        "^I",
			FullName:    "Conjunction Introduction",
			DisplayName: "∧I",
			LatexName:   `\conjI`,
			RuleType:    gentzen.RTintroduction,
			Premises: []string{
				"/L_1:s_1",
				"/L_2:s_2",
			},
			Conclusion: []string{
				"/L_1,/L_2:^s_1s_2",
			},
		},

		gentzen.InfRuleTemplate{
			Name:        "^E",
			DisplayName: "∧E",
			LatexName:   `\conjE`,
			FullName:    "Conjunction Elimination",
			Premises: []string{
				"/L:^s_1s_2",
			},
			Conclusion: []string{
				"/L:s_1",
				"/L:s_2",
			},
		},

		gentzen.InfRuleTemplate{
			Name:        "vI",
			DisplayName: "∨I",
			LatexName:   `\disjI`,
			FullName:    "Disjunction Introduction",
			RuleType:    gentzen.RTintroduction,
			Premises: []string{
				"/L_1:s_1",
			},
			Conclusion: []string{
				"/L_1:vs_1s_2",
				"/L_1:vs_2s_1",
			},
		},

		gentzen.InfRuleTemplate{
			Name:        "vE",
			DisplayName: "∨E",
			LatexName:   `\disjE`,
			FullName:    "Disjunction Elimination",
			Premises: []string{
				"/L_1:vs_1s_2",
				"s_1,/L_2:s_3",
				"s_2,/L_3:s_3",
			},
			Conclusion: []string{
				"/L_1,/L_2,/L_3:s_3",
			},
		},

		gentzen.InfRuleTemplate{
			Name:        ">I",
			DisplayName: "⊃I",
			LatexName:   `\condI`,
			FullName:    "Conditional Introduction",
			RuleType:    gentzen.RTintroduction,
			Premises: []string{
				"s_1,/L:s_2",
			},
			Conclusion: []string{
				"/L:>s_1s_2",
			},
		},

		gentzen.InfRuleTemplate{
			Name:        ">E",
			DisplayName: "⊃E",
			LatexName:   `\condE`,
			FullName:    "Conditional Elimination",
			Premises: []string{
				"/L_1:>s_1s_2",
				"/L_2:s_1",
			},
			Conclusion: []string{
				"/L_1,/L_2:s_2",
			},
		},

		gentzen.InfRuleTemplate{
			Name:        "-I",
			DisplayName: "¬I",
			LatexName:   `\negI`,
			FullName:    "Negation Introduction",
			RuleType:    gentzen.RTintroduction,
			Premises: []string{
				"s_1,/L_1:s_2",
				"s_1,/L_2:-s_2",
			},
			Conclusion: []string{
				"/L_1,/L_2:-s_1",
			},
		},

		gentzen.InfRuleTemplate{
			Name:        "-E",
			DisplayName: "¬I",
			LatexName:   `\negE`,
			FullName:    "Negation Elimination",
			Premises: []string{
				"/L:--s_1",
			},
			Conclusion: []string{
				"/L:s_1",
			},
		},

		gentzen.InfRuleTemplate{
			Name:     "M",
			FullName: "Monotonicity",
			RuleType: gentzen.RTfrontmanipulation,
			Premises: []string{
				"/L_1:s_1",
			},
			Conclusion: []string{
				"/L_1,/L_2:s_1",
			},
		},

		gentzen.InfRuleTemplate{
			Name:        "rewrite",
			DisplayName: " ",
			FullName:    "Sequent Rewrite",
			RuleType:    gentzen.RTfrontmanipulation,
		},

		gentzen.InfRuleTemplate{
			Name:        "UE",
			DisplayName: "∀E",
			LatexName:   `\uniE`,
			FullName:    "Universal Quantifier Elimination",
			RuleType:    gentzen.RTpredicateLogic,
			Premises: []string{
				"/L_1:Ux_1Fx_1",
			},
			Conclusion: []string{
				"/L_1:Fa",
			},
		},

		gentzen.InfRuleTemplate{
			Name:        "UI",
			DisplayName: "∀I",
			LatexName:   `\uniI`,
			FullName:    "Universal Quantifier Introduction",
			RuleType:    gentzen.RTpredicateLogic | gentzen.RTintroduction,
			Premises: []string{
				"/L_1:Fa",
			},
			Conclusion: []string{
				"/L_1:UxFx",
			},
			Spec: "constants unique",
		},

		gentzen.InfRuleTemplate{
			Name:        "XI",
			DisplayName: "∃I",
			LatexName:   `\exI`,
			FullName:    "Existential Quantifier Introduction",
			RuleType:    gentzen.RTpredicateLogic | gentzen.RTintroduction,
			Premises: []string{
				"/L_1:Fa",
			},
			Conclusion: []string{
				"/L_1:XxFx",
			},
		},

		gentzen.InfRuleTemplate{
			Name:        "XE",
			DisplayName: "∃E",
			LatexName:   `\exE`,
			FullName:    "Existential Quantifier Elimination",
			RuleType:    gentzen.RTpredicateLogic,
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

	gentzen.ClearInferenceRules()

	for i := range infrules {
		gentzen.DeclareInferenceRule(infrules[i])
	}
}
