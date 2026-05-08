# AI Detection Guide: Patterns That Give Away Machine-Generated Fiction

A practical reference for identifying and eliminating AI tells from novel prose.
Compiled from 20+ sources: Reddit writing communities, editor accounts, academic
stylometry research, book reviews, professional writers, and writing craft analysis.

---

## Summary: Most Common AI Tells (Ranked by Citation Frequency)

These are the patterns most frequently cited across sources, ordered by how often
they appear in complaints, reviews, and detection research.

1. **Monotonous sentence rhythm** -- uniform length and cadence, no variation
2. **Overused "AI vocabulary"** -- delve, tapestry, testament, vibrant, etc.
3. **Telling instead of showing** -- stated emotions, exposed subtext
4. **Em dash addiction** -- compulsive overuse for parenthetical insertions
5. **Purple prose / metaphor overload** -- stacked figurative language with no restraint
6. **Generic sensory details** -- "gentle breeze," "blooming flowers," "golden light"
7. **All characters sound the same** -- identical dialogue voice regardless of background
8. **Negative parallelism** -- "It's not X, it's Y" sentence construction
9. **Rule of three abuse** -- triplet adjectives/nouns in nearly every description
10. **Relentless positivity** -- avoids darkness, ambiguity, and genuine discomfort
11. **Cliche action beats** -- "heart pounded," "jaw clenched," "eyes widened" on repeat
12. **Consistency failures** -- character details, timeline, worldbuilding contradictions
13. **Missing subtext** -- characters say exactly what they mean, always
14. **Perfect grammar** -- no fragments, no sentence-starting "And" or "But," no contractions
15. **Formulaic paragraph structure** -- uniform length, same opening patterns

---

## Detailed Analysis of Each Tell

### 1. Monotonous Sentence Rhythm

**What it is:** AI-generated prose tends toward sentences of similar length and
cadence. The rhythm never shifts. Every sentence lands with the same weight. There
are no short punches. No long, winding sentences that build and defer their meaning
across multiple clauses before finally arriving at something unexpected.

**Why it gives away AI:** Human "burstiness" -- the natural variation in sentence
length and complexity -- is one of the most measurable differences between human and
AI text. Academic research (StyloAI, 2024) found that low burstiness is a primary
statistical indicator of machine generation. AI text maintains "uniform, consistently
low perplexity throughout" while human writing shows spikes and valleys.

**AI example:**
> The ship drifted through the void. Marcus checked the navigation console. The
> readings showed nothing unusual. He leaned back in his chair. The silence pressed
> against him. He thought about what lay ahead.

**Human prose does this instead:** Mixes sentence lengths deliberately. A paragraph
might open with a two-word fragment, follow with a complex compound sentence, then
drop to a medium-length declarative. The rhythm serves the emotional content -- short
sentences for tension, longer ones for reflection or description.

**Fix:** Read paragraphs aloud. If every sentence takes roughly the same time to say,
rewrite. Vary sentence length from 3 words to 30+. Use fragments. Use run-ons
(intentionally). Let rhythm follow emotion.

---

### 2. Overused "AI Vocabulary"

**What it is:** A specific set of words that appear at dramatically higher rates in
LLM output than in human writing. A Helsinki study documented measurable surges in
certain words post-ChatGPT.

**The worst offenders (fiction-relevant):**

*Nouns/descriptors:* tapestry, testament, landscape, realm, journey, beacon, enigma,
labyrinth, crucible, symphony, gossamer, nexus, mosaic, cascade

*Verbs:* delve, embark, unveil, navigate, foster, harness, illuminate, resonate,
reverberate, transcend, underscore, spearhead

*Adjectives:* vibrant, pivotal, profound, meticulous, seamless, compelling,
transformative, nuanced, robust, multifaceted, unwavering, invaluable

*Adverbs:* seamlessly, meticulously, profoundly, relentlessly, notably, arguably

**Why it gives away AI:** These words are statistically overrepresented in training
data reward signals. LLMs learned that these words score well on "helpful and
articulate" metrics. When three or more appear in a single paragraph, experienced
readers flag it immediately. Wikipedia editors maintain a specific watchlist of
these terms.

**AI example:**
> The vibrant tapestry of cultures in the station's lower decks was a testament
> to humanity's unwavering resilience. Marcus delved into the multifaceted dynamics
> of the crew, navigating the nuanced interplay between the factions.

**Fix:** Search-and-destroy. Grep your manuscript for every word on this list. Replace
with specific, concrete language. "Vibrant tapestry of cultures" becomes "the Mandarin
signs next to the Arabic graffiti next to the Portuguese prayer cards taped over a
ventilation grate."

---

### 3. Telling Instead of Showing (Exposed Subtext)

**What it is:** AI states emotions directly rather than dramatizing them through
action, dialogue, and sensory detail. Alexander Wales calls this "exposed subtext" --
the AI "plops the subtext into the text." CRAFT Literary's analysis puts it bluntly:
AI "cannot be oblique or indirect, cannot let details speak for themselves."

**Why it gives away AI:** Fiction works through implication. When a character
"felt a surge of grief mixed with anger and guilt," that's a therapist's case notes,
not a novel. Readers detect this because it reads like summary rather than scene.
The Poets & Writers analysis found AI writing is "gracelessly heavy-handed and
dripping with sentimentality" -- it cannot resist explaining its own thematic
significance.

**AI example:**
> Elena felt overwhelmed by a complex mixture of sadness and determination. She
> knew that the journey ahead would test her in ways she couldn't imagine, but her
> resolve was unshakeable. The weight of her responsibility pressed down on her,
> yet she found strength in her commitment to those she had lost.

**Human prose does this instead:** Shows the emotion through specific physical action,
thought, or sensory perception. Instead of "she felt grief," a human writer might
write about the character noticing a coffee mug left by someone who died, running her
thumb along the chip in its rim, then setting it back exactly where it was.

**Fix:** For every paragraph, ask: "Am I naming the emotion or dramatizing it?" Cut
every sentence that labels a feeling. Replace with a concrete moment that makes the
reader feel it themselves.

---

### 4. Em Dash Addiction

**What it is:** AI uses em dashes at 3-5x the rate of human writers. The pattern is
so consistent it appears on virtually every list of AI tells.

**Why it gives away AI:** Human writers use em dashes occasionally for emphasis or
interruption. AI uses them as its default parenthetical construction -- inserting
them everywhere -- because the training data rewarded this as "stylish" punctuation.
When every other sentence has an em-dash insertion, it creates a distinctive mechanical
rhythm. One analysis (Hunting the Muse) noted that ChatGPT specifically uses them
"between sentences instead of having proper spacing."

**AI example:**
> The station -- once a beacon of human achievement -- now drifted in silence. Marcus
> -- who had spent three years aboard -- could feel the change. The corridors --
> empty and cold -- echoed with memories of a time that would never return.

**Fix:** Count em dashes per page. If you have more than 2-3 per page, you have a
problem. Replace most with commas, periods, or parentheses. Reserve em dashes for
genuine interruptions or dramatic emphasis, max once or twice per chapter.

---

### 5. Purple Prose / Metaphor Overload

**What it is:** AI stacks metaphors on top of each other without restraint. One
reviewer of OpenAI's creative writing model counted "between 17 and 38 metaphors"
in a single short piece. Wales identifies this as a core failure mode: "excessive
metaphors and overwrought descriptions stacked together."

**Why it gives away AI:** Human writers choose their metaphors. They might use one
strong image per paragraph and let it breathe. AI generates figurative language the
way it generates everything else -- by statistical pattern completion -- so it has no
sense of when enough is enough. The result, as one critic put it, reads like prose
that is "as purple as an eggplant stomping grapes."

**AI example:**
> The void was a canvas of infinite darkness, painted with the brushstrokes of distant
> stars that whispered ancient secrets. Each nebula was a cathedral of light, its
> pillars of gas reaching upward like prayers from a dying civilization, while the
> cosmic winds sang a requiem for worlds that had been swallowed by the hungry throat
> of entropy.

