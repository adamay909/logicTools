### Online proof checker 

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
⊢ (κ1=κ2 ∧ φ(κ1)) ⊃ φ∗ (κ2) where φ∗ (κ2) is any formula you can
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

### Installation

This proof-checker is designed to run completely inside the browser so it is easy to host yourself so long as you are able to host static websites. You'll need the index.html file and the files in the assets/ folder.

### Copyright

The Go, HTML, and CSS sources for this proof checker written by Masahiro Yamada. Licensed under the MIT License. You can get the source code at:
