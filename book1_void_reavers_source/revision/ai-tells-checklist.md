# AI Tells Checklist

Run this against each chapter during revision. A single "yes" is fine.
Three or more in one chapter means that chapter needs a rewrite pass.

---

## Sentence-Level

- [ ] **Uniform sentence length?** Read a paragraph aloud. Do all sentences take
      roughly the same time to say? (Should vary from 3 words to 30+.)
- [ ] **No fragments or intentional rule-breaking?** Is every sentence grammatically
      perfect? (Fiction needs fragments, conjunctions as openers, run-ons.)
- [ ] **Low contraction rate?** Are characters and narration using "cannot," "do not,"
      "it is" where "can't," "don't," "it's" would sound natural?
- [ ] **Em dash overuse?** More than 2-3 em dashes per page? (Count them.)
- [ ] **"Not X, it's Y" pattern?** Does this negative parallelism appear more than
      once per chapter?
- [ ] **Rule of three on repeat?** Are adjectives/nouns grouped in threes more than
      twice per page?
- [ ] **Semicolons connecting simple clauses?** AI connects phrases with semicolons
      where conjunctions or periods would be more natural.

## Word Choice

- [ ] **AI vocabulary present?** Search for: tapestry, testament, beacon, vibrant,
      pivotal, profound, meticulous, seamless, nuanced, multifaceted, unwavering,
      delve, embark, unveil, navigate, foster, illuminate, resonate, transcend.
      Any in the chapter?
- [ ] **"Serves as" / "stands as" / "marks" / "represents"?** Instead of just
      using "is" or "was"?
- [ ] **Adverb clustering?** Seamlessly, meticulously, profoundly, relentlessly,
      notably appearing near each other?
- [ ] **Generic sensory language?** "Gentle breeze," "golden light," "towering
      spires," "crisp air" -- details that could appear in any book?

## Emotion and Character

- [ ] **Named emotions?** "She felt grief/anger/fear/relief" stated directly instead
      of shown through action and detail?
- [ ] **Cliche action beats?** Heart pounding, jaw clenching, eyes widening, breath
      catching, fists tightening, stomach dropping, knuckles whitening?
- [ ] **"Released a breath she didn't know she was holding"?** Or any variant of this
      specific cliche?
- [ ] **"The weight of [abstract noun]"?** Used to convey emotional burden?
- [ ] **"Something shifted in [body part]"?** Vague internal sensation?
- [ ] **Characters all sound the same in dialogue?** Cover the dialogue tags -- can
      you tell who is speaking from voice alone?
- [ ] **Characters say exactly what they mean?** No subtext, no lies, no evasion,
      no loaded silence?
- [ ] **Emotions explicitly labeled in dialogue?** "I'm scared." "I understand your
      fear." "That means a lot to me."

## Structure and Pacing

- [ ] **Uniform paragraph length?** Are most paragraphs 4-6 sentences with no
      variation? (Should range from 1 sentence to 10+.)
- [ ] **Same paragraph openings?** Do 3+ consecutive paragraphs start the same way?
- [ ] **Every scene wraps up neatly?** No loose threads, no unresolved tension
      carried forward?
- [ ] **Tone drift between chapters?** Does the mood shift without narrative cause?
      (Dark chapter followed by inexplicably upbeat chapter.)
- [ ] **Resolves ambiguity too fast?** Are mysteries or questions answered within the
      same scene they're raised?
- [ ] **No silence?** Is every moment filled with action/dialogue/description? Are
      there any pauses where a character simply exists?

## Purple Prose and Metaphor

- [ ] **Stacked metaphors?** More than one figurative comparison in the same
      sentence or consecutive sentences?
- [ ] **Metaphor-to-wrong-target?** Is the imagery evocative but attached to
      something that doesn't logically fit? (e.g., "constraints humming")
- [ ] **Abstract descriptions of concrete things?** "A symphony of colors," "a
      dance of light," "a tapestry of experiences" instead of specific imagery?
- [ ] **Everything is beautiful?** Relentless positivity in description? No
      ugliness, discomfort, or mundane detail?

## Consistency (Full-Manuscript)

- [ ] **Character descriptions stable?** Eye color, hair, height, distinguishing
      features consistent across chapters?
- [ ] **Dead characters stay dead?** No one reappears without explanation?
- [ ] **Timeline coherent?** "Three days ago" references match actual chapter
      chronology?
- [ ] **Worldbuilding rules consistent?** Tech, magic, physics, geography match
      earlier established facts?
- [ ] **Repeated scenes?** Any paragraphs or revelations that appear verbatim or
      near-verbatim in multiple chapters?

## Overall Chapter Test

- [ ] **Could any other author have written this chapter?** If yes, the voice is
      too generic.
- [ ] **Could these descriptions appear in any other book?** If yes, the details
      are too generic.
- [ ] **Does the chapter end too cleanly?** If it wraps up with a bow, it needs a
      thread left hanging.
- [ ] **Read one page to someone unfamiliar with AI writing. Do they say it sounds
      "weird" or "off" or "too smooth"?** Trust non-writer instincts.

---

## Quick Grep Commands

Run these against your manuscript to get fast counts:

```bash
# AI vocabulary scan
grep -icE 'tapestry|testament|beacon|vibrant|pivotal|profound|meticulous|seamless|nuanced|multifaceted|unwavering|delve|embark|unveil|navigate|foster|illuminate|resonate|transcend' chapter.md

# Em dash count
grep -o '---\|--\|—' chapter.md | wc -l

# "Not X, it's Y" pattern
grep -icE "(it'?s not|wasn'?t|isn'?t).{1,40}(it'?s|it was|but)" chapter.md

# Cliche action beats
grep -icE 'heart pound|jaw clench|eyes widen|breath caught|fist tight|stomach drop|knuckles whit|blood ran cold|shiver.{0,10}spine' chapter.md

# "Weight of" cliche
grep -ic 'weight of' chapter.md

# "Serves as" copula avoidance
grep -icE 'serves as|stands as|marks a|represents a' chapter.md

# Telling emotions
grep -icE 'felt (a |the )?(surge|wave|pang|rush|mixture|complex)' chapter.md
```

---

## Severity Levels

**Cosmetic (fix during line edit):**
Single instances of AI vocabulary, occasional cliche beat, one too many em dashes.

**Structural (fix during revision pass):**
Uniform sentence rhythm across multiple paragraphs, all characters sounding the same
in a dialogue scene, missing subtext in a key emotional exchange.

**Rewrite required:**
Exposed subtext throughout a scene (emotions named instead of dramatized), purple
prose stacking 3+ metaphors per paragraph, consistency failures with established
facts.