**Fix:** One metaphor per paragraph maximum. If you've used a figurative comparison,
the next 2-3 sentences should be concrete and literal. Never stack two metaphors
in the same sentence. Delete any metaphor that doesn't earn its place by revealing
something specific about the character or situation.

---

### 6. Generic Sensory Details

**What it is:** AI defaults to the most common sensory descriptions from its
training data. "Gentle breeze," "blooming flowers," "golden light," "crisp air,"
"towering spires." CRAFT Literary's analysis found that when asked to describe a
pond, AI writes "A peaceful day at a pond is a serene and tranquil experience" --
pure summary with zero specific observation.

**Why it gives away AI:** Human writers select details that are specific to their
character's perception and the story's emotional context. Jesmyn Ward writes about
"washing coal from her father's hair." A human writer describing a dying dog writes
about it "trying to flip over into a headstand." These details could only come from
one particular observer in one particular moment. AI details could appear in any
book, any scene, any genre.

**AI example:**
> The space station's observation deck offered a breathtaking view of the stars.
> The gentle glow of distant galaxies painted the viewport in hues of blue and
> purple. It was a sight that filled the soul with wonder and reminded one of
> humanity's small place in the vast universe.

**Human prose does this instead:** Picks one or two weird, specific, uncomfortable,
or unexpected details. The observation deck smells like burnt coffee and recycled
air. The viewport has a hairline crack that nobody has requisitioned parts to fix.
The stars don't inspire wonder -- they make the character think about her dead
mother, who used to name constellations wrong on purpose to make her laugh.

**Fix:** For every description, ask: "Could this appear in any other book?" If yes,
replace it. Use details that could only exist in this scene, seen by this character,
in this emotional state.

---

### 7. All Characters Sound the Same

**What it is:** Every character uses identical sentence structures, vocabulary level,
and speech patterns regardless of background, education, age, or personality. One
detailed analysis found that "all characters use the same sentence structures and
vocabulary regardless of background, status, or personality."

**Why it gives away AI:** AI generates from a single statistical model. It has one
"voice." Without heavy prompting, a street kid and a professor will use the same
diction. A soldier and a poet will construct sentences identically. Readers detect
this quickly in dialogue-heavy scenes because the conversations feel like one person
talking to themselves.

**AI example:**
> "We need to consider the implications of what we've discovered," said Marcus.
> "Indeed, the ramifications could be significant for the entire crew," replied Jara.
> "I understand the gravity of the situation, but we must proceed with caution," added
> Kez, the mechanic.

*(A mechanic does not say "ramifications" and "gravity of the situation.")*

**Human prose does this instead:** Characters have distinct vocabularies, sentence
lengths, verbal tics, and speech rhythms. A mechanic speaks in short declaratives
and technical jargon. A diplomat hedges and qualifies. A teenager interrupts and
trails off. Real dialogue captures how people actually talk, including
mispronunciations (Ward writes "orner" for "ornery"), incomplete sentences, and
non-sequiturs.

**Fix:** Write a voice sheet for each major character. Define their vocabulary range,
typical sentence length, verbal habits, and what they would never say. Read dialogue
aloud: could you tell which character is speaking without the dialogue tag?

---

### 8. Negative Parallelism: "It's Not X, It's Y"

**What it is:** One of the most frequently identified single sentence constructions
in AI writing. Variants include:
- "It's not X -- it's Y"
- "Not because X, but because Y"
- "X -- not Y"
- "Not X. Not Y. Just Z."

The AI Tropes gist (ossa-ma) calls the "It's not X -- it's Y" pattern "the single
most commonly identified AI writing tell."

**Why it gives away AI:** This construction creates a false sense of insight by
negating one framing and substituting another. Humans use it occasionally for
rhetorical effect. AI uses it constantly because training data rewarded it as
"insightful" language. When it appears more than once per chapter, it signals
machine generation.

**AI example:**
> The void wasn't empty. It was alive.
> Marcus didn't feel fear. He felt something older, something primal.
> This wasn't a mission anymore. It was a reckoning.

