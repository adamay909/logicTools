## In-Browser proof checker 



### The proof system

This proof checker is designed for working with the proof system Gentzen uses in his "Die Wiederspruchsfreiheit der reinen Zahlentheorie" (1936) that I use in my introductory formal logic course. A proof is stated as a series of sequents where each sequent must have exactly one formula as its succedent. It can be trivially converted to proofs in the style of Lemmon's *Beginning Logic* as well others influenced by him (e.g., Allen and Hand, *Logic Primer*). E.g., the following proof:
1. P ⊢ P
2. Q ⊢ Q
3. P,Q ⊢ P and Q

turns into the following in Lemmon's style:

1 (1) P  
2 (2) Q  
1,2 (3) P and Q  

What we do is replace the turnstile with the line  number, and replace formulas on the antecedent side of each sequent with the appropriate line numbers (of course, you need to add appropriate annotations). One thing Gentzen's allows is the use of placeholders on the antecedent side of a sequent which can be useful (there is an example of that below).

The proof system has 9 inference rules for sentential logic. In the system, you may:
- **Assumption Introduction (A)** Infer s ⊢ s.
- **Conjunction Introduction (∧I)** From Γ ⊢ s1 and ∆ ⊢ s2 , infer Γ, ∆ ⊢ s1 ∧ s2 .
- **Conjunction Elimination (∧E)** From Γ ⊢ s1 ∧ s2 , infer Γ ⊢ s1 as well as
Γ ⊢ s2 .
- **Disjunction Introduction (∨I)** From Γ ⊢ s1 , infer Γ ⊢ s1 ∨ s2 as well as
Γ ⊢ s2 ∨ s1 , for any s2 .
- **Disjunction Elimination (∨E)** From Γ ⊢ s1 ∨ s2 and s1 , ∆ ⊢ s3 and
s2 , Θ ⊢ s3 , infer Γ, ∆ , Θ ⊢ s3 .
- **Negation Introduction (¬I)** From Γ, s1 ⊢ s2 and ∆, s1 ⊢ ¬s2 , infer
Γ, ∆ ⊢ ¬s1 .
- **Negation Elimination (¬E)** From Γ ⊢ ¬¬s, infer Γ ⊢ s.
- **Conditional Introduction (⊃I)** From Γ , s1 ⊢ s2 , infer Γ ⊢ s1 ⊃ s2 .
- **Conditional Elimination (⊃E)** From Γ ⊢ s1 ⊃ s2 and ∆ ⊢ s1 , infer Γ, ∆ ⊢ s2 .

There are four more rules for predicate logic with quantifiers:
- **Universal Quantifier Introduction (∀I)** Given a constant κ, from Γ ⊢ φ(κ)
infer Γ ⊢ ∀υφ(υ), provided κ does not appear in any of the sentences
in Γ .
- **Universal Quantifier Elimination (∀E)** From Γ ⊢ ∀υφ(υ), infer Γ ⊢ φ(κ ),
for any constant κ.
- **Existential Quantifier Introduction (∃I)** Given a constant κ , infer from
Γ ⊢ φ( κ) to Γ ⊢ ∃υφ(υ).
- **Existential Quantifier Elimination (∃E)** From Γ ⊢ ∃υφ(υ) and ∆, φ (κ) ⊢ ψ,
infer Γ, ∆ ⊢ ψ , provided κ does not appear in any of Γ, ∆ , and ψ.

Finally, there are two rules for identity (Gentzen does not have these; he uses axioms governing the identity predicate):

- **Identity Introduction (=I)** For any constant κ, infer ⊢ κ=κ.
- **Identity Elimination (=E)** For any constants κ1 and κ2 , infer
⊢ (κ1=κ2 ∧ φ(κ1)) ⊃ φ∗(κ2) where φ∗(κ2) is any formula you can
obtain by substituting at least one instance of κ1 in φ with κ2 .


For annotations, the proof checker requires you to use the abbreviations given in parentheses.

There are three more rules for rewriting the antecedent side of sequents:
- You may reorder items within the antecedent of a sequent as you see fit.
- You may delete duplicate items within the antecedent of a sequent.
- You may add arbitrary items to the antecedent of a sequent.

These sequent rewrite rules have no names. When you use them, just give the relevant line numbers in the annotations. The proof checker allows you to use the first two silently. E.g., instead of:
1. P ⊢ P...A
2. Q ⊢ Q...A
3. P,Q ⊢ P∧Q...1,2,∧I
4. P,Q,P ⊢ (P∧Q)∧P...1,3,∧I
5. P,Q ⊢ (P∧Q)∧P...4

you may, but are not required to, write:
1. P ⊢ P...A
2. Q ⊢ Q...A
3. P,Q ⊢ P∧Q...1,2,∧I
4. P,Q ⊢ (P∧Q)∧P...1,3,∧I

But you are required to be explicit in the use of the third sequent rewrite rule. E.g., the following is not accepted by the proof checker:

1. &#x0393; ⊢P...premise
2. &#x0393; ⊢ Q⊃P...1,⊃I 

You must write the above as:
1. &#x0393; ⊢ P...premise
2. &#x0393;,Q ⊢ P...1
3. &#x0393; ⊢ Q⊃P...2,⊃I 

Notice the use of the keyword "premise" in the annotation. That is what you must use for a premise that is not an assumption (an assumption must take the form s ⊢ s). 


### Theorems