**Fix:** Search for "not...but" and "wasn't...was" patterns. Allow yourself one per
chapter at most. Replace others with direct statements, questions, or action.

---

### 9. Rule of Three Abuse

**What it is:** AI compulsively groups items in threes: three adjectives, three
actions, three examples. "Efficient, effective, and reliable." "The darkness,
the silence, and the cold." While the rule of three is a legitimate rhetorical
device, AI deploys it in nearly every description.

**Why it gives away AI:** The frequency is the tell. One source notes AI uses the
rule of three "every other sentence," making the rhythm predictable and mechanical.
Human writers use threes selectively for emphasis. AI uses them as a default
structural template.

**AI example:**
> The station was vast, ancient, and dying. Its corridors held secrets, shadows, and
> the echoes of a thousand voices. Marcus felt small, insignificant, and profoundly
> alone.

**Fix:** Vary your list lengths. Use pairs. Use singles. Occasionally use four or
five items. When you do use three, make sure the third item breaks the pattern
established by the first two (the classic comic triple).

---

### 10. Relentless Positivity / Redemption Bias

**What it is:** AI defaults to upbeat resolutions, redemption arcs, and happy
endings even when the story doesn't call for them. The Poets & Writers analysis
found that when writing about animals with parvovirus, "the puppies always survived,
always triumphed." AI "resisted darker possibilities."

**Why it gives away AI:** LLMs are trained with safety constraints and "helpful"
reward signals that bias them toward positive outcomes. They struggle to sustain
tragic, ambiguous, or genuinely dark conclusions. Fiction that resolves too cleanly
without earning it feels artificial to experienced readers, especially in genres like
horror, literary fiction, noir, and grimdark fantasy.

**AI example:**
> Despite everything they had endured, Marcus and the crew found renewed purpose.
> The station, once a symbol of humanity's failures, became a beacon of hope. They
> had lost much, but they had gained something far more valuable: each other.

**Fix:** Let bad things stay bad. Not every wound heals. Not every sacrifice pays
off. The most memorable fiction often ends with unanswered questions, unresolved
tension, and characters who are changed but not fixed.

---

### 11. Cliche Action Beats

**What it is:** AI cycles through the same physical reactions to convey emotion:
heart pounding, jaw clenching, eyes widening, breath catching, fists tightening,
stomach dropping, knuckles whitening. One analysis found "her heart pounding in
her chest" recurring constantly. Another found phrases with the word "walls"
(as in emotional walls) appearing once every seven pages.

**Fiction-specific cliches AI loves:**
- "released a breath he/she didn't know he/she was holding"
- "a shiver ran down her spine"
- "electricity crackled between them"
- "couldn't help but notice"
- "the weight of [something abstract]"
- "her walls crumbled"
- "his blood ran cold"
- "something shifted in her chest"
- "a beat of silence" (used as a paragraph transition)

**Why it gives away AI:** These phrases are everywhere in the training data because
they were already overused in published fiction. AI amplifies the problem by
selecting the most statistically common phrases. When a manuscript uses five or
more of these in a single chapter, it flags as AI-generated.

**Fix:** Ban these phrases outright. Instead, find a physical detail specific to
your character. A nervous engineer doesn't "clench her jaw" -- she taps the edge of
her datapad with her thumbnail or checks the seal on a compartment that's already
sealed. Anchor reactions in character-specific behavior.

---

### 12. Consistency Failures

**What it is:** AI loses track of established facts across long texts. Character
descriptions change (green eyes become brown). Dead characters reappear. Timeline
contradictions surface. Magic systems that worked one way in chapter 5 work
differently in chapter 22. Technology appears in settings where it shouldn't exist.

**Why it gives away AI:** Context window limitations mean AI literally cannot hold
an entire novel in memory simultaneously. One analysis documented identical scenes
appearing verbatim across multiple chapters, advisor introductions repeated word
for word, and plot revelations occurring multiple times "as though previous
revelations didn't occur."

**Reader impact:** Readers in genre fiction (especially fantasy and sci-fi) treat
consistency failures as serious quality defects. When a reader already suspects
AI generation, a single continuity error becomes "evidence of careless AI
generation" rather than a forgivable mistake.

**Fix:** Maintain a continuity ledger (you already have one). Track every character
description, timeline marker, and worldbuilding rule. Review each chapter against
the ledger before finalizing.

---

### 13. Missing Subtext

**What it is:** Characters say exactly what they mean. There is no gap between what
they say and what they feel. No lies, no evasions, no loaded silences. AI "cannot
generate meaningful absence" (CRAFT Literary). One analysis describes AI dialogue as
"chatbot having feelings" -- characters who are "overly articulate about emotions,
filling silences with explanation instead of implication."

**Why it gives away AI:** Real human conversation is built on what goes unsaid. A
married couple discussing dinner plans might really be negotiating whether their
marriage is over. A soldier reporting to a commander might be hiding fear beneath
clinical language. AI cannot write this naturally because it is "trained to generate
phrase patterns; it's not trained to generate silence."

**AI example:**
> "I'm scared," she said. "I don't know if I'm ready for this."
> "I understand your fear," he replied. "But I believe in your abilities."
> "Thank you. That means a lot to me."

**Human prose does this instead:**
> "The nav readings look clean," she said, not looking at him.
> He waited.
> "Just tell me when," she said.

**Fix:** In every dialogue exchange, identify what the characters are *not* saying.
Write around the unsaid thing. Let the reader infer what the characters won't admit.

---

### 14. Perfect Grammar / Excessive Formality

**What it is:** AI produces grammatically flawless prose. It never starts sentences
with "And" or "But." It rarely uses contractions. It never writes fragments
intentionally. It uses Oxford commas consistently. It defaults to formal register
even in casual scenes.

**Why it gives away AI:** Human prose, especially fiction, breaks rules constantly
and deliberately. Fragments create emphasis. Starting with "And" creates connection.
Contractions make narration feel natural. Run-on sentences mirror racing thoughts.
Grammatical perfection reads as sterile and processed.

**Fix:** Write in the register your POV character would think in. A teenager's
internal monologue should not read like a thesis. Break grammar rules where
breaking them serves voice and rhythm. Use contractions freely. Write fragments.
Start sentences with conjunctions.

---

### 15. Formulaic Paragraph Structure

**What it is:** AI paragraphs tend toward uniform length (4-6 sentences) with
identical opening patterns. Every paragraph in a description might start with a
noun. Every paragraph in an action sequence might start with a character name. The
"listicle in a trench coat" -- numbered or labeled points dressed as prose with
"The first... The second... The third..."

**Why it gives away AI:** Human paragraph length varies dramatically. A paragraph
might be one word. The next might be twelve sentences. Opening patterns shift
naturally. AI produces what one source calls "fractal summaries" -- summary at every
level of the text, creating a numbing regularity.

**Fix:** Vary paragraph length from 1 sentence to 10+. Never start three consecutive
paragraphs the same way. Use one-sentence paragraphs for impact. Use long paragraphs
for immersion.

---

## Additional Fiction-Specific Tells

### Tone Drift

AI's "helpful assistant" personality bleeds through. Chapter 3 might be atmospheric
and tense. Chapter 4 suddenly becomes upbeat and encouraging. The tone drifts because
the model cannot maintain emotional consistency across long texts without constant
course correction.

### Self-Contained Scenes

Every scene wraps up neatly. AI generates self-complete units when the story needs
fragments -- scenes that leave threads hanging for later chapters. Wales calls this
"self-containment syndrome."

### Missing Idiosyncratic Detail

AI cannot generate the unexpected, specific observations that mark genuine
authorship. A human writer notices things only that particular person would notice.
AI notices what every generic observer would notice.

### Resolved Ambiguity

AI wants to answer every question it raises. It cannot leave mysteries unresolved
or let contradictions stand. Stories lose tension because the AI resolves uncertainty
as fast as it creates it.

### Absent Silence