To make life easier, you can choose to allow the use of a few theorems. If you do use theorems, you must use their abbreviations in the annotations. The theorems are:

#### Theorems of Sentential Logic
- **Identity (ID)** ⊢ p ⊃ p
- **Non-Contradiction (NC)** ⊢ ¬(p ∧ ¬p)
- **Excluded Middle (EM)** ⊢ p ∨ ¬p
- **DeMorgan (DM)** ⊢ ¬(p ∨ q) ⊃ (¬p ∧ ¬q)
- **DeMorgan (DM)** ⊢ (¬p ∧ ¬q) ⊃ ¬(p ∨ q)
- **DeMorgan (DM)** ⊢ ¬(p ∧ q) ⊃ (¬p ∨ ¬q)
- **DeMorgan (DM)** ⊢ (¬p ∨ ¬q) ⊃ ¬(p ∧ q)
- **Implication (IM)** ⊢ (p ⊃ q) ⊃ (¬p ∨ q)
- **Elimination (EL)** ⊢ (p ∨ q) ⊃ (¬p ⊃ q)
- **Contraposition (CP)** ⊢ (p ⊃ q) ⊃ (¬q ⊃ ¬p)
- **Commutativity of Conjunction (CC)** ⊢ (p ∧ q) ⊃ (q ∧ p)
- **Commutativity of Disjunction (CD)** ⊢ (p ∨ q) ⊃ (q ∨ p)
- **Associativity of Conjunction (AC)** ⊢ [(p ∧ q) ∧ r] ⊃ [p ∧ (q ∧ r)]
- **Associativity of Conjunction (AC)** ⊢ [p ∧ (q ∧ r)] ⊃ [(p ∧ q) ∧ r]
- **Associativity of Disjunction (AD)** ⊢ [(p ∨ q) ∨ r] ⊃ [p ∨ (q ∨ r)]
- **Associativity of Disjunction (AD)** ⊢ [p ∨ (q ∨ r)] ⊃ [(p ∨ q) ∨ r]
- **Double Negation Introduction (DN)** ⊢ p ⊃ ¬¬p

#### Theorems of Predicate Logic
- **Quantifier Exchange (QE)** ⊢ ∃xFx ⊃ ¬∀x(¬Fx)
- **Quantifier Exchange (QE)** ⊢ ¬∀x(¬Fx) ⊃ ∃xFx
- **Quantifier Exchange (QE)** ⊢ ∀xFx ⊃ ¬∃x(¬Fx)
- **Quantifier Exchange (QE)** ⊢ ¬∃x(¬Fx) ⊃ ∀xFx
- **Quantifier Exchange (QE)** ⊢ ∃x(¬Fx) ⊃ ¬∀xFx
- **Quantifier Exchange (QE)** ⊢ ¬∀xFx ⊃ ∃x(¬Fx)
- **Quantifier Exchange (QE)** ⊢ ∀x(¬Fx) ⊃ ¬∃xFx
- **Quantifier Exchange (QE)** ⊢ ¬∃xFx ⊃ ∀x(¬Fx)

The proof checker will recognize instances of theorems. Here is an example of a use of EM:

1. Γ ⊢ P⊃Q...premise
2. ⊢ P∨¬P...EM
3. P ⊢ P...A
4. Γ,P ⊢ Q...1,3,⊃E
5. Γ,P ⊢ ¬P∨Q...4,∨I
6. ¬P ⊢ ¬P...A
7. ¬P ⊢ ¬P∨Q...6,∨I
8. Γ ⊢ ¬P∨Q...2,5,7,∨E


### The Proof Checker

 The proof checker checks you whether each line is in accordance with the proof system. But it does not check whether you have managed to show what you set out to show. You'll have to check that yourself---usually a matter of inspecting the last line of your derivation, possibly in combination with the premises.

You may have to go through several rounds of checking and fixing a derivation because the proof checker does not always list all the problems at once.

### The Editor

The editor is very primitive with limited functionality. It works like an old-school, keyboard-only editor. You can move the cursor around with the arrow keys, home and end for moving to the start or end of line, and delete and backspace should work more or less normally (sometimes more, sometimes less...). But no more advanced navigation around the input area, no positioning the cursor with your mouse,  and no copying and pasting and the like. Given the intended use, it should be enough (it works for me...).  You will need a physical keyboard. 


Apart from the above limitations, The editor is designed to be as transparent as possible: symbols should be easy to type and students should not have to worry about how what's on the screen corresponds to what they see in the course material. Some special key combinations are used to facilitate typing symbol. Check the help on how to input symbols. 


### Preservation of History

The proof checker will attempt to store the current state of the editor so that when you open the proof checker again, you will be presented with the last state of things before you quit (or the program crashed). The proof checker will also store a series of snapshots of the editor. This last happens whenever you clear the screen or make edits in the history. All history is stored in the browser as off-line data so how much history is stored for how long depends on your browser settings and the like. 

You can go back and forth in history using the 'back' and 'forward' buttons.
 
### Installation

This proof-checker is designed to run completely inside the browser so it is easy to host it yourself so long as you are able to host static websites. The docs folder of the GitHub repository contains all the files you need.

### Copyright

The Go, HTML, and CSS sources for this proof checker written by Masahiro Yamada. Licensed under the MIT License. You can get the source code at:
[https://github.com/adamay909/logicTools](https://github.com/adamay909/logicTools)