Scenes are filled wall-to-wall with action, dialogue, and description. There are no
pauses. No moments where a character simply exists in space without purpose. No dead
air in conversation. AI has no model for productive emptiness.

### Historical/Cultural Flattening

AI cannot synthesize multiple cultural or literary traditions the way human writers
can. Ward weaving Greek tragedy into contemporary Mississippi fiction, Pollack
connecting obsessive parenting to specific Jewish-American experience -- these
conceptual leaps require embodied understanding AI lacks.

---

## Measurable Linguistic Features (From Academic Research)

The following are measurable with text analysis tools, based on the StyloAI (2024)
and related stylometric studies:

| Feature | Human Writing | AI Writing |
|---------|--------------|------------|
| Type-Token Ratio (vocabulary diversity) | Balanced, context-appropriate | Either too narrow or artificially broad |
| Hapax Legomena Rate (single-use words) | Higher -- rich, varied vocabulary | Lower -- recycles vocabulary more |
| Sentence length standard deviation | High (varies a lot) | Low (uniform) |
| Stop word frequency | Natural, maintains flow | Deviates from natural speech rhythm |
| Unique bigram/trigram count | Higher -- more original word pairings | Lower -- relies on common phrases |
| Contraction frequency | Moderate to high in fiction | Very low |
| Pronoun-to-noun ratio | Higher (more pronouns, reference, perspective) | Lower (nominal loading, denser clauses) |
| Named entity usage | Specific, anchored to real knowledge | Sometimes fabricated or misaligned |
| Emotional word distribution | Varies naturally with narrative arc | Over-concentrates or under-uses |
| Perplexity variation (burstiness) | High -- surprising word choices mixed in | Low -- consistently predictable |

---

## The "Copula Avoidance" Pattern

A subtle but reliable tell identified by Wikipedia editors: AI avoids simple "is/has"
constructions in favor of:
- "serves as" instead of "is"
- "features" instead of "has"
- "marks" or "represents" instead of direct statement

In fiction, this manifests as:
> "The station served as a reminder of everything they had lost"
instead of:
> "The station was a reminder."

Or:
> "Her expression represented a complex emotional state"
instead of:
> "She looked sick."

---

## The "Alternating Modes" Pattern

Identified by the Record Crash analysis: AI prose cycles mechanically between modes.
A paragraph of "loquacious ponderous similes" followed by "tiny sentences made out
of cliches." Dense description block, then staccato dialogue, then dense description
again. The transitions between modes feel artificial because they follow a pattern
rather than responding to narrative need.

Human writers shift modes fluidly, often within a single paragraph. Description
bleeds into thought which triggers dialogue which returns to action. AI builds
blocks.

---

## Sources

- Neil Clarke / Clarkesworld Magazine editor accounts (2023-2025)
- Alexander Wales, "Adventures in AI Text Generation" (detailed failure mode taxonomy)
- Coyote Tracks, "Creative Writing and AI's Failure Modes" (extensive novel-length testing)
- Record Crash / Makin, "How to Identify AI-Written Web Fiction" (pattern analysis)
- CRAFT Literary, "Show, Don't Tell: What AI Can't Do" (2025)
- Poets & Writers / Eileen Pollack, "The Antithesis of Inspiration" (literary analysis)
- Max Read / Substack, OpenAI creative writing model review (2025)
- TechCrunch, OpenAI creative writing model criticism (2025)
- StyloAI academic paper (2024) -- stylometric feature analysis
- Pangram Labs, comprehensive AI pattern guide
- Wikipedia:Signs of AI writing (community-maintained reference)
- ossa-ma AI Writing Tropes gist (comprehensive trope taxonomy)
- Cornell Daily Sun, "Mass-Generating Literary Slop" (AI novel review, 2026)
- Novarrium, AI novel consistency analysis
- Hunting the Muse, "6 Elements of a Robot's Style"
- Academic Platypus / Michelle Kassorla, "Recognizing AI Structures in Writing"
- Multiple Reddit communities (r/writing, r/selfpublish, r/WritingWithAI)
- Jane Friedman publishing industry analysis
- Hacker News thread discussions on AI writing tells
